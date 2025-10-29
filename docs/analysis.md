# Analysis Document — Go Reloaded

---

## 1. Problem Description

**Go Reloaded** is a text-processing tool written in Go that performs **context-aware text correction and formatting**.  
Its purpose is to read an input text file, apply a series of predefined transformation rules, and write the corrected version to an output file.

The program behaves as a lightweight text editor and formatter, capable of detecting special markers such as `(hex)`, `(bin)`, `(up)`, `(low)`, `(cap)`, and punctuation errors, then applying the corresponding corrections automatically.

---

### Transformation Rules

Each rule targets a specific kind of transformation:

- **(hex)** → converts the previous *hexadecimal number* into decimal.  
- **(bin)** → converts the previous *binary number* into decimal.  
- **(up)**, **(low)**, **(cap)** → modify the *case* of the preceding word or a group of words when followed by a number, e.g. `(up, 2)`.  
- **Punctuation rules** ensure proper spacing: commas, periods, and exclamation marks must be adjacent to the preceding word but separated by one space from the next.  
- **Quotation marks `' '`** are repositioned to correctly wrap the intended words.  
- The article **“a”** becomes **“an”** before vowels or “h”.  

---

### Example Transformations

```text
"This is so fun (up, 2) !" → "This is SO FUN!"
"Simply add 42 (hex) and 10 (bin)" → "Simply add 66 and 2"
"There it was. A amazing rock!" → "There it was. An amazing rock!"
```

### Efficiency and Design

To make the process lightweight and memory-efficient, the program reads the input file in overlapping chunks rather than loading it entirely into memory. Each chunk is **CHUNK_BYTES** in size, with **OVERLAP_WORDS** words preserved from the previous chunk to maintain context for transformations that reference preceding words.

**Critical constraint**: each chunk boundary must align with complete UTF-8 runes to avoid corrupting multi-byte characters. The parser ensures chunks end at valid rune boundaries, never splitting Unicode characters.

#### Chunk Processing Workflow

1. **Parser** reads CHUNK_BYTES from input file, ensuring the chunk ends at a complete UTF-8 rune boundary
2. **FSM Transformer** processes the chunk, applying all transformation rules
3. **Word Separation**: The last OVERLAP_WORDS from the processed chunk are stored in memory
4. **Exporter** writes the remaining words (excluding the stored overlap) to the output file
5. **Next Iteration**: Parser reads the next CHUNK_BYTES starting from where the previous chunk ended
6. **Context Merging**: The stored OVERLAP_WORDS are prepended to the new chunk before FSM processing
7. Process repeats until end of file

This overlapping design ensures commands like `(up, n)` can access words from previous chunks while maintaining minimal memory usage.

---

### Objectives

Build a robust, efficient text manipulation system using Go’s filesystem (fs) and string handling libraries.

Ensure the program remains light on memory usage and fast in execution, even with large files.

Design a clear architecture that separates input reading, transformation logic, and output writing.

Demonstrate understanding of stateful processing, modular code design, and testing for text manipulation.

---

## 2. Rules and Examples

| **Rule**                            | **Description**                                 | **Example**                             |
| ----------------------------------- | ----------------------------------------------- | --------------------------------------- |
| `(hex)`                             | Converts previous hexadecimal number to decimal | `"1E (hex)" → "30"`                     |
| `(bin)`                             | Converts previous binary number to decimal      | `"10 (bin)" → "2"`                      |
| `(up)`                              | Uppercases previous word                        | `"go (up)" → "GO"`                      |
| `(low)`                             | Lowercases previous word                        | `"HELLO (low)" → "hello"`               |
| `(cap)`                             | Capitalizes previous word                       | `"bridge (cap)" → "Bridge"`             |
| `(up, n)` / `(low, n)` / `(cap, n)` | Affects previous *n* words                      | `"so exciting (up, 2)" → "SO EXCITING"` |
| Punctuation                         | Fixes spacing between punctuation and words     | `"Hello ,world !!" → "Hello, world!!"`  |
| Quotes `' '`                        | Moves quotation marks next to the correct words | `" ' awesome ' " → "'awesome'"`         |
| `a → an`                            | Changes “a” to “an” before vowels or “h”        | `"a apple" → "an apple"`                |

---

## 3. Architecture Comparison — Pipeline vs FSM
| **Aspect**            | **Pipeline**                  | **Finite State Machine (FSM)**                |
| --------------------- | ----------------------------- | --------------------------------------------- |
| **Flow**              | Linear, step-by-step          | Event-driven, based on states                 |
| **Control**           | Fixed sequence                | Dynamic transitions                           |
| **Context awareness** | Limited                       | Full (can look back and react to context)     |
| **Best for**          | Simple static transformations | Complex text processing and conditional logic |

----

Chosen Architecture: FSM  
The FSM approach allows the program to handle contextual transformations such as backtracking across words or maintaining punctuation state between chunks. This design supports the “chunked” reading strategy while keeping memory usage minimal.

---
### 4. FSM Logic Overview

The FSM operates across the main stages of the program:

 1. Parser → Reads data from file in chunks of CHUNK_BYTES bytes, ensuring chunk boundaries align with complete UTF-8 runes.

2. FSM Transformer → Applies transformation rules while maintaining context between chunks.

3. Exporter → Writes processed text to the output file progressively.

4. Controller → Orchestrates the workflow (Parser → FSM → Exporter).

---

**Main FSM States (conceptual):**

* READ_WORD

* APPLY_RULE

* HANDLE_PUNCTUATION

* WRITE_OUTPUT

* ERROR_STATE

Each transition is triggered by detecting markers like (hex), (cap), punctuation symbols, or chunk boundaries.

---

## 5. Future Improvements

* Add error logging and detailed reporting for invalid markers.
* Implement parallel chunk processing for very large files.
* Introduce custom configuration files for user-defined rules.
* Expand test coverage with randomized and stress test files.