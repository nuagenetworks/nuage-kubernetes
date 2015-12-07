Name: Nuage Kubernetes Monitor		
Version: 0.1
Release: 1%{?dist}
Summary: Nuage Kubernetes Monitor	

Group: System Environments/Daemons	
License: ALU EULA and ASL 2.0	
URL:
Source0:	

BuildRequires:  %{?go_compiler:compiler(go-compiler)}%{!?go_compiler:golang}
Requires: bridge-utils	

%description
Nuage 

%prep
%setup -q


%build
godep go build

%install
install --directory $RPM_BUILD_ROOT/usr/share/nuagekubemon
install --directory $RPM_BUILD_ROOT/usr/bin
install -m 755 nuagekubemon $RPM_BUILD_ROOT/usr/bin

%clean
rm -rf $RPM_BUILD_ROOT

%files

/usr/bin/nuagekubemon
%attr(755, root, nobody) /usr/share/nuagekubemon

%doc

%changelog

