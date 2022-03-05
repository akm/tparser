package runes_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func BenchmarkPickUpViaPointer(b *testing.B) {
	runesPtr := readText(b)
	length := len(*runesPtr)
	dest := devnull()
	defer dest.Close()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		x := i % length
		y := x + 100
		if y >= length {
			y = length - 1
		}
		sub := (*runesPtr)[x:y]
		fmt.Fprintf(dest, "%s\n", string(sub))
	}
}

func BenchmarkPickUpFromCopy(b *testing.B) {
	runes := *readText(b)
	length := len(runes)
	dest := devnull()
	defer dest.Close()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		x := i % length
		y := x + 100
		if y >= length {
			y = length - 1
		}
		sub := runes[x:y]
		fmt.Fprintf(dest, "%s\n", string(sub))
	}
}

func readText(b *testing.B) *[]rune {
	f, err := os.Open("../Grammer.md")
	if err != nil {
		b.Fatal(err)
	}
	defer f.Close()
	d, err := ioutil.ReadAll(f)
	if err != nil {
		b.Fatal(err)
	}
	s := string(d)
	l := len(s)
	m := 10_000_000 / l
	res := ""
	for i := 0; i < m; i++ {
		res += s
	}
	r := []rune(res)
	return &r
}

func devnull() *os.File {
	f, err := os.OpenFile("/dev/null", os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	return f
}
