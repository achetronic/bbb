package list

import (
	"fmt"
	"strings"
)

const (
	Reset   = "\033[0m"
	Bold    = "\033[1m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	White   = "\033[37m"
)

// Define a map for character substitutions
var substitutionMap = map[rune]string{
	'á': "a", 'à': "a", 'ä': "a", 'â': "a", 'ã': "a",
	'é': "e", 'è': "e", 'ë': "e", 'ê': "e",
	'í': "i", 'ì': "i", 'ï': "i", 'î': "i",
	'ó': "o", 'ò': "o", 'ö': "o", 'ô': "o", 'õ': "o",
	'ú': "u", 'ù': "u", 'ü': "u", 'û': "u",
	'ñ': "n",
	'ç': "c",
	// Add other substitutions as needed
}

// SubstituteCharacters replaces characters based on the substitution map
func SubstituteCharacters(input string) string {
	var output strings.Builder
	for _, char := range input {
		if replacement, found := substitutionMap[char]; found {
			output.WriteString(replacement)
		} else {
			// Use the original character if no substitution is found
			output.WriteRune(char)
		}
	}
	return output.String()
}

// TODO
func GenerateAbbreviation(text string) string {

	// Convert text to lower case and split by spaces
	lowerText := strings.ToLower(text)
	normalizedText := SubstituteCharacters(lowerText)
	words := strings.Fields(normalizedText)

	// Several words? get initial chars
	if len(words) > 1 {
		abbreviation := ""
		for _, word := range words {
			abbreviation += string(word[0:3])
		}
		return abbreviation
	}

	// Get first chars for single words
	if len(words) == 1 {
		word := words[0]
		if len(word) <= 4 {
			return word
		}
		return word[:4]
	}

	return ""
}

// TODO
func PrintSeparator(colWidths []int, left, middle, right, line string) {
	fmt.Print(left)
	for i, width := range colWidths {
		if i > 0 {
			fmt.Print(middle)
		}
		fmt.Print(strings.Repeat(line, width+2))
	}
	fmt.Println(right)
}

// TODO
func PrintTable(data [][]string) {

	// Calculate max width of each column
	colWidths := make([]int, len(data[0]))
	for _, row := range data {
		for i, cell := range row {
			if len(cell) > colWidths[i] {
				colWidths[i] = len(cell)
			}
		}
	}

	// Print table with fancy separators and borders
	PrintSeparator(colWidths, "┌", "┬", "┐", "─")
	for rowIndex, rowContent := range data {
		fmt.Print("│")
		for cellIndex, cellContent := range rowContent {
			if rowIndex == 0 {
				// Print fancy header for first row
				fmt.Print(" " + Bold + Blue + cellContent + Reset + strings.Repeat(" ", colWidths[cellIndex]-len(cellContent)) + " │")
			} else {
				fmt.Print(" " + cellContent + strings.Repeat(" ", colWidths[cellIndex]-len(cellContent)) + " │")
			}
		}

		fmt.Println()

		// Different separators according to the row to distinguish starting, middle and final ones
		if rowIndex == 0 {
			PrintSeparator(colWidths, "├", "┼", "┤", "─")
		} else if rowIndex == len(data)-1 {
			PrintSeparator(colWidths, "└", "┴", "┘", "─")
		} else {
			PrintSeparator(colWidths, "├", "┼", "┤", "─")
		}
	}
}
