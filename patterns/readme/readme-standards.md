# README Standards

A comprehensive guide to writing effective README documentation for software projects. This document serves as the knowledge base for the readme pattern.

## Table of Contents

1. [Core Principles](#core-principles)
2. [Essential Sections](#essential-sections)
3. [Writing Guidelines](#writing-guidelines)
4. [Section Templates](#section-templates)
5. [Project Type Variations](#project-type-variations)
6. [Quality Checklist](#quality-checklist)

---

## Core Principles

### The Three Questions

A great README answers three fundamental questions in the first few paragraphs:

| Question | Answer In |
|----------|-----------|
| What is this? | Title and first paragraph |
| Why should I care? | Overview/description |
| How do I use it? | Quick Start section |

### Progressive Disclosure

Structure information from general to specific:

| Time | User Should |
|------|-------------|
| 10 seconds | Understand what it is |
| 30 seconds | Know if it's relevant |
| 5 minutes | Get started using it |
| When needed | Find detailed documentation |

### Key Principles

| Principle | Description |
|-----------|-------------|
| Clarity | Use simple, direct language |
| Completeness | Include all essential sections |
| Accuracy | Keep information current |
| Scannability | Use headings, lists, and code blocks |

---

## Essential Sections

### Required Sections (In Order)

1. **Project Title** - Clear, descriptive name
2. **One-line Description** - What it does and who it's for
3. **Overview** - 2-3 sentences expanding on the description
4. **Installation** - Prerequisites and step-by-step instructions
5. **Quick Start** - Working example users can run immediately
6. **Usage** - Detailed usage examples
7. **License** - License information

### Recommended Sections

| Section | When to Include |
|---------|-----------------|
| Badges | Public/open source projects |
| Table of Contents | READMEs longer than 200 lines |
| Features | Projects with multiple capabilities |
| Configuration | Projects with config options |
| API Reference | Libraries with programmatic APIs |
| Examples | Complex projects |
| Contributing | Open source projects |
| Testing | Projects with test suites |
| Acknowledgments | Projects with contributors/dependencies |

---

## Writing Guidelines

### Language and Style

**DO:**
- Use simple, direct language in active voice
- Write short sentences and paragraphs
- Define technical terms on first use
- Provide context for decisions

**DON'T:**
- Use vague language ("the usual setup")
- Include placeholder text ("TODO", "Coming soon")
- Create walls of text
- Assume reader knowledge without explanation

### Code Examples

**Good code examples:**
```bash
# Install the package
npm install my-package

# Run the example
npm run example
```

**Bad code examples:**
```bash
# TODO: Add installation command
install-command-here
```

### Formatting

| Element | Use For |
|---------|---------|
| `# Title` | Main title (use once) |
| `## Section` | Major sections |
| `### Subsection` | Subsections |
| `**bold**` | Emphasis and key terms |
| `` `code` `` | Variables, commands, file names |
| `> quote` | Important notes and warnings |
| Tables | Structured data |

### Visual Hierarchy

- Place key metrics and KPIs prominently
- Group related information together
- Use consistent panel sizes for related content
- Create visual separation between sections
- Use rows/sections to organize

---

## Section Templates

### Installation Section

```markdown
## Installation

### Prerequisites

- Node.js 18 or higher
- npm 9 or higher
- Git 2.30+

### Quick Install

\`\`\`bash
npm install my-package
\`\`\`

### From Source

\`\`\`bash
git clone https://github.com/user/project.git
cd project
npm install
npm run build
\`\`\`

### Verification

\`\`\`bash
my-package --version
# Expected output: my-package v1.0.0
\`\`\`
```

### Quick Start Section

```markdown
## Quick Start

Create your first project in under 60 seconds:

\`\`\`bash
# Initialize a new project
my-package init my-project
cd my-project

# Run the development server
my-package serve
\`\`\`

Open http://localhost:3000 to see your project.
For detailed usage, see the [Usage](#usage) section.
```

### Configuration Section

```markdown
## Configuration

### Environment Variables

| Variable | Description | Required | Default |
|----------|-------------|----------|---------|
| `API_KEY` | API authentication key | Yes | - |
| `DEBUG` | Enable debug logging | No | `false` |
| `PORT` | Server port | No | `3000` |

### Configuration File

Create a `config.yaml` file:

\`\`\`yaml
server:
  port: 3000
  host: localhost

logging:
  level: info
  format: json
\`\`\`
```

### Contributing Section

```markdown
## Contributing

We welcome contributions! Here's how to get started:

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/amazing-feature`
3. Make your changes and add tests
4. Run the test suite: `npm test`
5. Commit your changes: `git commit -m 'Add amazing feature'`
6. Push to the branch: `git push origin feature/amazing-feature`
7. Open a Pull Request

### Development Setup

\`\`\`bash
git clone https://github.com/user/project.git
cd project
npm install
npm run dev
\`\`\`
```

---

## Project Type Variations

### CLI Tool

**Emphasize:**
- Command-line usage examples
- Available commands and options
- Input/output formats
- Shell completion setup

**Example sections:**
```markdown
## Commands

| Command | Description |
|---------|-------------|
| `init` | Initialize a new project |
| `build` | Build the project |
| `serve` | Start development server |

## Options

| Option | Description |
|--------|-------------|
| `-v, --verbose` | Enable verbose output |
| `-c, --config` | Specify config file |
```

### Library/Package

**Emphasize:**
- API documentation
- Code examples for common use cases
- Type definitions
- Integration examples

**Example sections:**
```markdown
## API Reference

### `myFunction(options)`

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `options.name` | `string` | Yes | Name of the resource |
| `options.timeout` | `number` | No | Timeout in ms (default: 5000) |

**Returns:** `Promise<Result>`
```

### Web Application

**Emphasize:**
- Live demo link
- Screenshots/GIFs
- Deployment instructions
- Browser support

**Example sections:**
```markdown
## Demo

Try the live demo at https://example.com/demo

## Screenshots

![Dashboard](./docs/screenshots/dashboard.png)
*Dashboard view showing key metrics*
```

### Framework/SDK

**Emphasize:**
- Architecture overview
- Getting started guide
- Plugin/extension system
- Migration guides

### Monorepo

**Emphasize:**
- Package structure
- Workspace commands
- Cross-package dependencies
- Individual package READMEs

---

## Quality Checklist

### Before Publishing

- [ ] Answers "what, why, how" in first paragraphs
- [ ] Installation instructions are complete with versions
- [ ] At least one working example in Quick Start
- [ ] All code blocks specify language
- [ ] Prerequisites include version numbers
- [ ] No placeholder text (TODO, Coming soon)
- [ ] No walls of text - uses lists, code blocks, headings
- [ ] Verification steps after installation
- [ ] License section present
- [ ] Links work correctly

### For Open Source

- [ ] Contributing guidelines
- [ ] Code of conduct (or link to one)
- [ ] Issue templates mentioned
- [ ] Badges for build status, version, license
- [ ] Acknowledgments for contributors

### For Enterprise

- [ ] Security considerations
- [ ] Support contact information
- [ ] SLA information (if applicable)
- [ ] Compliance/certification info

---

## Common Mistakes

### Too Vague

**Bad:**
```markdown
Install the tool using the standard method.
Configure it as needed for your environment.
```

**Good:**
```markdown
Install with npm:
\`\`\`bash
npm install -g mytool
\`\`\`

Configure by creating `~/.mytool.yaml`:
\`\`\`yaml
api_key: your-api-key
region: us-west-2
\`\`\`
```

### Missing Prerequisites

**Bad:**
```markdown
## Installation
Run `make install` to install.
```

**Good:**
```markdown
## Installation

### Prerequisites
- Go 1.21 or higher
- Make 4.0+
- PostgreSQL 14+ (for database features)

### Install
\`\`\`bash
make install
\`\`\`
```

### No Working Examples

**Bad:**
```markdown
## Usage
Use the API to do things.
```

**Good:**
```markdown
## Usage

\`\`\`python
from mylib import Client

client = Client(api_key="your-key")
result = client.process("Hello, World!")
print(result)
# Output: "HELLO, WORLD!"
\`\`\`
```

---

## References

- [Make a README](https://www.makeareadme.com/)
- [Awesome README](https://github.com/matiassingers/awesome-readme)
- [Standard Readme](https://github.com/RichardLitt/standard-readme)
- [Art of README](https://github.com/hackergrrl/art-of-readme)

---

*Last updated: 2026-01-10*
