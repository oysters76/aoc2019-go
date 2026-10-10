# What You Learned: Go Fundamentals Through Advent of Code Day 1

A study guide built from refactoring your fuel calculator from a working first draft into idiomatic, production-style Go.

---

## 1. Functions and Function Signatures

In Go, the **type comes after the name**, and the return type comes after the parameter list.

```go
func fuelForMass(mass int) int {
	return mass/3 - 2
}
```

Read it as: *a function named `fuelForMass`, taking `mass` of type `int`, returning an `int`.*

Functions can return **multiple values**. The most common pattern is a result plus an error:

```go
func readInts(r io.Reader) ([]int, error)
```

---

## 2. Integer Division Truncates Toward Zero

Dividing two `int`s throws away the fractional part. It does **not** round down for negative numbers.

```go
7 / 2   // 3
-7 / 2  // -3  (floor would be -4)
7 % 2   // 1
-7 % 2  // -1  (remainder takes the sign of the dividend)
```

- Need decimals? Convert first: `float64(7) / 2` → `3.5`
- Integer division by zero **panics**; float division by zero gives `Inf` or `NaN`.

Your fuel formula `mass/3 - 2` relies on this truncation. That's exactly what the puzzle asks for.

---

## 3. Functions Are Values

Go has no "function pointers," but functions are first-class values. The type of a function from `int` to `int` is written `func(int) int`.

```go
func sumBy(nums []int, fn func(int) int) int {
	sum := 0
	for _, n := range nums {
		sum += fn(n)
	}
	return sum
}

sumBy(masses, fuelForMass)                         // pass a named function
sumBy(masses, func(n int) int { return n * n })    // pass an anonymous one
```

**Key point:** `fuelForMass` passes the function; `fuelForMass(5)` calls it and passes the result.

For long or repeated function types, give them a name:

```go
type IntFunc func(int) int
```

---

## 4. Go Has Only One Loop: `for`

There is no `while` or `do-while`. `for` covers every case.

```go
// Classic three-part loop
for i := 0; i < 3; i++ { }

// While loop: condition only
for n < 100 { }

// Infinite loop: exit with break
for { if done { break } }

// Range over a slice
for _, n := range nums { }
```

### How the three-part loop runs

```
setup → check → body → step → check → body → step → check → ... stop
```

The condition is checked **before every pass, including right after the step**. This is the same in C, C++, Java, JavaScript, and C#.

That's why this loop never adds a negative value:

```go
for fuel := fuelForMass(mass); fuel > 0; fuel = fuelForMass(fuel) {
	total += fuel
}
```

| Step  | Action                 | `fuel` | `total` |
|-------|------------------------|--------|---------|
| Setup | `14/3 - 2`             | 2      | 0       |
| Check | `2 > 0` ✓              | 2      | 0       |
| Body  | `total += 2`           | 2      | 2       |
| Step  | `2/3 - 2`              | -2     | 2       |
| Check | `-2 > 0` ✗ → **stop**  | -2     | 2       |

The negative value is computed but rejected before the body runs again. This is why the refactor could drop your extra `if current_mass > 0` check.

Any three-part loop can be rewritten as a while loop with the same behavior:

```go
fuel := fuelForMass(mass)      // setup
for fuel > 0 {                 // check
	total += fuel              // body
	fuel = fuelForMass(fuel)   // step
}
```

---

## 5. Reading a File Line by Line

```go
scanner := bufio.NewScanner(r)
for scanner.Scan() {
	line := scanner.Text()   // current line, newline stripped
	// ...
}
if err := scanner.Err(); err != nil {
	// handle read error
}
```

- `scanner.Scan()` returns `true` while there's another line.
- `scanner.Text()` gives the line as a string.
- `strconv.Atoi(s)` converts a string to an `int`, returning an error for bad input.
- `strings.TrimSpace` removes surrounding whitespace.
- Default max line length is 64 KB. Raise it with `scanner.Buffer(...)` if needed.

### Why `scanner.Err()` is required

`Scan()` returns `false` **both** at end of file **and** on a read failure. The loop stops either way. Checking `scanner.Err()` afterward is the only way to tell a normal finish from a failure.

---

## 6. `if` With an Init Statement

```go
if err := scanner.Err(); err != nil {
	log.Fatal(err)
}
```

The part before the semicolon runs first; the part after is the condition. It's equivalent to:

```go
err := scanner.Err()
if err != nil {
	log.Fatal(err)
}
```

**Benefit:** `err` only exists inside the `if` block, so it doesn't leak into the rest of the function.

---

## 7. `defer`

```go
f, err := os.Open(path)
if err != nil {
	return err
}
defer f.Close()
```

`defer` schedules a call to run **when the surrounding function returns**, no matter how it returns (normal end, early `return`, or panic).

- Cleanup sits right next to the setup, so it's hard to forget.
- Multiple defers run in **reverse order** (last in, first out).
- Arguments are evaluated **immediately**, at the `defer` line.
- Place `defer` **after** the error check; if `Open` failed, there's nothing to close.
- `log.Fatal` and `os.Exit` **skip** deferred calls.

---

## 8. Type Assertions

A type assertion extracts the concrete value from an interface.

```go
var x interface{} = 42

n := x.(int)            // panics if x isn't an int
s, ok := x.(string)     // safe: ok is false, no panic

switch v := x.(type) {  // type switch for several possibilities
case int:
	fmt.Println("int", v)
case string:
	fmt.Println("string", v)
}
```

Go has **no built-in `assert`** for tests. You check conditions yourself and call `t.Errorf`.

---

## 9. Refactoring to Idiomatic Go

This was the big lesson. Here's what changed between your first draft and the refactor, and why each change matters.

### 9.1 Run `gofmt`

Go has one official format: tabs for indentation, sorted imports. Every company expects it. Let your editor run it on save. Never format by hand.

```go
// Before: imports unsorted, inconsistent indentation
import (
"fmt"
"bufio"
...
```

### 9.2 camelCase, never snake_case

| Before                       | After              |
|------------------------------|--------------------|
| `calc_fuel_for_mass`         | `fuelForMass`      |
| `calc_fuel_for_mass_ver2`    | `totalFuelForMass` |
| `run_op_on_function`         | `sumBy`            |
| `current_mass`               | `fuel`             |

Capitalization has meaning: `lowercase` = private to the package, `Uppercase` = exported.

### 9.3 Names describe purpose, not mechanics

`ver2` tells a reader nothing about what's different. `run_op_on_function` describes *how*, not *what*. Good names (`totalFuelForMass`, `sumBy`) explain intent. Drop filler words like `calc` and `get`.

### 9.4 One job per function

Your original `run_op_on_function` opened the file, parsed it, **and** did math, so the file was read **twice**.

| Function    | Single responsibility       |
|-------------|-----------------------------|
| `readInts`  | Parse input into `[]int`    |
| `sumBy`     | Apply a function and sum    |
| `run`       | Wire the pieces together    |

Now the file is read once, and each piece is easy to understand and test.

### 9.5 Accept `io.Reader`, not a filename

```go
func readInts(r io.Reader) ([]int, error)
```

This works with files, network streams, or plain strings. It makes testing trivial:

```go
readInts(strings.NewReader("12\n14\n"))
```

**Rule of thumb:** accept interfaces, return concrete types.

### 9.6 Wrap errors with context

```go
return nil, fmt.Errorf("line %d: %w", lineNum, err)
```

Turns a vague `invalid syntax` into:

```
day1.txt: line 37: strconv.Atoi: parsing "abc": invalid syntax
```

`%w` keeps the original error inside, so callers can still inspect it with `errors.Is` / `errors.As`.

### 9.7 Keep `main` tiny with a `run` function

```go
func main() {
	if err := run(path); err != nil {
		log.Fatal(err)
	}
}
```

`main` can't return an error, so the real work goes in `run() error`. This gives one exit point and ensures deferred cleanup (like `f.Close()`) actually runs before `log.Fatal`.

### 9.8 Consistent error names

```go
f, e := run_op_on_function(...)   // before: confusing
masses, err := readInts(f)        // after: always "err"
```

### 9.9 Doc comments

A comment directly above a function, starting with its name, is Go's documentation convention. `go doc` reads these.

```go
// fuelForMass returns the fuel needed to launch a module of the given mass.
func fuelForMass(mass int) int {
```

### 9.10 Accept input from the command line

```go
path := "day1.txt"
if len(os.Args) > 1 {
	path = os.Args[1]
}
```

Hardcoded paths make programs rigid. A sensible default plus an override is a small, professional touch.

---

## 10. Table-Driven Tests

The standard Go testing pattern. Put tests in `main_test.go` next to `main.go`.

```go
func TestFuelForMass(t *testing.T) {
	tests := []struct {
		mass, want int
	}{
		{12, 2},
		{14, 2},
		{1969, 654},
		{100756, 33583},
	}
	for _, tt := range tests {
		if got := fuelForMass(tt.mass); got != tt.want {
			t.Errorf("fuelForMass(%d) = %d, want %d", tt.mass, got, tt.want)
		}
	}
}
```

- Adding a case is one line.
- Failure messages show input, actual, and expected values.
- `t.Errorf` reports and continues; `t.Fatal` reports and stops the test.
- Test bad input too, not just the happy path.

---

## 11. The Go Philosophy

- **Simplicity over cleverness.** Usually one obvious way to do something.
- **Readability over writability.** Code is read far more than it's written.
- **Big standard library.** HTTP, JSON, crypto, testing, and more ship with Go.
- **"A little copying is better than a little dependency."** Write the three-line loop instead of importing a package.
- **No `map` / `filter` / `reduce` in the standard library**, on purpose. A plain loop is considered clearer. (Third-party `github.com/samber/lo` exists if a team wants it.)

### Why dependencies hurt less than in Node

| Go                                              | Node                                   |
|-------------------------------------------------|----------------------------------------|
| One shared module cache                         | `node_modules` per project             |
| Minimal version selection: picks the oldest compatible version | Often resolves to newest in range |
| `go.sum` checksums + public checksum database   | Lockfile only                          |
| No install scripts                              | `postinstall` can run arbitrary code   |
| Major versions get new import paths (`/v2`)     | Major bumps possible via ranges        |
| Go 1 compatibility promise since 2012           | Frequent ecosystem churn               |
| Single static binary to deploy                  | Runtime + `node_modules` on server     |

**The honest tradeoff:** more repetitive code, lots of `if err != nil`, and a smaller ecosystem for some areas.

---

## 12. Commands to Know

```bash
go mod init day1     # create a module (once per project)
gofmt -w .           # format all files
go vet ./...         # catch common mistakes
go test ./...        # run tests
go run . day1.txt    # run with an argument
go build             # produce a single binary
```

Many teams also use `golangci-lint`, which bundles dozens of checks.

---

## Code Review Checklist

Before calling Go code "done," ask:

- [ ] Has it been run through `gofmt`?
- [ ] Are all names camelCase and descriptive?
- [ ] Does each function do one thing?
- [ ] Do functions accept interfaces (`io.Reader`) where possible?
- [ ] Are errors returned (not swallowed) and wrapped with context?
- [ ] Is every `scanner.Err()` checked?
- [ ] Is every opened resource closed with `defer`?
- [ ] Is `main` tiny, delegating to `run() error`?
- [ ] Are there table-driven tests, including bad input?
- [ ] Does `go vet` pass cleanly?

---

## Practice Exercises

1. **Floor division.** Write `floorDiv(a, b int) int` that rounds toward negative infinity. Test it with `(-7, 2)` → `-4`.
2. **Recursion.** Rewrite `totalFuelForMass` recursively. Which version do you find clearer?
3. **Defer order.** Predict the output, then run it:
   ```go
   for i := 0; i < 3; i++ {
   	defer fmt.Println(i)
   }
   ```
4. **Error inspection.** Use `errors.As` to check whether an error from `readInts` contains a `*strconv.NumError`.
5. **Test the reader.** Add test cases to `readInts` for an empty input and for input with only blank lines.
6. **Next puzzle.** Solve Advent of Code 2019 Day 2 applying every item from the checklist above.

---

## Further Reading

- **Effective Go**: the official guide to writing idiomatic Go
- **Go Code Review Comments**: the style checklist most companies follow
- **Go Proverbs**: short sayings that capture the philosophy
- **A Tour of Go**: interactive basics, good for review