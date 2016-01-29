%define nuage-openshift-monitor_binary	nuage-openshift-monitor
%define nuage-openshift-monitor_service    nuage-openshift-monitor
%define nuage-openshift-monitor_datadir	/usr/share/%{nuage-openshift-monitor_binary}
%define nuage-openshift-monitor_yaml	nuage-openshift-monitor.yaml
%define nuage-openshift-monitor_yaml_path	%{nuage-openshift-monitor_datadir}/%{nuage-openshift-monitor_yaml}
%define nuage-openshift-monitor_logdir     /var/log/%{nuage-openshift-monitor_binary}
%define nuage-openshift-monitor_init_script scripts/nuage-openshift-monitor.init

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
	cp $RPM_BUILD_ROOT%{nuage-openshift-monitor_yaml_path} $RPM_BUILD_ROOT%{nuage-openshift-monitor_yaml_path}.orig
fi

%install
install --directory $RPM_BUILD_ROOT%{nuage-openshift-monitor_datadir}
install --directory $RPM_BUILD_ROOT%{nuage-openshift-monitor_logdir}
install --directory $RPM_BUILD_ROOT/usr/bin
install --directory $RPM_BUILD_ROOT/etc/init.d

install -m 755 %{nuage-openshift-monitor_binary} $RPM_BUILD_ROOT/usr/bin
install -m 755 %{nuage-openshift-monitor_init_script} $RPM_BUILD_ROOT/etc/init.d/%{nuage-openshift-monitor_service}
install -m 644 %{nuage-openshift-monitor_yaml}.template  $RPM_BUILD_ROOT%{nuage-openshift-monitor_yaml_path}

%post
if [ "$1" = "2" ]; then
	mv $RPM_BUILD_ROOT%{nuage-openshift-monitor_yaml_path}.orig $RPM_BUILD_ROOT%{nuage-openshift-monitor_yaml_path}
fi
/sbin/chkconfig --add %{nuage-openshift-monitor_service}

%preun
if [ "$1" = "0" ]; then
      /sbin/service %{nuage-openshift-monitor_service} stop > /dev/null 2>&1
      /sbin/chkconfig --del %{nuage-openshift-monitor_service}
fi

%postun
if [ "$1" = "0" ]; then
   rm -rf $RPM_BUILD_ROOT%{nuage-openshift-monitor_datadir}
   rm -rf $RPM_BUILD_ROOT%{nuage-openshift-monitor_logdir}
   rm -f  $RPM_BUILD_ROOT/etc/init.d/%{nuage-openshift-monitor_service}
fi

%clean
rm -rf $RPM_BUILD_ROOT

%files

/usr/bin/%{nuage-openshift-monitor_binary}
/etc/init.d/%{nuage-openshift-monitor_service}
%attr(755, root, nobody) %{nuage-openshift-monitor_logdir} 
%attr(755, root, nobody) %{nuage-openshift-monitor_datadir}
%attr(644, root, nobody) %{nuage-openshift-monitor_yaml_path} 
%doc

%changelog
