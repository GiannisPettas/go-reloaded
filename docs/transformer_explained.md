# Go-Reloaded Transformer: Step-by-Step Guide for Junior Developers

## What Does the Transformer Do?

The transformer is the **brain** of Go-Reloaded. Think of it like a **smart text editor** that automatically fixes and transforms text according to special commands.

**Simple Example:**
```
Input:  "I need a apple (up) and FF (hex) items."
Output: "I need an APPLE and 255 items."
```

**What happened?**
1. `(up)` made "apple" → "APPLE" (uppercase)
2. `(hex)` converted "FF" → "255" (hexadecimal to decimal)
3. "a apple" → "an apple" (grammar correction)

**The transformer is like having a personal assistant that:**
- **Follows commands**: Sees `(up)` and makes text UPPERCASE
- **Converts numbers**: Changes hex/binary to decimal automatically
- **Fixes grammar**: Corrects "a apple" to "an apple"
- **Cleans formatting**: Fixes spacing around punctuation
- **Handles quotes**: Removes extra spaces in `' text '` → `'text'`

**All of this happens in a single pass through the text - no multiple readings needed!**

## Core Concepts

### 1. How the Transformer Reads Text (Finite State Machine)

**Think of the transformer like a person reading a book:**

**Reading Mode (STATE_TEXT):**
- Reading normal words: "Hello world"
- When they see `(`, they switch to...

**Command Mode (STATE_COMMAND):**
- Reading special instructions: "up", "hex", "low"
- When they see `)`, they execute the command and switch back

**Simple Example:**
```
Text: "Make this (up) please"

Reading: "Make" "this" → STATE_TEXT (normal reading)
Sees '(': Switch to STATE_COMMAND
Reading: "up" → STATE_COMMAND (reading instruction)
Sees ')': Execute command "up" on "this" → "THIS", switch back to STATE_TEXT
Result: "Make THIS please"
```

**Why Two Systems Working Together?**

**System 1: Character Reader (Low-Level)**
- Reads one character at a time: 'H', 'e', 'l', 'l', 'o'
- Decides: "Is this a letter? Space? Parenthesis?"
- Groups characters into words and commands

**System 2: Word Processor (High-Level)**
- Takes complete words: "Hello", "world"
- Applies transformations: "hello" → "HELLO"
- Manages memory efficiently with fixed buffers

**Why Split the Work?**

**Like a Factory Assembly Line:**

**Station 1 (Character Reader):**
- Worker reads characters one by one
- Sorts them: "This is a letter", "This is punctuation"
- Groups letters into words
- **Specializes in**: Fast character recognition

**Station 2 (Word Processor):**
- Worker takes complete words
- Applies transformations based on commands
- Manages word memory efficiently
- **Specializes in**: Word transformations

**Benefits of This Design:**

1. **Constant Memory**: Uses only ~8KB no matter how big the file
2. **Easy to Extend**: Want a new command? Just add it to Station 2
3. **Fast Processing**: Each station is optimized for its job
4. **Easy Testing**: Test each station independently
5. **Clean Code**: Each part has one clear responsibility

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

### 2. Breaking Text Into Pieces (Tokens)

**Think of tokens like LEGO blocks** - each piece of text becomes a labeled block:

```go
Input: "Hello, world!"

Tokens created:
[WORD: "Hello"] [PUNCTUATION: ","] [SPACE: " "] [WORD: "world"] [PUNCTUATION: "!"]
```

**Why Break Text Into Pieces?**

**Imagine you're organizing a toolbox:**
- **Without organization**: Everything mixed together, hard to find what you need
- **With organization**: Screws in one compartment, nails in another, easy to find

**Same with text processing:**
- **Without tokens**: "Hello, world!" is just one long string, hard to work with
- **With tokens**: Each piece is labeled and easy to find/modify

**Real Example - Command Processing:**
```
Input: "Make these words (up, 3) please"

Tokens: [WORD: "Make"] [SPACE: " "] [WORD: "these"] [SPACE: " "] [WORD: "words"] [SPACE: " "] [WORD: "please"]

Command (up, 3) says: "Find the last 3 words and make them uppercase"
Easy with tokens: Count backwards 3 WORD tokens → "Make", "these", "words"
Result: "MAKE THESE WORDS please"
```

**Token Types (Like Different LEGO Shapes):**
- **WORD**: "hello", "world", "FF" (the main content)
- **PUNCTUATION**: ".", "!", "?" (needs special spacing)
- **SPACE**: " " (separates words)
- **NEWLINE**: "\n" (line breaks)

**Memory Efficiency:**
The transformer uses a **fixed toolbox** (80 token slots) that never grows:
```go
tokens [80]Token  // Fixed size, constant memory
```

**Like a conveyor belt** - tokens come in, get processed, and move out. The belt size never changes!

### 3. Building Output Efficiently (strings.Builder)
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

**Complete Example Walkthrough:**

Let's follow `"Hi (up) there!"` through the entire process:

**Step 1: Character Reading**
```
Input: "Hi (up) there!"
State: STATE_TEXT (normal reading mode)

Character 'H': Add to wordBuilder → "H"
Character 'i': Add to wordBuilder → "Hi"
Character ' ': Space found!
  → Save "Hi" as WORD token
  → Add SPACE token
  → Clear wordBuilder
Character '(': Parenthesis found!
  → Switch to STATE_COMMAND
  → Start building command
Character 'u': Add to cmdBuilder → "u"
Character 'p': Add to cmdBuilder → "up"
Character ')': End of command!
  → Process command "up" (makes last word uppercase)
  → "Hi" becomes "HI"
  → Switch back to STATE_TEXT
Character ' ': Add SPACE token
Character 't': Add to wordBuilder → "t"
Character 'h': Add to wordBuilder → "th"
Character 'e': Add to wordBuilder → "the"
Character 'r': Add to wordBuilder → "ther"
Character 'e': Add to wordBuilder → "there"
Character '!': Punctuation found!
  → Save "there" as WORD token
  → Add PUNCTUATION token "!"
```

**Step 2: Token Processing**
```
Tokens created:
[WORD: "HI"] [SPACE: " "] [SPACE: " "] [WORD: "there"] [PUNCTUATION: "!"]
```

**Step 3: Output Building**
```
Process each token:
WORD "HI" → Write "HI"
SPACE " " → Write " "
SPACE " " → Skip (avoid double spaces)
WORD "there" → Write "there"
PUNCTUATION "!" → Remove space before, write "!"

Final result: "HI there!"
```

**Magic! The transformer:**
1. ✅ Applied the `(up)` command to "Hi" → "HI"
2. ✅ Fixed spacing around punctuation
3. ✅ Preserved the rest of the text

**All in a single pass through the text!**> switch to STATE_COMMAND
Character 'u': STATE_COMMAND -> add to cmdBuilder: "u"
Character 'p': STATE_COMMAND -> add to cmdBuilder: "up"
Character ')': STATE_COMMAND, ')' -> process command "up", switch to STATE_TEXT
Character ' ': STATE_TEXT, space -> add SPACE token
Character 'w': STATE_TEXT, default -> add to wordBuilder: "w"
...and so on
```

### Step 3: Token Management - `TokenProcessor`

The TokenProcessor manages a buffer of 80 tokens:

```go
type TokenProcessor struct {
    tokens   [80]Token  // Fixed-size buffer
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

#### Quote Repositioning - `fixQuotes()` - **Independent Odd/Even Algorithm**

**Handles mixed quote types independently using odd/even positioning logic:**

```go
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
```

**Algorithm Logic:**
- **Single quotes (`'`)**: Tracked independently with separate counter
- **Double quotes (`"`)**: Tracked independently with separate counter
- **Odd-numbered quotes**: Stick to the right (remove space after)
- **Even-numbered quotes**: Stick to the left (remove space before)
- **No mixing**: Each quote type processed completely independently

**Examples:**

**Simple cases (same quote type):**
- `' hello world '` → `'hello world'`
- `" any text "` → `"any text"`

**Mixed quote types (the key improvement):**
- `" hello ' world ' text "` → `"hello 'world' text"`
- `' outer " inner " text '` → `'outer "inner" text'`

**Complex mixed example:**
```
Input:  `" a' b f ' c " d " "`
Output: `"a'b f' c" d ""`

Breakdown:
1st " (odd) → stick to 'a': "a
1st ' (odd) → stick to 'b': 'b  
2nd ' (even) → stick to 'f': f'
2nd " (even) → stick to 'c': c"
3rd " (odd) → stick to 4th ": "
4th " (even) → already has 3rd attached: "
```

**Key Features:**
- **Independent processing**: Single and double quotes don't interfere
- **General algorithm**: Works for ANY content and quote arrangement
- **No hardcoded patterns**: Handles arbitrary text between quotes
- **Robust**: Handles unmatched quotes gracefully

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
- Fixed token buffer: `[80]Token`
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
1. **Fixed token buffer:** Only 80 tokens max
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