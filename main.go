package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

//Question struct that stores question with answer
type Question struct {
	question string
	answer   string
}

func main() {

	fileName, timeLimit := readArguments()
	f, err := openFile(fileName)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	questions, err := readCSV(f)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if questions == nil {
		return
	}
	score, err := askQuestion(questions, timeLimit)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Your Score %d/%d\n\n", score, len(questions))

}
func readArguments() (fileName string, timeLimit int) {
	csvFileName := flag.String("csv", "problems.csv", "a csv file in the format 'question, answer'")
	limit := flag.Int("limit", 5, "the time limit for the quiz in seconds")
	flag.Parse()
	return *csvFileName, *limit
}

func readCSV(f io.Reader) ([]Question, error) {
	questions, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return nil, err
	}
	amount := len(questions)
	if amount == 0 {
		return nil, fmt.Errorf("no question in file")
	}
	var data []Question
	for _, line := range questions {
		ques := Question{
			question: line[0],
			answer:   strings.TrimSpace(line[1]),
		}
		data = append(data, ques)
	}
	return data, err
}

func openFile(fileName string) (io.Reader, error) {
	return os.Open(fileName)
}

func getInput(input chan string) {
	for {
		in := bufio.NewReader(os.Stdin)
		result, err := in.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		input <- result
	}
}

func askQuestion(questions []Question, timeLimit int) (int, error) {
	totalScore := 0
	timer := time.NewTimer(time.Duration(timeLimit) * time.Second)
	done := make(chan string)

	go getInput(done)

	for i := range questions {
		ans, err := eachQuestion(questions[i].question, questions[i].answer, timer.C, done)
		if err != nil && ans == -1 {
			return totalScore, nil
		}
		totalScore += ans

	}
	return totalScore, nil
}

func eachQuestion(Quest string, answer string, timer <-chan time.Time, done <-chan string) (int, error) {
	fmt.Printf("%s: ", Quest)
	for {
		select {
		case <-timer:
			return -1, fmt.Errorf("time out")
		case ans := <-done:
			score := 0
			if strings.Compare(strings.Trim(strings.ToLower(ans), "\n"), answer) == 0 {
				score = 1
			} else {
				return 0, fmt.Errorf("wrong Answer")
			}

			return score, nil
		}
	}
}
