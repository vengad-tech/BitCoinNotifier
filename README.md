## Introduction
Small notifier that periodically notifies , how much profit/loss incurred that is invested in bitcoin.

## Build
You need Go language compiler installed to compile and use this utility
first you need terminal-notifier command line utility installed on your mac os to show notification , to istall it use
```
[sudo] gem install terminal-notifier
```
Then install Go language , After installing Go , run
```
go build notifier.go
```
you will get notifier binary file created , then run it
```
notifier -btaddress="your bitcoin address" -buyingprice="your buying price in USD"
Example:
./notifier -btaddress="14ahC1Dd9oEB9guYG2G7B1pK3ZCTyJFPJM" -buyingprice=64.36
```
PS: Add notifier to your crontab to execute periodically , so you will get periodic bitcoin alert.
## Why so complex ??
Just created this utility to try the Go language , nothing else.!!