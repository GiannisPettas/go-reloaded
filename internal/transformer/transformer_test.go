package transformer

import (
	"testing"
)

// Purpose: Tests constants during development/CI

func TestProcessTextBasic(t *testing.T) {
	text := "hello world"
	result := ProcessText(text)
	expected := "hello world"

	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestProcessTextHex(t *testing.T) {
	text := "FF (hex) equals 255"
	result := ProcessText(text)
	expected := "255 equals 255"

	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestProcessTextBinary(t *testing.T) {
	text := "1010 (bin) equals 10"
	result := ProcessText(text)
	expected := "10 equals 10"

	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestProcessTextCaseUp(t *testing.T) {
	text := "hello (up) world"
	result := ProcessText(text)
	expected := "HELLO world"

	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestProcessTextCaseLow(t *testing.T) {
	text := "HELLO (low) world"
	result := ProcessText(text)
	expected := "hello world"

	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestProcessTextCaseCap(t *testing.T) {
	text := "hello (cap) world"
	result := ProcessText(text)
	expected := "Hello world"

	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestProcessTextMultiWord(t *testing.T) {
	text := "these three words (up, 3) test"
	result := ProcessText(text)
	expected := "THESE THREE WORDS test"

	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestProcessTextArticles(t *testing.T) {
	text := "I need a apple and an car"
	result := ProcessText(text)
	expected := "I need an apple and a car"

	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestProcessTextPunctuation(t *testing.T) {
	text := "Hello , world ! How are you ?"
	result := ProcessText(text)
	expected := "Hello, world! How are you?"

	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestProcessTextLineBreaks(t *testing.T) {
	text := "first line\nsecond line"
	result := ProcessText(text)
	expected := "first line\nsecond line"

	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestProcessTextChaining(t *testing.T) {
	text := "1010 (bin) (hex) result"
	result := ProcessText(text)
	expected := "16 result" // 1010 bin->10, 10 hex->16

	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestProcessTextEmpty(t *testing.T) {
	text := ""
	result := ProcessText(text)
	expected := ""

	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}
