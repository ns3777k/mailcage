# MailCage

Based on Mailhog.

## Features
1. Straight go, only one binary
2. One repo (not like mailhog)
3. Refactoring
4. maildir -> sqlite

## TODO
2. smtp graceful shutdown?
3. mailhog copyrights
4. move tests from mailhog
5. websockets
8. ui
9. sqlite context
10. ErrMessageNotFound sqlite

## Limitations for now
1. No search
2. No releasing
3. No download
4. UI Auth

```shell script
/home/nsafonov/go/src/github.com/mailhog/mhsendmail/mhsendmail test@mailhog.local <<EOF
To: Test <test@mailhog.local>
From: Nikita <ns3777k@gmail.com>
Subject: Test message

Some content!
EOF
```
