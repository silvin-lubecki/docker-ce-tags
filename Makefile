.PHONY: build
build:
	go build .

.PHONY: tags
tags:
	./docker-ce-tags diff-tags config.yml

.PHONY: branch
branch:
	./docker-ce-tags branch config.yml

.PHONE: cherry-pick
cherry-pick:
	go run cherrypick/main.go commits