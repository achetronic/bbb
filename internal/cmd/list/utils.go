package list

import "strings"

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

// GetScopesByScope TODO
func GetScopesByScope(scopes ListScopesResponseT) (result map[string][]ScopeT) {

	result = make(map[string][]ScopeT)

	for _, scope := range scopes.Items {
		result[scope.ScopeId] = append(result[scope.ScopeId], scope)
	}

	return result
}
