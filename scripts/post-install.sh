#!/bin/bash

BIN_DIR=/usr/bin
DATA_DIR=/var/lib/rush
LOG_DIR=/var/log/rush
SCRIPT_DIR=/usr/lib/rush/scripts

function install_init {
    cp -f $SCRIPT_DIR/init.sh /etc/init.d/rush
    chmod +x /etc/init.d/rush
}

function install_systemd {
    cp -f $SCRIPT_DIR/rush.service /lib/systemd/system/rush.service
}

function enable_systemd {
    systemctl enable rush
}

function enable_update_rcd {
    update-rc.d rush defaults
}

function enable_chkconfig {
    chkconfig --add rush
}

if ! id rush >/dev/null 2>&1; then
    useradd --system -U -M rush -s /bin/false -d $DATA_DIR
fi
chmod a+rX $BIN_DIR/rush*

mkdir -p $LOG_DIR
chown -R -L rush:rush $LOG_DIR
mkdir -p $DATA_DIR
chown -R -L rush:rush $DATA_DIR

test -f /etc/default/rush || touch /etc/default/rush

# Distribution-specific logic
if [[ -f /etc/redhat-release ]]; then
    # RHEL-variant logic
    if [[ "$(readlink /proc/1/exe)" == */systemd ]]; then
        install_systemd
        # Do not enable service
    else
        # Assuming SysV
        install_init
        # Do not enable service
    fi
elif [[ -f /etc/debian_version ]]; then
    # Debian/Ubuntu logic
    if [[ "$(readlink /proc/1/exe)" == */systemd ]]; then
        install_systemd
        enable_systemd
    else
        # Assuming SysV
        install_init
        # Run update-rc.d or fallback to chkconfig if not available
        if which update-rc.d &>/dev/null; then
            enable_update_rcd
        else
            enable_chkconfig
        fi
    fi
elif [[ -f /etc/os-release ]]; then
    source /etc/os-release
    if [[ $ID = "amzn" ]]; then
        # Amazon Linux logic
        install_init
        # Do not enable service
    fi
fi
