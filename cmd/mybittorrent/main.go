package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"unicode"
)

func decodeBencode(bencodedString string, start int, end int) ([]interface{}, bool) {
	if start > end {
		return []interface{}{}, false
	}
	if bencodedString[start] == 'l' {
		result, wrapper := decodeBencode(bencodedString, start+1, end-1)

		returnResult := []interface{}{}
		if wrapper {
			returnResult = append(returnResult, result)
		} else {
			i := 0
			for i < len(result) {
				returnResult = append(returnResult, result[i])
				i++
			}
		}
		return returnResult, true
	}
	if unicode.IsDigit(rune(bencodedString[start])) {
		var firstColonIndex int
		for i := start; i < len(bencodedString); i++ {
			if bencodedString[i] == ':' {
				firstColonIndex = i
				break
			}
		}
		lengthStr := bencodedString[start:firstColonIndex]
		length, _ := strconv.Atoi(lengthStr)

		result, wrapper := decodeBencode(bencodedString, firstColonIndex+1+length, end)
		returnResult := []interface{}{bencodedString[firstColonIndex+1 : firstColonIndex+1+length]}
		if wrapper {
			returnResult = append(returnResult, result)
		} else {
			i := 0
			for i < len(result) {
				returnResult = append(returnResult, result[i])
				i++
			}
		}
		return returnResult, false
	} else if bencodedString[start] == 'i' {
		lastIndex := 0
		for i := start + 1; i < len(bencodedString); i++ {
			if bencodedString[i] == 'e' {
				lastIndex = i
				break
			}
		}
		num, _ := strconv.Atoi(bencodedString[start+1 : lastIndex])
		result, wrapper := decodeBencode(bencodedString, lastIndex+1, end)
		returnResult := []interface{}{num}
		if wrapper {
			returnResult = append(returnResult, result)
		} else {
			i := 0
			for i < len(result) {
				returnResult = append(returnResult, result[i])
				i++
			}
		}
		return returnResult, false
	} else {
		return nil, true
	}
}

func main() {

	command := os.Args[1]

	if command == "decode" {
		bencodedValue := os.Args[2]
		decoded, wrapper := decodeBencode(bencodedValue, 0, len(bencodedValue)-1)
		jsonOutput, _ := json.Marshal(decoded)
		if !wrapper {
			jsonOutput, _ = json.Marshal(decoded[0])
		}
		fmt.Println(string(jsonOutput))
	} else {
		fmt.Println("Unknown command: " + command)
		os.Exit(1)
	}
}
