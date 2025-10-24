package domain

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"
)

// GetRandomValueS returns a random value from a slice based on given probabilities
func GetRandomValueS(lineCount int) (string, string) {
	productIndex := lineCount % len(cardProducts)
	modelIndex := lineCount % 10
	modelIndexVals := []int{2, 8}
	if slices.Contains(modelIndexVals, modelIndex) {
		modelIndex = 1
	} else {
		modelIndex = 0
	}
	return cardProducts[productIndex], cardModels[modelIndex]

}


// ReplaceFakeBinsLines processes a line by replacing specific substrings
func ReplaceFakeBinsLines(line string, lineCount int) (string, error) {
	line = strings.ReplaceAll(line, "(bin, bandeira)", "(bin, bandeira, modalidade_final, produto_final)")
	prod, mod := GetRandomValueS(lineCount)
	line = strings.ReplaceAll(line, ");", fmt.Sprintf(", '%s', %s);", mod, prod))
	return line, nil
}

// ReplicaFakeBinsFiles is a placeholder function for future implementation
func ReplicaFakeBinsFiles(filename string) error {
	// open the file
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	// loop reading lines
	scanner := bufio.NewScanner(file)
	count := 0
	for scanner.Scan() {
		line := scanner.Text()
		// for each line call ReplaceFakeBinsLines
		newLine, err := ReplaceFakeBinsLines(line, count)
		if err != nil {
			return err
		}
		// print new line
		fmt.Println(newLine)
		// handle errors
		count++
	}
	return nil

}


