# MailCage

Based on Mailhog.

## Why?
My company and I have been using `mailhog` for quite a while, but it's poorly maintained now.
I made `MailCage` on top of `mailhog` with the goal of actively maintaining it.

## Differences from Mailhog
- `Maildir` is replaced with `sqlite`
- One single repository
- Frontend is rewritten in React
- Improved logging
- Some bugs fixed

## Design
Current design is the simplest I could make :-) I'm waiting on a designer friend to make a new one :-)

## TODO
0. More readme and contributing guide.
1. Proper error handling
2. Mail downloads
3. Swagger
4. Tests
5. Search

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
