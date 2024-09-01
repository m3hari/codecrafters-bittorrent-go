package main

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

func decodeString(bencodedString string) (value string, remaining string, err error) {
	colonIndex := strings.Index(bencodedString, ":")
	if colonIndex == -1 {
		return "", "", fmt.Errorf("invalid bencode string. Missing colon")
	}

	length, err := strconv.Atoi(bencodedString[:colonIndex])
	if err != nil {
		return "", "", err
	}

	start := colonIndex + 1
	end := start + length
	if end > len(bencodedString)+1 {
		return "", "", fmt.Errorf("invalid bencode string. Length is greater than actual string length")
	}

	value = bencodedString[start:end]
	remaining = bencodedString[end:]
	return value, remaining, nil
}

func decodeInteger(bencodedString string) (value int, remaining string, err error) {
	endIndex := strings.Index(bencodedString, "e")
	if endIndex == -1 {
		return 0, "", fmt.Errorf("Invalid bencode integer. Missing closing 'e'")
	}
	value, err = strconv.Atoi(bencodedString[1:endIndex])
	if err != nil {
		return 0, "", err
	}

	return value, bencodedString[endIndex+1:], nil
}

func decodeList(bencodedString string) (value any, remaining string, err error) {
	if len(bencodedString) < 2 || bencodedString[0] != 'l' {
		return "", "", fmt.Errorf("invalid bencode list")
	}

	result := []any{}
	rest := bencodedString[1:]
	for rest[0] != 'e' {
		value, remaining, err = decodeBencode(rest)
		if err != nil {
			return "", "", fmt.Errorf("invalid bencode list")
		}
		result = append(result, value)
		rest = remaining

		if remaining == "" {
			return "", "", fmt.Errorf("invalid bencode list. No closing 'e'")
		}
	}

	return result, rest[1:], nil

}

func decodeDictionary(bencodedString string) (value any, remaining string, err error) {
	if len(bencodedString) < 2 {
		return "", "", fmt.Errorf("invalid bencode dictionary")
	}

	result := map[string]any{}
	rest := bencodedString[1:]
	for rest[0] != 'e' {
		key, remaining, err := decodeBencode(rest)
		if err != nil {
			return "", "", fmt.Errorf("invalid bencode dictionary. Failed to parse key")
		}
		if remaining == "" {
			return "", "", fmt.Errorf("invalid bencode dictionary. No value for key")
		}

		value, remaining, err = decodeBencode(remaining)
		if err != nil {
			return "", "", fmt.Errorf("invalid bencode dictionary. Failed to parse value for key '%v'", key)
		}

		strKey, ok := key.(string)
		if !ok {
			return "", "", fmt.Errorf("invalid bencode dictionary. Key must be string")
		}

		result[strKey] = value

		rest = remaining

		if remaining == "" {
			return "", "", fmt.Errorf("invalid bencode dictionary. No closing 'e'")
		}
	}

	return result, rest[1:], nil

}

func decodeBencode(bencodedString string) (value any, remain string, err error) {
	firstChar := bencodedString[0]
	switch {
	case unicode.IsDigit(rune(firstChar)):
		return decodeString(bencodedString)
	case firstChar == 'i':
		return decodeInteger(bencodedString)
	case firstChar == 'l':
		return decodeList(bencodedString)
	case firstChar == 'd':
		return decodeDictionary(bencodedString)
	default:
		return "", "", fmt.Errorf("invalid bencode input")
	}

}
