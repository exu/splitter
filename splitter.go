package main

import (
	"flag"
	"fmt"
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

  Parameters:

  [33m<reapeter>[0m - how many jobs generate (numbers passed to command string)
  [33m<command>[0m  - command template (could contain %d placeholder,
                    where number from repeater will be injected

  Options:

    [33m--exit[0m exit on error - default false
    [33m--pool[0m worker pool size - default 4

  Examples:

[32m
    $ splitter --pool=10 1560 get-item %d
[0m

[30m
    running concurrent 10 workers and executing foreach
    - get-item 1
    - get-item 2
    - ....
    - ....
    - get-item 1560
[0m

    You can use ranges in first repeater argument

[32m
    $ splitter 1,2,100-102 get-item %d
[0m

[30m
    running concurrent 4 default workers and executing foreach
    - get-item 1
    - get-item 2
    - get-item 100
    - get-item 101
    - get-item 102
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

func worker(i int, cmdTemplate string, jobs <-chan int, results chan<- int, exit *bool) {
	var cmd string

	for param := range jobs {

		start := time.Now()
		// injecting parameter for loop
		cmd = cmdTemplate
		if strings.Contains(cmdTemplate, "%") {
			cmd = fmt.Sprintf(cmdTemplate, param)
		}

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
			results <- 0

			if *exit {
				os.Exit(1)
			}

			return
		}

		log.Printf("Worker %d completed [32m`%s`[0m in %.3fs (PID:%d)",
			i, cmd, time.Since(start).Seconds(), ps.Pid())

		results <- ps.Pid()
	}
}

func parseRanges(output map[int]int, s string) map[int]int {
	if strings.Contains(s, "-") {
		ranges := strings.Split(s, "-")
		from, _ := strconv.Atoi(ranges[0])
		to, _ := strconv.Atoi(ranges[1])
		for i := from; i <= to; i++ {
			output[i] = i
		}
	} else {
		i, _ := strconv.Atoi(s)
		output[i] = i
	}

	return output
}

func parseRepeater(input string) map[int]int {
	output := map[int]int{}

	if strings.Contains(input, ",") {
		parts := strings.Split(input, ",")
		for _, s := range parts {
			parseRanges(output, s)
		}
	} else if strings.Contains(input, "-") {
		parseRanges(output, input)
	} else {
		num, _ := strconv.Atoi(input)
		for i := 1; i <= num; i++ {
			output[i] = i
		}
	}

	return output
}

func main() {
	flags.Usage = printUsage
	flags.Parse(os.Args[1:])
	argv := flags.Args()

	// repeater
	if len(argv) < 1 {
		check(fmt.Errorf("<repeater> required"))
	}

	repeater := parseRepeater(argv[0])

	if len(repeater) < 1 {
		check(fmt.Errorf("valid <repeater> required (e.g. 56 or 1,3,4,100-102"))
	}

	jobs := make(chan int, len(repeater))
	results := make(chan int, len(repeater))

	// command
	if len(argv) < 2 {
		check(fmt.Errorf("<command> required"))
	}

	// command name and args
	cmd := strings.Join(argv[1:], " ")

	// run workers in concurent subroutines
	log.Printf("Splitting jobs on %d workers", *pool)

	for i := 1; i <= *pool; i++ {
		go worker(i, cmd, jobs, results, exit)
	}

	// adding jobs to job channel
	for j := range repeater {
		jobs <- j
	}

	// collect results - blocks until all results will be filled
	for _, _ = range repeater {
		<-results
	}
}
