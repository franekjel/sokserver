SOK protocol
------------

This file described protocol used by Sok to communicate with clients. Sok listen to connections on port specified in config (default 19151). The communication scheme is:

- Client initialize tls connection (Sok force tls because of its authentication method)
- Client send command - yaml formatted string, max 128KB
- Server execute command and send response, then close connection

Command fields
-----
- login: user login. In case user isn't registered and command is create_accout this is desired user login. This field must be present. Example:
"login: franekjel"
- password: user password. If user isn't registered yet and command is create_accout this is new user password upon succesful registration. This field must be present. Example:
 "password: P@ssword"
 - command: desired operation. There is some command avaible (described below) that can require additional fields. This field must be present. Example
 "command: create_account"
- contest: contest ID at which the operation will be performend. Necessary for some commands. Example:
"contest: con1"
- round: round ID at which the operation will be performend (like contest). Necessary for some commands. Example:
"round: rnd1"
- task: task ID, necessary for some command. Example:
"task: task1"

Return message status
---
Return message always contains fields "status" if command was succesful
Status is "ok" if command was executed successful. Otherwise it contains error message

Commands:
 - create_account: Creates new account with specified login and password. If login is already used or not match criteria command may failed. Example:
    ```
    login: franekjel
    password: P@ssword
    command: create_account
    ```
Todo:
- submit
- contest_ranking
- round_ranking
- list_submissions
- last_submission
- get_submission