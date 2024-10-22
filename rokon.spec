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


%bcond check 0

# https://github.com/BrycensRanch/Rokon
%global goipath         github.com/brycensranch/rokon
%global forgeurl        https://github.com/BrycensRanch/Rokon
%global commit          d7945dc774573152b80125d74105c233fe5fdc45


%if 0%{?fedora}
%gometa -L -f
%endif


Name:           rokon
Version:        1.0.0
%if 0%{?fedora}
Release:        %autorelease -p
%else
Release:        17%{?dist}
%endif
Summary:        Control your Roku device with your desktop!


License:        AGPL-3.0-only AND BSD-3-Clause AND CC-BY-SA-4.0 AND ISC AND MIT AND MPL-2.0
URL:            https://github.com/BrycensRanch/Rokon
Source:         https://nightly.link/BrycensRanch/Rokon/workflows/publish/master/rokon-vendored-source.zip


BuildRequires:  go
%if 0%{?fedora}
BuildRequires:  go-vendor-tools
%endif

BuildRequires:  gcc
BuildRequires:  unzip
BuildRequires:  gcc-c++
BuildRequires:  gtk4-devel
BuildRequires:  gobject-introspection-devel
Requires:       gtk4

%description
Rokon is a GTK4 application that control your Roku.
Whether that be with your keyboard, mouse, or controller.

%prep
ls

%setup -c
ls

%generate_buildrequires
%if 0%{?fedora}
%go_vendor_license_buildrequires -c go-vendor-tools.toml
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
    BUILDTAGS="rpm_crashtraceback" \
    PACKAGEFORMAT=rpm \
    EXTRALDFLAGS="-compressdwarf=false -X main.rpmRelease=%{rpmRelease}" \
    EXTRAGOFLAGS="-x -mod=vendor -buildmode=pie -trimpath"

%install
%if 0%{?fedora}
%go_vendor_license_install -c go-vendor-tools.toml
%endif

# Why was this necessary, you ask?!?!  Because I kept getting
# File not found: /builddir/build/BUILDROOT/rokon-1.0.0-16.suse.tw.x86_64/usr/share/doc/packages/rokon


%if 0%{?suse_version}
%make_install PREFIX=%{_prefix} \
              DOCDIR="%buildroot/%_docdir/%name"
%else
%make_install PREFIX=%{_prefix} \
              DOCDIR=%{_docdir}
%endif

%check
%if 0%{?fedora}
%if %{with check}
%go_vendor_license_check -c go-vendor-tools.toml
%gocheck
%endif
%endif

%if 0%{?fedora}
%files -f %{go_vendor_license_filelist}
%else
%files
%endif
%{_bindir}/%{name}
%if 0%{?suse_version}
%{_docdir}/%name
%endif
%{_datadir}/applications/io.github.brycensranch.Rokon.desktop
%{_datadir}/metainfo/io.github.brycensranch.Rokon.metainfo.xml
%{_datadir}/dbus-1/services/io.github.brycensranch.Rokon.service
%{_datadir}/icons/hicolor/48x48/apps/io.github.brycensranch.Rokon.png
%{_datadir}/icons/hicolor/128x128/apps/io.github.brycensranch.Rokon.png
%{_datadir}/icons/hicolor/256x256/apps/io.github.brycensranch.Rokon.png
%{_datadir}/icons/hicolor/scalable/apps/io.github.brycensranch.Rokon.svg
%if 0%{?fedora}
%else
%license vendor/modules.txt LICENSE.md
%endif

# https://en.opensuse.org/openSUSE:Packaging_Conventions_RPM_Macros#%_docdir
%if 0%{?suse_version}
%else
%doc *.md
%endif


%if 0%{?fedora}
%changelog
%autochangelog
%else

%changelog
* Tue Sep 3 2024 Brycen <brycengranville@outlook.com> 1.0.0-6
- Removed sysinfo package decreasing binary size and portability and startup time
- Added metainfo file for appstream
-Added icons to package
- Added desktop entry
- Added license file to package
- Added documentation to package
* Mon Sep 2 2024 Brycen <brycengranville@outlook.com> 1.0.0-3
- Initial package
%endif

