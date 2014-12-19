package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const usage = `
[33m
            _______________
           <  Job splitter >
            ---------------
                   \   ^__^
                    \  (oo)\_______
                       (__)\       )\/\
                           ||----w |
                           ||     ||
[0m

  Usage: [32msplitter [options] <repeater> <command>[0m
     or: [32msplitter [options] <file name>[0m

  Parameters:

    [33m<reapeter>[0m  - how many jobs generate (numbers passed to command string)
    [33m<command>[0m   - command template (could contain %s placeholder,
               where number from repeater will be injected)
    [33m<file name>[0m - file with command in each line

  Options:

    [33m--exit[0m exit on error - default false
    [33m--pool[0m worker pool size - default 4

  Examples:
[32m
    $ splitter --pool=10 1560 get-item %s
[0m
[30m
    running concurrent 10 workers and executing

    - get-item 1
    - get-item 2
    - ....
    - ....
    - get-item 1560
[0m
    You can use ranges in first repeater argument
[32m
    $ splitter 100-102,mom,dad,son get-item %s
[0m
[30m
    running concurrent 4 default workers and executing

    - get-item 100
    - get-item 101
    - get-item 102
    - get-item mom
    - get-item dad
    - get-item son
[0m

`

var flags = flag.NewFlagSet("every", flag.ContinueOnError)
var exit = flags.Bool("exit", false, "")
var pool = flags.Int("pool", 4, "")

func printUsage() {
	fmt.Println(usage)
	os.Exit(0)
}

func check(err error) {
	if err != nil {
		fmt.Println(usage)
		log.Fatalf("Error: %s\n", err)
	}
}

func worker(i int, queue <-chan string, results chan<- int, exit *bool) {
	for cmd := range queue {
		start := time.Now()
		// injecting parameter for loop
		log.Printf("Worker %d runs [32m`%s`[0m", i, cmd)

		// running process
		proc := exec.Command("/bin/sh", "-c", cmd)
		proc.Stdout = os.Stdout
		proc.Stderr = os.Stderr
		proc.Start()

		err := proc.Wait()
		ps := proc.ProcessState

		if err != nil {
			log.Printf("pid %d failed with %s", ps.Pid(), ps.String())
			if *exit {
				os.Exit(1)
			}
		} else {
			log.Printf("Worker %d completed [32m`%s`[0m in %.3fs (PID:%d)",
				i, cmd, time.Since(start).Seconds(), ps.Pid())
		}

		results <- ps.Pid()
	}
}

func parseRanges(output map[int]string, idx int, s string) int {
	if strings.Contains(s, "-") {
		ranges := strings.Split(s, "-")
		from, _ := strconv.Atoi(ranges[0])
		to, _ := strconv.Atoi(ranges[1])
		for i := from; i <= to; i++ {
			idx++
			output[idx] = strconv.Itoa(i)
		}
	} else {
		idx++
		output[idx] = s
	}

	return idx
}

func getJobs(input string, cmdTemplate string) map[int]string {
	var cmd string

	output := map[int]string{}
	idx := 0

	if strings.Contains(input, ",") {
		parts := strings.Split(input, ",")
		for _, s := range parts {
			idx = parseRanges(output, idx, s)
		}
	} else if strings.Contains(input, "-") {
		idx = parseRanges(output, idx, input)
	} else {
		num, _ := strconv.Atoi(input)
		for i := 1; i <= num; i++ {
			idx++
			output[idx] = strconv.Itoa(i)
		}
	}

	// pass parameters to templates
	for idx, val := range output {
		cmd = cmdTemplate
		if strings.Contains(cmdTemplate, "%") {
			cmd = fmt.Sprintf(cmdTemplate, val)
		}

		output[idx] = cmd
	}

	return output
}

func fileJobs(file string) map[int]string {
	var output = map[int]string{}

	contents, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalf("Something wrong with '%s', pass another one", file)
	}

	for i, cmd := range strings.Split(string(contents), "\n") {
		if cmd != "" {
			output[i+1] = cmd
		}
	}

	return output
}

func main() {
	var jobs map[int]string

	flags.Usage = printUsage
	flags.Parse(os.Args[1:])
	argv := flags.Args()

	// we can pass filename or repeater with command combo
	if len(argv) == 1 {
		_, err := os.Open(argv[0]) // For read access.
		if err != nil {
			log.Fatal("File '%s' can't be opened")
		}

		jobs = fileJobs(argv[0])
	} else {
		// repeater
		if len(argv) < 1 {
			check(fmt.Errorf("<repeater> required"))
		}

		// command
		if len(argv) < 2 {
			check(fmt.Errorf("<command> required"))
		}

		cmdTemplate := strings.Join(argv[1:], " ")
		jobs = getJobs(argv[0], cmdTemplate)

		if len(jobs) < 1 {
			check(fmt.Errorf("valid <repeater> required (e.g. 56 or 1,3,4,100-102"))
		}
	}

	// prepeare channels
	queue := make(chan string, len(jobs))
	results := make(chan int, len(jobs))

	// run workers in concurent subroutines
	log.Printf("Splitting jobs on %d workers", *pool)

	for i := 1; i <= *pool; i++ {
		go worker(i, queue, results, exit)
	}

	// adding jobs to job queue
	for _, cmd := range jobs {
		queue <- cmd
	}

	// collect results - blocks until all results will be filled
	for _, _ = range jobs {
		<-results
	}
}
