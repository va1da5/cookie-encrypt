.PHONY: build
build:
	docker build -t va1da5/cookie-encrypt .

.PHONY: push
push:
	docker push va1da5/cookie-encrypt:latest