package bencode

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"unicode"

	"github.com/codecrafters-io/bittorrent-starter-go/internal/maps"
)

func Marshal(data interface{}) (string, error) {
	switch v := data.(type) {
	case string:
		return fmt.Sprintf("%d:%s", len(v), v), nil
	case int:
		return fmt.Sprintf("i%de", v), nil
	case []interface{}:
		result := "l"
		for _, item := range v {
			encodedItem, err := Marshal(item)
			if err != nil {
				return "", err
			}
			result += encodedItem
		}
		result += "e"
		return result, nil
	case map[string]interface{}:
		keys := maps.Keys(v)
		sort.Strings(keys)

		result := "d"
		for _, key := range keys {
			encodedKey, err := Marshal(key)
			if err != nil {
				return "", err
			}
			encodedValue, err := Marshal(v[key])
			if err != nil {
				return "", err
			}
			result += encodedKey + encodedValue
		}
		result += "e"
		return result, nil
	default:
		return "", fmt.Errorf("unsupported data type")
	}
}

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

	return bencodedString[start:end], bencodedString[end:], nil
}

func decodeInteger(bencodedString string) (value int, remaining string, err error) {
	endIndex := strings.Index(bencodedString, "e")
	if endIndex == -1 {
		return 0, "", fmt.Errorf("invalid bencode integer. Missing closing 'e'")
	}
	value, err = strconv.Atoi(bencodedString[1:endIndex])
	if err != nil {
		return 0, "", err
	}

	return value, bencodedString[endIndex+1:], nil
}

func decodeList(bencodedString string) (value interface{}, remaining string, err error) {
	if len(bencodedString) < 2 || bencodedString[0] != 'l' {
		return "", "", fmt.Errorf("invalid bencode list")
	}

	result := []interface{}{}
	rest := bencodedString[1:]
	for rest[0] != 'e' {
		value, remaining, err = decode(rest)
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

func decodeDictionary(bencodedString string) (value interface{}, remaining string, err error) {
	if len(bencodedString) < 2 {
		return "", "", fmt.Errorf("invalid bencode dictionary")
	}
	result := map[string]interface{}{}
	rest := bencodedString[1:]
	for rest[0] != 'e' {
		key, remaining, err := decode(rest)
		if err != nil {
			return "", "", fmt.Errorf("invalid bencode dictionary. Failed to parse key")
		}
		if remaining == "" {
			return "", "", fmt.Errorf("invalid bencode dictionary. No value for key")
		}

		value, remaining, err = decode(remaining)
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

func Unmarshal(bencodedString string) (value interface{}, err error) {
	result, _, err := decode(bencodedString)
	return result, err
}

func decode(bencodedString string) (value interface{}, remain string, err error) {
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

func ToBencodeDictionary(data interface{}) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	value := reflect.ValueOf(data)
	typ := reflect.TypeOf(data)
	for i := 0; i < value.NumField(); i++ {
		field := typ.Field(i)
		fieldName := field.Tag.Get("bencode")
		if fieldName == "" {
			fieldName = field.Name
		}
		fieldValue := value.Field(i).Interface()

		result[fieldName] = fieldValue
	}
	return result, nil
}
