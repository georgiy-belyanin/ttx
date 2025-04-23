# TTX - CLI tool for Tarantool developers

TTX is a CLI tool for [Tarantool](https://github.com/tarantool/tarantool) developers. It simplifies working with [Tarantool cluster configuration](https://www.tarantool.io/en/doc/latest/reference/configuration/configuration_reference/) and testing the clusters during the development.

## Usage

If you want to see how to use the utility run `ttx --help`.

Simply run `ttx` to seek for the Tarantool `config.yml` files nearby and start the cluster based on them.

```bash
ttx
# Or explicitly provide the configuration and start the cluster using -c flag.
ttx start -c config.yml
# Or instead of starting the whole cluster run only specific its parts.
ttx start replicaset-003 -c config.yaml
```

## Build & install

To build & install ttx run the following bash script.

```
git clone https://github.com/georgiy-belyanin/ttx.git
cd ttx
go mod tidy
go install
```

`ttx` will likely appear in `$GOPATH/bin/ttx` (i.e. `$HOME/go/bin/ttx`). Make sure this dir is added to your `$PATH`.
