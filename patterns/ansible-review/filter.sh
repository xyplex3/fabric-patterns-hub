#!/usr/bin/env bash

# Filter script for ansible-review pattern
# Cleans up LLM output to ensure consistent markdown formatting

set -euo pipefail

# Read all input
content=$(cat)

# Remove any markdown code fences that might wrap the entire output
content=$(echo "$content" | sed -e '1{/^```markdown$/d;}' -e '${/^```$/d;}')

# Remove any leading/trailing whitespace
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

# Remove placeholder text for empty sections
content=$(echo "$content" | sed -E '/^\[No (CRITICAL|HIGH|MEDIUM|LOW|INFO) (issues|findings)\]$/d')
content=$(echo "$content" | sed -E '/^\[None found\]$/d')
content=$(echo "$content" | sed -E '/^\[No issues found\]$/d')
content=$(echo "$content" | sed -E '/^- None$/d')
content=$(echo "$content" | sed -E '/^- No issues found$/d')

# Remove empty severity sections (section header followed by next section or end)
content=$(echo "$content" | awk '
BEGIN { buffer = ""; section_type = "" }
/^## (Critical Issues|Improvements|Positive Observations|Recommendations)/ {
    section_type = $0
    buffer = $0 "\n"
    getline
    # Read until we find content or the next section
    while (length($0) == 0 || $0 ~ /^---$/) {
        if ($0 ~ /^---$/) {
            # Separator found, check next line
            getline
            if ($0 ~ /^## (Critical Issues|Improvements|Positive Observations|Recommendations|Summary)/) {
                # Next section found, this section was empty
                section_type = ""
                buffer = ""
                print $0
                next
            } else {
                # Real content after separator
                print buffer
                print "---"
                print $0
                buffer = ""
                section_type = ""
                next
            }
        }
        buffer = buffer $0 "\n"
        getline
    }
    # Found content
    print buffer
    print $0
    buffer = ""
    section_type = ""
    next
}
{ print }
')

# Preserve code blocks properly - don't strip YAML code blocks
content=$(echo "$content" | awk '
/^```yaml$/ { in_code = 1; print; next }
/^```bash$/ { in_code = 1; print; next }
/^```jinja2$/ { in_code = 1; print; next }
/^```markdown$/ { in_code = 1; print; next }
/^```$/ && in_code { in_code = 0; print; next }
{ print }
')

# Ensure the output ends with exactly one newline
content=$(echo "$content" | sed -e :a -e '/^\n*$/{$d;N;ba' -e '}')

# Output the cleaned content
echo "$content"
echo
