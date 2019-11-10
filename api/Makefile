NAME:=$(shell basename `git rev-parse --show-toplevel`)
HASH:=$(shell git rev-parse --verify --short HEAD)

all: dockerbuild

clean:
	rm -rf pkg bin

up: dockerbuild down
	docker run -d --name $(NAME)_service -p 8080:80 $(NAME)

dockerbuild:
	GOOS=linux GOARCH=amd64 go build -ldflags "-X main.version=$(HASH)" -o $(NAME)
	docker build -t $(NAME) .

build:
	go build -ldflags "-X main.version=$(HASH)" -o $(NAME)

down:
	(docker stop $(NAME)_service || true) && (docker rm $(NAME)_service || true)
