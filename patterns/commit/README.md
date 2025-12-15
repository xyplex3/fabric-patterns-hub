# Commit Pattern

Generate clear, informative commit messages following the Conventional Commits specification.

## Usage

### Basic Usage

```bash
git diff --staged | fabric -p commit
```

### With Post-Processing Filter

To ensure clean output without empty sections or placeholder text:

```bash
git diff --staged | fabric -p commit | ./patterns/commit/filter.sh
```

### Direct Commit

Use the pattern output directly for committing:

```bash
# Generate commit message and review it
commit_msg=$(git diff --staged | fabric -p commit | ./patterns/commit/filter.sh)
echo "$commit_msg"

# Commit with the generated message
git commit -m "$commit_msg"
```

### With Git Hooks

You can integrate this into your git workflow using a prepare-commit-msg hook:

```bash
#!/bin/bash
# .git/hooks/prepare-commit-msg

COMMIT_MSG_FILE=$1
COMMIT_SOURCE=$2

# Only run for regular commits (not merge, squash, etc.)
if [ -z "$COMMIT_SOURCE" ]; then
    # Generate commit message from staged changes
    commit_msg=$(git diff --staged | fabric -p commit | ./patterns/commit/filter.sh)

    # Write to commit message file
    echo "$commit_msg" > "$COMMIT_MSG_FILE"
fi
```

## Features

- Follows Conventional Commits specification
- Organizes changes into Added, Changed, and Removed sections when needed
- Omits empty sections automatically
- Removes placeholder text
- No code block markers in output
- Optional body - only includes detailed sections when necessary

## Conventional Commit Types

- `feat`: A new feature
- `fix`: A bug fix
- `docs`: Documentation only changes
- `style`: Code formatting changes (no logic changes)
- `refactor`: Code restructuring without changing behavior
- `perf`: Performance improvements
- `test`: Adding or updating tests
- `build`: Build system or dependency changes
- `ci`: CI/CD configuration changes
- `chore`: Other changes (tooling, etc.)
- `revert`: Reverts a previous commit

## Output Format

### Simple Commit (Summary Only)

For straightforward changes:

```
feat: add user authentication module
```

### Detailed Commit (With Body)

For complex changes with multiple components:

```
feat: add user authentication module

**Added:**

- User authentication middleware with JWT support
- Login and registration endpoints in auth controller
- Password hashing utilities using bcrypt

**Changed:**

- Updated user model to include password and token fields
- Modified API routes to include authentication middleware

**Removed:**

- Deprecated session-based authentication code
```

Note: Only sections with actual content are included. Empty sections are completely omitted.

## Tips

### When to Use Body Sections

- **Use body sections** when:
  - Multiple files or components are affected
  - Changes need context or explanation
  - There are distinct additions, changes, and removals

- **Use summary only** when:
  - Single file change with obvious purpose
  - Simple bug fix or typo correction
  - Self-explanatory change

### Writing Good Commit Messages

- Keep the summary line under 80 characters
- Use present tense ("add" not "added")
- Start with lowercase after the type
- Be specific about what changed, not just which files
- Explain "why" when it's not obvious

## Examples

### Example 1: Simple Change

```bash
# Input: Minor documentation fix
git diff --staged | fabric -p commit | ./patterns/commit/filter.sh

# Output:
docs: fix typo in installation instructions
```

### Example 2: Feature Addition

```bash
# Input: New feature with multiple files
git diff --staged | fabric -p commit | ./patterns/commit/filter.sh

# Output:
feat: add email notification system

**Added:**

- Email service with SMTP configuration and templating
- Notification triggers for user registration and password reset
- Email templates for common user actions

**Changed:**

- Updated user controller to send welcome emails on registration
- Modified environment configuration to include SMTP settings
```

### Example 3: Refactoring

```bash
# Input: Code restructuring without new features
git diff --staged | fabric -p commit | ./patterns/commit/filter.sh

# Output:
refactor: reorganize database models into separate modules

**Changed:**

- Split monolithic models file into individual model files
- Updated import paths across application to reference new structure
- Improved model organization by domain (user, product, order)
```
