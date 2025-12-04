# Contributing to DBX – Open-Source Database Backup & Restore CLI

Thank you for your interest in contributing to **DBX**, an open-source, cross-platform tool for **database backups, restores, scheduling, compression, and cloud uploads**. Community contributions help DBX grow into a reliable DevOps-grade backup utility for MySQL, PostgreSQL, MongoDB, and SQLite.

This guide explains how to contribute code, documentation, bug reports, and feature suggestions.

---

## 🚀 How to Contribute

### 1. Fork the repository
```bash
git fork https://github.com/zfhassaan/dbx

```

### 2. Create a New Branch
```bash
git checkout -b feature/my-update
```

### 3. Make your changes
Follow clean Go coding practices and keep PRs focused and atomic.

Relevant areas include:
* Backup/restore improvements
* Cloud storage integrations (AWS S3, GCS, Azure)
* Logging, error handling, and scheduling enhancements
* Performance optimizations
* CLI improvements and new database features

### 4. Run tests
```bash
make test
```

### 5. Submit a Pull Request

Explain:
* What you changed
* Why the change is needed
* Any technical details or limitations

PRs that improve backup reliability, compression efficiency, or cloud upload stability are especially welcomed.

### 6. Testing Requirements

All new code must:
* Include unit tests
* Pass all existing tests
* Maintain or improve test coverage

DBX ships with 110+ tests — keep the suite healthy.

### Code Style Guidelines

* Prefer readability over clever tricks
* Avoid unnecessary dependencies
* Keep commits clean and well-described
* Follow idiomatic Go patterns


### Need Help?

Open an issue or discussion on GitHub.
We’re happy to guide new contributors as long as questions are clear and technical.

Thanks for helping improve DBX — your contribution strengthens the open-source DevOps ecosystem!
