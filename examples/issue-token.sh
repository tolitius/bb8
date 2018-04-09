#!/bin/bash

if [ "$#" -ne 1 ]; then
  echo "usage:   $0 student name" >&2
  echo "example: $0 john_malkovich" >&2
  exit 1
fi

STUDENT_NAME=$1

TOKEN=TIL
ISSUER=GBKOPETTWWVE7DM72YOVQ4M2UIY3JCKDYQBTSNLGLHI6L43K7XPDROID
INFLATION_POOL=GCCD6AJOYZCUAQLX32ZJF2MKFFAUJ53PVCFQI3RHWKL3V47QYE2BNAUT

KEYS_HOME=~/.stellar
STUDENT=$KEYS_HOME/$STUDENT_NAME

## assuming university has a trusline for $TOKEN by the $ISSUER
## if not, another call to "bb change-trust" is needed
UNIVERSITY=path-to-the-real-key-that-has-XLMs

echo generating a fresh pair of keys
bb gen-keys $STUDENT

echo creating a new student account on a Stellar network
bb create-account -s '{"source_account":"'$(cat $UNIVERSITY)'",
                       "new_account":"'$(cat $STUDENT.pub)'",
                       "amount":"3.0"}'

echo setting a trustline for $TOKEN to $STUDENT_NAME and setting federation/inflation options
bb change-trust '{"source_account": "'$(cat $STUDENT)'",
                  "code": "'$TOKEN'",
                  "issuer_address": "'$ISSUER'"}' | xargs \
bb set-options  '{"home_domain": "dotkam.com",
                  "inflation_destination": "'$INFLATION_POOL'"}' | xargs \
bb sign '["'$(cat $STUDENT)'"]' | xargs \
bb submit

echo sending 42,000 $TOKEN to $STUDENT_NAME
bb send-payment -s '{"from": "'$(cat $UNIVERSITY)'",
                     "to": "'$(cat $STUDENT.pub)'",
                     "token": "'$TOKEN'",
                     "amount": "42000.00",
                     "issuer": "'$ISSUER'"}'

echo checking $STUDENT_NAME balance
bb load-account $(cat $STUDENT.pub) | jq '.balances'
