FROM golang:latest as builder 
RUN mkdir -p /go/src/github.com/wtownse/gopxe
ADD . /go/src/github.com/wtownse/gopxe
WORKDIR /go/src/github.com/wtownse/gopxe
#RUN go test ./...
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o main .

FROM docker.io/dokken/rockylinux-9 as tftp

RUN dnf update -y && yum install -y dhcp* epel-release syslinux syslinux-tftpboot  \
grub2-efi-x64-modules grub2-tools-extra grub2-pc-modules glibc-devel \
 && yum clean all 
RUN grub2-mknetdir --net-directory /var/lib/tftpboot/
RUN mkdir -p /var/lib/tftpboot/pxelinux.cfg /opt/localrepo
RUN cp -r /usr/share/syslinux/pxelinux.0 /usr/share/syslinux/ldlinux* /var/lib/tftpboot
RUN cd /var/lib/tftpboot && wget http://boot.ipxe.org/ipxe.efi
ADD ./pxebootImages /var/lib/tftpboot

FROM docker.io/dokken/rockylinux-9
RUN mkdir -p /gopxe/public ; mkdir /gopxe/ksTempl
WORKDIR /gopxe
COPY --from=builder /go/src/github.com/wtownse/gopxe/main /gopxe/
RUN mkdir -p /var/lib/tftpboot /coredhcp
COPY --from=tftp /var/lib/tftpboot/ /var/lib/tftpboot/
ADD ./public /gopxe/public
ADD ./ksTempl /gopxe/ksTempl
ADD ./start-gopxe.sh /gopxe/
EXPOSE 67 67/udp 69/udp 9090 9090/udp
ENTRYPOINT ["/gopxe/start-gopxe.sh"]
