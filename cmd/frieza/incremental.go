package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	. "github.com/outscale/frieza/internal/common"
)

const (
	ADD_RESOURCE  = "y"
	SKIP_RESOURCE = "n"
	CANCEL        = "q"
	ADD_TYPE      = "a"
	SKIP_TYPE     = "d"
	HELP          = "?"

	InfoColor    = "\033[1;30m%s\033[0m"
	NoticeColor  = "\033[1;36m%s\033[0m"
	WarningColor = "\033[1;33m%s\033[0m"
	ErrorColor   = "\033[1;31m%s\033[0m"
	DebugColor   = "\033[0;32m%s\033[0m"
)

func printIncrementalUsage() {
	log.Printf(fmt.Sprintf(ErrorColor+"\n", "%v - %v"), ADD_RESOURCE, "add this resource")
	log.Printf(fmt.Sprintf(ErrorColor+"\n", "%v - %v"), SKIP_RESOURCE, "skip the resource")
	log.Printf(fmt.Sprintf(ErrorColor+"\n", "%v - %v"), CANCEL, "cancel addition of resources")
	log.Printf(
		fmt.Sprintf(ErrorColor+"\n", "%v - %v"),
		ADD_TYPE,
		"add this resource and all later resources of this type",
	)
	log.Printf(
		fmt.Sprintf(ErrorColor+"\n", "%v - %v"),
		SKIP_TYPE,
		"skip the resource and all later resources of this type",
	)
	log.Printf(fmt.Sprintf(ErrorColor+"\n", "%v - %v"), HELP, "print help")
}

func retrieveNextUserInput(message string, currentStage int, numberStage int) string {
	acceptedChar := []string{ADD_RESOURCE, SKIP_RESOURCE, CANCEL, ADD_TYPE, SKIP_TYPE, HELP}
	for {
		log.Printf("%v\n", message)
		log.Printf(
			fmt.Sprintf(InfoColor, "(%d/%d) Add this resource [%v]? "),
			currentStage,
			numberStage,
			strings.Join(acceptedChar[:], ","),
		)

		reader := bufio.NewReader(os.Stdin)
		response, _ := reader.ReadString('\n')
		response = strings.ReplaceAll(response, "\n", "")
		response = strings.ReplaceAll(response, "\r", "")
		switch response {
		case ADD_RESOURCE, SKIP_RESOURCE, CANCEL, ADD_TYPE, SKIP_TYPE:
			return response
		case HELP:
			printIncrementalUsage()
		default:
			printIncrementalUsage()
		}
	}
}

/*
y - stage this hunk
n - do not stage this hunk
q - cancel update
a - stage this hunk and all later hunks in the type
d - do not stage this hunk or any of the later hunks in the type
? - print help
*/
func incrementalChoice(typeName ObjectType, values []Object) (*[]Object, error) {
	selectObjects := []Object{}
	templateMessage := fmt.Sprintf(InfoColor+"\n"+DebugColor+"\n", "# Type : %v", "+ %v")
	numberStage := len(values)
	for i, value := range values {
		switch retrieveNextUserInput(fmt.Sprintf(templateMessage, typeName, value), i+1, numberStage) {
		case ADD_RESOURCE:
			selectObjects = append(selectObjects, value)
			continue
		case SKIP_RESOURCE:
			continue
		case CANCEL:
			return nil, nil
		case ADD_TYPE:
			selectObjects = append(selectObjects, values[i:]...)
			return &selectObjects, nil
		case SKIP_TYPE:
			return &selectObjects, nil
		default:
			return nil, fmt.Errorf("internal error")
		}
	}
	return &selectObjects, nil
}
