package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
	"sync"
)

func main() {
	if err := run(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
	log.Println("exit")
}

type line struct {
	value string
}

func run() error {
	keyWords, err := createSliceOfKeyWordsFromFile("key_words.txt")
	if err != nil {
		return err
	}

	r, err := os.Open("input_test")
	if err != nil {
		return err
	}
	defer r.Close()

	in := make(chan line)

	go read(r, in)

	out := make(chan map[string]int)

	var wg sync.WaitGroup
	goroutinesNumber := 2
	wg.Add(goroutinesNumber)
	for i := 0; i < goroutinesNumber; i++ {
		go process(in, out, keyWords, &wg)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	mapOfQuantity := create(out, keyWords)
	w, err := os.OpenFile("output_test", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer w.Close()

	for key, val := range mapOfQuantity {
		_, err := fmt.Fprintf(w, "%s: %d\n", key, val)
		if err != nil {
			return err
		}
	}
	return nil
}

func read(r io.Reader, out chan<- line) {
	s := bufio.NewScanner(r)
	for s.Scan() {
		out <- line{value: s.Text()}
	}
	close(out)
}

func process(in <-chan line, out chan<- map[string]int, keyWords []string, wg *sync.WaitGroup) {
	defer wg.Done()
	for l := range in {
		formatLine := strings.ToLower(l.value)
		mapOfQuantity := make(map[string]int)
		for _, keyWord := range keyWords {
			re := regexp.MustCompile(fmt.Sprintf("\\W*((?i)%s(?-i))\\W*", keyWord))
			count := len(re.FindAllString(formatLine, -1))
			mapOfQuantity[keyWord] = count
		}
		out <- mapOfQuantity
	}
}

func create(in <-chan map[string]int, keyWords []string) map[string]int {
	mapOfResults := make(map[string]int)
	for mapOfQuantity := range in {
		for _, keyWord := range keyWords {
			mapOfResults[keyWord] += mapOfQuantity[keyWord]
		}
	}
	totalScore := 0
	for _, val := range mapOfResults {
		totalScore += val
	}
	mapOfResults["всего"] = totalScore
	return mapOfResults
}

func createSliceOfKeyWordsFromFile(fileNameKeyWords string) ([]string, error) {
	fileWithKeyWords, err := os.Open(fileNameKeyWords)
	if err != nil {
		return nil, err
	}
	defer fileWithKeyWords.Close()

	var keyWords []string

	scanner := bufio.NewScanner(fileWithKeyWords)
	for scanner.Scan() {
		word := strings.ToLower(scanner.Text())
		keyWords = append(keyWords, word)
	}
	return keyWords, nil
}
