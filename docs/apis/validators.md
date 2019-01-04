
# Validators
> validator列表

* Uri:

   `/nodes/{chainId}/validators`

* Method:

    `GET`

* Param

| 参数名   |      类型      |是否必须|默认值|说明|
|----------|:-------------:|:-----:|:---:|:--:|
| chainId | string | yes | - | 链ID|
| height | number | false | 查询最新块的validators |块高度|


* Response
```
"result":[
    {
        "address": "C95E740DBD2285BB2060630FDE24FCD971ADE4C7",
        "pub_key": {
            "type": "tendermint/PubKeyEd25519",
            "value": "XsVLCKdI6FMhQex5gPzZqunXli8fREkZrOgRj6Lipdg="
        },
        "voting_power": "10",
        "accum": "0",
    }
]
```
