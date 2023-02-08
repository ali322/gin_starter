Gin API Starter
===

gin api service starter

## Get Started

```bash
#1. clone this repository

#2. install all dependencies
make install

#3. start project
make start
```

## Release binary

```bash
#1. change arch or os in Makefile if needed

#2. build new binary release
make build
```

## Deploy service

```bash
#1. upload binary file
scp -C api-starter config.yml api.service user@deploy.server:/path/to/app/

#2. copy service file into systemd config directory
cp app.service /etc/systemd/system/

#3. reload and start daemon  xr_scene service
systemctl daemon-reload && systemctl start api
```
