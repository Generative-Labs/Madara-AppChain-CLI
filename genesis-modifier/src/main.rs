use std::collections::BTreeMap;
use std::fs::File;
use std::io::prelude::*;
use std::ops::Add;
use std::path::PathBuf;

use anyhow::{anyhow, Result};
use clap::{arg, Command,value_parser};
use mp_felt::Felt252Wrapper;
use pallet_starknet::genesis_loader::{GenesisData, HexFelt};
use primitive_types::U256;
use serde::{Deserialize, Serialize};
use starknet_core::utils::get_storage_var_address;
use starknet_ff::FieldElement;

const DEMO_CUSTOM_GENESIS: &str = r#"
{
    "name": "MyToken",
    "symbol": "MT",
    "decimals": 12,
    "total_supply": "1000000",
    "balances": {
        "address1": "1000",
        "address2": "500",
        "address3": "2000"
    }
}
"#;

#[derive(Debug, Serialize, Deserialize)]
struct FeeToken {
    pub name: String,
    pub symbol: String,
    pub decimals: u8,
    pub total_supply: String,
    pub balances: BTreeMap<String, String>,
}

impl FeeToken {
    pub fn generate_contract_storage(&self, contract_address: HexFelt) -> Result<Vec<((HexFelt, HexFelt), HexFelt)>> {
        let mut storages = Vec::<((HexFelt, HexFelt), HexFelt)>::new();

        let name_storage_key: Felt252Wrapper =
            get_storage_var_address("ERC20_name", &[]).map_err(|e| anyhow!(e))?.into();
        let name_data: HexFelt = Felt252Wrapper::try_from(self.name.as_bytes()).map_err(|e| anyhow!(e))?.into();
        storages.push(((contract_address, name_storage_key.into()), name_data));

        let symbol_storage_key: Felt252Wrapper =
            get_storage_var_address("ERC20_symbol", &[]).map_err(|e| anyhow!(e))?.into();
        let symbol_data: HexFelt = Felt252Wrapper::try_from(self.symbol.as_bytes()).map_err(|e| anyhow!(e))?.into();
        storages.push(((contract_address, symbol_storage_key.into()), symbol_data));

        let decimals_storage_key: Felt252Wrapper =
            get_storage_var_address("ERC20_decimals", &[]).map_err(|e| anyhow!(e))?.into();
        let decimals_data: HexFelt = Felt252Wrapper::from(self.decimals).into();
        storages.push(((contract_address, decimals_storage_key.into()), decimals_data));

        for (address, str_balance) in self.balances.iter() {
            let balance = str_balance.parse::<U256>().map_err(|e| anyhow!(e))?;
            let low_balance: u128 = balance.as_u128();
            let high_balance: u128 = (balance >> 128).as_u128();
            let addr = FieldElement::from_hex_be(address).map_err(|e| anyhow!(e))?;
            let low_balance_key = get_storage_var_address("ERC20_balances", &[addr]).map_err(|e| anyhow!(e))?;
            let high_balance_key = low_balance_key.add(FieldElement::from(1u8));

            storages.push((
                (contract_address, Felt252Wrapper::from(low_balance_key).into()),
                Felt252Wrapper::from(low_balance).into(),
            ));
            storages.push((
                (contract_address, Felt252Wrapper::from(high_balance_key).into()),
                Felt252Wrapper::from(high_balance).into(),
            ));
        }

        Ok(storages)
    }
}

fn main() -> Result<()> {
    let matches = Command::new("madara app chain genesis modifier")
        .version("0.1.0")
        .arg(arg!(--custom_genesis_file <FILE>).short('c').help(DEMO_CUSTOM_GENESIS).required(true).value_parser(value_parser!(PathBuf)))
        .arg(arg!(--default_genesis_file <FILE>).short('d').help("madara default genesis file.").required(true).value_parser(value_parser!(PathBuf)))
        .get_matches();

    let custom_genesis_file = matches.get_one::<PathBuf>("custom_genesis_file").expect("can't get custom_genesis_file");
    let default_genesis_path =
        matches.get_one::<PathBuf>("default_genesis_file").expect("can't get default_genesis_file");

    let genesis_file_content = std::fs::read_to_string(default_genesis_path).map_err(|e| anyhow!(e))?;

    let mut default_genesis: GenesisData = serde_json::from_str(&genesis_file_content).map_err(|e| anyhow!(e))?;

    let custom_genesis_context = std::fs::read_to_string(custom_genesis_file).map_err(|e| anyhow!(e))?;

    let custom_genesis: FeeToken = serde_json::from_str(&custom_genesis_context).map_err(|e| anyhow!(e))?;

    // remove fee token storage in default genesis data.
    let default_fee_token_address = default_genesis.fee_token_address;
    default_genesis.storage.retain(|(key, _value)| !key.0.0.eq(&default_fee_token_address.0));

    // generate fee token storage
    let fee_token_storage =
        custom_genesis.generate_contract_storage(default_fee_token_address).map_err(|e| anyhow!(e))?;

    // update genesis
    for data in fee_token_storage.iter() {
        default_genesis.storage.push(*data);
    }

    let modify_genesis = serde_json::to_string(&default_genesis)?;

    File::create(default_genesis_path)
        .and_then(|mut f| f.write_all(modify_genesis.as_bytes()))
        .map_err(|e| anyhow!(e))?;

    Ok(())
}

#[cfg(test)]
mod test {
    use lazy_static::lazy_static;

    use super::*;

    lazy_static! {
        pub static ref EXPECTED_STORAGE: Vec<((Felt252Wrapper, Felt252Wrapper), Felt252Wrapper)> = {
        // The data is from madara unit test code & genesis file
        let expected_storage_data = vec![
            (
                (
                    "0x49d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7",
                    "0x341c1bdfd89f69748aa00b5742b03adbffd79b8e80cab5c50d91cd8c2a79be1",
                ),
                "0x4574686572", // ERC20_name: "Ether"
            ),
            (
                (
                    "0x49d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7",
                    "0x0b6ce5410fca59d078ee9b2a4371a9d684c530d697c64fbef0ae6d5e8f0ac72",
                ),
                "0x455448", // ERC20_symbol: "ETH"
            ),
            (
                (
                    "0x49d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7",
                    "0x1f0d4aa99431d246bac9b8e48c33e888245b15e9678f64f9bdfc8823dc8f979",
                ),
                "0xc", // ERC20_decimals: 12
            ),
            (
                (
                    "0x49d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7",
                    "0x02450cb55e7682ffca9e1db504e2de1263d587242b9b43e27a1eedf18b4bbabf",
                ),
                "0xFFFF", /* ERC20_balances (0x00a13c294af26c4e940d28b1db914e4bb28158638deeeb4ae9ca9b37ab3e4a97) low
                           * key */
            ),
            (
                (
                    "0x49d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7",
                    "0x02450cb55e7682ffca9e1db504e2de1263d587242b9b43e27a1eedf18b4bbac0",
                ),
                "0x0", // ERC20_balances (0x00a13c294af26c4e940d28b1db914e4bb28158638deeeb4ae9ca9b37ab3e4a97) high key
            ),
            (
                (
                    "0x49d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7",
                    "0x053e9b52cb5684cc371a58df3e04f74a96953eeca91400d47abc231983d08303",
                ),
                "0xFFFE", /* ERC20_balances (0x07e0f2302e5861f2617c03af9009a5da775f65ddea7a31d5fb8f85c72d87fe13) low
                           * key */
            ),
            (
                (
                    "0x49d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7",
                    "0x053e9b52cb5684cc371a58df3e04f74a96953eeca91400d47abc231983d08304",
                ),
                "0x0", // ERC20_balances (0x07e0f2302e5861f2617c03af9009a5da775f65ddea7a31d5fb8f85c72d87fe13) high key
            ),
        ];

            let mut expected_storage = Vec::new();
            for data in expected_storage_data.iter() {
                expected_storage.push((
                    (Felt252Wrapper::from_hex_be(data.0 .0).unwrap(), Felt252Wrapper::from_hex_be(data.0 .1).unwrap()),
                    Felt252Wrapper::from_hex_be(data.1).unwrap(),
                ));
            }
            expected_storage
        };
    }

    lazy_static! {
        pub static ref FEE_TOKEN_ADDRESS: HexFelt = {
            let fee_token_address: HexFelt =
                Felt252Wrapper::from_hex_be("0x49d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7")
                    .unwrap()
                    .into();
            fee_token_address
        };
    }

    #[test]
    fn generate_fee_token_storage_should_work() {
        let mut fee_token = FeeToken {
            name: "Ether".to_string(),
            symbol: "ETH".to_string(),
            decimals: 12,
            total_supply: "100000".to_string(),
            balances: BTreeMap::default(),
        };

        fee_token.balances.insert(
            "0x00a13c294af26c4e940d28b1db914e4bb28158638deeeb4ae9ca9b37ab3e4a97".to_string(),
            "0xFFFF".to_string(),
        );
        fee_token.balances.insert(
            "0x07e0f2302e5861f2617c03af9009a5da775f65ddea7a31d5fb8f85c72d87fe13".to_string(),
            "0xFFFE".to_string(),
        );

        let calculated_storages = fee_token.generate_contract_storage(*FEE_TOKEN_ADDRESS).unwrap();

        for (index, storage) in EXPECTED_STORAGE.iter().enumerate() {
            assert_eq!(storage.0.0.0, calculated_storages[index].0.0.0);
            assert_eq!(storage.0.1.0, calculated_storages[index].0.1.0);
            assert_eq!(storage.1.0, calculated_storages[index].1.0);
        }
    }
}
