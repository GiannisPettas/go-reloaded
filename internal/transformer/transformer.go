package transformer

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// TokenType represents the type of a token
type TokenType int

const (
	Word TokenType = iota
	Command
	Punctuation
	Quote
	LineBreak
)

// Token represents a parsed element from text
type Token struct {
	Type  TokenType
	Value string
}

// TokenizeText splits text into tokens (words, commands, punctuation, quotes, line breaks)
func TokenizeText(text string) []Token {
	if text == "" {
		return []Token{}
	}
	
	var tokens []Token
	
	// Regular expression to match tokens including line breaks
	// Matches: commands like (up), (hex), (up, 2), words, punctuation, quotes, newlines
	re := regexp.MustCompile(`\([^)]+\)|[a-zA-Z0-9]+|[,.!?;]|'|\n`)
	
	matches := re.FindAllString(text, -1)
	
	for _, match := range matches {
		token := Token{Value: match}
		
		// Determine token type
		if match == "\n" {
			token.Type = LineBreak
		} else if strings.HasPrefix(match, "(") && strings.HasSuffix(match, ")") {
			token.Type = Command
		} else if match == "'" {
			token.Type = Quote
		} else if strings.ContainsAny(match, ",.!?;") {
			token.Type = Punctuation
		} else {
			token.Type = Word
		}
		
		tokens = append(tokens, token)
	}
	
	return tokens
}

// ConvertHex converts hexadecimal numbers to decimal
func ConvertHex(tokens []Token) []Token {
	var result []Token
	
	for i := 0; i < len(tokens); i++ {
		if i > 0 && tokens[i].Type == Command && tokens[i].Value == "(hex)" {
			// Convert previous token if it's a valid hex number
			prevToken := tokens[i-1]
			if prevToken.Type == Word {
				if converted := convertHexToDecimal(prevToken.Value); converted != "" {
					// Replace previous token with converted value
					result[len(result)-1] = Token{Type: Word, Value: converted}
					// Skip the (hex) command
					continue
				}
			}
		}
		result = append(result, tokens[i])
	}
	
	return result
}

// convertHexToDecimal converts a hex string to decimal, returns empty string if invalid
func convertHexToDecimal(hex string) string {
	isNegative := false
	if strings.HasPrefix(hex, "-") {
		isNegative = true
		hex = hex[1:]
	}
	
	// Validate hex characters
	for _, char := range hex {
		if !((char >= '0' && char <= '9') || (char >= 'a' && char <= 'f') || (char >= 'A' && char <= 'F')) {
			return "" // Invalid hex
		}
	}
	
	if hex == "" {
		return "" // Empty after removing negative sign
	}
	
	// Convert hex to decimal
	var decimal int64 = 0
	for _, char := range hex {
		decimal *= 16
		if char >= '0' && char <= '9' {
			decimal += int64(char - '0')
		} else if char >= 'a' && char <= 'f' {
			decimal += int64(char - 'a' + 10)
		} else if char >= 'A' && char <= 'F' {
			decimal += int64(char - 'A' + 10)
		}
	}
	
	if isNegative {
		return "-" + strings.Trim(strings.Replace(fmt.Sprintf("%d", decimal), "-", "", 1), " ")
	}
	return fmt.Sprintf("%d", decimal)
}

// ConvertBinary converts binary numbers to decimal
func ConvertBinary(tokens []Token) []Token {
	var result []Token
	
	for i := 0; i < len(tokens); i++ {
		if i > 0 && tokens[i].Type == Command && tokens[i].Value == "(bin)" {
			// Convert previous token if it's a valid binary number
			prevToken := tokens[i-1]
			if prevToken.Type == Word {
				if converted := convertBinaryToDecimal(prevToken.Value); converted != "" {
					// Replace previous token with converted value
					result[len(result)-1] = Token{Type: Word, Value: converted}
					// Skip the (bin) command
					continue
				}
			}
		}
		result = append(result, tokens[i])
	}
	
	return result
}

// convertBinaryToDecimal converts a binary string to decimal, returns empty string if invalid
func convertBinaryToDecimal(binary string) string {
	isNegative := false
	if strings.HasPrefix(binary, "-") {
		isNegative = true
		binary = binary[1:]
	}
	
	// Validate binary characters (only 0 and 1)
	for _, char := range binary {
		if char != '0' && char != '1' {
			return "" // Invalid binary
		}
	}
	
	if binary == "" {
		return "" // Empty after removing negative sign
	}
	
	// Convert binary to decimal
	var decimal int64 = 0
	for _, char := range binary {
		decimal *= 2
		if char == '1' {
			decimal += 1
		}
	}
	
	if isNegative {
		return "-" + fmt.Sprintf("%d", decimal)
	}
	return fmt.Sprintf("%d", decimal)
}

// ApplyCaseTransform applies case transformations (up), (low), (cap)
func ApplyCaseTransform(tokens []Token) []Token {
	var result []Token
	
	for i := 0; i < len(tokens); i++ {
		if i > 0 && tokens[i].Type == Command {
			// Check for case transformation commands
			if cmd, count := parseCaseCommand(tokens[i].Value); cmd != "" {
				// Apply transformation to previous 'count' words
				wordCount := 0
				for j := len(result) - 1; j >= 0 && wordCount < count; j-- {
					if result[j].Type == Word {
						var transformed string
						switch cmd {
						case "up":
							transformed = strings.ToUpper(result[j].Value)
						case "low":
							transformed = strings.ToLower(result[j].Value)
						case "cap":
							transformed = strings.Title(strings.ToLower(result[j].Value))
						}
						result[j] = Token{Type: Word, Value: transformed}
						wordCount++
					}
				}
				// Skip the command
				continue
			}
		}
		result = append(result, tokens[i])
	}
	
	return result
}

// parseCaseCommand parses case commands like (up), (low, 2), (cap, 3)
// Returns command type and word count (1 if not specified)
func parseCaseCommand(cmd string) (string, int) {
	if cmd == "(up)" {
		return "up", 1
	}
	if cmd == "(low)" {
		return "low", 1
	}
	if cmd == "(cap)" {
		return "cap", 1
	}
	
	// Check for numbered commands like (up, 2)
	re := regexp.MustCompile(`\((up|low|cap),\s*(\d+)\)`)
	matches := re.FindStringSubmatch(cmd)
	if len(matches) == 3 {
		command := matches[1]
		if count, err := strconv.Atoi(matches[2]); err == nil {
			return command, count
		}
	}
	
	return "", 0
}

// FixPunctuationSpacing fixes spacing around punctuation marks
func FixPunctuationSpacing(tokens []Token) []Token {
	var result []Token
	
	for i := 0; i < len(tokens); i++ {
		if tokens[i].Type == Punctuation {
			// Attach punctuation to previous word if exists
			if len(result) > 0 && result[len(result)-1].Type == Word {
				// Collect consecutive punctuation marks
				punctuation := tokens[i].Value
				j := i + 1
				for j < len(tokens) && tokens[j].Type == Punctuation {
					punctuation += tokens[j].Value
					j++
				}
				
				// Attach to previous word
				result[len(result)-1].Value += punctuation
				
				// Skip processed punctuation tokens
				i = j - 1
			} else {
				// No previous word, keep punctuation as is
				result = append(result, tokens[i])
			}
		} else {
			result = append(result, tokens[i])
		}
	}
	
	return result
}

// RepositionQuotes moves single quotes to correct positions around words
func RepositionQuotes(tokens []Token) []Token {
	var result []Token
	
	for i := 0; i < len(tokens); i++ {
		if tokens[i].Type == Quote && tokens[i].Value == "'" {
			// Look for matching closing quote
			quoteStart := i
			quoteEnd := -1
			
			// Find the next quote
			for j := i + 1; j < len(tokens); j++ {
				if tokens[j].Type == Quote && tokens[j].Value == "'" {
					quoteEnd = j
					break
				}
			}
			
			if quoteEnd != -1 {
				// Found matching pair, reposition quotes
				var quotedContent []Token
				for k := quoteStart + 1; k < quoteEnd; k++ {
					quotedContent = append(quotedContent, tokens[k])
				}
				
				if len(quotedContent) > 0 {
					// Attach opening quote to first word
					if quotedContent[0].Type == Word {
						quotedContent[0].Value = "'" + quotedContent[0].Value
					}
					
					// Attach closing quote to last word
					if quotedContent[len(quotedContent)-1].Type == Word {
						quotedContent[len(quotedContent)-1].Value += "'"
					}
					
					// Add quoted content to result
					result = append(result, quotedContent...)
				}
				
				// Skip to after closing quote
				i = quoteEnd
			} else {
				// No matching quote found, keep as is
				result = append(result, tokens[i])
			}
		} else {
			result = append(result, tokens[i])
		}
	}
	
	return result
}

// CorrectArticles changes "a" to "an" before vowels and "h"
func CorrectArticles(tokens []Token) []Token {
	var result []Token
	
	for i := 0; i < len(tokens); i++ {
		if tokens[i].Type == Word && strings.ToLower(tokens[i].Value) == "a" {
			// Check if next token is a word starting with vowel or 'h'
			if i+1 < len(tokens) && tokens[i+1].Type == Word {
				nextWord := strings.ToLower(tokens[i+1].Value)
				if len(nextWord) > 0 {
					firstChar := nextWord[0]
					if firstChar == 'a' || firstChar == 'e' || firstChar == 'i' || firstChar == 'o' || firstChar == 'u' || firstChar == 'h' {
						// Change "a" to "an"
						if tokens[i].Value == "a" {
							tokens[i].Value = "an"
						} else if tokens[i].Value == "A" {
							tokens[i].Value = "An"
						}
					}
				}
			}
		}
		result = append(result, tokens[i])
	}
	
	return result
}

// ApplyAllTransformations applies all transformations in the correct order
// This ensures proper command chaining and left-to-right execution
func ApplyAllTransformations(tokens []Token) []Token {
	// Apply transformations in order:
	// 1. Numeric conversions (hex, bin) - these can chain
	// 2. Case transformations
	// 3. Article corrections
	// 4. Quote repositioning
	// 5. Punctuation spacing (last, as it modifies token structure)
	
	result := tokens
	
	// Apply numeric conversions multiple times to handle chaining
	for i := 0; i < 3; i++ { // Max 3 iterations to handle reasonable chaining
		prevLen := len(result)
		result = ConvertHex(result)
		result = ConvertBinary(result)
		// If no changes in length, break early
		if len(result) == prevLen {
			break
		}
	}
	
	// Apply other transformations
	result = ApplyCaseTransform(result)
	result = CorrectArticles(result)
	result = RepositionQuotes(result)
	result = FixPunctuationSpacing(result)
	
	return result
}

// ApplyAllTransformationsWithContext applies transformations with overlap context
// This handles cross-chunk word references for commands like (up, n)
func ApplyAllTransformationsWithContext(currentChunk, overlapContext string) []Token {
	// Merge overlap context with current chunk
	mergedText := overlapContext
	if overlapContext != "" && currentChunk != "" {
		mergedText += " " + currentChunk
	} else if currentChunk != "" {
		mergedText = currentChunk
	}
	
	// Tokenize merged text
	tokens := TokenizeText(mergedText)
	
	// Apply all transformations
	result := ApplyAllTransformations(tokens)
	
	// Return only the portion that corresponds to current chunk + context
	// For cross-chunk commands, we need to return the transformed context too
	return result
}

// TokensToString converts tokens back to string, preserving line breaks
func TokensToString(tokens []Token) string {
	var result strings.Builder
	
	for i, token := range tokens {
		if token.Type == LineBreak {
			result.WriteString(token.Value)
		} else {
			if i > 0 && tokens[i-1].Type != LineBreak && token.Type != LineBreak {
				result.WriteString(" ")
			}
			result.WriteString(token.Value)
		}
	}
	
	return result.String()
}