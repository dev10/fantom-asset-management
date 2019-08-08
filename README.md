# What is it
A stable coin issuance platform on top of Cosmos.
A new Cosmos chain that can Issue, Mint, Burn and Transfer custom tokens like 
Binance Asset Management https://docs.binance.org/tokens.html

# Intro
Assets are stored as tokens on Chain, and the below management actions are available.

## Issue
Issue is a transaction used to create a new asset. Anyone can issue a new token with fee paid. 
After issuing, the token would appear in the issuer's account as free balance.

An issuance transaction contains:

* **Source Address**: the sender address of the transaction and it will become the owner of the token, all created tokens will be in this account.
* **Token Name**: it is the long official name, such as "Fantom Coin". It is limited to 32 characters.
* **Symbol**: identifier of the token, limited to 8 alphanumeric characters and is case insensitive, for example, "FTM".

    "F" suffixed symbol is also allowed for migrating tokens that already exist on other chains.

    The symbol doesn't have to be unique, "-" followed by random 3 letters will be appended to the provided symbol to avoid uniqueness constraint.

    Those 3 letters are the first three letters of tx hash of the issue transaction.

    For example, "NNF-F90". Only FTM does not have this suffix.
* **Total Supply**: an int64 boosted by 1e8 for decimal part. The max total supply is 90 billion.
* **Mintable**: that means whether this token can be minted in the future. To set the tokens to be mintable, you need to add --mintable, otherwise just omit this field to set this token to be non-mintable.

### Example on **mainnet:**
```bash
 # To issue a NNB mintable token with total-supply 1 billion on mainnet
 ./famcli token issue --token-name "new token" --total-supply 100000000000000000 --symbol NNF --mintable --from alice --chain-id Fantom-Chain-Alpha  --node  https://data.mainnet.io:443 --trust-node
```
```bash
# To issue a NNB non-mintable token with total-supply 1 billion on mainnet
./famcli token issue --token-name "new token" --total-supply 100000000000000000 --symbol NNF --from alice  --chain-id Fantom-Chain-Alpha   --node  https://data.mainnet.io:443 --trust-node
```

### Example on **testnet:**
```bash
 # To issue a NNB mintable token with total-supply 1 billion on testnet
 ./tfamcli token issue --token-name "new token" --total-supply 100000000000000000 --symbol NNF --mintable --from alice --chain-id Fantom-Chain-Alpha  --node  https://data.testnet.io:80 --trust-node
 
 # Output:  Committed at block 1887 (tx hash: F77A055DDD570AE42A7050182993A0B4DBC81A0D, ... Issued NNF-F77...)
```
```bash
# To issue a NNB non-mintable token with total-supply 1 billion on testnet
./tfamcli token issue --token-name "new token" --total-supply 100000000000000000 --symbol NNF --from alice  --chain-id Fantom-Chain-Omega   --node  https://data.testnet.io:80 --trust-node

# Output: Committed at block 1887 (tx hash: F77A055DDD570AE42A7050182993A0B4DBC81A0D, ... Issued NNF-F77...)
```

## Mint

Tokens that is "mintable" (specified when issue) can use this function. The amount is boosted by **1e8** for decimal part. The total supply after mint is still restricted by 90 billion. 

Note only the `owner` of the token can use this transaction.

Example on **mainnet:**

```bash
./famcli token mint --amount 100000000000000000 --symbol NNF-F77 --from alice --chain-id Fantom-Chain-Alpha --node https://data.defibit.io:443 --trust-node
```

Example on **testnet**:
```bash
./tfamcli token mint --amount 100000000000000000 --symbol NNF-F77 --from alice --chain-id Fantom-Chain-Omega --node https://data.testnet.io:80 --trust-node
```

## Burn
Burn is to destroy certain amount of token, after which that amount of tokens will be subtracted from the operator's balance. The total supply will be updated at the same time. 

Notice that only the owner of the token has the permission to burn token. The amount is boosted by **1e8** for decimal part.
   
Example on **mainnet:**

```bash
./famcli token burn --amount 100000000000000000 --symbol NNF-F77 --from alice --chain-id Fantom-Chain-Alpha --node https://data.mainnet.io:443 --trust-node
```

Example on **testnet:**

```bash
./tfamcli token burn --amount 100000000000000000 --symbol NNF-F77 --from alice --chain-id Fantom-Chain-Omega --node https://data.testnet.io:443 --trust-node
```

## Freeze & Unfreeze
Freeze would move the specified amount of token into "frozen" status, so that these tokens can not transferred, spent in orders or any other transaction until they are unfreezed.

Anyone can (only) freeze or unfreeze tokens on their account with status in "free". The amount is boosted by 1e8 for decimal part.

Example on **mainnet:**
```bash
./famcli token freeze --amount 2000000 --symbol NNF-F77 --from alice --chain-id Fantom-Chain-Alpha --node https://data.mainnet.io:443 --trust-node
```
```bash
./famcli token unfreeze --amount 2000000 --symbol NNF-F77 --from alice --chain-id Fantom-Chain-Alpha --node https://data.mainnet.io:443 --trust-node
```

Example on **testnet:**
```bash
./famcli token freeze --amount 2000000 --symbol NNF-F77 --from alice --chain-id Fantom-Chain-Omega --node https://data.testnet.io:443 --trust-node
```
```bash
./tfamcli token unfreeze --amount 2000000 --symbol NNF-F77 --from alice --chain-id Fantom-Chain-Omega --node https://data.testnet.io:443 --trust-node
```