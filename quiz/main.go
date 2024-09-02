package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
)

type problem struct {
	question string
	answer   string
}

func main() {
	correct := 0
	incorrect := 0
	csvFileName := "problems.csv"
	// timeLimit := 30
	flag.StringVar(&csvFileName, "csv", csvFileName, "a csv file in the format of 'question,answer' (default problems.csv)")
	// flag.IntVar(&timeLimit, "limit", timeLimit, "the time limit for the quiz in seconds (default 30)")
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

	problems := parseLines(lines)

	for i, problem := range problems {
		fmt.Printf("Problem: #%d: %s = ", i+1, problem.question)
		var userAnswer string
		fmt.Scanln(&userAnswer)
		if userAnswer == problem.answer {
			correct++
		} else {
			incorrect++
		}
	}

	fmt.Printf("Correct: %d out of: %d\n", correct, len(lines))
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
