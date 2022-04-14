# Current Issues
1. In ```images/Dockerfile.node``` the encode url is currently set manually, may need to change if nodes are not peering with bootnode. Need to automatically configure this.
2. In all Dockerfiles --nat flag needs to be configured manually to reflect the allocated ip address

---

# Architecture
*The following are the network services used in designing the SUT:*

- Bootnode (Consensus Node)
    - `geth-bootnode`
- Node (Network Interface)
    - `geth-node`
- Miners (Proof-of-Work)
    - `geth-miner-01`
    - `geth-miner-02`
    - `geth-miner-03`

# Accounts
The following accounts were created via the [**geth-cli**](https://geth.ethereum.org/docs/interface/command-line-options) and stored in `volumes/geth-keystore`.

| Service | Address | Private Key | Password |
| - | - | - | - |
| `geth-bootnode` | `9c481da59a4115f91d35bba140955fd79cc4e478` | `de230ef38051468dfb7758c7cfb2cd136895224fd24fcb1525d646861d4c6fb3` | `pass` |
| `geth-node` | `e79914e459ac0cdb83e605cb55a2872379d2cdc7` | `de55ccd2922f839f88533ee7b5e46a7cbc1e338d57c7406a46d5a1064922cad7` | `pass` |
| `geth-miner-01` | `2d93c959cc6dd35ed737e31de42b42f0654b0a42` | `ced06f7c9ce499ec0bbc9d78a162fcbcaa78f0146a6362c41ed1c5168bb6a265` | `pass` |
| `geth-miner-02` | `b81a23bac045acaf16c4bcb1370e8a9b5c6e5297` | `d03399488878127a6eb734317ebe4d28b49704b40142768328ccbd1b2ba51ec1` | `pass` |
| `geth-miner-03` | `4cf5d669f0194624eddd257e407452a0e9941d8e` | `3ca88eee0a2341ff684646eb25e9d2dd77d1a02623ffc60e8cb3567eb8aea346` | `pass` |

# Deploy Local Network
To start the local ethereum network run:
```
$ docker-compose build
$ docker-compose up
```
The local network consists out of multiple parts:
* 1 Bootnode - registers existing nodes on the network, discovery service.
* 3 Miners - Also called **sealers** with proof-of-authority. They validate the blocks. No RPC is exposed as they are required to be unlocked.
* 1 Node - These serve as **transaction relay** and are fullnodes that do not mine, they are locked but have RPC exposed
* [2 Swarm nodes - These nodes make up the **peer-to-peer CDN**]