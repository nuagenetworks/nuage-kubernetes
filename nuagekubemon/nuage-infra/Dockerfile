FROM centos:centos7

RUN yum install -y iptables iproute

ADD nuage-infra /usr/bin/nuage-infra
ADD install-nuage-infra-pod.sh /usr/bin/install-nuage-infra-pod.sh

CMD ["/install-nuage-infra-pod.sh"]
