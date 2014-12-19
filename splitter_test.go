package main

import (
	"fmt"
	"testing"
)

type result map[int]string

func assertMapEquals(expected result, input result) {
	for i, _ := range expected {
		if input[i] != expected[i] {
			fmt.Errorf("Oups! %d indexes are not equal %s != %s", i, input[i], expected[i])
		}
	}
}

func ExampleHello() {
	a := getJobs("1", "ls %s")

	for _, cmd := range a {
		fmt.Print(cmd)
		//Output: ls 1
	}
}

func TestParseRanges(t *testing.T) {
	input := "1-5"
	idx := 1

	output := map[int]string{}
	parseRanges(output, idx, input)

	expected := map[int]string{
		1: "1",
		2: "2",
		3: "3",
		4: "4",
		5: "5",
	}

	assertMapEquals(expected, output)
}

func TestSimpleGetJobs(t *testing.T) {
	jobs := getJobs("1,2", "ls %s")

	cmd1 := jobs[1]
	if cmd1 != "ls 1" {
		t.Errorf("Bad command %s", cmd1)
	}

	cmd2 := jobs[2]
	if cmd2 != "ls 2" {
		t.Errorf("Bad command %s", cmd2)
	}
}

func TestGetJobsWithRangeAndEnum(t *testing.T) {
	output := getJobs("1,2,4-5,mom,dad,son", "ls %s")
	expected := map[int]string{}
	assertMapEquals(expected, output)
}

func TestGetJobsFromFile(t *testing.T) {
	jobs := fileJobs("commands_example.txt")

	if len(jobs) != 6 {
		t.Errorf("There should be 6 jobs in example file not %d", len(jobs))
	}

	// 1 indexed  should be firts command from file
	cmd1 := jobs[1]
	if cmd1 != "ls -la" {
		t.Errorf("Bad command %s - should be 'ls -la'", cmd1)
	}

	cmd5 := jobs[5]
	if cmd5 != "du -sh" {
		t.Errorf("Bad command %s - should be 'du -sh'", cmd5)
	}

}
