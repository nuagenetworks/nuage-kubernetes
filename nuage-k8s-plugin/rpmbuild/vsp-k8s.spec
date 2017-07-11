%define vsp_k8s_plugin vsp-k8s
%define vsp_k8s_plugin_dir /usr/libexec/kubernetes/kubelet-plugins/net/exec/%{vsp_k8s_plugin}
%define vsp_k8s_datadir /usr/share/%{vsp_k8s_plugin} 
%define vsp_k8s_yaml %{vsp_k8s_plugin}.yaml
%define vsp_k8s_yaml_path %{vsp_k8s_datadir}/%{vsp_k8s_yaml}
%define nuage_vrs_platform_script /usr/share/openvswitch/scripts/vrs-platform-lib 

Name: nuage-k8s-plugin 
Version: 0.0
Release: 1%{?dist}
Summary: Nuage Kubernetes Plugin 

Group: System Environments/Daemons	
License: ALU EULA and ASL 2.0	
Source0: nuage-k8s-plugin-%{version}.tar.gz

Requires: nuage-openvswitch, bridge-utils, python-yaml, python-requests >= 2.3.0

%description
%{summary}

%prep
%setup -q

%build

%pre
if [ "$1" = "2" ]; then
	cp $RPM_BUILD_ROOT%{vsp_k8s_yaml_path} $RPM_BUILD_ROOT%{vsp_k8s_yaml_path}.orig
	rm -f $RPM_BUILD_ROOT/%{vsp_k8s_plugin_dir}/%{vsp_k8s_plugin}
fi

%install
install --directory $RPM_BUILD_ROOT%{vsp_k8s_plugin_dir}
install --directory $RPM_BUILD_ROOT%{vsp_k8s_datadir}
install -m 755 %{vsp_k8s_plugin} $RPM_BUILD_ROOT%{vsp_k8s_plugin_dir}
install -m 644 %{vsp_k8s_yaml}.template  $RPM_BUILD_ROOT%{vsp_k8s_yaml_path}

%post 

if [ "$1" = "1" ]; then # first time install only.
test -e %{nuage_vrs_platform_script} || exit 0
. %{nuage_vrs_platform_script}
add_platform k8s 
fi

if [ "$1" = "2" ]; then
	cp $RPM_BUILD_ROOT%{vsp_k8s_yaml_path}.orig $RPM_BUILD_ROOT%{vsp_k8s_yaml_path}
fi

%preun

if [ "$1" = "0" ]; then     # $1 = 0 for uninstall
    test -e %{nuage_vrs_platform_script} || exit 0
    . %{nuage_vrs_platform_script}
    remove_platform k8s 
fi

%postun
if [ "$1" = "0" ]; then
   rm -rf $RPM_BUILD_ROOT%{vsp_k8s_datadir}
   rm -rf $RPM_BUILD_ROOT%{vsp_k8s_plugin_dir}
fi

%clean
rm -rf $RPM_BUILD_ROOT

%files
%{vsp_k8s_datadir}
%{vsp_k8s_plugin_dir}
%{vsp_k8s_plugin_dir}/%{vsp_k8s_plugin}
%attr(644, root, nobody) %{vsp_k8s_yaml_path}

%doc

%changelog

