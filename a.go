package main

import "strings"

func main() {
	route := "/aggg/bad/afwe"
	routeArr := strings.Split(route, "/")
	routeLen := len(routeArr)
	if routeLen < 3 {
		return
	}
	newRoute := ""
	for i := 1; i < routeLen; i++ {
		newRoute += strings.Title(routeArr[i])
		if i == 1 {
			newRoute += "_"
		}
	}
	print(newRoute)
}
