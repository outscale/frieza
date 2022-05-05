# Frieza clean action

This action will clean the resources created between the call to the end of the calling GitHub Action

## Inputs

## `access_key`

**Required** The 3DS Outscale Access Key.

## `secret_key`

**Required** The 3DS Outscale Secret Key.

## `region`

**Required** The 3DS Outscale region.

## `frieza_version`

**Required** The version of Frieza to use. Default is "`latest`".

## Example usage
```
uses: outscale-dev/github_action/frieza@master
with:
    access_key: ${{ secrets.OSC_ACCESS_KEY }}
    secret_key: ${{ secrets.OSC_SECRET_KEY }}
    region: ${{ secrets.OSC_REGION }}
```