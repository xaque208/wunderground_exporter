BIN_NAME=wunderground_exporter
PKG_SITE=gopkg.larch.space

build:
	go get -v
	go build -o $(BIN_NAME) -v

clean:
	rm -f $(BIN_NAME)

publish:
	scp -o StrictHostKeyChecking=no $(BIN_NAME) webdeploy@${PKG_SITE}:/usr/local/www/${PKG_SITE}/
