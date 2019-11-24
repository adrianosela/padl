# padl - CLI Reference

The padl Command Line Interface (CLI) is the client component of padl. It must be pointed to a running padl server in order to work properly. Below you will find a guide on building the CLI, as well as a reference of the available commands.

## Contents

* [Set-Up](#setting-up-tool)
	* [Build the CLI](#building-the-cli)
	* [Configure the CLI](#configuring-the-cli)
* [Commands](#commands-reference)
	* [Accounts](#account-commands)
	* [Projects](#project-commands)
	* [Users](#user-commands)
	* [Service Accounts](#service-account-commands)
	* [Secrets](#secret-commands)

## Setting Up Tool

### Building the CLI

Clone the padl repository and change directory into `/cli`. The Makefile target `make build` will build a padl binary in the current directory. You may then move the binary to the binaries path of your choice.

For UNIX systems (i.e. Linux, MacOS) you can build with `make`. This will move the built binary to the standard `/usr/local/bin` directory.

### Configuring the CLI

To point the padl CLI to a padl server, use the `padl config set` command, specifying the server URL with the `--url` flag as follows:

```
$ padl config set --url https://padl.adrianosela.com
padl configuration set successfully!
```

## Commands Reference

Below is usage information of all available commands in CLI. If you are comfortable with command line tools, you might instead want to use the CLI's built-in help menu available by appending the `--help` flag to any command or subcommand.

### Account Commands

// TODO

### Project Commands

// TODO

### User Commands

// TODO

### Service Account Commands

// TODO

### Secret Commands

// TODO