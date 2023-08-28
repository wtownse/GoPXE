#!/bin/bash 

DATE=`date +%Y-%m-%d-%H:%M:%S`
DHCPD_CONF="/etc/dhcp/dhcpd.conf"
TFTPD_CONF="/etc/xinetd.d/tftp"
TFTPD_BOOT_PATH="/var/lib/tftpboot"

cat > $DHCPD_CONF << EOF
default-lease-time 600;
max-lease-time 7200;

# If this DHCP server is the official DHCP server for the local
# network, the authoritative directive should be uncommented.
authoritative;

# Use this to enble / disable dynamic dns updates globally.
ddns-update-style none;

# No service will be given on this subnet, but declaring it helps the
# DHCP server to understand the network topology.

subnet ${SUBNET} netmask ${NETMASK} {
        option routers                  ${ROUTER};
        option subnet-mask              ${NETMASK};
        option domain-search            "${DOMAIN}";
        option domain-name-servers      ${DNS};
        option time-offset              -18000;     # Eastern Standard Time
        next-server                     $(hostname -I | awk '{print $1}');
        filename                        "/pxelinux.0";
        range  ${DHCPDRANGE};  # reserved DHCPD range e.g. 10.17.224.100 10.17.224.150
}
EOF

function log(){
    echo "$DATE $@" 
    return 0
}

function warn(){
    echo "$DATE $@" 
    return 1
}

function panic(){
    echo "$DATE $@" 
    exit 1
}

## Starting goPXE
log "gopxe is starting..."
/gopxe/main -ksURL $(hostname -I | awk '{print $1}') -wsHOST $(hostname -I | awk '{print $1}') -wsPORT "8080" & 

## Start dhcpd
log "starting dhcpd"
/usr/sbin/dhcpd -4 -f -d --no-pid -cf ${DHCPD_CONF}

## Start tftp
#log "starting tftpd"
#/usr/sbin/in.tftpd --foreground --address 0.0.0.0:69 --secure ${TFTPD_BOOT_PATH}
