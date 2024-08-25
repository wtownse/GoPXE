module github.com/wtownse/gopxe

go 1.22

toolchain go1.23.0

replace github.com/coreos/bbolt v1.3.11 => go.etcd.io/bbolt v1.3.11

replace go.etcd.io/bbolt v1.3.11 => github.com/coreos/bbolt v1.3.11

require (
	github.com/coredhcp/coredhcp v0.0.0-20240709092356-bd8c8089a5ab
	github.com/coreos/bbolt v1.3.11
	github.com/gorilla/mux v1.8.1
	github.com/mattn/go-sqlite3 v1.14.22
	github.com/sirupsen/logrus v1.9.3
	github.com/spf13/pflag v1.0.6-0.20201009195203-85dd5c8bc61c
	gopkg.in/yaml.v3 v3.0.1
	pack.ag/tftp v1.0.0
)

require (
	github.com/bits-and-blooms/bitset v1.13.0 // indirect
	github.com/chappjc/logrus-prefix v0.0.0-20180227015900-3a1d64819adb // indirect
	github.com/fsnotify/fsnotify v1.7.0 // indirect
	github.com/google/gopacket v1.1.19 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/insomniacslk/dhcp v0.0.0-20240227161007-c728f5dd21c8 // indirect
	github.com/josharian/native v1.1.0 // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mgutz/ansi v0.0.0-20200706080929-d51e80ef957d // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/nxadm/tail v1.4.11 // indirect
	github.com/pelletier/go-toml/v2 v2.2.2 // indirect
	github.com/pierrec/lz4/v4 v4.1.21 // indirect
	github.com/rifflock/lfshook v0.0.0-20180920164130-b9218ef580f5 // indirect
	github.com/sagikazarmark/locafero v0.4.0 // indirect
	github.com/sagikazarmark/slog-shim v0.1.0 // indirect
	github.com/sourcegraph/conc v0.3.0 // indirect
	github.com/spf13/afero v1.11.0 // indirect
	github.com/spf13/cast v1.6.0 // indirect
	github.com/spf13/viper v1.19.0 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	github.com/u-root/uio v0.0.0-20230305220412-3e8cd9d6bf63 // indirect
	go.etcd.io/bbolt v1.3.11 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/crypto v0.25.0 // indirect
	golang.org/x/exp v0.0.0-20240112132812-db7319d0e0e3 // indirect
	golang.org/x/net v0.27.0 // indirect
	golang.org/x/sys v0.22.0 // indirect
	golang.org/x/term v0.22.0 // indirect
	golang.org/x/text v0.16.0 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
)
