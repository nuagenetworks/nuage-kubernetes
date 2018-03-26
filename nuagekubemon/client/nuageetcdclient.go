/*
###########################################################################
#
#   Filename:           nuageetcdclient.go
#
#   Author:             Siva Teja, Areti
#   Created:            August 2, 2017
#
#   Description:        NuageKubeMon etcd Client Interface
#
###########################################################################
#
#              Copyright (c) 2017 Nuage Networks
#
###########################################################################
*/

package client

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/golang/glog"
	"github.com/nuagenetworks/nuage-kubernetes/nuagekubemon/api"
	"github.com/nuagenetworks/nuage-kubernetes/nuagekubemon/config"
	"io/ioutil"
	"path"
	"strconv"
	"strings"
	"time"
)

const (
	ETCD_BASE_PATH       = "/nuage.io/v1.0/"
	POD_METADATA_TREE    = ETCD_BASE_PATH + "pod_metadata/"
	POOL_CIDR_TREE       = ETCD_BASE_PATH + "pool_cidr/"
	SUBNET_METADATA_TREE = ETCD_BASE_PATH + "subnet_metadata/"
	ZONE_METADATA_TREE   = ETCD_BASE_PATH + "zone_metadata/"
	SCALE_UP_THRESHOLD   = 75
	SCALE_DOWN_THRESHOLD = 25
)

type etcdSubnetValue struct {
	ACTIVEIP int    `json:"activeip"`
	VSDID    string `json:"vsdid"`
	CIDR     string `json:"cidr"`
}

type NuageEtcdClient struct {
	etcdBaseURL       []string
	serverCA          string
	clientCertificate string
	clientKey         string
	maxIPCount        int
	subnetSize        int
	subnetIDCache     map[string]string
	client            *clientv3.Client
	clusterNetwork    *IPv4Subnet
}

//NewNuageEtcdClient creates a new etcd client
func NewNuageEtcdClient(nkmConfig *config.NuageKubeMonConfig) (*NuageEtcdClient, error) {
	nuageetcd := &NuageEtcdClient{
		etcdBaseURL:       nkmConfig.EtcdClientConfig.UrlList,
		serverCA:          nkmConfig.EtcdClientConfig.ServerCA,
		clientCertificate: nkmConfig.EtcdClientConfig.ClientCertificate,
		clientKey:         nkmConfig.EtcdClientConfig.ClientKey,
	}
	nuageetcd.subnetIDCache = make(map[string]string)
	return nuageetcd, nuageetcd.Init(nkmConfig)
}

//Init Initializes the etcd client
func (nuageetcd *NuageEtcdClient) Init(nkmConfig *config.NuageKubeMonConfig) error {
	var err error
	if len(nuageetcd.etcdBaseURL) == 0 {
		nuageetcd.etcdBaseURL = []string{"http://127.0.0.1:2379"}
	}
	nuageetcd.clusterNetwork, err = IPv4SubnetFromString(nkmConfig.MasterConfig.NetworkConfig.ClusterCIDR)
	if err != nil {
		return fmt.Errorf("Failure in getting cluster CIDR: %s\n", err)
	}
	nuageetcd.subnetSize = nkmConfig.MasterConfig.NetworkConfig.SubnetLength
	if nuageetcd.subnetSize < 2 || nuageetcd.subnetSize > 32 {
		glog.Errorf("Invalid hostSubnetLength of %d.  Using default value of 8",
			nuageetcd.subnetSize)
		nuageetcd.subnetSize = 8
	}
	if nuageetcd.subnetSize > (32 - nuageetcd.clusterNetwork.CIDRMask) {
		return fmt.Errorf("Cannot allocate %d bit subnets from %s. ",
			nuageetcd.subnetSize, nuageetcd.clusterNetwork.String())
	}
	if nuageetcd.maxIPCount == 0 {
		nuageetcd.maxIPCount = 1<<uint(nuageetcd.subnetSize) - 3
	}
	etcdConfig := clientv3.Config{
		Endpoints:   nuageetcd.etcdBaseURL,
		DialTimeout: 5 * time.Second,
	}
	//TODO TLS setup
	if nuageetcd.serverCA != "" {
		etcdConfig.TLS, err = nuageetcd.tls_setup()
		if err != nil {
			glog.Errorf("Error doing tls setup: %v", err)
			return err
		}
	}
	nuageetcd.client, err = clientv3.New(etcdConfig)
	if err != nil {
		glog.Errorf("Creating etcd client failed with error %v", err)
		return err
	}
	return nil
}

//AllocateSubnetForPod increments active ip count of a subnet in a namespace
func (nuageetcd *NuageEtcdClient) AllocateSubnetForPod(data *api.EtcdPodMetadata) (*api.EtcdPodSubnet, error) {
	pod := data.PodName
	ns := data.NamespaceName
	podSubnet := &api.EtcdPodSubnet{}
	var err error
	var subnetStr string
	var txnResp *clientv3.TxnResponse

	f := func() (bool, error) {
		podSubnet.ToCreate = ""
		podSubnet.ToUse = ""
		count := nuageetcd.maxIPCount
		ACTIVEIPCount := 0
		suffix := -1
		snet := &etcdSubnetValue{}
		allocatedSubnet := &etcdSubnetValue{}

		nuageEtcdRetry(
			func() error {
				txnResp, err = nuageetcd.client.KV.Txn(context.Background()).Then(
					clientv3.OpGet(POD_METADATA_TREE+ns+"/"+pod),
					clientv3.OpGet(SUBNET_METADATA_TREE+ns+"/"+ns, clientv3.WithPrefix())).Commit()
				return err
			})
		if err != nil {
			glog.Errorf("fetching pod and subnet metadata during allocating subnet failed: %v", err)
			return false, err
		}
		podResp := (*clientv3.GetResponse)(txnResp.Responses[0].GetResponseRange())
		subnetResp := (*clientv3.GetResponse)(txnResp.Responses[1].GetResponseRange())
		if len(podResp.Kvs) != 0 {
			podSubnet.ToUse = string(podResp.Kvs[0].Value)
			glog.Warningf("pod %s is already allocated to subnet %s", pod, podSubnet.ToUse)
			return true, nil
		}

		for _, kv := range subnetResp.Kvs {
			if err := json.Unmarshal(kv.Value, snet); err != nil {
				glog.Errorf("unmarshal kv pair(%s, %s) conversion failed: %v", kv.Key, kv.Value, err)
				continue
			} else {
				ACTIVEIPCount += snet.ACTIVEIP
				suffix = max(suffix, getSuffix(string(kv.Key)))
				if count == nuageetcd.maxIPCount && snet.ACTIVEIP < count {
					count = snet.ACTIVEIP
					subnetStr = string(kv.Key)
					allocatedSubnet = &etcdSubnetValue{
						CIDR:     snet.CIDR,
						VSDID:    snet.VSDID,
						ACTIVEIP: snet.ACTIVEIP,
					}
				}
			}
		}

		puts := []clientv3.Op{}
		compares := []clientv3.Cmp{}
		noOfSubnets := len(subnetResp.Kvs)
		if (ACTIVEIPCount+1)*100 > noOfSubnets*nuageetcd.maxIPCount*SCALE_UP_THRESHOLD {
			newSubnet := fmt.Sprintf("%s-%d", ns, suffix+1)
			snet := &etcdSubnetValue{ACTIVEIP: 0, VSDID: "", CIDR: "0"}
			b, err := json.Marshal(snet)
			if err != nil {
				glog.Errorf("marshal struct %v to json string failed: %v", snet, err)
			}
			puts = append(puts, clientv3.OpPut(SUBNET_METADATA_TREE+ns+"/"+newSubnet, string(b)))
			podSubnet.ToCreate = newSubnet
		}
		podSubnet.ToUse = path.Base(subnetStr)
		allocatedSubnet.ACTIVEIP += 1
		b, err := json.Marshal(allocatedSubnet)
		if err != nil {
			glog.Errorf("marshal struct %v to json string failed: %v", allocatedSubnet, err)
		}

		puts = append(puts, clientv3.OpPut(subnetStr, string(b)))
		puts = append(puts, clientv3.OpPut(POD_METADATA_TREE+ns+"/"+pod, path.Base(subnetStr)))
		for _, kv := range subnetResp.Kvs {
			compares = append(compares, clientv3.Compare(clientv3.ModRevision(string(kv.Key)), "=", kv.ModRevision))
		}
		nuageEtcdRetry(
			func() error {
				txnResp, err = nuageetcd.client.KV.Txn(context.Background()).If(
					compares...).Then(
					puts...).Commit()
				return err
			})
		if err != nil {
			glog.Errorf("writing transaction to etcd failed:", err)
			return false, err
		}
		return txnResp.Succeeded, nil
	}
	nuageetcd.nuageSTM(f)
	return podSubnet, nil
}

//DeAllocateSubnetFromPod decrements active ip count of a subnet in a namespace
func (nuageetcd *NuageEtcdClient) DeAllocateSubnetFromPod(data *api.EtcdPodMetadata) ([]*api.EtcdSubnetMetadata, error) {
	var err error
	var txnResp *clientv3.TxnResponse
	pod := data.PodName
	ns := data.NamespaceName
	var subnetList []*api.EtcdSubnetMetadata

	f := func() (bool, error) {
		ACTIVEIPCount := 0
		noOfSubnets := 0
		compares := []clientv3.Cmp{}
		ops := []clientv3.Op{}
		nuageEtcdRetry(
			func() error {
				txnResp, err = nuageetcd.client.KV.Txn(context.Background()).Then(
					clientv3.OpGet(POD_METADATA_TREE+ns+"/"+pod),
					clientv3.OpGet(SUBNET_METADATA_TREE+ns+"/"+ns, clientv3.WithPrefix())).Commit()
				return err
			})
		if err != nil {
			glog.Errorf("fetching pod and subnet metadata during deallocating subnet failed: %v", err)
			return false, err
		}
		podResp := (*clientv3.GetResponse)(txnResp.Responses[0].GetResponseRange())
		subnetResp := (*clientv3.GetResponse)(txnResp.Responses[1].GetResponseRange())
		if len(podResp.Kvs) == 0 {
			return true, fmt.Errorf("no pod found with given name(%s) inside etcd", pod)
		}
		subnetStr := string(podResp.Kvs[0].Value)
		snet := &etcdSubnetValue{}
		for _, kv := range subnetResp.Kvs {
			if path.Base(string(kv.Key)) == subnetStr {
				if err := json.Unmarshal(kv.Value, snet); err != nil {
					glog.Errorf("unmarshal kv pair(%s, %s) to struct failed: %v", kv.Key, kv.Value, err)
					return false, err
				}
				break
			}
		}
		snet.ACTIVEIP -= 1
		if b, err := json.Marshal(snet); err != nil {
			glog.Errorf("marshal struct %v to string failed: %v", snet, err)
			return false, err
		} else {
			ops = append(ops, clientv3.OpPut(SUBNET_METADATA_TREE+ns+"/"+subnetStr, string(b)))
		}
		for _, kv := range subnetResp.Kvs {
			if err := json.Unmarshal(kv.Value, snet); err != nil {
				glog.Errorf("unmarshal kv pair(%s, %s) to struct failed: %v", kv.Key, kv.Value, err)
			} else {
				ACTIVEIPCount += snet.ACTIVEIP
				noOfSubnets += 1
			}
		}
		if noOfSubnets != 1 && ACTIVEIPCount*100 < noOfSubnets*nuageetcd.maxIPCount*SCALE_DOWN_THRESHOLD {
			for _, kv := range subnetResp.Kvs {
				if err := json.Unmarshal(kv.Value, snet); err != nil {
					glog.Errorf("unmarshal kv pair(%s, %s) to struct failed: %v", kv.Key, kv.Value, err)
				} else {
					if snet.ACTIVEIP == 0 {
						ops = append(ops, clientv3.OpDelete(string(kv.Key)))
						subnetList = append(subnetList, &api.EtcdSubnetMetadata{ID: snet.VSDID, CIDR: snet.CIDR})
					}
				}
			}
		}
		ops = append(ops, clientv3.OpDelete(POD_METADATA_TREE+ns+"/"+pod))
		for _, kv := range subnetResp.Kvs {
			compares = append(compares, clientv3.Compare(clientv3.ModRevision(string(kv.Key)), "=", kv.ModRevision))
		}
		nuageEtcdRetry(
			func() error {
				txnResp, err = nuageetcd.client.KV.Txn(context.Background()).If(
					compares...).Then(
					ops...).Commit()
				return err
			})
		if err != nil {
			glog.Errorf("writing transaction to etcd failed:", err)
			return false, err
		}
		return txnResp.Succeeded, nil
	}
	nuageetcd.nuageSTM(f)
	return subnetList, nil
}

//CreateFirstSubnet creates the first subnet in a namespace
func (nuageetcd *NuageEtcdClient) CreateFirstSubnet(subnetInfo *api.EtcdSubnetMetadata) error {
	var err error
	var txnResp *clientv3.TxnResponse
	b := []byte{}
	ns := subnetInfo.Namespace
	key := SUBNET_METADATA_TREE + ns + "/" + subnetInfo.Name
	s := &etcdSubnetValue{ACTIVEIP: 0, VSDID: subnetInfo.ID, CIDR: subnetInfo.CIDR}
	b, err = json.Marshal(s)
	if err != nil {
		glog.Errorf("json marshal struct(%v) failed: %v", s, err)
		return err
	}
	nuageEtcdRetry(
		func() error {
			txnResp, err = nuageetcd.client.KV.Txn(context.Background()).If(
				clientv3.Compare(clientv3.ModRevision(key), "=", 0),
			).Then(
				clientv3.OpPut(key, string(b)),
			).Commit()
			return err
		})
	if err != nil {
		glog.Errorf("creating first subnet in namespace %s failed %v", ns, err)
		return err
	}
	if !txnResp.Succeeded {
		glog.Warningf("first subnet in namespace %s is already created", ns)
	}
	return nil
}

//DeleteLastSubnet deletes last subnet in a namespace
func (nuageetcd *NuageEtcdClient) DeleteLastSubnet(subnetInfo *api.EtcdSubnetMetadata) (*api.EtcdSubnetMetadata, error) {
	var err error
	var delResp *clientv3.DeleteResponse
	key := SUBNET_METADATA_TREE + subnetInfo.Namespace + "/" + subnetInfo.Name
	nuageEtcdRetry(
		func() error {
			delResp, err = nuageetcd.client.Delete(context.Background(), key, clientv3.WithPrefix(), clientv3.WithPrevKV())
			return err
		})
	if err != nil {
		glog.Errorf("delete on %s failed: %v", key, err)
		return nil, err
	}
	if len(delResp.PrevKvs) == 0 {
		return nil, fmt.Errorf("no matching ns(%s) found", subnetInfo.Namespace)
	}
	s := &etcdSubnetValue{}
	if err := json.Unmarshal(delResp.PrevKvs[0].Value, s); err != nil {
		return nil, fmt.Errorf("unmarshal json string %s to struct failed: %v", string(delResp.PrevKvs[0].Value), err)
	}

	subnet := &api.EtcdSubnetMetadata{Name: path.Base(string(delResp.PrevKvs[0].Key)), CIDR: s.CIDR, ID: s.VSDID}
	return subnet, nil
}

//AllocateSubnet marks an entry if the subnet is not already used
func (nuageetcd *NuageEtcdClient) AllocateSubnetCIDR(subnet *api.EtcdSubnetMetadata) (string, error) {
	var err error
	var txnResp *clientv3.TxnResponse
	networkCIDR := strings.Replace(subnet.CIDR, "/", "-", -1)
	key := POOL_CIDR_TREE + networkCIDR
	nuageEtcdRetry(
		func() error {
			txnResp, err = nuageetcd.client.KV.Txn(context.Background()).If(
				clientv3.Compare(clientv3.ModRevision(key), "=", 0),
			).Then(
				clientv3.OpPut(key, subnet.Name),
			).Else(
				clientv3.OpGet(key),
			).Commit()
			return err
		})
	if err != nil {
		glog.Errorf("create subnet cidr key(%s) txn failed:%v", networkCIDR, err)
		return "", err
	}

	if !txnResp.Succeeded {
		getResp := (*clientv3.GetResponse)(txnResp.Responses[0].GetResponseRange())
		allocatedSubnet := string(getResp.Kvs[0].Value)
		glog.Warningf("subnet cidr %s is already allocated to subnet %s:", networkCIDR, allocatedSubnet)
		return allocatedSubnet, nil
	}
	return "", err
}

//FreeSubnet delete the subnet entry
func (nuageetcd *NuageEtcdClient) FreeSubnetCIDR(subnet *api.EtcdSubnetMetadata) error {
	var err error
	networkCIDR := strings.Replace(subnet.CIDR, "/", "-", -1)
	key := POOL_CIDR_TREE + networkCIDR
	nuageEtcdRetry(
		func() error {
			_, err = nuageetcd.client.Delete(context.Background(), key)
			if err != nil {
				return err
			}
			return nil
		})
	if err != nil {
		glog.Errorf("deleting subnet cidr(%s) failed with error: %v", key, err)
		return err
	}
	return nil
}

//UpdateSubnetIDKey create a new key that maps to vsd network id
func (nuageetcd *NuageEtcdClient) UpdateSubnetInfo(subnetInfo *api.EtcdSubnetMetadata) error {
	nuageetcd.subnetIDCache[subnetInfo.Namespace+"/"+subnetInfo.Name] = subnetInfo.ID
	key := SUBNET_METADATA_TREE + subnetInfo.Namespace + "/" + subnetInfo.Name
	s := etcdSubnetValue{ACTIVEIP: 0, VSDID: subnetInfo.ID, CIDR: subnetInfo.CIDR}
	b, err := json.Marshal(s)
	if err != nil {
		glog.Errorf("marshal struct(%v) to string failed: %v", s, err)
		return err
	}
	nuageEtcdRetry(
		func() error {
			_, err = nuageetcd.client.Put(context.Background(), key, string(b))
			if err != nil {
				return err
			}
			return nil
		})
	if err != nil {
		glog.Errorf("updating subnet(%s) info failed: %v", key, err)
		return err
	}
	return nil
}

//GetSubnetInfo fetches subnet info from subnet metadata tree
func (nuageetcd *NuageEtcdClient) GetSubnetInfo(subnetInfo *api.EtcdSubnetMetadata) (*api.EtcdSubnetMetadata, error) {
	var err error
	var getResp *clientv3.GetResponse
	key := SUBNET_METADATA_TREE + subnetInfo.Namespace + "/" + subnetInfo.Name
	nuageEtcdRetry(
		func() error {
			getResp, err = nuageetcd.client.Get(context.Background(), key)
			return err
		})
	if err != nil {
		glog.Errorf("fetching subnet info failed: %v", err)
		return nil, err
	}
	if len(getResp.Kvs) == 0 {
		return nil, fmt.Errorf("subnet(%s) under namespace(%s) not found in etcd", subnetInfo.Name, subnetInfo.Namespace)
	}
	subnet := &etcdSubnetValue{}
	err = json.Unmarshal(getResp.Kvs[0].Value, subnet)
	if err != nil {
		glog.Errorf("json unmarshal in GetSubnetInfo failed: %v", err)
		return nil, err
	}
	return &api.EtcdSubnetMetadata{
		Name:      subnetInfo.Name,
		Namespace: subnetInfo.Namespace,
		ID:        subnet.VSDID,
		CIDR:      subnet.CIDR,
	}, nil

}

//GetSubnetID fetches the subnet ID. If it is "", it waits until the value is updated
func (nuageetcd *NuageEtcdClient) GetSubnetID(subnetInfo *api.EtcdSubnetMetadata) (string, error) {
	if subnetID, ok := nuageetcd.subnetIDCache[subnetInfo.Namespace+"/"+subnetInfo.Name]; ok {
		return subnetID, nil
	}
	key := SUBNET_METADATA_TREE + subnetInfo.Namespace + "/" + subnetInfo.Name
	transform := func(b []byte) string {
		subnet := &etcdSubnetValue{}
		_ = json.Unmarshal(b, subnet)
		return subnet.VSDID
	}
	id, err := nuageetcd.nuageWatch(key, transform)
	if err != nil {
		glog.Errorf("fetching subnet(%s) id failed: %v", key, err)
		return "", err
	}
	nuageetcd.subnetIDCache[subnetInfo.Namespace+"/"+subnetInfo.Name] = id
	return id, err
}

func (nuageetcd *NuageEtcdClient) AddZone(zoneInfo *api.EtcdZoneMetadata) (string, error) {
	var id string
	var err error
	var txnResp *clientv3.TxnResponse
	key := ZONE_METADATA_TREE + zoneInfo.Name
	nuageEtcdRetry(
		func() error {
			txnResp, err = nuageetcd.client.Txn(context.Background()).If(
				clientv3.Compare(clientv3.ModRevision(key), "=", 0),
			).Then(
				clientv3.OpPut(key, ""),
			).Else(
				clientv3.OpGet(key),
			).Commit()
			return err
		})
	if err != nil {
		glog.Errorf("adding key for zone %s failed: %v", zoneInfo.Name, err)
		return "", err
	}
	if !txnResp.Succeeded {
		getResp := (*clientv3.GetResponse)(txnResp.Responses[0].GetResponseRange())
		id = string(getResp.Kvs[0].Value)
		if id != "" {
			return id, nil
		}
	} else {
		return "", nil
	}

	return nuageetcd.nuageWatch(key, func(b []byte) string { return string(b) })
}

func (nuageetcd *NuageEtcdClient) UpdateZone(zoneInfo *api.EtcdZoneMetadata) error {
	var err error
	key := ZONE_METADATA_TREE + zoneInfo.Name
	nuageEtcdRetry(
		func() error {
			_, err = nuageetcd.client.Put(context.Background(), key, zoneInfo.ID)
			return err
		})
	if err != nil {
		glog.Errorf("updating zone(%s) with id(%s) failed: %v", key, zoneInfo.ID, err)
		return err
	}
	return nil
}

func (nuageetcd *NuageEtcdClient) DeleteZone(zoneInfo *api.EtcdZoneMetadata) error {
	var err error
	key := ZONE_METADATA_TREE + zoneInfo.Name
	nuageEtcdRetry(
		func() error {
			_, err = nuageetcd.client.Delete(context.Background(), key)
			return err
		})
	if err != nil {
		glog.Errorf("deleting zone(%s) from etcd failed: %v", key, err)
		return err
	}
	return nil
}

func (nuageetcd *NuageEtcdClient) GetZonesSubnets() (map[string]map[string]bool, error) {
	var err error
	var getResp *clientv3.GetResponse
	result := make(map[string]map[string]bool)
	nuageEtcdRetry(
		func() error {
			getResp, err = nuageetcd.client.Get(context.Background(), SUBNET_METADATA_TREE,
				clientv3.WithPrefix(), clientv3.WithKeysOnly())
			return err
		})
	if err != nil {
		glog.Errorf("getting subnet metadata tree keys failed: %v", err)
		return nil, err
	}
	for _, kv := range getResp.Kvs {
		elems := strings.Split(string(kv.Key), "/")
		n := len(elems)
		if n <= 2 {
			glog.Errorf("invalid key obtained %s", kv.Key)
			continue
		}
		if _, ok := result[elems[n-2]]; !ok {
			result[elems[n-2]] = make(map[string]bool)
		}
		result[elems[n-2]][elems[n-1]] = true
	}
	return result, nil
}

func (nuageetcd *NuageEtcdClient) nuageSTM(f func() (bool, error)) {
	for {
		if status, err := f(); err != nil {
			return
		} else if status {
			return
		}
	}
}

func (nuageetcd *NuageEtcdClient) nuageWatch(key string, transform func([]byte) string) (string, error) {
	var watchChan clientv3.WatchChan
	var err error
	var getResp *clientv3.GetResponse
	for {
		nuageEtcdRetry(
			func() error {
				getResp, err = nuageetcd.client.Get(context.Background(), key)
				return err
			})
		if err != nil {
			glog.Errorf("fetching key %s from etcd failed: %v", key, err)
			return "", err
		}
		if len(getResp.Kvs) != 0 && transform(getResp.Kvs[0].Value) != "" {
			return transform(getResp.Kvs[0].Value), nil
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		nuageEtcdRetry(
			func() error {
				watchChan = nuageetcd.client.Watch(ctx, key, clientv3.WithRev(getResp.Header.Revision+1))
				return nil
			})
		select {
		case watchResp := <-watchChan:
			cancel()
			if watchResp.Err() != nil {
				glog.Errorf("watch received an error: %v", watchResp.Err())
				return "", watchResp.Err()
			}
			//could not figure out why this is happening. sometimes getting zero events??
			//found out that this happens when etcd server is busy.so we will keep doing this
			//until we get the right response. happens very sporadically
			if len(watchResp.Events) == 0 {
				glog.Errorf("something went wrong. received zero updates when watching for key %s. will try again in a second", key)
				time.Sleep(time.Second)
				break
			}
			return transform(watchResp.Events[0].Kv.Value), nil
		case <-ctx.Done():
			break
		}
	}
	return "", fmt.Errorf("invalid scenario")
}

//retry the same request in case of a timeout
func nuageEtcdRetry(f func() error) {
	for {
		err := f()
		if err != nil {
			glog.Errorf("%v", err)
		}
		if err != nil && strings.Contains(err.Error(), "request timed out") {
			continue
		}
		if err != nil && strings.Contains(err.Error(), "unexpected end of JSON input") {
			continue
		}
		return
	}
}

//Run starts the nuage etcd client and listens for events
func (nuageetcd *NuageEtcdClient) Run(etcdChannel chan *api.EtcdEvent) {
	for {
		select {
		case etcdEvent := <-etcdChannel:
			nuageetcd.HandleEtcdEvent(etcdEvent)
		}
	}
}

//HandleEtcdEvent handles the events that came on the channel
func (nuageetcd *NuageEtcdClient) HandleEtcdEvent(event *api.EtcdEvent) {
	var data interface{}
	var err error
	switch event.Type {
	case api.EtcdIncActiveIPCount:
		data, err = nuageetcd.AllocateSubnetForPod(event.EtcdReqObject.(*api.EtcdPodMetadata))
	case api.EtcdDecActiveIPCount:
		data, err = nuageetcd.DeAllocateSubnetFromPod(event.EtcdReqObject.(*api.EtcdPodMetadata))
	case api.EtcdAddSubnet:
		err = nuageetcd.CreateFirstSubnet(event.EtcdReqObject.(*api.EtcdSubnetMetadata))
	case api.EtcdDelSubnet:
		data, err = nuageetcd.DeleteLastSubnet(event.EtcdReqObject.(*api.EtcdSubnetMetadata))
	case api.EtcdAllocSubnetCIDR:
		data, err = nuageetcd.AllocateSubnetCIDR(event.EtcdReqObject.(*api.EtcdSubnetMetadata))
	case api.EtcdFreeSubnetCIDR:
		err = nuageetcd.FreeSubnetCIDR(event.EtcdReqObject.(*api.EtcdSubnetMetadata))
	case api.EtcdUpdateSubnetID:
		err = nuageetcd.UpdateSubnetInfo(event.EtcdReqObject.(*api.EtcdSubnetMetadata))
	case api.EtcdGetSubnetInfo:
		data, err = nuageetcd.GetSubnetInfo(event.EtcdReqObject.(*api.EtcdSubnetMetadata))
	case api.EtcdGetSubnetID:
		data, err = nuageetcd.GetSubnetID(event.EtcdReqObject.(*api.EtcdSubnetMetadata))
	case api.EtcdAddZone:
		data, err = nuageetcd.AddZone(event.EtcdReqObject.(*api.EtcdZoneMetadata))
	case api.EtcdDeleteZone:
		err = nuageetcd.DeleteZone(event.EtcdReqObject.(*api.EtcdZoneMetadata))
	case api.EtcdUpdateZone:
		err = nuageetcd.UpdateZone(event.EtcdReqObject.(*api.EtcdZoneMetadata))
	case api.EtcdGetZonesSubnets:
		data, err = nuageetcd.GetZonesSubnets()
	}
	event.EtcdRespObjectChan <- &api.EtcdRespObject{EtcdData: data, Error: err}
}

func getSuffix(s string) int {
	list := strings.Split(s, "-")
	x, err := strconv.Atoi(list[len(list)-1])
	if err != nil {
		glog.Errorf("getSuffix: string(%s) to int conversion failed", list[len(list)-1])
	}
	return x
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (nuageetcd *NuageEtcdClient) tls_setup() (*tls.Config, error) {
	// Load client cert
	cert, err := tls.LoadX509KeyPair(nuageetcd.clientCertificate, nuageetcd.clientKey)
	if err != nil {
		glog.Errorf("Error loading client cert file to communicate with Nuage K8S monitor: %v", err)
		return nil, err
	}

	// Load CA cert
	caCert, err := ioutil.ReadFile(nuageetcd.serverCA)
	if err != nil {
		glog.Errorf("Error loading CA cert file to communicate with Nuage K8S monitor: %v", err)
		return nil, err
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Setup HTTPS client
	return &tls.Config{
		Certificates:       []tls.Certificate{cert},
		RootCAs:            caCertPool,
		InsecureSkipVerify: true,
	}, nil
}
