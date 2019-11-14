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
15:25 $ padl account login
Enter your email:
adrianosela@protonmail.com
Enter your password:
user adrianosela@protonmail.com logged in successfully!
```

After logging in, your padl configuration file (`~/.padl/config` by default) will contain your fresh access token

## Projects

// TODO

## Secrets

// TODO
