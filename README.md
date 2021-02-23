# gorex
Regular expression syntax builder to introduce basic concepts.

## usage
Regular expressions are an efficient means to search and match a pattern against a collection of text.

These expressions can be obtuse enough to discourage their usage. This package provides a simple collection of functions such that an expression can be generated by an easily understood sequence of function calls. The example found in `example_main.go` provides a simple e-mail regular expression.

This package creates a sequence of groups with each add command. These groups are represented in regular expressions in parentheses. In each group, gorex supports either a class of individual characters or a collection of fixed strings. Classes of characters are things like any upper-case letter (Uppers A through Z) or any numeral (Numerics 0 through 9). Fixed strings are things commonly found in a fixed sequence like the top-level domain of an e-mail address ('.com', '.net', etc).

Classes are added to the expression in square brackets expecting the use of a quantifier. Square brackets in these expressions list a collection of matching options **for a single character**. So if you add multiple classes to a particular group, any character that matches any class added will be accepted. These are commonly added using the `OneOrMore` quantifier, meaning as long as there's one character, there can be as many characters in a row, provided they match the classes in that group.

Fixed strings are added to the expression without square brackets, but may also support quantifiers.

## functionality
GolangExpression() (gorex, error) produces a gorex object and its related functions the following functions are all accessible via any gorex object created

gorex.AddClass(int) (gorex, error) produces a gorex object with a new class group of characters in the expression sequence. the classes consist of one or more of:
      NoClass, Uppers, Lowers, Numbers, Whitespaces, Punctuation, AlphaNumerics
  each is a simple integer flag, so if you need upper- and lower-case characters but not numbers, consider using `Uppers + Lowers`.
    produces a group like: `([A-Za-z])`

gorex.AddClassToLast(int) (gorex, error) produces a gorex object with a class group added to the previously created class sequence. So, instead of adding a class using `Uppers + Lowers`, consider `AddClass(Uppers); AddClassToLast(Lowers)`.
    produces a group like: `([A-Za-z])`

gorex.AddFixed(string) (gorex, error) produces a gorex object with a new fixed group of one or more strings. This expression only accepts one string.
    produces a group like: `(com)`

gorex.AddFixedToLast(string) (gorex, error) produces a gorex object with a new fixed string applied to the prior group using an OR operator.
    produces a group like: `(com|net)`

gorex.ApplyQuantityToLast(...int) (gorex, error) produces a gorex object with a quantifier applied to the last class or fixed token generated. Quantifiers can be any one of:
      Single (no quantifier), ZeroOrOne, ZeroOrMore, OneOrMore
  so, if AddClass(Uppers) is used, only a single character from A through Z matches the expression, but if ApplyQuantityToLast(OneOrMore) is used, any number of any characters that all fall within A through Z match as a group.
    produces a group like: `([A-Z]+)`

hopefuly you'll find that these function names are reasonably straight-forware, if they are, to some extent, verbose.

## example
The example_main.go application provides a simple e-mail verification regular expression generation. Note--the verbose errors are not necessary; they are in the example go code, but not shown here:
```
package main

import(
  "fmt"
  . "gorex" // gorex is added for ease of referencing. '.' causes access to gorex.go exports to be immediately accessible (otherwise, must use 'gorex.' in front of everything)
  "regexp"
)

func main() {
    var g Gorex

    validEmails := [...]string{ "joe@mail.org", "john_doe@co.net", "perry.@place.com" }
    invalidEmails := [...]string{ "_tobby@message.org", "goat@mail", "finn@.net" }

    // create expression object
    g, _ = GolangExpression()

    // add any combination or number of 'A-Za-z0-9+' for the user identifier of the e-mail to match any alphanumerics
    g, _ = g.AddClass(Uppers)               // adds A-Z; group is then ([A-Z])
    g, _ = g.AddClassToLast(Lowers)         // adds a-z; group is then ([A-Za-z])
    g, _ = g.AddClassToLast(Numbers)        // adds 0-9; group is then ([A-Za-z0-9])
    g, _ = g.ApplyQuantityToLast(OneOrMore) // necessary to have at least one alphanumberic; adds OneOrMore '+' flag; final group: ([A-Za-z0-9]+)

    // add optional single character '.' or '_' character in an e-mail
    g, _ = g.AddFixed(".")                  // adds '.'; group is then (.)
    g, _ = g.AddFixedToLast("_")            // adds '_'; group is then (.|_)
    g, _ = g.ApplyQuantityToLast(ZeroOrOne) // it's optional, OK if it's not there; adds ZeroOrOne '?' flag; final group: (.|_?)

    // add optional second any combination or number of 'A-Za-z0-9+' for the user identifier of the e-mail 
    g, _ = g.AddClass(AlphaNumerics)        // adds A-Za-z0-9; group is then ([A-Za-z0-9])
    g, _ = g.ApplyQuantityToLast(ZeroOrOne) // not necessary to have a second group of alphanumerics; adds ZeroOrOne '?' flag; final group: ([A-Za-z0-9]?)

    // add the '@' in the e-mail
    g, _ = g.AddFixed("@")                  // adds a necessary singular '@'; final group: (@)

    // add the institution identifier of any number of alphanumerics
    g, _ = g.AddClass(AlphaNumerics)        // adds A-Za-z0-9; group is then ([A-Za-z0-9])
    g, _ = g.ApplyQuantityToLast(OneOrMore) // necessary to have at least one alphanumeric; adds OneOrMore '+' flag; final group: ([A-Za-z0-9]+)

    // adds the '.' of the predecessor top-level domain in the e-mail
    g, _ = g.AddFixed(".")                  // adds a necessary singular '.'; final group: (.)

    // adds the top-level domain, supporting specific fixed options
    g, _ = g.AddFixed("com")                // adds 'com' as an option; group is then (com)
    g, _ = g.AddFixedToLast("net")          // adds 'net' as an option; group is then (com|net)
    g, _ = g.AddFixedToLast("org")          // adds 'org' as an option; final group: (com|net|org)

    // create an expression string
    exp, _ := g.Output()                    // Expected output: ([A-Za-z0-9]+)(.|_?)([A-Za-z0-9]?)(@)([A-Za-z0-9]+)(.)(com|net|org)

    var rex = regexp.MustCompile(exp)       // create the regular expression state machine

    fmt.Printf("Expression: %s\n", exp)
    // Output:
    // Expression: ([A-Za-z0-9]+)(.|_?)([A-Za-z0-9]?)(@)([A-Za-z0-9]+)(.)(com|net|org)

    var r string
    for _, r = range(validEmails) { // checks valid emails via the regexp state machine
        fmt.Printf("Attempt: %s, value: %#v\n", r, rex.MatchString(r))
    }
    // Output:
    // Attempt: joe@mail.org, value: true
    // Attempt: john_doe@co.net, value: true
    // Attempt: perry.@place.com, value: true
    // Attempt: _tobby@message.org, value: true

    for _, r = range(invalidEmails) { // checks invalid emails via the regexp state machine
        fmt.Printf("Attempt: %s, value: %#v\n", r, rex.MatchString(r))
    }
    // Output:
    // Attempt: goat@mail, value: false
    // Attempt: finn@.net, value: false
}
```
