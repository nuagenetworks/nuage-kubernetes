Name: nuagekubemon 
Version: 0.0
Release: 1%{?dist}
Summary: Nuage Kubernetes Monitor	

Group: System Environments/Daemons	
License: ALU EULA and ASL 2.0	
Source0: nuagekubemon-%{version}.tar.gz

BuildRequires:  %{?go_compiler:compiler(go-compiler)}%{!?go_compiler:golang}

%description
%{summary}

%prep
%setup -q


%build
godep go build
mv nuagekubemon-%{version} nuagekubemon

%install
install --directory $RPM_BUILD_ROOT/usr/share/nuagekubemon
install --directory $RPM_BUILD_ROOT/var/log/nuagekubemon
install --directory $RPM_BUILD_ROOT/usr/bin
install -m 755 nuagekubemon $RPM_BUILD_ROOT/usr/bin
install -m 644 nuagekubemon.yaml.template  $RPM_BUILD_ROOT/usr/share/nuagekubemon/nuagekubemon.yaml

%clean
rm -rf $RPM_BUILD_ROOT

%files

/usr/bin/nuagekubemon
%attr(755, root, nobody) /var/log/nuagekubemon
%attr(755, root, nobody) /usr/share/nuagekubemon
%attr(644, root, nobody) /usr/share/nuagekubemon/nuagekubemon.yaml

%doc

%changelog

