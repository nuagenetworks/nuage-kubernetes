#!/bin/sh

NET_CONF=""
MONITOR_CONF=""
RUNAS=root
if [ "$1" = "k8s" ]; then
    BIN=/usr/bin/nuagekubemon
    PIDFILE=/var/run/nuagekubemon.pid
    rm -irf /usr/bin/nuage-openshift-monitor
    MONITOR=nuagekubemon
    MONITOR_CONF='/usr/share/nuagekubemon/nuagekubemon.yaml'
    NET_CONF='/usr/share/nuagekubemon/net-config.yaml'
fi

if [ "$1" = "ose" ]; then
    BIN=/usr/bin/nuage-openshift-monitor
    PIDFILE=/var/run/nuage-openshift-monitor.pid
    rm -irf /usr/bin/nuagekubemon
    MONITOR=nuage-openshift-monitor
    MONITOR_CONF='/usr/share/nuage-openshift-monitor/nuage-openshift-monitor.yaml'
fi

if [ "$2" = "is_atomic" ]; then
    DIR='/var/usr/share'
    MONITOR_CONF=$DIR/$MONITOR/$MONITOR.yaml
fi

# Configuring Nuage monitor yaml file on master nodes using
# daemon sets
cat > $MONITOR_CONF << EOF
EOF
chmod 777 $MONITOR_CONF

if [ "${NUAGE_MASTER_VSP_CONFIG:-}" != "" ]; then
cat >$MONITOR_CONF <<EOF
${NUAGE_MASTER_VSP_CONFIG:-}
EOF
fi

# Configuring Nuage master netconfig yaml file on master nodes using
# daemon sets
if [ "${NUAGE_MASTER_NETWORK_CONFIG:-}" != "" ]; then
cat > $NET_CONF << EOF
EOF
chmod 777 $NET_CONF

cat >$NET_CONF <<EOF
${NUAGE_MASTER_NETWORK_CONFIG:-}
EOF
fi

# Starting Nuage monitor on master nodes
CMD="$BIN --config=$MONITOR_CONF"
su -c "$CMD" $RUNAS > "$PIDFILE"
