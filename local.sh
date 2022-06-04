#!/bin/bash


export WEB_CLIENT_KEY=$(gpg -d ~/creds/webclientaccount.rsa.asc | base64)
go run main.go "$@"
