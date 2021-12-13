FROM ubuntu:xenial

RUN apt-get update \
     && apt-get install -y wget software-properties-common \
     && rm -rf /var/lib/apt/lists/*

WORKDIR "/root"

RUN add-apt-repository -y ppa:ethereum/ethereum

ARG binary
RUN apt-get update \
     && apt-get install -y ethereum

ARG password
ARG privatekey
RUN echo $password > ~/.accountpassword
RUN echo $privatekey > ~/.privatekey
ADD ./genesis.json ./genesis.json
RUN geth --nousb --datadir data init genesis.json
RUN geth --nousb --datadir data account import --password ~/.accountpassword  ~/.privatekey

ENV address=""
ENV bootnodeId=""
ENV bootnodeIp=""
ENV nodekeyhex=""
ENV currentIP=""
ENV bootnodeEnode=""
ENV networkRestrictions=""

CMD exec geth --datadir data --port 30303 --networkid 1158 --bootnodes "$bootnodeEnode" --nousb --nat=extip:$currentIp --verbosity=4 --http --http.addr "0.0.0.0" --http.api "eth,web3,net,admin,debug,db" --http.corsdomain  "*" --syncmode=full --mine --miner.threads=1 --miner.gasprice "0" --miner.etherbase $address --unlock $address --password ~/.accountpassword --allow-insecure-unlock --netrestrict $networkRestrictions


# CMD exec geth --datadir data --port 30303 --networkid 1158 --bootnodes enode://dbf61f64c260ca5cadea14199c2ba8dc7daf04f055b3709885414f12ce773cf474ca1fd90e3eee86cf0305f574ae6907a4e30ab2abb747356900659faeb83d7d@172.25.0.101:30303 --nousb --nat=extip:172.25.0.105 --rpc --rpcaddr "0.0.0.0" --rpcapi "eth,web3,net,admin,debug,db" --rpccorsdomain "*" --syncmode="full"

EXPOSE 8545
EXPOSE 30303/udp
EXPOSE 30303/tcp