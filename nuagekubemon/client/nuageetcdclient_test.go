package client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/nuagenetworks/nuage-kubernetes/nuagekubemon/api"
	"github.com/nuagenetworks/nuage-kubernetes/nuagekubemon/config"
	"path"
	"testing"
	"time"
)

var nuageetcd *NuageEtcdClient
var etcdChan chan *api.EtcdEvent

const MAX_COUNT = 100

func init() {
	var err error
	conf := &config.NuageKubeMonConfig{}
	conf.MasterConfig.NetworkConfig.ClusterCIDR = "70.70.0.0/16"
	conf.MasterConfig.NetworkConfig.SubnetLength = 8
	nuageetcd, err = NewNuageEtcdClient(conf)
	if err != nil {
		fmt.Printf("Starting etcd client failed with error: %v", err)
	}
	nuageetcd.maxIPCount = 10
	etcdChan = make(chan *api.EtcdEvent)
	go nuageetcd.Run(etcdChan)
}

func TestSubnetPathCreationDeletion(t *testing.T) {
	firstSubnet := &api.EtcdSubnetMetadata{
		Name:      "ns-0",
		Namespace: "ns",
		ID:        "ahakghkabgbwbgbabviwrgbgrbgi",
		CIDR:      "0.0.0.0/0",
	}

	resp := api.EtcdChanRequest(etcdChan, api.EtcdAddSubnet, firstSubnet)
	if resp.Error != nil {
		t.Fatalf("Creating subnet path failed with error: %v", resp.Error)
	}

	getResp, err := nuageetcd.client.Get(context.Background(), SUBNET_METADATA_TREE+"ns/ns-0")
	if err != nil {
		t.Fatalf("Getting key from etcd failed with error: %v", err)
	}

	if len(getResp.Kvs) != 1 {
		t.Fatalf("received unexpected number(%d) of kv pairs for key %s", len(getResp.Kvs), SUBNET_METADATA_TREE+"ns/ns-0")
	}
	s := &etcdSubnetValue{}
	if err := json.Unmarshal(getResp.Kvs[0].Value, s); err != nil {
		t.Errorf("unmarshal string %s to struct failed: %v", getResp.Kvs[0].Value, err)
	}

	if s.ACTIVEIP != 0 {
		t.Fatalf("put and get values does not match")
	}

	resp = api.EtcdChanRequest(etcdChan, api.EtcdDelSubnet, firstSubnet)
	if resp.Error != nil {
		t.Fatalf("Deleting subnet path failed with error: %v", resp.Error)
	}
}

func TestIncrementDecrement(t *testing.T) {
	var subnet string
	var getResp *clientv3.GetResponse
	var err error
	firstSubnet := &api.EtcdSubnetMetadata{
		Name:      "ns-0",
		Namespace: "ns",
		ID:        "ahakghkabgbwbgbabviwrgbgrbgi",
		CIDR:      "0.0.0.0/0",
	}

	resp := api.EtcdChanRequest(etcdChan, api.EtcdAddSubnet, firstSubnet)
	if resp.Error != nil {
		t.Fatalf("Creating subnet path failed with error: %v", resp.Error)
	}

	podData := &api.EtcdPodMetadata{NamespaceName: "ns", PodName: "pod"}
	resp = api.EtcdChanRequest(etcdChan, api.EtcdIncActiveIPCount, podData)
	if resp.Error != nil {
		t.Fatalf("incrementing active ip count failed with error: %v", resp.Error)
	}
	resp = api.EtcdChanRequest(etcdChan, api.EtcdIncActiveIPCount, podData)
	if resp.Error != nil {
		t.Fatalf("incrementing active ip count failed with error: %v", resp.Error)
	}
	podSubnet := resp.EtcdData.(*api.EtcdPodSubnet)
	subnet = podSubnet.ToUse

	getResp, err = nuageetcd.client.Get(context.Background(), SUBNET_METADATA_TREE+"ns/"+subnet)
	if err != nil {
		t.Fatalf("Getting key from etcd failed with error: %v", err)
	}

	s := &etcdSubnetValue{}
	if err := json.Unmarshal(getResp.Kvs[0].Value, s); err != nil {
		t.Errorf("unmarshal string %s to struct failed: %v", getResp.Kvs[0].Value, err)
	}

	if s.ACTIVEIP != 1 {
		t.Fatalf("increment operation did not work")
	}

	resp = api.EtcdChanRequest(etcdChan, api.EtcdDecActiveIPCount, podData)
	if resp.Error != nil {
		t.Fatalf("decrementing active ip count failed with error: %v", resp.Error)
	}
	resp = api.EtcdChanRequest(etcdChan, api.EtcdDecActiveIPCount, podData)
	if resp.Error != nil {
		t.Fatalf("decrementing active ip count failed with error: %v", resp.Error)
	}

	subnetList := resp.EtcdData.([]*api.EtcdSubnetMetadata)
	if len(subnetList) != 0 {
		t.Fatalf("received a subnet to delete but there is only one subnet")
	}

	getResp, err = nuageetcd.client.Get(context.Background(), SUBNET_METADATA_TREE+"ns/"+subnet)
	if err != nil {
		t.Fatalf("Getting key from etcd failed with error: %v", err)
	}

	if len(getResp.Kvs) == 0 {
		t.Fatalf("did not receive any kv pairs")
	}

	if err := json.Unmarshal(getResp.Kvs[0].Value, s); err != nil {
		t.Errorf("unmarshal string %s to struct failed: %v", getResp.Kvs[0].Value, err)
	}

	if s.ACTIVEIP != 0 {
		t.Fatalf("increment operation did not work")
	}

	resp = api.EtcdChanRequest(etcdChan, api.EtcdDelSubnet, firstSubnet)
	if resp.Error != nil {
		t.Fatalf("Deleting subnet path failed with error: %v", resp.Error)
	}
}

func TestAutoScaling(t *testing.T) {
	var err error
	var getResp *clientv3.GetResponse
	firstSubnet := &api.EtcdSubnetMetadata{
		Name:      "ns-0",
		Namespace: "ns",
		ID:        "ahakghkabgbwbgbabviwrgbgrbgi",
		CIDR:      "0.0.0.0/0",
	}

	if err = nuageetcd.CreateFirstSubnet(firstSubnet); err != nil {
		t.Fatalf("Creating subnet path failed with error: %v", err)
	}
	//scale up and scale down twice
	for range []int{0, 1} {
		done := make(chan bool)
		start := time.Now()
		for i := 0; i < MAX_COUNT; i++ {
			go func(i int) {
				podData := &api.EtcdPodMetadata{NamespaceName: "ns", PodName: fmt.Sprintf("pod%d", i)}
				_, err = nuageetcd.AllocateSubnetForPod(podData)
				if err != nil {
					t.Fatalf("incrementing active ip count failed with error: %v", err)
				}
				done <- true
			}(i)
		}
		for i := 0; i < MAX_COUNT; i++ {
			<-done
		}
		end := time.Now()
		fmt.Printf("scale up time is %s\n", end.Sub(start))
		getResp, err = nuageetcd.client.Get(context.Background(), SUBNET_METADATA_TREE+"ns/ns", clientv3.WithPrefix())
		if err != nil {
			t.Fatalf("Getting key from etcd failed with error: %v", err)
		}

		if len(getResp.Kvs) <= 1 {
			t.Fatalf("Subnets not scaled automatically. should be more than one")
		}
		for i := 0; i < len(getResp.Kvs); i++ {
			s := &etcdSubnetValue{}
			if err := json.Unmarshal(getResp.Kvs[i].Value, s); err != nil {
				t.Errorf("unmarshal string %s to struct failed: %v", getResp.Kvs[i].Value, err)
			}

			fmt.Printf("UP Subnet%d: key = %s value = %d\n", i, getResp.Kvs[i].Key, s.ACTIVEIP)
		}

		start = time.Now()
		for i := 0; i < MAX_COUNT; i++ {
			go func(i int) {
				podData := &api.EtcdPodMetadata{NamespaceName: "ns", PodName: fmt.Sprintf("pod%d", i)}
				_, err = nuageetcd.DeAllocateSubnetFromPod(podData)
				if err != nil {
					t.Fatalf("incrementing active ip count failed with error: %v", err)
				}
				done <- true
			}(i)
		}
		for i := 0; i < MAX_COUNT; i++ {
			<-done
		}
		end = time.Now()
		fmt.Printf("scale down time is %s\n", end.Sub(start))
		getResp, err = nuageetcd.client.Get(context.Background(), SUBNET_METADATA_TREE+"ns/ns", clientv3.WithPrefix())
		if err != nil {
			t.Fatalf("Getting key from etcd failed with error: %v", err)
		}

		for i := 0; i < len(getResp.Kvs); i++ {
			s := &etcdSubnetValue{}
			if err := json.Unmarshal(getResp.Kvs[i].Value, s); err != nil {
				t.Errorf("unmarshal string %s to struct failed: %v", getResp.Kvs[i].Value, err)
			}
			fmt.Printf("DOWN Subnet%d: key = %s value = %d\n", i, getResp.Kvs[i].Key, s.ACTIVEIP)
		}

		if len(getResp.Kvs) != 1 {
			t.Fatalf("Subnets not scaled automatically. should be only one")
		}
	}

	firstSubnet.Name = path.Base(string(getResp.Kvs[0].Key))
	if _, err = nuageetcd.DeleteLastSubnet(firstSubnet); err != nil {
		t.Fatalf("Creating subnet path failed with error: %v", err)
	}
}

func TestParallelReadsWrites(t *testing.T) {
	var err error
	var getResp *clientv3.GetResponse
	firstSubnet := &api.EtcdSubnetMetadata{
		Name:      "ns-0",
		Namespace: "ns",
		ID:        "ahakghkabgbwbgbabviwrgbgrbgi",
		CIDR:      "0.0.0.0/0",
	}

	if err = nuageetcd.CreateFirstSubnet(firstSubnet); err != nil {
		t.Fatalf("Creating subnet path failed with error: %v", err)
	}
	//scale up and scale down twice
	done := make(chan bool)
	for i := 0; i < MAX_COUNT; i++ {
		go func(i int) {
			podData := &api.EtcdPodMetadata{NamespaceName: "ns", PodName: fmt.Sprintf("pod%d", i)}
			_, err = nuageetcd.AllocateSubnetForPod(podData)
			if err != nil {
				t.Fatalf("incrementing active ip count failed with error: %v", err)
			}
			done <- true
		}(i)
	}
	for i := 0; i < MAX_COUNT; i++ {
		<-done
	}

	getResp, err = nuageetcd.client.Get(context.Background(), SUBNET_METADATA_TREE+"ns/ns", clientv3.WithPrefix())
	if err != nil {
		t.Fatalf("Getting key from etcd failed with error: %v", err)
	}

	if len(getResp.Kvs) <= 1 {
		t.Fatalf("Subnets not scaled automatically. should be more than one")
	}
	for i := 0; i < len(getResp.Kvs); i++ {
		s := &etcdSubnetValue{}
		if err := json.Unmarshal(getResp.Kvs[i].Value, s); err != nil {
			t.Errorf("unmarshal string %s to struct failed: %v", getResp.Kvs[i].Value, err)
		}
		fmt.Printf("UP Subnet%d: key = %s value = %d\n", i, getResp.Kvs[i].Key, s.ACTIVEIP)
	}

	go func() {
		for i := 0; i < MAX_COUNT/2; i++ {
			go func(i int) {
				podData := &api.EtcdPodMetadata{NamespaceName: "ns", PodName: fmt.Sprintf("pod%d", i)}
				_, err = nuageetcd.DeAllocateSubnetFromPod(podData)
				if err != nil {
					t.Fatalf("incrementing active ip count failed with error: %v", err)
				}
				done <- true
			}(i)
		}
	}()
	go func() {
		for i := MAX_COUNT; i < 3*MAX_COUNT/2; i++ {
			go func(i int) {
				podData := &api.EtcdPodMetadata{NamespaceName: "ns", PodName: fmt.Sprintf("pod%d", i)}
				_, err = nuageetcd.AllocateSubnetForPod(podData)
				if err != nil {
					t.Fatalf("incrementing active ip count failed with error: %v", err)
				}
				done <- true
			}(i)
		}
	}()
	for i := 0; i < MAX_COUNT; i++ {
		<-done
	}

	getResp, err = nuageetcd.client.Get(context.Background(), SUBNET_METADATA_TREE+"ns/ns", clientv3.WithPrefix())
	if err != nil {
		t.Fatalf("Getting key from etcd failed with error: %v", err)
	}

	_count := 0
	for i := 0; i < len(getResp.Kvs); i++ {
		s := &etcdSubnetValue{}
		if err := json.Unmarshal(getResp.Kvs[i].Value, s); err != nil {
			t.Errorf("unmarshal string %s to struct failed: %v", getResp.Kvs[i].Value, err)
		}
		fmt.Printf("UP-DOWN Subnet%d: key = %s value = %d\n", i, getResp.Kvs[i].Key, s.ACTIVEIP)
		_count += s.ACTIVEIP
	}
	if _count != MAX_COUNT {
		t.Fatalf("ip count(%d) did not match pod count(%d)", _count, MAX_COUNT)
	}

	for i := MAX_COUNT / 2; i < 3*MAX_COUNT/2; i++ {
		go func(i int) {
			podData := &api.EtcdPodMetadata{NamespaceName: "ns", PodName: fmt.Sprintf("pod%d", i)}
			_, err = nuageetcd.DeAllocateSubnetFromPod(podData)
			if err != nil {
				t.Fatalf("incrementing active ip count failed with error: %v", err)
			}
			done <- true
		}(i)
	}
	for i := 0; i < MAX_COUNT; i++ {
		<-done
	}

	getResp, err = nuageetcd.client.Get(context.Background(), SUBNET_METADATA_TREE+"ns/ns", clientv3.WithPrefix())
	if err != nil {
		t.Fatalf("Getting key from etcd failed with error: %v", err)
	}
	if len(getResp.Kvs) != 1 {
		t.Fatalf("Subnets not scaled automatically. should be only one")
	}

	firstSubnet.Name = path.Base(string(getResp.Kvs[0].Key))
	if _, err = nuageetcd.DeleteLastSubnet(firstSubnet); err != nil {
		t.Fatalf("Creating subnet path failed with error: %v", err)
	}
}

func TestAllocateDeallocateSubnetsUnit(t *testing.T) {
	count := 0
	subnet := &api.EtcdSubnetMetadata{
		CIDR: "0.0.0.0/0",
		Name: "ns-0",
	}
	done := make(chan bool)
	for i := 0; i < 3; i++ {
		go func() {
			allocatedSubnet, err := nuageetcd.AllocateSubnetCIDR(subnet)
			if err != nil {
				t.Fatalf("allocating subnet %v failed: %v", subnet, err)
			}
			if allocatedSubnet == "" {
				count += 1
			} else {
				fmt.Printf("cidr allocated to subnet %s\n", allocatedSubnet)
			}
			done <- true
		}()
	}
	for i := 0; i < 3; i++ {
		<-done
	}
	if count != 1 {
		t.Fatalf("subnet cidr allocated more than once")
	}
	for i := 0; i < 3; i++ {
		go func() {
			if err := nuageetcd.FreeSubnetCIDR(subnet); err != nil {
				t.Fatalf("deallocating subnet failed")
			}
			done <- true
		}()
	}
	for i := 0; i < 3; i++ {
		<-done
	}
}

func TestAllocateDeallocateSubnetsScale(t *testing.T) {
	done := make(chan bool)
	for i := 0; i < MAX_COUNT; i++ {
		go func(i int) {
			subnet := &api.EtcdSubnetMetadata{
				CIDR: fmt.Sprintf("0.0.0.%d/0", i),
				Name: fmt.Sprintf("subnet%d", i),
			}
			if _, err := nuageetcd.AllocateSubnetCIDR(subnet); err != nil {
				t.Fatalf("allocating subnet failed")
			}
			done <- true
		}(i)
	}
	for i := 0; i < MAX_COUNT; i++ {
		<-done
	}
	for i := 0; i < MAX_COUNT; i++ {
		go func(i int) {
			subnet := &api.EtcdSubnetMetadata{
				CIDR: fmt.Sprintf("0.0.0.%d/0", i),
				Name: fmt.Sprintf("subnet%d", i),
			}
			if err := nuageetcd.FreeSubnetCIDR(subnet); err != nil {
				t.Fatalf("deleting subnet failed")
			}
			done <- true
		}(i)
	}
	for i := 0; i < MAX_COUNT; i++ {
		<-done
	}
}

func TestGetSubnetID(t *testing.T) {
	firstSubnet := &api.EtcdSubnetMetadata{
		Name:      "ns-0",
		Namespace: "ns",
		ID:        "ahakghkabgbwbgbabviwrgbgrbgi",
		CIDR:      "0.0.0.0/0",
	}
	if err := nuageetcd.CreateFirstSubnet(firstSubnet); err != nil {
		t.Fatalf("Creating subnet path failed with error: %v", err)
	}
	vsd_network_id := "test"
	done := make(chan bool)
	for i := 0; i < nuageetcd.maxIPCount; i++ {
		go func(i int) {
			podData := &api.EtcdPodMetadata{NamespaceName: "ns", PodName: fmt.Sprintf("pod%d", i)}
			podSubnet, err := nuageetcd.AllocateSubnetForPod(podData)
			if err != nil {
				t.Fatalf("incrementing active ip count failed with error: %v", err)
			}
			if podSubnet.ToCreate != "" {
				go func() {
					subnet := &api.EtcdSubnetMetadata{Namespace: "ns", Name: podSubnet.ToCreate}
					id, err := nuageetcd.GetSubnetID(subnet)
					if err != nil {
						t.Fatalf("Getting subnet id failed: %v", err)
						return
					}
					if id != vsd_network_id {
						t.Fatalf("Expected id = %s found = %s", vsd_network_id, id)
					}
					done <- true
				}()
				subnet := &api.EtcdSubnetMetadata{Namespace: "ns", Name: podSubnet.ToCreate, ID: vsd_network_id, CIDR: "0.0.0.0/24"}
				if err := nuageetcd.UpdateSubnetInfo(subnet); err != nil {
					t.Fatalf("updating subnet id failed: %v", err)
				}
			}
			done <- true
		}(i)
	}
	for i := 0; i < nuageetcd.maxIPCount+1; i++ {
		<-done
	}

	getResp, err := nuageetcd.client.Get(context.Background(), "/", clientv3.WithPrefix())
	if err != nil {
		t.Fatalf("Getting key from etcd failed with error: %v", err)
	}
	for i := 0; i < len(getResp.Kvs); i++ {
		fmt.Printf("UP Subnet%d: key = %s value = %s\n", i, getResp.Kvs[i].Key, getResp.Kvs[i].Value)
	}

	for i := 0; i < nuageetcd.maxIPCount; i++ {
		go func(i int) {
			podData := &api.EtcdPodMetadata{NamespaceName: "ns", PodName: fmt.Sprintf("pod%d", i)}
			_, err := nuageetcd.DeAllocateSubnetFromPod(podData)
			if err != nil {
				t.Fatalf("incrementing active ip count failed with error: %v", err)
			}
			done <- true
		}(i)
	}

	for i := 0; i < nuageetcd.maxIPCount; i++ {
		<-done
	}

	getResp, err = nuageetcd.client.Get(context.Background(), SUBNET_METADATA_TREE+"ns/ns", clientv3.WithPrefix())
	if err != nil {
		t.Fatalf("Getting key from etcd failed with error: %v", err)
	}

	if len(getResp.Kvs) != 1 {
		t.Fatalf("Subnets not scaled automatically. should be only one")
	}

	firstSubnet.Name = path.Base(string(getResp.Kvs[0].Key))
	if _, err = nuageetcd.DeleteLastSubnet(firstSubnet); err != nil {
		t.Fatalf("Creating subnet path failed with error: %v", err)
	}
}

func TestZoneCRUD(t *testing.T) {
	zoneInfo := &api.EtcdZoneMetadata{Name: "test-zone"}
	done := make(chan bool)
	count := 0
	for i := 0; i < 3; i++ {
		go func() {
			id, err := nuageetcd.AddZone(zoneInfo)
			if err != nil {
				t.Fatalf("adding zone failed: %v", err)
			}
			if id == "" {
				count += 1
				zoneInfo1 := &api.EtcdZoneMetadata{Name: "test-zone", ID: "test-id"}
				err := nuageetcd.UpdateZone(zoneInfo1)
				if err != nil {
					t.Fatalf("updating zone failed: %v", err)
				}
			} else {
				fmt.Printf("zone id received is: %s\n", id)
			}
			done <- true
		}()
	}
	for i := 0; i < 3; i++ {
		<-done
	}
	if count != 1 {
		t.Fatalf("more than one go routine updated the zone information")
	}
	err := nuageetcd.DeleteZone(zoneInfo)
	if err != nil {
		t.Fatalf("deleting zone failed:%v", err)
	}

	if getResp, err := nuageetcd.client.Get(context.Background(), ZONE_METADATA_TREE+zoneInfo.Name); err != nil {
		t.Fatalf("getting zone %s info failed: %v", zoneInfo.Name, err)
	} else if len(getResp.Kvs) != 0 {
		t.Fatalf("deleting zone from etcd failed. still received more than zero zones")
	}
}
