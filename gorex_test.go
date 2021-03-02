package gorex

import(
	"fmt"
	"regexp"
	"testing"
)

// test specific AddClass
func testClass(g *Gorex, c string, t *testing.T) (*testing.T, error) {
	var e error

	// creation of GolangExpression
	g, _ = GolangExpression()
	g, e = g.AddClass(c)
	if e != nil {
		t.Fatalf("AddClass(\"%s\") unexpected err %#v => %s\n", c, g, e)
	}

	// after AddClass, there are no classes generated
	if len(g.groups) == 0 { t.Fatalf("AddClass(0x%04x) unexpected error, empty groups\n", c) }

	// after AddClass, there are not tokens generated
	if len(g.groups[0].tokens) == 0 { t.Fatalf("AddClass(0x%04x) unexpected error, empty tokens\n", c) }

	// after AddClass, token marked as default NoClass
	if g.groups[0].tokens[0].class == NoClass { t.Fatalf("AddClassToLast(0x%04x) class data does not exist: %#v\n", c, g.groups[0].tokens[0].class) }

	// after AddClass, token marked for a fixed string
	if len(g.groups[0].tokens[0].fixed) != 0 { t.Fatalf("AddClass(0x%04x) unexpected error, fixed data exists: %#v\n", c, g.groups[0].tokens[0].fixed) }

	// after AddClass, expected class placed in intended variable
	w := g.groups[0].tokens[0].class
	if w != c {
		t.Fatalf("AddClass(\"%s\") non-equivalent %#v != %#v\n", c, c, w)
	}

	// tests passed
	return t, nil
}

// test several AddClass calls
func TestAddClass(t *testing.T) {
	var g *Gorex
	var e error

	// construct safe expression
	g, _ = GolangExpression()

	// test expected NoClass error
	_, e = g.AddClass(NoClass)
	if e == nil {
		t.Fatalf("AddClass(\"%s\") did not produce expected error: %s\n", NoClass, e)
	}

	// test invalid class
	_, e = g.AddClass("?!")
	if e == nil {
		t.Fatalf("AddClass(\"%s\") did not produce expected error: %s\n", "?!", e)
	}

	// test successes for defined classes
	testClass(&Gorex{ }, Ascii, t)
	testClass(&Gorex{ }, Blank, t)
	testClass(&Gorex{ }, Control, t)
	testClass(&Gorex{ }, Digits, t)
	testClass(&Gorex{ }, Graphical, t)
	testClass(&Gorex{ }, Lowers, t)
	testClass(&Gorex{ }, Printable, t)
	testClass(&Gorex{ }, Punctuation, t)
	testClass(&Gorex{ }, Whitespace, t)
	testClass(&Gorex{ }, Uppers, t)
	testClass(&Gorex{ }, HexDigits, t)
	testClass(&Gorex{ }, AlphaNumerics, t)
	testClass(&Gorex{ }, Alphabetics, t)

	// test unsafe built class
	g, _ = g.AddClass(Uppers + Lowers)
	if e == nil {
		t.Fatalf("AddClass(\"%s\") did not produce expected error: %s\n", Uppers + Lowers, e)
	}

	// test intentionally unsafe built class
	g, e = GolangExpression("Unsafe")
	if e != nil {
		t.Fatalf("GolangExpression(\"%s\") unexpected error: %s\n", "Unsafe", e)
	}
	testClass(g, Uppers + Lowers, t)
}

func TestAddClassToLast(t *testing.T) {
	var g *Gorex
	var e error

	// test options
	c := Uppers
	d := Lowers
	cd := c + d
	w := Alphabetics

	// construct safe expression
	g, _ = GolangExpression()

	// test ...ToLast for invalid class
	_, e = g.AddClassToLast(NoClass)
	if e == nil {
		t.Fatalf("AddClassToLast(0x%04x) did not produce expected error (invalid class): %s\n", NoClass, e)
	}

	// test ...ToLast without prior call
	_, e = g.AddClassToLast(c)
	if e == nil {
		t.Fatalf("AddClassToLast(0x%04x) did not produce expected error (no last): %s\n", c, e)
	}

	// Add Prior
	g, e = g.AddClass(c)
	if e != nil { t.Fatalf("AddClassToLast(0x%04x) unexpected error produced: %s\n", c, e) }

	// Valid ...ToLast
	g, e = g.AddClassToLast(d)
	if e != nil { t.Fatalf("AddClassToLast(0x%04x) unexpected error produced: %s\n", d, e)	}

	// one group added as expected
	if len(g.groups) != 1 { t.Fatalf("AddClassToLast(0x%04x) unexpected error, group count\n", d) }

	// one token added as expected
	if len(g.groups[0].tokens) != 1 { t.Fatalf("AddClassToLast(0x%04x) unexpected error, token count\n", d) }

	// no class token data
	if g.groups[0].tokens[0].class == NoClass { t.Fatalf("AddClass(0x%04x) class data does not exist: %#v\n", c, g.groups[0].tokens[0].class) }

	// fixed data exists when none requested
	if len(g.groups[0].tokens[0].fixed) != 0 { t.Fatalf("AddClassToLast(0x%04x) unexpected error, fixed data exists: %#v\n", c, g.groups[0].tokens[0].fixed) }

	// expected output
	if w != g.groups[0].tokens[0].class {
			t.Fatalf("AddClassToLast(0x%04x) tokens not equivalent %#v != %#v\n", cd, w, g.groups[0].tokens[0].class)
	}
}

func TestAddFixed(t *testing.T) {
	var g *Gorex
	var e error

	// construct expression
	g, _ = GolangExpression()

	// test out of range error
	g, e = g.AddFixed( "\xFF" )
	if e == nil { t.Fatalf("AddFixed(\\x255) expected error\n") }

	// valid request
	c := "com"
	g, e = g.AddFixed(c)

	// valid reqeust return error
	if e != nil { t.Fatalf("AddFixed(\"%s\") unexpected error: %s\n", c, e) }

	// one group added as expected
	if len(g.groups) != 1 { t.Fatalf("AddFixed(\"%s\") unexpected error, group count\n", c) }

	// one token added as expected
	if len(g.groups[0].tokens) != 1 { t.Fatalf("AddFixed(\"%s\") unexpected error, token count\n", c) }

	// class data exists when it should be empty
	if g.groups[0].tokens[0].class != NoClass { t.Fatalf("AddFixed(0x%04x) class data exists: %#v\n", c, g.groups[0].tokens[0].class) }

	// fixed string has data
	if len(g.groups[0].tokens[0].fixed) == 0 { t.Fatalf("AddFixed(%s) unexpected error, fixed data does not exists: %#v\n", c, g.groups[0].tokens[0].fixed) }

	// fixed data is identical
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
	var g *Gorex
	var e error
	c := [...]string{ "com", "net" }

	// construct expression
	g, _ = GolangExpression()

	// test out of range
	_, e = g.AddFixedToLast("\xFF")
	if e == nil { t.Fatalf("AddFixedToLast(\"\\x%2x\") did not produce expected error (invalid byte): %s\n", '\xFF', e)	}

	// ...ToLast without prior
	_, e = g.AddFixedToLast(c[0])
	if e == nil { t.Fatalf("AddFixedToLast(\"%s\") did not produce expected error (no last): %s\n", c[0], e) }

	// valid test
	g, e = g.AddFixed(c[0])
	if e != nil { t.Fatalf("AddFixedToLast(\"%s\") unexpected error produced: %s\n", c[0], e) }
	g, e = g.AddFixedToLast(c[1])
	if e != nil { t.Fatalf("AddFixedToLast(\"%s\") unexpected error produced: %s\n", c[1], e) }

	// one group generated
	if len(g.groups) != 1 { t.Fatalf("AddFixedToLast(\"%s\") unexpected error, group count\n", c[1]) }

	// one token generated
	if len(g.groups[0].tokens) != 2 { t.Fatalf("AddFixedToLast(\"%s\") unexpected error, empty tokens\n", c[1]) }

	// fixed string zero length
	if len(g.groups[0].tokens[0].fixed) == 0 { t.Fatalf("AddFixedToLast(\"%s\") class data does not exist: %#v\n", c, g.groups[0].tokens[0].fixed) }

	// class data found when none expected
	if g.groups[0].tokens[0].class != NoClass { t.Fatalf("AddFixedToLast(\"%s\") unexpected error, class data exists: %#v\n", c[1], g.groups[0].tokens[0].class) }
	if g.groups[0].tokens[1].class != NoClass { t.Fatalf("AddFixedToLast(\"%s\") unexpected error, class data exists: %#v\n", c[1], g.groups[0].tokens[1].class) }

	// fixed data matches intended
	for i, tk := range(g.groups[0].tokens) {
		if c[i] != tk.fixed {
			t.Fatalf("AddFixedToLast(\"%s\") tokens[%d] not equivalent %#v != %#v\n", c[1], i, c[i], tk.fixed)
		}
	}
}

func TestApplyQuantifier(t *testing.T) {
	var g *Gorex
	var e error
	var q Quantifier = OneOrMore
	var f string = "com"
	var c string = Uppers

	// construct expression
	g, _ = GolangExpression()

	// quantifier applied without prior, no group
	_, e = g.ApplyQuantifier(q)
	if e == nil { t.Fatalf("ApplyQuantifier(\"%s\") expected error invalid group index", q) }

	// quantifier applied without prior, no token
	g.groups = append(g.groups, rexGroup{ })
	_, e = g.ApplyQuantifier(q)
	if e == nil { t.Fatalf("ApplyQuantifier(\"%s\") expected error invalid token index", q) }

	// apply empty quantifier
	g, _ = GolangExpression()
	g, _ = g.AddFixed(f)
	g, e = g.ApplyQuantifier("")
	if e == nil { t.Fatalf("ApplyQuantifier(\"\") expected error invalid quantifier") }

	// apply incorrect arguments
	g, _ = GolangExpression()
	g, _ = g.AddClass(c)
	g, e = g.ApplyQuantifier(OneOrMore, 1, 2)
	if e == nil { t.Fatalf("ApplyQuantifier() expected error invalid quantifier") }

	g, _ = GolangExpression()
	g, _ = g.AddClass(c)
	g, e = g.ApplyQuantifier(MinToMax)
	if e == nil { t.Fatalf("ApplyQuantifier() expected error invalid quantifier") }
	g, e = g.ApplyQuantifier(MinToMax, 1)
	if e == nil { t.Fatalf("ApplyQuantifier() expected error invalid quantifier") }

	g, _ = GolangExpression()
	g, _ = g.AddClass(c)
	g, e = g.ApplyQuantifier(q, -2, -3)
	if e == nil { t.Fatalf("ApplyQuantifier() expected error invalid quantifier") }

	g, _ = GolangExpression()
	g, _ = g.AddClass(c)
	g, e = g.ApplyQuantifier(q, 0, 0)
	if e == nil { t.Fatalf("ApplyQuantifier() expected error invalid quantifier") }

	// test valid
	g, _ = GolangExpression()
	g, _ = g.AddClass(c)
	g, e = g.ApplyQuantifier(q)
	if e != nil { t.Fatalf("ApplyQuantifier(\"%s\") unexpected error invalid quantifier", q) }

	q = MinToMax
	g, _ = GolangExpression()
	g, _ = g.AddClass(c)
	g, e = g.ApplyQuantifier(q, 2, 3)
	if e != nil { t.Fatalf("ApplyQuantifier(\"%s\", %d, %d) unexpected error invalid quantifier", q, 2, 3) }
}

func TestFlags(t *testing.T) {
	var g *Gorex
	var e error
	var f string = "COM"
	var s string = "com"

	// construct expression
	g, _ = GolangExpression()

	// flag applied without prior, no group
	_, e = g.SetFlags(CaseInsensitive)
	if e == nil { t.Fatalf("SetFlags(\"%s\") expected error invalid group index", CaseInsensitive) }

	_, e = g.ClearFlags(CaseInsensitive)
	if e == nil { t.Fatalf("ClearFlags(\"%s\") expected error invalid group index", CaseInsensitive) }

	// test CaseInsensitive flag
	g.AddFixed(f) // add uppercase
	_, e = g.SetFlags(CaseInsensitive)
	if e != nil { t.Fatalf("SetFlags(\"%s\") unexpected error", CaseInsensitive) }
	g.AddFixed(s) // add lowercase
	_, e = g.ClearFlags(CaseInsensitive)
	if e != nil { t.Fatalf("ClearFlags(\"%s\") unexpected error", CaseInsensitive) }
	o, e := g.Output()
	r := regexp.MustCompile(o)
	// test against (?i)(COM)(?-i)(com) against "comcom"
	if !r.MatchString(s+s) { t.Fatalf("r.MatchString(\"%s\") failed to match", s) }
	// test against (?i)(COM)(?-i)(com) against "comCOM"
	if r.MatchString(s+f) { t.Fatalf("r.MatchString(\"%s\") unexpectedly matched", s) }

	// test MultiLineMode flag TODO--cannot test without ^$ support

	// test PeriodMatchesNewline
	g, _ = GolangExpression()
	g.AddFixed("com.org") // add string with period
	_, e = g.SetFlags(PeriodMatchesNewline)
	if e != nil { t.Fatalf("SetFlags(\"%s\") unexpected error", PeriodMatchesNewline) }
	o, e = g.Output()
	r = regexp.MustCompile(o)
	// test against (?i)(com.org) against "com\norg"
	if !r.MatchString("com\norg") { t.Fatalf("r.MatchString(\"%s\") failed to match %s", o, "com\\norg") }

	g, _ = GolangExpression()
	g.AddFixed("com.org") // add string with period
	o, e = g.Output()
	r = regexp.MustCompile(o)
	// test against (com.org) against "com\norg"
	if r.MatchString("com\norg") { t.Fatalf("r.MatchString(\"%s\") unexpectedly matched %s", o, "com\\norg") }

	// test Ungreedy
	g, _ = GolangExpression()
	g.AddClass(Digits) // add numerics filter
	g.ApplyQuantifier(OneOrMore) // add + prefer more
	_, e = g.SetFlags(UngreedySwap) // swaps ungreedy so OneOrMore => prefer fewer
	if e != nil { t.Fatalf("SetFlags(\"%s\") unexpected error", UngreedySwap) }
	o, _ = g.Output()
	r = regexp.MustCompile(o)
	// test against (?U)([0-9]+) against "0123456789" prefers fewer--expects '0'
	if len(r.Find([]byte("0123456789"))) > 1 { t.Fatalf("r.MatchString(\"%s\") failed to match %s", o, "0123456789") }

	_, e = g.ClearFlags(UngreedySwap) // swaps ungreedy so OneOrMore => prefer fewer
	if e != nil { t.Fatalf("SetFlags(\"%s\") unexpected error", UngreedySwap) }
	o, _ = g.Output()
	r = regexp.MustCompile(o)
	// test against ([0-9]+) against "0123456789" prefers more--expects ['0' '1' ... '9']
	if len(r.Find([]byte("0123456789"))) != len("0123456789") { t.Fatalf("r.MatchString(\"%s\") failed to match %s", o, "0123456789") }

}

func ExampleEmail() {
	var g *Gorex
	var e error

	validEmails := [...]string{ "joe@mail.org", "john_doe@co.net", "perry.@place.com" }
	invalidEmails := [...]string{ "_tobby@message.org", "goat@mail", "finn@.net" }

	if g, e = GolangExpression(); e != nil { fmt.Println("ExampleEmail failed to create Gorex") }
	if g, e = g.AddClass(Uppers); e != nil { fmt.Println("ExampleEmail failed to add Uppers") }
	if g, e = g.AddClassToLast(Lowers); e != nil { fmt.Println("ExampleEmail failed to add Lowers") }
	if g, e = g.AddClassToLast(Digits); e != nil { fmt.Println("ExampleEmail failed to add Numbers") }
	if g, e = g.ApplyQuantifier(OneOrMore); e != nil { fmt.Println("ExampleEmail failed to apply quantity OneOrMore") }

	if g, e = g.AddFixed("."); e != nil { fmt.Println("ExampleEmail failed to add \".\"") }
	if g, e = g.AddFixedToLast("_"); e != nil { fmt.Println("ExampleEmail failed to add \"_\"") }
	if g, e = g.ApplyQuantifier(ZeroOrOne); e != nil { fmt.Println("ExampleEmail failed to apply ZeroOrOne") }

	if g, e = g.AddClass(AlphaNumerics); e != nil { fmt.Println("ExampleEmail failed to add Alphanumerics") }
	if g, e = g.ApplyQuantifier(ZeroOrOne); e != nil { fmt.Println("ExampleEmail failed to apply OneOrZero") }

	if g, e = g.AddFixed("@"); e != nil { fmt.Println("ExampleEmail failed to add \"@\"") }

	if g, e = g.AddClass(AlphaNumerics); e != nil { fmt.Println("ExampleEmail failed to add Alphanumerics") }
	if g, e = g.ApplyQuantifier(ZeroOrOne); e != nil { fmt.Println("ExampleEmail failed to apply OneOrZero") }

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
