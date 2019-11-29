# padl - CLI Reference

The padl Command Line Interface (CLI) is the client component of padl. It must be pointed to a running padl server in order to work properly. Below you will find a guide on building the CLI, as well as a reference of the available commands.

## Contents

* [Set-Up](#setting-up-tool)
	* [Build the CLI](#building-the-cli)
	* [Configure the CLI](#configuring-the-cli)
* [Commands](#commands-reference)
	* [Accounts](#account-commands)
	 	* [create](#account-creation)
	 	* [login](#account-login)
	 	* [show](#account-show)
	 	* [rotate-key](#account-key-rotation)
	* [Projects](#project-commands)
	 	* [create](#project-creation)
	 	* [get](#project-description)
	 	* [list](#project-list)
	 	* [delete](#project-deletion)
	* [Users](#user-commands)
	 	* [add](#user-addition)
	 	* [remove](#user-removal)
	* [Service Accounts](#service-account-commands)
	 	* [create](#service-account-creation)
	 	* [remove](#service-account-removal)
	* [Secrets](#secret-commands)
	 	* [set](#set-a-secret)
	 	* [show](#see-a-secret)
	 	* [remove](#delete-a-secret)
	* [Padlfile](#padlfile-commands)
	 	* [pull](#synchronize-a-padlfile-with-a-padl-server)

* [Feed Your App Secrets](#passing-your-app-secrets)

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

The following commands deal with your padl user account.

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

The following commands deal padl projects

#### Project Creation

Create your first project by changing directory into your desired working directory (e.g. top level repo):

```
$ padl project create  --name demo-project --description "project for docs"
project demo-project initialized successfully!
```

Note that you may override the default project file (.padlfile) location with the ```--path``` flag

#### Project Description

To get a project by name you may use the ```padl project get``` command:

```
$ padl project get --project demo-project
+-------------+----------------------------------+
|    NAME     |           demo-project           |
| DESCRIPTION |         project for docs         |
|     KEY     | 49e9df18868c24225025558529a2188d |
|   MEMBERS   |   adrianosela@protonmail.com 2   |
+-------------+----------------------------------+
```

Note that the `--json` flag is available to print JSON formatted output instead

#### Project List

To get a list of all projects you are a member of, use the ```padl project list``` command:

```
$ padl project list
+--------------+--------------------------------+
|     NAME     |          DESCRIPTION           |
+--------------+--------------------------------+
| demo-project |        project for docs        |
|     mapp     |          project for           |
|              |  github.com/adrianosela/mapp   |
|    sslmgr    |          project for           |
|              | github.com/adrianosela/sslmgr  |
+--------------+--------------------------------+
```

Note that the `--json` flag is available to print JSON formatted output instead

#### Project Deletion

To delete a project, use the ```padl project delete``` command:

```
$ padl project delete --project sslmgr
project sslmgr deleted successfully!
```

### User Commands

The following commands deal with user account access to projects

#### User Addition

The ```padl project user add``` command adds a given user to a project:

```
$ padl project user add --project demo-project --email adrianosela@gmail.com --privilege 1
user adrianosela@gmail.com added to project demo-project successfully!
```

Privilege Levels: 

> 0 - READ ONLY: can only see a project
> 
> 1 - EDIT: can add and remove service accounts to the project
> 
> 2 - OWNER: can add and remove other users to the project

#### User Removal

The ```padl project user remove``` command removes a given user from a project:

```
$ padl project user remove --project demo-project --email adrianosela@gmail.com
user adrianosela@gmail.com removed from project demo-project successfully!
```

### Service Account Commands

The following commands deal with service account access to projects

#### Service Account Creation

The ```padl project service-account create``` command creates a service account:

```
$ padl project service-account create --project demo-project --name deploybot
---------------------- IMPORTANT NOTE ----------------------
>> Both the RSA private key and and auth token are secret <<
>> If either is disclosed you MUST delete the svc account <<
>> If both disclosed - your secrets have been compromised <<
------------------------------------------------------------

SERVICE ACCOUNT PRIVATE KEY:
-----BEGIN RSA PRIVATE KEY-----
IJKAIBAAKCAgEA2FLMK0qwyuCjUtgUa9skT
5m6IyWFX48IOLgyRypUvfEXTjjd\nlypctK
...				...				...
nPdfoqGGUN3clCBt6x1eBFkvq1Wd2XuWkHc
gHPyAJtw0MdjKN8EIUbTdAMfzjbPD1d3tod	
-----END RSA PRIVATE KEY-----

SERVICE ACCOUNT AUTH TOKEN:
eyJhbG7pTd7RBEUaj9Eyct4Tb670KQslzg69-78f6uX_QP8-0qxVG5OERlbkag1rsKgrPAbRZukPU8ilfu6K2Kt-aXcuI6QxyaX0PRgiOFRHpq4B1WLbjCJR1KKXsh--jmSQw
```

Note that the `--json` flag is available and can be used as follows:

```
18:13 $ padl project service-account create --project demo-project --name deploybot2 --json | jq -r .
{
  "private_key": "-----BEGIN RSA PRIVATE KEY-----\nMIIJKAIBAAKCAgEA2FLMK0qwyuCjUtgUa9skT5m6IyWFX48IOLgyRypUvfEXTjjd\nlypctKnPdfoqGGUN3clCBt6x1eBFkvq1Wd2XuWkHcgHPyAJtw0MdjKN8EIUbTdAMfzjbPD1d3toda3EJdOhNBJaE2XDUIjO+WfNAkFU61DjYjBnLaZ\nn91rSwaJDEcL53fwJo6H0Iz5xPE7Aulbm7Q0yae5enytnzI1RLJn1Ok2vII=\n-----END RSA PRIVATE KEY-----\n",
  "jwt": "eyJhbGciOiJSUzUxMiIjlkOGYzYjhmM2MzOTg2NTg5N2U5MzNjNjNlNWMxYjdkIiwidHlwIjoiSldUIn0.eyJhdWQiOiJkZWNyeXB0IiwiZXhwIjoxNjA2NTI5OTIyLCJqdGkiOiJjZGg"
}
```

#### Service Account Removal

The ```padl project service-account remove``` command removes a service account from a project:

```
$ padl project service-account remove --project demo-project --name deploybot2
service account deploybot2 removed from project demo-project successfully!
```

### Secret Commands

The following commands deal with secrets in a padlfile

#### Set a Secret

Set a secret in a padlfile with the ```padl file secret set``` command:

```
$ padl file secret set --name MONGODB_CONNSTR --secret "mongo://user:supersecretstuff@mymongoinstance.com"
padlfile updated!
```

#### See a Secret

To decrypt and see a secret in plaintext, use the ```padl file secret show``` command:

```
$ padl file secret show --name MONGODB_CONNSTR
mongo://user:supersecretstuff@mymongoinstance.com
```

#### Delete a Secret

To delete a secret from a padlfile, you may use the ```padl file secret remove``` command:

```
$ padl file secret remove --name MONGODB_CONNSTR
padlfile updated!
```

### Padlfile Commands

The following commands deal with a padlfile

#### Synchronize a padlfile With a Padl Server

Use the ```padl file pull``` command to pull any new encryption key ids from the server, e.g. to include a new user or service account in the encryption:

```
$ padl file pull
padlfile updated!
```

## Passing Your App Secrets

The padl CLI must be installed in the host machine

By passing the command/executable of your application as an argument to the ```padl run``` command, you can run your application with padl decrypted secrets in the runtime environment.

For example, here we run the unix command `env` within the padl wrapper, and we see that the environment has the secret variable MONGODB_CONNSTR we set earlier when [setting a secret](#set-a-secret):

```
18:46 $ padl run env
MONGODB_HOST=127.0.0.1
TERM_PROGRAM=iTerm.app
SHELL=/bin/bash
MONGODB_PORT=27017
MONGODB_CONNSTR=mongo://user:supersecretstuff@mymongoins.com
```

For a detailed walkthrough head over to [demos/simple](https://github.com/adrianosela/padl/tree/master/demos/simple).

