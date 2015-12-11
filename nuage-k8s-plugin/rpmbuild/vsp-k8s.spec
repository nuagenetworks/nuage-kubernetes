Name: nuage-k8s-plugin 
Version: 0.0
Release: 1%{?dist}
Summary: Nuage Kubernetes Plugin 

Group: System Environments/Daemons	
License: ALU EULA and ASL 2.0	
Source0: nuage-k8s-plugin-%{version}.tar.gz

Requires: nuage-openvswitch, bridge-utils, python-yaml, python-requests

%description
%{summary}

%build

%prep
%setup -q

%install
install --directory $RPM_BUILD_ROOT/usr/libexec/kubernetes/kubelet-plugins/net/exec/nuage~vsp-k8s
install --directory $RPM_BUILD_ROOT/usr/share/vsp-k8s
install -m 755 vsp-k8s $RPM_BUILD_ROOT/usr/libexec/kubernetes/kubelet-plugins/net/exec/nuage~vsp-k8s
install -m 644 vsp-k8s.yaml.template  $RPM_BUILD_ROOT/usr/share/vsp-k8s/vsp-k8s.yaml

%post 

test -e /usr/share/openvswitch/scripts/vrs-platform-lib || exit 0
. /usr/share/openvswitch/scripts/vrs-platform-lib
add_platform k8s 

%preun

if [ "$1" = "0" ]; then     # $1 = 0 for uninstall
    test -e /usr/share/openvswitch/scripts/vrs-platform-lib || exit 0
    . /usr/share/openvswitch/scripts/vrs-platform-lib
    remove_platform k8s 
fi

%clean
rm -rf $RPM_BUILD_ROOT

%files
/usr/share/vsp-k8s
/usr/libexec/kubernetes/kubelet-plugins/net/exec/nuage~vsp-k8s/vsp-k8s
%attr(644, root, nobody) /usr/share/vsp-k8s/vsp-k8s.yaml

%doc

%changelog

