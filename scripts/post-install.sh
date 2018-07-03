#!/bin/bash

BIN_DIR=/usr/bin
DATA_DIR=/var/lib/aiis
LOG_DIR=/var/log/aiis
SCRIPT_DIR=/usr/lib/aiis/scripts

function install_init {
    cp -f $SCRIPT_DIR/init.sh /etc/init.d/aiis
    chmod +x /etc/init.d/aiis
}

function install_systemd {
    cp -f $SCRIPT_DIR/aiis.service /lib/systemd/system/aiis.service
}

function enable_systemd {
    systemctl enable aiis
}

function enable_update_rcd {
    update-rc.d aiis defaults
}

function enable_chkconfig {
    chkconfig --add aiis
}

if ! id aiis >/dev/null 2>&1; then
    useradd --system -U -M aiis -s /bin/false -d $DATA_DIR
fi
chmod a+rX $BIN_DIR/aiis*

mkdir -p $LOG_DIR
chown -R -L aiis:aiis $LOG_DIR
mkdir -p $DATA_DIR
chown -R -L aiis:aiis $DATA_DIR

test -f /etc/default/aiis || touch /etc/default/aiis

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
