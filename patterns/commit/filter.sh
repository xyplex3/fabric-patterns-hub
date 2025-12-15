#!/bin/bash
# Post-processing filter for commit pattern output
# Removes empty sections, placeholder text, and code block markers

# Read input from stdin
input=$(cat)

# Remove code block markers (backticks)
output="${input//\`\`\`/}"

# Remove empty sections with placeholder text (case-insensitive)
output=$(echo "$output" | sed -E '/\*\*Added:\*\*/,/^$/{ /[Nn]o (content|features|code|changes)/d; /[Nn]othing (added|was added)/d; }')
output=$(echo "$output" | sed -E '/\*\*Changed:\*\*/,/^$/{ /[Nn]o (content|features|code|changes)/d; /[Nn]othing (changed|was changed)/d; }')
output=$(echo "$output" | sed -E '/\*\*Removed:\*\*/,/^$/{ /[Nn]o (content|features|code|changes)/d; /[Nn]othing (removed|was removed)/d; }')

# Remove section headers that have no content following them
output=$(echo "$output" | awk '
BEGIN {
    in_section=0
    section_header=""
    buffer=""
    first_line=1
}

# Keep the summary line (first non-empty line)
NR==1 || (first_line && NF>0) {
    if (first_line) {
        print $0
        first_line=0
    }
    next
}

# Detect section headers
/^\*\*Added:\*\*/ || /^\*\*Changed:\*\*/ || /^\*\*Removed:\*\*/ {
    # Print previous section if it had content
    if (in_section && buffer !~ /^[[:space:]]*$/) {
        print section_header
        printf "%s", buffer
    }
    section_header=$0
    buffer=""
    in_section=1
    next
}

# Collect content for current section
{
    if (in_section) {
        buffer=buffer $0 "\n"
    } else {
        print $0
    }
}

END {
    # Print last section if it had content
    if (in_section && buffer !~ /^[[:space:]]*$/) {
        print section_header
        printf "%s", buffer
    }
}
')

# Clean up excessive blank lines (max 2 consecutive)
output=$(echo "$output" | awk '
BEGIN { blank_count=0 }
/^[[:space:]]*$/ {
    blank_count++
    if (blank_count <= 2) print
    next
}
{
    blank_count=0
    print
}
')

# Trim trailing whitespace
output=$(echo "$output" | sed 's/[[:space:]]*$//')

echo "$output"
