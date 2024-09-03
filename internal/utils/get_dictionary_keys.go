package utils

func GetDictionaryKeys(data map[string]interface{}) []string {
	keys := []string{}
	for k := range data {
		keys = append(keys, k)
	}
	return keys
}
