# MailCage

Based on Mailhog.

## Features
1. Straight go, only one binary
2. One repo (not like mailhog)
3. Refactoring
4. maildir -> bbolt

## TODO
1. clean modules
2. smtp graceful shutdown?
3. mailhog copyrights
4. move tests from mailhog
5. websockets
6. replace logger in mailhog's smtp stuff
7. release servers
8. file storage
9. ui
10. sqlite context
11. ErrMessageNotFound sqlite

## Limitations for now
1. No search
2. No releasing

```shell script
/home/nsafonov/go/src/github.com/mailhog/mhsendmail/mhsendmail test@mailhog.local <<EOF
To: Test <test@mailhog.local>
From: Nikita <ns3777k@gmail.com>
Subject: Test message

Some content!
EOF
```
