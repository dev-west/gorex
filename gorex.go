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
//  rex.AddClassToLast(gorex.Numbers)         // modify group ([A-Za-z0-9])
//  rex.ApplyQuantityToLast(gorex.OneOrMore)  // modify group ([A-Za-z0-9]+)
//  rex.AddFixed(".")                         // add group (.)
//  rex.AddFixedToLast("_")                   // modify group (.|_)
//  rex.ApplyQuantityToLast(gorex.ZeroOrOne)  // modify group (.|_?)
//  rex.AddClass(gorex.Alphanumerics)         // add group ([A-Za-z0-9]+)
//  rex.ApplyQuantityToLast(gorex.ZeroOrOne)
//  rex.AddFixed("@")                         // add group (@)
//  rex.AddClass(gorex.Alphanumerics)         // add group ([A-Za-z0-9]+)
//  rex.ApplyQuantityToLast(gorex.OneOrMore)
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
	"strconv"
)

type Gorex struct {
	groups []rexGroup // expression details
}

type rexGroup struct {
	tokens []rexToken
}

type rexToken struct {
	fixed string
	class int
	quantity rexQuan
}

type rexQuan struct {
	flag int
	start int
	end int
}

// token quantifiers
const (
	Single int = 0
	ZeroOrOne int = -1
	ZeroOrMore int = -2
	OneOrMore int = -3
)

// token classes
const (
	NoClass int =     1 >> 1
	Uppers int =      1 << 0   // [A-Z]
	Lowers int =      1 << 1   // [a-z]
	Numbers int =     1 << 2   // [0-9]
	Whitespaces int = 1 << 3   // [ \t\n\r]
	Punctuation int = 1 << 4   // [][!\x22#$%&\x27()*+,./:;<=>?@_\x60{|}~-]
	AlphaNumerics int = Uppers + Lowers + Numbers
)

var eG Gorex = Gorex{ }

const validClasses int = Uppers + Lowers + Numbers + Whitespaces + Punctuation

const (
	empty string = ""
	num string = "0-9"
	lower string = "a-z"
	upper string = "A-Z"
	whitespace string = "[ \t]"
	punctuation string = "[][!\x22#$%&\x27()*+,./:;<=>?@_\x60{|}~-]"
)

func GolangExpression() (Gorex, error) {
	var r Gorex

	return r, nil
}

func (g Gorex) Output() (string, error) {
	o := bytes.NewBufferString("")
	for _, gr := range(g.groups) {
		o.WriteString("(")

		// add token data
		for i, tk := range(gr.tokens) {
			if tk.class != 0 {
				if len(tk.fixed) != 0 {
					return "", errors.New("Gorex @112: invalid token error")
				}
				o.WriteString("[")
				if tk.class & Uppers != 0 { o.WriteString(upper) }
				if tk.class & Lowers != 0 {	o.WriteString(lower) }
				if tk.class & Numbers != 0 { o.WriteString(num) }
				if tk.class & Whitespaces != 0 { o.WriteString(whitespace) }
				if tk.class & Punctuation != 0 { o.WriteString(punctuation) }
				o.WriteString("]")
			}

			if len(tk.fixed) != 0 {
				if tk.class != 0 {
					return "", errors.New("Gorex @125: invalid token error")
				}
				o.WriteString(tk.fixed)
				if(len(gr.tokens) > i + 1) { o.WriteString("|") }
			}

			// add class quantity
			switch tk.quantity.flag {
			case Single: break
			case ZeroOrOne:	o.WriteString("?")
			case ZeroOrMore: o.WriteString("*")
			case OneOrMore: o.WriteString("+")
			default: return "", errors.New("Gorex @137: invalid quantity error")
			}

			if tk.quantity.end != 0 {
				if tk.quantity.start > tk.quantity.end {
					return "", errors.New("Gorex @142: invalid quantity error")
				}
				o.WriteString("{")
				o.WriteString(strconv.Itoa(tk.quantity.start))
				o.WriteString(",")
				o.WriteString(strconv.Itoa(tk.quantity.end))
				o.WriteString("}")
			}
		}

		o.WriteString(")")
	}

	return o.String(), nil
}

func (g Gorex) AddClass(c int) (Gorex, error) {
	if c & validClasses == 0 { return eG, errors.New("Gorex @159: invalid class value") }
	id := len(g.groups)
	r := rexGroup{ }
	g.groups = append(g.groups, r)
	g.groups[id].tokens = append(g.groups[id].tokens, rexToken { "", c, rexQuan{ } } )

	return g, nil
}

func (g Gorex) AddClassToLast(c int) (Gorex, error) {
	if c & validClasses == 0 { return eG, errors.New("Gorex @169: invalid class value") }
	if(len(g.groups) == 0) { return eG, errors.New("Gorex @170: invalid group index") }
	gId := len(g.groups) - 1
	if len(g.groups[gId].tokens) == 0 { return eG, errors.New("Gorex @172: invalid token index") }
	tId := len(g.groups[gId].tokens) - 1
	g.groups[gId].tokens[tId].class = g.groups[gId].tokens[tId].class + c

	return g, nil
}

func (g Gorex) AddFixed(a string) (Gorex, error) {
	for _, b := range(a) {
		if byte(b) >= 128 { return eG, errors.New("Gorex @181: invalid byte error") }
	}
	id := len(g.groups)
	r := rexGroup{ }
	g.groups = append(g.groups, r)
	g.groups[id].tokens = append(g.groups[id].tokens, rexToken{ a, NoClass, rexQuan{ } } )

	return g, nil
}

func (g Gorex) AddFixedToLast(a string) (Gorex, error) {
	for _, b := range(a) {
		if byte(b) >= 128 { return eG, errors.New("Gorex @193: invalid byte value") }
	}
	if(len(g.groups) == 0) { return eG, errors.New("Gorex @195: invalid group index") }
	id := len(g.groups) - 1
	g.groups[id].tokens = append(g.groups[id].tokens, rexToken { a, NoClass, rexQuan{ } } )

	return g, nil
}

func (g Gorex) ApplyQuantityToLast(params ...int) (Gorex, error) {
	if len(g.groups) == 0 { return eG, errors.New("Gorex @203: invalid group index") }
	gId := len(g.groups) - 1
	if len(g.groups[gId].tokens) == 0 { return eG, errors.New("Gorex @205: invalid token index") }
	tId := len(g.groups[gId].tokens) - 1
	if g.groups[gId].tokens[tId].class == 0 && len(g.groups[gId].tokens[tId].fixed) == 0 { return eG, errors.New("Gorex @238: invalid quantifier") }
	if len(params) == 0 || len(params) > 2 { return eG, errors.New("Gorex @208: invalid quantifier") }
	switch(len(params)) {
	case 1:
		if params[0] >= 0 { return eG, errors.New("Gorex @211: invalid quantifier") }
		g.groups[gId].tokens[tId].quantity.flag = params[0]
		g.groups[gId].tokens[tId].quantity.start = 0
		g.groups[gId].tokens[tId].quantity.end = 0
	case 2:
		if params[0] < 0 || params[1] <= 0 { return eG, errors.New("Gorex @216: invalid quantifier") }
		g.groups[gId].tokens[tId].quantity.flag = 0
		g.groups[gId].tokens[tId].quantity.start = params[0]
		g.groups[gId].tokens[tId].quantity.end = params[1]
	}

	return g, nil
}
