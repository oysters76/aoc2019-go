# Go Learnings

Notes from a session working through the basics of Go while solving an Intcode puzzle (Advent of Code 2019, Day 2).

---

## 1. Strings

### Splitting

```go
parts := strings.Split("a,b,c", ",")      // ["a" "b" "c"]
parts := strings.SplitN("a,b,c,d", ",", 2) // ["a" "b,c,d"]
```

Edge cases to remember:

- `strings.Split("", ",")` returns `[""]`, a slice with one empty string, not an empty slice.
- Consecutive separators produce empty strings: `"a,,b"` gives `["a" "" "b"]`.
- To drop empty entries, use `strings.FieldsFunc(s, func(r rune) bool { return r == ',' })`.
- For real CSV with quoted fields, use `encoding/csv`.

### Trimming and stripping newlines

| Goal | Function |
|---|---|
| Strip leading/trailing whitespace (incl. `\n`, `\r`, tabs) | `strings.TrimSpace(s)` |
| Strip any trailing `\r` / `\n` characters | `strings.TrimRight(s, "\r\n")` |
| Remove exactly one trailing `\n` | `strings.TrimSuffix(s, "\n")` |
| Remove every `\n` anywhere | `strings.ReplaceAll(s, "\n", "")` |

`TrimRight` and `TrimLeft` treat their second argument as a **set of characters**, while `TrimSuffix` and `TrimPrefix` treat it as an **exact string**.

Watch out for Windows line endings (`\r\n`): a leftover `\r` is invisible when printed but breaks comparisons and `strconv.Atoi`.

Strings are **immutable**, so these functions return a new string. Always assign the result: `s = strings.TrimSpace(s)`.

### Converting strings to ints

```go
n, err := strconv.Atoi("42")
if err != nil {
	// handle bad input
}
```

Use `strconv.ParseInt(s, 10, 64)` when you need a specific bit size or a different base.

Converting a whole slice:

```go
func toInts(strs []string) ([]int, error) {
	nums := make([]int, len(strs)) // pre-allocate
	for i, s := range strs {
		n, err := strconv.Atoi(strings.TrimSpace(s))
		if err != nil {
			return nil, fmt.Errorf("element %d (%q): %w", i, s, err)
		}
		nums[i] = n
	}
	return nums, nil
}
```

---

## 2. Structs (Go's answer to classes)

Go has no classes. A **struct** holds data, and **methods** attached to it provide behavior.

```go
type IntList struct {
	values []int
}
```

- `[]int` is a **slice** (growable). `[5]int` is a fixed-size **array**. You almost always want slices.
- **Capitalization controls visibility**: `Values` is exported (public outside the package), `values` is private to the package. This applies to fields, functions, types and constants.
- You can also define a named type directly on a slice: `type IntList []int`.

### Methods and receivers

A method is a function with a **receiver** between `func` and the name:

```go
func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}
```

The receiver is like `this` or `self`, but you name it yourself, usually with the type's first letter.

| Receiver | Behavior | Use when |
|---|---|---|
| Value `(c Counter)` | Gets a copy; changes are lost | Read-only methods on small structs |
| Pointer `(c *Counter)` | Gets a reference; can modify the original | Method mutates state, or struct is large |

Rules of thumb:

- If any method needs a pointer receiver, make them all pointer receivers for consistency.
- You call methods the same way either way (`c.Increment()`); Go takes the address automatically.
- You can only define methods on types declared in your own package. To add methods to `int`, declare `type MyInt int`.

### Constructors

Go has no built-in constructors. The convention is a plain function named `NewTypeName`:

```go
func NewUser(name string, age int) (*User, error) {
	if name == "" {
		return nil, errors.New("name cannot be empty")
	}
	return &User{Name: name, Age: age}, nil
}
```

- Usually return a pointer.
- Validate inputs and return an `error` as the second value.
- Use constructors to initialize maps and slices; writing to a nil map **panics**.
- If the package is named after the type (e.g. `user`), call it `New` so callers write `user.New(...)`.

---

## 3. Pointers

Access fields on a pointer with the normal dot syntax. Go dereferences automatically; there's no `->`.

```go
u := &User{Name: "Alice"} // u is *User
u.Name = "Bob"            // same as (*u).Name = "Bob"
```

- `&x` gets a pointer to `x`.
- `*User` in a type means "pointer to User".
- `*u` gets the value back out (useful for making a copy).
- Pass a pointer to a function when the function should modify the caller's struct.
- Accessing a field on a **nil pointer panics**. Check `if u != nil` when a pointer might be nil.

---

## 4. Constants

```go
const Pi = 3.14159

const (
	StatusOK       = 200
	StatusNotFound = 404
)
```

### Enums with `iota`

`iota` starts at 0 and increments for each line in a `const` block:

```go
type Color int

const (
	Red   Color = iota // 0
	Green              // 1
	Blue               // 2
)
```

It works in expressions too:

```go
const (
	_  = iota             // skip 0
	KB = 1 << (10 * iota) // 1024
	MB                    // 1048576
)
```

Limitations:

- Constants must be known at **compile time** (no `const t = time.Now()`).
- Only numbers, strings, booleans and runes. No constant slices, maps or structs; use a package-level `var` instead.

---

## 5. Switch

```go
switch day {
case "Mon":
	fmt.Println("Start")
case "Tue", "Wed", "Thu": // multiple values
	fmt.Println("Midweek")
default:
	fmt.Println("Other")
}
```

- **No fall-through by default**, so no `break` is needed. Use the `fallthrough` keyword if you really want it.
- A switch with no value acts as a clean if/else chain: `switch { case score >= 90: ... }`.
- You can declare a scoped variable: `switch op := mem[ip]; op { ... }`.
- **Type switch** checks the concrete type of an interface value:

```go
switch x := v.(type) {
case int:
	fmt.Println("int", x)
case string:
	fmt.Println("string", x)
default:
	fmt.Printf("other: %T\n", x)
}
```

Implementing `String() string` on a type (often with a switch) controls how it prints with `fmt`.

---

## 6. Reading files

### Whole file at once

```go
data, err := os.ReadFile("input.txt")
if err != nil {
	return err
}
content := strings.TrimSpace(string(data))
```

`os.ReadFile` opens and closes the file for you. **No `Close` needed.**

Paths are relative to the directory you **run** the program from, not where the `.go` file lives.

### Line by line (for large files)

```go
file, err := os.Open("input.txt")
if err != nil {
	return err
}
defer file.Close()

scanner := bufio.NewScanner(file)
for scanner.Scan() {
	line := scanner.Text()
	// ...
}
if err := scanner.Err(); err != nil {
	return err
}
```

### Closing files and `defer`

- `os.Open` / `os.Create` give you a handle that **you must close**.
- Put `defer file.Close()` immediately **after** the error check (if `Open` fails, `file` is nil).
- `defer` runs when the **function** returns, not at the end of a loop iteration. Opening files in a loop with `defer` keeps them all open; move the per-file work into its own function.
- When **writing**, check the error from `Close`, because buffered data is flushed then: `return f.Close()`.

---

## 7. Optional parameters

Go has no optional or default parameters. Three patterns replace them:

**Variadic parameter** for a single optional value:

```go
func greet(name string, greeting ...string) {
	g := "Hello"
	if len(greeting) > 0 {
		g = greeting[0]
	}
	fmt.Println(g + ", " + name)
}
```

**Options struct** for several optional values. Zero values mean "not set", so apply defaults for them. Use pointer fields if zero is a legitimate value.

```go
func NewServer(cfg ServerConfig) *Server {
	if cfg.Port == 0 {
		cfg.Port = 8080
	}
	// ...
}
```

**Functional options** for flexible, extensible APIs (common in libraries):

```go
type Option func(*Server)

func WithPort(p int) Option { return func(s *Server) { s.port = p } }

func NewServer(opts ...Option) *Server {
	s := &Server{port: 8080} // defaults
	for _, opt := range opts {
		opt(s)
	}
	return s
}

NewServer(WithPort(9000))
```

---

## 8. Idiomatic Go principles (from the code review)

These came out of rewriting the Intcode solution.

### `error` is always the last return value

```go
// Not idiomatic
func HandleOp(ps *ProgramState) (*ProgramState, error, bool)

// Idiomatic
func (c *computer) step() (bool, error)
```

Linters flag anything else, and every Go reader expects it.

### Use methods instead of functions that take the struct

```go
// Before
GetCurrentInstruction(ps)

// After
c.step()
```

### Don't return a pointer you were given

If a method modifies the struct through a pointer receiver, the caller already sees the change. Returning the same pointer is redundant.

### Do expensive work once

The original brute-force search re-read and re-parsed the input file on all 10,000 attempts. Parse once, then give each run a fresh copy of memory:

```go
mem := make([]int, len(prog))
copy(mem, prog)
```

Slices share their underlying array, so without `copy` every run would mutate the same data.

### Remove duplication with functions as values

Add and multiply were identical except for the operator. Pass the operation instead:

```go
case opAdd:
	return false, c.apply(func(a, b int) int { return a + b })
case opMul:
	return false, c.apply(func(a, b int) int { return a * b })
```

### Validate everything you index with

Checking that parameters exist isn't enough; also check that the **addresses they point to** are in bounds. Return an error rather than letting the program panic.

### Naming conventions

- Use `camelCase` or `PascalCase`, never `ALL_CAPS`, even for constants: `opHalt`, not `HALT`.
- Uppercase means exported. In `package main` nothing is imported by anyone else, so lowercase names are the norm.
- Keep names short where scope is small: `c` for a receiver, `ip` for instruction pointer, `mem` for memory.
- No semicolons. Run `gofmt` (or let your editor do it) and it will clean up formatting for you.

### Handle each error exactly once

Either log it **or** return it, never both. Printing an error and then returning it causes it to be reported twice up the chain.

### Wrap errors with context

```go
return fmt.Errorf("part 1: %w", err)
```

`%w` preserves the original error (so `errors.Is` / `errors.As` still work) while adding context about where it happened.

### The `main` → `run() error` pattern

```go
func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func run() error {
	// use plain `return err` everywhere
}
```

This keeps error handling in one place, prints to stderr, and exits with a non-zero status code so scripts can detect failure.

### Keep output meaningful

Remove debug prints from hot loops (thousands of "Trying..." lines) and print answers rather than whole data dumps.

### Use the tooling

- `gofmt` formats code to the standard style.
- `go vet` catches suspicious constructs and common mistakes.