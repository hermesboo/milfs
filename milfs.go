package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"strings"
)

func shellinf() (string, string, string) {
	/* UN means username */
	un, _ := user.Current()
	/* WD is present working directory */
	wd, _ := os.Getwd()
	/* HN is hostname */
	hn, _ := os.Hostname()
	return un.Username, wd, hn
}

func main() {
	/* shell loop */
	shloop()
}

func shloop() {
	readInput := bufio.NewReader(os.Stdin)
	fmt.Printf("Starting My Intentionally Less Friendly Shell\n")

	var status int
	for status != 1 {
		un, wd, hn := shellinf()
		fmt.Printf("%v@%v > %v >> ", un, hn, wd)

		input := CleanInput(readInput)
		fmt.Printf("input %v\n", input)
		//check :=
		ParseInput(input)
		status = HandlingInput(input)
	}
}

/* this function removes the newline from the end of the input */
func CleanInput(reader *bufio.Reader) string {
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Error getting input to clean it: %v \n", err)
	}
	input = strings.TrimSuffix(input, "\n")

	return input
}

func ParseInput(input string) {
	pipe := strings.Contains(input, "|")
	if pipe == true {
		fmt.Printf("Contains at least 1 pipe\n")
	} else {
		fmt.Printf("No pipes here\n")
	}
}

func HandlingInput(input string) int {
	/* Sanitized input aka clean */
	SanInput := strings.Split(input, " ")

	fmt.Printf("saninput type: %T\n", SanInput)

	if SanInput[0] == "cd" {

		user, err := user.Current()
		LoneCD := strings.Join(SanInput[:], " ")
		LoneCD = strings.Trim(LoneCD, " ")
		if LoneCD == "cd" {
			if err == nil {
				os.Chdir(user.HomeDir)
				return 0
			} else {
				fmt.Printf("Couldnt get user, defaulting to /home\n")
				os.Chdir("/home")
			}
		} else {
			SanInput[1] = strings.Replace(SanInput[1], "~", user.HomeDir, 1)
			os.Chdir(SanInput[1])
			return 0
		}

	} else if SanInput[0] == "exit" {
		return 1

	} else if SanInput[0] == "pwd" {
		pwd, err := os.Getwd()
		if err != nil {
			fmt.Printf("Error getting PWD: %v\n", err)
			return 0
		}
		fmt.Printf("%v\n", pwd)
	} else {
		/*
		 * Comodoro is the command to execute
		 * the reason its called comodoro, is because
		 * the name started as command, then commandtoexec,
		 * then commandeer
		 * then comodoro
		 */
		Comodoro := exec.Command(SanInput[0], SanInput[1:]...)
		Comodoro.Stdin = os.Stdin
		Comodoro.Stdout = os.Stdout
		Comodoro.Stderr = os.Stderr

		err := Comodoro.Start()
		Comodoro.Wait()
		if err != nil {
			fmt.Printf("Error at starting process: %v\n", err)
		}
		return 0
	}
	return 0
}