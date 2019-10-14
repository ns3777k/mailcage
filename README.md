# MailCage

Based on Mailhog.

## Features
1. Straight go, only one binary
2. One repo (not like mailhog)
3. Refactoring
4. maildir -> sqlite

## TODO
1. Error handling

## Not yet
1. Search
2. Mail downloads
3. Swagger
4. Tests

## Mcsendmail

A fork of mhsendmail:

```shell script
./mcsendmail test@mailhog.local <<EOF
To: Test <test@mailhog.local>
From: Nikita <ns3777k@gmail.com>
Subject: Test message

Some content!
EOF
```
