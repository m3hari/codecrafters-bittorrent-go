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
	switch {
	case command == "decode":
		result, _, err := decodeBencode(args[1])
		if err != nil {
			return "", err
		}
		jsonOutput, err := json.Marshal(result)
		if err != nil {
			return "", err
		}
		return string(jsonOutput), nil

	case command == "info":
		result, err := getTorrentMetaInfo(args[1])
		if err != nil {
			return "", err
		}
		fmt.Printf("Tracker URL: %v\n", result.announce)
		fmt.Printf("Length: %v", result.info.length)

		return "", nil

	default:
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
