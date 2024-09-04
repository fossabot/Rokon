Name:           rokon
Version:        1.0.0
Release:        5%{?dist}
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
BuildRequires:  glib2-devel
BuildRequires:  gobject-introspection-devel

Provides:       %{name} = %{version}-%{release}

%description
Rokon is a GTK4 application that allows you to control your Roku device with your desktop. It is written in Golang and uses the Roku External Control API to communicate with your Roku device.

%global debug_package %{nil}

%prep
%autosetup -n Rokon-master



%build
go mod download all
ls
go build -v -o %{name}

%install
install -Dpm 0755 ./%{name} %{buildroot}%{_bindir}/%{name}


%files
%{_bindir}/%{name}
%license LICENSE.md
%doc *.md



%changelog
* Mon Sep 2 2024 Brycen <brycengranville@outlook.com> 1.0.0-3
- Initial package
* Tue Sep 3 2024 Brycen <brycengranville@outlook.com> 1.0.0-5
- Removed sysinfo package decreasing binary size and portability and startup time
