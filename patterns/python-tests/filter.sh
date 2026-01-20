#!/bin/bash
# Filter script to clean up LLM output for Python test files
#
# Usage: cat output.txt | ./filter.sh > clean_output.py
#
# This script:
# - Removes markdown code fences
# - Cleans leading/trailing whitespace
# - Removes introductory phrases
# - Ensures proper formatting

# Read all input
input=$(cat)

# Remove markdown code fences (```python, ```py, ```)
output=$(echo "$input" | sed 's/^```python$//' | sed 's/^```py$//' | sed 's/^```$//')

# Remove introductory phrases that LLMs often add
output=$(echo "$output" | sed '/^Here is/d' | sed '/^Here'\''s/d' | sed '/^Below is/d')

# Remove lines that are just "---" separators
output=$(echo "$output" | sed '/^---$/d')

# Clean up excessive blank lines (more than 2 consecutive)
output=$(echo "$output" | cat -s)

# Remove leading blank lines
output=$(echo "$output" | sed '/./,$!d')

# Remove trailing blank lines and ensure single newline at end
output=$(echo "$output" | sed -e :a -e '/^\n*$/{$d;N;ba' -e '}')

# Output the cleaned content
echo "$output"
