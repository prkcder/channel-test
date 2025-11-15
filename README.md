# Folder Structure
```
channel-test/
├── cmd/
│   └── scores-api/
│       └── main.go
│
├── internal/
│   ├── api/
│   │   ├── handlers.go
│   │   └── router.go
│   │
│   ├── consumer/
│   │   └── sse.go
│   │
│   └── store/
│       ├── memory.go
│       └── store.go
│
├── pkg/
│   └── models/
│       └── models.go
│
├── Dockerfile
├── docker-compose.yml
├── Makefile
├── go.mod
├── go.sum
└── README.md

```

# Implementation Summary

## Overview

This is a complete implementation of the test scores API take-home test. The application is written entirely in Go using only the standard library, demonstrating idiomatic Go code, clean architecture, and production-ready patterns.
