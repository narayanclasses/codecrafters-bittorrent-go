package main

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
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

func calculateSHA1(input []byte) string {
	sha1Hash := sha1.New()
	sha1Hash.Write(input)
	hashBytes := sha1Hash.Sum(nil)
	sha1String := fmt.Sprintf("%x", hashBytes)
	return sha1String
}

func getHexValue(input []byte) string {
	return fmt.Sprintf("%x", input)
}

var tracker string
var fileLength int
var pieceLength int
var piecesHash string
var infoHash string

func decodeString(bencodedValue string) string {
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
						if list[j].(string) == "announce" {
							tracker = list[j-1].(string)
						}
						if list[j].(string) == "length" {
							fileLength = list[j-1].(int)
						}
						if list[j].(string) == "piece length" {
							pieceLength = list[j-1].(int)
						}
						if list[j].(string) == "pieces" {
							for k := 0; k < len(list[j-1].(string)); k += 20 {
								piecesHash += "\n" + getHexValue([]byte((list[j-1].(string))[k:k+20]))
							}
						}
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
	return string(jsonOutput)
}

func fillInfo(fileName string) {
	content, _ := os.ReadFile(fileName)
	bencodedValue := string(content)
	decodeString(bencodedValue)
	for i := 0; i < len(bencodedValue); i++ {
		if bencodedValue[i:i+4] == "info" {
			infoHash = calculateSHA1([]byte(bencodedValue[i+4 : len(bencodedValue)-1]))
			break
		}
	}
}

func main() {

	command := os.Args[1]
	fileName := os.Args[2]
	if command == "decode" {
		bencodedValue := os.Args[2]
		fmt.Println(decodeString(bencodedValue))
	} else if command == "info" {
		fillInfo(fileName)
		fmt.Printf("Tracker URL: %s\nLength: %d\nInfo Hash: %s\nPiece Length: %d\nPiece Hashes:%s\n", tracker, fileLength, infoHash, pieceLength, piecesHash)
	} else if command == "peers" {
		fillInfo(fileName)
		// Query parameters
		params := url.Values{}
		params.Add("info_hash", "%d6%9f%91%e6%b2%ae%4c%54%24%68%d1%07%3a%71%d4%ea%13%87%9a%7f")
		params.Add("peer_id", "00112233445566778899")
		params.Add("port", "6881")
		params.Add("uploaded", "0")
		params.Add("downloaded", "0")
		params.Add("left", fmt.Sprint(fileLength))
		params.Add("compact", "1")

		// Construct the final URL with query parameters
		finalURL := fmt.Sprintf("%s?%s", tracker, params.Encode())

		// Making the GET request
		response, _ := http.Get(finalURL)
		defer response.Body.Close()
		body, _ := io.ReadAll(response.Body)
		fmt.Println(string(body))
	} else {
		fmt.Println("Unknown command: " + command)
		os.Exit(1)
	}
}
