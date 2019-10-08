# MailCage

Based on Mailhog.

## Features
1. Straight go, only one binary
2. One repo (not like mailhog)
3. Refactoring
4. maildir -> sqlite

## TODO
1. move tests from mailhog
3. ErrMessageNotFound sqlite

## TODO
1. Search
2. Releasing
3. Mail downloads?
4. UI and API Auth
5. Dockerfile
6. Swagger

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
