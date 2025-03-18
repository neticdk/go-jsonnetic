<p align="center">
    <img alt="jsonnetic logo" src="assets/logo.png" height="300" />
    <h3 align="center">Jsonnetic - a Jsonnet implementation</h3>
</p>

---
# Jsonnetic

CLI Tool based on [jsonnet](https://jsonnet.org), but with a few extra native functions.

## Features
The features added can be seen in the [docs](docs/jsonnetic.md). or as code in the [funcs.go](pkg/jsonnetic/native/funcs.go) file.

# Build Jsonnetic

```sh
git clone git@github.com:neticdk/go-jsonnetic.git
```

```sh
make build
```

Clean slate builds:

```sh
go clean --modcache
go get -u
go build
```
