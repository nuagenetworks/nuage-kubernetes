FROM centos:7

RUN yum install -y iptables

ADD vrs-resolve-cport /usr/bin/vrs-resolve-cport
ADD install-nuage-infra-pod.sh /usr/bin/install-nuage-infra-pod.sh

CMD ["/install-nuage-infra-pod.sh"]
