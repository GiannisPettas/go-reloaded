# Golden Test Set — Go Reloaded

This file contains the functional test cases used to verify the correctness of the **Go Reloaded** transformations.  
Each case includes the input, the expected output, and a short description of the rule being tested.

---

## T1 — Overlaping commands
**Description**  
commands that overlap eachother have to run in order from left to right

**Input:**   
Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea (up, 5) consequat. Duis (low, 6)  

**expected**  
Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris NISI ut aliquip ex ea consequat. duis 


## T2 — Mixed Punctuation and Capitalization

**Description:**  
Tests the FSM’s ability to handle punctuation spacing and capitalization rules simultaneously.  
Ensures that punctuation is attached correctly after transformations and that case changes don’t break sentence flow.

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
Validates the FSM’s ability to handle malformed or incomplete commands.  
If a command is missing its closing parenthesis or contains a punctuation error (like a period instead of a comma),  
the FSM should **ignore the faulty marker** and continue processing normally without breaking the text structure.

**Input:**  
It was a (up. 2 b(up, 2))eautiful day and everything felt calm (low  

**Expected Output:**      
It was a (up. 2 B )eautiful day and everything felt calm (low

## T5 — Punctuation Spacing and Grouping

**Description:**  
Checks correct spacing for mixed punctuation marks and ensures that multiple punctuation symbols remain grouped together without breaking word flow.

**Input:**  
Wait ,what ?! This can't be real ;or can it ? Look over there ...no ,behind you !!


**Expected Output:**  
Wait, what?! This can't be real; or can it? Look over there... no, behind you!!
