%define nuage_openshift_monitor_binary	nuage-openshift-monitor
%define nuage_openshift_monitor_service    nuage-openshift-monitor
%define nuage_openshift_monitor_datadir	/usr/share/%{nuage_openshift_monitor_binary}
%define nuage_openshift_monitor_yaml	nuage-openshift-monitor.yaml
%define nuage_openshift_monitor_yaml_path	%{nuage_openshift_monitor_datadir}/%{nuage_openshift_monitor_yaml}
%define nuage_openshift_monitor_logdir     /var/log/%{nuage_openshift_monitor_binary}
%define nuage_openshift_monitor_init_script scripts/nuage-openshift-monitor.init

Name: nuage-openshift-monitor 
Version: 0.0
Release: 1%{?dist}
Summary: Nuage OpenShift Monitor	
Group: System Environments/Daemons	
License: ALU EULA and ASL 2.0	
Source0: nuage-openshift-monitor-%{version}.tar.gz

BuildRequires:  %{?go_compiler:compiler(go-compiler)}%{!?go_compiler:golang}
Requires: atomic-openshift-master

%description
%{summary}

%prep
%setup -q

%build

%pre
if [ "$1" = "2" ]; then
	cp $RPM_BUILD_ROOT%{nuage_openshift_monitor_yaml_path} $RPM_BUILD_ROOT%{nuage_openshift_monitor_yaml_path}.orig
fi

%install
install --directory $RPM_BUILD_ROOT%{nuage_openshift_monitor_datadir}
install --directory $RPM_BUILD_ROOT%{nuage_openshift_monitor_logdir}
install --directory $RPM_BUILD_ROOT/usr/bin
install --directory $RPM_BUILD_ROOT/etc/init.d

install -m 755 %{nuage_openshift_monitor_binary} $RPM_BUILD_ROOT/usr/bin
install -m 755 %{nuage_openshift_monitor_init_script} $RPM_BUILD_ROOT/etc/init.d/%{nuage_openshift_monitor_service}
install -m 644 %{nuage_openshift_monitor_yaml}.template  $RPM_BUILD_ROOT%{nuage_openshift_monitor_yaml_path}

%post
if [ "$1" = "2" ]; then
	mv $RPM_BUILD_ROOT%{nuage_openshift_monitor_yaml_path}.orig $RPM_BUILD_ROOT%{nuage_openshift_monitor_yaml_path}
fi
/sbin/chkconfig --add %{nuage_openshift_monitor_service}

%preun
if [ "$1" = "0" ]; then
      /sbin/service %{nuage_openshift_monitor_service} stop > /dev/null 2>&1
      /sbin/chkconfig --del %{nuage_openshift_monitor_service}
fi

%postun
if [ "$1" = "0" ]; then
   rm -rf $RPM_BUILD_ROOT%{nuage_openshift_monitor_datadir}
   rm -rf $RPM_BUILD_ROOT%{nuage_openshift_monitor_logdir}
   rm -f  $RPM_BUILD_ROOT/etc/init.d/%{nuage_openshift_monitor_service}
fi

%clean
rm -rf $RPM_BUILD_ROOT

%files

/usr/bin/%{nuage_openshift_monitor_binary}
/etc/init.d/%{nuage_openshift_monitor_service}
%attr(755, root, nobody) %{nuage_openshift_monitor_logdir} 
%attr(755, root, nobody) %{nuage_openshift_monitor_datadir}
%attr(644, root, nobody) %{nuage_openshift_monitor_yaml_path} 
%doc

%changelog
