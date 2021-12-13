FROM ubuntu:xenial

RUN apt-get update \
  && apt-get install -y wget software-properties-common \
  && rm -rf /var/lib/apt/lists/* \
  && add-apt-repository -y ppa:ethereum/ethereum \
  && apt-get update \
  && apt-get install -y ethereum

WORKDIR "/root"

ARG password
ARG privatekey
RUN echo $password > ~/.accountpassword
RUN echo $privatekey > ~/.privatekey
ADD ./genesis.json ./genesis.json
RUN geth --nousb --datadir data init genesis.json
RUN geth --nousb --datadir data account import --password ~/.accountpassword  ~/.privatekey

ENV address=""
ENV bootnodeIp=""
ENV nodekeyhex=""
ENV networkRestrictions=""
ENV networkId=""

CMD exec geth --datadir data --nodekeyhex $nodekeyhex --networkid 1158 --nat upnp --nousb --nat=extip:172.25.0.101 --netrestrict $networkRestrictions

# CMD exec geth --datadir data --networkid 1158 --nat extip:172.25.0.101 --nousb

EXPOSE 30301/udp
EXPOSE 30303/udp
EXPOSE 30303/tcp
