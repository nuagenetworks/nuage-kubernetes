/*
###########################################################################
#
#   Filename:           vsdclient_test.go
#
#   Author:             Ryan Fredette
#   Created:            August 10, 2015
#
#   Description:        tests of functionality implemented in
#                       nuagevsdclient.go
#
###########################################################################
#
#              Copyright (c) 2015 Nuage Networks
#
###########################################################################

*/

package client

import (
	"testing"

	"github.com/nuagenetworks/nuage-kubernetes/nuagekubemon/api"
)

func TestGetEnterpriseID(t *testing.T) {
	if vsdClient.enterpriseID == "" {
		t.Skip("Needs VSD connection")
	}
	myEnterpriseName := "openshift-test-enterprise"
	// Verify that the enterprise we're trying to create doesn't already exist
	_, err := vsdClient.GetEnterpriseID(myEnterpriseName)
	if err != nil && err.Error() != "Enterprise not found" {
		t.Fatal("Unexpected error:", err)
	}
}

func TestCreateTemplates(t *testing.T) {
	if vsdClient.enterpriseID == "" {
		t.Skip("Needs VSD connection")
	}
	myEnterpriseName := "openshift-test-enterprise"
	// Verify that the enterprise we're trying to create doesn't already exist
	enterpriseID, err := vsdClient.GetEnterpriseID(myEnterpriseName)
	if err != nil && err.Error() != "Enterprise not found" {
		t.Fatal("Unexpected error:", err)
	}
	// Create a domain template
	domainID, err := vsdClient.CreateDomainTemplate(enterpriseID, "domain-template")
	if err != nil {
		t.Fatal(err)
	}
	// Verify that the template exists
	id, err := vsdClient.GetDomainTemplateID(enterpriseID, "domain-template")
	if err != nil {
		t.Fatal(err)
	}
	if id != domainID {
		t.Fatalf("Domain template ID mismatch! CreateDomainTemplate() "+
			"returned %v, GetDomainTemplateID() returned %v.", domainID, id)
	}
}

func TestCreateDomain(t *testing.T) {
	if vsdClient.enterpriseID == "" {
		t.Skip("Needs VSD connection")
	}
	myEnterpriseName := "openshift-test-enterprise"
	// Verify that the enterprise we're trying to create doesn't already exist
	enterpriseID, err := vsdClient.GetEnterpriseID(myEnterpriseName)
	if err != nil && err.Error() != "Enterprise not found" {
		t.Fatal("Unexpected error:", err)
	}
	// Create a domain template
	domainTemplateID, err := vsdClient.CreateDomainTemplate(enterpriseID, "domain-template")
	if err != nil {
		t.Fatal(err)
	}
	// Instantiate a domain
	domainID, err := vsdClient.CreateDomain(enterpriseID, domainTemplateID, "test-domain")
	if err != nil {
		t.Fatal(err)
	}
	// Guarantee that the domain will be deleted even in error cases
	defer vsdClient.DeleteDomain(domainID)
	// Verify that it was instantiated
	id, err := vsdClient.GetDomainID(enterpriseID, "test-domain")
	if err != nil {
		t.Fatal(err)
	}
	if id != domainID {
		t.Fatalf("Domain ID mismatch! CreateDomain() returned %v, "+
			"GetDomainID() returned %v.", domainID, id)
	}
	// Verify that Address Translation (PAT) was enabled
	result := make([]api.VsdDomain, 1)
	e := api.RESTError{}
	response, err := vsdClient.session.Get(vsdClient.url+"domains/"+domainID, nil, &result, &e)
	if err != nil {
		t.Fatalf("Failed GET on %s: %s", vsdClient.url+"domains/"+domainID, err)
	}
	if response.Status() != 200 {
		t.Fatalf("Got unexpected response to GET on %s: code %d\nraw text:\n%s",
			vsdClient.url+"domains/"+domainID, response.Status(), response.RawText())
	}
	if result[0].UnderlayEnabled != "ENABLED" {
		t.Fatalf("Domain PATEnabled status mismatch! Expected \"ENABLED\", got %q",
			result[0].UnderlayEnabled)
	}
}

func TestDeleteDomain(t *testing.T) {
	if vsdClient.enterpriseID == "" {
		t.Skip("Needs VSD connection")
	}
	myEnterpriseName := "openshift-test-enterprise"
	// Verify that the enterprise we're trying to create doesn't already exist
	enterpriseID, err := vsdClient.GetEnterpriseID(myEnterpriseName)
	if err != nil && err.Error() != "Enterprise not found" {
		t.Fatal("Unexpected error:", err)
	}
	domainTemplateID, err := vsdClient.CreateDomainTemplate(enterpriseID, "domain-template")
	if err != nil {
		t.Fatal(err)
	}
	// Instantiate a domain
	domainID, err := vsdClient.CreateDomain(enterpriseID, domainTemplateID, "test-domain")
	if err != nil {
		t.Fatal(err)
	}
	// Verify that it was instantiated
	id, err := vsdClient.GetDomainID(enterpriseID, "test-domain")
	if err != nil {
		t.Fatal(err)
	}
	if id != domainID {
		t.Fatal("Domain was not instantiated! Aborting test.")
	}
	vsdClient.DeleteDomain(domainID)
	id, err = vsdClient.GetDomainID(enterpriseID, "test-domain")
	if err == nil || err.Error() != "Domain not found" {
		t.Fatal(err)
	}
	if id == domainID {
		t.Fatal("Domain was not deleted!")
	}
}

func TestCreateZone(t *testing.T) {
	if vsdClient.enterpriseID == "" {
		t.Skip("Needs VSD connection")
	}
	// Create an enterprise
	myEnterpriseName := "openshift-test-enterprise"
	// Verify that the enterprise we're trying to create doesn't already exist
	enterpriseID, err := vsdClient.GetEnterpriseID(myEnterpriseName)
	if err != nil && err.Error() != "Enterprise not found" {
		t.Fatal("Unexpected error:", err)
	}
	domainTemplateID, err := vsdClient.CreateDomainTemplate(enterpriseID, "domain-template")
	if err != nil {
		t.Fatal(err)
	}
	// Instantiate a domain from the domain template
	domainID, err := vsdClient.CreateDomain(enterpriseID, domainTemplateID, "test-domain")
	if err != nil {
		t.Fatal(err)
	}
	// Guarantee that it's deleted when we're done too
	defer vsdClient.DeleteDomain(domainID)
	// Create a zone inside the domain
	zoneID, err := vsdClient.CreateZone(domainID, "test-zone")
	if err != nil {
		t.Fatal(err)
	}
	// Guarantee that the zone will be deleted even in error cases
	defer vsdClient.DeleteZone(zoneID)
	// Verify that it was instantiated
	id, err := vsdClient.GetZoneID(domainID, "test-zone")
	if err != nil {
		t.Fatal(err)
	}
	if id != zoneID {
		t.Fatalf("Zone ID mismatch! CreateZone() returned %v, "+
			"GetZoneID() returned %v.", zoneID, id)
	}
}

func TestDeleteZone(t *testing.T) {
	if vsdClient.enterpriseID == "" {
		t.Skip("Needs VSD connection")
	}
	// Create an enterprise
	myEnterpriseName := "openshift-test-enterprise"
	// Verify that the enterprise we're trying to create doesn't already exist
	enterpriseID, err := vsdClient.GetEnterpriseID(myEnterpriseName)
	if err != nil && err.Error() != "Enterprise not found" {
		t.Fatal("Unexpected error:", err)
	}
	// Create a domain template
	domainTemplateID, err := vsdClient.CreateDomainTemplate(enterpriseID, "domain-template")
	if err != nil {
		t.Fatal(err)
	}
	// Instantiate a domain from the domain template
	domainID, err := vsdClient.CreateDomain(enterpriseID, domainTemplateID, "test-domain")
	if err != nil {
		t.Fatal(err)
	}
	// Guarantee that it's deleted when we're done too
	defer vsdClient.DeleteDomain(domainID)
	// Create a zone inside the domain
	zoneID, err := vsdClient.CreateZone(domainID, "test-zone")
	if err != nil {
		t.Fatal(err)
	}
	// Verify that it was instantiated
	id, err := vsdClient.GetZoneID(domainID, "test-zone")
	if err != nil {
		t.Fatal(err)
	}
	if id != zoneID {
		t.Fatal("Zone was not instantiated! Aborting test.")
	}
	vsdClient.DeleteZone(zoneID)
	id, err = vsdClient.GetZoneID(domainID, "test-zone")
	if err == nil || err.Error() != "Zone not found" {
		t.Fatal(err)
	}
	if id == zoneID {
		t.Fatal("Zone was not deleted!")
	}
}

func TestCreateSubnet(t *testing.T) {
	if vsdClient.enterpriseID == "" {
		t.Skip("Needs VSD connection")
	}
	// Create an enterprise
	myEnterpriseName := "openshift-test-enterprise"
	// Verify that the enterprise we're trying to create doesn't already exist
	enterpriseID, err := vsdClient.GetEnterpriseID(myEnterpriseName)
	if err != nil && err.Error() != "Enterprise not found" {
		t.Fatal("Unexpected error:", err)
	}
	// Create a domain template
	domainTemplateID, err := vsdClient.CreateDomainTemplate(enterpriseID, "domain-template")
	if err != nil {
		t.Fatal(err)
	}
	// Instantiate a domain from the domain template
	domainID, err := vsdClient.CreateDomain(enterpriseID, domainTemplateID, "test-domain")
	if err != nil {
		t.Fatal(err)
	}
	// Guarantee that it's deleted when we're done too
	defer vsdClient.DeleteDomain(domainID)
	// Get the ID of the zone that was instantiated with the domain
	zoneID, err := vsdClient.CreateZone(domainID, "zone")
	if err != nil {
		t.Fatal(err)
	}
	defer vsdClient.DeleteZone(zoneID)
	// Create a subnet with specific parameters in the zone
	subnet, err := IPv4SubnetFromString("10.1.1.0/24")
	if err != nil {
		t.Fatal(err)
	}
	subnetID, err := vsdClient.CreateSubnet("test-subnet", zoneID, subnet)
	if err != nil {
		t.Fatal(err)
	}
	// Guarantee that the subnet gets deleted when we're done too
	defer vsdClient.DeleteSubnet(subnetID)
	// Verify that it was created as defined
	id, err := vsdClient.GetSubnetID(zoneID, "test-subnet")
	if err != nil {
		t.Fatal(err)
	}
	if subnetID != id {
		t.Fatalf("Subnet ID mismatch! CreateSubnet() returned %v, "+
			"GetSubnetID() returned %v.", subnetID, id)
	}
	result := make([]api.VsdSubnet, 1)
	e := api.RESTError{}
	response, err := vsdClient.session.Get(vsdClient.url+"subnets/"+subnetID, nil, &result, &e)
	if err != nil {
		t.Fatalf("Failed GET on %s: %s", vsdClient.url+"subnets/"+subnetID, err)
	}
	if response.Status() != 200 {
		t.Fatalf("Got unexpected response to GET on %s: code %d\nraw text:\n%s",
			vsdClient.url+"subnets/"+subnetID, response.Status(), response.RawText())
	}
	if result[0].UnderlayEnabled != "INHERITED" {
		t.Fatalf("Subnet PATEnabled status mismatch! Expected \"INHERITED\", got %q",
			result[0].UnderlayEnabled)
	}
}

func TestDeleteSubnet(t *testing.T) {
	if vsdClient.enterpriseID == "" {
		t.Skip("Needs VSD connection")
	}
	// Create an enterprise
	myEnterpriseName := "openshift-test-enterprise"
	// Verify that the enterprise we're trying to create doesn't already exist
	enterpriseID, err := vsdClient.GetEnterpriseID(myEnterpriseName)
	if err != nil && err.Error() != "Enterprise not found" {
		t.Fatal("Unexpected error:", err)
	}
	// Create a domain template
	domainTemplateID, err := vsdClient.CreateDomainTemplate(enterpriseID, "domain-template")
	if err != nil {
		t.Fatal(err)
	}
	// Instantiate a domain from the domain template
	domainID, err := vsdClient.CreateDomain(enterpriseID, domainTemplateID, "test-domain")
	if err != nil {
		t.Fatal(err)
	}
	// Guarantee that it's deleted when we're done too
	defer vsdClient.DeleteDomain(domainID)
	// Get the ID of the zone that was instantiated with the domain
	zoneID, err := vsdClient.CreateZone(domainID, "zone")
	if err != nil {
		t.Fatal(err)
	}
	defer vsdClient.DeleteZone(zoneID)
	// Create a subnet with specific parameters in the zone
	subnet, err := IPv4SubnetFromString("10.1.1.0/24")
	if err != nil {
		t.Fatal(err)
	}
	subnetID, err := vsdClient.CreateSubnet("test-subnet", zoneID, subnet)
	if err != nil {
		t.Fatal(err)
	}
	// Delete it
	err = vsdClient.DeleteSubnet(subnetID)
	if err != nil {
		t.Fatal(err)
	}
	// Verify that it no longer exists
	_, err = vsdClient.GetSubnetID(zoneID, "test-subnet")
	if err == nil {
		t.Fatal("Subnet not deleted!")
	}
}
