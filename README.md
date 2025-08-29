# piiredact: Enterprise-Grade PII Redaction for Go

[![Go Reference](https://pkg.go.dev/badge/github.com/rmasci/piiredact.svg)](https://pkg.go.dev/github.com/rmasci/piiredact)
[![Go Report Card](https://goreportcard.com/badge/github.com/rmasci/piiredact)](https://goreportcard.com/report/github.com/rmasci/piiredact)
[![License](https://img.shields.io/github/license/rmasci/piiredact)](https://github.com/rmasci/piiredact/blob/main/LICENSE)

`piiredact` is a high-performance, enterprise-grade library for detecting and redacting Personally Identifiable Information (PII) from text data. It's designed for applications that process sensitive information, such as transcription services, chat logs, or document processing systems.

## Features

- **Comprehensive PII Detection**: Identifies multiple types of sensitive information:
    - Social Security Numbers (SSN)
    - Credit Card Numbers
    - Phone Numbers
    - Bank Routing Numbers (ABA)
    - Driver's License Numbers
    - Email Addresses
    - IP Addresses
    - Passport Numbers
    - Dates of Birth
    - Custom patterns

- **High Accuracy**: Reduces false positives through:
    - Precise regex patterns
    - Validation algorithms (Luhn check for credit cards, checksum for routing numbers)
    - Format validation for structured identifiers

- **High Performance**:
    - Concurrent processing with configurable worker pools
    - Optimized pattern matching
    - Efficient string manipulation

- **Enterprise Features**:
    - Configurable behavior
    - Performance metrics
    - Optional logging
    - Custom pattern support
    - Flexible redaction formatting

## Installation

```bash
go get github.com/rmasci/piiredact
