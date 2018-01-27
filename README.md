# Stellar Mission Control Center

By default stellar-mc is using [Stellar](https://www.stellar.org/) test network a.k.a. testnet.

- [Buttons to Push](#buttons-to-push)
- [Create and Fund a Test Account](#create-and-fund-a-test-account)
- [Issuing a New Token](#issuing-a-new-token)
- [Sending Payments](#sending-payments)
- [License](#license)

## Buttons to Push

``` sh
$ ./mc --help
Usage of ./mc:
  -fund string
    	funds a test account. example: --fund address
  -gen-keys string
    	creates a pair of keys (in two files "file-path" and "file-path.pub"). example: --gen-keys file-path
  -issue-new-token string
    	issue new token (asset). example: --issue-new-token token issuer-seed distributor-seed [limit]
  -send-payment string
    	send payment from one account to another. example: --send-payment '{"from": "seed", "to": "address", "token": "BTC", "amount": "42.0", "issuer": "address"}'
  -submit-tx string
    	submits a base64 encoded transaction. example: --submit-tx txn
```

## Create and Fund a Test Account

``` sh
$ ./mc --gen-keys foo; ./mc --fund "$(cat foo.pub)"

2018/01/26 00:11:49 keys are created and stored in: foo.pub and foo
2018/01/26 00:11:50 successfully funded GDXSI3GFROEMAK3K77633RZEXTFTJPR2RQVIM4S647MAWS3TW7PQUBSM.
balances: [{Balance:10000.0000000 Limit: Asset:{Type:native Code: Issuer:}}]
more details: https://horizon-testnet.stellar.org/accounts/GDXSI3GFROEMAK3K77633RZEXTFTJPR2RQVIM4S647MAWS3TW7PQUBSM
```

## Issuing a New Token

There are usually at least two accounts that participate in issuing a new token (a.k.a. "asset"):

* an issuer account which signs a new asset
* a distribution account that sets a trustline for this "asset" and this "issuer", and is later used as an account that would distribute this asset to other accounts

A distribution account is just a concept, and does not have to exist. Once the issuer signs an asset, any other account on the Stellar network can create a trustline: a declaration that it trusts a particular asset from a particular issuer.

But usually keeping a separate distribution account works well: it is easier to track funds since the money sent back to it won't disappear and would still remain in circulation while any money sent back directly to the issuer account would disappear.

The official name for the distribution account is [specialized issuing account](https://www.stellar.org/developers/guides/issuing-assets.html#specialized-issuing-accounts) as per Stellar documentation.

In this example we would assume no accounts exist so we'll generate issuer and distributor key pairs:

``` sh
$ ./mc --gen-keys issuer
2018/01/26 14:59:21 keys are created and stored in: issuer.pub and issuer
$ ./mc --gen-keys distributor
2018/01/26 14:59:24 keys are created and stored in: distributor.pub and distributor
```

In order to process transactions these accounts need to have a few `XLM`s for minimum balances, trustline and transaction fees, but let's be _very_ generous and fund them `10,000` each:

``` sh
$ ./mc --fund "$(cat issuer.pub)"

2018/01/26 16:27:25 successfully funded account: GBJYH4JSSPHVIJSNU3OFNX2XQUX23N6EV3IPMDLRB2SIWXTUMFEVNY4D.
balances: [{Balance:10000.0000000 Limit: Asset:{Type:native Code: Issuer:}}]
more details: https://horizon-testnet.stellar.org/accounts/GBJYH4JSSPHVIJSNU3OFNX2XQUX23N6EV3IPMDLRB2SIWXTUMFEVNY4D
```

``` sh
$ ./mc --fund "$(cat distributor.pub)"

2018/01/26 16:27:39 successfully funded account: GBUV4AVA53R75U3TYI3KC4GHJ7YPWSKSXZB76ZKTRJHZPKOFM476EY6V.
balances: [{Balance:10000.0000000 Limit: Asset:{Type:native Code: Issuer:}}]
more details: https://horizon-testnet.stellar.org/accounts/GBUV4AVA53R75U3TYI3KC4GHJ7YPWSKSXZB76ZKTRJHZPKOFM476EY6V
```

Now we are ready to issue a new token, let's call it `YUM`:

``` sh
$ ./mc --issue-new-token YUM $(cat issuer) $(cat distributor)

2018/01/26 16:45:55 issued trust for YUM to account: GBUV4AVA53R75U3TYI3KC4GHJ7YPWSKSXZB76ZKTRJHZPKOFM476EY6V.
balances: [{Balance:0.0000000 Limit:922337203685.4775807 Asset:{Type:credit_alphanum4 Code:YUM Issuer:GBJYH4JSSPHVIJSNU3OFNX2XQUX23N6EV3IPMDLRB2SIWXTUMFEVNY4D}} {Balance:9999.9999900 Limit: Asset:{Type:native Code: Issuer:}}]
more details: https://horizon-testnet.stellar.org/accounts/GBUV4AVA53R75U3TYI3KC4GHJ7YPWSKSXZB76ZKTRJHZPKOFM476EY6V
```

`mc --issue-new-token` does several things:

* creates a new asset (in this case YUM)
* signs it with an issuer seed (private key)
* creates a new transaction where it sets a trustline between the receiving account ("distributor") and this asset
* signs this transaction with distributor's seed
* submits this transaction to Stellar

Setting a trustline is called a "[Change Trust](https://www.stellar.org/developers/guides/concepts/list-of-operations.html#change-trust)" operation in Stellar speak. By default this operation would allow the distributor account to receive up to `922337203685.4775807` (`MaxInt64  = 1<<63 - 1`) YUMs. But it has an additional `limit` parameter that sets a cap on how much YUMs an account may get.

For example let's set a cap of `42` YUMs for the distributor account:

``` sh
$ ./mc --issue-new-token YUM $(cat issuer) $(cat distributor) 42

2018/01/26 16:46:00 issued trust for YUM to account: GBUV4AVA53R75U3TYI3KC4GHJ7YPWSKSXZB76ZKTRJHZPKOFM476EY6V.
balances: [{Balance:0.0000000 Limit:42.0000000 Asset:{Type:credit_alphanum4 Code:YUM Issuer:GBJYH4JSSPHVIJSNU3OFNX2XQUX23N6EV3IPMDLRB2SIWXTUMFEVNY4D}} {Balance:9999.9999800 Limit: Asset:{Type:native Code: Issuer:}}]
more details: https://horizon-testnet.stellar.org/accounts/GBUV4AVA53R75U3TYI3KC4GHJ7YPWSKSXZB76ZKTRJHZPKOFM476EY6V
```

notice `Limit:42.0000000` for `Asset:{Code:YUM}`.

All the YUMmy details could be seen on any ledger interface. For example this is the distributor account on [testnet.stellarchain.io](http://testnet.stellarchain.io/address/GBUV4AVA53R75U3TYI3KC4GHJ7YPWSKSXZB76ZKTRJHZPKOFM476EY6V):

<img src="doc/img/yum-42.png">

notice a "Change Trust" operation and zero balance (for now).

## Sending Payments

In order to send a payment of a non native assset, which is any token on a Stellar network besides 'XLM', we need to know several things:

* sender's private key in order to sign this payment transaction
* address of the recepient (i.e. receiving account's public key)
* token code
* amount
* token issuer (i.e. issuer's public key)

To continue the [issuing a new token](#issuing-a-new-token) example, we'll send `42.0` YUMs from the issuer to distributor. Since there are quite a few parameters, the `--send-payment` option takes them as a JSON map with these keys: `"from", "to", "token", "amount", "issuer"`:

```bash
$ ./mc --send-payment '{"from": "'$(cat issuer)'", "to": "'$(cat distributor.pub)'", "token": "YUM", "amount": "42.0", "issuer": "'$(cat issuer.pub)'"}'

2018/01/27 03:01:47 sending 42.0 YUM from GCP4KWV2UOZXEKHV2RV4BO5C5SY6HERCUON7YNCUYQ4D5XMWAKBP5B5L to GCBGULNYR75FIFLS4YKWFW42WRN5NYRHWXCRI672UQCJDCAMXRU6PY6D
2018/01/27 03:01:50 ... [payment sent]
account: GCBGULNYR75FIFLS4YKWFW42WRN5NYRHWXCRI672UQCJDCAMXRU6PY6D.
balances: [{Balance:42.0000000 Limit:42.0000000 Asset:{Type:credit_alphanum4 Code:YUM Issuer:GCP4KWV2UOZXEKHV2RV4BO5C5SY6HERCUON7YNCUYQ4D5XMWAKBP5B5L}} {Balance:9999.9999800 Limit: Asset:{Type:native Code: Issuer:}}]
more details: https://horizon-testnet.stellar.org/accounts/GCBGULNYR75FIFLS4YKWFW42WRN5NYRHWXCRI672UQCJDCAMXRU6PY6D
```

Now if we check the distribution account on [testnet.stellarchain.io](http://testnet.stellarchain.io/address/GCBGULNYR75FIFLS4YKWFW42WRN5NYRHWXCRI672UQCJDCAMXRU6PY6D) we'll see it has `42.0` `YUM`s sent to it from the issuer:

<img src="doc/img/yum-42-42.png">

Source (i.e. `"from"`) of the payment could be any account, not just the issuer, as long as this account has YUMs to send. The reason it was the issuer in the example above is that it was only account that had YUMs. The issuer's public key though should still be used to identify the asset: "this YUM token you are getting was indeed signed by me".

## License

Copyright © 2018 tolitius

Distributed under the Eclipse Public License either version 1.0 or (at
your option) any later version.
