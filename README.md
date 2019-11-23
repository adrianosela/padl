# padl - secrets management as-a-service

[![Go Report Card](https://goreportcard.com/badge/github.com/adrianosela/padl)](https://goreportcard.com/report/github.com/adrianosela/padl)
[![Documentation](https://godoc.org/github.com/adrianosela/padl?status.svg)](https://godoc.org/github.com/adrianosela/padl)
[![license](https://img.shields.io/github/license/adrianosela/padl.svg)](https://github.com/adrianosela/padl/blob/master/LICENSE)
[![Generic badge](https://img.shields.io/badge/UBC-CPEN442-RED.svg)](https://blogs.ubc.ca/cpen442/about/)

Padl is an attempt at simplyfing secrets management for inexperienced developers and teams looking to quickly prototype solutions while spending very little time at securing secrets.

Our design is inspired by two popular secrets management tools:

* Mozilla's [SOPS](https://github.com/mozilla/sops) - except our keys are purely RSA (not PGP), we do not (yet) have cloud KMS integrations, and we require users connect to a server-side component
* CyberArk's [Conjur](https://github.com/cyberark/conjur) - except our server **never** sees plaintext secrets, as obfuscation occurs at the client side through splitting with Shamir's Secret Sharing algorithm

The goal is to create a secrets management solution with minimal set-up time, high-level abstractions, and great usability.


## Contents

* [Before You Begin](#prereading)
* [Running the API](#getting-started-with-the-api)
* [Command Line Interface (client)](#cli-reference)
* [Disclaimer](#disclaimer-and-recommendations)

## Prereading

Padl is composed of two main components, a REST API and a Command Line Interface (CLI). All software provided in this repository is and will always remain free and Open-Source as per our [MPL-2.0 License](https://github.com/adrianosela/padl/blob/master/LICENSE).

A globally available API instance is provided at **https://padl.adrianosela.com**, use of this API will remain **free at reasonable request rates** (which are yet to be defined).

You can either use the global API, or run your own instance(s) on-prem as described by the rest of this document.

The CLI **must** be able to communicate with a Padl API and can not be used independently without a running API.

## Getting Started with the API

Start by cloning this git repository:
```git clone https://github.com/adrianosela/padl```

### Configuration

The API reads configuration from a config.yaml file in the [/api/config](https://github.com/adrianosela/padl/blob/master/api/config) subdirectory. To be able to run the API, all variables with a \`yaml\` tag in the struct in [/api/config/config.go](https://github.com/adrianosela/padl/blob/master/api/config/config.go) must be defined in the yaml file.

### Build the API

The API can be built with the `go build` command or with the Makefile target:

```
$ make build
go build -ldflags "-X main.version=26f5980" -o padl
```

To build a binary for a specific operating system, you may populate the GOOS and GOARCH environment variables when running `go build`:

Example for Linux: ```GOOS=linux GOARCH=amd64 go build```

### Containerize the API (optional)

The Padl API can be containerized using the Dockerfile present in the top level of the repository. 

Note that the base image is Linux-based and thus you must build a Linux binary for it to be ran within the container.

To cross-compile a binary for a Linux OS and then build the Docker container you may use the Makefile target:

```
$ make dockerbuild
GOOS=linux GOARCH=amd64 go build -ldflags "-X main.version=26f5980" -o padl
docker build -t padl .
Sending build context to Docker daemon  53.37MB
Step 1/5 : FROM alpine:latest
 ---> cdf98d1859c1
Step 2/5 : RUN apk add --update bash curl && rm -rf /var/cache/apk/*
 ---> Using cache
 ---> 884e3b9ebb3a
Step 3/5 : COPY . .
 ---> 420d4ef81582
Step 4/5 : EXPOSE 80
 ---> Running in f3752836d34a
Removing intermediate container f3752836d34a
 ---> 1ce25429a488
Step 5/5 : CMD ["./padl"]
 ---> Running in 4f33b6664024
Removing intermediate container 4f33b6664024
 ---> 1e9790a9190d
Successfully built 1e9790a9190d
Successfully tagged padl:latest
```

### Run the API

If you wish to run the API on your local machine (and not within a Docker container), you may run the built binary directly:

```
$ ./padl
2019/11/22 12:43:12 [info] successfully connected to MongoDB
2019/11/22 12:43:13 [info] successfully connected to MongoDBKeystore
```

If you wish to run the API within a Docker container, you may use the Makefile target:

```
$ make up
docker run -d --name padl_service -p 8080:80 padl
ef6640424b08a98677d4d4fec8dd134e9b3df054f3e9b45255ab6f0b27928bdc
```

## CLI Reference

For detailed usage information on the CLI, head over to our official [CLI reference](https://github.com/adrianosela/padl/tree/master/cli/README.md)

## Disclaimer and Recommendations

Although the design is provably secure and uses strong encryption to enforce policies, **we do not recommend it be used in any production environment**. Tools out there are written by security professionals and have undergone fierce scrutiny by the security community as a whole.

As per our [MPL-2.0 License](https://github.com/adrianosela/padl/blob/master/LICENSE), this software is provided as-is and the creators/developers are free from all liability.
