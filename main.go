//go:generate peg -strict crontab/crontab.peg
//go:generate peg -strict diskstats/diskstats.peg
//go:generate peg -strict fstab/fstab.peg
//go:generate peg -strict golang/golang.peg
//go:generate peg -strict group/group.peg
//go:generate peg -strict gshadow/gshadow.peg
//go:generate peg -strict maps/maps.peg
//go:generate peg -strict offside/offside.peg
//go:generate peg -strict passwd/passwd.peg
//go:generate peg -strict peg/peg.peg
//go:generate peg -strict peg1/peg1.peg
//go:generate peg -strict popl04/F1.peg
//go:generate peg rfc1459/rfc1459.peg
//go:generate peg rfc2812/rfc2812.peg
//go:generate peg -strict shadow/shadow.peg

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/heyitsanthony/grammars/crontab"
	"github.com/heyitsanthony/grammars/diskstats"
	"github.com/heyitsanthony/grammars/fstab"
	"github.com/heyitsanthony/grammars/golang"
	"github.com/heyitsanthony/grammars/group"
	"github.com/heyitsanthony/grammars/gshadow"
	"github.com/heyitsanthony/grammars/maps"
	"github.com/heyitsanthony/grammars/offside"
	"github.com/heyitsanthony/grammars/passwd"
	"github.com/heyitsanthony/grammars/peg"
	"github.com/heyitsanthony/grammars/peg1"
	"github.com/heyitsanthony/grammars/popl04"
	"github.com/heyitsanthony/grammars/rfc1459"
	"github.com/heyitsanthony/grammars/rfc2812"
	"github.com/heyitsanthony/grammars/shadow"
)

var (
	astFlag = flag.Bool("ast", true, "display AST")
	all     = flag.Bool("all", false, "apply all matching grammars")
	grammar = flag.String("grammar", "", "parse with grammar")
)

type g interface {
	Parse(...int) error
	Init()
	PrintSyntaxTree()
}

func newCrontab(s string) g   { return &crontab.Grammar{Buffer: s, Pretty: true} }
func newDiskstats(s string) g { return &diskstats.Grammar{Buffer: s, Pretty: true} }
func newFstab(s string) g     { return &fstab.Grammar{Buffer: s, Pretty: true} }
func newGo(s string) g        { return &golang.Grammar{Buffer: s, Pretty: true} }
func newGroup(s string) g     { return &group.Grammar{Buffer: s, Pretty: true} }
func newGShadow(s string) g   { return &gshadow.Grammar{Buffer: s, Pretty: true} }
func newMaps(s string) g      { return &maps.Grammar{Buffer: s, Pretty: true} }
func newOffside(s string) g   { return &offside.Grammar{Buffer: s, Pretty: true} }
func newPasswd(s string) g    { return &passwd.Grammar{Buffer: s, Pretty: true} }
func newPeg(s string) g       { return &peg.Grammar{Buffer: s, Pretty: true} }
func newPeg1(s string) g      { return &peg1.Grammar{Buffer: s, Pretty: true} }
func newPopl04(s string) g    { return &popl04.Grammar{Buffer: s, Pretty: true} }
func newRFC1459(s string) g   { return &rfc1459.Grammar{Buffer: s, Pretty: true} }
func newRFC2812(s string) g   { return &rfc2812.Grammar{Buffer: s, Pretty: true} }
func newShadow(s string) g    { return &shadow.Grammar{Buffer: s, Pretty: true} }

var grammars = map[string](func(string) g){
	"crontab":   newCrontab,
	"diskstats": newDiskstats,
	"fstab":     newFstab,
	"go":        newGo,
	"group":     newGroup,
	"gshadow":   newGShadow,
	"maps":      newMaps,
	"offside":   newOffside,
	"passwd":    newPasswd,
	"peg":       newPeg,
	"peg1":      newPeg1,
	"popl04":    newPopl04,
	"rfc1459":   newRFC1459,
	"rfc2812":   newRFC2812,
	"shadow":    newShadow,
}

func main() {
	flag.Parse()
	buffer, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	s := string(buffer)
	var errs []error
	apply := func(l string, f func(string) g) {
		g := f(s)
		g.Init()
		if err := g.Parse(); err != nil {
			errs = append(errs, err)
			return
		}
		if *astFlag {
			g.PrintSyntaxTree()
		} else {
			fmt.Println(l)
		}
		if !*all {
			os.Exit(0)
		}
	}
	if *grammar != "" {
		if f, ok := grammars[*grammar]; ok {
			apply(*grammar, f)
		} else {
			log.Fatalf("unknown grammar %q", *grammar)
		}
		if len(errs) != 0 {
			log.Fatal(errs)
		}
	} else {
		for l, f := range grammars {
			apply(l, f)
		}
		if len(errs) == len(grammars) {
			log.Fatal(errs)
		}
	}
}
