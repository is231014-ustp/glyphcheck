# Glyphcheck
Glyphcheck is a prototype CI/CD security tool developed as part of a bachelor's thesis.
It detects potentially dangerous Unicode characters in source code to help prevent homoglyph attacks.

**This project is currently a Proof of Concept / MVP and is not production-ready!!!**

## Overview
Modern programming languages allow Unicode characters in identifiers and strings. While useful, this also enables attackers to introduce visually similar characters (homoglyphs) that can hide malicious code.

Glyphcheck scans files and flags:

* Disallowed characters — characters outside your configured allow-list

* Suspicious characters — known problematic glyphs from the allowed script often used in attacks

It is designed to be integrated into CI/CD pipelines so repositories can automatically reject unsafe characters before merge or release.

## Features

* Allow-list based validation

* Configurable scripts, categories, and characters

* Suspicious glyph detection [WIP]

* Unicode normalization (NFC)

* GitHub Actions-compatible error output format

* Recursive directory scanning

* Extension filtering

* Directory exclusion support

## Current Status

| Feature                           | Status                 |
| --------------------------------- | ---------------------- |
| Mixed-script detection            | ✅ Implemented          |
| Single-script homoglyph detection | ⚠️ Hardcoded prototype |
| Custom suspicious glyph list      | ✅                      |
| Config file support               | ✅                      |
| Production stability              | ❌ Not yet              |

## Installation
Requires Go ≥ 1.24

```
git clone <repo>
cd glyphcheck
go build -o glyphcheck .
```

## Usage

Run in repository root:
```
./glyphcheck
```
The tool looks for:
```
.glyphcheck.yaml
```
If the file is missing, defaults are used. It is highly recommended to use the config file to adjust to your specific needs.

Exit codes:

| Code | Meaning                      |
| ---- | ---------------------------- |
| 0    | No violations                |
| 1    | Violations or runtime errors |

## Example CI Integration (GitHub Actions)
```
- name: Build glyphcheck
  run: go build -o glyphcheck .

- name: Run glyphcheck
  run: ./glyphcheck
```

Errors are printed in GitHub Actions annotation format:

```
::error file=main.go,line=5,col=10::disallowed unicode character U+01C3 (ǃ)
```

## Configuration

Example .glyphcheck.yaml:

```
root: .

allowed_extensions:
  - .go

excluded_directories:
  - .git

allowed:
  scripts:
    - Latin
    - Common

  categories:
    - Letter
    - Number
    - Punct

  characters:
    - "\n"
    - "\r"
    - "\t"
    - " "
    - "="
    - "+"
    - "-"
    - "`"
    - "<"
    - ">"

suspicious:
  enabled: true
  characters:
    - "ᴀ"
```

### Config Fields

```root``` - Directory to scan recursively.

```allowed_extensions``` - File extensions to analyze.

```excluded_directories``` - Directories skipped during scan.

```allowed.scripts``` - Unicode scripts allowed (currently supported):

* ```Latin```
* ```Common```

```allowed.categories``` - Unicode categories allowed:

* ```Letter```
* ```Number```
* ```Punct```

```allowed.characters``` - Explicitly allowed characters regardless of script/category.

```suspicious.enabled``` - Enable warning mode for suspicious glyphs.

```suspicious.characters``` - Additional characters to flag as suspicious.

## How Detection Works

1. Files are read line-by-line

2. Input is normalized using Unicode NFC

3. Each rune is checked:
    * If not in allow-list → error
    * If suspicious → warning

4. Any error causes non-zero exit status

## Security Model

Glyphcheck follows an allow-list security model:

Only explicitly allowed Unicode scripts, categories, and characters are permitted.

This is more secure than deny-list approaches because new Unicode homoglyphs cannot bypass the filter.

## Limitations

Current prototype limitations:

* Only two Unicode scripts supported in config parsing
* Limited category support
* Single-script detection not fully implemented
* Suspicious list is partially hardcoded
* No performance optimization yet
* No incremental scanning
* No IDE integration

**Detected suspicious characters currently do not fail the pipeline (they do not affect the exit code).** 

## Roadmap

Planned improvements:

* Full Unicode script support
* Advanced homoglyph similarity detection
* Proper single-script analysis
* Performance optimization for large repos
* Adding option in config to fail pipeline if suspicious characters are detected

## Motivation

Unicode security issues are increasingly relevant in:

* supply-chain attacks
* open-source contributions
* dependency spoofing
* typosquatting
* code review bypasses

Glyphcheck demonstrates a practical CI-enforced defense approach.

## Dependencies

* ```golang.org/x/text``` — Unicode normalization
* ```gopkg.in/yaml.v3``` — configuration parsing

## License

------------------!!!!!ADD LICENSE DON'T FORGET!!!!!------------------
