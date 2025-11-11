%define debug_package   %{nil}
%define _build_id_links none
%define _name vaultreader
%define _prefix /opt
%define _version 1.40.00
%define _rel 0
%define _arch x86_64
%define _binaryname vaultreader

Name:       vaultreader
Version:    %{_version}
Release:    %{_rel}
Summary:    Lightweight Hashicorp Vault secret reader

Group:      CI/CD utils
License:    GPL2.0
URL:        https://git.famillegratton.net:3000/devops/vaultreader

Source0:    %{name}-%{_version}.tar.gz
BuildArchitectures: x86_64

%description
Hashicorp Vault client

%prep
%autosetup

%build
cd %{_sourcedir}/%{_name}-%{_version}/src
PATH=$PATH:/opt/go/bin CGO_ENABLED=0 go build -trimpath -ldflags="-s -w -buildid=" -o %{_sourcedir}/%{_binaryname} .

%clean
rm -rf $RPM_BUILD_ROOT

%pre
if getent group vaultreader > /dev/null; then
  exit 0
else
  if getent group 3000 > /dev/null; then
    groupadd vaultreader
  else
    groupadd -g 3000 vaultreader
  fi
fi

%install
install -Dpm 2755 %{_sourcedir}/%{_binaryname} %{buildroot}%{_bindir}/%{_binaryname}

%post
BIN="%{_prefix}/bin/%{_binaryname}"
LOG="/var/log/%{_name}.log"
install -m 0644 -o root -g vaultreader /dev/null "$LOG" || :

%preun

%postun

%files
%defattr(-,root,root,-)
%attr(2755,root,vaultreader) %{_prefix}/bin/%{_binaryname}


%changelog
* Tue Nov 11 2025 Binary package builder <builder@famillegratton.net> 1.40.00-0
- Fixed error handling and wrong env variable (jean-
  francois@famillegratton.net)
- Fixes in APK scripts (jean-francois@famillegratton.net)

* Mon Nov 10 2025 Binary package builder <builder@famillegratton.net> 1.30.00-0
- Automatic commit of package [vaultreader] release [1.30.00-0].
  (builder@famillegratton.net)
- Automatic commit of package [vaultreader] release [1.30.00-0].
  (builder@famillegratton.net)
- Automatic commit of package [vaultreader] release [1.30.00-0].
  (builder@famillegratton.net)
- Automatic commit of package [vaultreader] release [1.30.00-0].
  (builder@famillegratton.net)
- Automatic commit of package [vaultreader] release [1.30.00-0].
  (builder@famillegratton.net)
- Fixed debian script (builder@famillegratton.net)
- other round of buildscript fix (jean-francois@famillegratton.net)
- build scripts fixes (jean-francois@famillegratton.net)
- Completed verbosity feature addition (jean-francois@famillegratton.net)
- Updated packaging scripts to accomodate the new logfile path (jean-
  francois@famillegratton.net)
- migrated my helper and error packages to v3 (jean-
  francois@famillegratton.net)
- interim sync (jean-francois@famillegratton.net)
- post-install script cleanup (jean-francois@famillegratton.net)
- New way of packaging APK (jean-francois@famillegratton.net)
- Forgot to comment out dependency (jean-francois@famillegratton.net)
- Updated link options for Alpine (jean-francois@famillegratton.net)
- builddedps update (builder@famillegratton.net)

* Sun Nov 02 2025 Binary package builder <builder@famillegratton.net> 1.23.00-1
- Fixed branch drift (jean-francois@famillegratton.net)

* Sun Nov 02 2025 Binary package builder <builder@famillegratton.net> 1.23.00-0
- Now trimming the binary at build stage, code completion, go version bump
  (jean-francois@famillegratton.net)
- admin stub -- take 2 (jean-francois@famillegratton.net)
- admin subcommands stub (jean-francois@famillegratton.net)


* Fri Jul 25 2025 Binary package builder <builder@famillegratton.net> 1.22.00-1
- Package version bump (jean-francois@famillegratton.net)
- Moved to a new logging facility, no more need using root (jean-
  francois@famillegratton.net)

* Tue Jul 01 2025 Binary package builder <builder@famillegratton.net> 1.21.00-1
- release bump (builder@famillegratton.net)
- disabled CGO when building packages (jean-francois@famillegratton.net)

* Tue Jul 01 2025 Binary package builder <builder@famillegratton.net> 1.21.00-0
- Version bump (jean-francois@famillegratton.net)
- added more error handling (jean-francois@famillegratton.net)
- fixed perm on rpmbuild-deps.sh (builder@famillegratton.net)

* Tue Jul 01 2025 Binary package builder <builder@famillegratton.net> 1.20.00-0
- Added a error message when a secret path is wrong (jean-
  francois@famillegratton.net)
- GO version bump (jean-francois@famillegratton.net)
- Added a dependency to APK package (jean-francois@famillegratton.net)

* Thu Jun 05 2025 APK Builder <builder@famillegratton.net> 1.10.01-0
- Environment var fix (jean-francois@famillegratton.net)

* Thu Jun 05 2025 APK Builder <builder@famillegratton.net> 1.01.00-0
- Completed logging capabilities (jean-francois@famillegratton.net)
- Doc and version update (jean-francois@famillegratton.net)
- formatting output to json (jean-francois@famillegratton.net)
- refactored output (jean-francois@famillegratton.net)
- preparing to manage output (jean-francois@famillegratton.net)
- Version string cosmetic change (jean-francois@famillegratton.net)
- post-install fixes (builder@famillegratton.net)



