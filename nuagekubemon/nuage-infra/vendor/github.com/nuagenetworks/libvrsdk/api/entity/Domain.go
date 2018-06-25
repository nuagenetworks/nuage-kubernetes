package entity

// Domain the different virtualization platforms supported by Nuage
type Domain int

// Constants defining different platforms
const (
	KVM     Domain = 0
	ESXI    Domain = 1
	Xen     Domain = 2
	HyperV  Domain = 3
	Docker  Domain = 4
	LXC     Domain = 5
	ESXIM   Domain = 6
	QEMU    Domain = 7
	Unknown Domain = 8
)
