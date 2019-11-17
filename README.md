# padl - secrets management as-a-service

[![Go Report Card](https://goreportcard.com/badge/github.com/adrianosela/padl)](https://goreportcard.com/report/github.com/adrianosela/padl)
[![Documentation](https://godoc.org/github.com/adrianosela/padl?status.svg)](https://godoc.org/github.com/adrianosela/padl)
[![license](https://img.shields.io/github/license/adrianosela/padl.svg)](https://github.com/adrianosela/padl/blob/master/LICENSE)
[![Generic badge](https://img.shields.io/badge/UBC-CPEN412-RED.svg)](https://blogs.ubc.ca/cpen442/about/)

padl is an attempt at simplyfing secrets management for inexperienced developers and teams looking to quickly prototype solutions while spending very little time at securing secrets. The goal is to create a secrets management tool with great usability.

## Contents

### [Config](#configuration)
* [Server Establishment](#setting-padl-server)
* [Verify Configuration](#verify-configuration-file)

### [Account](#accounts)
* [Create / Register](#create-a-new-account)
* [Login](#login-to-Your-Account)

### [KMS](#keys)
* [Get Public Key](#get-a-public-key)
// TODO

### [Project](#projects)

## Configuration:

<b>Note:</b> All configuration commands can have the default configuration file path (`~/.padl`) overriden with the `--path` flag

### Setting Padl Server:

Point your CLI to the padl server:

```
$ padl config set --url ${PADL_SERVER_URL}
```

### Verify Configuration File:

Verify your CLI can find the configuration file correctly:

```
$ padl config show
+----------+---------------------------+
| HOST_URL | https://api.padl.com:8080 |
+----------+---------------------------+
```

## Accounts:

### Create a new account:

Create an account by providing your email and choice of password:

```
$ padl account create --email adrianosela@protonmail.com --password @V3rYs3cuR3Pa$$w0rd
registered user adrianosela@protonmail.com successfully!
```

<b>Note:</b> both of `--email` and `--password` are optional and the CLI will prompt if not given

```
$ padl account create
Enter your email:
adrianosela@protonmail.com
Enter your password:
registered user adrianosela@protonmail.com successfully!
```

After registering, your padl configuration directory (`~/.padl` by default) will contain your new account's private key.

### Login to Your Account:

```
$ padl account login --email adrianosela@protonmail.com --password @V3rYs3cuR3Pa$$w0rd
user adrianosela@protonmail.com logged in successfully!
```

<b>Note:</b> both of `--email` and `--password` are optional and the CLI will prompt if not given

```
$ padl account login
Enter your email:
adrianosela@protonmail.com
Enter your password:
user adrianosela@protonmail.com logged in successfully!
```

After logging in, your padl configuration file (`~/.padl/config` by default) will contain your fresh access token

## Keys

### Get a Public Key

```
$ padl kms public --id 7852b5df86cdb7c45780f8aa0cd12f44 --json | jq -r .
{
  "id": "7852b5df86cdb7c45780f8aa0cd12f44",
  "pem": "-----BEGIN RSA PUBLIC KEY-----\nMIIBCgKCAQEAxMI/lQfvoLt3OMy0v9OIorq+Kdng3AKZ8feMOMefLmDvtGqnRvRL\nat3nQq+3mqyVlG/LFefliDnhEXqP1hjH64CEUJEpYwoDTwW7apWs3+T/1o492AGs\nLmeyAgU+SpWnTp4qD4pvbe4QEKTFXMRi3NLJwX9QD7Q7z1+pf6laNwVFdeoVKkDR\nrjFjIpCqmtkpxdDmE9aWeyCTAo4dXBMWu8zRzlIKL936OiqP4v0RSu5l3Y/fASSK\nayUf4qF+FhkT2zPohiRXdoX7trQEmF8oXL2R1xFhRxR9BckJWn0bf8NdFMhWJhTk\noSg0xWYnpU2tvf1fsKBfDcpyIwz6BvzcnwIDAQAB\n-----END RSA PUBLIC KEY-----\n"
}
```
// TODO

## Projects

// TODO
