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
	// Note: Do we need to delete the admin user as well, or does it get cleaned
	//       up at enterprise deletion?
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
