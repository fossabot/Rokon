Name:           rokon
Version:        1.0.0
Release:        8%{?dist}
Summary:        Control your Roku device with your desktop!

License:        AGPL-3.0-or-later
URL:            https://github.com/BrycensRanch/Rokon

%define _disable_source_fetch 0
Source:         %{url}/archive/master.tar.gz

BuildRequires:  git
BuildRequires:  go
BuildRequires:  gcc
BuildRequires:  gcc-c++
BuildRequires:  gtk4-devel
BuildRequires:  gobject-introspection-devel
Requires:       gtk4

Provides:       %{name} = %{version}-%{release}

%description
Rokon is a GTK4 application that allows you to control your Roku device with your desktop or controller!

%global debug_package %{nil}

%prep
%autosetup -n Rokon-master

%build
go mod download all
ls
go build -v -o %{name}

%install
tree .
install -Dpm 0755 ./%{name} %{buildroot}%{_bindir}/%{name}
install -Dpm 0644 ./usr/share/applications/io.github.brycensranch.Rokon.desktop %{buildroot}%{_datadir}/applications/io.github.brycensranch.Rokon.desktop
install -Dpm 0644 ./usr/share/icons/hicolor/48x48/apps/io.github.brycensranch.Rokon.png %{buildroot}%{_datadir}/icons/hicolor/48x48/apps/io.github.brycensranch.Rokon.png
install -Dpm 0644 ./usr/share/icons/hicolor/256x256/apps/io.github.brycensranch.Rokon.png %{buildroot}%{_datadir}/icons/hicolor/256x256/apps/io.github.brycensranch.Rokon.png
install -Dpm 0644 ./usr/share/icons/hicolor/scalable/apps/io.github.brycensranch.Rokon.svg %{buildroot}%{_datadir}/icons/hicolor/scalable/apps/io.github.brycensranch.Rokon.svg
install -Dpm 0644 ./usr/share/metainfo/io.github.brycensranch.Rokon.metainfo.xml %{buildroot}%{_datadir}/metainfo/io.github.brycensranch.Rokon.metainfo.xml

%files
%{_bindir}/%{name}
%{_datadir}/applications/io.github.brycensranch.Rokon.desktop
%{_datadir}/icons/hicolor/48x48/apps/io.github.brycensranch.Rokon.png
%{_datadir}/icons/hicolor/256x256/apps/io.github.brycensranch.Rokon.png
%{_datadir}/icons/hicolor/scalable/apps/io.github.brycensranch.Rokon.svg
%{_datadir}/metainfo/io.github.brycensranch.Rokon.metainfo.xml
%{buildroot}%{_docdir}/%{name}/
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
