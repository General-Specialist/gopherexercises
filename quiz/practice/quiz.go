package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	csvFilename := flag.String("csv", "../problems.csv", "name of problem csv")
	timeLimit := flag.Int("limit", 10, "time limit to solve problems in seconds")
	flag.Parse()

	file, err := os.Open(*csvFilename)
	if err != nil {
		exit(fmt.Sprintf("Failed to open file %s", *csvFilename))
	}
	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		exit(fmt.Sprintf("Failed to read file %s", *csvFilename))
	}

	problems := parseProblems(lines)
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)
	answerCh := make(chan string)
	correct := 0

problemLoop:
	for _, p := range problems {
		fmt.Printf("%s=", p.q)
		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerCh <- answer
		}()
		select {
		case <-timer.C:
			fmt.Println()
			break problemLoop
		case answer := <-answerCh:
			if answer == p.a {
				correct++
			}
		}
	}
	fmt.Printf("You scored %d out of %d", correct, len(problems))
}

type problem struct {
	q string
	a string
}

func parseProblems(lines [][]string) []problem {
	ret := make([]problem, len(lines))
	for i, line := range lines {
		ret[i] = problem{
			q: line[0],
			a: strings.TrimSpace(line[1]),
		}
	}
	return ret
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
