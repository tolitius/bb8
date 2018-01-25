# Stellar Mission Control Center

## buttons to push

```bash
$ ./mc
  -fund string
    	funds a test account. example: --fund address
  -gen-keys string
    	creates a pair of keys (in two files "file-path" and "file-path.pub"). example: --gen-keys file-path
```

## create and fund a test account

```bash
$ ./mc --gen-keys foo; ./mc --fund $(cat foo.pub)
2018/01/25 00:59:18 keys are created and stored in: foo.pub and foo
2018/01/25 00:59:21 funded GAVRWDUEFADAN7GIL46UF473S2QUDYIJ64HQS55GU4JLQ7TA2VVRLTFD, horizon said:
```
```json
{
  "_links": {
    "transaction": {
      "href": "https://horizon-testnet.stellar.org/transactions/e0a3f4e63c05c4edf92fc08bdc20aa452d169c2bae498f3a7ef6cb1dfdabdbed"
    }
  },
  "hash": "e0a3f4e63c05c4edf92fc08bdc20aa452d169c2bae498f3a7ef6cb1dfdabdbed",
  "ledger": 6900295,
  "envelope_xdr": "AAAAABB90WssODNIgi6BHveqzxTRmIpvAFRyVNM+Hm2GVuCcAAAAZABiwhcAAEEPAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAKxsOhCgGBvzIXz1C8/uWoUHhCfcPCXempxK4fmDVaxUAAAAXSHboAAAAAAAAAAABhlbgnAAAAEAuVaQRM73muCeYUEsgJWCCzRwrVcbAAKi+QeA146F0LXOFJBRGc3SmjObN9Vll03mdWesDvPgcZK6ZHseXCBwJ",
  "result_xdr": "AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAA=",
  "result_meta_xdr": "AAAAAAAAAAEAAAADAAAAAABpSkcAAAAAAAAAACsbDoQoBgb8yF89QvP7lqFB4Qn3Dwl3pqcSuH5g1WsVAAAAF0h26AAAaUpHAAAAAAAAAAAAAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAAAAAAAAAAwBpSkcAAAAAAAAAABB90WssODNIgi6BHveqzxTRmIpvAFRyVNM+Hm2GVuCcAAOdJEEyVGQAYsIXAABBDwAAAAAAAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAAAAAAAAAAQBpSkcAAAAAAAAAABB90WssODNIgi6BHveqzxTRmIpvAFRyVNM+Hm2GVuCcAAOdDPi7bGQAYsIXAABBDwAAAAAAAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAAAAA"
}
```

### checking balance

```bash
$ curl -s https://horizon-testnet.stellar.org/accounts/GAVRWDUEFADAN7GIL46UF473S2QUDYIJ64HQS55GU4JLQ7TA2VVRLTFD | jq '.balances'
```
```json
[
  {
    "balance": "10000.0000000",
    "asset_type": "native"
  }
]
```

## License

Copyright © 2018 tolitius

Distributed under the Eclipse Public License either version 1.0 or (at
your option) any later version.
