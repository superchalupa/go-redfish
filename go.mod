module github.com/superchalupa/sailfish

require (
	github.com/BurntSushi/toml v0.3.1 // indirect
	github.com/Knetic/govaluate v3.0.1-0.20171022003610-9aa49832a739+incompatible
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/etcd-io/bbolt v1.3.2
	github.com/fsnotify/fsnotify v1.4.7
	github.com/go-stomp/stomp v2.0.5+incompatible
	github.com/golang/protobuf v1.3.5
	github.com/gorilla/context v1.1.1 // indirect
	github.com/gorilla/mux v1.6.2
	github.com/inconshreveable/log15 v0.0.0-20180818164646-67afb5ed74ec
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/jmoiron/sqlx v1.2.0
	github.com/json-iterator/go v1.1.7
	github.com/looplab/eventhorizon v0.4.0
	github.com/mattn/go-colorable v0.0.9 // indirect
	github.com/mattn/go-isatty v0.0.4 // indirect
	github.com/mattn/go-sqlite3 v1.10.0
	github.com/mitchellh/go-homedir v1.0.0 // indirect
	github.com/mitchellh/mapstructure v1.0.0
	github.com/prometheus/client_golang v1.1.0 // indirect
	github.com/spacemonkeygo/openssl v0.0.0-20180913232302-a66df3e4f582 // indirect
	github.com/spf13/cobra v0.0.3 // indirect
	github.com/spf13/pflag v1.0.2
	github.com/spf13/viper v1.2.0
	github.com/stretchr/testify v1.3.0
	golang.org/x/xerrors v0.0.0-20191011141410-1b5146add898
)

// replace github.com/looplab/eventhorizon => github.com/superchalupa/eventhorizon v0.0.1

replace github.com/looplab/eventhorizon => github.com/looplab/eventhorizon v0.2.1-0.20180328082012-7067a22f516d

go 1.13
