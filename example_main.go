package main

import(
	"fmt"
	. "gorex"
	"regexp"
)

func main() {
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
    if g, e = g.ApplyQuantityToLast(ZeroOrOne); e != nil { fmt.Println("ExampleEmail failed to apply OneOrZero") }

    if g, e = g.AddFixed("@"); e != nil { fmt.Println("ExampleEmail failed to add \"@\"") }

    if g, e = g.AddClass(AlphaNumerics); e != nil { fmt.Println("ExampleEmail failed to add Alphanumerics") }
    if g, e = g.ApplyQuantityToLast(OneOrMore); e != nil { fmt.Println("ExampleEmail failed to apply OneOrZero") }

    if g, e = g.AddFixed("."); e != nil { fmt.Println("ExampleEmail failed to add \".\"") }

    if g, e = g.AddFixed("com"); e != nil { fmt.Println("ExampleEmail failed to add \"com\"") }
    if g, e = g.AddFixedToLast("net"); e != nil { fmt.Println("ExampleEmail failed to add \"net\"") }
    if g, e = g.AddFixedToLast("org"); e != nil { fmt.Println("ExampleEmail failed to add \"org\"") }

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
