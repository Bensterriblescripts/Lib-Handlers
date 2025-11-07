# 🪟 Windows Utility Library

[![Go Version](https://img.shields.io/badge/Go-%3E%3D1.22-blue?logo=go)](https://golang.org/)
[![Platform](https://img.shields.io/badge/Platform-Windows-lightgrey?logo=windows)]()
[![License](https://img.shields.io/badge/License-MIT-green.svg)]()
[![Status](https://img.shields.io/badge/Status-Experimental-orange)]()

> This library is **Windows-only**. It’s heavily tailored to my personal needs — but if you’re brave enough to use it, here’s what you should know before giving up after the first compiler error.

---

This library provides helper functions similar to what you’d find in other standard libraries, plus some shortcuts to simplify Go’s often verbose error handling.

It’s intended for use in personal or internal automation tools on Windows environments, with built-in logging and configuration helpers.

---

Right after entering `main()`, declare your required variables and call `InitLogs()`:

```go
package main

func main() {
    AppName = "YourApp"
    ExecutableName = "yourexe"
    InitLogs()
}
```
