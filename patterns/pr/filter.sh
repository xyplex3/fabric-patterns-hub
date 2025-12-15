#!/bin/bash
# Post-processing filter for PR pattern output
# Removes empty sections, placeholder text, and code block markers

# Read input from stdin
input=$(cat)

# Remove code block markers (backticks)
output="${input//\`\`\`/}"

# Remove empty sections with placeholder text
output=$(echo "$output" | sed -E '/\*\*Added:\*\*/,/^$/{ /No (content|features|code|changes)/d; /Nothing (added|was added)/d; }')
output=$(echo "$output" | sed -E '/\*\*Changed:\*\*/,/^$/{ /No (content|features|code|changes)/d; /Nothing (changed|was changed)/d; }')
output=$(echo "$output" | sed -E '/\*\*Removed:\*\*/,/^$/{ /No (content|features|code|changes)/d; /Nothing (removed|was removed)/d; }')

# Remove section headers that have no content following them (up to next section or end)
output=$(echo "$output" | awk '
BEGIN { in_section=0; section_name=""; buffer="" }
/^\*\*Added:\*\*/ || /^\*\*Changed:\*\*/ || /^\*\*Removed:\*\*/ {
    if (in_section && buffer ~ /^[[:space:]]*$/) {
        # Previous section was empty, skip it
        buffer=""
    } else if (in_section) {
        print section_header
        printf "%s", buffer
    }
    section_header=$0
    buffer=""
    in_section=1
    next
}
/^\*\*Key Changes:\*\*/ {
    if (in_section && buffer ~ /^[[:space:]]*$/) {
        buffer=""
    } else if (in_section) {
        print section_header
        printf "%s", buffer
    }
    print $0
    in_section=0
    buffer=""
    next
}
{
    if (in_section) {
        buffer=buffer $0 "\n"
    } else {
        print $0
    }
}
END {
    if (in_section && buffer !~ /^[[:space:]]*$/) {
        print section_header
        printf "%s", buffer
    }
}
')

echo "$output"
