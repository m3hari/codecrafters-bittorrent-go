package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"unicode"
)

// Example:
// - 5:hello -> hello
// - 10:hello12345 -> hello12345

func decodeBencode(bencodedString string) (interface{}, error) {
	lengthOfStr := len(bencodedString)
	firstChar := bencodedString[0]
	lastChar := bencodedString[lengthOfStr-1:][0]

	if unicode.IsDigit(rune(firstChar)) {
		var firstColonIndex int
		for i := 0; i < lengthOfStr; i++ {
			if bencodedString[i] == ':' {
				firstColonIndex = i
				break
			}
		}

		lengthStr := bencodedString[:firstColonIndex]

		length, err := strconv.Atoi(lengthStr)
		if err != nil {
			return "", err
		}

		return bencodedString[firstColonIndex+1 : firstColonIndex+1+length], nil
	} else if firstChar == 'i' && lastChar == 'e' {
		possibleNumber := bencodedString[1 : lengthOfStr-1]
		value, err := strconv.Atoi(possibleNumber)
		if err != nil {
			fmt.Println("Invalid value for integer encoding")
			return nil, err
		}
		return value, nil

	} else {
		return "", fmt.Errorf("Only strings & numbers are supported at the moment")
	}
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: mybittorrent <command> <argument>")
		os.Exit(1)
	}

	command := os.Args[1]
	if command == "decode" {
		bencodedValue := os.Args[2]
		decoded, err := decodeBencode(bencodedValue)
		if err != nil {
			fmt.Println(err)
			return
		}

		jsonOutput, _ := json.Marshal(decoded)
		fmt.Println(string(jsonOutput))
	} else {
		fmt.Println("Unknown command: " + command)
		os.Exit(1)
	}
}
