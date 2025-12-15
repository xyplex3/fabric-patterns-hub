# PR Pattern

Generate comprehensive pull request descriptions following the Conventional Commits specification.

## Usage

### Basic Usage

```bash
git diff main | fabric -p pr
```

### With Post-Processing Filter

To ensure clean output without empty sections or placeholder text:

```bash
git diff main | fabric -p pr | ./patterns/pr/filter.sh
```

### Creating a PR with GitHub CLI

Extract title from first line and use rest as body:

```bash
# Generate PR description
pr_content=$(git diff main | fabric -p pr | ./patterns/pr/filter.sh)

# Extract title and body
title=$(echo "$pr_content" | head -n 1)
body=$(echo "$pr_content" | tail -n +3)

# Create PR
gh pr create --title "$title" --body "$body"
```

## Features

- Follows Conventional Commits specification
- Organizes changes into Key Changes, Added, Changed, and Removed sections
- Omits empty sections automatically
- Removes placeholder text
- No code block markers in output

## Output Format

```
<type>: <brief title>

**Key Changes:**

- Most important change
- Second most important change
- Third most important change

**Added:**

- Description of additions

**Changed:**

- Description of changes

**Removed:**

- Description of removals
```

Note: Only sections with actual content are included. Empty sections are completely omitted.
