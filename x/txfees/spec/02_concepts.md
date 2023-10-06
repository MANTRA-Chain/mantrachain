<!-- order: 2 -->

# Concepts

The tx fees modules is responsible for a couple of things:

- Creating a mapping between a fee token and a liquidity pair.
- Sending all the gas fees to the chain admin wallet address.

## Fee Token

A fee token is a token that is used to pay fees(gas) for transactions instead of the native token.
When create a fee token you specify the denom which you want to use as a fee token, i.e. as alternative to the native token. You also specify a liquidity pair id. The liquidity pair id is used to identify the liquidity pair which will be used to get the last price for that pair. By having the last price the module calculates the amount of the fee token corresponding the amount of the native token to pay the gas price with. The liquidity pair must exist in the liquidity module. Only the `chain admin` can create/update/delete a fee token.

## Gas Fees

In the tx fees module there is an ante hadler which overwrites the default cosmos-sdk ante handler. The tx fees ante handler is responsible for sending all the gas fees to the chain admin wallet address.
