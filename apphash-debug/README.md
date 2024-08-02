# Apphash debugging procedure

The purpose of this doc is to document the procedure on how to find the root cause of an apphash mismatch error. There doesn't seem to any standard procedure which documents it (atleast I could not find it). In many networks when the the chain halts due to an apphash mismatch we use snapshots to recover and sync a node depending on which hash has the majority of VP on it without understanding what caused it in the first place. Using this procedure we can pinpoint the exact source of non-determinism in the code.

I'll be simulating an apphash mismatch halted network and underline the procedure on how to find out the reason for the halt.

## Network specification.

I used [Juno](https://github.com/CosmosContracts/juno/) binaries to setup the network. I created a 2 validator devnet with equal voting power.Validator 1 was running [v22.0.0](https://github.com/CosmosContracts/juno/releases/tag/v22.0.0) binary and Validator 2 was running [v22.0.1](https://github.com/CosmosContracts/juno/releases/tag/v22.0.0) binary. I used these binaries because they were tagged as state compatible binaries which they were not and had led to a chain halt on juno mainnet at height [17328771](https://www.mintscan.io/juno/block/17328771). Many validators had not upgraded to `v22.0.1` on mainnet and a state transition triggered the apphash mismatch.

Using the binaries on the validators I triggered an apphash mismatch on the devnet. I'll come back to how I triggered the apphash mismatch later on in the doc. These were in the logs of the validators when the chain halted.

Validator 1 log:

```
ERR prevote step: consensus deems this block invalid; prevoting nil err="wrong Block.Header.AppHash.  Expected 6A84014DFE03092C897C17B979E47B2374C7DF6F0B6A993ABEE1FEE2EE146CE7, got 94800FC8DC32DB79388B5033D40D091B21B870865889A9D0B9E7563DB93C44B9" height=104 module=consensus round=7
```

Validator 2 log:

```
ERR prevote step: consensus deems this block invalid; prevoting nil err="wrong Block.Header.AppHash.  Expected 94800FC8DC32DB79388B5033D40D091B21B870865889A9D0B9E7563DB93C44B9, got 6A84014DFE03092C897C17B979E47B2374C7DF6F0B6A993ABEE1FEE2EE146CE7" height=104 module=consensus round=7
```

## Debugging procedure

### 1) Identifying the offending module

First I used this [go program](https://gist.github.com/freak12techno/845a3061ed65295667c145c05ffd3b23) written by @freak12techno to calculate the apphash of each module in the `applicaton.db` of both the validators.

Apphash calculator on validator 1:

```
# go run main.go ~/.juno/data/
got commitInfo with 34 stores
 index 0: store name 08-wasm, hash e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
 index 1: store name acc, hash 83f40ec1165c0f6573ece501cbe77482a0b8448ee7fb6c2ef227e6ba78d1ab17
 index 2: store name authz, hash e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
 index 3: store name bank, hash 8699243f40720dc1c10f65457f04f7be3e88305e1d571917ab38807dc22d733f
 index 4: store name builder, hash ac382c5935bac36073bc641509301c0e8a708144f0ee38a23cf2fefbbfc7fe4e
 index 5: store name capability, hash f4ecf19b0d1653d817932d41992d82bd99ef2dca7872b57e299a9d416003e2d9
 index 6: store name clock, hash d038eed9140ca6fbcbf311a2dba44445b74a8b3633bd1b0004910db8dc084ff2
 index 7: store name consensus, hash a997a76d689f9e32a4485a829614bd80c535b84ffbfeea28336fc3d189243c4b
 index 8: store name crisis, hash 38eb3cafc1b78ea8ea0179067f2a52c69bcb8c0520bbc587fbf809891966d680
 index 9: store name cw-hooks, hash d0fe4b56fbd485647bea04c48b108700e6a78c0cde7b3d1cec42630cc6030518
 index 10: store name distribution, hash c245598bf8e839bfe0e8a22e17b8ca12a454c6f7d9a69b64188104e586252e13
 index 11: store name drip, hash 7582c07b9b82ea42fb1d6fa35e8ed8f739d770d4827b9124f2345d0809ba7aee
 index 12: store name evidence, hash e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
 index 13: store name feegrant, hash e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
 index 14: store name feeibc, hash e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
 index 15: store name feepay, hash fc68a5cbae2520c0c9cb259354c590524b234135e75f95888fd193a8c7105122
 index 16: store name feeshare, hash 837a3659c6d96c72c107945ba5f7b55f4b048f0c7901fcd476f846d4192d1edf
 index 17: store name globalfee, hash 334164fa5fc292680097c5ff4fef750b10859fb2c4612d66fe6abb438108af4d
 index 18: store name gov, hash 24d34fe9c73dead0f4feb0d3256b4f67036dcfaede3178a5097de11c372f4bac
 index 19: store name hooks-for-ibc, hash e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
 index 20: store name ibc, hash 36ef8c974edeb2e89a446e66030f7585e65ab87d0f8a64020893390f056cd9da
 index 21: store name icacontroller, hash e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
 index 22: store name icahost, hash 3441002ff98191aa7e392618f41135abcb1cd1505175bd9be2a6468f9227009b
 index 23: store name interchainquery, hash aceb3aa23f3679d60e310bcc4879976af4ce2806bd393da9a6616c4d9065b4f9
 index 24: store name mint, hash aa6bf24135001363824b3d43f7ce18f5d3436919f5fc6bf509eb6957c7735da6
 index 25: store name nft, hash e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
 index 26: store name packetfowardmiddleware, hash a4802765c29252d63291835953abd4cf8eb0ef060a08fe015fed79845d18d522
 index 27: store name params, hash c9c2a8aeba6d8c1b922be7e9006a63e37a97e6486ecf7c22a41a2b184e1e2105
 index 28: store name slashing, hash 99c1983d4e2b4ab9318f10125b82ceec376793a76f55ab00240e43da2f3764be
 index 29: store name staking, hash c77433f598167163421aa32ae73a7b5e02b57fa14bd6d82b87d9182f277460d5
 index 30: store name tokenfactory, hash 39ff94a9f1703cbbb52afeff89c94a4a9a0a490c7d916ffbd3fc35158738566c
 index 31: store name transfer, hash 259ace33f3d819fd9cc751968cc83d503a92e8a7e6fd5e3298068249565fd706
 index 32: store name upgrade, hash b3257b87b36c34a36024748a70698d01b03344d7860664c50745159396afdcf8
 index 33: store name wasm, hash c0aaf62b3fcd6220e2420d8eb46ba0622baa7a6491a1608ac13d20fc37a1088d
hash: 6a84014dfe03092c897c17b979e47b2374c7df6f0b6a993abee1fee2ee146ce7
```

Apphash calculator on validator 2:

```
got commitInfo with 34 stores
 index 0: store name 08-wasm, hash e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
 index 1: store name acc, hash 83f40ec1165c0f6573ece501cbe77482a0b8448ee7fb6c2ef227e6ba78d1ab17
 index 2: store name authz, hash e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
 index 3: store name bank, hash eea82414b653f802a9e1127add39e4aac569e72b4b36beb07036a88b3ae726c8
 index 4: store name builder, hash ac382c5935bac36073bc641509301c0e8a708144f0ee38a23cf2fefbbfc7fe4e
 index 5: store name capability, hash f4ecf19b0d1653d817932d41992d82bd99ef2dca7872b57e299a9d416003e2d9
 index 6: store name clock, hash d038eed9140ca6fbcbf311a2dba44445b74a8b3633bd1b0004910db8dc084ff2
 index 7: store name consensus, hash a997a76d689f9e32a4485a829614bd80c535b84ffbfeea28336fc3d189243c4b
 index 8: store name crisis, hash 38eb3cafc1b78ea8ea0179067f2a52c69bcb8c0520bbc587fbf809891966d680
 index 9: store name cw-hooks, hash d0fe4b56fbd485647bea04c48b108700e6a78c0cde7b3d1cec42630cc6030518
 index 10: store name distribution, hash c245598bf8e839bfe0e8a22e17b8ca12a454c6f7d9a69b64188104e586252e13
 index 11: store name drip, hash 7582c07b9b82ea42fb1d6fa35e8ed8f739d770d4827b9124f2345d0809ba7aee
 index 12: store name evidence, hash e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
 index 13: store name feegrant, hash e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
 index 14: store name feeibc, hash e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
 index 15: store name feepay, hash fc68a5cbae2520c0c9cb259354c590524b234135e75f95888fd193a8c7105122
 index 16: store name feeshare, hash 837a3659c6d96c72c107945ba5f7b55f4b048f0c7901fcd476f846d4192d1edf
 index 17: store name globalfee, hash 334164fa5fc292680097c5ff4fef750b10859fb2c4612d66fe6abb438108af4d
 index 18: store name gov, hash 24d34fe9c73dead0f4feb0d3256b4f67036dcfaede3178a5097de11c372f4bac
 index 19: store name hooks-for-ibc, hash e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
 index 20: store name ibc, hash 36ef8c974edeb2e89a446e66030f7585e65ab87d0f8a64020893390f056cd9da
 index 21: store name icacontroller, hash e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
 index 22: store name icahost, hash 3441002ff98191aa7e392618f41135abcb1cd1505175bd9be2a6468f9227009b
 index 23: store name interchainquery, hash aceb3aa23f3679d60e310bcc4879976af4ce2806bd393da9a6616c4d9065b4f9
 index 24: store name mint, hash aa6bf24135001363824b3d43f7ce18f5d3436919f5fc6bf509eb6957c7735da6
 index 25: store name nft, hash e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
 index 26: store name packetfowardmiddleware, hash a4802765c29252d63291835953abd4cf8eb0ef060a08fe015fed79845d18d522
 index 27: store name params, hash c9c2a8aeba6d8c1b922be7e9006a63e37a97e6486ecf7c22a41a2b184e1e2105
 index 28: store name slashing, hash 99c1983d4e2b4ab9318f10125b82ceec376793a76f55ab00240e43da2f3764be
 index 29: store name staking, hash c77433f598167163421aa32ae73a7b5e02b57fa14bd6d82b87d9182f277460d5
 index 30: store name tokenfactory, hash 39ff94a9f1703cbbb52afeff89c94a4a9a0a490c7d916ffbd3fc35158738566c
 index 31: store name transfer, hash 259ace33f3d819fd9cc751968cc83d503a92e8a7e6fd5e3298068249565fd706
 index 32: store name upgrade, hash b3257b87b36c34a36024748a70698d01b03344d7860664c50745159396afdcf8
 index 33: store name wasm, hash c0aaf62b3fcd6220e2420d8eb46ba0622baa7a6491a1608ac13d20fc37a1088d
hash: 94800fc8dc32db79388b5033d40d091b21b870865889a9d0b9e7563db93c44b9
```

From this we I could see that all the module hashes match with each other across the validators except the `bank` module hash. Validator 1 has `6a84014dfe03092c897c17b979e47b2374c7df6f0b6a993abee1fee2ee146ce7` for the `bank` and validator 2 has `94800fc8dc32db79388b5033d40d091b21b870865889a9d0b9e7563db93c44b9` for the `bank` module.

So the `bank` module is the offending module I should focus on.

### 2) Identifying the offending kv store

Now that I know that the `bank` module is the one I should focus on, I used [iaviewer](https://github.com/cosmos/iavl/tree/master/cmd/iaviewer) to inspect the shape of the iavl tree of it.

I piped the output of `iaviewer` to a file on both the validators.

Validator 1:
`iaviewer shape ~/.juno/data/application.db s/k:bank/  > shape-1.txt`

Validator 2:
`iaviewer shape ~/.juno/data/application.db s/k:bank/  > shape-2.txt`

**Note:-** the module name has to be mentioned in a specific format. I wanted to query the bank module so specified it as `s/k:bank/`. If any other module needs to be mentioned like `staking` or `slashing` it needs to be specified as `s/k:staking/`, `s/k:slashing/`etc. If the module name is not mentioned in that format then you'll see the output as

```
Got version: 0
Error reading data: version does not exist
```

After piping the outputs I used the `diff` command to get the difference between both the shapes.

```
# diff shape-1.txt shape-2.txt

1a2,3
>         *4 00666163746F72792F6A756E6F31686B6B706A6B6D6779676A33777579397473637474616E74766B793239787379747374356D7A2F74657374
>       -3 B2B690FE501397C91C4ED9B17A04F45F480EF71D98F401CDAF10BF792753E3A3
3c5
<       -3 DF426B01E303332BC40C69CBE92BEBA05E926EB73D2A765857E453B87F2C906A
---
>     -2 36A2B0B9EFB8576635CC592EB8A2767FA1DA8D0D41185E75CCD86BA5EF8EEA9F
5,7c7,9
<     -2 BEADC1C866B478792DC4C5BE3B737438673D5DA3C542B43F9ED01B8A85942FB2
<       *3 02141CA54005142878EC59D5426B352F19524A806A167374616B65
<   -1 F4D4EF086852E6096F07AC20690A6FDB923EBC5E3D18AE849FED4A6EBB571FEC
---
>       -3 B746B3645DF5B7FFC0CC0AA1EDA3D6C8CD6A9E206DE5486D3765386045A3BA85
>         *4 02141CA54005142878EC59D5426B352F19524A806A167374616B65
>   -1 89CAD3490194AF4E93AACAEB43B37E227EFFD143C52F56C1EFD7796B2DAB0370
11,13c13,19
```

From this I could see that there are a few differences in the shapes.

So I converted the hexadecimal string to ascii to see what it was.

```
$ echo "00666163746F72792F6A756E6F31686B6B706A6B6D6779676A33777579397473637474616E74766B793239787379747374356D7A2F74657374" | xxd -r -p
factory/juno1hkkpjkmgygj3wuy9tscttantvky29xsytst5mz/test
```

And that is the culprit of the chain halt.

## How the apphash was triggered

I first created a new denom on the running devnet using the `tokenfactory` module:

```
junod tx tokenfactory create-denom test --from juno1hkkpjkmgygj3wuy9tscttantvky29xsytst5mz --chain-id test
```

Once the new denom was created I tried minting tokens of the new denom using:

```
junod tx tokenfactory mint 1000000factory/juno1hkkpjkmgygj3wuy9tscttantvky29xsytst5mz/test --from juno1hkkpjkmgygj3wuy9tscttantvky29xsytst5mz --chain-id test
```

Note: I had created the denom as `test` but the `tokenfactory` module registers it as `factory/<creator-addr>/<denom-name>`. So to mint tokens from that denom required me to use the enire path in the command.

This `mint` command was the tx that triggered the apphash mismatch between the two validators. The debugging procedure I used to find the root cause of the apphash also points to this. Comparing both the [binary versions](https://github.com/CosmosContracts/juno/compare/v22.0.0...v22.0.1) I can see that the there are changes in the bankactions.go file which caused the apphash mismatch.
