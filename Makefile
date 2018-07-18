NAME := auth-server
VERSION := 0.1.0
RELEASE := 1
PKG_HOME_ROOT := /opt

BUILD_DATE := $(shell LANG=c date)
GO_VERSION := $(shell go version | awk '{print $$3" "$$4}')
GIT_REVISION := $(shell git rev-list -1 HEAD)

.PHONY: all test clean build build-deb dep clean-build clean-deb

build:
	mkdir -p build

	go build -o build/$(NAME) -ldflags="-X main.VERSION=$(VERSION)-$(RELEASE) -X 'main.BUILD_DATE=$(BUILD_DATE)' \
		-X 'main.GIT_REVISION=$(GIT_REVISION)' -X 'main.GO_VERSION=$(GO_VERSION)'"

build-deb:
	rm -rf $(NAME)-$(VERSION)
	cp -a deb-root $(NAME)-$(VERSION)
	sed -e "s#%VERSION%#$(VERSION)#" \
	        -e "s#%NAME%#$(NAME)#" \
	        -e "s#%RELEASE%#$(RELEASE)#" deb-root/DEBIAN/control.in > $(NAME)-$(VERSION)/DEBIAN/control

	mkdir -p $(NAME)-$(VERSION)/$(PKG_HOME_ROOT)/$(NAME)
	if [ ! -e "build/$(NAME)"  ]; then \
		echo "Error: build/$(NAME) not found. run 'make build' first"; \
		exit 1; \
	fi
	strip build/$(NAME)
	cp build/$(NAME) $(NAME)-$(VERSION)/$(PKG_HOME_ROOT)/$(NAME)
	fakeroot dpkg-deb --build $(NAME)-$(VERSION) && mv $(NAME)-$(VERSION).deb $(NAME)_$(VERSION)-$(RELEASE)_all.deb

dep:
	dep ensure

all: build build-deb

clean-build:
	@rm -rf build $(NAME)

clean-deb:
	@rm -rf $(NAME)-$(VERSION) *.deb

clean: clean-deb clean-build
