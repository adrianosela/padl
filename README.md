# padl - secrets management as-a-service

padl is an attempt at simplyfing secrets management for inexperienced developers and teams looking to quickly prototype solutions while spending very little time at securing secrets. The goal is to create a secrets management tool with great usability.

## Contents

### [Accounts](#accounts)
* [New User Registration](#new-user-registration)
* [User Key Rotation](#user-key-rotation)

### [Projects](#projects)

### [Secrets](#secrets)

## Accounts

### New User Registration:

- generate your personal PGP key following [these steps](http://irtfweb.ifa.hawaii.edu/~lockhart/gpg/). To be accepted by padl, the key must satisfy the following constraints:
	- 4096-bit RSA
	- you have access to the email address on the key

```
gpg --key-gen
``` 

- create an account using the padl CLI
	- you will receive a confirmation email which must be confirmed within 24 hours, or the account will be purged

```
padl account create --pub ${PATH_TO_YOUR_PGP_PUBLIC_KEY}
```

- once you have confirmed your account via email, you will be able to manage projects and secrets on padl


### User Key Rotation:

- generate your new personal PGP key following [these steps](http://irtfweb.ifa.hawaii.edu/~lockhart/gpg/). To be accepted by padl, the key must satisfy the following constraints:
	- 4096-bit RSA
	- you have access to the email address on the key

```
padl account rotate --old ${PATH_TO_OLD_PUB} --new ${PATH_TO_NEW_PUB}
```

- you will be emailed for confirmation; your new key will not be in place until you have confirmed the key rotation request

## Projects

// TODO

## Secrets

// TODO