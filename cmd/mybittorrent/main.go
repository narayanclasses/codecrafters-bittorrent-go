package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"unicode"
)

func decodeBencode(bencodedString string, start int) (interface{}, int) {
	if unicode.IsDigit(rune(bencodedString[start])) {
		var firstColonIndex int

		for i := start; i < len(bencodedString); i++ {
			if bencodedString[i] == ':' {
				firstColonIndex = i
				break
			}
		}

		lengthStr := bencodedString[:firstColonIndex]

		length, _ := strconv.Atoi(lengthStr)
		return bencodedString[firstColonIndex+1 : firstColonIndex+1+length], firstColonIndex + 1 + length

	} else if bencodedString[start] == 'i' {
		lastIndex := 0
		for i := start + 1; i < len(bencodedString); i++ {
			if bencodedString[i] == 'e' {
				lastIndex = i
				break
			}
		}
		num, _ := strconv.Atoi(bencodedString[start+1 : lastIndex])
		return num, lastIndex + 1
	} else {
		return "", len(bencodedString)
	}
}

func main() {

	command := os.Args[1]

	if command == "decode" {
		bencodedValue := os.Args[2]

		var decoded interface{}
		slice := []interface{}{}

		i := 0
		nexti := 0
		for i < len(bencodedValue) {
			if bencodedValue[i] == 'l' {
				decoded, nexti = decodeBencode(bencodedValue, i+1)
			} else if bencodedValue[i] != 'e' {
				decoded, nexti = decodeBencode(bencodedValue, i)
			} else {
				nexti = i + 1
			}
			i = nexti
			if decoded != "" {
				slice = append(slice, decoded)
			}
		}
		jsonOutput, _ := json.Marshal(decoded)
		fmt.Println(string(jsonOutput))
	} else {
		fmt.Println("Unknown command: " + command)
		os.Exit(1)
	}
}
