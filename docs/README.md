# Frieza

Cleanup your cloud ressources!

Frieza can remove all resources from a cloud account or resources which are not part of a "snapshot".

# Features

- Ready to support multiple providers
- Can store resources from multiple profiles in one snapshot

# Installation

You can go to release page and download the latest frieza binary or build it yourself.

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

To delete ALL resources of a profile:
```
frieza destroy-all myDevAccount
```

To delete all newly created resources since myFristSnap:
```
frieza destroy myFristSnap
```

Note that plan is show before any action and confirmation is asked by default.
You can overide those behavior with `--auto-approve` option.

# License

> Copyright Outscale SAS
>
> BSD-3-Clause

This `LICENSE` contain raw licenses terms following spdx naming.
You can check which license apply to which copyright owner through `.reuse/dep5` specification.
You can test [reuse](https://reuse.software/.) compliance by running `make test-reuse`.