#
# spec file for package rokon
#
# Copyright (c) 2024 Brycen Granville <brycengranville@outlook.com>
#
# All modifications and additions to the file contributed by third parties
# remain the property of their copyright owners, unless otherwise agreed
# upon. The license for this file, and modifications and additions to the
# file, is the same license as for the pristine package itself (unless the
# license for the pristine package is not an Open Source License, in which
# case the license is the MIT License). An "Open Source License" is a
# license that conforms to the Open Source Definition (Version 1.9)
# published by the Open Source Initiative.

# Please submit bugfixes or comments via https://github.com/BrycensRanch/Rokon/issues as I am the developer

Name:           rokon
Version:        1.0.0
Release:        12%{?dist}
Summary:        Control your Roku device with your desktop!
License:        AGPL-3.0-or-later
URL:            https://github.com/BrycensRanch/Rokon
Source:         %{url}/archive/master.tar.gz

BuildRequires:  git
BuildRequires:  go
BuildRequires:  gcc
BuildRequires:  gcc-c++
BuildRequires:  gtk4-devel
BuildRequires:  gobject-introspection-devel
Requires:       gtk4

%description
Rokon is a GTK4 application that control your Roku.
Whether that be with your keyboard, mouse, or controller.

%prep
%autosetup -n Rokon-master

%build
go mod download all
ls
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
go build -v -ldflags="-X main.version=%{version} -X main.commit=$(git rev-parse --short HEAD) -X main.packaged=true -X main.packageFormat=rpm -X main.rpmRelease=%{rel} -X main.branch=$(git rev-parse --abbrev-ref HEAD) -X main.date=$(date -u +%Y-%m-%d)" -o %{name}
=======
make TARGET=%{name} PACKAGED=true PACKAGEFORMAT=rpm EXTRALDFLAGS="-s -w -X main.rpmRelease=%{rel}" EXTRAGOFLAGS="-trimpath" build
>>>>>>> 939117b (build: stop including debug symbols in native packages)

%install
install -Dpm 0755 ./%{name} %{buildroot}%{_bindir}/%{name}
install -Dpm 0644 ./usr/share/applications/io.github.brycensranch.Rokon.desktop %{buildroot}%{_datadir}/applications/io.github.brycensranch.Rokon.desktop
install -Dpm 0644 ./usr/share/icons/hicolor/48x48/apps/io.github.brycensranch.Rokon.png %{buildroot}%{_datadir}/icons/hicolor/48x48/apps/io.github.brycensranch.Rokon.png
install -Dpm 0644 ./usr/share/icons/hicolor/256x256/apps/io.github.brycensranch.Rokon.png %{buildroot}%{_datadir}/icons/hicolor/256x256/apps/io.github.brycensranch.Rokon.png
install -Dpm 0644 ./usr/share/icons/hicolor/scalable/apps/io.github.brycensranch.Rokon.svg %{buildroot}%{_datadir}/icons/hicolor/scalable/apps/io.github.brycensranch.Rokon.svg
install -Dpm 0644 ./usr/share/metainfo/io.github.brycensranch.Rokon.metainfo.xml %{buildroot}%{_datadir}/metainfo/io.github.brycensranch.Rokon.metainfo.xml
=======
make TARGET=%{name} PACKAGED=true PACKAGEFORMAT=rpm EXTRALDFLAGS="-X main.rpmRelease=%{rel}" EXTRAGOFLAGS="-trimpath" build

%install
<<<<<<< HEAD
make NODOCUMENTATION="1" PREFIX=%{buildroot}/usr install
>>>>>>> c08ac53 (refactor: fix build on opensuse & remove unnecessary comments)
=======
%if "%{?dist}" == "opensuse"
=======
# Rokon's Makefile still respects any CFLAGS LDFLAGS CXXFLAGS passed to it. It is compliant.
# https://lists.fedoraproject.org/archives/list/devel@lists.fedoraproject.org/thread/PK5PEKWE65UC5XQ6LTLSMATVPIISQKQS/
# Do not compress the DWARF debug information, it causes the build to fail!
# As of Go 1.11, debug information is compressed by default. We're disabling that.

%if 0%{?mageia}
    # Setup the correct compilation flags for the environment
    %set_build_flags
    # Fixes RPM build errors:
    # error: Empty %files file /builddir/build/BUILD/Rokon-master/debugsourcefiles.list
    # Empty %files file /builddir/build/BUILD/Rokon-master/debugsourcefiles.list
    %define _debugsource_template %{nil}
%endif

make TARGET=%{name} PACKAGED=true PACKAGEFORMAT=rpm EXTRALDFLAGS="-compressdwarf=false -X main.rpmRelease=%{rel}" EXTRAGOFLAGS="-buildmode=pie -trimpath" build

%install
%if 0%{?suse_version}
>>>>>>> 9b029b0 (build(spec): produce proper debug package)
    make NODOCUMENTATION="1" PREFIX=%{buildroot}/usr install
%else
    make NODOCUMENTATION="0" PREFIX=%{buildroot}/usr install
%endif
>>>>>>> 105b547 (build(spec): fix building on fedora)

%files
%{_bindir}/%{name}
%{_datadir}/applications/io.github.brycensranch.Rokon.desktop
%{_datadir}/icons/hicolor/48x48/apps/io.github.brycensranch.Rokon.png
%{_datadir}/icons/hicolor/256x256/apps/io.github.brycensranch.Rokon.png
%{_datadir}/icons/hicolor/scalable/apps/io.github.brycensranch.Rokon.svg
%{_datadir}/metainfo/io.github.brycensranch.Rokon.metainfo.xml
%license LICENSE.md
%doc *.md



%changelog
* Tue Sep 3 2024 Brycen <brycengranville@outlook.com> 1.0.0-6
- Removed sysinfo package decreasing binary size and portability and startup time
- Added metainfo file for appstream
- Added icons to package
- Added desktop entry
- Added license file to package
- Added documentation to package
* Mon Sep 2 2024 Brycen <brycengranville@outlook.com> 1.0.0-3
- Initial package
