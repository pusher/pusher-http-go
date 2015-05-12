## 0.2.2 / 2015-05-12

Socket_ids are now validated upon Trigger*Exclusive and channel authentication.

## 0.2.1 / 2015-04-30

Webhook validation uses hmac.Equals to guard against timing attacks.

0.2.0 / 2015-03-30
==================

* A HTTP client is shared between requests to allow configuration. If none is set by the user, the library supplies a default. Allows for pipelining or to change the transport.

0.1.0 / 2015-03-26
==================

*Instantiation of client from credentials, URL or environment variables.
* User can trigger Pusher events on single channels, multiple channels, and exclude recipients
* Authentication of private and presence channels
* Pusher webhook validation
* Querying application state
* Cluster configuration, HTTPS support, timeout configuration.
