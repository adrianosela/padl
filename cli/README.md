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

```
$ make
go build -ldflags "-X main.version=0.1.0-c1838a7" -o padl
cp padl /usr/local/bin

$ padl --version
padl version 0.1.0-c1838a7
```

### Configuring the CLI

To point the padl CLI to a padl server, use the `padl config set` command, specifying the server URL with the `--url` flag as follows:

```
$ padl config set --url https://padl.adrianosela.com
padl configuration set successfully!
```

## Commands Reference

Below is usage information of all available commands in CLI. If you are comfortable with command line tools, you might instead want to use the CLI's built-in help menu available by appending the `--help` flag to any command or subcommand.

### Account Commands

The following commands deal with padl user account.

#### Account Creation

You can start off by creating an account with the `padl account create` command:

```
$ padl account create
Enter your email:
adrianosela@protonmail.com
Enter your password:
registered user adrianosela@protonmail.com successfully!
```

Note that one may skip the interactive prompt by populating the `--email` and `--password` flags. However, not providing the `--password` flag will use the "silent" prompt to hide your password.

#### Account Login

Log into your padl account through the `padl account login` command:

```
$ padl account login
Enter your email:
adrianosela@protonmail.com
Enter your password:
user adrianosela@protonmail.com logged in successfully!
```
Note that one may skip the interactive prompt by populating the `--email` and `--password` flags. However, not providing the `--password` flag will use the "silent" prompt to hide your password.


#### Account Show

To view the claims in your access token (...and under the hood make a call to check their validity) you may use the `padl account show` command:

```
$ padl account show
+-----+--------------------------------------+
| aud |                                  api |
| iss |                 padl.adrianosela.com |
| sub |           adrianosela@protonmail.com |
| iat |                           1574988578 |
| exp |                           1575031778 |
| jti | 2d351405-e3d7-468a-826d-d342faf552fe |
+-----+--------------------------------------+
```
Note that the `--json` flag is available for JSON output.

#### Account Key Rotation

Rotate your user private key with the ```padl account rotate-key``` command:

```
$ padl account rotate-key
rotated user key successfully!
```
Important Considerations: 
> Any padlfile encrypted with your old key can still be decrypted with that key if and only if the holder of the key has the user's active session token. (Or else secrets theft will be halted by the need to provide padl login credentials)
> 
> If your machine was compromised while you had an active padl session token, your secrets have been compromised, and they must also be rotated
>
> Note that when rotating a key, you will still need access to the old key if you still want to decrypt secrets in existing padlfiles. Otherwise have another user update the padlfile to include your new key ID, (and newly encrypted secrets), and push to version control

### Project Commands

// TODO

### User Commands

// TODO

### Service Account Commands

// TODO

### Secret Commands

// TODO