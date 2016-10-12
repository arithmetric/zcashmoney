# zcashmoney

This repository contains the following tools for interacting with a Zcash node:

- `privatize` sends unspent transactions to a private address.

- `zreceived` lists all addresses, their balance and received transactions.

- `zsend` is an interactive shell tool for sending zcash.

These tools are written in Go and communicate with a Zcash node using JSON-RPC
with [zcashrpcclient](https://github.com/arithmetric/zcashrpcclient).

This code is provided for demonstration purposes only under the MIT license.
