VERSION=$(shell cat ./ottersacle/DEBIAN/control | grep 'Version' | awk '{print $$2}')

build_deb:
	chmod +x ./ottersacle/DEBIAN/postinst
	chmod +x ./ottersacle/DEBIAN/postrm
	chmod +x ./ottersacle/DEBIAN/preinst
	chmod +x ./ottersacle/DEBIAN/prerm
	chmod +x ./ottersacle/usr/bin/openhdc
	dpkg-deb --build ottersacle ottersacle-v$(VERSION)-jammy-amd64.deb
