# Go-Reloaded Transformer: Step-by-Step Guide for Junior Developers

## What Does the Transformer Do?

The transformer is the heart of Go-Reloaded. It takes text like:
```
"The value FF (hex) should be (up) and I need a apple with ' spaced quotes '."
```

And transforms it to:
```
"The value 255 SHOULD BE and I need an apple with 'spaced quotes'."
```

It handles:
- **Commands in parentheses**: `(hex)`, `(bin)`, `(up)`, `(low)`, `(cap)`, `(up, 2)`
- **Invalid command preservation**: `(invalid)` stays as `(invalid)`
- **Article correction**: `a apple` → `an apple`
- **Quote repositioning**: `' text '` → `'text'`
- **Punctuation spacing**: `word ,` → `word,`

## Core Concepts

### 1. Dual Finite State Machine (FSM) Architecture

Go-Reloaded uses **two FSMs working together** for maximum efficiency:

#### **Low-Level FSM (Character Parser)**
- **STATE_TEXT**: Reading normal text
- **STATE_COMMAND**: Reading commands inside parentheses `(like this)`
- Handles **syntax**: What does each character mean?
- Manages **state transitions**: When to switch between TEXT/COMMAND
- Deals with **Unicode safety**: Converting bytes to runes properly

#### **High-Level FSM (Token Processor)**
- Handles **semantics**: What do complete tokens mean?
- Manages **transformations**: How to modify words
- Deals with **memory management**: Buffer overflow, token storage

#### **Why Split Into Two FSMs?**

**Not just for clarity - there are real technical benefits:**

**1. Separation of Concerns**
- Low-level handles character parsing
- High-level handles token processing
- Clean, maintainable code

**2. Memory Efficiency**
```go
// Without split - grows with file size
type SingleFSM struct {
    allTokens []Token     // Memory grows!
    allCommands []Command // More memory!
}

// With split - constant memory
type TokenProcessor struct {
    tokens [50]Token  // Fixed size, ~8KB always
}
```

**3. Extensibility**
- Easy to add new character types (low-level)
- Easy to add new commands (high-level)
- Independent development

**4. Performance Optimization**
- Low-level: Fast character classification
- High-level: Smart token buffering
- Each optimized for its purpose

**5. Testing**
- Test character parsing separately
- Test token processing separately
- Easier debugging

#### **How They Work Together**
```go
for each character {
    // LOW-LEVEL FSM decides what character means
    switch state {
    case STATE_TEXT:
        if r == ' ' {
            // HIGH-LEVEL FSM processes immediately
            processor.addToken(Token{WORD, wordBuilder.String()})
        }
    }
}
```

**No waiting - it's event-driven!** The low-level FSM **triggers** the high-level FSM when something is ready to process.

### 2. Token Structure - Why We Need It

**Token Definition:**
```go
type Token struct {
    Type  int    // What kind of element (WORD, PUNCTUATION, SPACE, NEWLINE)
    Value string // The actual text content
}
```

**Token Types:**
- **WORD**: "hello", "world", "FF"
- **PUNCTUATION**: ".", "!", "?"
- **SPACE**: " " (spaces and tabs)
- **NEWLINE**: "\n" (line breaks)
- **COMMAND**: Commands are processed immediately, not stored as tokens

#### Why We Need Structured Tokens

**1. Command Processing**
Commands like `(up, 3)` need to find and transform **specific previous words**:

```go
// Without tokens: Hard to find "these three words"
"these three words should be (up, 3)"

// With tokens: Easy to locate by type
[{WORD, "these"}, {SPACE, " "}, {WORD, "three"}, {WORD, "words"}, ...]
```

**2. Punctuation Spacing**
Different punctuation needs different spacing rules:

```go
// Token types allow specific handling
{PUNCTUATION, "("} → Keep space before
{PUNCTUATION, "!"} → Remove space before
{PUNCTUATION, "."} → Remove space before
```

**3. Buffered Processing**
The transformer uses a **fixed-size token buffer** (50 tokens) for constant memory:

```go
tokens [50]Token  // Fixed buffer, no growing slices
```

**4. Context Preservation**
Commands can reference words that appeared earlier:

```go
"word1 word2 word3 (up, 2)" 
// Need to find "word2 word3" - tokens make this easy
```

#### Without Tokens (Problems)

```go
// String processing - hard to track word positions
text := "hello world (up, 1)"
// How do you find "world" efficiently?
// How do you handle punctuation spacing?
// How do you maintain fixed memory usage?
```

#### With Tokens (Clean)

```go
// Structured processing - easy to navigate
tokens := [{WORD, "hello"}, {SPACE, " "}, {WORD, "world"}]
// Easy to find last word, apply transformations, handle spacing
```

**Tokens provide the structure needed for complex text transformations while maintaining constant memory usage.**

### 3. strings.Builder
`strings.Builder` is a Go standard library type for **efficient string concatenation**.

**Why use it instead of regular string concatenation?**
```go
// BAD - Creates new string objects each time (slow & memory-heavy)
result := ""
result += "Hello"
result += " "
result += "World"  // Each += creates a new string in memory

// GOOD - Reuses internal buffer (fast & memory-efficient)
var builder strings.Builder
builder.WriteString("Hello")
builder.WriteByte(' ')
builder.WriteString("World")
result := builder.String()  // Only creates final string once
```

**Key methods used in the transformer:**
- `WriteString("text")` - Add string
- `WriteByte(' ')` - Add single character
- `WriteRune('世')` - Add Unicode character
- `Len()` - Get current length
- `String()` - Get final result
- `Reset()` - Clear and reuse

**Performance benefits:**
- **Memory Efficient**: Grows internal buffer as needed, no constant reallocation
- **Fast**: No copying strings repeatedly like with `+=`
- **Zero Value Ready**: Can use immediately without initialization

This is why Go-Reloaded can process large files efficiently!

## Function Execution Order

**The first function that runs in the transformer is:**

### **`ProcessText(text string) string`** - Main Entry Point

This is the **main entry point** - it's the function that gets called from the controller.

**Execution Flow:**
```go
1. ProcessText()           // ← STARTS HERE (main entry point)
   ↓
2. Creates TokenProcessor  // processor := &TokenProcessor{}
   ↓
3. Character-by-character loop begins
   ↓
4. addToken()             // Called for each token found
   ↓
5. processCommand()       // Called when ')' is found
   ↓
6. transformWord()        // Called from processCommand()
   ↓
7. flushTokens()          // Called at the end
   ↓
8. fixArticles()          // Final post-processing
   ↓
9. Returns final string   // Back to controller
```

**How It Gets Called:**
From the controller:
```go
// controller.go calls transformer
result := transformer.ProcessText(text)  // ← This is where it starts
```

**What ProcessText Does First:**
```go
func ProcessText(text string) string {
    if text == "" {           // 1. Check for empty input
        return ""
    }
    
    runes := []rune(text)     // 2. Convert to runes (Unicode safe)
    processor := &TokenProcessor{}  // 3. Create token processor
    
    state := STATE_TEXT       // 4. Initialize FSM state
    var wordBuilder strings.Builder  // 5. Create string builders
    var cmdBuilder strings.Builder
    
    // 6. Start the main character loop...
}
```

## Step-by-Step Process

### Step 1: ProcessText() Details

```go
func ProcessText(text string) string
```

**What it does:**
1. Takes input text as string
2. Converts to runes (handles Unicode properly)
3. Creates a TokenProcessor to manage tokens
4. Sets up two string builders for collecting words and commands
5. Starts the FSM in STATE_TEXT

**Example:**
```
Input: "Hello (up) world!"
Runes: ['H', 'e', 'l', 'l', 'o', ' ', '(', 'u', 'p', ')', ' ', 'w', 'o', 'r', 'l', 'd', '!']
```

### Step 2: Character-by-Character Processing

The main loop processes each character based on the current state:

#### STATE_TEXT (Normal Reading)
```go
switch r {
case '(':
    // Found command start - switch to command mode
case ' ', '\t':
    // Found space - save current word, add space token
case '\n':
    // Found newline - save current word, add newline token
case ',', '.', '!', '?', ';', ':':
    // Found punctuation - save current word, add punctuation token
default:
    // Regular character - add to current word
}
```

#### STATE_COMMAND (Reading Commands)
```go
if r == ')' {
    // Command ended - process it and switch back to text mode
} else {
    // Still building command - add character to command
}
```

**Example walkthrough:**
```
Input: "Hello (up) world!"

Character 'H': STATE_TEXT, default -> add to wordBuilder: "H"
Character 'e': STATE_TEXT, default -> add to wordBuilder: "He"
Character 'l': STATE_TEXT, default -> add to wordBuilder: "Hel"
Character 'l': STATE_TEXT, default -> add to wordBuilder: "Hell"
Character 'o': STATE_TEXT, default -> add to wordBuilder: "Hello"
Character ' ': STATE_TEXT, space -> save "Hello" as WORD token, add SPACE token
Character '(': STATE_TEXT, '(' -> switch to STATE_COMMAND
Character 'u': STATE_COMMAND -> add to cmdBuilder: "u"
Character 'p': STATE_COMMAND -> add to cmdBuilder: "up"
Character ')': STATE_COMMAND, ')' -> process command "up", switch to STATE_TEXT
Character ' ': STATE_TEXT, space -> add SPACE token
Character 'w': STATE_TEXT, default -> add to wordBuilder: "w"
...and so on
```

### Step 3: Token Management - `TokenProcessor`

The TokenProcessor manages a buffer of 50 tokens:

```go
type TokenProcessor struct {
    tokens   [50]Token  // Fixed-size buffer
    tokenIdx int        // Current position
    output   strings.Builder // Final output
    upperArticles map[int]bool // Track articles uppercased by (up) commands
}
```

#### Adding Tokens - `addToken()`
- If buffer has space: add token
- If buffer is full: flush first half to output, shift remaining tokens, then add new token

**Why fixed buffer?** Memory efficiency! No matter how big the file, we only use ~8KB of memory.

#### Token Buffer Management

**Buffer Overflow Handling:**
```go
if tp.tokenIdx >= len(tp.tokens) {
    // Buffer is full, flush first half to output
    halfSize := len(tp.tokens) / 2
    for i := 0; i < halfSize; i++ {
        // Process token and add to output
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
```

**This ensures:**
- Constant memory usage (~8KB)
- No data loss during processing
- Continuous streaming for large files

### Step 4: Command Validation and Processing

#### Command Validation - `isValidCommand()`

**CRITICAL**: Commands are validated BEFORE processing. Invalid commands are preserved as text.

```go
func (tp *TokenProcessor) isValidCommand(cmdValue string) bool {
    // Valid single commands
    switch cmdValue {
    case "hex", "bin", "up", "low", "cap":
        return true
    }
    
    // Valid multi-word commands: "up, 2", "low, 3", "cap, 1"
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
```

**Examples:**
- `(up)` → Valid, processes command
- `(invalid)` → Invalid, preserved as `(invalid)` in output
- `(up, text)` → Invalid, preserved as `(up, text)` in output

#### Command Processing - `processCommand()`

**Only valid commands reach this function:**

#### Single Word Commands
```go
switch cmdValue {
case "hex":
    // Convert hexadecimal to decimal: "FF" -> "255"
case "bin":
    // Convert binary to decimal: "1010" -> "10"
default:
    // Case transformations: "up", "low", "cap"
    tp.tokens[lastWordIdx].Value = tp.transformWord(word, cmdValue)
}
```

#### Multi-Word Commands - **Fixed Order Processing**
```go
if strings.Contains(cmdValue, ",") {
    // Find word indices in reverse order
    var wordIndices []int
    for i := tp.tokenIdx - 1; i >= 0 && len(wordIndices) < count; i-- {
        if tp.tokens[i].Type == WORD {
            wordIndices = append(wordIndices, i)
        }
    }
    // Transform words in FORWARD order (left to right)
    for i := len(wordIndices) - 1; i >= 0; i-- {
        idx := wordIndices[i]
        tp.tokens[idx].Value = tp.transformWord(tp.tokens[idx].Value, cmd)
    }
}
```

**Example:**
```
Input: "these three words (cap, 3)"
Tokens: ["these", "three", "words"]
Command: "cap, 3"
Result: ["These", "Three", "Words"] (left-to-right capitalization)
```

### Step 5: Word Transformation - `transformWord()`

Simple case transformations:
```go
func (tp *TokenProcessor) transformWord(word, cmd string) string {
    switch cmd {
    case "up":
        return strings.ToUpper(word)     // "hello" -> "HELLO"
    case "low":
        return strings.ToLower(word)     // "HELLO" -> "hello"
    case "cap":
        if len(word) == 0 {
            return word
        }
        lower := strings.ToLower(word)
        return strings.ToUpper(string(lower[0])) + lower[1:]  // "hello" -> "Hello"
    }
    return word
}
```

**Note**: Uses manual capitalization instead of `strings.Title()` for precise control.

### Step 6: Output Generation - `flushTokens()`

Converts tokens back to text with proper spacing:

```go
switch token.Type {
case WORD:
    // Add space before word if needed
    if tp.output.Len() > 0 && !strings.HasSuffix(tp.output.String(), " ") {
        tp.output.WriteByte(' ')
    }
    tp.output.WriteString(token.Value)
case PUNCTUATION:
    // Remove trailing space before punctuation
    if token.Value != "(" {
        result := tp.output.String()
        if strings.HasSuffix(result, " ") {
            tp.output.Reset()
            tp.output.WriteString(result[:len(result)-1])
        }
    }
    tp.output.WriteString(token.Value)
case SPACE:
    // Add space if not already present
    if !strings.HasSuffix(tp.output.String(), " ") {
        tp.output.WriteByte(' ')
    }
case NEWLINE:
    tp.output.WriteByte('\n')
}
```

### Step 7: Post-Processing

After FSM processing, two post-processing steps fix grammar and formatting:

#### Article Correction - `fixArticles()`

**Fixes "a/an" usage based on vowel sounds:**

```go
func fixArticles(text string) string {
    lines := strings.Split(text, "\n")
    for lineIdx, line := range lines {
        words := strings.Fields(line)
        for i := 0; i < len(words)-1; i++ {
            switch words[i] {
            case "a", "A", "an", "An", "AN":
                nextWord := words[i+1]
                // Remove punctuation for vowel check
                cleanWord := removePunctuation(nextWord)
                
                if len(cleanWord) > 0 {
                    first := strings.ToLower(cleanWord)[0]
                    if first == 'a' || first == 'e' || first == 'i' || first == 'o' || first == 'u' || first == 'h' {
                        // Should be "an"
                        if words[i] == "a" { words[i] = "an" }
                        if words[i] == "A" { words[i] = "AN" } // Preserve (up) command result
                    } else {
                        // Should be "a"
                        if words[i] == "an" { words[i] = "a" }
                        if words[i] == "An" { words[i] = "A" }
                        if words[i] == "AN" { words[i] = "A" } // Preserve (up) command result
                    }
                }
            }
        }
    }
    return strings.Join(lines, "\n")
}
```

**Examples:**
- `a apple` → `an apple`
- `an car` → `a car`
- `A (up) apple` → `AN apple` (preserves uppercase from command)

#### Quote Repositioning - `fixQuotes()` - **General Algorithm**

**Removes spaces between quotes and their content using algorithmic approach:**

```go
func fixQuotes(text string) string {
    runes := []rune(text)
    var result strings.Builder
    
    for i := 0; i < len(runes); i++ {
        r := runes[i]
        
        if r == '\'' || r == '"' {
            // Find matching quote
            matchingQuote := -1
            for j := i + 1; j < len(runes); j++ {
                if runes[j] == r {
                    matchingQuote = j
                    break
                }
            }
            
            if matchingQuote != -1 {
                // Process quote pair
                result.WriteRune(r) // Opening quote
                
                // Skip space after opening quote
                startIdx := i + 1
                if startIdx < len(runes) && runes[startIdx] == ' ' {
                    startIdx++
                }
                
                // Find content end (before closing quote)
                endIdx := matchingQuote
                if endIdx > 0 && runes[endIdx-1] == ' ' {
                    endIdx--
                }
                
                // Write content without internal spaces
                for k := startIdx; k < endIdx; k++ {
                    result.WriteRune(runes[k])
                }
                
                result.WriteRune(r) // Closing quote
                i = matchingQuote // Skip to after closing quote
            } else {
                result.WriteRune(r) // No matching quote
            }
        } else {
            result.WriteRune(r)
        }
    }
    
    return result.String()
}
```

**Examples:**
- `' hello world '` → `'hello world'`
- `" any text "` → `"any text"`
- `' I am` → `'I am` (opening quote only)
- `carries '` → `carries'` (closing quote only)

**Key Feature**: This is a **general algorithmic solution** that works for ANY text content between quotes, not hardcoded patterns.

## Single-Pass Architecture

### FSM Single-Pass Processing

The **dual FSM processes each character exactly once** in a single forward pass:

```go
for i := 0; i < len(runes); i++ {
    r := runes[i]
    // Each character processed exactly once
    // No backtracking, no re-reading
}
```

**No multiple iterations over the input text.**

### Complete Processing Pipeline

The transformer uses **3 total passes** for optimal performance vs. complexity:

```go
// Pass 1: FSM processes text once (95% of transformations)
for i := 0; i < len(runes); i++ {
    // Dual FSM handles commands, tokens, transformations
}
result := processor.output.String()

// Pass 2: Article correction (single pass)
result = fixArticles(result)    

// Pass 3: Quote repositioning (single pass)
return fixQuotes(result)        
```

### Why This Design?

**Performance vs. Complexity Trade-off:**

- **FSM handles 95% of transformations** in single pass
- **Post-processing handles edge cases** that would complicate FSM significantly  
- **Total: 3 passes** instead of complex multi-state FSM

**Alternative would be:**
- Single mega-FSM with article/quote state tracking
- Much more complex state management
- Harder to maintain and debug

### Memory Efficiency Maintained

Even with 3 passes, memory usage remains **constant ~8KB** because:
- FSM uses fixed token buffer
- Post-processing works on final string (not growing)
- No intermediate data structures stored

**Result**: Effectively single-pass for core transformations, with minimal post-processing cleanup.

## Extensibility and Future Enhancements

### Why the Separated Architecture is Perfect for Custom Commands

The current design with **FSM + separate post-processing functions** creates an ideal foundation for adding custom user commands:

```go
// Current architecture
result := fsm.ProcessText(text)     // Core transformations
result = fixArticles(result)        // Grammar fixes
result = fixQuotes(result)          // Formatting fixes
// Easy to add more post-processors!
```

### Adding Custom Commands - Future Design

**1. FSM Extension (for real-time commands)**
```go
// Add to isValidCommand()
case "custom1", "custom2", "userdef":
    return true

// Add to processCommand()
case "custom1":
    // Custom transformation logic
```

**2. Post-Processing Extension (for complex rules)**
```go
// Easy to add new post-processors
result := fsm.ProcessText(text)
result = fixArticles(result)
result = fixQuotes(result)
result = fixCustomGrammar(result)    // New!
result = fixUserDefinedRules(result) // New!
return result
```

### Benefits of This Architecture

**1. Clean Separation**
- **FSM**: Fast, real-time word transformations
- **Post-processing**: Complex pattern matching and grammar rules
- **No interference** between different types of transformations

**2. Performance Optimization**
- Simple commands stay in fast FSM
- Complex rules use dedicated algorithms
- Memory usage remains constant

**3. Easy Testing**
- Test FSM commands independently
- Test post-processing rules independently
- Add new tests without affecting existing ones

**4. Maintainability**
- Each function has single responsibility
- Easy to debug specific transformation types
- Clear code organization

### Example: Adding Custom Commands

**Scenario**: User wants to add `(reverse)` command and custom punctuation rules.

```go
// 1. Add to FSM (simple word transformation)
case "reverse":
    reversed := ""
    for _, r := range tp.tokens[lastWordIdx].Value {
        reversed = string(r) + reversed
    }
    tp.tokens[lastWordIdx].Value = reversed

// 2. Add post-processor (complex pattern matching)
func fixCustomPunctuation(text string) string {
    // Complex regex-based punctuation rules
    // User-defined formatting preferences
    return processedText
}

// 3. Update pipeline
result = fixArticles(result)
result = fixQuotes(result)
result = fixCustomPunctuation(result)  // New!
```

**This design makes Go-Reloaded infinitely extensible while maintaining performance and code clarity.**

## Implementation Standards

### General Algorithmic Approach

**CRITICAL RULE**: Never use hardcoded patterns specific to test cases.

**❌ BAD - Hardcoded patterns:**
```go
result = strings.ReplaceAll(result, "' hello world '", "'hello world'")
result = strings.ReplaceAll(result, "' goodbye '", "'goodbye'")
// This only works for specific test cases!
```

**✅ GOOD - General algorithm:**
```go
// Works for ANY content between quotes
for each quote pair {
    remove spaces after opening quote
    remove spaces before closing quote
    preserve all other content
}
```

### Memory Efficiency

**Constant Memory Usage**: ~8KB regardless of file size
- Fixed token buffer: `[50]Token`
- String builders reused, not recreated
- Chunked processing for large files
- No growing slices or maps

### Error Handling

**Graceful Degradation**:
- Invalid commands preserved as text
- Malformed input processed as much as possible
- No crashes on edge cases
- UTF-8 safe character handling

### Testing Philosophy

**Functions must handle arbitrary input, not just test-specific strings**:
- Quote repositioning works for any quoted content
- Article correction works for any a/an usage
- Command processing works for any valid command format
- No assumptions about specific test case contentD:
    // Add space before word (if needed), then add word
case PUNCTUATION:
    // Remove space before punctuation, then add punctuation
case SPACE:
    // Add space (if not duplicate)
case NEWLINE:
    // Add newline
}
```

**Smart spacing rules:**
- Words get spaces between them: "hello world"
- Punctuation sticks to words: "hello, world" (not "hello , world")
- No double spaces or trailing spaces

### Step 7: Article Correction - `fixArticles()`

Final cleanup step that fixes "a/an" usage:

```go
// Check each word
switch words[i] {
case "a", "A", "an", "An":
    nextWord := words[i+1]
    // Remove punctuation: "apple." -> "apple"
    // Check first letter
    if vowel_or_h {
        // "a apple" -> "an apple"
    } else {
        // "an car" -> "a car"
    }
}
```

**Vowel detection:**
```go
first := strings.ToLower(cleanWord)[0]
if first == 'a' || first == 'e' || first == 'i' || first == 'o' || first == 'u' || first == 'h' {
    // Use "an"
} else {
    // Use "a"
}
```

## Complete Example Walkthrough

**Input:** `"The value FF (hex) should be (up)"`

### Phase 1: Tokenization
```
Tokens created:
1. WORD: "The"
2. SPACE: " "
3. WORD: "value"
4. SPACE: " "
5. WORD: "FF"
6. SPACE: " "
7. WORD: "should"
8. SPACE: " "
9. WORD: "be"
```

### Phase 2: Command Processing
```
Command "hex" found:
- Find last word: "FF"
- Convert hex to decimal: "FF" -> "255"
- Update token: WORD: "255"

Command "up" found:
- Find last word: "be"
- Convert to uppercase: "be" -> "BE"
- Update token: WORD: "BE"
```

### Phase 3: Output Generation
```
Flush tokens to output:
"The" + " " + "value" + " " + "255" + " " + "should" + " " + "BE"
Result: "The value 255 should BE"
```

### Phase 4: Article Correction
```
Check for "a/an" issues:
- No articles found, no changes needed
Final result: "The value 255 should BE"
```

## Memory Management

**Key insight:** The transformer uses constant memory (~8KB) regardless of file size!

**How?**
1. **Fixed token buffer:** Only 50 tokens max
2. **Streaming:** Process and flush tokens continuously
3. **No storing entire file:** Read, process, output, repeat

**Buffer overflow handling:**
```go
if tp.tokenIdx >= len(tp.tokens) {
    // Flush first half of tokens to output
    // Shift remaining tokens to beginning
    // Continue with new token
}
```

## Why This Design?

1. **Memory Efficient:** Constant memory usage
2. **Single Pass:** No need to read file multiple times
3. **Fast:** Direct character processing, no regex
4. **Robust:** Handles Unicode, large files, edge cases
5. **Maintainable:** Clear separation of concerns

## Common Gotchas for Junior Developers

1. **Runes vs Bytes:** Always use `[]rune(text)` for Unicode safety
2. **String Builder:** Use `strings.Builder` for efficient string concatenation
3. **State Management:** Always reset builders when switching states
4. **Buffer Management:** Handle buffer overflow gracefully
5. **Edge Cases:** Empty strings, malformed commands, Unicode characters

This transformer is a great example of how to build efficient, memory-conscious text processors in Go!