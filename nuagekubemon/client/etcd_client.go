package client

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/docker/libkv"
	"github.com/docker/libkv/store"
	"github.com/docker/libkv/store/etcd"
	"github.com/golang/glog"
	"github.com/nuagenetworks/nuage-kubernetes/nuagekubemon/config"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

const (
	aclLockKey = "nuageACLLock"
)

type EtcdClient struct {
	valid   bool
	store   store.Store
	aclLock store.Locker
}

func NewEtcdClient(nkmConfig *config.NuageKubeMonConfig) *EtcdClient {
	glog.Infof("Initializing etcd client with config %+v", nkmConfig.EtcdClientConfig)

	etcd.Register()

	etcdConfig := nkmConfig.EtcdClientConfig

	if _, err := os.Stat(etcdConfig.CA); err != nil {
		glog.Warningf("etcd CA cert not found skipping etcd client config")
		return nil
	}

	if _, err := os.Stat(etcdConfig.CertFile); err != nil {
		glog.Warningf("etcd cert not found skipping etcd client config")
		return nil
	}

	if _, err := os.Stat(etcdConfig.KeyFile); err != nil {
		glog.Warningf("etcd key not found skipping etcd client config")
		return nil
	}

	tlsConfig := new(tls.Config)
	tlsConfig.InsecureSkipVerify = true
	caCertBin, err := ioutil.ReadFile(nkmConfig.EtcdClientConfig.CA)
	if err != nil {
		glog.Warningf("Unable to get etcd cert %s, skipping etcd client config",
			nkmConfig.EtcdClientConfig.CA)
		return nil
	}

	caPool := x509.NewCertPool()
	caPool.AppendCertsFromPEM(caCertBin)
	tlsConfig.RootCAs = caPool

	var clientCert tls.Certificate
	if clientCert, err = tls.LoadX509KeyPair(etcdConfig.CertFile, etcdConfig.KeyFile); err != nil {
		glog.Warningf("Unable to load client cert %s and key %s", etcdConfig.CertFile, etcdConfig.KeyFile)
		return nil
	}

	tlsConfig.Certificates = []tls.Certificate{clientCert}

	storeConfig := new(store.Config)
	storeConfig.TLS = tlsConfig

	var editedURLs []string
	for _, url := range nkmConfig.EtcdClientConfig.URLs {
		editedURLs = append(editedURLs, strings.Split(url, "//")[1])
	}

	etcdClient := new(EtcdClient)
	if etcdClient.store, err = libkv.NewStore(store.ETCD, editedURLs, storeConfig); err != nil {
		glog.Errorf("Unable to create libkv store for etcd")
		return nil
	}

	if etcdClient.aclLock, err = etcdClient.store.NewLock(aclLockKey, &store.LockOptions{Value: []byte(aclLockKey),
		TTL: 2 * time.Second}); err != nil {
		glog.Errorf("Unable to create ACL lock")
		return nil
	}

	etcdClient.valid = true

	glog.Info("Created etcd store client")
	return etcdClient
}

func (etcdClient *EtcdClient) ACLLock() error {
	glog.Info("Trying for etcd ACL Lock")
	_, err := etcdClient.aclLock.Lock(nil)
	glog.Infof("Got etcd ACL Lock with err %+v", err)
	return err
}

func (etcdClient *EtcdClient) ACLUnLock() error {
	glog.Info("Unlocking etcd ACL Lock")

	if err := etcdClient.store.Delete(aclLockKey); err != nil {
		glog.Errorf("Error deleting etcd lock for unlocking %s", aclLockKey)
		return fmt.Errorf("Unable to delete the etcd lock")
	}

	etcdClient.aclLock.Unlock()
	return nil
}
