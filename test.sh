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

## this does exist on a public Stellar network, to test STELLAR_NETWORK switch
BB8_ACCOUNT=GBKOPETTWWVE7DM72YOVQ4M2UIY3JCKDYQBTSNLGLHI6L43K7XPDROID
BB8_FEDERATION_ADDRESS=bb8*dotkam.com

## making sure tests are not run against someone's / public network
unset STELLAR_NETWORK

## get the latest build
rm ./bb8
go build

assert_balance() {

    pkey_file=$1
	expected_balance=\"$2\"
	fail_msg=$3
	asset=$4

	if [ -z "$asset" ]; then
		balance=`$bb load-account $(cat $pkey_file) | jq '.balances[] | select(.asset_type == "native").balance'`
	else
		balance=`$bb load-account $(cat $pkey_file) | jq '.balances[] | select(.asset_code == "'$4'").balance'`
	fi

	if [ "$balance" != "$expected_balance" ]; then
		echo "[FAIL] $fail_msg, expecting $expected_balance, but got $balance instead"
		exit 1
	fi
}

assert_option() {
	pkey_file=$1
	option=$2
	expected=$3

    actual=`$bb load-account $(cat $pkey_file) | jq '.'"$option"''`

	if [ "$expected" != "$actual" ]; then
		echo "[FAIL] expected option $option to be $expected, but got $actual instead"
		exit 1
	fi
}

## TEST create key files
echo
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
echo
echo TEST: fund a test account

$bb fund $(cat $pub_file)

assert_balance $pub_file "10000.0000000" "could not fund account"

## TEST create account
echo
echo TEST: create account

$bb gen-keys $tmp/bar
$bb create-account -s '{"source_account":"'$seed'",
                        "new_account":"'$(cat $tmp/bar.pub)'",
						"amount":"1.5"}'

assert_balance $tmp/bar.pub "1.5000000" "could not create account"

## TEST change trust
echo
echo TEST: change trust

$bb gen-keys $tmp/xyz
$bb fund $(cat $tmp/xyz.pub)
$bb change-trust -s '{"source_account": "'$(cat $tmp/xyz)'",
                      "code": "XYZ",
                      "issuer": "'$pub'"}'

assert_balance $tmp/xyz.pub "0.0000000" "could not change trust" "XYZ"

## TEST send custom asset payment
echo
echo TEST: send custom asset payment

$bb send-payment -s '{"from": "'$seed'",
                      "to": "'$(cat $tmp/xyz.pub)'",
                      "token": "XYZ",
                      "amount": "42.0",
                      "issuer": "'$pub'"}'

assert_balance $tmp/xyz.pub "42.0000000" "could not send custom asset payment" "XYZ"

## TEST send native payment
echo
echo TEST: send native payment

$bb send-payment -s '{"from": "'$seed'",
                      "to": "'$(cat $tmp/xyz.pub)'",
                      "amount": "42.0"}'

assert_balance $tmp/xyz.pub "10041.9999900" "could not send native payment"

## TEST manage data
echo
echo TEST: manage data: add value

$bb manage-data -s '{"source_account": "'$seed'",
                     "name": "answer to the ultimate question",
                     "value": "42"}'

assert_option $pub_file "data.\"answer to the ultimate question\"" "\"NDI=\""  ##  "echo NDI= | base64 -D" will result in "42"

echo
echo TEST: manage data: remove value

$bb manage-data -s '{"source_account": "'$seed'",
                     "name": "answer to the ultimate question"}'

assert_option $pub_file "data.\"answer to the ultimate question\"" "null"

## TEST compose transaction
echo
echo TEST: compose transaction

$bb change-trust '{"source_account": "'$(cat $tmp/xyz)'",
                   "code": "ABC",
                   "issuer": "'$pub'"}' | xargs \
  $bb set-options  '{"home_domain": "dotkam.com",
                    "max_weight": 1}' | xargs \
  $bb sign '["'$(cat $tmp/xyz)'"]' | xargs \
  $bb submit

assert_balance $tmp/xyz.pub "0.0000000" "could not compose a transaction" "ABC"
assert_option $tmp/xyz.pub "home_domain" "\"dotkam.com\""

## TEST set options
echo
echo TEST: set options

$bb set-options -s '{"source_account": "'$(cat $tmp/xyz)'",
                     "inflation_destination": "GCCD6AJOYZCUAQLX32ZJF2MKFFAUJ53PVCFQI3RHWKL3V47QYE2BNAUT",
                     "thresholds": {"low": 42, "high": 3},
                     "home_domain": "dotkam.com"}'

assert_option $tmp/xyz.pub "home_domain" "\"dotkam.com\""
assert_option $tmp/xyz.pub "inflation_destination" "\"GCCD6AJOYZCUAQLX32ZJF2MKFFAUJ53PVCFQI3RHWKL3V47QYE2BNAUT\""
assert_option $tmp/xyz.pub "thresholds.low_threshold" 42
assert_option $tmp/xyz.pub "thresholds.high_threshold" 3

## TEST account merge
echo
echo TEST: account merge

$bb account-merge -s '{"source_account": "'$(cat $tmp/bar)'",
                       "destination":"'$pub'"}'

assert_balance $pub_file "9957.9999400" "could not merge two native accounts"


## TEST federation
echo
echo TEST: federation address lookup

account=`$bb federation --address "$BB8_FEDERATION_ADDRESS"`

if [ "$account" != "$BB8_ACCOUNT" ]; then
  echo "[FAIL] could not lookup $BB8_FEDERATION_ADDRESS address on the federation server"
  exit 1
fi

## TEST federation
echo
echo 'TEST: federation account lookup and network switch (via STELLAR_NETWORK)'

## the BB-8 account address only exists on the Stellar public network
export STELLAR_NETWORK=public

address=`$bb federation --account $BB8_ACCOUNT`

if [ "$address" != "$BB8_FEDERATION_ADDRESS" ]; then
  echo "[FAIL] could not lookup $BB8_ACCOUNT account on the federation server"
  unset STELLAR_NETWORK
  exit 1
fi

unset STELLAR_NETWORK

echo
echo "==================="
echo "all tests... [PASS]"

## cleanup

# rm $seed_file $pub_file $tmp/bar*
