DIR=deployment/docker
COMPOSE=${DIR}/docker-compose.yaml

DIND_PREFIX ?= $(HOME)

PREFIX=$(shell echo $(PWD) | sed -e s:$(HOME):$(DIND_PREFIX):)

UID=$(shell whoami) 

ifeq ($(CACHE_PREFIX),)
	CACHE_PREFIX=/tmp
endif

# untuk mendapatkan go yang sudah pernah dimasukan sebelumnya
# untuk mendapatkan apk yang sudah diinstall didalam alpine sebelumnya

test: 
	docker run \
		--network api_default \
		-v $(CACHE_PREFIX)/cache/go:/go/pkg/mod \
		-v $(CACHE_PREFIX)/cache/apk:/etc/apk/cache \
		-v $(PREFIX)/deployment/docker/build:/build \
		-v $(PREFIX)/:/src \
		-v $(PREFIX)/migrations:/migrations \
		-v $(PREFIX)/scripts/test.sh:/test.sh \
		-e UID=$(UID) \
		golang:1.13-alpine /test.sh 

network: 
	docker network create -d bridge api_default; /bin/true

