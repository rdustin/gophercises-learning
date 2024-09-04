package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

type problem struct {
	question string
	answer   string
}

var Correct = 0

func main() {
	csvFileName := "problems.csv"
	timeLimit := 30
	flag.StringVar(&csvFileName, "csv", csvFileName, "a csv file in the format of 'question,answer' (default problems.csv)")
	flag.IntVar(&timeLimit, "limit", timeLimit, "the time limit for the quiz in seconds (default 30)")
	flag.Parse()

	file, err := os.Open(csvFileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	lines, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}

	fmt.Println("Press enter to start")
	fmt.Scanln()
	q := make(chan bool)
	t := make(chan bool)

	go func() {
		problems := parseLines(lines)
		doQuiz(problems)
		q <- true
	}()

	go func() {
		timer(timeLimit)
		t <- true
	}()

	select {
	case <-q:
		fmt.Printf("\nYou scored %d out of %d\n", Correct, len(lines))
	case <-t:
		fmt.Printf("\nTime is up. You scored %d out of %d\n", Correct, len(lines))
	}
}

func parseLines(lines [][]string) []problem {
	result := make([]problem, len(lines))

	for i, line := range lines {
		result[i] = problem{
			question: line[0],
			answer:   strings.TrimSpace(line[1]),
		}
	}
	return result
}

func doQuiz(problems []problem) {
	for i, problem := range problems {
		fmt.Printf("Problem: #%d: %s = ", i+1, problem.question)
		var userAnswer string
		fmt.Scanln(&userAnswer)
		if userAnswer == problem.answer {
			Correct++
		}
	}
}

func timer(timeLimit int) {
	time.Sleep(time.Duration(timeLimit) * time.Second)
}
