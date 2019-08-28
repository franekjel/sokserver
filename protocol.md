SOK protocol
------------

This file described protocol used by Sok to communicate with clients. Sok listen to connections on port specified in config (default 19151). The communication scheme is:

- Client initialize tls connection (Sok force tls because of its authentication method)
- Client send command - yaml formatted string, max 128KB
- Server execute command and send response, then close connection

### Command fields

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

### Return message status

Return message always contains fields "status" if command was succesful
Status is "ok" if command was executed successful. Otherwise it contains error message

Commands:
 - create_account: Creates new account with specified login and password. If login is already used or not match criteria command may failed. Example:
    ```
    login: franekjel
    password: P@ssword
    command: create_account
    ```
 - get_task: Get problem statement for given task. Requires contest round and task field. Return message has fields filename and data which contains name of file and file content encoded in base64. Example of return message:
 	```
	status: ok
	filename: taskzad.txt
	data: UHJvYmxlbSBzdGF0ZW1lbnQgZm9yIHRhc2s=
 	```
 - submit: Submits solution. Requires contest, round, task and data. Data field contains code. Example:
	```
	login: franekjel
	password: P@ssword
	command: submit
	contest: con1
	round: round1
	task: task1
	data: '#include <stdio.h>\n int main(){printf("answer");}'
	```
 - contest_ranking: Get the ranking of given contests. Requires contest field. Return message contains additonal map field "contest_ranking":
	```
	status: ok
	contest_ranking:
		foo: 45
		franekjel: 120
		bar: 80
	``` 
 - round_ranking: Get the ranking of given round. Requires contest and round field. Return message contains additional fields: "tasks" - tasks names, "users" - user names and "round_ranking" - two dimensional array of points. Example:
 	```
	status: ok
    tasks:
    - task1
    - task2
    - task3
    users:
    - franekjel
    - foo
    - bar
    round_ranking: [[100, 0, 50], [0, 60, 0], [0, 0, 20]]
 	```
 	In this example user franekjel get 100 points for task1 and 50 for task3, user foo get 60 points for task2 and user bar 20 points for task3

 - list_submissions: Get list of submissions in given round. Requires contest and round fields. Return message contains additional list field submissions. 
 Each of elements is list with three values - submission ID, status (eg. OK, TIMEOUT) and points. If resuls are not present yet points will be 0.
 	```
 	status: ok
 	submissions: [['15be7c9cec0ef768', 'OK', 100], ['15be7c9cec0ab023', 'TIMEOUT', 78]]
 	``` 
 - get_submission: Get given submission. Requires contest, round, task and data fields. Data field contains submission ID. Example:
 	```
	login: franekjel
	password: P@ssword
	command: submit
	contest: con1
	round: round1
	task: task1
	data: '15be7c9cec0ef768'
 	``` 
	Return message contains yamled submission struct (as in tasks/submissions.go).