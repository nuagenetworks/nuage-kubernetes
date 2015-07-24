/*
###########################################################################
#
#   Filename:           main.go
#
#   Author:             Aniket Bhat
#   Created:            July 20, 2015
#
#   Description:        Main For NuageKubeMon
#
###########################################################################
#
#              Copyright (c) 2015 Nuage Networks
#
###########################################################################
*/

/* package main is the main entry point for nuagekubemon */

package main

import (
	"github.com/nuagenetworks/openshift-integration/nuagekubemon/monitor"
	flag "github.com/spf13/pflag"
)

func main() {
	mNuageKubeMon := monitor.NewNuageKubeMonitor()
	var flagSet = flag.CommandLine
	mNuageKubeMon.ParseArgs(flagSet)  //define flags.
	flag.Parse()                      //parse the flags.
	mNuageKubeMon.SetLogging(flagSet) //setup logging framework
	mNuageKubeMon.Run()               //start the monitor.
}
