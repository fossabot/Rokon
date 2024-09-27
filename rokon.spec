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

%global goipath github.com/BrycensRanch/Rokon

%global forgeurl https://github.com/BrycensRanch/Rokon

Name:           rokon
Version:        1.0.0
%if 0%{?fedora}
	Release:        %autorelease -p
%else
	Release:        13%{?dist}
%endif
Summary:        Control your Roku device with your desktop!
License:        AGPL-3.0-or-later
URL:            https://github.com/BrycensRanch/Rokon
%if 0%{?opensuse_bs}
	Source:         Rokon.tar.xz
%else
	Source:         %{url}/archive/master.tar.gz
%endif

%if 0%{?fedora}
%gometa -f
%endif

BuildRequires:  git
BuildRequires:  go
BuildRequires:  gcc
BuildRequires:  gcc-c++
BuildRequires:  gtk4-devel
BuildRequires:  gobject-introspection-devel
Requires:       gtk4
%if 0%{?opensuse_bs}
	# Logic specific to openSUSE Build Service. I imagine this will make it extremely difficult to build the spec locally on OBS.
	Source1:        vendor.tar.zst
	BuildRequires:  golang-packaging
	BuildRequires:  zstd
%endif

%generate_buildrequires
%if 0%{?fedora}
	%go_generate_buildrequires
%endif

%description
Rokon is a GTK4 application that control your Roku.
Whether that be with your keyboard, mouse, or controller.

%if 0%{?fedora}
	%goprep
%else
    %prep
%endif
%if 0%{?opensuse_bs}
	%autosetup -n Rokon
%else
	%autosetup -n Rokon-master
%endif

%build
# Setup the correct compilation flags for the environment
# Not all distributions do this automatically
%if 0%{?fedora}
    # Fedora specific behavior (no-op or something else)
    # Do nothing, since Fedora 33 the build flags are already set
%else
    %set_build_flags
%endif

ls
# Rokon's Makefile still respects any CFLAGS LDFLAGS CXXFLAGS passed to it. It is compliant.
# https://lists.fedoraproject.org/archives/list/devel@lists.fedoraproject.org/thread/PK5PEKWE65UC5XQ6LTLSMATVPIISQKQS/
# Do not compress the DWARF debug information, it causes the build to fail!
# As of Go 1.11, debug information is compressed by default. We're disabling that.

%if 0%{?mageia}
    # Fixes RPM build errors:
    # error: Empty %files file /builddir/build/BUILD/Rokon-master/debugsourcefiles.list
    # Empty %files file /builddir/build/BUILD/Rokon-master/debugsourcefiles.list
    %define _debugsource_template %{nil}
%endif

%define rpmRelease %{?dist}

%make_build \
    PACKAGED=true \
    PACKAGEFORMAT=rpm \
    EXTRALDFLAGS="-compressdwarf=false -X main.rpmRelease=%{rpmRelease}" \
    EXTRAGOFLAGS="-mod=vendor -buildmode=pie -trimpath"

%install
%if 0%{?suse_version}
    %make_install NODOCUMENTATION="1" PREFIX=%{buildroot}/usr
%else
    %make_install NODOCUMENTATION="0" PREFIX=%{buildroot}/usr
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

%if 0%{?fedora}
	%autochangelog
%else

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
%endif
