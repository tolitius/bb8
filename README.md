# Stellar Mission Control Center

## buttons to push

```bash
$ ./mc --help
Usage of ./mc:
  -fund string
    	funds a test account. example: --fund address
  -gen-keys string
    	creates a pair of keys (in two files "file-path" and "file-path.pub"). example: --gen-keys file-path
  -submit-tx string
    	submits a base64 encoded transaction. example: --submit-tx txn
```

## create and fund a test account

```bash
$ ./mc --gen-keys foo; ./mc --fund "$(cat foo.pub)"
2018/01/26 00:11:49 keys are created and stored in: foo.pub and foo
2018/01/26 00:11:50 successfully funded GDXSI3GFROEMAK3K77633RZEXTFTJPR2RQVIM4S647MAWS3TW7PQUBSM.
balances: [{Balance:10000.0000000 Limit: Asset:{Type:native Code: Issuer:}}]
more details: https://horizon-testnet.stellar.org/accounts/GDXSI3GFROEMAK3K77633RZEXTFTJPR2RQVIM4S647MAWS3TW7PQUBSM
```

## License

Copyright © 2018 tolitius

Distributed under the Eclipse Public License either version 1.0 or (at
your option) any later version.
