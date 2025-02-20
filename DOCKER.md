Docker Instructions:
====================

The Dockerfile in this repo assumes a completely self-contained setup, where I2P resides in the same container as gitea.
This is purely for simplicity's sake.
In order to build it, use:

```sh
docker build -t go-i2p/go-gittisane .
```

then in order to run it, use:

```sh
docker run --name i2p-gittisane -d go-i2p/go-gittisane
docker log i2p-gittinsane
```