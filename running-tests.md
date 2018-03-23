when contributing to BB-8, please (add your own if needed and) run tests to make sure all is well:

```bash
[bb8]$ ./test.sh

TEST: create key files
2018/03/23 15:31:31 keys are created and stored in: /tmp/foo.pub and /tmp/foo

TEST: fund a test account

TEST: create account
2018/03/23 15:31:37 keys are created and stored in: /tmp/bar.pub and /tmp/bar
2018/03/23 15:31:37 submitting transaction to horizon at https://horizon-testnet.stellar.org

TEST: change trust
2018/03/23 15:31:42 keys are created and stored in: /tmp/xyz.pub and /tmp/xyz
2018/03/23 15:31:46 submitting transaction to horizon at https://horizon-testnet.stellar.org

TEST: send custom asset payment
2018/03/23 15:31:52 sending 42.0 XYZ from GAFOL4AJDIGAPWY5DMHS3FWHOWT2F7DBH4VU4XTHPVGOLRE5MMSA2SDN to GA53VEL47MTSR3ZVDQ7EQW6QIGETFCWBQB5IQKHHHLIREXSJFGWA33U4
2018/03/23 15:31:53 submitting transaction to horizon at https://horizon-testnet.stellar.org

TEST: send native payment
2018/03/23 15:31:57 sending 42.0 XLM from GAFOL4AJDIGAPWY5DMHS3FWHOWT2F7DBH4VU4XTHPVGOLRE5MMSA2SDN to GA53VEL47MTSR3ZVDQ7EQW6QIGETFCWBQB5IQKHHHLIREXSJFGWA33U4
2018/03/23 15:31:58 submitting transaction to horizon at https://horizon-testnet.stellar.org

TEST: manage data
2018/03/23 15:32:03 submitting transaction to horizon at https://horizon-testnet.stellar.org

TEST: compose transaction
2018/03/23 15:32:07 submitting transaction to horizon at https://horizon-testnet.stellar.org

TEST: set options
2018/03/23 15:32:12 submitting transaction to horizon at https://horizon-testnet.stellar.org

TEST: account merge
2018/03/23 15:32:18 submitting transaction to horizon at https://horizon-testnet.stellar.org

TEST: federation address lookup

TEST: federation account lookup and network switch (via STELLAR_NETWORK)

===================
all tests... [PASS]
```
