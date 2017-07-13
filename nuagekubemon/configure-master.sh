#!/bin/sh

RUNAS=root
if [ "$1" = "k8s" ]; then
    BIN=/usr/bin/nuagekubemon
    PIDFILE=/var/run/nuagekubemon.pid
    rm -irf /usr/bin/nuage-openshift-monitor
    MONITOR=nuagekubemon
    TMP_CONF='/tmp/nuagekubemon.yaml'
    MONITOR_CONF='/usr/share/nuagekubemon/nuagekubemon.yaml'
fi

if [ "$1" = "ose" ]; then
    BIN=/usr/bin/nuage-openshift-monitor
    PIDFILE=/var/run/nuage-openshift-monitor.pid
    rm -irf /usr/bin/nuagekubemon
    MONITOR=nuage-openshift-monitor
    TMP_CONF='/tmp/nuage-openshift-monitor.yaml'
    MONITOR_CONF='/usr/share/nuage-openshift-monitor/nuage-openshift-monitor.yaml'
fi

if [ "$2" = "is_atomic" ]; then
    DIR='/var/usr/share'
    MONITOR_CONF=$DIR/$MONITOR/$MONITOR.yaml
fi

# Configuring Nuage monitor yaml file on master nodes using
# daemon sets
cat > $TMP_CONF << EOF
EOF
chmod 777 $TMP_CONF

if [ "${NUAGE_MASTER_VSP_CONFIG:-}" != "" ]; then
cat >$TMP_CONF <<EOF
${NUAGE_MASTER_VSP_CONFIG:-}
EOF
fi

mv $TMP_CONF $MONITOR_CONF

# Starting Nuage monitor on master nodes
CMD="$BIN --config=$MONITOR_CONF"
su -c "$CMD" $RUNAS > "$PIDFILE"
