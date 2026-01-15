package main

import (
	"fmt"
	"os"
)

func main() {
	// Basic Options
	code := 0
	err := DefiningOptionalValues()
	if err != nil {
		fmt.Println("DefiningOptionalValues example failed")
		code = 1
	}

	err = InspectingValues()
	if err != nil {
		fmt.Println("InspectingValues example failed")
		code = 1
	}

	err = MarshalingExamples()
	if err != nil {
		fmt.Println("MarshalingExamples example failed")
		code = 1
	}

	err = TransformationExamples()
	if err != nil {
		fmt.Println("TransformationExamples example failed")
		code = 1
	}

	if code == 0 {
		fmt.Println("")
		fmt.Println("Examples ran successfully")
	} else {
		fmt.Println("")
		fmt.Println("At least one example failed")
	}
	os.Exit(code)
}
