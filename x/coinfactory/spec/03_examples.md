<!-- order: 3 -->

# Examples

To create a new token, use the create-denom command from the tokenfactory module. The following example uses the address mantra1axznhnm82lah8qqvp9hxdad49yx3s5dcj66qka from mylocalwallet as the default admin for the new token.

## Creating a token

To create a new token we can use the create-denom command.

```sh
aumegad tx tokenfactory create-denom ufoo --keyring-backend=test --from mylocalwallet
```

## Mint a new token

Once a new token is created, it can be minted using the mint command in the tokenfactory module. Note that the complete tokenfactory address, in the format of factory/{creator address}/{subdenom}, must be used to mint the token.

```sh
aumegad tx tokenfactory mint 100000000000factory/mantra1axznhnm82lah8qqvp9hxdad49yx3s5dcj66qka/ufoo --keyring-backend=test --from mylocalwallet
```

## Checking Token metadata

To view a token's metadata, use the denom-metadata command in the bank module. The following example queries the metadata for the token factory/mantra1axznhnm82lah8qqvp9hxdad49yx3s5dcj66qka/ufoo:

```sh
aumegad query bank denom-metadata --denom factory/mantra1axznhnm82lah8qqvp9hxdad49yx3s5dcj66qka/ufoo
```

## Check the tokens created by an account

To see a list of tokens created by a specific account, use the denoms-from-creator command in the tokenfactory module. The following example shows tokens created by the account mantra1axznhnm82lah8qqvp9hxdad49yx3s5dcj66qka:

```sh
aumegad query tokenfactory denoms-from-creator mantra1axznhnm82lah8qqvp9hxdad49yx3s5dcj66qka
```
