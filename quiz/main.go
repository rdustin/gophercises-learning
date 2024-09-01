package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
)

func main() {
	correct := 0
	incorrect := 0
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

	for _, line := range lines {
		question := line[0]
		answer := line[1]
		fmt.Println(question)
		var userAnswer string
		fmt.Scanln(&userAnswer)
		if userAnswer == answer {
			correct++
		} else {
			incorrect++
		}
	}

	fmt.Printf("Correct: %d out of: %d\n", correct, len(lines))
}
