#!/usr/bin/env bash

# Filter script for python-review pattern
# Cleans up LLM output to ensure consistent formatting

set -euo pipefail

# Read all input
content=$(cat)

# Remove any markdown code fences that might wrap the entire output
content=$(echo "$content" | sed -e '1{/^```markdown$/d;}' -e '${/^```$/d;}')

# Remove any leading/trailing whitespace lines
content=$(echo "$content" | sed -e :a -e '/./,$!d;/^\n*$/{$d;N;};/\n$/ba')

# Ensure exactly one blank line before each major section header
content=$(echo "$content" | awk '
BEGIN { blank_count = 0; prev_was_separator = 0 }
/^---$/ {
    if (!prev_was_separator) {
        print ""
        print $0
        prev_was_separator = 1
        blank_count = 0
    }
    next
}
/^#/ {
    if (NR > 1 && blank_count == 0) {
        print ""
    }
    print $0
    blank_count = 0
    prev_was_separator = 0
    next
}
/^$/ {
    blank_count++
    if (blank_count == 1) {
        print $0
    }
    prev_was_separator = 0
    next
}
{
    print $0
    blank_count = 0
    prev_was_separator = 0
}
')

# Remove any "Here is..." or "Here's..." introductory phrases
content=$(echo "$content" | sed -E '/^(Here is|Here'\''s|Below is) (the|an?) .*:?$/d')

# Ensure the output ends with exactly one newline
content=$(echo "$content" | sed -e :a -e '/^\n*$/{$d;N;ba' -e '}')

# Output the cleaned content
echo "$content"
echo
