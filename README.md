# Stellar Mission Control Center

A command line interface to [Stellar](https://www.stellar.org/) network.

- [Why](#why)
- [Buttons to Push](#buttons-to-push)
- [Create Account Keys](#create-account-keys)
- [Funding a Test Account](#funding-a-test-account)
- [Account Details](#account-details)
- [Issuing a New Token](#issuing-a-new-token)
  - [Issuer and Distributor](#issuer-and-distributor)
  - [Creating and funding accounts](#creating-and-funding-accounts)
  - [Do You Trust Me?](#do-you-trust-me)
  - [Limitting Trustline](#limitting-trustline)
- [Sending Payments](#sending-payments)
- [Transaction Options](#transaction-options)
  - [Adding Discoverablity and Meta Information](#adding-discoverablity-and-meta-information)
  - [Setting Inflation Destination](#setting-inflation-destination)
- [License](#license)

## Why

There are already mutiple ways to interact with Stellar:

* [REST API](https://www.stellar.org/developers/horizon/reference/index.html)
* [SDK libraries](https://www.stellar.org/developers/reference/#libraries) in several languages
* [Stellar Labratory](https://www.stellar.org/laboratory/)

Stellar Mission Control Center adds a command line / terminal capabilities to the Stellar family of tools. This is useful for exploration as well as the real world interaction with Stellar network.

## Buttons to Push

By default stellar-mc is using [Stellar](https://www.stellar.org/) test network a.k.a. testnet.
> _//TODO: add a note about switching between networks_

``` sh
$ ./mc --help
Usage of ./mc:
  -account-details string
    	load and return account details. example: --account-details address
  -fund string
    	funds a test account. example: --fund address
  -gen-keys string
    	creates a pair of keys (in two files "file-path" and "file-path.pub"). example: --gen-keys file-path
  -issue-new-token string
    	issue new token/asset. example ("limit" param is optional): --issue-new-token '{"code": "XYZ", "issuer-address": "address", "distributor-seed":"seed", "limit": "42.0"}'
  -new-tx string
    	build and submit a new transaction. "operations" and "signers" are optional, if there are no "signers", the "source-account" seed will be used to sign this transaction. example: --new-tx '{"source-account": "address or seed", {"operations": "trust": {"code": "XYZ", "issuer-address": "address"}}, "signers": ["seed1", "seed2"]}'
  -send-payment string
    	send payment from one account to another. example: --send-payment '{"from": "seed", "to": "address", "token": "BTC", "amount": "42.0", "issuer": "address"}'
  -submit-tx string
    	submits a base64 encoded transaction. example: --submit-tx txn
  -tx-options string
    	add one or more transaction options. example: --tx-options '{"home-domain": "stellar.org", "max-weight": 1, "inflation-destination": "address"}'
```

## Create Account Keys

Every Stellar account has a pair of keys:

* a public key that is also known as account's `address`
* a private key that is also known as `seed`

`stellar-mc` has a `--gen-keys` option to generate this pair of keys:

``` sh
$ ./mc --gen-keys foo
2018/01/30 15:15:46 keys are created and stored in: foo.pub and foo
```

`foo` in this case is a path to a pair of files where these keys will be stored:

```sh
$ echo address: $(cat foo.pub); echo seed: $(cat foo)
address: GBTG5ZSVAG6LNKA5ZGMD4SJYJX4AQI347WYURPHZV7A6DC7FCG44DOMY
seed: SAWDCNZF7Y67RWT5FC6YSAHFRH23OFU2OFYFXFNWYETL7S7J72CWK6JG
```

An account `seed` should be kept private, hence the "`private`" key.
It is later used to _sign_ Stellar transactions to confirm that it is "_really you_" and "_you approve_" the transaction.

## Funding a Test Account

Stellar has a friendly utility called [Friendbot](https://www.stellar.org/developers/horizon/reference/tutorials/follow-received-payments.html#funding-your-account) that funds a new account on the Stellar test network. When a new account is created (e.g. a pair of keys we created above), this account has no balance and does not exist in the ledger until it is funded. Friendbot fixes that problem.

`stellar-mc` has `--fund` option that takes an account's address and funds it a good amount of lumens:

``` sh
$ ./mc --fund $(cat foo.pub)
```

here we used a `foo.pub` address that we generated above. Next, we'll look at the account on the real, distributed Stellar ledger.

## Account Details

In order to look at the account in the ledger `stellar-mc` provides an `--account-details` option that takes an account address and returns all the details known to Stellar:

```sh
$ ./mc --account-details $(cat foo.pub)
```
```json
{
  "_links": {
    "self": {
      "href": "https://horizon-testnet.stellar.org/accounts/GBTG5ZSVAG6LNKA5ZGMD4SJYJX4AQI347WYURPHZV7A6DC7FCG44DOMY"
    },
    "transactions": {
      "href": "https://horizon-testnet.stellar.org/accounts/GBTG5ZSVAG6LNKA5ZGMD4SJYJX4AQI347WYURPHZV7A6DC7FCG44DOMY/transactions{?cursor,limit,order}",
      "templated": true
    },
    "operations": {
      "href": "https://horizon-testnet.stellar.org/accounts/GBTG5ZSVAG6LNKA5ZGMD4SJYJX4AQI347WYURPHZV7A6DC7FCG44DOMY/operations{?cursor,limit,order}",
      "templated": true
    },
    "payments": {
      "href": "https://horizon-testnet.stellar.org/accounts/GBTG5ZSVAG6LNKA5ZGMD4SJYJX4AQI347WYURPHZV7A6DC7FCG44DOMY/payments{?cursor,limit,order}",
      "templated": true
    },
    "effects": {
      "href": "https://horizon-testnet.stellar.org/accounts/GBTG5ZSVAG6LNKA5ZGMD4SJYJX4AQI347WYURPHZV7A6DC7FCG44DOMY/effects{?cursor,limit,order}",
      "templated": true
    },
    "offers": {
      "href": "https://horizon-testnet.stellar.org/accounts/GBTG5ZSVAG6LNKA5ZGMD4SJYJX4AQI347WYURPHZV7A6DC7FCG44DOMY/offers{?cursor,limit,order}",
      "templated": true
    }
  },
  "id": "GBTG5ZSVAG6LNKA5ZGMD4SJYJX4AQI347WYURPHZV7A6DC7FCG44DOMY",
  "paging_token": "",
  "account_id": "GBTG5ZSVAG6LNKA5ZGMD4SJYJX4AQI347WYURPHZV7A6DC7FCG44DOMY",
  "sequence": "30170921188720640",
  "subentry_count": 0,
  "thresholds": {
    "low_threshold": 0,
    "med_threshold": 0,
    "high_threshold": 0
  },
  "flags": {
    "auth_required": false,
    "auth_revocable": false
  },
  "balances": [
    {
      "balance": "10000.0000000",
      "asset_type": "native"
    }
  ],
  "signers": [
    {
      "public_key": "GBTG5ZSVAG6LNKA5ZGMD4SJYJX4AQI347WYURPHZV7A6DC7FCG44DOMY",
      "weight": 1,
      "key": "GBTG5ZSVAG6LNKA5ZGMD4SJYJX4AQI347WYURPHZV7A6DC7FCG44DOMY",
      "type": "ed25519_public_key"
    }
  ],
  "data": {}
}
```

notice the `balances` section:

```sh
$ ./mc --account-details $(cat foo.pub) | jq '.balances'
```

```json
{
  "balance": "10000.0000000",
  "asset_type": "native"
}
```

it is telling us that the Friendbot from the step above did what we asked and funded this account with `10,000` lumens (a.k.a. as "native" Stellar currency).

## Issuing a New Token

> _One of Stellar’s most powerful features is the ability to trade any kind of asset, US dollars, Nigerian naira, bitcoins, special coupons, ICO tokens or just about anything you like. ([Stellar Developers Guide](https://www.stellar.org/developers/guides/issuing-assets.html))_

### Issuer and Distributor

There are usually at least two accounts that participate in issuing a new token (a.k.a. "asset"):

* an issuer account which "signs a new asset into existence"
* a distribution account that sets a [trustline](https://www.stellar.org/developers/guides/concepts/assets.html#trustlines) for both: this "asset" and this "issuer", and is later used as an account that would distribute this asset to other accounts

A distribution account is just a concept, and does not have to exist. Once the issuer signs an asset, any other account on the Stellar network can create a trustline: a declaration that it trusts a particular asset from a particular issuer.

But usually keeping a separate distribution account works well: it is easier to track funds since the money sent back to it won't disappear and would still remain in circulation while any money sent back directly to the issuer account would disappear.

The official name for the distribution account is [specialized issuing account](https://www.stellar.org/developers/guides/issuing-assets.html#specialized-issuing-accounts) as per Stellar documentation.

### Creating and funding accounts

In this example we would assume no accounts exist so we'll generate issuer and distributor key pairs:

``` sh
$ ./mc --gen-keys issuer
2018/01/30 15:42:48 keys are created and stored in: issuer.pub and issuer

$ ./mc --gen-keys distributor
2018/01/30 15:42:52 keys are created and stored in: distributor.pub and distributor
```

In order to process transactions these accounts need to have a few `XLM`s for minimum balances, trustline and transaction fees. We'll use Stellar's Friendbot to fund these accounts:

``` sh
$ ./mc --fund "$(cat issuer.pub)"
$ ./mc --fund "$(cat distributor.pub)"
```

Let's make sure it worked by checking their balances:

```sh
$ ./mc --account-details $(cat issuer.pub) | jq '.balances'
```
```json
{
  "balance": "10000.0000000",
  "asset_type": "native"
}
```

```sh
$ ./mc --account-details $(cat distributor.pub) | jq '.balances'
```
```json
{
  "balance": "10000.0000000",
  "asset_type": "native"
}
```

### Do You Trust Me?

Now we are ready to issue a new token, let's call it `YUM`.

`stellar-mc` has an `--issue-new` option that takes a token's "code", issuer's address and distributor's seed (to sign the transaction) as JSON:

``` sh
$ ./mc --issue-new-token '{"code": "YUM",
                           "issuer-address": "'$(cat issuer.pub)'",
                           "distributor-seed":"'$(cat distributor)'"}'
```

`mc --issue-new-token` does several things:

* creates a new asset (in this case `YUM`)
* signs it with an issuer's seed (private key)
* creates a new transaction where it sets a trustline between the receiving account ("distributor") and this asset
* signs this transaction with distributor's seed
* submits this transaction to Stellar

Let's check that `YUM` is now an existing token that was issued by issuer's address and that the distributor has successfully _created a trustline_ for it:

```sh
$ ./mc --account-details $(cat distributor.pub) | jq '.balances'
```

```json
{
  "balance": "0.0000000",
  "limit": "922337203685.4775807",
  "asset_type": "credit_alphanum4",
  "asset_code": "YUM",
  "asset_issuer": "GBW2U2GEWVD7GDTQPPJSDGE4SRYXN3USYZKNNI6EPVHUHROS47S6NUZJ"
},
{
  "balance": "9999.9999800",
  "asset_type": "native"
}
```

Nice! We have established a trustline for YUMs and almost ready to distribute them to other accounts. Notice the native balance is no longer 10,000 lumens. This is due to the fees the distribution account had to pay: 100 [stroops](https://www.stellar.org/developers/guides/concepts/assets.html#one-stroop-multiple-stroops) for the transaction processing and another 100 stroops for setting up a trustline which is one of the transaction operations.

Setting a trustline is called a "[Change Trust](https://www.stellar.org/developers/guides/concepts/list-of-operations.html#change-trust)" operation in Stellar speak. By default this operation would allow the distributor account to receive up to `922337203685.4775807` (`MaxInt64  = 1<<63 - 1`) YUMs. But it has an additional `limit` parameter that sets a cap on how much YUMs an account may get.

### Limitting Trustline

The `--issue-new` options takes an optional `limit` parameter to set such a cap. For example let's set a cap of `42` YUMs for the distribution account:

``` sh
$ ./mc --issue-new-token '{"code": "YUM",
                           "issuer-address": "'$(cat issuer.pub)'",
                           "distributor-seed":"'$(cat distributor)'",
                           "limit": "42.0"}'
```

```sh
./mc --account-details $(cat distributor.pub) | jq '.balances'
```
```json
{
  "balance": "0.0000000",
  "limit": "42.0000000",
  "asset_type": "credit_alphanum4",
  "asset_code": "YUM",
  "asset_issuer": "GBW2U2GEWVD7GDTQPPJSDGE4SRYXN3USYZKNNI6EPVHUHROS47S6NUZJ"
},
{
  "balance": "9999.9999600",
  "asset_type": "native"
}
```

notice `"limit": "42.0000000"` for YUMs.

All the YUMmy details could be seen on any ledger interface. For example this is the distribution account on [testnet.stellarchain.io](http://testnet.stellarchain.io/address/GDPKQGOY33DYUPC3PXX222FRZOLQD4L6CMXGJV5I4W2GB4UOT4MCJCO5):

<img src="doc/img/yum-42.png">

notice a "Change Trust" operation and zero balance (for now).

## Sending Payments

In order to send a payment of a non native assset, which is any token on a Stellar network besides `XLM`, we need to know several things:

* sender's private key in order to sign this payment transaction
* address of the recepient (i.e. receiving account's public key)
* token code
* amount
* token issuer (i.e. issuer's public key)

To continue the [issuing a new token](#issuing-a-new-token) example, we'll send `42.0` YUMs from the issuer to distributor.

`stellar-mc` has a `--send-payment` option that takes a JSON map with these keys: `"from", "to", "token", "amount", "issuer"`:

```sh
$ ./mc --send-payment '{"from": "'$(cat issuer)'",
                        "to": "'$(cat distributor.pub)'",
                        "token": "YUM",
                        "amount": "42.0",
                        "issuer": "'$(cat issuer.pub)'"}'

2018/01/30 16:11:56 sending 42.0 YUM from GBW2U2GEWVD7GDTQPPJSDGE4SRYXN3USYZKNNI6EPVHUHROS47S6NUZJ to GDPKQGOY33DYUPC3PXX222FRZOLQD4L6CMXGJV5I4W2GB4UOT4MCJCO5
```

Let's check the balance now:

```sh
$ ./mc --account-details $(cat distributor.pub) | jq '.balances'
```
```json
{
  "balance": "42.0000000",
  "limit": "42.0000000",
  "asset_type": "credit_alphanum4",
  "asset_code": "YUM",
  "asset_issuer": "GBW2U2GEWVD7GDTQPPJSDGE4SRYXN3USYZKNNI6EPVHUHROS47S6NUZJ"
},
{
  "balance": "9999.9999600",
  "asset_type": "native"
}
```

Here is the prettier version of `42.0` YUMs on [testnet.stellarchain.io](http://testnet.stellarchain.io/address/GDPKQGOY33DYUPC3PXX222FRZOLQD4L6CMXGJV5I4W2GB4UOT4MCJCO5):

<img src="doc/img/yum-42-42.png">

Source (i.e. `"from"`) of the payment could be any account, not just the issuer, as long as this account has YUMs to send. The reason it was the issuer in the example above is that it was only account that had YUMs. The issuer's public key though should still be used to identify the asset: "this YUM token you are getting was indeed signed by me".

Excellent, we are now ready to distribute YUMs. We can use the same `--send-payment` option with different account addresses to do that.

## Transaction Options

When submitting a transaction to Stellar there are several [transaction options](https://www.stellar.org/developers/guides/concepts/list-of-operations.html#set-options) that could be set.

`stellar-mc` has a `--tx-options` option that takes this options as JSON and sets them on a transaction before it is submitted.

### Adding Discoverablity and Meta Information

To continue the [issuing a new token](#issuing-a-new-token) example, whenever a new token/asset is introduced to Stellar network it is important to provide clear information about what this token/asset represents. This info can be discovered and displayed by clients so users know exactly what they are getting when they hold your asset. Here is [more about it](https://www.stellar.org/developers/guides/issuing-assets.html#discoverablity-and-meta-information) from Stellar documentation.

In order to discover information about a particular token Stellar would look at a "home domain" property of an account and then will try to read a "[stellar.toml](https://www.stellar.org/developers/guides/concepts/stellar-toml.html)" file at "`https://home-domain/.well-known/stellar.toml`".

Since we issued a brand new `YUM` token, we can create a "`stellar.toml`" file to describe it make it reachable at "`https://home-domain/.well-known/stellar.toml`", and let Stellar know to look for it there by setting a "home domain" transaction option on the issuer's account by `--tx-options`:

```sh
$ ./mc --new-tx '{"source-account": "'$(cat issuer)'"}' \
       --tx-options '{"home-domain": "dotkam.com"}'
```

> _will discuss `--new-tx` later as it is still work in progress to include other transaction operations_

and now this link to the domain is there on the Stellar network:

```sh
$ ./mc --account-details $(cat issuer.pub) | jq '.home_domain'
"dotkam.com"
```
### Setting Inflation Destination

Another example of using Stellar transaction options would be setting an inflation destination on the account.

> _The Stellar distributed network has a built-in, fixed, nominal inflation mechanism. New lumens are added to the network at the rate of 1% each year. Each week, the protocol distributes these lumens to any account that gets over .05% of the “votes” from other accounts in the network_ (from [Stellar documentation](https://www.stellar.org/developers/guides/concepts/inflation.html))

Inflation destination can be set via `--tx-options`. For example let's set an inflation destination on the distributor account from the examples above:

```sh
$ ./mc --new-tx '{"source-account": "'$(cat distributor)'"}' \
       --tx-options '{"inflation-destination": "GCCD6AJOYZCUAQLX32ZJF2MKFFAUJ53PVCFQI3RHWKL3V47QYE2BNAUT"}'
```

We can combine other options, let's add a home domain as well:
```sh
./mc --new-tx '{"source-account": "'$(cat distributor)'"}' \
     --tx-options '{"home-domain": "dotkam.com",
                    "inflation-destination": "GCCD6AJOYZCUAQLX32ZJF2MKFFAUJ53PVCFQI3RHWKL3V47QYE2BNAUT"}'
```

and we can check that both options were set successfully:

```sh
$ ./mc --account-details $(cat distributor.pub) | jq '.home_domain, .inflation_destination'
"dotkam.com"
"GCCD6AJOYZCUAQLX32ZJF2MKFFAUJ53PVCFQI3RHWKL3V47QYE2BNAUT"
```

Here is a prettier version of the options that were set in this transaction on [testnet.stellarchain.io](http://testnet.stellarchain.io/address/GDPKQGOY33DYUPC3PXX222FRZOLQD4L6CMXGJV5I4W2GB4UOT4MCJCO5):

<img src="doc/img/tx-options.png">

## License

Copyright © 2018 tolitius

Distributed under the Eclipse Public License either version 1.0 or (at
your option) any later version.
