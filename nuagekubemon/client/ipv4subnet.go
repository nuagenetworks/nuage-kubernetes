/*
###########################################################################
#
#   Filename:           ipv4subnet.go
#
#   Author:             Ryan Fredette
#   Created:            August 24, 2015
#
#   Description:        IPv4 Address, Subnet, and Subnet Pool management
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
	"github.com/nuagenetworks/nuage-kubernetes/nuagekubemon/api"
)

type IPv4SubnetNode struct {
	subnet *api.IPv4Subnet
	next   *IPv4SubnetNode
}

type IPv4SubnetPool [33]*IPv4SubnetNode

/* A subnet pool is an array of linked lists.  Each list consists only of
 * subnets with the same CIDR netmask (/0 - /32).  When allocating a subnet
 * with netmask X, the pool will first attempt to pick a subnet of the exact
 * size.  If one is not available, it will get a subnet with netmask X-1, then
 * split it to create 2 subnets with netmask X.  It will return 1 of those
 * subnets to the pool, then return the other one.
 */
func (pool *IPv4SubnetPool) Alloc(size int) (*api.IPv4Subnet, error) {
	if size < 0 || size > 32 {
		return nil, errors.New("Invalid subnet size. Expected between /0 and /32")
	}
	// If there's already at least 1 subnet of the intended size, remove it
	// from the list and return it.
	if pool[size] != nil {
		node := pool[size]
		pool[size] = node.next
		return node.subnet, nil
	}
	// If not, get a larger subnet (1 CIDR mask less), and split it to create 2
	// subnets of the expected size.
	bigSubnet, err := pool.Alloc(size - 1)
	if err != nil {
		return nil, err
	}
	loSubnet, hiSubnet, err := bigSubnet.Split()
	if err != nil {
		pool.Free(bigSubnet)
		return nil, err
	}
	// Of the two subnets from the split, only one is needed, so release the other.
	err = pool.Free(hiSubnet)
	if err != nil {
		pool.Free(bigSubnet)
		return nil, err
	}
	return loSubnet, nil
}

/* Attempt to allocate a specific subnet from the pool.  If the subnet is not
 * available, return an error.
 */
func (pool *IPv4SubnetPool) AllocSpecific(subnet *api.IPv4Subnet) error {
	// If the subnet is available without splitting anything, just remove it
	// from the list and return
	if pool[subnet.CIDRMask] != nil {
		node := pool[subnet.CIDRMask]
		// If the subnet is the first item in the list, removing it requires a
		// special case
		if node.subnet.Compare(subnet) == 0 {
			pool[subnet.CIDRMask] = node.next
			return nil
		} else {
			// If the subnet was not the first item, traverse the list until
			// it's found or there are no items remaining
			for prev, curr := node, node.next; curr != nil; prev, curr = curr, curr.next {
				if curr.subnet.Compare(subnet) == 0 {
					// If we found it, remove the subnet from the list (and let
					// go GC it)
					prev.next = curr.next
					return nil
				}
			}
		}
	}
	// Walk the pool until you find the subnet that contains the intended
	// subnet, then split it until the intended subnet is found.
	size := subnet.CIDRMask - 1
	var bigSubnet *api.IPv4Subnet
	for size >= 0 && bigSubnet != nil {
		if pool[size] != nil {
			if pool[size].subnet.Contains(subnet) {
				// If we found the containing subnet, remove it from the list
				bigSubnet = pool[size].subnet
				pool[size] = pool[size].next
			} else {
				for prev, curr := pool[size], pool[size].next; curr != nil; prev, curr = curr, curr.next {
					if curr.subnet.Contains(subnet) {
						// If we found the containing subnet, remove it from the list
						bigSubnet = curr.subnet
						prev.next = curr.next
						// Then stop traversing the list
						break
					}
				}
			}
		}
		size--
	}
	if bigSubnet != nil {
		// If we found the subnet during the previous loop, split it until we
		// get the exact subnet we're looking for (and return other subnets to
		// the pool along the way)
		for bigSubnet.Compare(subnet) != 0 && bigSubnet.CIDRMask < subnet.CIDRMask {
			loSubnet, hiSubnet, err := bigSubnet.Split()
			if err != nil {
				// If we hit an error, return the entire subnet to the pool,
				// then abort
				pool.Free(bigSubnet)
				return errors.New("Subnet " + subnet.String() +
					" not found in pool")
			}
			if loSubnet.Contains(subnet) {
				bigSubnet = loSubnet
				pool.Free(hiSubnet)
			} else {
				bigSubnet = hiSubnet
				pool.Free(loSubnet)
			}
		}
		if bigSubnet.Compare(subnet) == 0 {
			return nil
		}
	}
	return errors.New("Subnet " + subnet.String() + " not found in pool")
}

/* When freeing a subnet, first the pool should be checked for another subnet
 * with the same netmask that it can be merged with (e.g. 10.0.0.0/25 and
 * 10.0.0.128/25 can be merged into 10.0.0.0/24).  If a merge can be done, both
 * subnets should temporarily be allocated, the subnets merged, then the merged
 * subnet should be freed.
 *
 * I've had some issues with figuring out a fast way to check if they can be
 * merged, so for the current version, no merge checks are made.  In the
 * current implementation, we will always request a /24 subnet, so eventually
 * the entire pool will gravitate toward fragmenting at the /24 level.  Because
 * that's the size we care about, it shouldn't be an issue until the
 * implementation requires bigger subnets to be available.
 */
func (pool *IPv4SubnetPool) Free(subnet *api.IPv4Subnet) error {
	if subnet.CIDRMask < 0 || subnet.CIDRMask > 32 {
		return errors.New(fmt.Sprintf("Cannot free bad subnet %s", subnet))
	}
	var prev, curr *IPv4SubnetNode
	curr = pool[subnet.CIDRMask]
	// If there's nothing in the list, or the current subnet would sort before
	// the one at the beginning of this list, insert it first.
	if curr == nil || subnet.Compare(curr.subnet) < 0 {
		pool[subnet.CIDRMask] = &IPv4SubnetNode{subnet, curr}
		return nil
	}
	prev = curr
	curr = curr.next
	for curr != nil {
		switch {
		case subnet.Compare(curr.subnet) == 0:
			return errors.New(fmt.Sprintf("Double free of %s", subnet))
		case subnet.Compare(curr.subnet) < 0:
			prev.next = &IPv4SubnetNode{subnet, curr}
			return nil
		}
		prev = curr
		curr = curr.next
	}
	// We reached the end of the list (prev.next is nil), so add the subnet to
	// the end of it.
	prev.next = &IPv4SubnetNode{subnet, nil}
	return nil
}
