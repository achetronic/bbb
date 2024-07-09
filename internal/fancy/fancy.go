package fancy

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"unicode/utf8"
)

const (
	//
	Reset  = "\033[0m"
	Bold   = "\033[1m"
	Italic = "\033[3m"
	Under  = "\033[4m"

	//
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	White   = "\033[37m"
)

var (
	StyleMap = map[string]string{
		"Reset":  Reset,
		"Bold":   Bold,
		"Italic": Italic,
		"Under":  Under,

		"Red":     Red,
		"Green":   Green,
		"Yellow":  Yellow,
		"Blue":    Blue,
		"Magenta": Magenta,
		"Cyan":    Cyan,
		"White":   White,
	}
)

// GetReplacedTextFromTemplate get a text with some patterns like '{Something}' inside.
// It returns the text with those patterns replaced with mapped values from styleMap
func GetReplacedTextFromTemplate(template string, styleMap map[string]string) string {
	re := regexp.MustCompile(`\{(\w+)\}`)

	return re.ReplaceAllStringFunc(template, func(match string) string {

		key := match[1 : len(match)-1] // Remove the curly braces
		if val, found := styleMap[key]; found {
			return val
		}
		return match // If the key is not found, return the original match
	})
}

// TODO
func Printf(format string, a ...any) {

	parsedText := GetReplacedTextFromTemplate(format+"{Reset}", StyleMap)
	parsedText = strings.ReplaceAll(parsedText, "\t", "")
	parsedText = strings.TrimSpace(parsedText)

	fmt.Printf(parsedText, a...)
}

// TODO
func Fatalf(format string, a ...any) {
	Printf(format, a...)
	os.Exit(0)
}

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

// GenerateAbbreviation receive a text and return a low-caps abbreviation that can be used as a keyword
// For example, 'Málaga Datacenter' will be abbreviated as 'maldat'
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

// CalculateVisibleLength return the length of a sentence counting only on human-visible chars
func CalculateVisibleLength(s string) int {
	ansiEscape := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	stripped := ansiEscape.ReplaceAllString(s, "")
	return utf8.RuneCountInString(stripped)
}

// TODO
func printSeparator(colWidths []int, left, middle, right, line string) {
	fmt.Print(left)
	for i, width := range colWidths {
		if i > 0 {
			fmt.Print(middle)
		}
		fmt.Print(strings.Repeat(line, width+2))
	}
	fmt.Println(right)
}

// PrintTable prints (header + content) received as parameters as a fancy table
func PrintTable(header string, data [][]string) {
	if len(data) == 0 {
		return
	}

	// Calculate max width of each column
	colWidths := make([]int, len(data[0]))
	for _, rowContent := range data {
		for cellIndex, cellContent := range rowContent {
			cellLength := CalculateVisibleLength(cellContent)
			if cellLength > colWidths[cellIndex] {
				colWidths[cellIndex] = cellLength
			}
		}
	}

	// Print a header row when passed on arguments
	if header != "" {
		totalLength := 0
		for _, colWidth := range colWidths {
			totalLength += colWidth + 3 // Adding space for padding and separator
		}

		// Print header row with different separators as it must be integrated into the main table
		printSeparator(colWidths, "┌", "─", "┐", "─")
		fmt.Println("│ " + Bold + Magenta + header + Reset + strings.Repeat(" ", totalLength-CalculateVisibleLength(header)-3) + " │")
		printSeparator(colWidths, "├", "┬", "┤", "─")
	} else {
		// Print table starting symbols when header is not passed
		printSeparator(colWidths, "┌", "┬", "┐", "─")
	}

	// Print table with fancy separators and borders
	for rowIndex, rowContent := range data {
		fmt.Print("│")
		for cellIndex, cellContent := range rowContent {
			cellLength := CalculateVisibleLength(cellContent)
			if rowIndex == 0 {
				// Print fancy header for first row
				fmt.Print(" " + Bold + Blue + cellContent + Reset + strings.Repeat(" ", colWidths[cellIndex]-cellLength) + " │")
			} else {
				fmt.Print(" " + cellContent + strings.Repeat(" ", colWidths[cellIndex]-cellLength) + " │")
			}
		}
		fmt.Println()

		// Different separators according to the row to distinguish starting, middle and final ones
		if rowIndex == 0 {
			printSeparator(colWidths, "├", "┼", "┤", "─")
		} else if rowIndex == len(data)-1 {
			printSeparator(colWidths, "└", "┴", "┘", "─")
		} else {
			printSeparator(colWidths, "├", "┼", "┤", "─")
		}
	}
}
