package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/go-chi/chi/middleware"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal(server("localhost:7711"))
	}
	dict := loadDict("german.txt")
	letters := strings.ToLower(os.Args[1])
	words := findValidWords(letters, dict)
	os.Stdout.Write(formatWords(words, len(letters)))
}

func server(addr string) error {
	r := http.NewServeMux()
	r.Handle("/", &handler{loadDict("german.txt")})
	srv := &http.Server{
		WriteTimeout: time.Second * 10,
		Addr:         addr,
		Handler:      middleware.Logger(r),
	}
	return srv.ListenAndServe()
}

type handler struct {
	dict []string
}

var lettersRe = regexp.MustCompile(`^[a-zäöü]{1,50}$`)

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	letters := strings.ToLower(r.URL.Path[1:])
	fmt.Println(letters)
	if !lettersRe.MatchString(letters) {
		http.Error(w, "kys", 400)
		return
	}
	words := findValidWords(letters, h.dict)
	w.Write(formatWords(words, len(letters)))
}

func formatWords(words []string, length int) []byte {
	buf := &bytes.Buffer{}
	pools := make([][]string, length)
	for _, word := range words {
		pool := len(word) - 1
		pools[pool] = append(pools[pool], word)
	}
	for i, pool := range pools {
		sort.Strings(pool)
		fmt.Fprintf(buf, "\n%d\n", i+1)
		for _, word := range pool {
			fmt.Fprintln(buf, word)
		}
	}
	return buf.Bytes()
}

func loadDict(path string) []string {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	reader := bufio.NewReaderSize(file, 1024*1000*8)
	lines := []string{}
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return lines
		}
		word := strings.ToLower(line[:len(line)-2])
		lines = append(lines, word)
	}
}

func findValidWords(letters string, dict []string) []string {
	validWords := []string{}
	for _, word := range dict {
		if validWord(letters, word) {
			validWords = append(validWords, word)
		}
	}
	return validWords
}

func validWord(letters string, word string) bool {
	if len(word) > len(letters) {
		return false
	}
	remaining := []rune(letters)
	for _, letter := range word {
		found := false
		for i, r := range remaining {
			if r == letter {
				remaining = append(remaining[:i], remaining[i+1:]...)
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}
