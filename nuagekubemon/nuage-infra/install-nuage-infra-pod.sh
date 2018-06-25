#!/bin/bash

### Go to idle state.
function go_to_idle_state () {
    echo "Nuage infra pod going to idle state..."
    while :
    do
        sleep 10;
    done
}

cleanup="false"
route_table_id="501"
veth1_name="nuage-infra-1"
veth2_name="nuage-infra-2"
router_pod_traffic_mark="0xabc/0xabc"
vport_resolve_config_file="/tmp/config.yaml"
vport_resolve_bin_file="/usr/bin/nuage-infra"

while [ "$1" != "" ]; do
    case $1 in
        -c | --cleanup )
            cleanup="true"
            ;;
         * )
            echo "invalid args $1 specified"
            exit
    esac
    shift
done

if [ "${PERSONALITY}" != "EVDF" || "${PERSONALITY}" != "evdf" ]
then
    if [ ${cleanup} == "true" ]
    then
        exit
    fi
    go_to_idle_state
fi

if [ X"${POD_NETWORK_CIDR}" == X ]
then
  echo "Pod Network CIDR not found in env.. Exiting.."
  exit
fi

if [ X"${VSP_USER}" == X ]
then
  echo "vspk user not set in env.. Exiting.."
  exit
fi

if [ X"${VSP_ENTERPRISE}" == X ]
then
  echo "vspk enterprise not set in env.. Exiting.."
  exit
fi

if [ X"${VSP_DOMAIN}" == X ]
then
  echo "vspk domain not set in env.. Exiting.."
  exit
fi

args=${vport_resolve_config_file}
if [ ${cleanup} == "true" ]
then
    args="${args} -cleanup"
fi

cat > ${vport_resolve_config_file} << EOF
---
name: "nuage-infra"
uuid: `uuidgen`
metadata:
    username: ${VSP_USER}
    enterprise: ${VSP_ENTERPRISE}
    domain: ${VSP_DOMAIN}
    zone: "default"
    subnet: "default-0"
    networktype: "ipv4"
interface:
    veth1: ${veth1_name}
    veth2: ${veth2_name}
EOF

gateway=`${vport_resolve_bin_file} ${args}`

if [ ${cleanup} == "true" ]
then
    ####################################################################
    ### Delete route table entries to redirect traffic through
    ####################################################################
    /usr/sbin/iptables -t mangle -D OUTPUT -d ${POD_NETWORK_CIDR} -j MARK --set-mark ${router_pod_traffic_mark} >& /dev/null
    /sbin/ip route del table ${route_table_id} ${gateway} dev ${veth2_name} >& /dev/null
    /sbin/ip route del table ${route_table_id} ${POD_NETWORK_CIDR} via ${gateway} dev ${veth2_name} >& /dev/null
    /sbin/ip rule del fwmark ${router_pod_traffic_mark} table ${route_table_id} >& /dev/null
    /usr/sbin/iptables -t nat -D POSTROUTING -o ${veth2_name} -j MASQUERADE >& /dev/null
    exit
fi

if [ X"${gateway}" == X ]
then
  echo "Error resolving infra vport. Exiting..."
  exit
fi

####################################################################
### Create route table entries to redirect traffic through
### new overlay vport
####################################################################
/usr/sbin/iptables -t mangle -A OUTPUT -d ${POD_NETWORK_CIDR} -j MARK --set-mark ${router_pod_traffic_mark} >& /dev/null
/sbin/ip rule add fwmark ${router_pod_traffic_mark} table ${route_table_id} >& /dev/null
/sbin/ip route add ${gateway} dev ${veth2_name} >& /dev/null
/sbin/ip route add ${POD_NETWORK_CIDR} via ${gateway} dev ${veth2_name} >& /dev/null
/sbin/ip route add table ${route_table_id} ${gateway} dev ${veth2_name} >& /dev/null
/sbin/ip route add table ${route_table_id} ${POD_NETWORK_CIDR} via ${gateway} dev ${veth2_name} >& /dev/null
/sbin/ip route del ${gateway} dev ${veth2_name} >& /dev/null
/sbin/ip route del ${POD_NETWORK_CIDR} via ${gateway} dev ${veth2_name} >& /dev/null
/sbin/ip route flush cache >& /dev/null
/usr/sbin/iptables -t nat -A POSTROUTING -o ${veth2_name} -j MASQUERADE >& /dev/null

####################################################################
### Disable reverse path filter check for this interface. This is 
### being done for other pat tap interfaces too
####################################################################
/usr/sbin/sysctl -w net.ipv4.conf.all.rp_filter=0 >& /dev/null
/usr/sbin/sysctl net.ipv4.conf.${veth1_name}.rp_filter=0 >& /dev/null
/usr/sbin/sysctl net.ipv4.conf.${veth2_name}.rp_filter=0 >& /dev/null

go_to_idle_state
