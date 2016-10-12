# zcashmoney

This repository contains the following tools for interacting with a Zcash node:

- `privatize` sends unspent transactions to a private address.

- `zreceived` lists all addresses, their balance and received transactions.

- `zsend` is an interactive shell tool for sending zcash.

These tools are written in Go and communicate with a Zcash node using JSON-RPC
with [zcashrpcclient](https://github.com/arithmetric/zcashrpcclient).

When running the tools, you should specify the hostname/port, username, and
password for the RPC service in the following environment variables:

- `ZCASH_RPC_HOST`: The hostname and port for the Zcash node. For example,
`testzcash.example.com:18232`.

- `ZCASH_RPC_USER`: The username for the Zcash node.

- `ZCASH_RPC_PASS`: The password for the Zcash node.

This code is provided for demonstration purposes only under the MIT license.
