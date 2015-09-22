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
	"errors"
	"fmt"
	"github.com/nuagenetworks/openshift-integration/nuagekubemon/api"
	"testing"
)

func deleteEnterprise(t *testing.T, vsdClient *NuageVsdClient, id string) error {
	result := make([]struct{}, 1)
	e := api.RESTError{}
	resp, err := vsdClient.session.Delete(vsdClient.url+"enterprises/"+
		id+"?responseChoice=1", &result, &e)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Status() != 204 {
		t.Error("Bad response status from VSD Server when deleting enterprise")
		t.Errorf("\t Status:  %v\n", resp.Status())
		t.Errorf("\t Message: %v\n", e.Message)
		t.Errorf("\t Errors: %v\n", e.Message)
		return errors.New("Unexpected error code: " +
			fmt.Sprintf("%v", resp.Status()))
	}
	return nil
}

func TestCreateEnterprise(t *testing.T) {
	if vsdClient == nil {
		t.Skip("Needs VSD connection")
	}
	myEnterpriseName := "openshift-test-enterprise"
	// Verify that the enterprise we're trying to create doesn't already exist
	_, err := vsdClient.GetEnterpriseID(myEnterpriseName)
	if err != nil && err.Error() != "Enterprise not found" {
		t.Fatal("Unexpected error:", err)
	}
	if err == nil {
		t.Fatalf("Enterprise %q already exists!", myEnterpriseName)
	}
	// Create it
	enterpriseID, err := vsdClient.CreateEnterprise(myEnterpriseName)
	if err != nil {
		t.Fatal(err)
	}
	// Guarantee that the enterprise gets deleted even in error cases
	defer deleteEnterprise(t, vsdClient, enterpriseID)
	// Verify that it exists now
	id, err := vsdClient.GetEnterpriseID(myEnterpriseName)
	if err != nil && err.Error() != "Enterprise not found" {
		t.Fatal("Unexpected error:", err)
	}
	if id != enterpriseID {
		t.Fatalf("Enterprise ID mismatch! CreateEnterprise() returned %v, "+
			"GetEnterpriseID() returned %v.", enterpriseID, id)
	}
}

func TestCreateAdminUser(t *testing.T) {
	if vsdClient == nil {
		t.Skip("Needs VSD connection")
	}
	myEnterpriseName := "openshift-test-enterprise"
	// Verify that the enterprise we're trying to create doesn't already exist
	_, err := vsdClient.GetEnterpriseID(myEnterpriseName)
	if err != nil && err.Error() != "Enterprise not found" {
		t.Fatal("Unexpected error:", err)
	}
	// Create it
	enterpriseID, err := vsdClient.CreateEnterprise(myEnterpriseName)
	if err != nil {
		t.Fatal(err)
	}
	// Guarantee that the enterprise gets deleted even in error cases
	defer deleteEnterprise(t, vsdClient, enterpriseID)
	adminID, err := vsdClient.CreateAdminUser(enterpriseID, "admin", "admin")
	if err != nil {
		t.Fatal(err)
	}
	// Verify that the admin user exists now
	id, err := vsdClient.GetAdminID(enterpriseID, "admin")
	if err != nil && err.Error() != "User not found" {
		t.Fatal("Unexpected error:", err)
	}
	if id != adminID {
		t.Fatalf("Admin ID mismatch! CreateAdminUser() returned %v, "+
			"GetAdminID() returned %v.", adminID, id)
	}
	// Verify that the admin is in the ORGADMIN group
	groupID, err := vsdClient.GetAdminGroupID(enterpriseID)
	if err != nil && err.Error() != "User not found" {
		t.Fatal("Unexpected error:", err)
	}
	result := make([]api.VsdUser, 1)
	e := api.RESTError{}
	vsdClient.session.Header.Add("X-Nuage-Filter", `userName == "admin"`)
	resp, err := vsdClient.session.Get(vsdClient.url+"groups/"+groupID+"/users", nil, &result, &e)
	vsdClient.session.Header.Del("X-Nuage-Filter")
	if err != nil {
		t.Fatal(err)
	}
	if resp.Status() != 200 {
		t.Errorf("\t Status:  %v\n", resp.Status())
		t.Errorf("\t Message: %v\n", e.Message)
		t.Errorf("\t Errors: %v\n", e.Message)
		t.Fatal("Bad response status from VSD Server")
	}
	if result[0].ID != adminID {
		t.Fatal("Admin not in ORGADMIN group!")
	}
}

func TestCreateTemplates(t *testing.T) {
	if vsdClient == nil {
		t.Skip("Needs VSD connection")
	}
	myEnterpriseName := "openshift-test-enterprise"
	// Verify that the enterprise we're trying to create doesn't already exist
	_, err := vsdClient.GetEnterpriseID(myEnterpriseName)
	if err != nil && err.Error() != "Enterprise not found" {
		t.Fatal("Unexpected error:", err)
	}
	// Create it
	enterpriseID, err := vsdClient.CreateEnterprise(myEnterpriseName)
	if err != nil {
		t.Fatal(err)
	}
	// Guarantee that the enterprise gets deleted even in error cases
	defer deleteEnterprise(t, vsdClient, enterpriseID)
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
	// Create a zone template
	zoneID, err := vsdClient.CreateZoneTemplate(domainID, "zone-template")
	if err != nil {
		t.Fatal(err)
	}
	// Verify that the template exists
	id, err = vsdClient.GetZoneTemplateID(domainID, "zone-template")
	if err != nil {
		t.Fatal(err)
	}
	if id != zoneID {
		t.Fatalf("Zone template ID mismatch! CreateZoneTemplate() returned "+
			"%v, GetZoneTemplateID() returned %v.", zoneID, id)
	}
}

func TestCreateDomain(t *testing.T) {
	if vsdClient == nil {
		t.Skip("Needs VSD connection")
	}
	myEnterpriseName := "openshift-test-enterprise"
	// Verify that the enterprise we're trying to create doesn't already exist
	_, err := vsdClient.GetEnterpriseID(myEnterpriseName)
	if err != nil && err.Error() != "Enterprise not found" {
		t.Fatal("Unexpected error:", err)
	}
	// Create it
	enterpriseID, err := vsdClient.CreateEnterprise(myEnterpriseName)
	if err != nil {
		t.Fatal(err)
	}
	// Guarantee that the enterprise gets deleted even in error cases
	defer deleteEnterprise(t, vsdClient, enterpriseID)
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
	defer vsdClient.DeleteDomain("test-domain", domainID)
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
	if result[0].PATEnabled != "ENABLED" {
		t.Fatalf("Domain PATEnabled status mismatch! Expected \"ENABLED\", got %q",
			result[0].PATEnabled)
	}
}

func TestDeleteDomain(t *testing.T) {
	if vsdClient == nil {
		t.Skip("Needs VSD connection")
	}
	myEnterpriseName := "openshift-test-enterprise"
	// Verify that the enterprise we're trying to create doesn't already exist
	_, err := vsdClient.GetEnterpriseID(myEnterpriseName)
	if err != nil && err.Error() != "Enterprise not found" {
		t.Fatal("Unexpected error:", err)
	}
	// Create it
	enterpriseID, err := vsdClient.CreateEnterprise(myEnterpriseName)
	if err != nil {
		t.Fatal(err)
	}
	// Guarantee that the enterprise gets deleted even in error cases
	defer deleteEnterprise(t, vsdClient, enterpriseID)
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
	// Verify that it was instantiated
	id, err := vsdClient.GetDomainID(enterpriseID, "test-domain")
	if err != nil {
		t.Fatal(err)
	}
	if id != domainID {
		t.Fatal("Domain was not instantiated! Aborting test.")
	}
	vsdClient.DeleteDomain("test-domain", domainID)
	id, err = vsdClient.GetDomainID(enterpriseID, "test-domain")
	if err == nil || err.Error() != "Domain not found" {
		t.Fatal(err)
	}
	if id == domainID {
		t.Fatal("Domain was not deleted!")
	}
}

func TestCreateSubnet(t *testing.T) {
	if vsdClient == nil {
		t.Skip("Needs VSD connection")
	}
	// Create an enterprise
	myEnterpriseName := "openshift-test-enterprise"
	// Verify that the enterprise we're trying to create doesn't already exist
	_, err := vsdClient.GetEnterpriseID(myEnterpriseName)
	if err != nil && err.Error() != "Enterprise not found" {
		t.Fatal("Unexpected error:", err)
	}
	// Create it
	enterpriseID, err := vsdClient.CreateEnterprise(myEnterpriseName)
	if err != nil {
		t.Fatal(err)
	}
	// Guarantee that it's deleted when we're done
	defer deleteEnterprise(t, vsdClient, enterpriseID)
	// Create a domain template
	domainTemplateID, err := vsdClient.CreateDomainTemplate(enterpriseID, "domain-template")
	if err != nil {
		t.Fatal(err)
	}
	// Create a zone template
	_, err = vsdClient.CreateZoneTemplate(domainTemplateID, "zone-template")
	if err != nil {
		t.Fatal(err)
	}
	// Instantiate a domain from the domain template
	domainID, err := vsdClient.CreateDomain(enterpriseID, domainTemplateID, "test-domain")
	if err != nil {
		t.Fatal(err)
	}
	// Guarantee that it's deleted when we're done too
	defer vsdClient.DeleteDomain("test-domain", domainID)
	// Get the ID of the zone that was instantiated with the domain
	zoneID, err := vsdClient.GetZoneID(domainID, "zone-template")
	if err != nil {
		t.Fatal(err)
	}
	// Create a subnet with specific parameters in the zone
	subnet, err := IPv4SubnetFromString("10.1.1.0/24")
	if err != nil {
		t.Fatal(err)
	}
	subnetID, err := vsdClient.CreateSubnet(zoneID, subnet)
	if err != nil {
		t.Fatal(err)
	}
	// Guarantee that the subnet gets deleted when we're done too
	defer vsdClient.DeleteSubnet(subnetID)
	// Verify that it was created as defined
	id, err := vsdClient.GetSubnetID(zoneID, subnet)
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
	if result[0].PATEnabled != "ENABLED" {
		t.Fatalf("Subnet PATEnabled status mismatch! Expected \"ENABLED\", got %q",
			result[0].PATEnabled)
	}
}

func TestDeleteSubnet(t *testing.T) {
	if vsdClient == nil {
		t.Skip("Needs VSD connection")
	}
	// Create an enterprise
	myEnterpriseName := "openshift-test-enterprise"
	// Verify that the enterprise we're trying to create doesn't already exist
	_, err := vsdClient.GetEnterpriseID(myEnterpriseName)
	if err != nil && err.Error() != "Enterprise not found" {
		t.Fatal("Unexpected error:", err)
	}
	// Create it
	enterpriseID, err := vsdClient.CreateEnterprise(myEnterpriseName)
	if err != nil {
		t.Fatal(err)
	}
	// Guarantee that it's deleted when we're done
	defer deleteEnterprise(t, vsdClient, enterpriseID)
	// Create a domain template
	domainTemplateID, err := vsdClient.CreateDomainTemplate(enterpriseID, "domain-template")
	if err != nil {
		t.Fatal(err)
	}
	// Create a zone template
	_, err = vsdClient.CreateZoneTemplate(domainTemplateID, "zone-template")
	if err != nil {
		t.Fatal(err)
	}
	// Instantiate a domain from the domain template
	domainID, err := vsdClient.CreateDomain(enterpriseID, domainTemplateID, "test-domain")
	if err != nil {
		t.Fatal(err)
	}
	// Guarantee that it's deleted when we're done too
	defer vsdClient.DeleteDomain("test-domain", domainID)
	// Get the ID of the zone that was instantiated with the domain
	zoneID, err := vsdClient.GetZoneID(domainID, "zone-template")
	if err != nil {
		t.Fatal(err)
	}
	// Create a subnet with specific parameters in the zone
	subnet, err := IPv4SubnetFromString("10.1.1.0/24")
	if err != nil {
		t.Fatal(err)
	}
	subnetID, err := vsdClient.CreateSubnet(zoneID, subnet)
	if err != nil {
		t.Fatal(err)
	}
	// Delete it
	err = vsdClient.DeleteSubnet(subnetID)
	if err != nil {
		t.Fatal(err)
	}
	// Verify that it no longer exists
	_, err = vsdClient.GetSubnetID(zoneID, subnet)
	if err == nil {
		t.Fatal("Subnet not deleted!")
	}
}
