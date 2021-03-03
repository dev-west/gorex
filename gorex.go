// gorex package MIT license
// github.com/aggero/gorex
// builds basic regular expression from understandable blocks
//
// usage
// 1. create a new gorex constructor
// 2. define content filter
// 3. define modifiers for that content filter
// 4. repeat steps 2-3
// 5. output regex string at any point
//
//  var rex = gorex.GolangExpression()
//  rex.AddClass(gorex.Uppers)                // add group ([A-Z])
//  rex.AddClassToLast(gorex.Lowers)          // modify group ([A-Za-z])
//  rex.AddClassToLast(gorex.Digits)          // modify group ([A-Za-z0-9])
//  rex.ApplyQuantifier(gorex.OneOrMore)      // modify group ([A-Za-z0-9]+)
//  rex.AddFixed(".")                         // add group (.)
//  rex.AddFixedToLast("_")                   // modify group (.|_)
//  rex.ApplyQuantifier(gorex.ZeroOrOne)      // modify group (.|_?)
//  rex.AddClass(gorex.Alphanumerics)         // add group ([A-Za-z0-9]+)
//  rex.ApplyQuantifier(gorex.ZeroOrOne)
//  rex.AddFixed("@")                         // add group (@)
//  rex.AddClass(gorex.Alphanumerics)         // add group ([A-Za-z0-9]+)
//  rex.ApplyQuantifier(gorex.OneOrMore)
//  rex.AddFixed(".")                         // add group (.)
//  rex.AddFixed("com")                       // add group (com)
//  rex.AddFixedToLast("net")                 // modify gorup (com|net)
//  rex.AddFixedToLast("org")                 // modify gorup (com|net|org)
//  var validEmail = regexp.MustCompile(rex.Output())
//  fmt.Println(validEmail.MatchString("adam@gmail.com"))
//
// -- Creates regular expressions in the form: (g1)(g2)(g3)...
//    wherein each group is implemented as a range (class) OR a
//    fixed value ("com")
// -- The example above is expected to create the following exp
//    ([A-Za-z0-9]+)(.?)(_?)([A-Za-z0-9]+)(@)([A-Za-z0-9+])(.)(com|net|org)

package gorex

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
)

type Gorex struct {
	groups []rexGroup // expression details
	unsafe bool
}

type rexGroup struct {
	tokens []rexToken
	flags rexFlag
	anchor Anchor
}

type rexToken struct {
	fixed string
	class string
	quantifier rexQuan
}

// token modifiers
type Anchor string

const (
	atBeginning = "^"
	atEnd = "&"
)

// token quantifiers (variable data)
type rexQuan struct {
	regexp Quantifier
	argv [2]int
}

// quantifier definitions
type Quantifier string

const (
	Single Quantifier = ""
	ZeroOrMore Quantifier = "*"
	OneOrMore Quantifier = "+"
	ZeroOrOne Quantifier = "?"
	MinToMax Quantifier = "{%d,%d}"
	MinOrMore Quantifier = "{%d,}"
	Exactly Quantifier = "{%d}"
	ZeroOrMorePrefFewer Quantifier = "*?"
	OneOrMorePrefFewer Quantifier = "+?"
	ZeroOrOnePrefFewer Quantifier = "??"
	MinToMaxPrefFewer Quantifier = "{%d,%d}?"
	MinOrMorePrefFewer Quantifier = "{%d,}?"
	ExactlyPrefFewer Quantifier = "{%d}?"
)

// class definitions
const (
	NoClass string = ""
	Ascii string = "\x00-\x7F"
	Blank string = "\t "
	Control string = "\x00-\x1F\x7F"
	Digits string = "0-9"
	Graphical string = "!-~"
	Lowers string = "a-z"
	Printable string = " -~"
	Punctuation string = "!-/:-@[-`{-~"
	Whitespace string = "\t\n\v\f\r "
	Uppers string = "A-Z"
	Words string = "0-9A-Za-z_"
	HexDigits string = "0-9A-Fa-f"
	AlphaNumerics string = "0-9A-Za-z"
	Alphabetics string = "A-Za-z"
)

// flag definitions
type Flag string

type rexFlag struct {
	i bool
	m bool
	s bool
	U bool
}

const (
	CaseInsensitive string = "i"
	MultiLineMode string = "m"
	PeriodMatchesNewline string = "s"
	UngreedySwap string = "U"
)

// options
const (
	Unsafe string = "Unsafe"
)

func GolangExpression(opts ...string) (*Gorex, error) {
	var r = &Gorex{ }
	r.unsafe = false
	if len(opts) > 1 {
		return &Gorex{ }, errors.New("Gorex @116: invalid GolangExpression options")
	} else {
		for _, op := range(opts) {
			if op == "Unsafe" {
				r.unsafe = true
				continue
			}

			return &Gorex{ }, errors.New("Gorex @124: invalid GolangExpression options")
		}
	}

	return r, nil
}

func (g *Gorex) Output() (string, error) {
	o := bytes.NewBufferString("")
	activeFlags := rexFlag{ false, false, false, false }
	for _, gr := range(g.groups) {
		flagParens := false
		if (gr.flags.i || gr.flags.m || gr.flags.s || gr.flags.U) ||
				((!gr.flags.i || !gr.flags.m || !gr.flags.s || !gr.flags.U) &&
				(activeFlags.i || activeFlags.m || activeFlags.s || activeFlags.U)) {
			flagParens = true
			o.WriteString("(?")
		}
		if gr.flags.i {
			activeFlags.i = true
			o.WriteString(CaseInsensitive)
		}
		if gr.flags.m {
			activeFlags.m = true
			o.WriteString(MultiLineMode)
		}
		if gr.flags.s {
			activeFlags.s = true
			o.WriteString(PeriodMatchesNewline)
		}
		if gr.flags.U {
			activeFlags.U = true
			o.WriteString(UngreedySwap)
		}
		if (!gr.flags.i && activeFlags.i) || (!gr.flags.m && activeFlags.m) || (!gr.flags.s && activeFlags.s) || (!gr.flags.U && activeFlags.U) {
			o.WriteString("-")
			if !gr.flags.i && activeFlags.i {
				activeFlags.i = false
				o.WriteString(CaseInsensitive)
			}
			if !gr.flags.m && activeFlags.m {
				activeFlags.m = false
				o.WriteString(MultiLineMode)
			}
			if !gr.flags.s && activeFlags.s {
				activeFlags.s = false
				o.WriteString(PeriodMatchesNewline)
			}
			if !gr.flags.U && activeFlags.U {
				activeFlags.U = false
				o.WriteString(UngreedySwap)
			}
		}
		if flagParens { o.WriteString(")") }

		o.WriteString("(")
		// add token data
		if gr.anchor != "" { o.WriteString(string(gr.anchor)) }
		for i, tk := range(gr.tokens) {
			if tk.class != NoClass {
				if len(tk.fixed) != 0 {
					return "", errors.New("Gorex @140: invalid token error")
				}
				o.WriteString("[")
				o.WriteString(tk.class)
				o.WriteString("]")
			}

			if len(tk.fixed) != 0 {
				if tk.class != NoClass {
					return "", errors.New("Gorex @149: invalid token error")
				}
				o.WriteString(tk.fixed)
				if(len(gr.tokens) > i + 1) { o.WriteString("|") }
			}

			// add class quantity
			argCount := regexp.MustCompile("%d") // only permits numbers
			var argc int
			if argCount.FindAllString(string(tk.quantifier.regexp), -1) == nil {
				argc = 0
			} else {
				argc = len(argCount.FindAllString(string(tk.quantifier.regexp), -1))
			}
			switch(argc) {
			case 0:
				o.WriteString(fmt.Sprintf(string(tk.quantifier.regexp)))
			case 1:
				o.WriteString(fmt.Sprintf(string(tk.quantifier.regexp), tk.quantifier.argv[0]))
			case 2:
				o.WriteString(fmt.Sprintf(string(tk.quantifier.regexp), tk.quantifier.argv[0], tk.quantifier.argv[1]))
			default:
				return "", errors.New("Gorex @171: invalid argument count")
			}
		}

		o.WriteString(")")
	}

	return o.String(), nil
}

func verifyClass(a string) bool {
	if		a == NoClass ||
			a == Ascii ||
			a == Blank ||
			a == Control ||
			a == Digits ||
			a == Graphical ||
			a == Lowers ||
			a == Printable ||
			a == Punctuation ||
			a == Whitespace ||
			a == Uppers ||
			a == Words ||
			a == HexDigits ||
			a == AlphaNumerics ||
			a == Alphabetics {
		return true
	}

	return false
}

func (g *Gorex) AddClass(c string) error {
	if c == NoClass {
		return errors.New("Gorex @205: invalid class")
	}

	if !g.unsafe { // is safe
		if !verifyClass(c) { return errors.New("Gorex @209: invalid class") }
	}

	id := len(g.groups)
	r := rexGroup{ }
	g.groups = append(g.groups, r)

	g.groups[id].tokens = append(g.groups[id].tokens, rexToken { "", c, rexQuan{ } } )

	return nil
}

func (g *Gorex) AddClassToLast(c string) error {
	if !g.unsafe { // is safe
		if !verifyClass(c) { return errors.New("Gorex @223: invalid class") }
	}

	if len(g.groups) == 0 { return errors.New("Gorex @226: invalid group index") }
	gId := len(g.groups) - 1
	if len(g.groups[gId].tokens) == 0 { return errors.New("Gorex @228: invalid token index") }
	tId := len(g.groups[gId].tokens) - 1

	g.groups[gId].tokens[tId].class = g.groups[gId].tokens[tId].class + c

	return nil
}

func (g *Gorex) AddFixed(a string) error {
	for _, b := range(a) {
		if byte(b) >= 128 { return errors.New("Gorex @238: invalid byte error") }
	}
	id := len(g.groups)
	r := rexGroup{ }
	g.groups = append(g.groups, r)
	g.groups[id].tokens = append(g.groups[id].tokens, rexToken{ a, NoClass, rexQuan{ } } )

	return nil
}

func (g *Gorex) AddFixedToLast(a string) error {
	for _, b := range(a) {
		if byte(b) >= 128 { return errors.New("Gorex @250: invalid byte value") }
	}
	if(len(g.groups) == 0) { return errors.New("Gorex @252: invalid group index") }
	id := len(g.groups) - 1
	g.groups[id].tokens = append(g.groups[id].tokens, rexToken{ a, NoClass, rexQuan{ } } )

	return nil
}

func verifyQuantifier(q Quantifier, args []int) bool {
	argc := len(args)
	if argc > 2 { return false }
	switch(argc) {
	case 0:
		if		q == Single ||
				q == ZeroOrMore ||
				q == OneOrMore ||
				q == ZeroOrOne ||
				q == ZeroOrMorePrefFewer ||
				q == OneOrMorePrefFewer ||
				q == ZeroOrOnePrefFewer {
			return true }
	case 1:
		if		q == MinOrMore ||
				q == Exactly ||
				q == MinOrMorePrefFewer ||
				q == ExactlyPrefFewer {
			return true }
	case 2:
		if		q == MinToMax ||
				q == MinToMaxPrefFewer {
			return true }
	}

	return false
}

func (g *Gorex) ApplyQuantifier(q Quantifier, args ...int) error {
	if string(q) == "" { return errors.New("Gorex @288: invalid quantifier") }
	if len(g.groups) == 0 { return errors.New("Gorex @289: invalid group index") }
	gId := len(g.groups) - 1
	if len(g.groups[gId].tokens) == 0 { return errors.New("Gorex @291: invalid token index") }
	tId := len(g.groups[gId].tokens) - 1
	if g.groups[gId].tokens[tId].class == NoClass && len(g.groups[gId].tokens[tId].fixed) == 0 { return errors.New("Gorex @293: invalid quantifier") }
	if len(args) > 2 { return errors.New("Gorex @294: invalid quantifier") }
	if !verifyQuantifier(q, args) { return errors.New("Gorex @295: invalid quantifier") }

	g.groups[gId].tokens[tId].quantifier.regexp = q

	switch(len(args)) {
	case 0:
		g.groups[gId].tokens[tId].quantifier.argv[0] = 0
		g.groups[gId].tokens[tId].quantifier.argv[1] = 0
	case 1:
		g.groups[gId].tokens[tId].quantifier.argv[0] = args[0]
		g.groups[gId].tokens[tId].quantifier.argv[1] = 0
	case 2:
		g.groups[gId].tokens[tId].quantifier.argv[0] = args[0]
		g.groups[gId].tokens[tId].quantifier.argv[1] = args[1]
	}

	return nil
}

func verifyAnchor(m Anchor) bool {
	if m == atBeginning { return true }
	if m == atEnd { return true }

	return false
}

func (g *Gorex) ApplyAnchor(m Anchor) error {
	if string(m) == "" { return errors.New("Gorex @389: invalid anchor") }
	if len(g.groups) == 0 { return errors.New("Gorex @390: invalid group index") }
	gId := len(g.groups) - 1
	if !verifyAnchor(m) { return errors.New("Gorex @392: invalid anchor") }

	g.groups[gId].anchor = m

	return nil
}

func verifyFlags(c string) bool {
	for _, ch := range(c) {
		if	string(ch) != CaseInsensitive &&
			string(ch) != MultiLineMode &&
			string(ch) != PeriodMatchesNewline &&
			string(ch) != UngreedySwap { return false }
	}

	return true
}

func (g *Gorex) SetFlags(c string) error {
	if c == "" { return errors.New("Gorex @326: invalid flag") }
	if len(g.groups) == 0 { return errors.New("Gorex @327: invalid group index") }
	gId := len(g.groups) - 1

	if !verifyFlags(c) { return errors.New("Gorex @342: invalid flag") }
	for _, ch := range(c) {
		if string(ch) == CaseInsensitive { g.groups[gId].flags.i = true }
		if string(ch) == MultiLineMode { g.groups[gId].flags.m = true }
		if string(ch) == PeriodMatchesNewline { g.groups[gId].flags.s = true }
		if string(ch) == UngreedySwap { g.groups[gId].flags.U = true }
	}

	return nil
}

func (g *Gorex) ClearFlags(c string) error {
	if c == "" { return errors.New("Gorex @326: invalid flag") }
	if len(g.groups) == 0 { return errors.New("Gorex @327: invalid group index") }
	gId := len(g.groups) - 1

	if !verifyFlags(c) { return errors.New("Gorex @342: invalid flag") }

	for _, ch := range(c) {
		if string(ch) == CaseInsensitive { g.groups[gId].flags.i = false }
		if string(ch) == MultiLineMode { g.groups[gId].flags.m = false }
		if string(ch) == PeriodMatchesNewline { g.groups[gId].flags.s = false }
		if string(ch) == UngreedySwap { g.groups[gId].flags.U = false }
	}

	return nil
}
