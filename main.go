//go:generate peg golang/golang.peg
//go:generate peg peg/peg.peg
//go:generate peg peg1/peg1.peg
//go:generate peg popl04/F1.peg
//go:generate peg rfc1459/rfc1459.peg
//go:generate peg rfc2812/rfc2812.peg

package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/heyitsanthony/grammars/golang"
	"github.com/heyitsanthony/grammars/peg"
	"github.com/heyitsanthony/grammars/peg1"
	"github.com/heyitsanthony/grammars/popl04"
	"github.com/heyitsanthony/grammars/rfc1459"
	"github.com/heyitsanthony/grammars/rfc2812"
)

var (
	astFlag = flag.Bool("ast", true, "display AST")
	grammar = flag.String("grammar", "", "parse with grammar")
)

type g interface {
	Parse(...int) error
	Init()
	PrintSyntaxTree()
}

func newGo(s string) g      { return &golang.Grammar{Buffer: s, Pretty: true} }
func newPeg(s string) g     { return &peg.Grammar{Buffer: s, Pretty: true} }
func newPeg1(s string) g    { return &peg1.Grammar{Buffer: s, Pretty: true} }
func newPopl04(s string) g  { return &popl04.Grammar{Buffer: s, Pretty: true} }
func newRFC1459(s string) g { return &rfc1459.Grammar{Buffer: s, Pretty: true} }
func newRFC2812(s string) g { return &rfc2812.Grammar{Buffer: s, Pretty: true} }

var grammars = map[string](func(string) g){
	"go":      newGo,
	"peg":     newPeg,
	"peg1":    newPeg1,
	"popl04":  newPopl04,
	"rfc1459": newRFC1459,
	"rfc2812": newRFC2812,
}

func do(g g) error {
	g.Init()
	if err := g.Parse(); err != nil {
		return err
	}
	if *astFlag {
		g.PrintSyntaxTree()
	}
	return nil
}

func main() {
	flag.Parse()
	buffer, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	s := string(buffer)
	var errs []error
	apply := func(f func(string) g) {
		if err := do(f(s)); err == nil {
			os.Exit(0)
		} else {
			errs = append(errs, err)
		}
	}
	if *grammar != "" {
		if f, ok := grammars[*grammar]; ok {
			apply(f)
		} else {
			log.Fatalf("unknown grammar %q", *grammar)
		}
	} else {
		for _, f := range grammars {
			apply(f)
		}
	}
	log.Fatal(errs)
}
