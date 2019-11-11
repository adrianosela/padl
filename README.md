# padl - secrets management as-a-service

padl is an attempt at simplyfing secrets management for inexperienced developers and teams looking to quickly prototype solutions while spending very little time at securing secrets. The goal is to create a secrets management tool with great usability.

## Contents

### [Configure CLI](#configuration)
* [Server Establishment](#setting-padl-server)
* [Verify Configuration](#verify-configuration-file)

### [Accounts](#accounts)

### [Projects](#projects)

### [Secrets](#secrets)

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

// TODO

## Projects

// TODO

## Secrets

// TODO
