package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func run(args []string) (string, error) {
	if len(args) < 2 {
		return "", fmt.Errorf("usage: <command> <argument>")
	}

	command := args[0]
	bencodedValue := args[1]
	if command == "decode" {
		decoded, _, err := decodeBencode(bencodedValue)
		if err != nil {
			return "", err
		}
		jsonOutput, err := json.Marshal(decoded)
		if err != nil {
			return "", nil
		}
		return string(jsonOutput), nil
	} else {
		return "", fmt.Errorf("Unknown command: " + command)
	}
}

func main() {
	result, err := run(os.Args[1:])
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(result)
	}
}
