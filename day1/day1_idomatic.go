package main 

import (
	"bufio" 
	"fmt"
	"io" 
	"log"
	"os"
	"strconv" 
	"strings"
)

func fuelForMass(mass int) int {
	return mass/3 - 2 
}

func totalFuelForMass(mass int) int {
	total := 0 
	for fuel := fuelForMass(mass); fuel > 0; fuel = fuelForMass(fuel) {
		total += fuel 
	}
	return total 
}

func readInts(r io.Reader) ([]int, error) {
	var nums []int
	scanner := bufio.NewScanner(r) 
	for lineNum := 1; scanner.Scan(); lineNum++ {
		line := strings.TrimSpace(scanner.Text()) 
		if line == "" {
			continue
		}
		n, err := strconv.Atoi(line) 
		if err != nil {
			return nil, fmt.Errorf("line %d: %w", lineNum, err) 
		}
		nums = append(nums, n)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("reading input: %w", err) 
	}
	return nums, nil 
}

func sumBy(nums []int, fn func(int) int) int {
	sum := 0 
	for _, n := range nums {
		sum += fn(n) 
	}
	return sum 
}

func run(path string) error {
	f, err := os.Open(path) 
	if err != nil {
		return err 
	}
	defer f.Close() 

	masses, err := readInts(f) 
	if err != nil {
		return fmt.Errorf("%s: %w", path, err) 
	}

	fmt.Println("Part 1:", sumBy(masses, fuelForMass)) 
	fmt.Println("Part 2:", sumBy(masses, totalFuelForMass)) 
	return nil 
}

func main() {
	path := "day1.txt" 
	if len(os.Args) > 1 {
		path = os.Args[1] 
	}
	if err := run(path); err != nil {
		log.Fatal(err)
	}
}