module kuskcore

go 1.16

replace (
	github.com/tendermint/ed25519 => ./lib/github.com/tendermint/ed25519
	github.com/tendermint/go-wire => github.com/tendermint/go-amino v0.6.2
	github.com/zondax/ledger-goclient => github.com/Zondax/ledger-cosmos-go v0.1.0
	golang.org/x/crypto => ./lib/golang.org/x/crypto
	golang.org/x/net => ./lib/golang.org/x/net
	gonum.org/v1/gonum/mat => github.com/gonum/gonum/mat v0.9.1
)

require (
	github.com/btcsuite/btcutil v1.0.2 // indirect
	github.com/btcsuite/go-socks v0.0.0-20170105172521-4720035b7bfd
	github.com/cenkalti/backoff v2.2.1+incompatible // indirect
	github.com/cespare/cp v1.1.1
	github.com/davecgh/go-spew v1.1.1
	github.com/denisenkom/go-mssqldb v0.0.0-20191124224453-732737034ffd // indirect
	github.com/erikstmartin/go-testdb v0.0.0-20160219214506-8d10e4a1bae5 // indirect
	github.com/fortytw2/leaktest v1.3.0 // indirect
	github.com/fsnotify/fsnotify v1.4.9 // indirect
	github.com/go-kit/kit v0.9.0 // indirect
	github.com/go-logfmt/logfmt v0.5.0 // indirect
	github.com/go-sql-driver/mysql v1.5.0
	github.com/go-stack/stack v1.8.0 // indirect
	github.com/gofrs/uuid v4.4.0+incompatible // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da
	github.com/golang/protobuf v1.4.3
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/uuid v1.2.0
	github.com/gopherjs/gopherjs v0.0.0-20181017120253-0766667cb4d1 // indirect
	github.com/gorilla/websocket v1.5.0
	github.com/grandcat/zeroconf v0.0.0-20190424104450-85eadb44205c
	github.com/hashicorp/go-version v1.3.0
	github.com/holiman/uint256 v1.2.0
	github.com/jinzhu/gorm v1.9.2
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.0.1 // indirect
	github.com/johngb/langreg v0.0.0-20150123211413-5c6abc6d19d2
	github.com/jonboulle/clockwork v0.2.2 // indirect
	github.com/jtolds/gls v4.20.0+incompatible // indirect
	github.com/kr/secureheader v0.2.0
	github.com/lestrrat-go/file-rotatelogs v2.4.0+incompatible
	github.com/lestrrat-go/strftime v1.0.4 // indirect
	github.com/lib/pq v1.1.1 // indirect
	github.com/magiconair/properties v1.8.5 // indirect
	github.com/mattn/go-sqlite3 v1.13.0 // indirect
	github.com/miekg/dns v1.1.15 // indirect
	github.com/mitchellh/mapstructure v1.4.1 // indirect
	github.com/onsi/ginkgo v1.11.0 // indirect
	github.com/pborman/uuid v1.2.1
	github.com/pelletier/go-toml v1.9.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/prometheus/prometheus v1.8.2
	github.com/sirupsen/logrus v1.8.1
	github.com/smartystreets/assertions v0.0.0-20180927180507-b2de0cb4f26d // indirect
	github.com/smartystreets/goconvey v1.6.3 // indirect
	github.com/spf13/afero v1.6.0 // indirect
	github.com/spf13/cast v1.3.1 // indirect
	github.com/spf13/cobra v0.0.5
	github.com/spf13/jwalterweatherman v1.1.0
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/spf13/viper v1.3.2
	github.com/stretchr/testify v1.8.4
	github.com/syndtr/goleveldb v1.0.0
	github.com/tendermint/ed25519 v0.0.0-20171027050219-d8387025d2b9
	github.com/tendermint/go-crypto v0.2.0
	github.com/tendermint/go-wire v0.16.0
	github.com/tendermint/tmlibs v0.9.0
	github.com/toqueteos/webbrowser v1.2.0
	golang.org/x/crypto v0.17.0
	golang.org/x/sync v0.1.0
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/fatih/set.v0 v0.1.0
	gopkg.in/karalabe/cookiejar.v2 v2.0.0-20150724131613-8dcd6a7f4951
	gopkg.in/yaml.v2 v2.4.0 // indirect
)
