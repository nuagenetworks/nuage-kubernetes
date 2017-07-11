# Setup etcd

  ::

    cp etcd_template.env etcd.env
    _Fill in node/cluster specific configuration_
    vim etcd.env
    chmod +x setup-etcd.sh
    ./setup-etcd.sh

# Setup base for master servers
cp master_template.env master.env
_Fill in node/cluster specific configuration_
vim master.env
chmod +x setup-master.sh
./setup-master.sh

# Initialize first master server
_only needed on first master server and on first configuration_
chmod +x init-master.sh
./init-master.sh

# Configure other master servers
_copy needed files from initial master_
_copy /etc/kubernetes/*_
