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
	"flag"

	"github.com/nuagenetworks/nuage-kubernetes/nuagekubemon/monitor"
)

func main() {
	mNuageKubeMon := monitor.NewNuageKubeMonitor()
	flagSet := flag.CommandLine
	mNuageKubeMon.ParseArgs(flagSet) //define flags.
	flag.Parse()                     //parse the flags.
	mNuageKubeMon.Run()              //start the monitor.
}
