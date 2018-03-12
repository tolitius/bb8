#!/bin/bash

## in order to run it, "jq" needs to be installed (i.e. brew install jq, apt-get install jq, etc.)
## by default would run against testnet

## TODO: stand up a test horizon in docker and run test against it instead of testnet

seed_name=foo
pub_name=$seed_name.pub
tmp=/tmp
pub_file=$tmp/$pub_name
seed_file=$tmp/$seed_name

bb=./bb8

## get the latest build
rm ./bb8
go build

assert_balance() {

    pkey_file=$1
	expected_balance=\"$2\"
	fail_msg=$3
	asset=$4

	if [ -z "$asset" ]; then
		balance=`$bb load-account $(cat $pkey_file) | jq '[.balances[0].balance][]'`
	else
		balance=`$bb load-account $(cat $pkey_file) | jq '.balances[] | select(.asset_code == "'$4'").balance'`
	fi

	if [ "$balance" != "$expected_balance" ]; then
		echo "[FAIL] $fail_msg, expecting $expected_balance, but got $balance instead"
		exit 1
	fi
}

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

assert_balance $pub_file "10000.0000000" "could not fund account"

## TEST create account
echo TEST: create account

$bb gen-keys $tmp/bar
$bb create-account -s '{"source_account":"'$seed'",
                        "new_account":"'$(cat $tmp/bar.pub)'",
						"amount":"1.5"}'

assert_balance $tmp/bar.pub "1.5000000" "could not create account"

## TEST change trust
echo TEST: change trust

$bb gen-keys $tmp/xyz
$bb fund $(cat $tmp/xyz.pub)
$bb change-trust -s '{"source_account": "'$(cat $tmp/xyz)'",
                      "code": "XYZ",
                      "issuer_address": "'$pub'"}'

assert_balance $tmp/xyz.pub "0.0000000" "could change trust" "XYZ"

## TEST send payment
echo TEST: send payment

$bb send-payment -s '{"from": "'$seed'",
                      "to": "'$(cat $tmp/xyz.pub)'",
                      "token": "XYZ",
                      "amount": "42.0",
                      "issuer": "'$pub'"}'

assert_balance $tmp/xyz.pub "42.0000000" "could change trust" "XYZ"

echo "all tests... [PASS]"

## cleanup

# rm $seed_file $pub_file $tmp/bar*
