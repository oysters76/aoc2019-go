package main

import (
	"fmt"
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

func calc_fuel_for_mass(mass int) int {
	return (mass/3)-2
}

func calc_fuel_for_mass_ver2(mass int) int {
	current_mass := mass 
	s := 0 
	for current_mass > 0 {
		current_mass = calc_fuel_for_mass(current_mass) 
		if current_mass > 0 {
			s += current_mass 
		}
	}
	return s 
}

func run_op_on_function(filename string, fn func(int) int) (int, error) {
	f, err := os.Open(filename) 
	if err != nil {
		return 0, err 
	}
	defer f.Close() 

	sum := 0 
	scanner := bufio.NewScanner(f) 
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text()) 
		if line == "" {
			continue 
		}
		n, err := strconv.Atoi(line) 
		if err != nil {
			return 0, err
		}
		sum += fn(n)
	}
	return sum, scanner.Err() 
}

func main() {
	fuel, err := run_op_on_function("day1.txt", calc_fuel_for_mass) 
	if err != nil {
		log.Fatal(err) 
	}
	fmt.Println("Fuel ", fuel)  

	f,e := run_op_on_function("day1.txt", calc_fuel_for_mass_ver2) 
	if e != nil {
		log.Fatal(e) 
	}
	fmt.Println("Fuel 2 ", f) 
}
