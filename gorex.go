// gorex package MIT license
// github.com/aggero/gorex
// builds basic regular expression from understandable blocks
//
// usage
// 1. create a new gorex constructor with optional standard definition
// 2. define possibility blocks (p-blocks)
// 3. define connecting blocks (c-blocks) between those p-blocks
// 4. optimize & output regex string at any point
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
}

type rexToken struct {
	fixed string
	class string
	quantifier rexQuan
}

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

var eG *Gorex = &Gorex{ }

// options
const (
	Unsafe string = "Unsafe"
)

func GolangExpression(opts ...string) (*Gorex, error) {
	var r Gorex
	r.unsafe = false
	if len(opts) > 1 {
		return eG, errors.New("Gorex @116: invalid GolangExpression options")
	} else {
		for _, op := range(opts) {
			if op == "Unsafe" {
				r.unsafe = true
				continue
			}

			return eG, errors.New("Gorex @124: invalid GolangExpression options")
		}
	}

	return &r, nil
}

func (g *Gorex) Output() (string, error) {
	o := bytes.NewBufferString("")
	for _, gr := range(g.groups) {
		o.WriteString("(")

		// add token data
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

func (g *Gorex) AddClass(c string) (*Gorex, error) {
	if c == NoClass {
		return eG, errors.New("Gorex @205: invalid class")
	}

	if !g.unsafe { // is safe
		if !verifyClass(c) { return eG, errors.New("Gorex @209: invalid class") }
	}

	id := len(g.groups)
	r := rexGroup{ }
	g.groups = append(g.groups, r)

	g.groups[id].tokens = append(g.groups[id].tokens, rexToken { "", c, rexQuan{ } } )

	return g, nil
}

func (g *Gorex) AddClassToLast(c string) (*Gorex, error) {
	if !g.unsafe { // is safe
		if !verifyClass(c) { return eG, errors.New("Gorex @223: invalid class") }
	}

	if len(g.groups) == 0 { return eG, errors.New("Gorex @226: invalid group index") }
	gId := len(g.groups) - 1
	if len(g.groups[gId].tokens) == 0 { return eG, errors.New("Gorex @228: invalid token index") }
	tId := len(g.groups[gId].tokens) - 1

	g.groups[gId].tokens[tId].class = g.groups[gId].tokens[tId].class + c

	return g, nil
}

func (g *Gorex) AddFixed(a string) (*Gorex, error) {
	for _, b := range(a) {
		if byte(b) >= 128 { return eG, errors.New("Gorex @238: invalid byte error") }
	}
	id := len(g.groups)
	r := rexGroup{ }
	g.groups = append(g.groups, r)
	g.groups[id].tokens = append(g.groups[id].tokens, rexToken{ a, NoClass, rexQuan{ } } )

	return g, nil
}

func (g *Gorex) AddFixedToLast(a string) (*Gorex, error) {
	for _, b := range(a) {
		if byte(b) >= 128 { return eG, errors.New("Gorex @250: invalid byte value") }
	}
	if(len(g.groups) == 0) { return eG, errors.New("Gorex @252: invalid group index") }
	id := len(g.groups) - 1
	g.groups[id].tokens = append(g.groups[id].tokens, rexToken{ a, NoClass, rexQuan{ } } )

	return g, nil
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

func (g *Gorex) ApplyQuantifier(q Quantifier, args ...int) (*Gorex, error) {
	if string(q) == "" { return eG, errors.New("Gorex @288: invalid quantifier") }
	if len(g.groups) == 0 { return eG, errors.New("Gorex @289: invalid group index") }
	gId := len(g.groups) - 1
	if len(g.groups[gId].tokens) == 0 { return eG, errors.New("Gorex @291: invalid token index") }
	tId := len(g.groups[gId].tokens) - 1
	if g.groups[gId].tokens[tId].class == NoClass && len(g.groups[gId].tokens[tId].fixed) == 0 { return eG, errors.New("Gorex @293: invalid quantifier") }
	if len(args) > 2 { return eG, errors.New("Gorex @294: invalid quantifier") }
	if !verifyQuantifier(q, args) { return eG, errors.New("Gorex @295: invalid quantifier") }

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

	return g, nil
}
