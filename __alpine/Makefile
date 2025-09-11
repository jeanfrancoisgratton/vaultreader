all: 
	cd ../src
	CGO_ENABLED=0 go build . -o ./vaultreader.exe

prepare:
	sudo /opt/bin/install_golang.sh 1.25.1 amd64
	sudo apk update --no-cache

install:
	test -d $(DESTDIR)$(PREFIX)/bin || mkdir -p $(DESTDIR)$(PREFIX)/bin
	install -pm 755 $(BUILDDIR)/vaultreader.exe $(DESTDIR)$(PREFIX)/bin/vaultreader
	strip $(DESTDIR)$(PREFIX)/bin/vaultreader
