# go-bitpay

Small cli utility for bitpay API

## install

`go install github.com/eknkc/go-bitpay`

## usage

```
NAME:
   go-bitpay - go bitpay pair|create|get

USAGE:
   go-bitpay [global options] command [command options] [arguments...]

VERSION:
   0.1.0

AUTHOR:
  Ekin Koc - <ekin@eknkc.com>

COMMANDS:
   pair   pair <pairCode> - Pair with BitPay API and write bitpay.json to current dir
   create create <price> - Create a new invoice for <price> usd
   get    get <invoiceId> - Get the status of an invoice
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h   show help
   --version, -v  print the version
```
