package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type FormatCommand struct {
	FileName string
}

func (fc *FormatCommand) Execute() {
	file, err := os.ReadFile(fc.FileName)
	if err != nil {
		fmt.Println(err)
		return
	}
	if filepath.Ext(fc.FileName) != ".txt" {
		fmt.Println("file extension isn't .txt")
		return
	}

	text := string(file)

	formatText := strings.ReplaceAll(text, "\n\n", "\n\n\t")
	formatTextWithDot := addDotBeforeUppercase(formatText)

	fileNameWithoutExt := fileNameWithoutExtTrimSuffix(fc.FileName)
	changedFileName := fileNameWithoutExt + "_format.txt"
	newFile, err := os.Create(changedFileName)
	defer newFile.Close()
	_, err = newFile.WriteString(formatTextWithDot)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("formatted text is saved to ", changedFileName)
}

func addDotBeforeUppercase(input string) string {
	reSpace := regexp.MustCompile(`([a-zа-яё]) ([A-ZА-ЯЁ])`)
	outputSpace := reSpace.ReplaceAllString(input, "${1}. ${2}")
	reEndLine := regexp.MustCompile(`([a-zа-яё])\n`)
	output := reEndLine.ReplaceAllString(outputSpace, "${1}.\n${2}")
	return output
}

func fileNameWithoutExtTrimSuffix(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}
