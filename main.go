package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	// open CSV file
	fd, error := os.Open("moxfield_haves_2023-09-19-0119Z.csv")
	if error != nil {
		fmt.Println(error)
	}
	fmt.Println("Successfully opened the CSV file")
	defer fd.Close()

	// read CSV file
	fileReader := csv.NewReader(fd)
	records, error := fileReader.ReadAll()
	if error != nil {
		fmt.Println(error)
	}
	collection := make(map[string]int, 0)
	for _, record := range records {
		numberOfCopies, err := strconv.Atoi(record[0])
		if err != nil {
			fmt.Println(error)
		}

		if collection[record[2]] == 0 {
			collection[record[2]] = numberOfCopies
		} else {
			collection[record[2]] = collection[record[2]] + numberOfCopies
		}

	}

	file, err := os.Open("card_wishlist_updated_091823.txt")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	wantList := make(map[string]int, 0)
	r := regexp.MustCompile(`^(.*)(\s\dx)$`)
	scanner := bufio.NewScanner(file)
	// optionally, resize scanner's capacity for lines over 64K, see next example
	var lineNumber int
	for scanner.Scan() {
		lineNumber++
		lineContents := scanner.Text()
		if r.MatchString(lineContents) {
			matches := r.FindStringSubmatch(scanner.Text())
			if len(matches) == 3 && matches[2] != "" {
				numberOfCopies, err := strconv.Atoi(matches[2][len(matches[2])-2 : 2])
				if err != nil {
					fmt.Println(err)
				}
				trimmedCardName := strings.Trim(matches[1], " ,")
				if trimmedCardName != "" {
					if _, ok := wantList[trimmedCardName]; ok {
						fmt.Printf("Duplicate requested card: %s on line number %d\n", trimmedCardName, lineNumber)
					} else {
						wantList[trimmedCardName] = numberOfCopies
					}
				}
			}
		} else {
			trimmedFullLine := strings.Trim(lineContents, " ,")
			if trimmedFullLine != "" {
				if _, ok := wantList[trimmedFullLine]; ok {
					fmt.Printf("Duplicate requested card: %s on line number %d\n", trimmedFullLine, lineNumber)
				} else {
					wantList[trimmedFullLine] = 1
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}

	for k, v := range wantList {
		if collection[k] > 0 {
			fmt.Printf("You asked for %d copies of %s, but you already own %d\n", v, k, collection[k])
		}
	}

}
