# Keyvault in GoLang

## Install dependencies
```sh
$ dep ensure
```

## Build
```sh
$ go build -o keyvault-get-secret main.go
```
## Run

```sh
export AZURE_TENANT_ID="xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
export AZURE_CLIENT_ID="xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
 export AZURE_CLIENT_SECRET="YOUR-CLIENT-SECRET-ID"
export VAULT_BASE_URL="https://example.vault.azure.net/"
export SECRET_NAME="supersecret"

./keyvault-get-secret

2019/05/05 13:36:01 Authorization: 		Bearer ey(....)xA

2019/05/05 13:36:02 Secret value is:    I feel the need - the need for speed!
```

## Release

The release script uses `gox` for fast parralel cross compiling of the binaries. Install `gox` if it itn't available on your system yet.
```
$ go get -u github.com/mitchellh/gox
```

```
./release.sh
```