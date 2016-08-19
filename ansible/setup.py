#!/usr/bin/env python

import os

def run_cmd(cmd):
	print cmd
	os.system(cmd)

stages = ["etcd", "mastes", "nodes", "network-service-install"]
map(lambda x: run_cmd("ansible-playbook -vvvv -i nodes cluster.yml --tags=%s" % x), stages)
