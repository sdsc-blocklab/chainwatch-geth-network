# Current Issues
* In ```node/Dockerfile``` the encode url is currently set manually, may need to change if nodes are not peering with bootnode. Need to automatically configure this.
* In all Dockerfiles --nat flag needs to be configured automaticall, as to reflect the given ip

# Local Ethereum Network
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