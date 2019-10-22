# MailCage

Based on [Mailhog](https://github.com/mailhog/MailHog).

- Download and run MailCage
- Configure your outgoing SMTP server
- View your outgoing email in a web UI
- Release it to a real mail server

Built with Go - MailCage runs without installation on multiple platforms.

## Overview

MailCage is an email testing tool for developers:

- Configure your application to use MailCage for SMTP delivery
- View messages in the web UI, or retrieve them with the JSON API
- Optionally release messages to real SMTP servers for delivery

## Running with docker

The most simple way to run the application:

```shell script
$ docker run --rm ns3777k/mailcage
```

## Running with pre-built binaries

Just [download the corresponding binary from releases](https://github.com/ns3777k/mailcage/releases) and run it :-)

## Why?
My company and I have been using `mailhog` for quite a while, but it's poorly maintained now.
I made `MailCage` on top of `mailhog` with the goal of actively maintaining it.

## Differences from Mailhog
- `Maildir` is replaced with `sqlite`
- One single repository
- Frontend is rewritten in React
- Improved logging
- Some bugs fixed

## TODO
0. More readme
1. Proper error handling
2. Mail downloads
3. Swagger
4. Tests
5. Search

## Mcsendmail

A fork of mhsendmail:

```shell script
./mcsendmail test@mailcage.local <<EOF
To: Test <test@mailcage.local>
From: Nikita <ns3777k@gmail.com>
Subject: Test message

Some content!
EOF
```
