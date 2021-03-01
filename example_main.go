package main

import(
	"fmt"
	. "gorex"
	"regexp"
)

func main() {
    var g *Gorex
    var e error

    validEmails := [...]string{ "joe@mail.org", "john_doe@co.net", "perry.@place.com" }
    invalidEmails := [...]string{ "_tobby@message.org", "goat@mail", "finn@.net" }

    if g, e = GolangExpression(); e != nil { fmt.Println("ExampleEmail failed to create Gorex: %s", e) }

    if g, e = g.AddClass(Uppers); e != nil { fmt.Println("ExampleEmail failed to add Uppers: %s", e) }
    if g, e = g.AddClassToLast(Lowers); e != nil { fmt.Println("ExampleEmail failed to add Lowers: %s", e) }
    if g, e = g.AddClassToLast(Digits); e != nil { fmt.Println("ExampleEmail failed to add Numbers %s", e) }
    if g, e = g.ApplyQuantifier(OneOrMore); e != nil { fmt.Println("ExampleEmail failed to apply quantity OneOrMore %s", e) }

    if g, e = g.AddFixed("."); e != nil { fmt.Println("ExampleEmail failed to add \".\": %s", e) }
    if g, e = g.AddFixedToLast("_"); e != nil { fmt.Println("ExampleEmail failed to add \"_\": %s", e) }
    if g, e = g.ApplyQuantifier(ZeroOrOne); e != nil { fmt.Println("ExampleEmail failed to apply ZeroOrOne: %s", e) }

    if g, e = g.AddClass(AlphaNumerics); e != nil { fmt.Println("ExampleEmail failed to add Alphanumerics: %s", e) }
    if g, e = g.ApplyQuantifier(ZeroOrMore); e != nil { fmt.Println("ExampleEmail failed to apply OneOrZero: %s", e) }

    if g, e = g.AddFixed("@"); e != nil { fmt.Println("ExampleEmail failed to add \"@\": %s", e) }

    if g, e = g.AddClass(AlphaNumerics); e != nil { fmt.Println("ExampleEmail failed to add Alphanumerics: %s", e) }
    if g, e = g.ApplyQuantifier(OneOrMore); e != nil { fmt.Println("ExampleEmail failed to apply OneOrZero: %s", e) }

    if g, e = g.AddFixed("."); e != nil { fmt.Println("ExampleEmail failed to add \".\": %s", e) }

    if g, e = g.AddFixed("com"); e != nil { fmt.Println("ExampleEmail failed to add \"com\": %s", e) }
    if g, e = g.AddFixedToLast("net"); e != nil { fmt.Println("ExampleEmail failed to add \"net\": %s", e) }
    if g, e = g.AddFixedToLast("org"); e != nil { fmt.Println("ExampleEmail failed to add \"org\": %s", e) }

    var exp string
    var r string
    if exp, e = g.Output(); e != nil { fmt.Println("ExampleEmail failed to Output: %s", e) }

    var rex = regexp.MustCompile(exp)

    fmt.Printf("\nExpression: %s\n", exp)

    for _, r = range(validEmails) {
        fmt.Printf("Attempt: %s, value: %#v\n", r, rex.MatchString(r))
    }

    for _, r = range(invalidEmails) {
        fmt.Printf("Attempt: %s, value: %#v\n", r, rex.MatchString(r))
    }

	g, _ = GolangExpression()
	g.AddFixed("n")
    if g, e = g.ApplyQuantifier(MinToMax, 2, 3); e != nil { fmt.Println("ExampleEmail failed to apply MinToMax, %d, %d: %s", 2, 3, e) }

    if exp, e = g.Output(); e != nil { fmt.Println("ExampleEmail failed to Output: %s", e) }

	rex = regexp.MustCompile(exp)

    fmt.Printf("\nExpression: %s\n", exp)

    for _, r = range(validEmails) {
        fmt.Printf("Attempt: %s, value: %#v\n", r, rex.MatchString(r))
    }

    for _, r = range(invalidEmails) {
        fmt.Printf("Attempt: %s, value: %#v\n", r, rex.MatchString(r))
    }
}
