# Golden Test Set — Go Reloaded

This file contains the functional test cases used to verify the correctness of the **Go Reloaded** transformations.  
Each case includes the input, the expected output, and a short description of the rule being tested.

---

## T1 — Overlaping commands

**Description:**  
commands that overlap eachother have to run in order from left to right

**Input:**  
Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea (up, 5) consequat. Duis (low, 6)

**Expected Output:**  
Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris NISI ut aliquip ex ea consequat. duis 


## T2 — Mixed Punctuation and Capitalization

**Description:**  
Tests the FSM's ability to handle punctuation spacing and capitalization rules simultaneously.  
Ensures that punctuation is attached correctly after transformations and that case changes don't break sentence flow.

**Input:**  
He said : ' this is incredible (cap, 2) ! ' can you believe it , though ?

**Expected Output:**  
He said: 'This Is incredible!' can you believe it, though?

## T3 — Chained Numeric Commands

**Description:**  
Verifies that multiple numeric commands can apply sequentially to the same word.  
The FSM should execute them **from left to right**, passing the result of the first conversion into the next.  
This ensures proper chaining logic and value propagation between transformations.

**Input:**  
Simply add 1010 (bin) (hex) , and check the total !

**Expected Output:**  
Simply add 16, and check the total!

## T4 — Broken or Incomplete Command

**Description:**  
Validates the FSM's ability to handle malformed or incomplete commands.  
If a command is missing its closing parenthesis or contains a punctuation error (like a period instead of a comma),  
the FSM should **ignore the faulty marker** and continue processing normally without breaking the text structure.

**Input:**  
It was a (up, 2 b(up2)) beautiful day and everything felt calm (low

**Expected Output:**  
It was a (up, 2 b(up2)) beautiful day and everything felt calm (low

## T5 — Punctuation Spacing and Grouping

**Description:**  
Checks correct spacing for mixed punctuation marks and ensures that multiple punctuation symbols remain grouped together without breaking word flow.

**Input:**  
Wait ,what ? ! This can't be real ;or can it ? Look over there ...no ,behind you !!

**Expected Output:**  
Wait, what?! This can't be real; or can it? Look over there... no, behind you!!

## T6 — Hexadecimal Conversion

**Description:**  
Tests conversion of hexadecimal numbers to decimal format using the (hex) marker.

**Input:**  
The value is 1E (hex) and another is FF (hex) .

**Expected Output:**  
The value is 30 and another is 255.

## T7 — Binary Conversion

**Description:**  
Tests conversion of binary numbers to decimal format using the (bin) marker.

**Input:**  
Binary 1010 (bin) equals decimal and 11111111 (bin) is maximum byte .

**Expected Output:**  
Binary 10 equals decimal and 255 is maximum byte.

## T8 — Article Correction (a to an)

**Description:**  
Tests automatic correction of "a" to "an" before vowels and "h".

**Input:**  
It was a amazing day with a elephant and a honest person .

**Expected Output:**  
It was an amazing day with an elephant and an honest person.

## T9 — Quote Repositioning

**Description:**  
Tests proper positioning of single quotes around words.

**Input:**  
He said ' hello world ' and then ' goodbye ' .

**Expected Output:**  
He said 'hello world' and then 'goodbye'.

## T10 — Complex Case Transformations

**Description:**  
Tests multiple case transformation commands with different word counts.

**Input:**  
this is amazing (cap, 3) and that was great (up, 1) but now VERY LOUD (low, 2) .

**Expected Output:**  
This Is Amazing and that was GREAT but now very loud.

## T11 — Mixed Numeric Conversions

**Description:**  
Tests multiple numeric conversions in sequence.

**Input:**  
First A (hex) then 101 (bin) and finally 1F (hex) .

**Expected Output:**  
First 10 then 5 and finally 31.

## T12 — Edge Case: Zero Values

**Description:**  
Tests conversion of zero values in different bases.

**Input:**  
Zero in hex is 0 (hex) and binary 0 (bin) .

**Expected Output:**  
Zero in hex is 0 and binary 0.

## T13 — Punctuation with Transformations

**Description:**  
Tests punctuation spacing combined with case transformations.

**Input:**  
wow (cap) ,this is amazing (up, 2) !great .

**Expected Output:**  
Wow, this IS AMAZING! great.

## T14 — Invalid Commands Ignored

**Description:**  
Tests that malformed commands are ignored and treated as regular text.

**Input:**  
This (invalid) and ( up, text) should remain unchanged .

**Expected Output:**  
This (invalid) and ( up, text) should remain unchanged.

## T15 — Large Number Conversions

**Description:**  
Tests conversion of larger hexadecimal and binary numbers.

**Input:**  
Large hex FFFF (hex) and binary 11111111 (bin) .

**Expected Output:**  
Large hex 65535 and binary 255.

## T16 — Multiple Quotes in Sentence

**Description:**  
Tests handling of multiple quote pairs in the same sentence.

**Input:**  
She said ' first quote ' then ' second quote ' and ' third quote ' .

**Expected Output:**  
She said 'first quote' then 'second quote' and 'third quote'.

## T17 — Case Sensitivity in Hex

**Description:**  
Tests that hexadecimal conversion works with both uppercase and lowercase letters.

**Input:**  
Lowercase abc (hex) and uppercase ABC (hex) .

**Expected Output:**  
Lowercase 2748 and uppercase 2748.

## T18 — Boundary Case: Single Character

**Description:**  
Tests transformations on single character words.

**Input:**  
a (up) b (cap) C (low) .

**Expected Output:**  
A B c.

## T19 — Complex Overlapping Commands

**Description:**  
Tests complex scenarios with overlapping transformation ranges.

**Input:**  
start here now (up, 3) MIDDLE SECTION (low, 2) end .

**Expected Output:**  
START HERE NOW middle section end.

## T20 — All Punctuation Types

**Description:**  
Tests spacing correction for all punctuation types mentioned in analysis.

**Input:**  
Hello ,world !How are you ?Fine ;thanks .Great ...

**Expected Output:**  
Hello, world! How are you? Fine; thanks. Great...

## T21 — Negative Hex and Binary Numbers

**Description:**  
Tests handling of negative hexadecimal and binary numbers. The command should ignore the '-' sign, convert the number, then append the '-' to the result.

**Input:**  
Negative -1A (hex) and -101 (bin) should convert .

**Expected Output:**  
Negative -26 and -5 should convert.

## T22 — Cross-Chunk Command Processing

**Description:**  
Tests FSM's ability to handle commands that reference more words than available in current chunk. The FSM should maintain state across chunk boundaries to apply transformations correctly.

**Input:**  
word1 word2 word3 word4 word5 (up, 10) remaining text .

**Expected Output:**  
WORD1 WORD2 WORD3 WORD4 WORD5 remaining text.

## T23 — Preserve Line Endings

**Description:**  
Tests that the system preserves original line endings and paragraph structure from input to output.

**Input:**  
First line with transformation (up).
Second line here.

Third line after blank line.
Final line with number A (hex).

**Expected Output:**  
First line with TRANSFORMATION.
Second line here.

Third line after blank line.
Final line with number 10.

## T24 — Contractions and Apostrophes

**Description:**  
Tests that contractions with apostrophes are handled correctly without adding extra spaces.

**Input:**  
Let's test contractions like don't, won't, and can't properly.

**Expected Output:**  
Let's test contractions like don't, won't, and can't properly.

## T25 — Invalid Commands at Start

**Description:**  
Tests that invalid commands at the beginning of text are ignored and removed since they have no preceding words to transform.

**Input:**  
(low, 2) Welcome to the comprehensive test of go-reloaded text processing system.

**Expected Output:**  
Welcome to the comprehensive test of go-reloaded text processing system.

## T26 — Invalid Commands After Numeric Conversions

**Description:**  
Tests that case transformation commands are removed when they follow numeric conversions that consume all preceding words, leaving no words to transform.

**Input:**  
The value 20 (hex) (bin) (low, 2)

**Expected Output:**  
The value 32

## T27 — Large File Chunked Processing

**Description:**  
Tests chunked processing with file size above 4KB to verify FSM works correctly across chunk boundaries. Commands should execute properly even when split across chunks, and transformations should maintain context.

**Input:**  
Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum. Sed ut perspiciatis unde omnis iste natus error sit voluptatem accusantium doloremque laudantium, totam rem aperiam, eaque ipsa quae ab illo inventore veritatis et quasi architecto beatae vitae dicta sunt explicabo. Nemo enim ipsam voluptatem quia voluptas sit aspernatur aut odit aut fugit, sed quia consequuntur magni dolores eos qui ratione voluptatem sequi nesciunt. Neque porro quisquam est, qui dolorem ipsum quia dolor sit amet, consectetur, adipisci velit, sed quia non numquam eius modi tempora incidunt ut labore et dolore magnam aliquam quaerat voluptatem. Ut enim ad minima veniam, quis nostrum exercitationem ullam corporis suscipit laboriosam, nisi ut aliquid ex ea commodi consequatur. Quis autem vel eum iure reprehenderit qui in ea voluptate velit esse quam nihil molestiae consequatur, vel illum qui dolorem eum fugiat quo voluptas nulla pariatur. At vero eos et accusamus et iusto odio dignissimos ducimus qui blanditiis praesentium voluptatum deleniti atque corrupti quos dolores et quas molestias excepturi sint occaecati cupiditate non provident, similique sunt in culpa qui officia deserunt mollitia animi, id est laborum et dolorum fuga. Et harum quidem rerum facilis est et expedita distinctio. Nam libero tempore, cum soluta nobis est eligendi optio cumque nihil impedit quo minus id quod maxime placeat facere possimus, omnis voluptas assumenda est, omnis dolor repellendus. Temporibus autem quibusdam et aut officiis debitis aut rerum necessitatibus saepe eveniet ut et voluptates repudiandae sint et molestiae non recusandae. Itaque earum rerum hic tenetur a sapiente delectus, ut aut reiciendis voluptatibus maiores alias consequatur aut perferendis doloribus asperiores repellat. The hexadecimal number FF (hex) should convert to decimal. Binary number 1010 (bin) should also convert. These three words (up, 3) should be uppercase. Article correction: I need a apple and a elephant for testing. Punctuation spacing test ,with commas !and exclamation ?marks should be fixed .Chain test: 1111 (bin) (hex) should work across chunks. More text to ensure we exceed 4KB limit and trigger chunked processing. Additional content with various transformations to test the FSM robustness across chunk boundaries. The system should handle all commands correctly regardless of chunk splits. Final test with more (cap, 2) words and A0 (hex) conversion.

**Expected Output:**  
Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum. Sed ut perspiciatis unde omnis iste natus error sit voluptatem accusantium doloremque laudantium, totam rem aperiam, eaque ipsa quae ab illo inventore veritatis et quasi architecto beatae vitae dicta sunt explicabo. Nemo enim ipsam voluptatem quia voluptas sit aspernatur aut odit aut fugit, sed quia consequuntur magni dolores eos qui ratione voluptatem sequi nesciunt. Neque porro quisquam est, qui dolorem ipsum quia dolor sit amet, consectetur, adipisci velit, sed quia non numquam eius modi tempora incidunt ut labore et dolore magnam aliquam quaerat voluptatem. Ut enim ad minima veniam, quis nostrum exercitationem ullam corporis suscipit laboriosam, nisi ut aliquid ex ea commodi consequatur. Quis autem vel eum iure reprehenderit qui in ea voluptate velit esse quam nihil molestiae consequatur, vel illum qui dolorem eum fugiat quo voluptas nulla pariatur. At vero eos et accusamus et iusto odio dignissimos ducimus qui blanditiis praesentium voluptatum deleniti atque corrupti quos dolores et quas molestias excepturi sint occaecati cupiditate non provident, similique sunt in culpa qui officia deserunt mollitia animi, id est laborum et dolorum fuga. Et harum quidem rerum facilis est et expedita distinctio. Nam libero tempore, cum soluta nobis est eligendi optio cumque nihil impedit quo minus id quod maxime placeat facere possimus, omnis voluptas assumenda est, omnis dolor repellendus. Temporibus autem quibusdam et aut officiis debitis aut rerum necessitatibus saepe eveniet ut et voluptates repudiandae sint et molestiae non recusandae. Itaque earum rerum hic tenetur a sapiente delectus, ut aut reiciendis voluptatibus maiores alias consequatur aut perferendis doloribus asperiores repellat. The hexadecimal number 255 should convert to decimal. Binary number 10 should also convert. THESE THREE WORDS should be uppercase. Article correction: I need an apple and an elephant for testing. Punctuation spacing test, with commas! and exclamation? marks should be fixed. Chain test: 21 should work across chunks. More text to ensure we exceed 4KB limit and trigger chunked processing. Additional content with various transformations to test the FSM robustness across chunk boundaries. The system should handle all commands correctly regardless of chunk splits. Final test With More words and 160 conversion.

## T28 — Complex Quote and Article Transformations

**Description:**  
Tests complex scenarios combining quote repositioning, article correction, case transformations, and punctuation spacing in a single sentence.

**Input:**  
harold wilson (cap, 2) a ' I am a (cap) optimist ,but a (up) optimist who carries ' a raincoat .

**Expected Output:**  
Harold Wilson a 'I am AN optimist, but AN optimist who carries' a raincoat.