# DefiFundr: Issue and Pull Request Guidelines

## Overview

This document outlines the guidelines for creating issues and pull requests (PRs) in DefiFundr repositories. These guidelines help maintain consistency, clarity, and efficiency in our development workflow.

## Table of Contents

- [Issue Guidelines](#issue-guidelines)
  - [Issue Types](#issue-types)
  - [Issue Creation](#issue-creation)
  - [Issue Lifecycle](#issue-lifecycle)
- [Pull Request Guidelines](#pull-request-guidelines)
  - [Branch Naming](#branch-naming)
  - [PR Creation](#pr-creation)
  - [PR Approval Process](#pr-approval-process)
- [Commit Message Guidelines](#commit-message-guidelines)
- [Code Review Guidelines](#code-review-guidelines)
- [Release Process](#release-process)

---

## Issue Guidelines

### Issue Types

DefiFundr uses the following issue types:

1. **Bug Report** (`bug`) - Issues reporting defects or unexpected behavior in existing functionality
2. **Feature Request** (`feature`) - New features or enhancements to existing functionality
3. **Security Issue** (`security`) - Security vulnerabilities or concerns
4. **Documentation** (`docs`) - Documentation improvements or corrections
5. **Technical Debt** (`tech-debt`) - Code refactoring, test improvements, performance optimizations
6. **Design** (`design`) - UI/UX design-related issues

### Issue Creation

When creating an issue:

1. **Use Templates** - Select the appropriate issue template
2. **Clear Title** - Use a concise, descriptive title following the format: `[TYPE]: Brief description`
3. **Complete Information** - Fill out all required fields in the template
4. **Reproducibility** - For bugs, provide clear reproduction steps
5. **Acceptance Criteria** - Define what "done" looks like for this issue
6. **Labels** - Add appropriate labels that categorize the issue
7. **Priority** - Indicate priority (P0-Critical, P1-High, P2-Medium, P3-Low)
8. **Assignees** - Leave unassigned unless you're working on it immediately

#### Example Bug Report Title:
```
[BUG]: Wallet connection fails on Safari mobile browsers
```

#### Example Feature Request Title:
```
[FEATURE]: Add multi-signature wallet support for DAO funding proposals
```

### Issue Lifecycle

Issues follow this general workflow:

1. **Triage** (`needs-triage`) - Initial assessment
2. **Backlog** - Prioritized but not actively worked on
3. **Ready for Development** (`ready-for-dev`) - Well-defined and ready to be worked on
4. **In Progress** (`in-progress`) - Being actively worked on
5. **Review** (`in-review`) - Associated PR under review
6. **Done** - Issue resolved and closed

---

## Pull Request Guidelines

### Branch Naming

Branch names should follow this convention:
```
<type>/<issue-number>-<brief-description>
```

Where `<type>` is one of:
- `feat` - New feature
- `fix` - Bug fix
- `docs` - Documentation changes
- `style` - Formatting, missing semicolons, etc; no code change
- `refactor` - Code refactoring
- `test` - Adding or refactoring tests; no production code change
- `chore` - Updating build tasks, package manager configs, etc; no production code change
- `perf` - Performance improvements

#### Examples:
```
feat/123-add-multi-sig-support
fix/456-safari-wallet-connection
docs/789-update-api-documentation
```

### PR Creation

When creating a pull request:

1. **Reference Issues** - Link related issues with keywords (`Fixes #123`, `Resolves #456`)
2. **Clear Title** - Use a concise, descriptive title matching commit message format
3. **Complete Description** - Follow the PR template and provide comprehensive information
4. **Size Limit** - Keep PRs focused and limited to ~300-500 lines of code when possible
5. **Tests** - Include appropriate tests for new functionality
6. **Documentation** - Update relevant documentation
7. **CI Checks** - Ensure all CI checks pass before requesting review

### PR Approval Process

1. **Reviewers** - Minimum of 2 approvals required
2. **Code Owner Review** - Required for core components (designated by CODEOWNERS file)
3. **CI Verification** - All automated checks must pass
4. **Merge Strategy** - We use "Squash and merge" for most PRs
5. **Branch Cleanup** - Delete branches after merging

---

## Commit Message Guidelines

DefiFundr follows the [Conventional Commits](https://www.conventionalcommits.org/) specification for commit messages:

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

Where `<type>` follows the same conventions as branch types.

#### Example Commit Messages:
```
feat(wallet): add support for multi-signature wallets

This implementation adds support for multi-signature wallets in the funding 
proposal flow, allowing DAOs to require multiple signers for transactions.

Resolves #123
```

```
fix(ui): correct wallet connection issue on Safari mobile

Resolves #456
```

---

## Code Review Guidelines

### For Authors

1. **Self-Review** - Review your own code before submitting
2. **Description** - Explain your implementation approach
3. **Test Cases** - Detail what was tested and how
4. **Concerns** - Highlight areas where you'd like specific feedback

### For Reviewers

1. **Timeliness** - Review PRs within 1 business day when possible
2. **Constructive Feedback** - Be specific, constructive, and kind
3. **Beyond Style** - Focus on logic, security, and architecture
4. **Review Scope** - Review the entire PR, not just changed lines
5. **Approval Standards** - Only approve if you're confident in the changes

---

## Release Process

1. **Versioning** - We follow [Semantic Versioning](https://semver.org/)
2. **Release Candidates** - Tagged with `-rc` suffix for testing
3. **Release Notes** - Generated from conventional commits
4. **Hotfixes** - Created from the main branch and merged back

---

## Continuous Improvement

These guidelines are not set in stone. If you have suggestions for improving our workflow, please create an issue with the `process` label.

## Security Policy

For security issues, please follow our [Security Policy](SECURITY.md) and report vulnerabilities privately rather than as public issues.