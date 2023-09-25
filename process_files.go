package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
	"sync"
)

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		log.Println(ctx, err)
		os.Exit(1)
	}

	log.Println(ctx, "exit")
}

type line struct {
	value string
}

func run(ctx context.Context) error {
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

	go read(ctx, r, in)

	out := make([]chan map[string]int, 2)
	ctx = context.WithValue(ctx, "keyWords", keyWords)
	for i := 0; i < len(out); i++ {
		out[i] = make(chan map[string]int)
		go process(ctx, in, out[i])
	}

	mapOfQuantity := create(ctx, merge(out))

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

func merge(in []chan map[string]int) chan map[string]int {
	out := make(chan map[string]int)

	var wg sync.WaitGroup

	wg.Add(len(in))
	for _, ch := range in {
		go func(ch chan map[string]int) {
			defer wg.Done()
			for mapOfQuantity := range ch {
				out <- mapOfQuantity
			}
		}(ch)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func read(ctx context.Context, r io.Reader, out chan<- line) {
	s := bufio.NewScanner(r)
	for s.Scan() {
		out <- line{value: s.Text()}
	}
	close(out)
}

func process(ctx context.Context, in <-chan line, out chan<- map[string]int) {

	for l := range in {
		formatLine := strings.ToLower(l.value)
		mapOfQuantity := make(map[string]int)
		for _, keyWord := range ctx.Value("keyWords").([]string) {
			re := regexp.MustCompile(fmt.Sprintf("\\W*((?i)%s(?-i))\\W*", keyWord))
			count := len(re.FindAllString(formatLine, -1))
			mapOfQuantity[keyWord] = count
		}
		out <- mapOfQuantity
	}
	close(out)
}

func create(ctx context.Context, in <-chan map[string]int) map[string]int {
	mapOfResults := make(map[string]int)
	for mapOfQuantity := range in {
		for _, keyWord := range ctx.Value("keyWords").([]string) {
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
		wordWithDelimiter := strings.Split(scanner.Text(), " ")
		for _, word := range wordWithDelimiter {
			if word != "-" {
				keyWords = append(keyWords, strings.ToLower(word))
			}
		}
	}
	return keyWords, nil
}
