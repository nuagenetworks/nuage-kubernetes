#!/bin/bash
# Parse CLI options
for i in "$@"; do
    case $i in
        --ca=*)
            CA_CERT="${i#*=}"
        ;;
        --server=*)
            SERVER="${i#*=}"
        ;;
        --output=*)
            CONFIG_FILE="${i#*=}"
        ;;
        --adminconfig=*)
            ADMIN_FILE="${i#*=}"
        ;;
    esac
done

# If any are missing, print the usage and exit
if [ -z $CA_CERT ] || [ -z $SERVER ] || [ -z $CONFIG_FILE ] || [ -z $ADMIN_FILE ]; then
	echo "Invalid syntax: $@"
	echo "Usage:"
	echo "  $0 --ca=/path/to/ca.crt --server=<address>:<port> --output=/path/to/nuage.kubeconfig --adminconfig=/path/to/admin.kubeconfig"
	echo "--ca:     Certificate for the OpenShift CA"
	echo "--server: Address of Kubernetes API server (default port is 8443)"
	echo "--output: File to save configuration to"
	echo "--adminconfig: kubeconfig with system:admin, used to create the service account and assign necessary privileges"
	echo ""
	echo "All options are required"
	exit 1
fi

# Login as admin so that we can create the service account
oc login -u system:admin --config=$ADMIN_FILE || exit 1
oc project default --config=$ADMIN_FILE

ACCOUNT_CONFIG='
{
  "apiVersion": "v1",
  "kind": "ServiceAccount",
  "metadata": {
    "name": "nuage"
  }
}
'

# Create the account with the included info
echo $ACCOUNT_CONFIG|oc create --config=$ADMIN_FILE -f -

# Get one of the account tokens
TOKEN_NAME=$(oc describe serviceaccount nuage --config=$ADMIN_FILE|awk '/^Tokens:/{ print $2 }'|head -n 1)
TOKEN=$(oc describe secret $TOKEN_NAME --config=$ADMIN_FILE|awk '/^token:/{ print $2 }'|head -n 1)

# Add the cluser-reader role, which allows this service account read access to
# everything in the cluster except secrets
oadm policy add-cluster-role-to-user cluster-reader system:serviceaccount:default:nuage --config=$ADMIN_FILE

# Login as the user to create the config file (kubeconfig data is automatically
# written to the file specified with --config)
oc login --api-version='v1' --certificate-authority=$CA_CERT --server=$SERVER --token=$TOKEN --config=$CONFIG_FILE

# Verify that we logged in correctly
if ! [ $(oc whoami --config=$CONFIG_FILE) == 'system:serviceaccount:default:nuage' ]; then
	echo "Service account creation failed!"
	exit 1
fi

# The kubeconfig file references the file specified in $CA_CERT by default, but
# we want the cert string in the kubeconfig so that it isn't necessary to keep
# the $CA_CERT file around, or keep the kubeconfig on the machine it was
# generated on.  `oc config view --flatten` will fix that
oc config view --flatten --config=$CONFIG_FILE > ${CONFIG_FILE}.tmp
mv ${CONFIG_FILE}.tmp $CONFIG_FILE

# Reverify the finalized kubeconfig
if ! [ $(oc whoami --config=$CONFIG_FILE) == 'system:serviceaccount:default:nuage' ]; then
	echo "Service account creation failed!"
	exit 1
fi
