package transformer

import (
	"reflect"
	"testing"
)

func TestTokenizeTextBasic(t *testing.T) {
	text := "hello world"
	tokens := TokenizeText(text)
	
	expected := []Token{
		{Type: Word, Value: "hello"},
		{Type: Word, Value: "world"},
	}
	
	if !reflect.DeepEqual(tokens, expected) {
		t.Errorf("Expected %+v, got %+v", expected, tokens)
	}
}

func TestTokenizeTextWithCommands(t *testing.T) {
	text := "hello (up) world (hex)"
	tokens := TokenizeText(text)
	
	expected := []Token{
		{Type: Word, Value: "hello"},
		{Type: Command, Value: "(up)"},
		{Type: Word, Value: "world"},
		{Type: Command, Value: "(hex)"},
	}
	
	if !reflect.DeepEqual(tokens, expected) {
		t.Errorf("Expected %+v, got %+v", expected, tokens)
	}
}

func TestTokenizeTextWithPunctuation(t *testing.T) {
	text := "hello, world! How are you?"
	tokens := TokenizeText(text)
	
	expected := []Token{
		{Type: Word, Value: "hello"},
		{Type: Punctuation, Value: ","},
		{Type: Word, Value: "world"},
		{Type: Punctuation, Value: "!"},
		{Type: Word, Value: "How"},
		{Type: Word, Value: "are"},
		{Type: Word, Value: "you"},
		{Type: Punctuation, Value: "?"},
	}
	
	if !reflect.DeepEqual(tokens, expected) {
		t.Errorf("Expected %+v, got %+v", expected, tokens)
	}
}

func TestTokenizeTextWithLineBreaks(t *testing.T) {
	text := "first line\nsecond line\n\nthird line"
	tokens := TokenizeText(text)
	
	expected := []Token{
		{Type: Word, Value: "first"},
		{Type: Word, Value: "line"},
		{Type: LineBreak, Value: "\n"},
		{Type: Word, Value: "second"},
		{Type: Word, Value: "line"},
		{Type: LineBreak, Value: "\n"},
		{Type: LineBreak, Value: "\n"},
		{Type: Word, Value: "third"},
		{Type: Word, Value: "line"},
	}
	
	if !reflect.DeepEqual(tokens, expected) {
		t.Errorf("Expected %+v, got %+v", expected, tokens)
	}
}

func TestTokenizeTextWithContractions(t *testing.T) {
	text := "Let's test don't and won't contractions"
	tokens := TokenizeText(text)
	
	expected := []Token{
		{Type: Word, Value: "Let's"},
		{Type: Word, Value: "test"},
		{Type: Word, Value: "don't"},
		{Type: Word, Value: "and"},
		{Type: Word, Value: "won't"},
		{Type: Word, Value: "contractions"},
	}
	
	if !reflect.DeepEqual(tokens, expected) {
		t.Errorf("Expected %+v, got %+v", expected, tokens)
	}
}

func TestConvertHexBasic(t *testing.T) {
	tokens := []Token{
		{Type: Word, Value: "1E"},
		{Type: Command, Value: "(hex)"},
	}
	
	result := ConvertHex(tokens)
	
	expected := []Token{
		{Type: Word, Value: "30"},
	}
	
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %+v, got %+v", expected, result)
	}
}

func TestConvertBinaryBasic(t *testing.T) {
	tokens := []Token{
		{Type: Word, Value: "1010"},
		{Type: Command, Value: "(bin)"},
	}
	
	result := ConvertBinary(tokens)
	
	expected := []Token{
		{Type: Word, Value: "10"},
	}
	
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %+v, got %+v", expected, result)
	}
}

func TestApplyCaseTransformUp(t *testing.T) {
	tokens := []Token{
		{Type: Word, Value: "hello"},
		{Type: Command, Value: "(up)"},
	}
	
	result := ApplyCaseTransform(tokens)
	
	expected := []Token{
		{Type: Word, Value: "HELLO"},
	}
	
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %+v, got %+v", expected, result)
	}
}

func TestApplyCaseTransformMultipleWords(t *testing.T) {
	tokens := []Token{
		{Type: Word, Value: "hello"},
		{Type: Word, Value: "world"},
		{Type: Command, Value: "(up, 2)"},
	}
	
	result := ApplyCaseTransform(tokens)
	
	expected := []Token{
		{Type: Word, Value: "HELLO"},
		{Type: Word, Value: "WORLD"},
	}
	
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %+v, got %+v", expected, result)
	}
}

func TestFixPunctuationSpacingBasic(t *testing.T) {
	tokens := []Token{
		{Type: Word, Value: "hello"},
		{Type: Punctuation, Value: ","},
		{Type: Word, Value: "world"},
	}
	
	result := FixPunctuationSpacing(tokens)
	
	expected := []Token{
		{Type: Word, Value: "hello,"},
		{Type: Word, Value: "world"},
	}
	
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %+v, got %+v", expected, result)
	}
}

func TestFixPunctuationSpacingMultiple(t *testing.T) {
	tokens := []Token{
		{Type: Word, Value: "word"},
		{Type: Punctuation, Value: "!"},
		{Type: Punctuation, Value: "!"},
	}
	
	result := FixPunctuationSpacing(tokens)
	
	expected := []Token{
		{Type: Word, Value: "word!!"},
	}
	
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %+v, got %+v", expected, result)
	}
}

func TestFixPunctuationSpacingAllTypes(t *testing.T) {
	tokens := []Token{
		{Type: Word, Value: "Hello"},
		{Type: Punctuation, Value: ","},
		{Type: Word, Value: "world"},
		{Type: Punctuation, Value: "!"},
		{Type: Word, Value: "How"},
		{Type: Punctuation, Value: "?"},
		{Type: Word, Value: "Fine"},
		{Type: Punctuation, Value: ";"},
		{Type: Word, Value: "thanks"},
		{Type: Punctuation, Value: "."},
	}
	
	result := FixPunctuationSpacing(tokens)
	
	expected := []Token{
		{Type: Word, Value: "Hello,"},
		{Type: Word, Value: "world!"},
		{Type: Word, Value: "How?"},
		{Type: Word, Value: "Fine;"},
		{Type: Word, Value: "thanks."},
	}
	
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %+v, got %+v", expected, tokens)
	}
}

func TestRepositionQuotesBasic(t *testing.T) {
	tokens := []Token{
		{Type: Word, Value: "He"},
		{Type: Word, Value: "said"},
		{Type: Quote, Value: "'"},
		{Type: Word, Value: "hello"},
		{Type: Word, Value: "world"},
		{Type: Quote, Value: "'"},
		{Type: Word, Value: "today"},
	}
	
	result := RepositionQuotes(tokens)
	
	expected := []Token{
		{Type: Word, Value: "He"},
		{Type: Word, Value: "said"},
		{Type: Word, Value: "'hello"},
		{Type: Word, Value: "world'"},
		{Type: Word, Value: "today"},
	}
	
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %+v, got %+v", expected, result)
	}
}

func TestRepositionQuotesSingleWord(t *testing.T) {
	tokens := []Token{
		{Type: Word, Value: "word"},
		{Type: Quote, Value: "'"},
		{Type: Word, Value: "quoted"},
		{Type: Quote, Value: "'"},
		{Type: Word, Value: "end"},
	}
	
	result := RepositionQuotes(tokens)
	
	expected := []Token{
		{Type: Word, Value: "word"},
		{Type: Word, Value: "'quoted'"},
		{Type: Word, Value: "end"},
	}
	
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %+v, got %+v", expected, result)
	}
}

func TestCorrectArticlesVowels(t *testing.T) {
	tokens := []Token{
		{Type: Word, Value: "It"},
		{Type: Word, Value: "was"},
		{Type: Word, Value: "a"},
		{Type: Word, Value: "amazing"},
		{Type: Word, Value: "day"},
	}
	
	result := CorrectArticles(tokens)
	
	expected := []Token{
		{Type: Word, Value: "It"},
		{Type: Word, Value: "was"},
		{Type: Word, Value: "an"},
		{Type: Word, Value: "amazing"},
		{Type: Word, Value: "day"},
	}
	
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %+v, got %+v", expected, result)
	}
}

func TestCorrectArticlesHonest(t *testing.T) {
	tokens := []Token{
		{Type: Word, Value: "a"},
		{Type: Word, Value: "honest"},
		{Type: Word, Value: "person"},
	}
	
	result := CorrectArticles(tokens)
	
	expected := []Token{
		{Type: Word, Value: "an"},
		{Type: Word, Value: "honest"},
		{Type: Word, Value: "person"},
	}
	
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %+v, got %+v", expected, result)
	}
}

func TestCorrectArticlesNoChange(t *testing.T) {
	tokens := []Token{
		{Type: Word, Value: "a"},
		{Type: Word, Value: "car"},
	}
	
	result := CorrectArticles(tokens)
	
	expected := []Token{
		{Type: Word, Value: "a"},
		{Type: Word, Value: "car"},
	}
	
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %+v, got %+v", expected, result)
	}
}

func TestApplyAllTransformationsChaining(t *testing.T) {
	tokens := []Token{
		{Type: Word, Value: "1010"},
		{Type: Command, Value: "(bin)"},
		{Type: Command, Value: "(hex)"},
	}
	
	result := ApplyAllTransformations(tokens)
	
	expected := []Token{
		{Type: Word, Value: "16"},
	}
	
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %+v, got %+v", expected, result)
	}
}

func TestApplyAllTransformationsComplete(t *testing.T) {
	tokens := []Token{
		{Type: Word, Value: "hello"},
		{Type: Punctuation, Value: ","},
		{Type: Word, Value: "world"},
		{Type: Command, Value: "(up)"},
		{Type: Punctuation, Value: "!"},
	}
	
	result := ApplyAllTransformations(tokens)
	
	expected := []Token{
		{Type: Word, Value: "hello,"},
		{Type: Word, Value: "WORLD!"},
	}
	
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %+v, got %+v", expected, result)
	}
}

func TestApplyAllTransformationsInvalidCommands(t *testing.T) {
	tokens := []Token{
		{Type: Word, Value: "This"},
		{Type: Command, Value: "(invalid)"},
		{Type: Word, Value: "and"},
		{Type: Command, Value: "(up, text)"},
		{Type: Word, Value: "should"},
		{Type: Word, Value: "remain"},
		{Type: Word, Value: "unchanged"},
	}
	
	result := ApplyAllTransformations(tokens)
	
	// Invalid commands should remain as tokens
	expected := []Token{
		{Type: Word, Value: "This"},
		{Type: Command, Value: "(invalid)"},
		{Type: Word, Value: "and"},
		{Type: Command, Value: "(up, text)"},
		{Type: Word, Value: "should"},
		{Type: Word, Value: "remain"},
		{Type: Word, Value: "unchanged"},
	}
	
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %+v, got %+v", expected, result)
	}
}

func TestApplyAllTransformationsWithContext(t *testing.T) {
	// Simulate cross-chunk scenario where command needs words from previous chunk
	overlapContext := "word1 word2 word3"
	currentChunk := "word4 word5 (up, 5) remaining"
	
	result := ApplyAllTransformationsWithContext(currentChunk, overlapContext)
	
	// Should transform all 5 words: 3 from context + 2 from current
	expected := []Token{
		{Type: Word, Value: "WORD1"},
		{Type: Word, Value: "WORD2"},
		{Type: Word, Value: "WORD3"},
		{Type: Word, Value: "WORD4"},
		{Type: Word, Value: "WORD5"},
		{Type: Word, Value: "remaining"},
	}
	
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %+v, got %+v", expected, result)
	}
}

func TestPreserveLineEndingsBasic(t *testing.T) {
	text := "First line with transformation (up).\nSecond line here."
	tokens := TokenizeText(text)
	result := ApplyAllTransformations(tokens)
	output := TokensToString(result)
	
	expected := "First line with TRANSFORMATION.\nSecond line here."
	
	if output != expected {
		t.Errorf("Expected %q, got %q", expected, output)
	}
}

func TestPreserveLineEndingsMultiple(t *testing.T) {
	text := "First line.\n\nThird line after blank.\nFinal line with A (hex)."
	tokens := TokenizeText(text)
	result := ApplyAllTransformations(tokens)
	output := TokensToString(result)
	
	expected := "First line.\n\nThird line after blank.\nFinal line with 10."
	
	if output != expected {
		t.Errorf("Expected %q, got %q", expected, output)
	}
}

func TestContractionsPreserved(t *testing.T) {
	text := "Let's test don't and won't contractions"
	tokens := TokenizeText(text)
	result := ApplyAllTransformations(tokens)
	output := TokensToString(result)
	
	expected := "Let's test don't and won't contractions"
	
	if output != expected {
		t.Errorf("Expected %q, got %q", expected, output)
	}
}