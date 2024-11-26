# [Metamod Go Plugin Example](https://github.com/et-nik/metamod-go-example)

This is example of usage of [Metamod-Go](https://github.com/et-nik/metamod-go) library.

This plugin adds `entinfo` and `traceline` commands to a game server.
* `entinfo` prints information about entities in the game.
* `traceline` prints information about a trace line from entity to forward direction.

## Build plugin

Build using Docker `golang:1.23.3-bookworm` image:
```
dpkg --add-architecture i386
apt-get update
apt-get install gcc-multilib

CGO_ENABLED=1 GOOS=linux GOARCH=386 go build -o example.so -buildvcs=false -buildmode=c-shared
```