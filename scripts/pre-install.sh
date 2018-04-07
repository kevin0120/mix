#!/bin/bash

if [[ -d /etc/opt/rush ]]; then
    # Legacy configuration found
    if [[ ! -d /etc/rush ]]; then
        # New configuration does not exist, move legacy configuration to new location
        echo -e "Please note, Rush's configuration is now located at '/etc/rush' (previously '/etc/opt/rush')."
        mv -vn /etc/opt/rush /etc/rush

        if [[ -f /etc/rush/rush.conf ]]; then
            backup_name="rush.conf.$(date +%s).backup"
            echo "A backup of your current configuration can be found at: /etc/rush/${backup_name}"
            cp -a "/etc/rush/rush.conf" "/etc/rush/${backup_name}"
        fi
    fi
fi
