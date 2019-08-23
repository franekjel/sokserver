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
This folder contains tasks. Sok use task format similar to sio2 (oioioi). Each task folder constains:
	- config.yml - task configuration. If missed defaults values will be used. Fields:
		- title (string) - task full title
		- memory_limit (uint) - task global memory limit in KB. If not present default_memory_limit value is used
		- memory_limits (map[string]uint) - holds memory limit for given test. If test doesn't have custom memory_limit it use memory_limit value
		- time_limit (uint) - like memory_limit, holds global time limit in ms
		- time_limits (map[string]uint) - like memory_limits, holds custom time limits in ms
	- in - This folder contains test inputs. User submission will be invoked with this data. Each file must have name as $TASK$TEST.in, eg for task "task": task0a.in, task0b.in, task11.in etc.
	- out - This folder contains tests outputs. User submission output will be comparted with this data. Each file must have name like corresponding input file, but with .out extension
	- doc - doc folder contains problem statement named ${TASK}zad.(pdf|html|txt|whatever)
	- prog (optional) - folder contains task solutions 
### Contests folder ($SOK/contests)
This folder contains contests folder.

### Contest folder ($SOK/contests/$CONTEST_NAME)
This folder contains contest configuration file (contests.yml) and rounds folders.
contest.yml - contests config. Contains fields:
    - name (string) - full contest name
    - key (string) - key needed for user to join contest

### Round folder ($SOK/contests)
TODO