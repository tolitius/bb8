#!/bin/bash

TOKEN=YUM

echo issuing a new token "$TOKEN"
echo

echo gen new account keys
bb gen-keys issuer
bb gen-keys distributor

echo funding accounts
bb fund $(cat issuer.pub)
bb fund $(cat distributor.pub)

echo changing trust and setting options
bb change-trust '{"source_account": "'$(cat distributor)'",
                  "code": "'$TOKEN'",
                  "issuer_address": "'$(cat issuer.pub)'"}' | xargs \
bb set-options  '{"home_domain": "dotkam.com",
                  "max_weight": 1}' | xargs \
bb sign '["'$(cat distributor)'"]' | xargs \
bb submit

echo funding distributor with new token
bb send-payment -s '{"from": "'$(cat issuer)'",
                     "to": "'$(cat distributor.pub)'",
                     "token": "'$TOKEN'",
                     "amount": "42000.00",
                     "issuer": "'$(cat issuer.pub)'"}'

echo "check your funds at: http://testnet.stellarchain.io/address/$(cat distributor.pub)"
