# Horizon Protocol Changelog

Any changes to the Horizon Public API should be included in this doc.

# SDK support

We started tracking SDK support at version 0.12.3. Support for 0.12.3 means that SDK can correctly:

* Send requests using all available query params / POST params / headers,
* Parse all fields in responses structs and headers.

For each new version we will only track changes from the previous version.

| Resource                                      | Changes                                      | Go SDK <sup>1</sup>            | JS SDK             | Java SDK                                          |
|:----------------------------------------------|:---------------------------------------------|:-------------------------------|:-------------------|:--------------------------------------------------|
| **0.12.3**                                    |                                              |                                |                    |                                                   |
| `GET /`                                       |                                              | +<br />(some `_links` missing) | -                  | -                                                 |
| `GET /metrics`                                |                                              | -                              | -                  | -                                                 |
| `GET /ledgers`                                |                                              | -                              | 0.8.2              | 0.2.0                                             |
| `GET /ledgers` SSE                            |                                              | +                              | 0.8.2              | 0.2.0                                             |
| `GET /ledgers/{ledger_id}`                    |                                              | -                              | 0.8.2              | 0.2.0                                             |
| `GET /ledgers/{ledger_id}/transactions`       |                                              | -                              | 0.8.2              | 0.2.0                                             |
| `GET /ledgers/{ledger_id}/transactions` SSE   |                                              | -                              | 0.8.2              | 0.2.0                                             |
| `GET /ledgers/{ledger_id}/operations`         |                                              | -                              | 0.8.2              | 0.2.0                                             |
| `GET /ledgers/{ledger_id}/operations` SSE     |                                              | -                              | 0.8.2              | -                                                 |
| `GET /ledgers/{ledger_id}/payments`           |                                              | -                              | 0.8.2              | 0.2.0                                             |
| `GET /ledgers/{ledger_id}/payments` SSE       |                                              | -                              | 0.8.2              | 0.2.0                                             |
| `GET /ledgers/{ledger_id}/effects`            |                                              | -                              | 0.8.2              | 0.2.0<br />(no support for data, inflation types) |
| `GET /ledgers/{ledger_id}/effects` SSE        |                                              | -                              | 0.8.2              | 0.2.0<br />(no support for data, inflation types) |
| `GET /accounts/{account_id}`                  |                                              | +                              | 0.8.2              | 0.2.0                                             |
| `GET /accounts/{account_id}/transactions`     |                                              | -                              | 0.8.2              | 0.2.0                                             |
| `GET /accounts/{account_id}/transactions` SSE |                                              | -                              | 0.8.2              | 0.2.0                                             |
| `GET /accounts/{account_id}/operations`       |                                              | -                              | 0.8.2              | 0.2.0                                             |
| `GET /accounts/{account_id}/operations` SSE   |                                              | -                              | 0.8.2              | -                                                 |
| `GET /accounts/{account_id}/payments`         |                                              | -                              | 0.8.2              | 0.2.0                                             |
| `GET /accounts/{account_id}/payments` SSE     |                                              | -                              | 0.8.2              | 0.2.0                                             |
| `GET /accounts/{account_id}/effects`          |                                              | -                              | 0.8.2              | 0.2.0<br />(no support for data, inflation types) |
| `GET /accounts/{account_id}/effects` SSE      |                                              | -                              | 0.8.2              | 0.2.0<br />(no support for data, inflation types) |
| `GET /accounts/{account_id}/offers`           |                                              | +                              | -                  | 0.2.0                                             |
| `GET /accounts/{account_id}/trades`           |                                              | -                              | 0.8.2              | -                                                 |
| `GET /accounts/{account_id}/data/{key}`       |                                              | -                              | -                  | -                                                 |
| `POST /transactions`                          |                                              | -                              | 0.8.2              | 0.2.0                                             |
| `GET /transactions`                           |                                              | +                              | 0.8.2              | 0.2.0                                             |
| `GET /transactions` SSE                       |                                              | +                              | 0.8.2              | 0.2.0                                             |
| `GET /transactions/{tx_id}`                   |                                              | +                              | 0.8.2              | 0.2.0                                             |
| `GET /transactions/{tx_id}/operations`        |                                              | -                              | 0.8.2              | 0.2.0                                             |
| `GET /transactions/{tx_id}/operations` SSE    |                                              | -                              | 0.8.2              | -                                                 |
| `GET /transactions/{tx_id}/payments`          |                                              | -                              | 0.8.2              | 0.2.0                                             |
| `GET /transactions/{tx_id}/effects`           |                                              | -                              | 0.8.2              | 0.2.0<br />(no support for data, inflation types) |
| `GET /transactions/{tx_id}/effects` SSE       |                                              | -                              | 0.8.2              | 0.2.0<br />(no support for data, inflation types) |
| `GET /operations`                             |                                              | -                              | 0.8.2              | 0.2.0                                             |
| `GET /operations` SSE                         |                                              | -                              | 0.8.2              | -                                                 |
| `GET /operations/{op_id}`                     |                                              | -                              | 0.8.2              | 0.2.0                                             |
| `GET /operations/{op_id}/effects`             |                                              | -                              | 0.8.2              | 0.2.0<br />(no support for data, inflation types) |
| `GET /operations/{op_id}/effects` SSE         |                                              | -                              | 0.8.2              | 0.2.0<br />(no support for data, inflation types) |
| `GET /payments`                               |                                              | -                              | 0.8.2              | 0.2.0                                             |
| `GET /payments` SSE                           |                                              | +                              | 0.8.2              | 0.2.0                                             |
| `GET /effects`                                |                                              | -                              | 0.8.2              | 0.2.0<br />(no support for data, inflation types) |
| `GET /effects` SSE                            |                                              | -                              | 0.8.2              | 0.2.0<br />(no support for data, inflation types) |
| `GET /trades`                                 |                                              | +                              | 0.8.2              | 0.2.0<br />(no `price` field)                     |
| `GET /trades_aggregations`                    |                                              | +                              | 0.8.2              | 0.2.0                                             |
| `GET /offers/{offer_id}/trades`               |                                              | -                              | 0.8.2              | -                                                 |
| `GET /order_book`                             |                                              | +                              | 0.8.2              | 0.2.0                                             |
| `GET /order_book` SSE                         |                                              | -                              | 0.8.2              | -                                                 |
| `GET /paths`                                  |                                              | -                              | 0.8.2              | 0.2.0                                             |
| `GET /assets`                                 |                                              | -                              | 0.8.2              | 0.2.0                                             |
| [**0.13.0**](#0130) (changes only)            |                                              |                                |                    |                                                   |
| `GET /assets`                                 | `amount` can be larger than `MAX_INT64`/10^7 | +                              | 0.8.2 <sup>2</sup> | 0.2.0                                             |
| `GET /ledgers/{ledger_id}/effects`            | `created_at` field added                     | +                              | 0.8.2 <sup>2</sup> | -                                                 |
| `GET /ledgers/{ledger_id}/effects` SSE        | `created_at` field added                     | +                              | 0.8.2 <sup>2</sup> | -                                                 |
| `GET /accounts/{account_id}/effects`          | `created_at` field added                     | +                              | 0.8.2 <sup>2</sup> | -                                                 |
| `GET /accounts/{account_id}/effects` SSE      | `created_at` field added                     | +                              | 0.8.2 <sup>2</sup> | -                                                 |
| `GET /transactions/{tx_id}/effects`           | `created_at` field added                     | +                              | 0.8.2 <sup>2</sup> | -                                                 |
| `GET /transactions/{tx_id}/effects` SSE       | `created_at` field added                     | +                              | 0.8.2 <sup>2</sup> | -                                                 |
| `GET /operations/{op_id}/effects`             | `created_at` field added                     | +                              | 0.8.2 <sup>2</sup> | -                                                 |
| `GET /operations/{op_id}/effects` SSE         | `created_at` field added                     | +                              | 0.8.2 <sup>2</sup> | -                                                 |
| `GET /effects`                                | `created_at` field added                     | +                              | 0.8.2 <sup>2</sup> | -                                                 |
| `GET /effects` SSE                            | `created_at` field added                     | +                              | 0.8.2 <sup>2</sup> | -                                                 |

<sup>1</sup> We don't do proper versioning for GO SDK yet. `+` means implemented in `master` branch.

<sup>2</sup> Native JSON support in JS, no changes needed.

# Changes

## 0.13.0

- `amount` field in `/assets` is now a String (to support Stellar amounts larger than `int64`).
- Effect resource contains a new `created_at` field.
