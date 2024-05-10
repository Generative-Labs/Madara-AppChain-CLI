# Madara-AppChain-CLI
Madara-AppChain-CLI

## Installation

Run the following command to install `madcli`:

```bash
curl -H 'Pragma: no-cache' -L https://raw.githubusercontent.com/Generative-Labs/Madara-AppChain-CLI/main/madcliup/install | bash
```

## Install

```bash
madcli install madara
```

## Start

```bash
# setup
madara setup --chain dev --from-remote

# start madara dev
madara --dev --unsafe-rpc-external --rpc-methods Unsafe --rpc-cors all
```
