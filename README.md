**mooZ** - Peer-to-Peer Voice Chat in the Browser, Over I2P
===========================================================

**EXPERIMENTAL.** *Like really, really experimental*. This is an application
which sets up a WebRTC voice chat inside of I2P, with services for peer discovery
and TURN built in. Just running it and sharing the URL is enough to set up your own
group meetings entirely inside I2P(Peer discovery won't work outside of I2P, it is
configured to never contact the clearnet by building all connections through the
[SAMv3 API](https://geti2p.net/en/docs/api/samv3)). It works in both TCP-like(I2P Streaming)
and UDP-like(I2P Datagram) modes, so if you have a SOCKS5 proxy configured which supports
UDP you can use repliable I2P Datagrams to incur less overhead.

- [Enable the SAMv3 API](https://geti2p.net/en/docs/api/samv3).

**DISCLAIMER: It is *deeply ill advised* to use this software for anything right now**
--------------------------------------------------------------------------------------

```sh
go build
./mooz
```
