package main

import "encoding/json"

func marshallToString(obj interface{}) string {
	marshalledBytes, _ := json.Marshal(obj)
	return string(marshalledBytes)
}

func removeEmptyFromArray(s []string) []string {
	var r []string
	for _, str := range s {
			if str != "" {
					r = append(r, str)
			}
	}
	return r
}
