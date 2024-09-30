# Firejail profile for rokon
# Description: Simple GTK4 Roku remote

# Note: According to my testing, this breaks under prime-run. It will trigger undesirable software rendering.

# Start with a strict sandbox
include default.profile

# Persistent global definitions
include globals.local

# Allow network access (local and external)
netfilter


mkdir ${HOME}/.cache/rokon
mkdir ${HOME}/.local/share/rokon
whitelist ${HOME}/.config/rokon     # Allow access to app config
whitelist ${HOME}/.local/share/rokon  # Allow access to app data
whitelist ${DOCUMENTS}
whitelist ${DOWNLOADS}
whitelist /usr/share/rokon
include whitelist-common.inc
include whitelist-usr-share-common.inc
include whitelist-var-common.inc

ignore noexec ${HOME}

# Noexec for the XDG Download directory
noexec ${HOME}/Downloads

# Read-only access to system config
read-only /etc

# Allow access to machine ID for telemetry
whitelist /etc/machine-id

# Allow D-Bus communication for GTK4
# dbus-user
# dbus-system

# Restrict access to sensitive directories
blacklist /media
blacklist /mnt
blacklist /srv
blacklist /proc/sys
blacklist /proc/acpi

apparmor
caps.drop all
nodvd
nogroups
nonewprivs
noroot
notv
nou2f
novideo
protocol unix
seccomp
tracelog

disable-mnt
private-bin rokon
private-cache
private-dev
private-etc @x11,gconf
private-tmp

# Drop all capabilities, keep only networking
caps.keep net_bind_service
caps.drop net_admin,net_raw

# AppImage-specific settings
whitelist /tmp/.mount_rokon*
read-write ${HOME}/.cache/rokon
read-write ${HOME}/.local/share/rokon

restrict-namespaces
