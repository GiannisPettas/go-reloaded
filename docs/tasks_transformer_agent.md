# Transformer Agent Tasks — Go Reloaded

Tasks for the Transformer Agent focused on text transformations and FSM implementation.

---

## Task 6: Transformer - Word Tokenization
**Functionality**: Split text into words and identify transformation markers  
**TDD Steps**:
1. **Red**: Write tests for word splitting and marker detection
   - Test splitting text into words with various separators
   - Test identifying transformation markers (hex), (bin), (up), etc.
   - Test handling punctuation and special characters
2. **Green**: Implement tokenizer that identifies words and commands
   - Create `internal/transformer/transformer.go`
   - Implement `TokenizeText(text string) []Token`
   - Define Token struct with Type (Word, Command, Punctuation)
3. **Refactor**: Optimize tokenization performance
4. **Validate**: Correct parsing of text with various markers

**Dependencies**: Config Agent (Task 1)

---

## Task 7: Transformer - Hexadecimal Conversion
**Functionality**: Convert hex numbers to decimal using (hex) marker  
**TDD Steps**:
1. **Red**: Write tests for hex conversion including edge cases
   - Test basic hex conversion: "1E (hex)" → "30"
   - Test zero value: "0 (hex)" → "0"
   - Test large values: "FF (hex)" → "255"
   - Test negative hex: "-1A (hex)" → "-26"
   - Test invalid hex values
2. **Green**: Implement hex-to-decimal transformation
   - Add `ConvertHex(tokens []Token) []Token`
   - Handle negative values by preserving sign
   - Skip invalid hex values
3. **Refactor**: Optimize hex conversion logic
4. **Validate**: Conversion accuracy with golden test T6, T17, T21

**Dependencies**: Task 6

---

## Task 8: Transformer - Binary Conversion
**Functionality**: Convert binary numbers to decimal using (bin) marker  
**TDD Steps**:
1. **Red**: Write tests for binary conversion including edge cases
   - Test basic binary: "1010 (bin)" → "10"
   - Test zero: "0 (bin)" → "0"
   - Test large binary: "11111111 (bin)" → "255"
   - Test negative binary: "-101 (bin)" → "-5"
   - Test invalid binary values
2. **Green**: Implement binary-to-decimal transformation
   - Add `ConvertBinary(tokens []Token) []Token`
   - Handle negative values by preserving sign
   - Skip invalid binary values
3. **Refactor**: Optimize binary conversion logic
4. **Validate**: Conversion accuracy with golden test T7, T21

**Dependencies**: Task 6

---

## Task 9: Transformer - Case Transformations (Single Word)
**Functionality**: Apply (up), (low), (cap) to single words  
**TDD Steps**:
1. **Red**: Write tests for single word case transformations
   - Test "word (up)" → "WORD"
   - Test "WORD (low)" → "word"
   - Test "word (cap)" → "Word"
   - Test edge cases with punctuation and numbers
2. **Green**: Implement uppercase, lowercase, capitalize functions
   - Add `ApplyCaseTransform(tokens []Token) []Token`
   - Handle (up), (low), (cap) commands
3. **Refactor**: Optimize case transformation logic
4. **Validate**: Transformations with golden test T18

**Dependencies**: Task 6

---

## Task 10: Transformer - Case Transformations (Multiple Words)
**Functionality**: Apply (up, n), (low, n), (cap, n) to multiple words  
**TDD Steps**:
1. **Red**: Write tests for multi-word case transformations
   - Test "word1 word2 (up, 2)" → "WORD1 WORD2"
   - Test boundary cases where n > available words
   - Test overlapping commands
2. **Green**: Implement word count handling and range application
   - Extend `ApplyCaseTransform` for numbered commands
   - Handle word counting and range validation
3. **Refactor**: Optimize multi-word processing
4. **Validate**: With golden tests T1, T10, T19, T22

**Dependencies**: Task 9

---

## Task 11: Transformer - Punctuation Spacing
**Functionality**: Fix spacing around punctuation marks  
**TDD Steps**:
1. **Red**: Write tests for punctuation spacing rules
   - Test "word ,punct" → "word, punct"
   - Test "word !" → "word!"
   - Test multiple punctuation "word !!" → "word!!"
2. **Green**: Implement punctuation detection and spacing correction
   - Add `FixPunctuationSpacing(tokens []Token) []Token`
   - Handle commas, periods, exclamation marks, question marks
3. **Refactor**: Optimize punctuation processing
4. **Validate**: With golden tests T2, T5, T13, T20

**Dependencies**: Task 6

---

## Task 12: Transformer - Quote Repositioning
**Functionality**: Move single quotes to correct positions  
**TDD Steps**:
1. **Red**: Write tests for quote pair detection and repositioning
   - Test "word ' quote ' word" → "word 'quote' word"
   - Test multiple quote pairs
   - Test unmatched quotes
2. **Green**: Implement quote movement logic
   - Add `RepositionQuotes(tokens []Token) []Token`
   - Handle quote pair matching and repositioning
3. **Refactor**: Optimize quote processing
4. **Validate**: With golden tests T9, T16

**Dependencies**: Task 6

---

## Task 13: Transformer - Article Correction
**Functionality**: Change "a" to "an" before vowels and "h"  
**TDD Steps**:
1. **Red**: Write tests for article correction rules
   - Test "a apple" → "an apple"
   - Test "a honest" → "an honest"
   - Test "a car" → "a car" (no change)
2. **Green**: Implement vowel/h detection and article replacement
   - Add `CorrectArticles(tokens []Token) []Token`
   - Handle vowel and 'h' detection
3. **Refactor**: Optimize article correction
4. **Validate**: With golden test T8

**Dependencies**: Task 6

---

## Task 14: Transformer - Command Chaining
**Functionality**: Handle multiple commands on same word  
**TDD Steps**:
1. **Red**: Write tests for command chaining scenarios
   - Test "1010 (bin) (hex)" → "16"
   - Test left-to-right execution order
2. **Green**: Implement left-to-right command execution
   - Modify transformation pipeline for chaining
   - Ensure proper order of operations
3. **Refactor**: Optimize command chaining
4. **Validate**: With golden test T3

**Dependencies**: Tasks 7, 8

---

## Task 15: Transformer - Invalid Command Handling
**Functionality**: Ignore malformed commands gracefully  
**TDD Steps**:
1. **Red**: Write tests for various invalid command formats
   - Test "(invalid)" commands
   - Test incomplete commands "(up,"
   - Test malformed syntax
2. **Green**: Implement error-tolerant command parsing
   - Add validation to command parsing
   - Skip invalid commands without errors
3. **Refactor**: Improve error handling
4. **Validate**: With golden tests T4, T14

**Dependencies**: Task 6

---

## Task 17: Transformer - Cross-Chunk Context
**Functionality**: Apply transformations using words from previous chunks  
**TDD Steps**:
1. **Red**: Write tests for commands referencing previous chunk words
   - Test commands that need more words than available in current chunk
   - Test context preservation across chunk boundaries
2. **Green**: Implement context merging in transformation logic
   - Modify transformer to accept overlap context
   - Handle cross-chunk word references
3. **Refactor**: Optimize context handling
4. **Validate**: With golden test T22

**Dependencies**: All previous transformer tasks, Parser Agent Task 16

---

## Task 23: Transformer - Preserve Line Endings
**Functionality**: Maintain original line breaks and paragraph structure  
**TDD Steps**:
1. **Red**: Write tests for line ending preservation
   - Test single line breaks are preserved
   - Test multiple consecutive line breaks (paragraphs)
   - Test line breaks with transformations
2. **Green**: Implement line break detection and preservation
   - Modify tokenizer to recognize and preserve \n characters
   - Ensure transformations don't remove line breaks
   - Handle line breaks in token stream
3. **Refactor**: Optimize line break handling
4. **Validate**: With golden test T23

**Dependencies**: Task 6 (Word Tokenization)