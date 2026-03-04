# GO Patterns

[![Go Report Card](https://goreportcard.com/badge/github.com/DucTran999/go-patterns)](https://goreportcard.com/report/github.com/DucTran999/go-patterns)
[![Go](https://img.shields.io/badge/Go-1.24.5-blue?logo=go)](https://golang.org)
[![codecov](https://codecov.io/gh/DucTran999/go-patterns/graph/badge.svg?token=5XBMMBKCPD)](https://codecov.io/gh/DucTran999/go-patterns)
[![License](https://img.shields.io/github/license/DucTran999/go-patterns)](LICENSE)

This repository serves as a comprehensive collection of design patterns and best practices in Go, covering a wide range of topics including concurrency, behavioral patterns, creational patterns, and resilience strategies. Each pattern is implemented with clear examples and explanations to help developers understand when and how to use them effectively in their Go projects.

---

## Project Structure

The repository is organized into the following directories:

- `structural/`: Illustrates Go's structural patterns, including Adapter, Decorator, and Composite.

- `behavioral/`: Implements common behavioral design patterns such as Strategy, Bridge, and Command.

- `creational/`: Contains examples of creational design patterns like Factory, Singleton, and Builder.

- `concurrency/`: Demonstrates advanced concurrency patterns such as Worker Pools, Pipelines, and Fan-in/Fan-out for efficient task orchestration.

- `resilience/`: Features fault-tolerance patterns including Circuit Breaker, Retry, and Bulkhead to ensure system stability and graceful degradation.

## Prerequisites

Ensure the following tools are installed on your machine:

- [**Go 1.25+**](https://go.dev/dl/) — The project requires Go version 1.25 or later.
- [**Taskfile CLI**](https://taskfile.dev/) — Used for task automation and scripting.

---

## Installation

Clone the repository:

```bash
git clone https://github.com/ductran999/go-patterns.git
cd go-patterns
```

---

## 🧪 Testing

This project uses Go's built-in testing framework with mocks and table-driven tests.

```bash
go test -v ./...
```

---

## 📄 License

This project is licensed under the [MIT License](LICENSE).

---

## 🤝 Contributing

Contributions are welcome! If you have suggestions for improvements or new patterns to include, please open an issue or submit a pull request.
