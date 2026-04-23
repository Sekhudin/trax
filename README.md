# Trax

**Automate your TypeScript/JavaScript project workflows with Go.**

Trax is a high-performance CLI tool designed to eliminate manual boilerplate by automating route discovery and generating type-safe helpers for your web projects.

---

### ⚠️ Alpha Warning

**Trax is currently in early ALPHA.** This project is experimental and under active development. Automated tests are not yet fully implemented. **Do not use this in production environments** as breaking changes may occur frequently.

---

## Features

- **Route Discovery**: Automatically scans and maps your project structure.
- **Multi-Language Support**: Generates artifacts for both **TypeScript** and **JavaScript**.
- **Type-Safe Helpers**: Create route helpers that prevent broken links and manual string concatenation.
- **Lightning Fast**: Built with Go for near-instant execution.

---

## Quick Start & Examples

Explore our official example repository to see how Trax handles complex project layouts:

**[trax-example](https://github.com/sekhudin/trax-example)**

The example provides a complete reference for configuration, folder scanning, and artifact generation.

---

## Installation

Install the latest version using the Go toolchain:

```bash
go install github.com/sekhudin/trax@latest
```

**Note**: Ensure your `$GOPATH/bin` is in your system's PATH to run the `trax` command globally.

---

## Project Initialization

Generate a default configuration file in your project root:

```bash
trax g config
```

This will create a `trax.toml` (default) file where you can define your routing strategies. Trax also supports `.json` and `.yaml` formats.

---

### Usage

To see all available commands and flags:

```
trax -h
```

---

# Support

If you find this repository useful, you can support the development.

<p align"center">
  <a href="https://trakteer.id/syaikhu" target="_blank">
    <img id="wse-buttons-preview" src="https://edge-cdn.trakteer.id/images/embed/trbtn-red-1.png?v=14-05-2025" height="40" style="border:0px;height:40px;" alt="Trakteer Saya">
  </a>
</p>

Support is completely optional but always appreciated.

---
