FROM golang:latest as builder 
RUN mkdir -p /go/src/github.com/wtownse/gopxe
ADD . /go/src/github.com/wtownse/gopxe
WORKDIR /go/src/github.com/wtownse/gopxe
#RUN go test ./...
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

FROM docker.io/dokken/centos-8

RUN dnf update -y && yum install -y kernel-headers xinetd dhcp* epel-release syslinux syslinux-tftpboot  \
grub2-efi-x64-modules grub2-tools-extra grub2-pc-modules \
shim-ia32 && yum clean all 
RUN grub2-mknetdir --net-directory /var/lib/tftpboot/
EXPOSE 67 67/udp 69/udp 9090 9090/udp
RUN mkdir -p /var/lib/tftpboot/pxelinux.cfg /opt/localrepo
RUN cp -r /usr/share/syslinux/pxelinux.0 /usr/share/syslinux/ldlinux* /var/lib/tftpboot
ADD ./pxebootImages /var/lib/tftpboot
ADD ./pxeFiles/* /var/lib/tftpboot
RUN mkdir -p /gopxe/public ; mkdir /gopxe/ksTempl
WORKDIR /gopxe
COPY --from=builder /go/src/github.com/wtownse/gopxe/main /gopxe/
ADD ./public /gopxe/public
ADD ./ksTempl /gopxe/ksTempl
ADD ./start-gopxe.sh /gopxe/
HEALTHCHECK --interval=4m --timeout=60s CMD curl --fail http://localhost:9090/health || exit 1
ENTRYPOINT ["/gopxe/start-gopxe.sh"]
