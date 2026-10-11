package main

import (
	"fmt"
	"strconv"
	"strings"
	"os"
)

const (
	HALT = 99 
	ADD  = 1 
	MULT = 2
)

const (
	TARGET = 19690720
	P1_LIMIT = 99
	P2_LIMIT = 99
)

type ProgramState struct {
	Values []int; 
	InstructionPointer int; 
}

func toInts(strs []string) ([]int, error) {
	nums := make([]int, len(strs)) 
	for i, s := range strs {
		n, err := strconv.Atoi(s) 
		if err != nil {
			return nil, fmt.Errorf("element %d (%q) : %w", i, s, err) 
		}
		nums[i] = n 
	}
	return nums, nil 
}

func CreateProgramState(input string) (*ProgramState, error) {
	input = strings.TrimSpace(input)
	parts := strings.Split(input, ",") 
	nums, err := toInts(parts) 
	if err != nil {
		return nil, err
	} 
	return &ProgramState{
		Values: nums,
		InstructionPointer: 0,
	}, nil
}

func GetCurrentInstruction(ps *ProgramState) (int, error) {
	if ps.InstructionPointer >= len(ps.Values) {
		return 0, fmt.Errorf("instruction pointer %d out of bounds", ps.InstructionPointer) 
	}
	return ps.Values[ps.InstructionPointer], nil 
}

func GetParamsForOp(ps *ProgramState, numParams int) ([]int, error) {
	if ps.InstructionPointer+numParams >= len(ps.Values) {
		return nil, fmt.Errorf("not enough parameters for operation at position %d", ps.InstructionPointer) 
	}
	params := make([]int, numParams) 
	for i := 0; i < numParams; i++ {
		params[i] = ps.Values[ps.InstructionPointer+i+1] 
	}
	return params, nil 
}

func HandleAddOp(ps *ProgramState) (*ProgramState, error, bool) {
	params,err := GetParamsForOp(ps, 3)
	if err != nil {
		return nil, err, false
	}
	ps.Values[params[2]] = ps.Values[params[0]] + ps.Values[params[1]] 
	ps.InstructionPointer += 4 
	return ps, nil, false
}

func HandleMultOp(ps *ProgramState) (*ProgramState, error, bool) {
	params,err := GetParamsForOp(ps, 3)
	if err != nil {
		return nil, err, false
	}
	ps.Values[params[2]] = ps.Values[params[0]] * ps.Values[params[1]] 
	ps.InstructionPointer += 4 
	return ps, nil, false
}

func HandleOp(ps *ProgramState) (*ProgramState, error, bool) {
	instr, err := GetCurrentInstruction(ps); 
	if err != nil {
		return nil, err, true
	}
	switch instr {
	case HALT: 
		return ps, nil, true
	case ADD:
		return HandleAddOp(ps) 
	case MULT:
		return HandleMultOp(ps) 
	default:
		return nil, fmt.Errorf("unknown instruction %d at position %d", instr, ps.InstructionPointer), true
	}	
}

func ReadFileToStr(filename string) (string, error) {
	data, err := os.ReadFile(filename) 
	if err != nil {
		fmt.Println("error reading file:", err) 
		return "", err
	}
	content := string(data) 
	return content, nil
}

func RunProgram(filename string, p1 int, p2 int) (*ProgramState, error) {
	content, err := ReadFileToStr(filename)
	if err != nil {
		return nil, err
	}
	ps,err := CreateProgramState(content) 
	if err != nil {
		return nil, err
	}
	ps.Values[1] = p1
	ps.Values[2] = p2
	ps, err, done := HandleOp(ps)
	if err != nil {
		return nil, err
	}
	for !done {
		ps, err, done = HandleOp(ps)
		if err != nil {
			return nil, err
		}
	}
	return ps, nil	
}

func BruteForceSearch(filename string, target int, p1Lim int, p2Lim int) (int, int, error) {
	for p1 := 0; p1 <= p1Lim; p1++ {
		for p2 := 0; p2 <= p2Lim; p2++ {
			fmt.Printf("Trying p1=%d, p2=%d\n", p1, p2)
			ps, err := RunProgram(filename, p1, p2)
			if err != nil {
				return 0, 0, err
			}
			if ps.Values[0] == target {
				return p1, p2, nil
			}
		}
	}
	return 0, 0, fmt.Errorf("target %d not found", target)
}

func main() {
	filename := "day2.txt"
	if len(os.Args) > 1 {
		filename = os.Args[1] 
	}
	ps, err := RunProgram(filename, 12, 2) 
	if err != nil {
		fmt.Println("error running program:", err) 
		return
	}
	fmt.Println("Final state:", ps.Values)

	p1, p2, err := BruteForceSearch(filename, TARGET, P1_LIMIT, P2_LIMIT)
	if err != nil {
		fmt.Println("error during brute force search:", err) 
		return
	}
	fmt.Printf("Found p1=%d, p2=%d, puzzle answer: %d\n", p1, p2, 100*p1+p2)
}