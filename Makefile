.PHONY: build
build:
	go build .

.PHONY: tags
tags:
	./docker-ce-tags diff-tags config.yml

branch:
	./docker-ce-tags branch config.yml
