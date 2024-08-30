#!/bin/bash 

DATE=`date +%Y-%m-%d-%H:%M:%S`
DHCPD_CONF="/etc/dhcp/dhcpd.conf"
TFTPD_CONF="/etc/xinetd.d/tftp"
TFTPD_BOOT_PATH="/var/lib/tftpboot"

if [[ -z "${TFTPD_NEXT_SERVER}" ]]; then
TFTPD_NEXT_SERVER=$(hostname -I | awk '{print $1}')
fi

if [[ -z "${WSHOST}" ]]; then
WSHOST="$(hostname -I | awk '{print $1}')"
fi

if [[ -z "${WSPORT}" ]]; then
WSPORT="8080"
fi

if [[ -z "${PXEFILE64}" ]]; then
PXEFILE64="boot/grub2/x86_64-efi/core.efi"
fi

if [[ -z "${PXEFILE32}" ]]; then
PXEFILE32="boot/grub2/x86_64-efi/core.efi"
fi

if [[ -z "${BOOTFILEPATH}" ]]; then
BOOTFILEPATH="/coreos/"
fi

if [[ -z "${WFILEPATH}" ]]; then
WFILEPATH="/coreos/"
fi

if [[ -z "${LISTENADDRESS}" ]]; then
LISTENADDRESS="0.0.0.0"
fi

cat > $DHCPD_CONF << EOF
default-lease-time 300;
max-lease-time 600;

# If this DHCP server is the official DHCP server for the local
# network, the authoritative directive should be uncommented.
authoritative;

# Use this to enble / disable dynamic dns updates globally.
ddns-update-style none;

# No service will be given on this subnet, but declaring it helps the
# DHCP server to understand the network topology.

option architecture code 93 = unsigned integer 16 ;
subnet ${SUBNET} netmask ${NETMASK} {
        option routers                  ${ROUTER};
        option subnet-mask              ${NETMASK};
        option domain-search            "${DOMAIN}";
        option domain-name-servers      ${DNS};
        option time-offset              -18000;     # Eastern Standard Time
        range  ${DHCPDRANGE};  # reserved DHCPD range e.g. 10.17.224.100 10.17.224.150
          class "pxeclients" {
                  match if substring (option vendor-class-identifier, 0, 9) = "PXEClient";

        next-server                     $TFTPD_NEXT_SERVER;
        if option architecture = 00:06 {
          filename "${PXEFILE32}";
          option boot-size 4328;
        } elsif option architecture = 00:07 {
          filename "${PXEFILE64}";
          option boot-size 4328;
        } elsif option architecture = 00:09 {
          filename "${PXEFILE64}";
          option boot-size 4328;
        } else {
          filename "pxelinux.0";
        }
  }
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
#/gopxe/main -ksURL $(hostname -I | awk '{print $1}') -wsHOST $WSHOST -wsPORT $WSPORT -bootFILEPath $BOOTFILEPATH
/gopxe/main  -K $(hostname -I | awk '{print $1}') -i "$LISTENADDRESS" -R "$DHCPRANGE" -M $NETMASK -G $ROUTER -D $DNS --port "9090" -s $WSHOST -r $WSPORT  -U -a /var/lib/tftpboot/coreos/agent-config.yaml -b $BOOTFILEPATH -w $WFILEPATH

## Start dhcpd
#log "starting dhcpd"
#/usr/sbin/dhcpd -4 -f -d --no-pid -cf ${DHCPD_CONF}

## Start tftp
#log "starting tftpd"
#/usr/sbin/in.tftpd --foreground --address 0.0.0.0:69 --secure ${TFTPD_BOOT_PATH}
