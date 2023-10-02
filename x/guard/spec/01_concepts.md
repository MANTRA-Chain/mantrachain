<!-- order: 1 -->

# Concepts

It use of bitwise operations, with each bit representing a **permission/category** for the user. So in this way we are able to create a multi permission requirement performant and cost effectively.
Currently, thee only supported permissions are related to tokens transfer for coins minted by the `x/coinfactory` module.

> The privileges are set to **0**, which means that the user and the custom tokens has no permissions by default and they cannot be traded if the user doesn't pass the KYC and the token groups privileges are not being specified by the operator.
