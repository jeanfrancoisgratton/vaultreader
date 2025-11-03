all: 
	cd ../src
	CGO_ENABLED=0 go build -trimpath -ldflags="-s -w -linkmode=internal -buildid=" . -o ./vaultreader.exe

prepare:
	sudo /opt/bin/install_golang.sh 1.25.3 amd64
	sudo apk update --no-cache

install:
	test -d $(DESTDIR)$(PREFIX)/bin || mkdir -p $(DESTDIR)$(PREFIX)/bin
	install -pm 755 $(BUILDDIR)/vaultreader.exe $(DESTDIR)$(PREFIX)/bin/vaultreader
