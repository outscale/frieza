# Frieza

Cleanup your cloud ressources!

# Usecases

The main usecase is to free all resources inside a cloud account (e.g. `frieza nuke regionEu2`)

An other usecase is to use Frieza for cleaning additional resources since a known state:
1. You want to keep important resources on your account (virtual machines, volumes, etc)
2. Make a "snapshot" (e.g. `frieza snap new cleanAccountState regionEu2`
3. Run some experiment which create a number of resources
4. Once done, cleanup those additional resources with `frieza clean cleanAccountState`

# Features

- Support multiple providers, see [list of all providers and supported objects](providers.md)
- Can store resources from multiple profiles in one snapshot

# Installation

You can go to [release page](https://github.com/outscale-dev/frieza/releases) and download the latest frieza binary.

Alternatively, you can also install frieza from sources:
```
git clone https://github.com/outscale-dev/frieza.git
cd frieza
make install
```

# Building

To build frieza, you will need golang and Make utilities:
- run `make build`
- binary is ready in `cmd/frieza/frieza`

# Usage

Type `frieza` to list all sub command, use `--help` parameter for more details.

## Manage Profiles

The subcommand `profile` allow you to manage all your provider configuration and test them.
Profiles are stored by default with frieza configuration in `~/.frieza/config.json`.

```
frieza profile new outscale_oapi --help
frieza profile new outscale_oapi myDevAccount --region=eu-west-2 --ak=XXX --sk=YYY
frieza profile test myDevAccount
frieza profile list
frieza profile describe myDevAccount
frieza profile rm myDevAccount
```

## Manage Snapshots

Frieza snapshots are only a listing of resources from one or more profiles at a specific time.
Snapshots are stored by default in `~/.frieza/snapshots/`.

```
frieza snapshot new myFristSnap myDevAccount myOtherAccount
frieza snapshot list
frieza snapshot describe myFristSnap
frieza snapshot rm myFristSnap
```

## Delete resources

To delete all newly created resources since myFristSnap:
```
frieza clean myFristSnap
```

To delete ALL resources of a profile:
```
frieza nuke myDevAccount
```

Note that a listing of deleted resources is show before any action.

Confirmation is asked by default but you can overide this behavior with `--auto-approve` option.

# License

> Copyright Outscale SAS
>
> BSD-3-Clause

`LICENSE` folder contain raw licenses terms following spdx naming.

You can check which license apply to which copyright owner through `.reuse/dep5` specification.

You can test [reuse](https://reuse.software/.) compliance by running `make test-reuse`.
