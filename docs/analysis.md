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

To make the process lightweight and memory-efficient, the program reads the input file in chunks rather than loading it entirely into memory. Each chunk is processed independently by the FSM, which maintains minimal state information — just enough to preserve context between words and sentences. This streaming approach not only saves memory but also makes the program faster, as it avoids unnecessary allocations and allows the exporter to write results progressively.

---

### Objectives

Build a robust, efficient text manipulation system using Go’s filesystem (fs) and string handling libraries.

Ensure the program remains light on memory usage and fast in execution, even with large files.

Design a clear architecture that separates input reading, transformation logic, and output writing.

Demonstrate understanding of stateful processing, modular code design, and testing for text manipulation.

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


## 3. Architecture Comparison — Pipeline vs FSM
| **Aspect**            | **Pipeline**                  | **Finite State Machine (FSM)**                |
| --------------------- | ----------------------------- | --------------------------------------------- |
| **Flow**              | Linear, step-by-step          | Event-driven, based on states                 |
| **Control**           | Fixed sequence                | Dynamic transitions                           |
| **Context awareness** | Limited                       | Full (can look back and react to context)     |
| **Best for**          | Simple static transformations | Complex text processing and conditional logic |


Chosen Architecture: FSM  
The FSM approach allows the program to handle contextual transformations such as backtracking across words or maintaining punctuation state between chunks. This design supports the “chunked” reading strategy while keeping memory usage minimal.

### 4. FSM Logic Overview

The FSM operates across the main stages of the program:

 1. Parser → Reads data from file in chunks of fixed size.

2. FSM Transformer → Applies transformation rules while maintaining context between chunks.

3. Exporter → Writes processed text to the output file progressively.

4. Controller → Orchestrates the workflow (Parser → FSM → Exporter).

**Main FSM States (conceptual):**

* READ_WORD

* APPLY_RULE

* HANDLE_PUNCTUATION

* WRITE_OUTPUT

* ERROR_STATE

Each transition is triggered by detecting markers like (hex), (cap), punctuation symbols, or chunk boundaries.