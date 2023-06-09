package utils

func GetKeys(inputMap map[string]bool) []*string {
	keys := make([]*string, 0, len(inputMap))
	for k, _ := range inputMap {
		keys = append(keys, &k)
	}
	return keys
}
