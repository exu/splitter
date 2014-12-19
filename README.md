```
            _______________
           <  Job splitter >
            ---------------
                   \   ^__^
                    \  (oo)\_______
                       (__)\       )\/\
                           ||----w |
                           ||     ||

```


## Usage:

`splitter [options] <repeater> <command>`

or

`splitter [options] <file name>`

if first parameter is file name command and reapeater are not needed!

### Parameters:

- `<reapeter>`  - decide how many jobs to run (numbers passed to command string)
- `<command>`   - command template (could contain %s placeholder,
                    where number from repeater will be injected
- `<file name>` - file with command in each line, if you pass file parameter


### Options:

- `--exit` exit on error - default false
- `--pool` worker pool size - default 4
           <repeater> and <command> is not needed anymore


It splits your work to workers pool if you need run some command several times in parallel.
You can pass `<reapeter>` parameter which is responsible for jobs count or runned job command.
Look below for more details.


### Repeater parameter

Repeater can be in different forms:
- `1233` - (int) splitter run 1233 times passed command: `splitter 1233 ls -la`

you can pass %s parameter to command too.

- `1,2,10,20` - int list `splitter mom,dad,son ls -la %s-directory`

splitter will run 4 jobs:
```
ls -la mom-directory
ls -la dad-directory
ls -la son-directory


```
- `1-20` - num range (`splitter 1-20 ls -la directory-%s`)
spliiter will run 20 jobs:
```
ls -la directory-1
ls -la directory-2
...
...
ls -la directory-19
ls -la directory-20
```


- You can merge above two in number chains:
- `1-3,6,8,10-12,extra` - num range (`splitter 1-3,6,8,10-12 ls -la directory-%s`)
```
ls -la directory-1
ls -la directory-2
ls -la directory-3
ls -la directory-6
ls -la directory-8
ls -la directory-10
ls -la directory-11
ls -la directory-12
ls -la directory-extra
```

For each of above examples you can pass `%s` parameter to shell command
splitter will inject repeater part into it before executing.



## Installation

```
go get github.com/exu/splitter
```

run `splitter` without parameters if you need help


## Additional informations

- no string ranges (for me unnecessary)
- remember that shell processed are spwned in goroutines
try to run `splitter 199 echo %s >> out.txt` You'll see
that number aren't ordered
- examples are trivial, it'll be better to run some long running tasks :)
