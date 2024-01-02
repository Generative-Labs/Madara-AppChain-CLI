# madara-genesis-modifier
This command-line program modifies the Madara genesis. It takes an ERC20 FeeToken contract's initialization state and converts it into a Cairo storage list. Then, it applies this transformed state to the genesis configuration file of the Madara pallet-starknet.

# Usage

step1: compile binary
```
    cargo build --release    
```

step2: crate custom fee token information json.
```
    {
    "name": "MyToken",
    "symbol": "MT",
    "decimals": 12,
    "total_supply": "0xFFFFF",
    "balances": {
        "0x00a13c294af26c4e940d28b1db914e4bb28158638deeeb4ae9ca9b37ab3e4a97": "0xFFFF",
        "0x07e0f2302e5861f2617c03af9009a5da775f65ddea7a31d5fb8f85c72d87fe13": "0xFFFF"
    }
}
```

step3: run binary
```
    ./target/release/genesis-modifier -c ./fee.json -d ./genesis.json
```

step4: check the modifier genesis.json
You should check the genesis.json again to verify if the modifications have taken effect.