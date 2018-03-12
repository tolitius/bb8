#!/bin/bash

## in order to run it, "jq" needs to be installed (i.e. brew install jq, apt-get install jq, etc.)
## by default would run against testnet

## TODO: stand up a test horizon in docker and run it against it

seed_name=foo
pub_name=$seed_name.pub
tmp=/tmp
pub_file=$tmp/$pub_name
seed_file=$tmp/$seed_name

bb=./bb8

## get the latest build
rm ./bb8
go build

## TEST create key files
echo TEST: create key files

$bb gen-keys $seed_file

if [ ! -f "$pub_file" ]; then
    echo "[FAIL] expected: $pub_file public key file, but found no such thing"
	exit 1
fi

if [ ! -f "$seed_file" ]; then
    echo "[FAIL] expected: $seed_file seed file, but found no such thing"
	exit 1
fi

## set pub and seed to use them going forward
pub=$(cat $pub_file)
seed=$(cat $seed_file)

## TEST fund on testnet
echo TEST: fund a test account

$bb fund $(cat $pub_file)
balance=`$bb load-account $(cat $pub_file) | jq '[.balances[0].balance][]'`

expected_balance="\"10000.0000000\""

if [ "$balance" != "$expected_balance" ]; then
    echo "[FAIL] could not fund a test account, expecting $expected_balance, but got $balance instead"
	exit 1
fi


## TEST create account
echo TEST: create account

$bb gen-keys $tmp/bar
$bb create-account -s '{"source_account":"'$seed'",
                        "new_account":"'$(cat $tmp/bar.pub)'",
						"amount":"1.5"}'


echo "all tests... [PASS]"

## cleanup

# rm $seed_file $pub_file $tmp/bar*
