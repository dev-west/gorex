package gorex

import(
	"errors"
	"fmt"
	"regexp"
	"testing"
)

// AddClass Test
func testClass(c int, t *testing.T) (*testing.T, error) {
	if c & validClasses == 0 {
		t.Fatalf("AddClass(0x%04x & 0x%04x => 0x%04x) attempt to test invalid class\n", c, validClasses, c & validClasses)
		return t, errors.New("AddClass(0x%04x) attempt to test invalid class\n")
	}

	var g Gorex
	var e error

	g, _ = GolangExpression()
	g, e = g.AddClass(c)
	if e != nil {
		t.Fatalf("AddClass(0x%04x) unexpected err %#v => %s\n", c, g, e)
	}
	if len(g.groups) == 0 { t.Fatalf("AddClass(0x%04x) unexpected error, empty groups\n", c) }
	if len(g.groups[0].tokens) == 0 { t.Fatalf("AddClass(0x%04x) unexpected error, empty tokens\n", c) }
	if g.groups[0].tokens[0].class == 0 { t.Fatalf("AddClassToLast(0x%04x) class data does not exist: %#v\n", c, g.groups[0].tokens[0].class) }
	if len(g.groups[0].tokens[0].fixed) != 0 { t.Fatalf("AddClass(0x%04x) unexpected error, fixed data exists: %#v\n", c, g.groups[0].tokens[0].fixed) }
	q := g.groups[0].tokens[0].quantity
	if q.flag != 0 { t.Fatalf("AddClass(0x%04x) unexpected error, quantity flag data exists: %#v\n", c, q) }
	if q.start != 0 { t.Fatalf("AddClass(0x%04x) unexpected error, quantity start data exists: %#v\n", c, q) }
	if q.end != 0 { t.Fatalf("AddClass(0x%04x) unexpected error, quantity end data exists: %#v\n", c, q) }
	w := g.groups[0].tokens[0].class
	if w != c {
		t.Fatalf("AddClass(0x%04x) non-equivalent %#v != %#v\n", c, c, w)
	}

	return t, nil
}

func TestAddClass(t *testing.T) {
	var g Gorex
	var e error

	// construct expression
	g, _ = GolangExpression()

	// test error
	_, e = g.AddClass(NoClass)
	if e == nil {
		t.Fatalf("AddClass(0x%04x &^ 0x%04x => 0x%04x) did not produce expected error: %s\n", NoClass, validClasses, NoClass & validClasses, e)
	}

	// test success
	testClass(Uppers, t)
	testClass(Lowers, t)
	testClass(Numbers, t)
	testClass(Whitespaces, t)
	testClass(Punctuation, t)
	testClass(Uppers + Lowers, t)
	testClass(Uppers + Lowers + Numbers, t)
}

func TestAddClassToLast(t *testing.T) {
	var g Gorex
	var e error

	// construct expression
	g, _ = GolangExpression()

	// test errors
	_, e = g.AddClassToLast(NoClass)
	if e == nil {
		t.Fatalf("AddClassToLast(0x%04x) did not produce expected error (invalid class): %s\n", NoClass, e)
	}

	_, e = g.AddClassToLast(Uppers)
	if e == nil {
		t.Fatalf("AddClassToLast(0x%04x) did not produce expected error (no last): %s\n", Uppers, e)
	}

	// test success
	g, e = g.AddClass(Uppers)
	if e != nil { t.Fatalf("AddClassToLast(0x%04x) unexpected error produced: %s\n", Uppers, e) }
	g, e = g.AddClassToLast(Lowers)
	if e != nil { t.Fatalf("AddClassToLast(0x%04x) unexpected error produced: %s\n", Lowers, e)	}
	c := Uppers + Lowers
	if len(g.groups) == 0 { t.Fatalf("AddClassToLast(0x%04x) unexpected error, empty groups\n", c) }
	if len(g.groups[0].tokens) == 0 { t.Fatalf("AddClassToLast(0x%04x) unexpected error, empty tokens\n", c) }
	if g.groups[0].tokens[0].class == 0 { t.Fatalf("AddClassToLast(0x%04x) class data does not exist: %#v\n", c, g.groups[0].tokens[0].class) }
	if len(g.groups[0].tokens[0].fixed) != 0 { t.Fatalf("AddClassToLast(0x%04x) unexpected error, fixed data exists: %#v\n", c, g.groups[0].tokens[0].fixed) }
	q := g.groups[0].tokens[0].quantity
	if q.flag != 0 { t.Fatalf("AddClass(0x%04x) unexpected error, quantity flag data exists: %#v\n", c, q) }
	if q.start != 0 { t.Fatalf("AddClass(0x%04x) unexpected error, quantity start data exists: %#v\n", c, q) }
	if q.end != 0 { t.Fatalf("AddClass(0x%04x) unexpected error, quantity end data exists: %#v\n", c, q) }
	w := Uppers + Lowers
	if w != g.groups[0].tokens[0].class {
			t.Fatalf("AddClassToLast(0x%04x) tokens not equivalent %#v != %#v\n", c, w, g.groups[0].tokens[0].class)
	}
}

func TestAddFixed(t *testing.T) {
	var g Gorex
	var e error

	// construct expression
	g, _ = GolangExpression()

	// test errors
	g, e = g.AddFixed( "\xFF" )
	if e == nil { t.Fatalf("AddFixed(\\x255) expected error\n") }

	// test success
	c := "com"
	g, e = g.AddFixed(c)
	if e != nil { t.Fatalf("AddFixed(\"%s\") unexpected error: %s\n", c, e) }
	if len(g.groups) == 0 { t.Fatalf("AddFixed(\"%s\") unexpected error, empty groups\n", c) }
	if len(g.groups[0].tokens) == 0 { t.Fatalf("AddFixed(\"%s\") unexpected error, empty tokens\n", c) }
	if g.groups[0].tokens[0].class != 0 { t.Fatalf("AddFixed(0x%04x) class data exists: %#v\n", c, g.groups[0].tokens[0].class) }
	if len(g.groups[0].tokens[0].fixed) == 0 { t.Fatalf("AddFixed(%s) unexpected error, fixed data does not exists: %#v\n", c, g.groups[0].tokens[0].fixed) }
	q := g.groups[0].tokens[0].quantity
	if q.flag != 0 { t.Fatalf("AddClass(0x%04x) unexpected error, quantity flag data exists: %#v\n", c, q) }
	if q.start != 0 { t.Fatalf("AddClass(0x%04x) unexpected error, quantity start data exists: %#v\n", c, q) }
	if q.end != 0 { t.Fatalf("AddClass(0x%04x) unexpected error, quantity end data exists: %#v\n", c, q) }
	if len(c) != len(g.groups[0].tokens[0].fixed) {
		t.Fatalf("AddFixed(%s) unexpected error, fixed data does not match due to lengths %d != %d\n", c, len(c), len(g.groups[0].tokens[0].fixed)) }
	for i, s := range(c) {
		if g.groups[0].tokens[0].fixed[i] != byte(s) {
			t.Fatalf("AddFixed(%s) unexpected error, fixed data does not match: %s != %s\n", c, c, g.groups[0].tokens[0].fixed)
			break
		}
	}
}

func TestAddFixedToLast(t *testing.T) {
	var g Gorex
	var e error
	var q rexQuan
	c := [...]string{ "com", "net" }

	// construct expression
	g, _ = GolangExpression()

	// test errors
	_, e = g.AddFixedToLast("\xFF")
	if e == nil { t.Fatalf("AddFixedToLast(\"\\x%2x\") did not produce expected error (invalid byte): %s\n", '\xFF', e)	}

	_, e = g.AddFixedToLast(c[0])
	if e == nil { t.Fatalf("AddFixedToLast(\"%s\") did not produce expected error (no last): %s\n", c[0], e) }

	// test success
	g, e = g.AddFixed(c[0])
	if e != nil { t.Fatalf("AddFixedToLast(\"%s\") unexpected error produced: %s\n", c[0], e) }
	g, e = g.AddFixedToLast(c[1])
	if e != nil { t.Fatalf("AddFixedToLast(\"%s\") unexpected error produced: %s\n", c[1], e) }
	if len(g.groups) == 0 { t.Fatalf("AddFixedToLast(\"%s\") unexpected error, empty groups\n", c[1]) }
	if len(g.groups[0].tokens) == 0 { t.Fatalf("AddFixedToLast(\"%s\") unexpected error, empty tokens\n", c[1]) }
	if len(g.groups[0].tokens[0].fixed) == 0 { t.Fatalf("AddFixedToLast(\"%s\") class data does not exist: %#v\n", c, g.groups[0].tokens[0].fixed) }
	if g.groups[0].tokens[0].class != 0 { t.Fatalf("AddFixedToLast(\"%s\") unexpected error, class data exists: %#v\n", c[1], g.groups[0].tokens[0].class) }
	if g.groups[0].tokens[1].class != 0 { t.Fatalf("AddFixedToLast(\"%s\") unexpected error, class data exists: %#v\n", c[1], g.groups[0].tokens[1].class) }
	q = g.groups[0].tokens[0].quantity
	if q.flag != 0 { t.Fatalf("AddFixed(0x%04x) unexpected error, quantity flag data exists: %#v\n", c[1], q) }
	if q.start != 0 { t.Fatalf("AddFixed(0x%04x) unexpected error, quantity start data exists: %#v\n", c[1], q) }
	if q.end != 0 { t.Fatalf("AddFixed(0x%04x) unexpected error, quantity end data exists: %#v\n", c[1], q) }
	q = g.groups[0].tokens[1].quantity
	if q.flag != 0 { t.Fatalf("AddFixed(0x%04x) unexpected error, quantity flag data exists: %#v\n", c[1], q) }
	if q.start != 0 { t.Fatalf("AddFixed(0x%04x) unexpected error, quantity start data exists: %#v\n", c[1], q) }
	if q.end != 0 { t.Fatalf("AddFixed(0x%04x) unexpected error, quantity end data exists: %#v\n", c[1], q) }
	for i, tk := range(g.groups[0].tokens) {
		if c[i] != tk.fixed {
			t.Fatalf("AddFixedToLast(\"%s\") tokens[%d] not equivalent %#v != %#v\n", c[1], i, c[i], tk.fixed)
		}
	}
}

func TestApplyQuantityToLast(t *testing.T) {
	var g Gorex
	var e error
	var q int = OneOrMore
	var f string = "com"
	var c int = Uppers

	// construct expression
	g, _ = GolangExpression()

	// test errors
	_, e = g.ApplyQuantityToLast(q)
	if e == nil { t.Fatalf("ApplyQuantityToLast(%d) expected error invalid group index", q) }

	g.groups = append(g.groups, rexGroup{ })
	_, e = g.ApplyQuantityToLast(q)
	if e == nil { t.Fatalf("ApplyQuantityToLast(%d) expected error invalid token index", q) }

	g, _ = GolangExpression()
	g, _ = g.AddFixed(f)
	g, e = g.ApplyQuantityToLast(0)
	if e == nil { t.Fatalf("ApplyQuantityToLast(%d) expected error invalid quantifier", q) }

	g, _ = GolangExpression()
	g, _ = g.AddClass(c)
	g, e = g.ApplyQuantityToLast()
	if e == nil { t.Fatalf("ApplyQuantityToLast() expected error invalid quantifier") }

	g, _ = GolangExpression()
	g, _ = g.AddClass(c)
	g, e = g.ApplyQuantityToLast(q, q, q)
	if e == nil { t.Fatalf("ApplyQuantityToLast() expected error invalid quantifier") }

	g, _ = GolangExpression()
	g, _ = g.AddClass(c)
	g, e = g.ApplyQuantityToLast(2)
	if e == nil { t.Fatalf("ApplyQuantityToLast() expected error invalid quantifier") }

	g, _ = GolangExpression()
	g, _ = g.AddClass(c)
	g, e = g.ApplyQuantityToLast(-2, -3)
	if e == nil { t.Fatalf("ApplyQuantityToLast() expected error invalid quantifier") }

	g, _ = GolangExpression()
	g, _ = g.AddClass(c)
	g, e = g.ApplyQuantityToLast(0, 0)
	if e == nil { t.Fatalf("ApplyQuantityToLast() expected error invalid quantifier") }

	// test success
	g, _ = GolangExpression()
	g, _ = g.AddClass(c)
	g, e = g.ApplyQuantityToLast(q)
	if e != nil { t.Fatalf("ApplyQuantityToLast(%d) unexpected error invalid quantifier", q) }

	g, _ = GolangExpression()
	g, _ = g.AddClass(c)
	g, e = g.ApplyQuantityToLast(2, 3)
	if e != nil { t.Fatalf("ApplyQuantityToLast(%d, %d) unexpected error invalid quantifier", 2, 3) }
}

func ExampleEmail() {
	var g Gorex
	var e error

	validEmails := [...]string{ "joe@mail.org", "john_doe@co.net", "perry.@place.com" }
	invalidEmails := [...]string{ "_tobby@message.org", "goat@mail", "finn@.net" }

	if g, e = GolangExpression(); e != nil { fmt.Println("ExampleEmail failed to create Gorex") }
	if g, e = g.AddClass(Uppers); e != nil { fmt.Println("ExampleEmail failed to add Uppers") }
	if g, e = g.AddClassToLast(Lowers); e != nil { fmt.Println("ExampleEmail failed to add Lowers") }
	if g, e = g.AddClassToLast(Numbers); e != nil { fmt.Println("ExampleEmail failed to add Numbers") }
	if g, e = g.ApplyQuantityToLast(OneOrMore); e != nil { fmt.Println("ExampleEmail failed to apply quantity OneOrMore") }

	if g, e = g.AddFixed("."); e != nil { fmt.Println("ExampleEmail failed to add \".\"") }
	if g, e = g.AddFixedToLast("_"); e != nil { fmt.Println("ExampleEmail failed to add \"_\"") }
	if g, e = g.ApplyQuantityToLast(ZeroOrOne); e != nil { fmt.Println("ExampleEmail failed to apply ZeroOrOne") }

	if g, e = g.AddClass(AlphaNumerics); e != nil { fmt.Println("ExampleEmail failed to add Alphanumerics") }
	if g, e = g.ApplyQuantityToLast(ZoreOrMore); e != nil { fmt.Println("ExampleEmail failed to apply ZeroOrMore") }

	if g, e = g.AddFixed("@"); e != nil { fmt.Println("ExampleEmail failed to add \"@\"") }

	if g, e = g.AddClass(AlphaNumerics); e != nil { fmt.Println("ExampleEmail failed to add Alphanumerics") }
	if g, e = g.ApplyQuantityToLast(ZeroOrOne); e != nil { fmt.Println("ExampleEmail failed to apply OneOrZero") }

	if g, e = g.AddFixed("."); e != nil { fmt.Println("ExampleEmail failed to add \".\"") }

	if g, e = g.AddFixed("com"); e != nil { fmt.Println("ExampleEmail failed to add \"com\"") }
	if g, e = g.AddFixed("net"); e != nil { fmt.Println("ExampleEmail failed to add \"net\"") }
	if g, e = g.AddFixed("org"); e != nil { fmt.Println("ExampleEmail failed to add \"org\"") }

	var exp string
	var r string
	if exp, e = g.Output(); e != nil { fmt.Println("ExampleEmail failed to Output") }
	var rex = regexp.MustCompile(exp)

	fmt.Printf("Expression: %s\n", exp)

	for _, r = range(validEmails) {
		fmt.Printf("Attempt: %s, value: %#v\n", r, rex.MatchString(r))
	}

	for _, r = range(invalidEmails) {
		fmt.Printf("Attempt: %s, value: %#v\n", r, rex.MatchString(r))
	}
}
