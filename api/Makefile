NAME:=$(shell basename `git rev-parse --show-toplevel`)
HASH:=$(shell git rev-parse --verify --short HEAD)

all: dockerbuild

clean:
	rm -rf pkg bin

deploy: decrypt-conf-prod
	gcloud app deploy

up: dockerbuild down
	docker run -d --name $(NAME)_service -p 8080:80 $(NAME)

dockerbuild: decrypt-conf-dev
	GOOS=linux GOARCH=amd64 go build -ldflags "-X main.version=$(HASH)" -o $(NAME)
	docker build -t $(NAME) .

build: decrypt-conf-dev
	go build -ldflags "-X main.version=$(HASH)" -o $(NAME)

decrypt-conf-prod:
	sops -d config/encrypted-config-prod.yaml > config/config.yaml

decrypt-conf-dev:
	sops -d config/encrypted-config-dev.yaml > config/config.yaml

down:
	(docker stop $(NAME)_service || true) && (docker rm $(NAME)_service || true)
