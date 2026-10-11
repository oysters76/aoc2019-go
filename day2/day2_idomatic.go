package main 

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	opAdd = 1 
	opMult = 2 
	opHalt = 99
)

const target = 19690720 

type computer struct {
	mem []int 
	ip int 
}

func parse(input string) ([]int, error) {
	fields := strings.Split(strings.TrimSpace(input), ",") 
	prog := make([]int, len(fields)) 
	for i, f := range fields {
		n, err := strconv.Atoi(strings.TrimSpace(f)) 
		if err != nil {
			return nil, fmt.Errorf("error parsing field %d (%q): %w", i, f, err) 
		} 
		prog[i] = n
	}
	return prog, nil
}

func newComputer(prog []int) *computer {
	mem := make([]int, len(prog)) 
	copy(mem, prog) 
	return &computer{
		mem: mem,
	}
}

func (c *computer) inBounds(addr int) bool {
	return addr >= 0 && addr < len(c.mem) 
}

func (c *computer) run() error {
	for {
		halted, err := c.step() 
		if err != nil {
			return err
		}
		if halted {
			return nil 
		}
	}
}

func (c *computer) step() (bool, error) {
	if !c.inBounds(c.ip) {
		return false, fmt.Errorf("instruction pointer %d out of bounds", c.ip)
	}
	switch op := c.mem[c.ip]; op {
	case opHalt:
		return true, nil 
	case opAdd:
		return false, c.apply(func(a,b int) int {return a+b})
	case opMult:
		return false, c.apply(func(a,b int) int {return a*b})
	default:
		return false, fmt.Errorf("unknown opcode %d at position %d", op, c.ip)
	}
}

func (c *computer) apply(op func(int,int) int) error {
	if c.ip+3 >= len(c.mem) {
		return fmt.Errorf("not enough parameters for operation at position %d", c.ip)
	}
	a,b,dst := c.mem[c.ip+1], c.mem[c.ip+2], c.mem[c.ip+3] 
	for _, addr := range []int{a,b,dst} {
		if !c.inBounds(addr) {
			return fmt.Errorf("address %d out of bounds at position %d", addr, c.ip)
		}
	}
	c.mem[dst] = op(c.mem[a], c.mem[b]) 
	c.ip += 4 
	return nil
}

func runWith(prog []int, noun, verb int) (int, error) {
	if len(prog) < 3 {
		return 0, errors.New("program too short to set noun and verb")
	} 
	c := newComputer(prog)
	c.mem[1], c.mem[2] = noun, verb 
	if err := c.run(); err != nil {
		return 0, err
	}
	return c.mem[0], nil
}

func search(prog []int, want int) (int, int, error) {
	for noun := 0; noun <= 99; noun++ {
		for verb := 0; verb <= 99; verb++ {
			got, err := runWith(prog, noun, verb) 
			if err != nil {
				return 0, 0, fmt.Errorf("error running with noun=%d, verb=%d: %w", noun, verb, err)
			}
			if got == want {
				return noun, verb, nil 
			}
		}
	}
	return 0,0,fmt.Errorf("no noun and verb found that produce %d", want)
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err) 
		os.Exit(1) 
	}
}

func run() error {
	filename := "day2.txt" 
	if len(os.Args) > 1 {
		filename = os.Args[1]
	}

	data,err := os.ReadFile(filename) 
	if err != nil {
		return err 
	}
	prog, err := parse(string(data)) 
	if err != nil {
		return err
	}

	part1, err := runWith(prog, 12, 2) 
	if err != nil {
		return fmt.Errorf("error running part 1: %w", err)
	}
	fmt.Println("Part 1:", part1) 

	noun, verb, err := search(prog, target) 
	if err != nil {
		return fmt.Errorf("error running part 2: %w", err)
	}
	fmt.Printf("Part 2: noun=%d, verb=%d, answer=%d\n", noun, verb, 100*noun+verb)
	return nil
}