VERSION = 0.0-$(shell date +"%Y%m%d-%H%M")-$(shell git log -n 1 --pretty="format:%h")

deb:
	mkdir -p djwheel_$(VERSION)/usr/local/bin
	go build -o djwheel_$(VERSION)/usr/local/bin/djwheel .
	cp -r DEBIAN/ djwheel_$(VERSION)/
	sed -i 's/DEBVERSION/$(VERSION)/g' djwheel_$(VERSION)/DEBIAN/control
	dpkg-deb --build djwheel_$(VERSION) .
	rm -rf djwheel_$(VERSION)
