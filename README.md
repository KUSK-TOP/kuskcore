Kusk
======

[![Build](https://github.com/KUSK-TOP/kuskcore/actions/workflows/main.yml/badge.svg?branch=main)](https://github.com/KUSK-TOP/kuskcore/actions)
[![Supports Windows](https://img.shields.io/badge/support-Windows-blue?logo=Windows)](https://github.com/KUSK-TOP/kuskcore/releases/latest)
[![Supports Linux](https://img.shields.io/badge/support-Linux-yellow?logo=Linux)](https://github.com/KUSK-TOP/kuskcore/releases/latest)
[![License](https://img.shields.io/github/license/KUSK-TOP/kuskcore)](https://github.com/KUSK-TOP/kuskcore/blob/master/LICENSE)
[![Latest Release](https://img.shields.io/github/v/release/KUSK-TOP/kuskcore?label=latest%20release)](https://github.com/KUSK-TOP/kuskcore/releases/latest)
[![Downloads](https://img.shields.io/github/downloads/KUSK-TOP/kuskcore/total)](https://github.com/KUSK-TOP/kuskcore/releases)
[![KUSK Community](https://img.shields.io/discord/1217075571528564736?label=EIYARO%20Project%20Discord&logo=discord)](https://discord.gg/V4ue4CVMKY)


## What is Kusk?

Kusk is a network of nodes designed to operate and connect to a highly scalable blockchain network and confirms the Kusk blockchain protocol, which allows participants to share and transfer digital assets. For more details, see the [White Paper](https://kusk.top/KUSK-WhitePaper.html) for more details.

In the current state `kusk` is able to:

- Manage key, account as well as asset
- Send transactions


## Building from source

### Requirements

- [Go](https://golang.org/doc/install) version 1.8 or higher, with `$GOPATH` set to your preferred directory

### Installation

Ensure Go with the supported version is installed properly:

```bash
$ wget -c https://go.dev/dl/go1.22.2.linux-amd64.tar.gz -O - | sudo tar -xz -C /usr/local

$ export PATH=$PATH:/usr/local/go/bin

$ source ~/.profile

$ go version
```

- Get the source code

``` bash
$ git clone https://github.com/KUSK-TOP/kuskcore.git $GOPATH/src/github.com/kuskcore
```

- Build source code

``` bash
$ cd $GOPATH/src/github.com/kuskcore
$ make kuskd    # go build kuskd
$ make kuskcli  # go build kuskcli
```

When successfully building the project, the `kuskd` and `kuskcli` binary should be present in `cmd/kuskd` and `cmd/kuskcli` directory, respectively.

### Executables

The Kusk project comes with several executables found in the `cmd` directory.

| Command      | Description                                                  |
| ------------ | ------------------------------------------------------------ |
| **kuskd**   | kuskd command can help to initialize and launch kusk domain by custom parameters. `kuskd --help` for command line options. |
| **kuskcli** | Our main Kusk CLI client. It is the entry point into the Kusk network (main-, test- or private net), capable of running as a full node archive node (retaining all historical state). It can be used by other processes as a gateway into the Kusk network via JSON RPC endpoints exposed on top of HTTP, WebSocket and/or IPC transports. `kuskcli --help` |

## Running kusk

Currently, kusk is still in active development and a ton of work needs to be done, but we also provide the following content for these eager to do something with `kusk`. This section won't cover all the commands of `kuskd` and `kuskcli` at length, for more information, please the help of every command, e.g., `kuskcli help`.

### Initialize

First of all, initialize the node:

```bash
$ cd ./cmd/kuskd
$ ./kuskd init --chain_id mainnet
```

There are three options for the flag `--chain_id`:

- `mainnet`: connect to the mainnet.
- `testnet`: connect to the testnet wisdom.
- `solonet`: standalone mode.

After that, you'll see `config.toml` generated, then launch the node.

### launch

``` bash
$ nohup ./kuskd node &
```

available flags for `kuskd node`:

```
Flags:
      --auth.disable                     Disable rpc access authenticate
      --chain_id string                  Select network type
  -h, --help                             help for node
      --log_file string                  Log output file (default "log")
      --log_level string                 Select log level(debug, info, warn, error or fatal)
      --p2p.dial_timeout int             Set dial timeout (default 3)
      --p2p.handshake_timeout int        Set handshake timeout (default 30)
      --p2p.keep_dial string             Peers addresses try keeping connecting to, separated by ',' (for example "1.1.1.1:46657;2.2.2.2:46658")
      --p2p.laddr string                 Node listen address. (0.0.0.0:0 means any interface, any port) (default "tcp://0.0.0.0:46656")
      --p2p.lan_discoverable             Whether the node can be discovered by nodes in the LAN (default true)
      --p2p.max_num_peers int            Set max num peers (default 50)
      --p2p.node_key string              Node key for p2p communication
      --p2p.proxy_address string         Connect via SOCKS5 proxy (eg. 127.0.0.1:1086)
      --p2p.proxy_password string        Password for proxy server
      --p2p.proxy_username string        Username for proxy server
      --p2p.seeds string                 Comma delimited host:port seed nodes
      --p2p.skip_upnp                    Skip UPNP configuration
      --prof_laddr string                Use http to profile kuskd programs
      --vault_mode                       Run in the offline enviroment
      --wallet.disable                   Disable wallet
      --wallet.rescan                    Rescan wallet
      --wallet.txindex                   Save global tx index
      --web.closed                       Lanch web browser or not
      --ws.max_num_concurrent_reqs int   Max number of concurrent websocket requests that may be processed concurrently (default 20)
      --ws.max_num_websockets int        Max number of websocket connections (default 25)

Global Flags:
      --home string   root directory for config and data
  -r, --root string   DEPRECATED. Use --home (default "/Users/zcc/Library/Application Support/Kusk")
      --trace         print out full stack trace on errors
```

Given the `kuskd` node is running, the general workflow is as follows:

- create key, then you can create account.
- send transaction, i.e., build, sign and submit transaction.
- query all kinds of information, let's say, avaliable key, account, key, balances, transactions, etc.

### Dashboard

Access the dashboard:

```
$ open http://localhost:9888/
```


## Contributing

Thank you for considering helping out with the source code! Any contributions are highly appreciated, and we are grateful for even the smallest of fixes!

If you run into an issue, feel free to [kusk issues](https://github.com/KUSK-TOP/kuskcore/issues/) in this repository. We are glad to help!

## License

[AGPL v3](./LICENSE)
