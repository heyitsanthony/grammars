//go:generate peg peg.peg
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type hideExtent struct {
	begin int
	end   int
	text  string
}

type hideExtents struct {
	hides []bool
}

func newHideExtents(l int) *hideExtents { return &hideExtents{make([]bool, l)} }

func (he *hideExtents) hide(begin, end int, text string) {
	for i := begin; i < end; i++ {
		he.hides[i] = true
	}
}

func (he *hideExtents) apply(buf string) string {
	v := []rune{}
	for i, r := range buf {
		if !he.hides[i] {
			v = append(v, r)
		}
	}
	return string(v)
}

func main() {
	buffer, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	s := string(buffer)
	p := &PegHide{hideExtents: newHideExtents(len(s)), Buffer: s}
	p.Init()
	if err := p.Parse(); err != nil {
		log.Fatal(err)
	}
	p.Execute()
	fmt.Println(p.apply(s))
}
