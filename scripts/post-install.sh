#!/bin/bash

BIN_DIR=/usr/bin
LOG_DIR=/var/log/rush
SCRIPT_DIR=/usr/lib/rush/scripts
LOGROTATE_DIR=/etc/logrotate.d

function install_init {
    cp -f $SCRIPT_DIR/init.sh /etc/init.d/rush
    chmod +x /etc/init.d/rush
}

function install_systemd {
    cp -f $SCRIPT_DIR/rush.service $1
    systemctl enable rush || true
    systemctl daemon-reload || true
}

function install_update_rcd {
    update-rc.d rush defaults
}

function install_chkconfig {
    chkconfig --add rush
}

if ! grep "^rush:" /etc/group &>/dev/null; then
    groupadd -r rush
fi

if ! id rush &>/dev/null; then
    useradd -r -M rush -s /bin/false -d /etc/rush -g rush
fi

test -d $LOG_DIR || mkdir -p $LOG_DIR
chown -R -L rush:rush $LOG_DIR
chmod 755 $LOG_DIR

# Remove legacy symlink, if it exists
if [[ -L /etc/init.d/rush ]]; then
    rm -f /etc/init.d/rush
fi
# Remove legacy symlink, if it exists
if [[ -L /etc/systemd/system/rush.service ]]; then
    rm -f /etc/systemd/system/rush.service
fi

# Add defaults file, if it doesn't exist
if [[ ! -f /etc/default/rush ]]; then
    touch /etc/default/rush
fi

# Add .d configuration directory
if [[ ! -d /etc/rush/rush.d ]]; then
    mkdir -p /etc/rush/rush.d
fi

# Distribution-specific logic
if [[ -f /etc/redhat-release ]] || [[ -f /etc/SuSE-release ]]; then
    # RHEL-variant logic
    if [[ "$(readlink /proc/1/exe)" == */systemd ]]; then
        install_systemd /usr/lib/systemd/system/rush.service
    else
        # Assuming SysVinit
        install_init
        # Run update-rc.d or fallback to chkconfig if not available
        if which update-rc.d &>/dev/null; then
            install_update_rcd
        else
            install_chkconfig
        fi
    fi
elif [[ -f /etc/debian_version ]]; then
    # Debian/Ubuntu logic
    if [[ "$(readlink /proc/1/exe)" == */systemd ]]; then
        install_systemd /lib/systemd/system/rush.service
        deb-systemd-invoke restart rush.service || echo "WARNING: systemd not running."
    else
        # Assuming SysVinit
        install_init
        # Run update-rc.d or fallback to chkconfig if not available
        if which update-rc.d &>/dev/null; then
            install_update_rcd
        else
            install_chkconfig
        fi
        invoke-rc.d rush restart
    fi
elif [[ -f /etc/os-release ]]; then
    source /etc/os-release
    if [[ $ID = "amzn" ]]; then
        # Amazon Linux logic
        install_init
        # Run update-rc.d or fallback to chkconfig if not available
        if which update-rc.d &>/dev/null; then
            install_update_rcd
        else
            install_chkconfig
        fi
    fi
fi
