build: clean
	mkdir build
	go build -o build/boli cmd/main/main.go

clean:
	rm -rf build

install: build
	mkdir -p $(PREFIX)/bin
	cp -f build/boli $(PREFIX)/bin
	rm -rf $(PREFIX)/lib/boli/modules
	mkdir -p $(PREFIX)/lib/boli/modules
	cp -rf modules/* $(PREFIX)/lib/boli/modules

uninstall:
	rm -f $(PREFIX)/bin/boli
	rm -rf $(PREFIX)/lib/boli/modules