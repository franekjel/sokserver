Sok files and folders
---
This file described files and folders used by Sok. All Sok files use yaml to store data so it is easy to modify they manually

### Main folder ($SOK)
This is main Sok data folder specified by -p flag. Default it is the Sok executable location (it is probably simples solution, which allows to cerate portable Sok installation). It contains Sok global config (sok.yml) and folders contests, queue, tasks and users.

sok.yml - sok global config. Contains fields:
    - port (uint16) - port on which Sok listen - default 19151
    - workers (uint16) - number of threads Sok use to check submissions (default 4)
    - default_memory_limit (uint) - memory limit (in KB) for task that has no config in (default 32768)
    - default_time_limit (uint) - time limit (in ms) for task that has no config (default 1000)
In case there is no config it will be generated automatically with default values

### Queue folder ($SOK/queue)
This folder contains submissions queued to check. It can contains unchecked submissions files.

### Users folder ($SOK/users)
This folder contains users folders.

### User folders ($SOK/users/$LOGIN)
This folder contains user data and submissions. TODO

### Tasks folder ($SOK/tasks)
This folder contains tasks. Sok supports task format used by sio2 (oioioi) TODO

### Contests folder ($SOK/contests)
This folder contains contests folder.

### Contest folder ($SOK/contests/$CONTEST_NAME)
Rhis folder contains contest configuration file (contests.yml) and rounds folders.
contest.yml - contests config. Contains fields:
    - name (string) - full contest name
    - key (string) - key needed for user to join contest

### Round folder ($SOK/contests)
TODO