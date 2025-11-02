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
	tokens   [50]Token
	tokenIdx int
	output   strings.Builder
}

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
				// Remove trailing space before punctuation
				result := tp.output.String()
				if strings.HasSuffix(result, " ") {
					tp.output.Reset()
					tp.output.WriteString(result[:len(result)-1])
				}
				tp.output.WriteString(token.Value)
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
				// Transform multiple words
				wordsTransformed := 0
				for i := tp.tokenIdx - 1; i >= 0 && wordsTransformed < count; i-- {
					if tp.tokens[i].Type == WORD {
						tp.tokens[i].Value = tp.transformWord(tp.tokens[i].Value, cmd)
						wordsTransformed++
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
			tp.tokens[lastWordIdx].Value = tp.transformWord(tp.tokens[lastWordIdx].Value, cmdValue)
		}
	}
}

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
			// Remove trailing space before punctuation
			result := tp.output.String()
			if strings.HasSuffix(result, " ") {
				tp.output.Reset()
				tp.output.WriteString(result[:len(result)-1])
			}
			tp.output.WriteString(token.Value)
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

// ProcessText - Single pass dual FSM implementation
func ProcessText(text string) string {
	if text == "" {
		return ""
	}

	runes := []rune(text)
	processor := &TokenProcessor{}

	state := STATE_TEXT
	var wordBuilder strings.Builder
	var cmdBuilder strings.Builder

	for i := 0; i < len(runes); i++ {
		r := runes[i]

		switch state {
		case STATE_TEXT:
			switch r {
			case '(':
				// Flush current word
				if wordBuilder.Len() > 0 {
					processor.addToken(Token{WORD, wordBuilder.String()})
					wordBuilder.Reset()
				}
				state = STATE_COMMAND
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
				// Process command
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

	// Post-process articles
	result := processor.output.String()
	return fixArticles(result)
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
			case "a", "A", "an", "An":
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
								words[i] = "An"
							}
						} else {
							// Should be "a"
							switch words[i] {
							case "an":
								words[i] = "a"
							case "An":
								words[i] = "A"
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

// Legacy compatibility functions
func TokenizeText(text string) []string {
	return strings.Fields(ProcessText(text))
}

func ApplyAllTransformations(tokens []string) []string {
	return strings.Fields(ProcessText(strings.Join(tokens, " ")))
}

func ApplyAllTransformationsWithContext(currentChunk, overlapContext string) []string {
	mergedText := overlapContext
	if overlapContext != "" && currentChunk != "" {
		mergedText += " " + currentChunk
	} else if currentChunk != "" {
		mergedText = currentChunk
	}
	return strings.Fields(ProcessText(mergedText))
}

func TokensToString(tokens []string) string {
	return strings.Join(tokens, " ")
}

// StreamProcessor maintains FSM state across chunks
type StreamProcessor struct {
	buffer    strings.Builder
	processor *TokenProcessor
	output    strings.Builder
}

func NewStreamProcessor() *StreamProcessor {
	return &StreamProcessor{
		processor: &TokenProcessor{},
	}
}

func (sp *StreamProcessor) ProcessChunk(data []byte) string {
	// Add chunk to buffer
	sp.buffer.Write(data)
	bufferText := sp.buffer.String()

	// If buffer is getting large, process it in segments
	if sp.buffer.Len() > config.CHUNK_BYTES*2 {
		// Find last complete sentence or word boundary
		lastBoundary := strings.LastIndexAny(bufferText, ".!?")
		if lastBoundary == -1 {
			lastBoundary = strings.LastIndex(bufferText, " ")
		}

		if lastBoundary > 0 {
			completeText := bufferText[:lastBoundary+1]
			remaining := bufferText[lastBoundary+1:]

			// Reset buffer with remaining text
			sp.buffer.Reset()
			sp.buffer.WriteString(remaining)

			// Process and accumulate output
			processed := ProcessText(completeText)
			sp.output.WriteString(processed)

			// Return accumulated output
			result := sp.output.String()
			sp.output.Reset()
			return result
		}
	}

	return ""
}

func (sp *StreamProcessor) Flush() string {
	// Process any remaining text in buffer
	remaining := sp.buffer.String()
	sp.buffer.Reset()

	// Add any accumulated output
	accumulated := sp.output.String()
	sp.output.Reset()

	if remaining != "" {
		processed := ProcessText(remaining)
		if accumulated != "" {
			return accumulated + processed
		}
		return processed
	}

	return accumulated
}
