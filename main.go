package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"unsafe"

	"vimagination.zapto.org/errors"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func run() error {
	var (
		incomplete bool
		plural     bool
	)
	args := os.Args
Loop:
	for {
		args = args[1:]
		if len(args) == 0 {
			return errors.New("no words")
		}
		switch args[0] {
		case "-p":
			incomplete = true
		case "-s":
			plural = true
		default:
			break Loop
		}
	}
	bs, err := ioutil.ReadFile("/usr/share/dict/words")
	if err != nil {
		return errors.WithContext("error reading dictionary: ", err)
	}
	ss := strings.Split(strings.ToLower(byteSliceToString(bs)), "\n")
	if plural {
		for _, word := range ss {
			if !strings.HasSuffix(word, "s") {
				ss = append(ss, word+"s")
			}
		}
	}
	sort.Strings(ss)
	wordsList := make([]map[byte][]string, len(args))
	for n, words := range args {
		wordsList[n] = make(map[byte][]string, strings.Count(words, " ")+1)
		for _, word := range strings.Split(words, " ") {
			var first byte
			if len(word) > 0 {
				first = strings.ToLower(word[:1])[0]
			}
			wordsList[n][first] = append(wordsList[n][first], word)
			sort.Strings(wordsList[n][first])
		}
	}
	buildAnagrams(ss, wordsList, make([]int, 0, len(wordsList)), make([]byte, 0, len(wordsList)), make([][]string, 0, len(wordsList)), incomplete)
	sort.Sort(res)
	var last string
	for _, result := range res {
		if result.Word != last {
			fmt.Println(result.Word)
			last = result.Word
		}
		for _, m := range result.Makeup {
			fmt.Print("->")
			for _, w := range m {
				fmt.Print(" ", w)
			}
			fmt.Println()
		}
	}
	return nil
}

func buildAnagrams(ss []string, wordsList []map[byte][]string, done []int, sofar []byte, words [][]string, incomplete bool) {
	if len(done) < len(wordsList) {
	Loop:
		for n, wordsL := range wordsList {
			for _, d := range done {
				if n == d {
					continue Loop
				}
			}
			d := append(done, n)
			for b, word := range wordsL {
				if b == 0 {
					buildAnagrams(ss, wordsList, d, sofar, words, incomplete)
				} else {
					buildAnagrams(ss, wordsList, d, append(sofar, b), append(words, word), incomplete)
				}
			}
		}
	}
	if (incomplete || len(done) == len(wordsList)) && len(sofar) > 0 {
		s := byteSliceToString(sofar)
		if pos := sort.SearchStrings(ss, s); pos < len(ss) && ss[pos] == s {
			currWords = nil
			printWords(words, make([]int, 0, len(words)))
			res = append(res, result{
				Word:   string(sofar),
				Makeup: currWords,
			})
		}
	}
}

type result struct {
	Word   string
	Makeup [][]string
}

type results []result

var res results

func (r results) Len() int {
	return len(r)
}

func (r results) Less(i, j int) bool {
	return r[i].Word < r[j].Word
}

func (r results) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

var currWords [][]string

func printWords(w [][]string, pos []int) {
	if len(pos) == len(w) {
		words := make([]string, len(pos))
		for n, p := range pos {
			words[n] = w[n][p]
		}
		currWords = append(currWords, words)
		return
	}
	for n := range w[len(pos)] {
		printWords(w, append(pos, n))
	}
}

func byteSliceToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
