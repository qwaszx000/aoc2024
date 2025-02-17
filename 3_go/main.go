package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// Parse mul(num,num) parts in line, multipy and sum results
// num can be up to 3 chars
/*func parse_line(line string) uint {
	var mul_word_status uint8 = 0
	const required_mul_word_status = 4

	var num1, num2 uint = 0, 0

	for _, c := range line {
		if mul_word_status == 0 && c == 'm' {
			mul_word_status += 1
			continue
		} else if mul_word_status == 1 && c == 'u' {

		}
	}
}*/

func parse_line(line string) uint64 {
	var result_sum uint64 = 0

	re := regexp.MustCompile(`mul\((\d{1,3}),(\d{1,3})\)`)

	for _, matches := range re.FindAllStringSubmatch(line, -1) {
		num1, _ := strconv.ParseUint(matches[1], 10, 0)
		num2, _ := strconv.ParseUint(matches[2], 10, 0)

		result_sum += num1 * num2
	}

	return result_sum
}

func parse_line2(line string, is_allowed bool) (uint64, bool) {
	var result_sum uint64 = 0
	const (
		DONT_STR = `don't()`
		DO_STR   = `do()`
	)

	for current_index := 0; current_index < len(line); {
		next_do_index := strings.Index(line[current_index:], DO_STR)
		next_dont_index := strings.Index(line[current_index:], DONT_STR)

		if next_dont_index == -1 && is_allowed {
			result_sum += parse_line(line[current_index:])
			break
		}

		if next_do_index == -1 && !is_allowed {
			break
		}

		next_do_index += current_index
		next_dont_index += current_index

		if is_allowed {
			result_sum += parse_line(line[current_index:next_dont_index])
			current_index = next_dont_index + len(DONT_STR)
			is_allowed = false
		} else if next_do_index < next_dont_index { //small optimization
			result_sum += parse_line(line[next_do_index:next_dont_index])
			current_index = next_dont_index + len(DONT_STR)
			is_allowed = false
		} else {
			current_index = next_do_index + len(DO_STR)
			is_allowed = true
		}
	}

	return result_sum, is_allowed
}

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("Usage: %s <input_file_path>\n", os.Args[0])
	}

	//Read file name from argv
	input_filepath_rel := os.Args[1]
	fmt.Printf("Input file path: %s\n", input_filepath_rel)

	//Open input file
	file, err := os.Open(input_filepath_rel)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	var part1_answer uint64 = 0

	var part2_answer uint64 = 0
	var part2_is_allowed bool = true

	//Read file line by line
	for scanner.Scan() {
		line := scanner.Text()

		part1_answer += parse_line(line)

		var part2_to_add uint64 = 0
		part2_to_add, part2_is_allowed = parse_line2(line, part2_is_allowed)
		part2_answer += part2_to_add
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	//Print results and exit
	fmt.Printf("Part1 result: %d\n", part1_answer)
	fmt.Printf("Part2 result: %d\n", part2_answer)
}
