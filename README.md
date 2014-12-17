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


### Parameters:

- `<reapeter>` - how many jobs generate (numbers passed to command string)
- `<command>`  - command template (could contain %d placeholder,
                    where number from repeater will be injected

### Options:

- `--exit` exit on error - default false
- `--pool` worker pool size - default 4


It splits your work to workers pool if you need run some command several times in parallel.
You can pass `<reapeter>` parameter which is responsible for jobs count or runned job command.
Look below for more details.


### Repeater parameter

Repeater can be in different forms:
- `1233` - (int) splitter run 1233 jobs: `splitter 1233 ls -la`
- `1,2,10,20` - int list `splitter 1,2,10,20 ls -la directory-%d`
splitter will run 4 jobs:
```
ls -la directory-1
ls -la directory-2
ls -la directory-10
ls -la directory-20

```
- `1-20` - num range (`splitter 1-20 ls -la directory-%d`)
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
- `1-3,6,8,10-12` - num range (`splitter 1-3,6,8,10-12 ls -la directory-%d`)
```
ls -la directory-1
ls -la directory-2
ls -la directory-3
ls -la directory-6
ls -la directory-8
ls -la directory-10
ls -la directory-11
ls -la directory-12
```

For each of above examples you can pass `%d` parameter to shell command
aplitter will inject it before executing.



## Installation

```
go get github.com/exu/splitter
```


run `splitter` without parameters if you need help
