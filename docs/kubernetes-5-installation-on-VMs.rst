.. _Kubernetes-5-installation-on-VMs:

.. include:: ../lib/doc-includes/VSDA-icons.inc

============================================
Kubernetes Installation on Virtual Machines
============================================

.. contents::
   :local:
   :depth: 3
   

Overview
==========

There are deployment scenarios where Kubernetes is run on VMs. For example, there can be an IaaS environment (like OpenStack) or simply a virtual environment where all the hosts are in the form of VMs. This gives the flexibility of deploying of hybrid application environments with both VMs and Containers/Pods.   

The Linux VMs which are used as Kubernetes Master or Nodes can be deployed using OpenStack or deployed manually.

.. Note:: The Linux VMs need to be deployed in a separate VSD Domain/Zone/Subnet to avoid any conflict with the entities created by Nuage Kubernetes monitor.

Nuage Ansible installer can be used to install Nuage VRS packages on the VMs. You must make sure that the VRS VMs are able to communicate with the VSC (The VSC can be the same as the one used for the VRS on the bare metal hypervisors). In this case the Kubernetes VMs are resolved by the Nuage VRS running on the bare metal hypervisor resulting in double overlay deployments.

.. Note:: For the use case where we want  Pod to VM communication the VMs need to be in the same Domain as the Kubernetes Pods. 

Installing Kubernetes Enterprise on VMs
========================================

To install Kubernetes Enterprise on the VMs, refer to the `Kubernetes Installation <Kubernetes-4-installation.html#>`_ section of this document. Make sure that the Kubernetes master(s) can reach all nodes deployed as VMs.

Installing Kubernetes Enterprise on Nuage-Resolved VMs
========================================================

While deploying pods on Nuage-resolved VMs running on bare metal hypervisors, the data traffic from the pods are encapsulated twice. Therefore, the MTU on the pod must be decreased to accommodate multiple VXLAN headers. For example, if the pod MTU is set to 1460, the VM MTU should be 1500, and the hypervisor MTU should be 1540.

In order to set pod MTU using Ansible, refer to the configuration file mentioned in the `Ansible Installation <Kubernetes-4-installation.html#ansible-installation#>`_ section of this document. 


Pod to VM Communication within a Nuage Domain
----------------------------------------------

This workflow is used for connectivity between pods in a namespace and other VMs resolved in the Nuage domain.


.. figure:: graphics/pod-to-vm.png
   
   
   Pod to VM Communication


Prerequisites and Assumptions
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

* Kubernetes pods and VMs are resolved in the same Nuage domain
* VMs and pods can reside on same or different compute resource


Pods to VM communication Workflow
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

1. Resolve the already created VMs using bottom-up or split activation workflow.
2. Go to the `Detailed Kubernetes Workflows <Kubernetes-3-workflows.html>`_ section and follow the steps in the Developer or Operations workflow to resolve Kubernetes pods. 
3. Make sure that the VM Zones and Pod Zones that need to communicate with each other are put in the same VSP Domain and appropriate ACLs are needed to enable inter-zone communications.
4. Set the rp_filter flag to 2 on the interface that the hypervisor uses to reach the other hypervisors because the outgoing and incoming routes are different for these packets.

