FROM centos:7

ADD configure-bgp.sh /configure-bgp.sh
ADD nuage-*.rpm /tmp/
RUN yum -y install /tmp/nuage-*.rpm && rm -f /tmp/nuage-*.rpm

CMD ["/configure-bgp.sh"]
