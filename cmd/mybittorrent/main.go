package main

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"unicode"
)

// Stack represents a stack that holds a slice of empty interfaces (to allow for different types)
type Stack struct {
	elements []interface{}
}

// Push adds an element to the top of the stack
func (s *Stack) Push(element interface{}) {
	s.elements = append(s.elements, element)
}

// Pop removes and returns the top element of the stack. Returns an error if the stack is empty.
func (s *Stack) Pop() interface{} {
	topElement := s.elements[len(s.elements)-1]
	s.elements = s.elements[:len(s.elements)-1]
	return topElement
}

// Peek returns the top element of the stack without removing it. Returns an error if the stack is empty.
func (s *Stack) Peek() interface{} {
	return s.elements[len(s.elements)-1]
}

// IsEmpty checks if the stack is empty
func (s *Stack) IsEmpty() bool {
	return len(s.elements) == 0
}

// Size returns the number of elements in the stack
func (s *Stack) Size() int {
	return len(s.elements)
}

func reverse(slice *[]interface{}) {
	length := len(*slice)
	for i := 0; i < length/2; i++ {
		j := length - i - 1
		(*slice)[i], (*slice)[j] = (*slice)[j], (*slice)[i]
	}
}

func main() {

	command := os.Args[1]

	if command == "decode" {
		bencodedValue := os.Args[2]
		stack := &Stack{}
		i := 0
		for i < len(bencodedValue) {
			if bencodedValue[i] == 'l' || bencodedValue[i] == 'd' {
				stack.Push(bencodedValue[i])
				i = i + 1
			} else if bencodedValue[i] == 'e' {
				list := []interface{}{}
				for {
					if reflect.TypeOf(stack.Peek()).Kind() == reflect.Uint8 && stack.Peek().(uint8) == 'd' {
						benMap := make(map[string]interface{})
						for j := len(list) - 1; j >= 0; j -= 2 {
							benMap[list[j].(string)] = list[j-1]
						}
						stack.Pop()
						stack.Push(benMap)
						break
					} else if reflect.TypeOf(stack.Peek()).Kind() == reflect.Uint8 && stack.Peek().(uint8) == 'l' {
						stack.Pop()
						reverse(&list)
						stack.Push(list)
						break
					} else {
						list = append(list, stack.Peek())
						stack.Pop()
					}
				}
				i = i + 1
			} else if unicode.IsDigit(rune(bencodedValue[i])) {
				var firstColonIndex int
				for j := i; j < len(bencodedValue); j++ {
					if bencodedValue[j] == ':' {
						firstColonIndex = j
						break
					}
				}
				lengthStr := bencodedValue[i:firstColonIndex]
				length, _ := strconv.Atoi(lengthStr)

				letter := bencodedValue[firstColonIndex+1 : firstColonIndex+1+length]
				stack.Push(letter)
				i = firstColonIndex + 1 + length
			} else if bencodedValue[i] == 'i' {
				lastIndex := 0
				for j := i + 1; j < len(bencodedValue); j++ {
					if bencodedValue[j] == 'e' {
						lastIndex = j
						break
					}
				}
				num, _ := strconv.Atoi(bencodedValue[i+1 : lastIndex])
				stack.Push(num)
				i = lastIndex + 1
			}
		}
		jsonOutput, _ := json.Marshal(stack.Peek())
		fmt.Println(string(jsonOutput))
	} else {
		fmt.Println("Unknown command: " + command)
		os.Exit(1)
	}
}
