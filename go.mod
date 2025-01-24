module github.com/wtownse/GoPXE

go 1.23.0

replace github.com/coreos/bbolt v1.3.11 => go.etcd.io/bbolt v1.3.11

replace go.etcd.io/bbolt v1.3.11 => github.com/coreos/bbolt v1.3.11

require (
	github.com/coredhcp/coredhcp v0.0.0-20250113163832-cbc175753a45
	github.com/coreos/bbolt v1.3.11
	github.com/gorilla/mux v1.8.1
	github.com/mattn/go-sqlite3 v1.14.24
	github.com/sirupsen/logrus v1.9.3
	github.com/spf13/pflag v1.0.6-0.20201009195203-85dd5c8bc61c
	github.com/wtownse/coredhcp v0.0.0-20250124134241-9b593250de94
	github.com/wtownse/gopxe v0.0.0-20250123220952-fc680b24dbde
	gopkg.in/yaml.v3 v3.0.1
	pack.ag/tftp v1.0.0
)

require (
	github.com/bits-and-blooms/bitset v1.20.0 // indirect
	github.com/chappjc/logrus-prefix v0.0.0-20180227015900-3a1d64819adb // indirect
	github.com/fsnotify/fsnotify v1.8.0 // indirect
	github.com/google/gopacket v1.1.19 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/insomniacslk/dhcp v0.0.0-20250109001534-8abf58130905 // indirect
	github.com/magiconair/properties v1.8.9 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mgutz/ansi v0.0.0-20200706080929-d51e80ef957d // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/pelletier/go-toml/v2 v2.2.3 // indirect
	github.com/pierrec/lz4/v4 v4.1.22 // indirect
	github.com/rifflock/lfshook v0.0.0-20180920164130-b9218ef580f5 // indirect
	github.com/sagikazarmark/locafero v0.6.0 // indirect
	github.com/sagikazarmark/slog-shim v0.1.0 // indirect
	github.com/sourcegraph/conc v0.3.0 // indirect
	github.com/spf13/afero v1.11.0 // indirect
	github.com/spf13/cast v1.7.1 // indirect
	github.com/spf13/viper v1.19.0 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	github.com/u-root/uio v0.0.0-20240224005618-d2acac8f3701 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/crypto v0.32.0 // indirect
	golang.org/x/exp v0.0.0-20241210194714-1829a127f884 // indirect
	golang.org/x/net v0.34.0 // indirect
	golang.org/x/sys v0.29.0 // indirect
	golang.org/x/term v0.28.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7 // indirect
)
