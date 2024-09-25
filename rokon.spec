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

# Please submit bugfixes or comments via http://bugs.opensuse.org/ or https://github.com/BrycensRanch/Rokon/issues

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
# Rokon's Makefile still respects any CFLAGS LDFLAGS CXXFLAGS passed to it. It is compliant.
# https://lists.fedoraproject.org/archives/list/devel@lists.fedoraproject.org/thread/PK5PEKWE65UC5XQ6LTLSMATVPIISQKQS/
# Do not compress the DWARF debug information, it causes the build to fail!
# As of Go 1.11, debug information is compressed by default. We're disabling that.

%if 0%{?mageia}
# Setup the correct compilation flags for the environment
%set_build_flags

# Use the %gobuild macro for Mageia
%gobuild \
    -o %{_bindir}/%{name} \
    -trimpath \
    -ldflags "-compressdwarf=false -X main.rpmRelease=%{rel} %{build_ldflags}" \
    %{name}

%else
# For other systems, use the regular make command with custom flags
make TARGET=%{name} \
     PACKAGED=true \
     PACKAGEFORMAT=rpm \
     EXTRALDFLAGS="-compressdwarf=false -X main.rpmRelease=%{rel}" \
     EXTRAGOFLAGS="-trimpath" \
     build
%endif

%install
%if 0%{?suse_version}
    make NODOCUMENTATION="1" PREFIX=%{buildroot}/usr install
%else
    make NODOCUMENTATION="0" PREFIX=%{buildroot}/usr install
%endif

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
