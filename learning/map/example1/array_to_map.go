package example1

func arr2map(arr [4]string) map[int]string {
	result := make(map[int]string, len(arr))
	for i, val := range arr {
		result[i] = val
	}
	return result
}
