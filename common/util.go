package common

// 判断对象是否在字符串数组中
func InStringArr(s string, arr []string) bool {
	for _, val := range arr {
		if val == s {
			return true
		}
	}
	return false
}
