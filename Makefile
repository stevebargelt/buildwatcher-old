SOURCEDIR=.
SOURCES = $(shell find $(SOURCEDIR) -name '*.go')
VERSION=$(shell git describe --always --tags)
BINARY=bin/buildwatcher

.PHONY:pi
pi:
	env GOOS=linux GOARCH=arm go build -o $(BINARY) -ldflags "-s -w -X main.Version=$(VERSION)" commands/*.go
	cp build/buildwatcher.yml bin/config.yml	

.PHONY:pi-zero
pi-zero:
	env GOARM=6 GOOS=linux GOARCH=arm go build -o $(BINARY) -ldflags "-s -w -X main.Version=$(VERSION)" *.go

.PHONY:deploy
deploy:
	scp bin/buildwatcher pi@192.168.1.6:buildwatcher
	scp bin/config.yml pi@192.168.1.6:config.yml
test:
	go test -cover -v -race ./...

# .PHONY: go-get
# go-get:
# 	go get -u github.com/boltdb/bolt/...
# 	go get -u gopkg.in/yaml.v2
# 	go get -u github.com/stretchr/gomniauth
# 	go get -u github.com/stretchr/gomniauth/providers/google
# 	go get -u github.com/kidoman/embd
# 	go get -u github.com/gorilla/mux
# 	go get -u gopkg.in/robfig/cron.v2
# 	go get -u github.com/dustin/go-humanize
# 	go get -u github.com/ranjib/adafruitio

.PHONY: vet
vet:
	go vet ./...

# .PHONY: build
# build: clean go-get test bin

# .PHONY: deb
# deb:
# 	mkdir -p dist/var/lib/buildwatcher/assets dist/usr/bin dist/etc/buildwatcher
# 	cp bin/buildwatcher dist/usr/bin/buildwatcher
# 	cp build/buildwatcher.yml dist/etc/buildwatcher/config.yml
# 	bundle exec fpm -t deb -s dir -a armhf -n buildwatcher -v $(VERSION) -m steve@bargelt.com --deb-systemd build/buildwatcher.service -C dist  -p buildwatcher-$(VERSION).deb .

 .PHONY: clean
 clean:
	-rm -rf *.deb
	-rm -rf dist
	-rm *.db
