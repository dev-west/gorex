package main

import(
  "fmt"
  // gorex is added for ease of referencing. '.' causes access to gorex.go exports to be immediately accessible (otherwise, must use 'gorex.' in front of everything from that package)
  . "github.com/dev-west/gorex"
  "regexp"
)

func main() {
    var g *Gorex

    validEmails := [...]string{ "joe@mail.org", "john_doe@co.net", "perry.@place.com" }
    invalidEmails := [...]string{ "_tobby@message.org", "goat@mail", "finn@.net" }

    // create expression object
    g, _ = GolangExpression()

    // add any combination or number of 'A-Za-z0-9+' for the user identifier of the e-mail to match any alphanumerics
    g.AddClass(Uppers)               // adds A-Z; group is then ([A-Z])
    g.AddClassToLast(Lowers)         // adds a-z; group is then ([A-Za-z])
    g.AddClassToLast(Digits)         // adds 0-9; group is then ([A-Za-z0-9])
    g.ApplyQuantifier(OneOrMore) // necessary to have at least one alphanumberic; adds OneOrMore '+' flag; final group: ([A-Za-z0-9]+)

    // add optional single character '.' or '_' character in an e-mail
    g.AddFixed(".")                  // adds '.'; group is then (.)
    g.AddFixedToLast("_")            // adds '_'; group is then (.|_)
    g.ApplyQuantifier(ZeroOrOne) // it's optional, OK if it's not there; adds ZeroOrOne '?' flag; final group: (.|_?)

    // add optional second any combination or number of 'A-Za-z0-9+' for the user identifier of the e-mail 
    g.AddClass(AlphaNumerics)        // adds A-Za-z0-9; group is then ([A-Za-z0-9])
    g.ApplyQuantifier(ZeroOrMore) // not necessary to have a second group of alphanumerics; adds ZeroOrMore '*' flag; final group: ([A-Za-z0-9]*)

    // add the '@' in the e-mail
    g.AddFixed("@")                  // adds a necessary singular '@'; final group: (@)

    // add the institution identifier of any number of alphanumerics
    g.AddClass(AlphaNumerics)        // adds A-Za-z0-9; group is then ([A-Za-z0-9])
    g.ApplyQuantifier(OneOrMore) // necessary to have at least one alphanumeric; adds OneOrMore '+' flag; final group: ([A-Za-z0-9]+)

    // adds the '.' of the predecessor top-level domain in the e-mail
    g.AddFixed(".")                  // adds a necessary singular '.'; final group: (.)

    // adds the top-level domain, supporting specific fixed options
    g.AddFixed("com")                // adds 'com' as an option; group is then (com)
    g.AddFixedToLast("net")          // adds 'net' as an option; group is then (com|net)
    g.AddFixedToLast("org")          // adds 'org' as an option; final group: (com|net|org)

    // create an expression string
    exp, _ := g.Output()                    // Expected output: ([A-Za-z0-9]+)(.|_?)([A-Za-z0-9]*)(@)([A-Za-z0-9]+)(.)(com|net|org)

    var rex = regexp.MustCompile(exp)       // create the regular expression state machine

    fmt.Printf("Expression: %s\n", exp)
    // Output:
    // Expression: ([A-Za-z0-9]+)(.|_?)([A-Za-z0-9]*)(@)([A-Za-z0-9]+)(.)(com|net|org)

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
