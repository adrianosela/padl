NAME:=$(shell basename `git rev-parse --show-toplevel`)
HASH:=$(shell git rev-parse --verify --short HEAD)

all: dockerbuild

clean:
	rm -rf pkg bin

deploy: decrypt-conf-prod
	gcloud app deploy

up: dockerbuild down
	docker run -d --name $(NAME)_service -p 8080:80 $(NAME)

dockerbuild:
	GOOS=linux GOARCH=amd64 go build -ldflags "-X main.version=$(HASH)" -o $(NAME)
	docker build -t $(NAME) .

run: build
	./$(NAME)

build:
	go build -ldflags "-X main.version=$(HASH)" -o $(NAME)

decrypt-conf-prod:
	sops -d api/config/encrypted-config-prod.yaml > api/config/config.yaml

decrypt-conf-dev:
	sops -d api/config/encrypted-config-dev.yaml > api/config/config.yaml

down:
	(docker stop $(NAME)_service || true) && (docker rm $(NAME)_service || true)
