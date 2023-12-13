<!-- order: 3 -->

# State Transitions

This document describes the state transitions pertaining a [DidDocument](02_state.md#diddocument) according to the [did operations](https://www.w3.org/TR/did-core/#method-operations):

1. [Create](03_state_transitions.md#Create)
2. [Resolve](03_state_transitions.md#Resolve)
3. [Update](03_state_transitions.md#Update)
4. [Deactivate](03_state_transitions.md#Deactivate)

A [DidMetadata](02_state.md#didmetadata) lifecycle follows the lifecycle of a  [DidDocument](02_state.md#diddocument)

## Create

[DidDocument](02_state.md#diddocument) are created via the rpc method `CreateDidDocument` that accepts a [MsgCreateDidDocument](./04_messages.md#MsgCreateDidDocument) messages as parameter.

The operation will fail if:

- the signer account has insufficient funds
- the did is malformed
- a did document with the same did exists
- verifications
  - the verification method is invalid (according to the verification method specifications)
  - there is more than one verification method with the same id
  - relationships are empty
  - relationships contain unsupported values (according to the did method specifications)
- services are invalid (according to the services specifications)

Example:

<!-- 

aumegad tx did create-did \
 900d82bc-2bfe-45a7-ab22-a8d11773568e \
 --from vasp --node https://127.0.0.1:443 --chain-id aumega
-->

```javascript
/* gRPC message */
CreateDidDocument(
    MsgCreateDidDocument(
        "did:cosmos:aumega:900d82bc-2bfe-45a7-ab22-a8d11773568e",
        [], // controller
        [   // verifications
            {
                "relationships": ["authentication"],
                {
                    "controller": "did:cosmos:aumega:900d82bc-2bfe-45a7-ab22-a8d11773568e",
                    "id": "did:cosmos:aumega:900d82bc-2bfe-45a7-ab22-a8d11773568e#mantra1axznhnm82lah8qqvp9hxdad49yx3s5dcj66qka",
                    "publicKeyHex": "0248a5178d7a90ec187b3c3d533a4385db905f6fcdaac5026859ca5ef7b0b1c3b5",
                    "type": "EcdsaSecp256k1VerificationKey2019"
                },
                [],
            },
        ],
        [], // services
        "mantra1axznhnm82lah8qqvp9hxdad49yx3s5dcj66qka" // signer
    )
)

/* Resolved DID document */
{
  "didDocument": {
    "context": [
      "https://www.w3.org/ns/did/v1"
    ],
    "id": "did:cosmos:aumega:900d82bc-2bfe-45a7-ab22-a8d11773568e",
    "controller": [],
    "verificationMethod": [
      {
        "controller": "did:cosmos:aumega:900d82bc-2bfe-45a7-ab22-a8d11773568e",
        "id": "did:cosmos:aumega:900d82bc-2bfe-45a7-ab22-a8d11773568e#mantra1axznhnm82lah8qqvp9hxdad49yx3s5dcj66qka",
        "publicKeyHex": "0248a5178d7a90ec187b3c3d533a4385db905f6fcdaac5026859ca5ef7b0b1c3b5",
        "type": "EcdsaSecp256k1VerificationKey2019"
      }
    ],
    "service": [],
    "authentication": [
      "did:cosmos:aumega:900d82bc-2bfe-45a7-ab22-a8d11773568e#mantra1axznhnm82lah8qqvp9hxdad49yx3s5dcj66qka"
    ],
    "assertionMethod": [],
    "keyAgreement": [],
    "capabilityInvocation": [],
    "capabilityDelegation": []
  },
  "didMetadata": {
    "versionId": "571615b8146082deaac90fa01afc8ff88e5a71b4c9c29bcaffef2d11b39a0437",
    "created": "2021-08-23T08:24:26.972761898Z",
    "updated": "2021-08-23T08:24:26.972761898Z",
    "deactivated": false
  }
}

```

### Resolve

[DidDocument](02_state.md#diddocument) are resolved via the rpc method `QueryDidDocument` that accepts a [QueryDidDocumentRequest](./04_messages.md#QueryDidDocumentRequest) messages as parameter.

The operation will fail if:

- the did does not exists

Example:

<!--
aumegad query did did did:aumega:900d82bc-2bfe-45a7-ab22-a8d11773568e \
 --from vasp --node https://127.0.0.1:443 --chain-id aumega \
 --output=json | jq
-->

```javascript
/* gRPC message */
QueryDidDocument(
    QueryDidDocumentRequest(
        "did:cosmos:aumega:900d82bc-2bfe-45a7-ab22-a8d11773568e"
    )
)

/* Resolved DID Document */
{
  "didDocument": {
    "context": [
      "https://www.w3.org/ns/did/v1"
    ],
    "id": "did:cosmos:aumega:900d82bc-2bfe-45a7-ab22-a8d11773568e",
    "controller": [],
    "verificationMethod": [
      {
        "controller": "did:cosmos:aumega:900d82bc-2bfe-45a7-ab22-a8d11773568e",
        "id": "did:cosmos:aumega:900d82bc-2bfe-45a7-ab22-a8d11773568e#mantra1axznhnm82lah8qqvp9hxdad49yx3s5dcj66qka",
        "publicKeyHex": "0248a5178d7a90ec187b3c3d533a4385db905f6fcdaac5026859ca5ef7b0b1c3b5",
        "type": "EcdsaSecp256k1VerificationKey2019"
      }
    ],
    "service": [],
    "authentication": [
      "did:cosmos:aumega:900d82bc-2bfe-45a7-ab22-a8d11773568e#mantra1axznhnm82lah8qqvp9hxdad49yx3s5dcj66qka"
    ],
    "assertionMethod": [],
    "keyAgreement": [],
    "capabilityInvocation": [],
    "capabilityDelegation": []
  },
  "didMetadata": {
    "versionId": "571615b8146082deaac90fa01afc8ff88e5a71b4c9c29bcaffef2d11b39a0437",
    "created": "2021-08-23T08:24:26.972761898Z",
    "updated": "2021-08-23T08:24:26.972761898Z",
    "deactivated": false
  }
}

```

### Update

[DidDocument](02_state.md#diddocument) are updated via the rpc methods:

- UpdateDidDocument
- AddVerification
- RevokeVerification
- SetVerificationRelationships
- AddService
- DeleteService

All the operations will fail if:

- the signer account has insufficient funds
- the signer account address doesn't match the verification method listed in the `Authorization` verification relationships
- the target did does not exists

The following sections provide specific details for each method invocation.

#### UpdateDidDocument

The  `UpdateDidDocument` method will is used to **overwrite** the  [DidDocument](02_state.md#diddocument) controllers. It accepts a [MsgUpdateDidDocument](./04_messages.md#MsgUpdateDidDocument) as a parameter.

The operation will fail if:

- any of the did provided controllers is not a valid did

<!-- 

aumegad tx did update-did-document \
 900d82bc-2bfe-45a7-ab22-a8d11773568e \
 mantra19xvxf55pzvs8ayeuce843sp3rjgmssh77e83fq \
 --from vasp --node https://127.0.0.1:443 --chain-id aumega
-->

```javascript
/* gRPC message */
UpdateDidDocument(
    MsgUpdateDidDocument(
        "did:cosmos:aumega:900d82bc-2bfe-45a7-ab22-a8d11773568e",
        ["did:cosmos:key:mantra19xvxf55pzvs8ayeuce843sp3rjgmssh77e83fq"],
        "mantra1axznhnm82lah8qqvp9hxdad49yx3s5dcj66qka"
    )
)

/* Resolved DID Document */

{
  "didDocument": {
    "context": [
      "https://www.w3.org/ns/did/v1"
    ],
    "id": "did:cosmos:aumega:900d82bc-2bfe-45a7-ab22-a8d11773568e",
    "controller": [
      "did:cosmos:key:mantra19xvxf55pzvs8ayeuce843sp3rjgmssh77e83fq"
    ],
    "verificationMethod": [
      {
        "controller": "did:cosmos:aumega:900d82bc-2bfe-45a7-ab22-a8d11773568e",
        "id": "did:cosmos:aumega:900d82bc-2bfe-45a7-ab22-a8d11773568e#mantra1axznhnm82lah8qqvp9hxdad49yx3s5dcj66qka",
        "publicKeyHex": "0248a5178d7a90ec187b3c3d533a4385db905f6fcdaac5026859ca5ef7b0b1c3b5",
        "type": "EcdsaSecp256k1VerificationKey2019"
      }
    ],
    "service": [],
    "authentication": [
      "did:cosmos:aumega:900d82bc-2bfe-45a7-ab22-a8d11773568e#mantra1axznhnm82lah8qqvp9hxdad49yx3s5dcj66qka"
    ],
    "assertionMethod": [],
    "keyAgreement": [],
    "capabilityInvocation": [],
    "capabilityDelegation": []
  },
  "didMetadata": {
    "versionId": "9f7c547dc852af60c9da1fd514e1497d407b6a3d8ae3e52b626d536519dc8f4c",
    "created": "2021-08-23T08:24:26.972761898Z",
    "updated": "2021-08-24T13:27:50.024635302Z",
    "deactivated": false
  }
}
```

#### AddVerification

The `AddVerification` method is used to add new [verification methods](https://w3c.github.io/did-core/#verification-methods) and [verification relationships](https://w3c.github.io/did-core/#verification-relationships) to a [DidDocument](02_state.md#diddocument). It accepts a [MsgAddVerification](./04_messages.md#MsgAddVerification) as a parameter.

The operation will fail if:

- the verification method is invalid (according to the verification method specifications)
- the verification method id already exists for the did document
- the verification relationships are empty
- the verification relationships contain unsupported values (according to the did method specification)

<!-- 

aumegad tx did add-verification-method \
 900d82bc-2bfe-45a7-ab22-a8d11773568e \
 cosmospub1addwnpepqduxp90pt6ez3a8p26fwmfhqvparz0xqsxk4f4564hg46527xpzeq82cerm \
 --from vasp --node https://127.0.0.1:443 --chain-id aumega
-->

```javascript
/* gRPC message */
AddVerification(
    MsgAddVerification(
        "did:cosmos:aumega:900d82bc-2bfe-45a7-ab22-a8d11773568e",
        {
            "relationships": ["authentication"],
            {
                "controller": "did:cosmos:aumega:900d82bc-2bfe-45a7-ab22-a8d11773568e",
                "id": "did:cosmos:aumega:900d82bc-2bfe-45a7-ab22-a8d11773568e#cosmos1lvl2s8x4pta5f96appxrwn3mypsvumukvk7ck2",
                "publicKeyHex": "03786095e15eb228f4e15692eda6e0607a313cc081ad54d69aadd15d515e304590",
                "type": "EcdsaSecp256k1VerificationKey2019"
            },
            [],
        },
        "mantra1axznhnm82lah8qqvp9hxdad49yx3s5dcj66qka" // signer
    )
)

/* Resolved DID Document */

{
  "didDocument": {
    "context": [
      "https://www.w3.org/ns/did/v1"
    ],
    "id": "did:cosmos:aumega:900d82bc-2bfe-45a7-ab22-a8d11773568e",
    "controller": [
      "did:cosmos:key:mantra19xvxf55pzvs8ayeuce843sp3rjgmssh77e83fq"
    ],
    "verificationMethod": [
      {
        "controller": "did:cosmos:aumega:900d82bc-2bfe-45a7-ab22-a8d11773568e",
        "id": "did:cosmos:aumega:900d82bc-2bfe-45a7-ab22-a8d11773568e#mantra1axznhnm82lah8qqvp9hxdad49yx3s5dcj66qka",
        "publicKeyHex": "0248a5178d7a90ec187b3c3d533a4385db905f6fcdaac5026859ca5ef7b0b1c3b5",
        "type": "EcdsaSecp256k1VerificationKey2019"
      },
      {
        "controller": "did:cosmos:aumega:900d82bc-2bfe-45a7-ab22-a8d11773568e",
        "id": "did:cosmos:aumega:900d82bc-2bfe-45a7-ab22-a8d11773568e#cosmos1lvl2s8x4pta5f96appxrwn3mypsvumukvk7ck2",
        "publicKeyHex": "03786095e15eb228f4e15692eda6e0607a313cc081ad54d69aadd15d515e304590",
        "type": "EcdsaSecp256k1VerificationKey2019"
      }
    ],
    "service": [],
    "authentication": [
      "did:cosmos:aumega:900d82bc-2bfe-45a7-ab22-a8d11773568e#mantra1axznhnm82lah8qqvp9hxdad49yx3s5dcj66qka",
      "did:cosmos:aumega:900d82bc-2bfe-45a7-ab22-a8d11773568e#cosmos1lvl2s8x4pta5f96appxrwn3mypsvumukvk7ck2"
    ],
    "assertionMethod": [],
    "keyAgreement": [],
    "capabilityInvocation": [],
    "capabilityDelegation": []
  },
  "didMetadata": {
    "versionId": "37467691e6ad832534f5d13df0be3362ec6aeb8cce1f252bb448879e1847de77",
    "created": "2021-08-23T08:24:26.972761898Z",
    "updated": "2021-08-24T14:09:11.322038045Z",
    "deactivated": false
  }
}

```

#### RevokeVerification

The `RevokeVerification` method is used to remove existing [verification methods](https://w3c.github.io/did-core/#verification-methods) and [verification relationships](https://w3c.github.io/did-core/#verification-relationships) from a [DidDocument](02_state.md#diddocument). It accepts a [MsgRevokeVerification](./04_messages.md#MsgRevokeVerification) as a parameter.

The operation will fail if:

- the verification method id is not found

<!--

aumegad tx did revoke-verification-method \
 900d82bc-2bfe-45a7-ab22-a8d11773568e \
 900d82bc-2bfe-45a7-ab22-a8d11773568e#cosmos1lvl2s8x4pta5f96appxrwn3mypsvumukvk7ck2 \
 --from vasp --node https://127.0.0.1:443 --chain-id aumega

-->

```javascript
/* gRPC message */
RevokeVerification(
    MsgRevokeVerification(
        "did:cosmos:aumega:900d82bc-2bfe-45a7-ab22-a8d11773568e",
        "did:cosmos:aumega:900d82bc-2bfe-45a7-ab22-a8d11773568e#cosmos1lvl2s8x4pta5f96appxrwn3mypsvumukvk7ck2",
        "mantra1axznhnm82lah8qqvp9hxdad49yx3s5dcj66qka" // signer
    )
)

/* Resolved DID Document */
{
  "didDocument": {
    "context": [
      "https://www.w3.org/ns/did/v1"
    ],
    "id": "did:cosmos:aumega:900d82bc-2bfe-45a7-ab22-a8d11773568e",
    "controller": [
      "did:cosmos:key:mantra19xvxf55pzvs8ayeuce843sp3rjgmssh77e83fq"
    ],
    "verificationMethod": [
      {
        "controller": "did:cosmos:aumega:900d82bc-2bfe-45a7-ab22-a8d11773568e",
        "id": "did:cosmos:aumega:900d82bc-2bfe-45a7-ab22-a8d11773568e#mantra1axznhnm82lah8qqvp9hxdad49yx3s5dcj66qka",
        "publicKeyHex": "0248a5178d7a90ec187b3c3d533a4385db905f6fcdaac5026859ca5ef7b0b1c3b5",
        "type": "EcdsaSecp256k1VerificationKey2019"
      }
    ],
    "service": [],
    "authentication": [
      "did:cosmos:aumega:900d82bc-2bfe-45a7-ab22-a8d11773568e#mantra1axznhnm82lah8qqvp9hxdad49yx3s5dcj66qka"
    ],
    "assertionMethod": [],
    "keyAgreement": [],
    "capabilityInvocation": [],
    "capabilityDelegation": []
  },
  "didMetadata": {
    "versionId": "d89461469fcac09d7f126c94493af54f58bbac27aae946aeed443b9ac669993d",
    "created": "2021-08-23T08:24:26.972761898Z",
    "updated": "2021-08-24T14:28:31.821486259Z",
    "deactivated": false
  }
}

```


#### SetVerificationRelationships


The `SetVerificationRelationships` method is used to **overwrite** existing [verification relationships](https://w3c.github.io/did-core/#verification-relationships) for a [verification methods](https://w3c.github.io/did-core/#verification-methods) in a [DidDocument](02_state.md#diddocument). It accepts a [MsgSetVerificationRelationships](./04_messages.md#MsgSetVerificationRelationships) as a parameter.

The operation will fail if:

- the verification method id is not found for the target did document
- the verification relationships are empty
- the verification relationships contain unsupported values (according to the did method specification)

<!--

aumegad tx did set-verification-relationships \
 900d82bc-2bfe-45a7-ab22-a8d11773568e \
 900d82bc-2bfe-45a7-ab22-a8d11773568e#mantra1axznhnm82lah8qqvp9hxdad49yx3s5dcj66qka \
 --relationship capabilityInvocation \
 --from vasp --node https://127.0.0.1:443 --chain-id aumega

-->

```javascript
/* gRPC message */
SetVerificationRelationships(
    MsgSetVerificationRelationships(
        "did:cosmos:aumega:900d82bc-2bfe-45a7-ab22-a8d11773568e",
        "did:cosmos:aumega:900d82bc-2bfe-45a7-ab22-a8d11773568e#mantra1axznhnm82lah8qqvp9hxdad49yx3s5dcj66qka",
        ["authentication", "capabilityInvocation"]
        "mantra1axznhnm82lah8qqvp9hxdad49yx3s5dcj66qka" // signer
    )
)

/* Resolved DID Document */
{
  "didDocument": {
    "context": [
      "https://www.w3.org/ns/did/v1"
    ],
    "id": "did:cosmos:aumega:900d82bc-2bfe-45a7-ab22-a8d11773568e",
    "controller": [
      "did:cosmos:key:mantra19xvxf55pzvs8ayeuce843sp3rjgmssh77e83fq"
    ],
    "verificationMethod": [
      {
        "controller": "did:cosmos:aumega:900d82bc-2bfe-45a7-ab22-a8d11773568e",
        "id": "did:cosmos:aumega:900d82bc-2bfe-45a7-ab22-a8d11773568e#mantra1axznhnm82lah8qqvp9hxdad49yx3s5dcj66qka",
        "publicKeyHex": "0248a5178d7a90ec187b3c3d533a4385db905f6fcdaac5026859ca5ef7b0b1c3b5",
        "type": "EcdsaSecp256k1VerificationKey2019"
      }
    ],
    "service": [],
    "authentication": [
        "did:cosmos:aumega:900d82bc-2bfe-45a7-ab22-a8d11773568e#mantra1axznhnm82lah8qqvp9hxdad49yx3s5dcj66qka"
    ],
    "assertionMethod": [],
    "keyAgreement": [],
    "capabilityInvocation": [
      "did:cosmos:aumega:900d82bc-2bfe-45a7-ab22-a8d11773568e#mantra1axznhnm82lah8qqvp9hxdad49yx3s5dcj66qka"
    ],
    "capabilityDelegation": []
  },
  "didMetadata": {
    "versionId": "4f0f8857ab36bdeee0ddb541ea7e7b9cb509d29e1103cc7def44d3d1e8220c22",
    "created": "2021-08-23T08:24:26.972761898Z",
    "updated": "2021-08-24T15:54:40.902858856Z",
    "deactivated": false
  }
}

```

#### AddService

The `AddService` method is used to add a [service](https://w3c.github.io/did-core/#services) in a [DidDocument](02_state.md#diddocument). It accepts a [MsgAddService](./04_messages.md#MsgAddService) as a parameter.

The operation will fail if:

- a service with the same id already present in the did document
- the service definition is invalid (according to the did services specification)

<!--

aumegad tx did add-service \
 900d82bc-2bfe-45a7-ab22-a8d11773568e \
 
TODO

 --from vasp --node https://127.0.0.1:443 --chain-id aumega

-->

```javascript
/* gRPC message */
AddService(
    MsgAddService(
        "did:cosmos:aumega:900d82bc-2bfe-45a7-ab22-a8d11773568e",
        {
            "agent:xyz",
            "DIDCommMessaging",
            "https://agent.xyz/1234",
        }
        "mantra1axznhnm82lah8qqvp9hxdad49yx3s5dcj66qka" // signer
    )
)

/* Resolved DID Document */
{
  "didDocument": {
    "context": [
      "https://www.w3.org/ns/did/v1"
    ],
    "id": "did:cosmos:aumega:900d82bc-2bfe-45a7-ab22-a8d11773568e",
    "controller": [
      "did:cosmos:key:mantra19xvxf55pzvs8ayeuce843sp3rjgmssh77e83fq"
    ],
    "verificationMethod": [
      {
        "controller": "did:cosmos:aumega:900d82bc-2bfe-45a7-ab22-a8d11773568e",
        "id": "did:cosmos:aumega:900d82bc-2bfe-45a7-ab22-a8d11773568e#mantra1axznhnm82lah8qqvp9hxdad49yx3s5dcj66qka",
        "publicKeyHex": "0248a5178d7a90ec187b3c3d533a4385db905f6fcdaac5026859ca5ef7b0b1c3b5",
        "type": "EcdsaSecp256k1VerificationKey2019"
      }
    ],
    "service": [
        {
            "agent:xyz",
            "DIDCommMessaging",
            "https://agent.xyz/1234",
        }
    ],
    "authentication": [
        "did:cosmos:aumega:900d82bc-2bfe-45a7-ab22-a8d11773568e#mantra1axznhnm82lah8qqvp9hxdad49yx3s5dcj66qka"
    ],
    "assertionMethod": [],
    "keyAgreement": [],
    "capabilityInvocation": [
      "did:cosmos:aumega:900d82bc-2bfe-45a7-ab22-a8d11773568e#mantra1axznhnm82lah8qqvp9hxdad49yx3s5dcj66qka"
    ],
    "capabilityDelegation": []
  },
  "didMetadata": {
    "versionId": "3021b47687e682bdd31dac8996537dea14bd0d4e7d90dc618a7f400a3024c048",
    "created": "2021-08-23T08:24:26.972761898Z",
    "updated": "2021-08-24T16:24:40.902858856Z",
    "deactivated": false
  }
}

```

#### DeleteService

The `DeleteService` method is used to remove a [service](https://w3c.github.io/did-core/#services) from a [DidDocument](02_state.md#diddocument). It accepts a [MsgDeleteService](./04_messages.md#MsgDeleteService) as a parameter.

The operation will fail if:

- the service id does not match any service in the did document
<!--

aumegad tx did add-service \
 900d82bc-2bfe-45a7-ab22-a8d11773568e \
 
TODO

 --from vasp --node https://127.0.0.1:443 --chain-id aumega

-->

```javascript
/* gRPC message */
DeleteService(
    MsgDeleteService(
        "did:cosmos:aumega:900d82bc-2bfe-45a7-ab22-a8d11773568e",
        "agent:xyz",
        "mantra1axznhnm82lah8qqvp9hxdad49yx3s5dcj66qka" // signer
    )
)

/* Resolved DID Document */
{
  "didDocument": {
    "context": [
      "https://www.w3.org/ns/did/v1"
    ],
    "id": "did:cosmos:aumega:900d82bc-2bfe-45a7-ab22-a8d11773568e",
    "controller": [
      "did:cosmos:key:mantra19xvxf55pzvs8ayeuce843sp3rjgmssh77e83fq"
    ],
    "verificationMethod": [
      {
        "controller": "did:cosmos:aumega:900d82bc-2bfe-45a7-ab22-a8d11773568e",
        "id": "did:cosmos:aumega:900d82bc-2bfe-45a7-ab22-a8d11773568e#mantra1axznhnm82lah8qqvp9hxdad49yx3s5dcj66qka",
        "publicKeyHex": "0248a5178d7a90ec187b3c3d533a4385db905f6fcdaac5026859ca5ef7b0b1c3b5",
        "type": "EcdsaSecp256k1VerificationKey2019"
      }
    ],
    "service": [],
    "authentication": [
        "did:cosmos:aumega:900d82bc-2bfe-45a7-ab22-a8d11773568e#mantra1axznhnm82lah8qqvp9hxdad49yx3s5dcj66qka"
    ],
    "assertionMethod": [],
    "keyAgreement": [],
    "capabilityInvocation": [
      "did:cosmos:aumega:900d82bc-2bfe-45a7-ab22-a8d11773568e#mantra1axznhnm82lah8qqvp9hxdad49yx3s5dcj66qka"
    ],
    "capabilityDelegation": []
  },
  "didMetadata": {
    "versionId": "5b3fc976d1393bf4a144cdd4d99612b813777a60ca6368dcd396cc687f58c872",
    "created": "2021-08-23T08:24:26.972761898Z",
    "updated": "2021-08-24T17:51:40.902858856Z",
    "deactivated": false
  }
}

```

### Deactivate

Currently not supported
