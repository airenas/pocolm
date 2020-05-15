
build-pocolm:
	scripts/check_dependencies.sh
	$(MAKE) -C src

build-tools:
	cd tools && $(MAKE) build

all: build-pocolm build-tools

