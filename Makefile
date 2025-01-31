MODULE_DIR=$(PREFIX)/lib/boli/modules

build: clean
	mkdir build
	go build -o build/boli cmd/main/main.go

clean:
	rm -rf build

install: build
	mkdir -p $(PREFIX)/bin
	cp -f build/boli $(PREFIX)/bin
	rm -rf $(MODULE_DIR)
	mkdir -p $(MODULE_DIR)
	cp -rf modules/* $(MODULE_DIR)

uninstall:
	rm -f $(PREFIX)/bin/boli
	rm -rf $(PREFIX)/lib/boli