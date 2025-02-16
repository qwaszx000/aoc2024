package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const lines_count = 1000 //from `wc -l`

type LevelGrowDirection uint8

const (
	Undefined_init   LevelGrowDirection = iota
	Undefined_second                    //used as indicator of second iteration of for loop
	Increasing
	Decreasing
	Same
)

func abs(num int) int {
	if num < 0 {
		return -num
	}

	return num
}

func calculate_growth_dir(last_level, current_level int) LevelGrowDirection {
	if current_level > last_level {
		return Increasing
	} else if current_level < last_level {
		return Decreasing
	} else {
		return Same
	}
}

type SafetyReportState struct {
	last_level int
	growth_dir LevelGrowDirection
}

func (state *SafetyReportState) is_level_safe(level_delta int, calculated_growth_dir LevelGrowDirection) bool {
	//growth direction
	if state.growth_dir > Undefined_second && state.growth_dir != calculated_growth_dir {
		return false
	}

	//delta
	if state.growth_dir > Undefined_init {
		level_delta = abs(level_delta)
		if level_delta < 1 || level_delta > 3 {
			return false
		}
	}

	return true
}

func (state *SafetyReportState) update_init_state(level int, calculated_growth_dir LevelGrowDirection) bool {
	//Use LevelGrowDirection's Undefined_init and Undefined_second as current state indicators
	//Undefined_init means it's first iteration and we can't determine growth_dir or difference with last_level(because it has init value yet)
	//Undefined_second means it's second iteration - we can determine growth_dir and difference, but haven't yet
	if state.growth_dir == Undefined_init {
		state.growth_dir = Undefined_second
		state.last_level = level
		return true
	}

	if state.growth_dir == Undefined_second {
		if calculated_growth_dir == Same {
			//they must differ by at least 1
			return false
		} else {
			state.growth_dir = calculated_growth_dir
		}
	}

	return true
}

func (state *SafetyReportState) perform_level_check_iteration(level int) bool {
	level_delta := state.last_level - level
	calculated_growth_dir := calculate_growth_dir(state.last_level, level)

	//Order is important, if we place update_init_state before is_level_safe -> it will break logic
	if state.is_level_safe(level_delta, calculated_growth_dir) == false {
		return false
	}

	if state.update_init_state(level, calculated_growth_dir) == false {
		return false
	}

	state.last_level = level

	return true
}

func is_safe_report_line(report_line string) bool {
	levels := strings.Split(report_line, " ")

	line_state := SafetyReportState{0, Undefined_init}
	for _, level_str := range levels {

		level, _ := strconv.Atoi(level_str)
		if line_state.perform_level_check_iteration(level) == false {
			return false
		}
	}

	return true
}

func is_safe_report_line_with_pd_skip_current(report_line string) bool {
	levels := strings.Split(report_line, " ")

	var is_pb_available bool = true
	line_state := SafetyReportState{0, Undefined_init}
	for _, level_str := range levels {
		level, _ := strconv.Atoi(level_str)

		state_copy := line_state
		//fmt.Printf("%p != %p\n", &line_state, &state_copy) //make sure the addrs are different -> we need to copy values, not ptr

		if line_state.perform_level_check_iteration(level) == false {

			if is_pb_available == true {
				is_pb_available = false
				//restore prev state, because we skip current level, but perform_level_check_iteration can change state
				line_state = state_copy
				continue
			} else {
				return false
			}
		}
	}

	return true
}

func is_safe_report_line_with_pd_skip_prev(report_line string) bool {
	levels := strings.Split(report_line, " ")

	var is_pb_available bool = true
	line_state := SafetyReportState{0, Undefined_init}
	var prev_state [2]SafetyReportState //hack, but we need it

	fmt.Println(report_line)
	for i := 0; i < len(levels); i++ {
		fmt.Println(i, levels[i], line_state)
		level, _ := strconv.Atoi(levels[i])

		if line_state.perform_level_check_iteration(level) == false {

			if is_pb_available == true {
				fmt.Println("PB")
				is_pb_available = false
				//restore prev state, because we skip current level, but perform_level_check_iteration can change state
				line_state = prev_state[1]
				i -= 1 //step back
				continue
			} else {
				fmt.Println("RIP")
				return false
			}
		}

		prev_state[1] = prev_state[0]
		prev_state[0] = line_state
	}
	fmt.Println("-----------")

	return true
}

// Part 2 - with problem dumper
func is_safe_report_line_with_pd(report_line string) bool {
	levels := strings.Split(report_line, " ")

	if is_safe_report_line(report_line) == true {
		return true
	}

	skip_first_str := strings.Join(levels[1:], " ")
	if is_safe_report_line(skip_first_str) == true {
		//fmt.Println("Safe skip first", report_line)
		return true
	}

	if is_safe_report_line_with_pd_skip_current(report_line) == true {
		//fmt.Println("Safe skip one", report_line)
		return true
	}

	if is_safe_report_line_with_pd_skip_prev(report_line) == true {
		fmt.Println("Safe skip one prev", report_line)
		return true
	}

	return false
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

	var part1_answer int = 0
	var part2_answer int = 0

	//Read file line by line
	var current_index = 0
	for scanner.Scan() {
		line := scanner.Text()

		//part1
		if is_safe_report_line(line) {
			//fmt.Printf("Safe line: %s\n", line)
			part1_answer += 1
		}

		//part2
		if is_safe_report_line_with_pd(line) {
			//fmt.Printf("Safe2 line: %s\n", line)
			part2_answer += 1
		}

		current_index += 1
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	//Print results and exit
	fmt.Printf("Part1 result: %d\n", part1_answer)
	fmt.Printf("Part2 result: %d\n", part2_answer)
}
