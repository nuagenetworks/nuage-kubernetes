/*
###########################################################################
#
#   Filename:           ipv4subnet_test.go
#
#   Author:             Ryan Fredette
#   Created:            August 24, 2015
#
#   Description:        tests of functionality implemented in ipv4subnet.go
#
###########################################################################
#
#              Copyright (c) 2015 Nuage Networks
#
###########################################################################

*/

package client

import (
	"bytes"
	"testing"
)

func TestSplitIPv4Subnet(t *testing.T) {
	input := []IPv4Subnet{
		// Lowest valid input (endpoint checking)
		{IPv4Address{0, 0, 0, 0}, 0},
		// A value in the middle
		{IPv4Address{192, 168, 1, 0}, 24},
		// A subnet definition with non-zero bits outside the mask.  It should
		// return exactly the same as the 192.168.1.0/24 test.
		{IPv4Address{192, 168, 1, 23}, 24},
		{IPv4Address{192, 168, 2, 0}, 23},
		// A subnet definition with non-zero bits outside the mask.  It should
		// return exactly the same as the 192.168.2.0/23 test.
		{IPv4Address{192, 168, 3, 0}, 23},
	}
	output := [][2]IPv4Subnet{
		//{IPv4Address{0, 0, 0, 0}, 0},
		[2]IPv4Subnet{
			{IPv4Address{0, 0, 0, 0}, 1},
			{IPv4Address{128, 0, 0, 0}, 1},
		},
		//{IPv4Address{192, 168, 1, 0}, 24},
		[2]IPv4Subnet{
			{IPv4Address{192, 168, 1, 0}, 25},
			{IPv4Address{192, 168, 1, 128}, 25},
		},
		//{IPv4Address{192, 168, 1, 23}, 24},
		[2]IPv4Subnet{
			{IPv4Address{192, 168, 1, 0}, 25},
			{IPv4Address{192, 168, 1, 128}, 25},
		},
		//{IPv4Address{192, 168, 2, 0}, 23},
		[2]IPv4Subnet{
			{IPv4Address{192, 168, 2, 0}, 24},
			{IPv4Address{192, 168, 3, 0}, 24},
		},
		//{IPv4Address{192, 168, 3, 0}, 23},
		[2]IPv4Subnet{
			{IPv4Address{192, 168, 2, 0}, 24},
			{IPv4Address{192, 168, 3, 0}, 24},
		},
	}
	for i, inSubnet := range input {
		loSubnet, hiSubnet, err := inSubnet.Split()
		if err != nil {
			t.Errorf("Split() failed. Error: %s", err)
			t.FailNow()
		}
		if loSubnet.Address != output[i][0].Address {
			t.Errorf("loSubnet Address mismatch! Expected: %s, got %s",
				output[i][0].Address, loSubnet.Address)
			t.Fail()
		}
		if loSubnet.CIDRMask != output[i][0].CIDRMask {
			t.Errorf("loSubnet CIDRMask mismatch! Expected: %v, got %v",
				output[i][0].CIDRMask, loSubnet.CIDRMask)
			t.Fail()
		}
		if hiSubnet.Address != output[i][1].Address {
			t.Errorf("hiSubnet Address mismatch! Expected: %s, got %s",
				output[i][1].Address, hiSubnet.Address)
			t.Fail()
		}
		if hiSubnet.CIDRMask != output[i][1].CIDRMask {
			t.Errorf("hiSubnet CIDRMask mismatch! Expected: %v, got %v",
				output[i][1].CIDRMask, hiSubnet.CIDRMask)
			t.Fail()
		}
	}
	// Test an invalid split.  This should return an error.
	inSubnet := IPv4Subnet{IPv4Address{192, 168, 1, 1}, 32}
	loSubnet, hiSubnet, err := inSubnet.Split()
	if loSubnet != nil {
		t.Errorf("Split() on %s succeeded! Produced loSubnet %s",
			inSubnet, loSubnet)
		t.Fail()
	}
	if hiSubnet != nil {
		t.Errorf("Split() on %s succeeded! Produced hiSubnet %s",
			inSubnet, hiSubnet)
		t.Fail()
	}
	if err == nil {
		t.Errorf("Split() on %s returned no result, but also no error!", inSubnet)
		t.Fail()
	}
}

func TestIPv4SubnetFromString(t *testing.T) {
	input := []string{
		"0.0.0.0/0",
		"10.0.0.0/8",
		"192.168.122.0/24",
		"172.30.1.1/32",
	}
	output := []IPv4Subnet{
		//"0.0.0.0/0",
		IPv4Subnet{IPv4Address{0, 0, 0, 0}, 0},
		//"10.0.0.0/8",
		IPv4Subnet{IPv4Address{10, 0, 0, 0}, 8},
		//"192.168.122.0/24",
		IPv4Subnet{IPv4Address{192, 168, 122, 0}, 24},
		//"172.30.1.1/32",
		IPv4Subnet{IPv4Address{172, 30, 1, 1}, 32},
	}
	for i, inString := range input {
		outSubnet, err := IPv4SubnetFromString(inString)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
			t.FailNow()
		}
		if bytes.Compare(outSubnet.Address[:], output[i].Address[:]) != 0 {
			t.Errorf("Address mismatch! Expected %s, got %s",
				output[i].Address, outSubnet.Address)
			t.Fail()
		}
		if outSubnet.CIDRMask != output[i].CIDRMask {
			t.Errorf("CIDRMask mismatch! Expected %v, got %v",
				output[i].CIDRMask, outSubnet.CIDRMask)
			t.Fail()
		}
	}
}
