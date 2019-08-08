SOK protocol
------------

This file described protocol used by Sok to communicate with clients. Sok listen to connections on port specified in config (default 19151). The communication scheme is:

- Client initialize tls connection (Sok force tls because of its authentication method)
- Client send command - yaml formatted string, max 128KB
- Server execute command and send response, then close connection

Command specification:
-----
TODO