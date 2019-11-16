# padl - secrets management as-a-service

[![Go Report Card](https://goreportcard.com/badge/github.com/adrianosela/padl)](https://goreportcard.com/report/github.com/adrianosela/padl)
[![Documentation](https://godoc.org/github.com/adrianosela/padl?status.svg)](https://godoc.org/github.com/adrianosela/padl)
[![Generic badge](https://img.shields.io/badge/UBC-CPEN 412-RED.svg)](https://blogs.ubc.ca/cpen442/about/)

padl is an attempt at simplyfing secrets management for inexperienced developers and teams looking to quickly prototype solutions while spending very little time at securing secrets. The goal is to create a secrets management tool with great usability.

## Contents

### [Config](#configuration)
* [Server Establishment](#setting-padl-server)
* [Verify Configuration](#verify-configuration-file)

### [Account](#accounts)
* [Create / Register](#create-a-new-account)
* [Login](#login-to-Your-Account)

### [KMS](#keys)
* [Create Key](#create-a-key)
* [Get Private Key](#get-a-private-key)
* [Get Public Key](#get-a-public-key)
* [Add User To Key](#add-user-to-private-key-ACL)
* [Remove User From Key](#remove-user-from-private-key-ACL)

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

### Create a Key

```
$ padl kms create --name "my key" --description "mock key" --json | jq -r .
{
  "id": "7852b5df86cdb7c45780f8aa0cd12f44",
  "name": "my key",
  "description": "mock key",
  "users": {
    "adrianosela@protonmail.com": 2
  },
  "pem": "-----BEGIN RSA PRIVATE KEY-----\nMIIEogIBAAKCAQEAxMI/lQfvoLt3OMy0v9OIorq+Kdng3AKZ8feMOMefLmDvtGqn\nRvRLat3nQq+3mqyVlG/LFefliDnhEXqP1hjH64CEUJEpYwoDTwW7apWs3+T/1o49\n2AGsLmeyAgU+SpWnTp4qD4pvbe4QEKTFXMRi3NLJwX9QD7Q7z1+pf6laNwVFdeoV\nKkDRrjFjIpCqmtkpxdDmE9aWeyCTAo4dXBMWu8zRzlIKL936OiqP4v0RSu5l3Y/f\nASSKayUf4qF+FhkT2zPohiRXdoX7trQEmF8oXL2R1xFhRxR9BckJWn0bf8NdFMhW\nJhTkoSg0xWYnpU2tvf1fsKBfDcpyIwz6BvzcnwIDAQABAoIBADrxejy6KOY84sVo\nRcmlpCwjx24gMEWYneen4iDsZFpvfb/Np5kQ/DriiTIoE9fJVfIm328Ljm6V8D/d\nOJPJzrJVSM4d/okF6eHVdMTEXAqivqXW7N31+k/YjrIeQf/z/zAFH9KSBTmodLWX\ntuxIhNlkaD6IVkKuGrDQFqYA5N7QNey8zSf/f6cFQhTE8pl0V8/Vd10l2lxl3Ft7\ntVaLTAh8pD8m3J+o5fqDBUbyNdTKdKuSy0TKflFXXutvO9oZOGkwnjIlf3ofX4Hl\n37pKkRYlelhfYRrSjv9I2Dy8STPcncMqHfw67/bSy5kIKHwIsqVa3h8mPtECoc6j\n5CjiOHkCgYEA2e2OXTm9wb3+1hxcVD1ivC1tQSQ3Vs8CTXbfJVVStj/Kj5tHgliK\nH+zEj7FT32T+QZvdORSG2AB6bVAmOMA57i3DljW8FXGtvQwX7dmfx00bEi1qsr6E\nUR1WmfwKvb6mYtMunanBJniUY1MqI3iceGWK8ZtdHRYEZgQw0cV2EvsCgYEA5yHw\nOWIJTl0QQZrQCNQONvwIHnrnnCtuSKhUc9tmkWJuFXOILhJtgXXfUiVvz95zJCOO\n2365YL2nG7yH+vsuvsrUMoOugV5f0aQ92IuulGMrhZ+j+faHGTmKUPDU8aTSAbRc\nVG5I6uzpsHWJW8huTtzXTwKf9E15kwwgYZ5fy60CgYBuXnBWaJLg1z+D4nMkOr6R\nfRQzBIt+THLnFofm2XJ9WItW9ZZevkad6oSWHYHTxss6IR0F9o5gQMXALPJelYQB\nS24d2fL6jUsnTkOkMy5HepZ2O0gpZHGQvyIH9GzgMfkEXd3i/YET4ceNEiZqNoBQ\nPWUD/eJHg8oQfJjY9H9bFwKBgAIGPBJkl2xGSGQqtPO+17kHkBKkRO8LOlYMk2DI\nZSeU0x4A+wpcQvVFUQVpKoeJjTydyxyFCZ6dSp9lkVNTa99j62Pd32NmrjQp2hjR\ncGAAVls/QLJpxFkmNd3rnhHXvbciG0TqCl10Yb+X5/IT2VN7f69DeJ8tJolxK79v\nIaupAoGAXDA6/uRS36OT4KqyNTPFGwcbKXvLVi88V6OT7SBa20TK198clklGQZax\nFmnPaCFD6CJcf8wg4non+l+kjKr6DNiLXV2qlEbZw64KGAeuY10zHELNLcSpMmTv\ncEK50RbaXJj2VqJX2oO0hGIl68fOMkERBjJG8D27ImLCg27gHbY=\n-----END RSA PRIVATE KEY-----\n"
}
```

<b>Note:</b> the default size of the RSA key is 2048 bits, this can be overriden with the `--bits` flag

### Get a Private Key

```
20:46 $ padl kms private --id 7852b5df86cdb7c45780f8aa0cd12f44 --json | jq -r .
{
  "id": "7852b5df86cdb7c45780f8aa0cd12f44",
  "name": "my key",
  "description": "mock key",
  "users": {
    "adrianosela@protonmail.com": 2
  },
  "pem": "-----BEGIN RSA PRIVATE KEY-----\nMIIEogIBAAKCAQEAxMI/lQfvoLt3OMy0v9OIorq+Kdng3AKZ8feMOMefLmDvtGqn\nRvRLat3nQq+3mqyVlG/LFefliDnhEXqP1hjH64CEUJEpYwoDTwW7apWs3+T/1o49\n2AGsLmeyAgU+SpWnTp4qD4pvbe4QEKTFXMRi3NLJwX9QD7Q7z1+pf6laNwVFdeoV\nKkDRrjFjIpCqmtkpxdDmE9aWeyCTAo4dXBMWu8zRzlIKL936OiqP4v0RSu5l3Y/f\nASSKayUf4qF+FhkT2zPohiRXdoX7trQEmF8oXL2R1xFhRxR9BckJWn0bf8NdFMhW\nJhTkoSg0xWYnpU2tvf1fsKBfDcpyIwz6BvzcnwIDAQABAoIBADrxejy6KOY84sVo\nRcmlpCwjx24gMEWYneen4iDsZFpvfb/Np5kQ/DriiTIoE9fJVfIm328Ljm6V8D/d\nOJPJzrJVSM4d/okF6eHVdMTEXAqivqXW7N31+k/YjrIeQf/z/zAFH9KSBTmodLWX\ntuxIhNlkaD6IVkKuGrDQFqYA5N7QNey8zSf/f6cFQhTE8pl0V8/Vd10l2lxl3Ft7\ntVaLTAh8pD8m3J+o5fqDBUbyNdTKdKuSy0TKflFXXutvO9oZOGkwnjIlf3ofX4Hl\n37pKkRYlelhfYRrSjv9I2Dy8STPcncMqHfw67/bSy5kIKHwIsqVa3h8mPtECoc6j\n5CjiOHkCgYEA2e2OXTm9wb3+1hxcVD1ivC1tQSQ3Vs8CTXbfJVVStj/Kj5tHgliK\nH+zEj7FT32T+QZvdORSG2AB6bVAmOMA57i3DljW8FXGtvQwX7dmfx00bEi1qsr6E\nUR1WmfwKvb6mYtMunanBJniUY1MqI3iceGWK8ZtdHRYEZgQw0cV2EvsCgYEA5yHw\nOWIJTl0QQZrQCNQONvwIHnrnnCtuSKhUc9tmkWJuFXOILhJtgXXfUiVvz95zJCOO\n2365YL2nG7yH+vsuvsrUMoOugV5f0aQ92IuulGMrhZ+j+faHGTmKUPDU8aTSAbRc\nVG5I6uzpsHWJW8huTtzXTwKf9E15kwwgYZ5fy60CgYBuXnBWaJLg1z+D4nMkOr6R\nfRQzBIt+THLnFofm2XJ9WItW9ZZevkad6oSWHYHTxss6IR0F9o5gQMXALPJelYQB\nS24d2fL6jUsnTkOkMy5HepZ2O0gpZHGQvyIH9GzgMfkEXd3i/YET4ceNEiZqNoBQ\nPWUD/eJHg8oQfJjY9H9bFwKBgAIGPBJkl2xGSGQqtPO+17kHkBKkRO8LOlYMk2DI\nZSeU0x4A+wpcQvVFUQVpKoeJjTydyxyFCZ6dSp9lkVNTa99j62Pd32NmrjQp2hjR\ncGAAVls/QLJpxFkmNd3rnhHXvbciG0TqCl10Yb+X5/IT2VN7f69DeJ8tJolxK79v\nIaupAoGAXDA6/uRS36OT4KqyNTPFGwcbKXvLVi88V6OT7SBa20TK198clklGQZax\nFmnPaCFD6CJcf8wg4non+l+kjKr6DNiLXV2qlEbZw64KGAeuY10zHELNLcSpMmTv\ncEK50RbaXJj2VqJX2oO0hGIl68fOMkERBjJG8D27ImLCg27gHbY=\n-----END RSA PRIVATE KEY-----\n"
}
```

### Get a Public Key

```
$ padl kms public --id 7852b5df86cdb7c45780f8aa0cd12f44 --json | jq -r .
{
  "id": "7852b5df86cdb7c45780f8aa0cd12f44",
  "pem": "-----BEGIN RSA PUBLIC KEY-----\nMIIBCgKCAQEAxMI/lQfvoLt3OMy0v9OIorq+Kdng3AKZ8feMOMefLmDvtGqnRvRL\nat3nQq+3mqyVlG/LFefliDnhEXqP1hjH64CEUJEpYwoDTwW7apWs3+T/1o492AGs\nLmeyAgU+SpWnTp4qD4pvbe4QEKTFXMRi3NLJwX9QD7Q7z1+pf6laNwVFdeoVKkDR\nrjFjIpCqmtkpxdDmE9aWeyCTAo4dXBMWu8zRzlIKL936OiqP4v0RSu5l3Y/fASSK\nayUf4qF+FhkT2zPohiRXdoX7trQEmF8oXL2R1xFhRxR9BckJWn0bf8NdFMhWJhTk\noSg0xWYnpU2tvf1fsKBfDcpyIwz6BvzcnwIDAQAB\n-----END RSA PUBLIC KEY-----\n"
}
```

### Add User to Private Key ACL

```
$ padl kms add-user --id 7852b5df86cdb7c45780f8aa0cd12f44 --email felipe@protonmail.com --json | jq -r .
{
  "id": "7852b5df86cdb7c45780f8aa0cd12f44",
  "name": "my key",
  "description": "mock key",
  "users": {
    "adrianosela@protonmail.com": 2,
    "felipe@protonmail.com": 0
  },
  "pem": "RSA PRIVATE KEY HIDDEN"
}
```

<b>Note:</b> the user will be added as a read-only user to the key by default. This can be changed by providing a 1 (Editor), or 2 (Owner) to the `--privilege` flag


### Remove User From Private Key ACL

```
$ padl kms remove-user --id 7852b5df86cdb7c45780f8aa0cd12f44 --email felipe@protonmail.com --json | jq -r .
{
  "id": "7852b5df86cdb7c45780f8aa0cd12f44",
  "name": "my key",
  "description": "mock key",
  "users": {
    "adrianosela@protonmail.com": 2
  },
  "pem": "RSA PRIVATE KEY HIDDEN"
}
```

## Projects

// TODO
