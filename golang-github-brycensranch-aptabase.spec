# Generated by go2rpm 1.14.0
%bcond check 1
%bcond bootstrap 0

%global debug_package %{nil}
%if %{with bootstrap}
%global __requires_exclude %{?__requires_exclude:%{__requires_exclude}|}^golang\\(.*\\)$
%endif

# https://github.com/brycensranch/go-aptabase
%global goipath         github.com/brycensranch/go-aptabase/pkg
%global commit          b987899e04dd7d42345e74a6427ce11a1abec24e

%gometa -L -f

%global common_description %{expand:
Golang SDK for Aptabase: Open Source, Privacy-First and Simple Analytics for
Mobile, Desktop and Web Apps.}

%global golicenses      LICENSE
%global godocs          example README.md

Name:           golang-github-brycensranch-aptabase
Version:        0
Release:        %autorelease -p
Summary:        Golang SDK for Aptabase: Open Source, Privacy-First and Simple Analytics for Mobile, Desktop and Web Apps

License:        MIT
URL:            %{gourl}
Source:         %{gosource}

%description %{common_description}

%gopkg

%prep
%goprep -A
%autopatch -p1

%if %{without bootstrap}
%generate_buildrequires
%go_generate_buildrequires
%endif

%install
%gopkginstall

%if %{without bootstrap}
%if %{with check}
%check
%gocheck
%endif
%endif

%gopkgfiles

%changelog
%autochangelog