.PHONY: build
build:
	go build .

.PHONY: tags
tags:
	./docker-ce-tags diff-tags config.yml

commits:
	./docker-ce-tags commits config.yml
