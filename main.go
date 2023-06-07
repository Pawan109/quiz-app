package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"time"
)

type Quiz struct {
	Que string
	Ans string
}

//problemPUller -> puri csv file ko read krne ke liye
//parseProblem -> csv file mein ek ek problem ko read krne ke liye

func problemPuller(fileName string) ([]Quiz, error) {
	//1.open the file
	csvFile, err := os.Open(fileName) //first opening the csv file
	if err != nil {
		return nil, fmt.Errorf("error in opening the %s file; %s ", fileName, err.Error())
	} else {
		fmt.Println("successfully opened the csv file")
		//2, we wil create a new reader
		csvR := csv.NewReader(csvFile)
		csvLines, err := csvR.ReadAll() //3. it will need to read the file

		// means if you can read all the files than parse the problems
		if err == nil {
			return parseProblem(csvLines), nil //4. call the parseproblem fn -> which will then read each line and seperate que & ans also
		} else {
			return nil, fmt.Errorf("error in reading data in csv"+"format from %s file; %s", fileName, err.Error())
		}

	}
}

func parseProblem(lines [][]string) []Quiz {
	//this funciton is just going to go over the lines and parse them , with Quiz struct

	r := make([]Quiz, len(lines)) //this will make slices of the csv file of length 13 . as we have 13 no. of lines in csv

	for i := 0; i < len(lines); i++ {
		r[i] = Quiz{Que: lines[i][0], Ans: lines[i][1]}
	}

	return r

}

func main() {

	//1. input the name of the file   -> for this we are using flag.String -> first parameter = name (to show in the msg) ,
	fName := flag.String("f", "problems.csv", "path of csv file")

	//2. set the duration of the timer   ( for timer -> go routines & channels)
	timer := flag.Int("t", 30, "timer for the quiz ") //name, default value, and usage of the string

	//now parse these flags
	flag.Parse() //-> so now golang has a file  you want to work with and the duration for how long it has to work

	//3. calling our problem puller function   -> how do we take input from the user  ?
	quizLines, err := problemPuller(*fName)

	//3.1 handle the err
	if err != nil {
		exit(fmt.Sprintf("something went wrong: %s", err.Error()))
	}

	//4. count correct answers
	correctAns := 0

	//5. using duration of timer , we want to initialize the timer
	toObj := time.NewTimer(time.Duration(*timer) * time.Second)

	ansChan := make(chan string)

	//6. loop through the problems we will print questions  & accept answers
problemLoop:

	for i, p := range quizLines {
		var answer string
		fmt.Printf("problem no. %d: %s =", i+1, p.Que)

		//ab answer ko user se lene ke liye -> we need answer channel -> so we will quickly do a go routine for answer -
		go func() {
			fmt.Scanf("%s", &answer) //user ka i/p answer read krke ansChan mein daal diya
			ansChan <- answer
		}()
		select {
		case <-toObj.C: //means if 30secs are over
			fmt.Println()
			break problemLoop

		case iAns := <-ansChan:
			if iAns == p.Ans {
				correctAns++
			}
			if i == len(quizLines)-1 { //means if you answered all the questions .-> close the answer channel
				close(ansChan)
			}
		}

	}
	//7. calculate and print result
	fmt.Printf("your result is %d out of %d\n", correctAns, len(quizLines))

	fmt.Println("press enter to exit")
	<-ansChan //this is your enter . // this means that every time problem 1 is completed  you have already inserted some ans in your ansChan
	//to go to problem 2 , you have to press enter -> so , it means to enter another entry in ansChan

}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
