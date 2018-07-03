#!/bin/bash

function uninstall_init {
    rm -f /etc/init.d/aiis
}

function uninstall_systemd {
    rm -f /lib/systemd/system/aiis.service
}

function disable_systemd {
    systemctl disable aiis
}

function disable_update_rcd {
    update-rc.d -f aiis remove
}

function disable_chkconfig {
    chkconfig --del aiis
}

if [[ -f /etc/redhat-release ]]; then
    # RHEL-variant logic
    if [[ "$1" = "0" ]]; then
        # Aiis is no longer installed, remove from init system
        rm -f /etc/default/aiis

        if [[ "$(readlink /proc/1/exe)" == */systemd ]]; then
            disable_systemd
            uninstall_systemd
        else
            # Assuming SysV
            # Run update-rc.d or fallback to chkconfig if not available
            if which update-rc.d &>/dev/null; then
                disable_update_rcd
            else
                disable_chkconfig
            fi
            uninstall_init
        fi
    fi
elif [[ -f /etc/debian_version ]]; then
    # Debian/Ubuntu logic
    if [[ "$1" != "upgrade" ]]; then
        # Remove/purge
        rm -f /etc/default/aiis

        if [[ "$(readlink /proc/1/exe)" == */systemd ]]; then
            disable_systemd
            uninstall_systemd
        else
            # Assuming SysV
            # Run update-rc.d or fallback to chkconfig if not available
            if which update-rc.d &>/dev/null; then
                disable_update_rcd
            else
                disable_chkconfig
            fi
            uninstall_init
        fi
    fi
elif [[ -f /etc/os-release ]]; then
    source /etc/os-release
    if [[ $ID = "amzn" ]]; then
        # Amazon Linux logic
        if [[ "$1" = "0" ]]; then
            # Aiis is no longer installed, remove from init system
            rm -f /etc/default/aiis

            # Run update-rc.d or fallback to chkconfig if not available
            if which update-rc.d &>/dev/null; then
                disable_update_rcd
            else
                disable_chkconfig
            fi
            uninstall_init
        fi
    fi
fi
