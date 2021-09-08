# aws-token-generator

Tool that request a token to be used with a MFA account in AWS.

# Installation

## Download
Download the binary file from _Github_ and move it to a directory added to your _PATH_.

```
$ cd /tmp
$ wget https://github.com/Gigigotrek/aws-token-generator/releases/download/latest/aws-token-generator
```
Or if you want to download an old version, 

```
$ cd /tmp
$ wget https://github.com/Gigigotrek/aws-token-generator/releases/download/{VERSION_TO_INSTALL}/aws-token-generator
```

## Move the binary file

```
$ sudo chmod +x /tmp/aws-token-generator
$ sudo mv /tmp/aws-token-generator /usr/local/bin/
```

# Usage

## Parameters

```
+-----------------------+--------+----------+--------------------------------------------+-----------------------------------------------------+
|   Input Variable      |  type  | Required |                  default                   |                     Description                     |
+-----------------------+--------+----------+--------------------------------------------+-----------------------------------------------------+
| token                 | string | yes      |                                            | Token provided by the MFA device.                   |
| profile               | string | yes      |                                            | Local profile to use to connect to AWS.             |
| expiration            | int64  | no       | 3600                                       | Expiration time for the token.                      |
| region                | string | no       | eu-west-1                                  | Region to set the environment profile.              |
+-----------------------+--------+----------+--------------------------------------------+-----------------------------------------------------+
```

## Example
```
aws-token-generator -t 327308 --profile stage -e 3000
```