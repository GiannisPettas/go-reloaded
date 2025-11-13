package transformer

import (
	"go-reloaded/internal/config"
	"strconv"
	"strings"
)

// Token types
const (
	WORD = iota
	COMMAND
	PUNCTUATION
	SPACE
	NEWLINE
)

type Token struct {
	Type  int
	Value string
}

// Low-level FSM states
const (
	STATE_TEXT = iota
	STATE_COMMAND
)

// High-level FSM for token processing
type TokenProcessor struct {
	tokens   []Token
	tokenIdx int
	output   strings.Builder
}

// ProcessText - Single pass dual FSM implementation
// main entry point
func ProcessText(text string) string {
	if text == "" {
		return ""
	}

	runes := []rune(text)
	processor := NewTokenProcessor()

	state := STATE_TEXT
	var wordBuilder strings.Builder // Accumulates characters for current word
	var cmdBuilder strings.Builder  // Accumulates characters for current command

	for i := 0; i < len(runes); i++ {
		r := runes[i]

		switch state {
		case STATE_TEXT:
			switch r {
			case '(':
				// Look ahead to see if this is a valid command (max 10 chars)
				if i+1 < len(runes) {
					// Find the closing parenthesis within 10 characters
					closeParen := -1
					maxLookAhead := i + 11 // i+1 + 10 chars max, assuming that no command has more than 10 runes
					if maxLookAhead > len(runes) {
						maxLookAhead = len(runes)
					}
					for j := i + 1; j < maxLookAhead; j++ {
						if runes[j] == ')' {
							closeParen = j
							break
						}
					}

					if closeParen != -1 {
						// Extract potential command
						potentialCmd := string(runes[i+1 : closeParen])
						if processor.isValidCommand(potentialCmd) {
							// Valid command - flush current word and switch to command state
							if wordBuilder.Len() > 0 {
								processor.addToken(Token{WORD, wordBuilder.String()})
								wordBuilder.Reset()
							}
							state = STATE_COMMAND
							break
						} else {
							// Invalid command - treat entire thing as word
							wordBuilder.WriteString(string(runes[i : closeParen+1]))
							i = closeParen // Skip to after closing paren
							break
						}
					}
				}

				// No closing paren found - treat as regular character
				wordBuilder.WriteRune(r)
			case ' ', '\t':
				// Flush word and add space
				if wordBuilder.Len() > 0 {
					processor.addToken(Token{WORD, wordBuilder.String()})
					wordBuilder.Reset()
				}
				processor.addToken(Token{SPACE, " "})
			case '\n':
				// Flush word and add newline
				if wordBuilder.Len() > 0 {
					processor.addToken(Token{WORD, wordBuilder.String()})
					wordBuilder.Reset()
				}
				processor.addToken(Token{NEWLINE, "\n"})
			case ',', '.', '!', '?', ';', ':':
				// Flush word and add punctuation
				if wordBuilder.Len() > 0 {
					processor.addToken(Token{WORD, wordBuilder.String()})
					wordBuilder.Reset()
				}
				processor.addToken(Token{PUNCTUATION, string(r)})
			default:
				wordBuilder.WriteRune(r)
			}

		case STATE_COMMAND:
			if r == ')' {
				// Process valid command
				processor.processCommand(cmdBuilder.String())
				cmdBuilder.Reset()
				state = STATE_TEXT
			} else {
				cmdBuilder.WriteRune(r)
			}
		}
	}

	// Flush remaining word
	if wordBuilder.Len() > 0 {
		processor.addToken(Token{WORD, wordBuilder.String()})
	}

	// Flush all tokens to output
	processor.flushTokens()

	// Post-process articles and quotes
	result := processor.output.String()
	result = fixArticles(result)
	return fixQuotes(result)
}

// --------------- CORE PROCESSING FUNCTIONS  ---------------
func (tp *TokenProcessor) addToken(token Token) {
	if tp.tokenIdx < len(tp.tokens) {
		tp.tokens[tp.tokenIdx] = token
		tp.tokenIdx++
	} else {
		// Buffer is full, flush first half to output
		halfSize := len(tp.tokens) / 2
		for i := 0; i < halfSize; i++ {
			token := tp.tokens[i]
			switch token.Type {
			case WORD:
				if tp.output.Len() > 0 && !strings.HasSuffix(tp.output.String(), " ") && !strings.HasSuffix(tp.output.String(), "\n") {
					tp.output.WriteByte(' ')
				}
				tp.output.WriteString(token.Value)
			case PUNCTUATION:
				// Handle spacing for punctuation
				if token.Value == "(" {
					// Opening parenthesis preserves space before it
					tp.output.WriteString(token.Value)
				} else {
					// Other punctuation - remove trailing space
					result := tp.output.String()
					if strings.HasSuffix(result, " ") {
						tp.output.Reset()
						tp.output.WriteString(result[:len(result)-1])
					}
					tp.output.WriteString(token.Value)
				}
			case SPACE:
				if tp.output.Len() > 0 && !strings.HasSuffix(tp.output.String(), " ") && !strings.HasSuffix(tp.output.String(), "\n") {
					tp.output.WriteByte(' ')
				}
			case NEWLINE:
				tp.output.WriteByte('\n')
			}
		}

		// Shift remaining tokens to beginning
		for i := 0; i < halfSize; i++ {
			tp.tokens[i] = tp.tokens[halfSize+i]
		}
		tp.tokenIdx = halfSize

		// Add new token
		tp.tokens[tp.tokenIdx] = token
		tp.tokenIdx++
	}
}

func (tp *TokenProcessor) processCommand(cmdValue string) {
	// Check if command is valid before processing
	if !tp.isValidCommand(cmdValue) {
		// Invalid command - ignore it completely
		return
	}

	if tp.tokenIdx == 0 {
		return // No words to transform
	}

	// Find last word token
	lastWordIdx := -1
	for i := tp.tokenIdx - 1; i >= 0; i-- {
		if tp.tokens[i].Type == WORD {
			lastWordIdx = i
			break
		}
	}

	if lastWordIdx == -1 {
		return
	}

	// Parse command
	if strings.Contains(cmdValue, ",") {
		parts := strings.Split(cmdValue, ",")
		if len(parts) == 2 {
			cmd := strings.TrimSpace(parts[0])
			countStr := strings.TrimSpace(parts[1])
			if count, err := strconv.Atoi(countStr); err == nil && count > 0 {
				// Find word indices to transform (in reverse order)
				var wordIndices []int
				for i := tp.tokenIdx - 1; i >= 0 && len(wordIndices) < count; i-- {
					if tp.tokens[i].Type == WORD {
						wordIndices = append(wordIndices, i)
					}
				}
				// Transform words in forward order
				for i := len(wordIndices) - 1; i >= 0; i-- {
					idx := wordIndices[i]
					// Mark articles transformed by (up) command
					if cmd == "up" && (tp.tokens[idx].Value == "a" || tp.tokens[idx].Value == "an") {
						tp.tokens[idx].Value = "UP_" + tp.transformWord(tp.tokens[idx].Value, cmd)
					} else {
						tp.tokens[idx].Value = tp.transformWord(tp.tokens[idx].Value, cmd)
					}
				}
			}
		}
	} else {
		// Single word command
		switch cmdValue {
		case "hex":
			if val, err := strconv.ParseInt(tp.tokens[lastWordIdx].Value, 16, 64); err == nil {
				tp.tokens[lastWordIdx].Value = strconv.FormatInt(val, 10)
			}
		case "bin":
			if val, err := strconv.ParseInt(tp.tokens[lastWordIdx].Value, 2, 64); err == nil {
				tp.tokens[lastWordIdx].Value = strconv.FormatInt(val, 10)
			}
		default:
			// Mark articles transformed by (up) command
			if cmdValue == "up" && (tp.tokens[lastWordIdx].Value == "a" || tp.tokens[lastWordIdx].Value == "an") {
				tp.tokens[lastWordIdx].Value = "UP_" + tp.transformWord(tp.tokens[lastWordIdx].Value, cmdValue)
			} else {
				tp.tokens[lastWordIdx].Value = tp.transformWord(tp.tokens[lastWordIdx].Value, cmdValue)
			}
		}
	}
}

// --------------- POST-PROCESSING PIPELINE ---------------
func fixQuotes(text string) string {
	runes := []rune(text)
	var result strings.Builder

	singleQuoteCount := 0
	doubleQuoteCount := 0

	for i := 0; i < len(runes); i++ {
		r := runes[i]

		if r == '\'' {
			singleQuoteCount++
			if singleQuoteCount%2 == 1 {
				// Odd quote - stick to right letter
				result.WriteRune(r)
				// Skip space after quote if present
				if i+1 < len(runes) && runes[i+1] == ' ' {
					i++ // Skip the space
				}
			} else {
				// Even quote - stick to left letter
				// Remove space before quote if present
				resultStr := result.String()
				if strings.HasSuffix(resultStr, " ") {
					result.Reset()
					result.WriteString(resultStr[:len(resultStr)-1])
				}
				result.WriteRune(r)
			}
		} else if r == '"' {
			doubleQuoteCount++
			if doubleQuoteCount%2 == 1 {
				// Odd quote - stick to right letter
				result.WriteRune(r)
				// Skip space after quote if present
				if i+1 < len(runes) && runes[i+1] == ' ' {
					i++ // Skip the space
				}
			} else {
				// Even quote - stick to left letter
				// Remove space before quote if present
				resultStr := result.String()
				if strings.HasSuffix(resultStr, " ") {
					result.Reset()
					result.WriteString(resultStr[:len(resultStr)-1])
				}
				result.WriteRune(r)
			}
		} else {
			result.WriteRune(r)
		}
	}

	return result.String()
}

func fixArticles(text string) string {
	// Process line by line to preserve line breaks
	lines := strings.Split(text, "\n")
	for lineIdx, line := range lines {
		if line == "" {
			continue
		}

		words := strings.Fields(line)
		for i := 0; i < len(words)-1; i++ {
			switch words[i] {
			case "a", "A", "an", "An", "AN", "UP_A", "UP_AN":
				nextWord := words[i+1]
				if len(nextWord) > 0 {
					// Remove punctuation for vowel check
					cleanWord := nextWord
					for strings.HasSuffix(cleanWord, ".") || strings.HasSuffix(cleanWord, ",") || strings.HasSuffix(cleanWord, "!") || strings.HasSuffix(cleanWord, "?") || strings.HasSuffix(cleanWord, ";") || strings.HasSuffix(cleanWord, ":") {
						cleanWord = cleanWord[:len(cleanWord)-1]
					}

					if len(cleanWord) > 0 {
						first := strings.ToLower(cleanWord)[0]
						if first == 'a' || first == 'e' || first == 'i' || first == 'o' || first == 'u' || first == 'h' {
							// Should be "an"
							switch words[i] {
							case "a":
								words[i] = "an"
							case "A":
								words[i] = "An" // From (cap) command

							case "UP_A":
								words[i] = "AN" // From (up) command
							case "AN":
								words[i] = "AN" // Already fully uppercase
							case "UP_AN":
								words[i] = "AN"
							case "UP_An":
								words[i] = "An"
							}
						} else {
							// Should be "a"
							switch words[i] {
							case "an":
								words[i] = "a"
							case "An":
								words[i] = "A"

							case "UP_A":
								words[i] = "A" // From (up) command
							case "AN":
								words[i] = "A" // Preserve uppercase from (up) command
							case "UP_AN":
								words[i] = "AN"
							case "UP_An":
								words[i] = "An"
							}
						}
					}
				}
			}
		}
		lines[lineIdx] = strings.Join(words, " ")
	}
	return strings.Join(lines, "\n")
}

// --------------- helper functions ---------------

// creates a new TokenProcessor with preallocated token buffer
func NewTokenProcessor() *TokenProcessor {
	tokenBufferSize := config.OVERLAP_WORDS * 4 // 4x OVERLAP_WORDS for consistency
	return &TokenProcessor{
		tokens: make([]Token, tokenBufferSize),
	}
}

// validates command syntax before processing to prevent invalid transformations
func (tp *TokenProcessor) isValidCommand(cmdValue string) bool {
	// Check for valid single commands
	switch cmdValue {
	case "hex", "bin", "up", "low", "cap":
		return true
	}

	// Check for valid multi-word commands.
	if strings.Contains(cmdValue, ",") {
		parts := strings.Split(cmdValue, ",")
		if len(parts) == 2 {
			cmd := strings.TrimSpace(parts[0])
			countStr := strings.TrimSpace(parts[1])
			if cmd == "up" || cmd == "low" || cmd == "cap" {
				if _, err := strconv.Atoi(countStr); err == nil {
					return true
				}
			}
		}
	}

	return false
}

// applies case transformations to individual words.
func (tp *TokenProcessor) transformWord(word, cmd string) string {
	switch cmd {
	case "up":
		return strings.ToUpper(word)
	case "low":
		return strings.ToLower(word)
	case "cap":
		if len(word) == 0 {
			return word
		}
		lower := strings.ToLower(word)
		return strings.ToUpper(string(lower[0])) + lower[1:]
	}
	return word
}

// writes remaining tokens to output buffer with proper spacing and resets token buffer
func (tp *TokenProcessor) flushTokens() {
	for i := 0; i < tp.tokenIdx; i++ {
		token := tp.tokens[i]
		switch token.Type {
		case WORD:
			if tp.output.Len() > 0 && !strings.HasSuffix(tp.output.String(), " ") && !strings.HasSuffix(tp.output.String(), "\n") {
				tp.output.WriteByte(' ')
			}
			tp.output.WriteString(token.Value)
		case PUNCTUATION:
			// Handle spacing for punctuation
			if token.Value == "(" {
				// Opening parenthesis preserves space before it
				tp.output.WriteString(token.Value)
			} else {
				// Other punctuation - remove trailing space
				result := tp.output.String()
				if strings.HasSuffix(result, " ") {
					tp.output.Reset()
					tp.output.WriteString(result[:len(result)-1])
				}
				tp.output.WriteString(token.Value)
			}
		case SPACE:
			if tp.output.Len() > 0 && !strings.HasSuffix(tp.output.String(), " ") && !strings.HasSuffix(tp.output.String(), "\n") {
				tp.output.WriteByte(' ')
			}
		case NEWLINE:
			tp.output.WriteByte('\n')
		}
	}
	tp.tokenIdx = 0
}
