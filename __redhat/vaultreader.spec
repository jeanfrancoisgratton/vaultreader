%define debug_package   %{nil}
%define _build_id_links none
%define _name vaultreader
%define _prefix /opt
%define _version 2.3.0
%define _rel 2
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
Requires: bash-completion

%description
Hashicorp Vault client

%prep
%autosetup

%build
cd src
CGO_ENABLED=0 /opt/go/bin/go build -trimpath -ldflags="-s -w -buildid=" -o %{_builddir}/%{name}-%{version}/%{_binaryname} .

%clean
rm -rf $RPM_BUILD_ROOT

%pre

%install
rm -rf %{buildroot}
install -Dpm 0755 %{_builddir}/%{name}-%{version}/%{_binaryname} %{buildroot}%{_bindir}/%{_binaryname}

%post

%preun

%postun

%files
%defattr(-,root,root,-)
%attr(0755,root,root) %{_prefix}/bin/%{_binaryname}


%changelog
* Sat Aug 01 2026 Binary package builder <builder@famillegratton.net> 2.3.0-1
- Merge remote-tracking branch 'refs/remotes/origin/develop' into develop
- packaging revamp, no more hardcoded go path
- Errors are now sent to stderr instead of stdout
- buildDeps updates, gover update, APKBUILDER: workflow enhancements
- chore: update changelog for 2.2.0-0
- fixed Makefile indentation

* Fri Jun 26 2026 Binary package builder <builder@famillegratton.net> 2.2.0-0
- fixed Makefile indentation

* Wed Jun 03 2026 Binary package builder <builder@famillegratton.net> 2.01.00-0
- refreshed code against refactored vaultLib package, go version upgrade

* Wed May 13 2026 Binary package builder <builder@famillegratton.net> 2.00.01-0
- Removed useless Sprintf() call
- Removed no-longer needed dependency
- Commented out removed files
- further post/pre install scripts refactoring
- moved redhat support in its own dir, added archlinux support, go upgrade to 1.26.3
- Merge remote-tracking branch 'refs/remotes/origin/develop' into develop
- version output fix

* Sat Apr 25 2026 Binary package builder <builder@famillegratton.net> 2.00.00-1
- release number fix, doc update (jean-francois@famillegratton.net)

* Sat Apr 25 2026 Binary package builder <builder@famillegratton.net> 2.00.00-0
- Completed the migration to the library (jean-francois@famillegratton.net)
- removed logging, prepping to use vaultLib (jean-francois@famillegratton.net)

* Mon Jan 19 2026 Binary package builder <builder@famillegratton.net> 1.41.00-0
- Suppresed final exit glyph (jean-francois@famillegratton.net)
- Made the tool quieter for when running in scripts (jean-
  francois@famillegratton.net)

* Fri Nov 28 2025 Binary package builder <builder@famillegratton.net> 1.40.04-0
- Fixed buildscript, removed un-needed dirs (jean-francois@famillegratton.net)
- Release bump (RPM)
* Fri Nov 28 2025 Binary package builder <builder@famillegratton.net> 1.40.03-0
- builddep bump, helperFunctions migration to v4 (jean-
  francois@famillegratton.net)

* Fri Nov 14 2025 Binary package builder <builder@famillegratton.net> 1.40.02-0
- Fixed final output when o=json (jean-francois@famillegratton.net)
- fixed typo in APK post-install script (jean-francois@famillegratton.net)

* Tue Nov 11 2025 Binary package builder <builder@famillegratton.net> 1.40.01-0
- Fixed wrong logfile perms, fixed wrong newline output in final function call
  (jean-francois@famillegratton.net)

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



