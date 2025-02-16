package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
)

const lines_count = 1000 //from `wc -l`

func _cmp_int(a, b int) int {
	return a - b
}

func part1(slice1, slice2 []int) int {
	var total_distance int = 0
	for i := 0; i < lines_count; i++ {
		if slice1[i] > slice2[i] {
			total_distance += slice1[i] - slice2[i]
		} else {
			total_distance += slice2[i] - slice1[i]
		}
	}

	return total_distance
}

func count_occurrences_in_sorted_slice(sorted_slice []int, needle int) int {
	var occurrences int = 0

	for _, num := range sorted_slice {
		if num < needle {
			continue
		}

		//slice is sorted, so we can break if number > needle without fear of skipping possible findings
		if num > needle {
			break
		}

		occurrences += 1
	}

	return occurrences
}

func part2(slice1, slice2 []int) int {
	var total_similarity int = 0
	for i := 0; i < lines_count; i++ {
		total_similarity += slice1[i] * count_occurrences_in_sorted_slice(slice2, slice1[i])
	}

	return total_similarity
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

	//Read input data in 2 arrays(slices)
	slice1 := make([]int, lines_count)
	slice2 := make([]int, lines_count)

	scanner := bufio.NewScanner(file)

	//Read file line by line
	var current_index = 0
	for scanner.Scan() {
		line := scanner.Text()
		str_nums := strings.Split(line, "   ")

		slice1[current_index], err = strconv.Atoi(str_nums[0])
		if err != nil {
			panic(err)
		}

		slice2[current_index], err = strconv.Atoi(str_nums[1])
		if err != nil {
			panic(err)
		}

		current_index += 1
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	//Sort them
	slices.SortStableFunc(slice1, _cmp_int)
	slices.SortStableFunc(slice2, _cmp_int)

	//Print results and exit
	fmt.Printf("Part1 result: %d\n", part1(slice1, slice2))
	fmt.Printf("Part2 result: %d\n", part2(slice1, slice2))
}
