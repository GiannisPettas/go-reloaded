# Go-Reloaded Transformer: Step-by-Step Guide for Junior Developers

## What Does the Transformer Do?

The transformer is the heart of Go-Reloaded. It takes text like:
```
"The value FF (hex) should be (up)"
```

And transforms it to:
```
"The value 255 SHOULD BE"
```

It handles commands in parentheses and fixes grammar automatically.

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

### 2. Tokens
Everything gets broken into tokens (pieces):
- **WORD**: "hello", "world", "FF"
- **PUNCTUATION**: ".", "!", "?"
- **SPACE**: " " (spaces and tabs)
- **NEWLINE**: "\n" (line breaks)

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
}
```

#### Adding Tokens - `addToken()`
- If buffer has space: add token
- If buffer is full: flush first half to output, shift remaining tokens, then add new token

**Why fixed buffer?** Memory efficiency! No matter how big the file, we only use ~8KB of memory.

### Step 4: Command Processing - `processCommand()`

When we hit a `)`, we process the command:

#### Single Word Commands
```go
switch cmdValue {
case "hex":
    // Convert hexadecimal to decimal: "FF" -> "255"
case "bin":
    // Convert binary to decimal: "1010" -> "10"
case "up", "low", "cap":
    // Transform case of the last word
}
```

#### Multi-Word Commands
```go
// Format: "(up, 3)" means uppercase the last 3 words
if strings.Contains(cmdValue, ",") {
    parts := strings.Split(cmdValue, ",")
    cmd := "up"      // command
    count := 3       // number of words
    // Apply command to last 'count' words
}
```

**Example:**
```
Tokens: ["These", "three", "words", "should", "be"]
Command: "up, 3"
Result: ["These", "THREE", "WORDS", "SHOULD", "be"]
```

### Step 5: Word Transformation - `transformWord()`

Simple case transformations:
```go
switch cmd {
case "up":   return strings.ToUpper(word)     // "hello" -> "HELLO"
case "low":  return strings.ToLower(word)     // "HELLO" -> "hello"
case "cap":  return strings.Title(word)       // "hello" -> "Hello"
}
```

### Step 6: Output Generation - `flushTokens()`

Converts tokens back to text with proper spacing:

```go
switch token.Type {
case WORD:
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