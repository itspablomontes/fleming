---
trigger: always
description: Core project philosophy and engineering principles. Always active.
---

# Fleming Project Rules – General Principles

> **Fleming is a real, open-source, production-grade project. Treat every contribution as if it will be deployed to thousands of users tomorrow.**

## 1. Project Philosophy

### This is NOT a Toy Project
- Fleming is designed for real-world use: people will trust it with their health data.
- Low effort, "just make it work" solutions are **unacceptable**.
- Every change must be **production-ready**, **secure**, and **maintainable**.
- We prioritize **correctness** over speed. Bugs in production are more expensive than time spent getting it right.

### Data Sovereignty & Privacy
- Users own their data. We never store unencrypted PII on our servers.
- Treat every data flow as a potential security audit target.
- When in doubt, **don't store it**.

### Open Source Ethos
- Code should be **readable by strangers**. Write as if you're explaining to a junior developer who will maintain this in 5 years.
- Document "why", not just "what".
- Prefer boring, battle-tested solutions over clever hacks.

### The Fleming Protocol (Mental Model)

> **The Protocol is the foundation. Applications are built on top of it.**

```
┌─────────────────────────────────────────────────┐
│           Applications (apps/)                  │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────┐  │
│  │   Backend   │  │     Web     │  │ Future  │  │
│  │  (Go API)   │  │   (React)   │  │ (Mobile)│  │
│  └──────┬──────┘  └──────┬──────┘  └────┬────┘  │
│         │                │              │       │
│         ▼                ▼              ▼       │
│  ┌─────────────────────────────────────────────┐│
│  │         Fleming Protocol (pkg/protocol)     ││
│  │  • Data Types & Schemas                     ││
│  │  • Encryption Primitives                    ││
│  │  • Timeline Event Definitions               ││
│  │  • Consent Semantics                        ││
│  └─────────────────────────────────────────────┘│
└─────────────────────────────────────────────────┘
```

| Layer | Location | Responsibility |
|:---|:---|:---|
| **Protocol** | `pkg/protocol/` | Source of truth for data shapes, encryption, consent |
| **Backend** | `apps/backend/` | HTTP API that implements the protocol |
| **Web** | `apps/web/` | UI that consumes the protocol via API |

**Implications**:
1. **Protocol types are canonical**: Go types in `pkg/protocol/` define what data looks like. TypeScript types mirror them.
2. **Apps depend on protocol, not vice versa**: Never import from `apps/` into `pkg/`.
3. **Third-party compatibility**: Future apps (mobile, CLI, external tools) implement the same protocol.
4. **Schema evolution**: Protocol changes require migration plans for all consumers.

---

## 2. Critical Thinking Requirements

### Before Writing Any Code
1. **Understand the existing code**: Read it carefully. What exists? Why does it exist?
2. **Clarify the requirement**: If the task is ambiguous, ask clarifying questions.
3. **Consider edge cases**: What happens on empty input? Invalid data? Network failure?
4. **Think about security**: Could this introduce a vulnerability?
5. **Evaluate alternatives**: Is there a simpler way? A library that already does this?

### While Writing Code
1. **Justify your decisions**: Be ready to explain why you chose one approach over another.
2. **Review your own work**: Before submitting, read your code as if someone else wrote it.
3. **Test thoroughly**: Automated tests are mandatory for business logic.

### Red Flags to Watch For
- "This works on my machine" → Not acceptable.
- "I'll add tests later" → Tests come **with** the code.
- "This is just a quick fix" → Quick fixes become permanent. Fix it properly.
- "Nobody will do that" → Someone will. Handle it.

---

## 3. Software Engineering Principles (Non-Negotiable)

### SOLID
| Principle | Application |
|-----------|-------------|
| **S**ingle Responsibility | One function / class / module = one job. |
| **O**pen/Closed | Extend via abstraction, not modification. Use composition. |
| **L**iskov Substitution | Subtypes must be substitutable for their base types. |
| **I**nterface Segregation | Depend on minimal, focused interfaces. |
| **D**ependency Inversion | Depend on abstractions, not concretions. Inject dependencies. |

### DRY (Don't Repeat Yourself)
- If you write the same logic twice, extract it.
- But don't over-abstract. If two things are similar today, they might diverge tomorrow.

### KISS (Keep It Simple, Stupid)
- The simplest solution that works is usually the best.
- Complexity is a liability. Every line of code is a potential bug.

### YAGNI (You Aren't Gonna Need It)
- Don't build features "for the future".
- Build what's needed *now*, make it extensible if cheap.

### Clean Code
- **Meaningful names**: `startTime` not `st`, `userRepository` not `ur`.
- **Small functions**: < 30 lines. One level of abstraction per function.
- **No comments for bad code**: Rewrite to be self-explanatory.
- **No dead code**: Delete it. Git remembers.

### Design Patterns
- Use patterns when they **simplify** the problem, not to show off.
- Common useful patterns: Repository, Service, Factory, Strategy, Observer.
- Avoid over-engineering: Not everything needs a pattern.

---

## 4. Agentic Development Best Practices

> These rules optimize AI-assisted development for safety and quality.

### Plan Before Executing
1. **Create a plan**: Before touching code, outline what you will change and why.
2. **Get approval**: Complex changes require user review before implementation.
3. **Break down tasks**: Large changes should be split into smaller, reviewable chunks.

### Verify Before Completing
1. **Run tests**: Every change must pass existing tests.
2. **Check for regressions**: Did something break that was working before?
3. **Verify in browser/runtime**: Don't just trust the compiler.

### Avoid Common AI Pitfalls
- **Don't hallucinate**: If you're unsure about an API or fact, verify it.
- **Don't over-generate**: Minimal changes are safer. Don't refactor the world.
- **Don't ignore errors**: Every error message is a clue. Investigate.
- **Don't skip context**: Read surrounding code, imports, and tests before editing.

### Communication
- **Be explicit**: State your assumptions clearly.
- **Show your reasoning**: Explain *why* you made a choice.
- **Ask when stuck**: It's better to ask than to guess wrong.

---

## 5. Security & Quality Mindset

### Every Change Could Be the One
- The one that introduces a vulnerability.
- The one that leaks user data.
- The one that crashes production.

### Security First
- Never trust user input. Validate everything.
- Never log secrets, tokens, or PII.
- Use parameterized queries. Never string-concat SQL.
- Secrets come from environment variables, never hardcoded.

### Quality Metrics
- Code coverage: 80%+ on business logic.
- Zero lint errors at commit time.
- Zero `any` in TypeScript. Zero `//nolint:` without justification in Go.

---

## 6. Collaboration & Communication

### Commit Messages
Follow Conventional Commits:
```
feat: add user profile endpoint
fix: handle null avatar in timeline card
docs: update API documentation
refactor: extract auth middleware
test: add integration tests for timeline
chore: upgrade dependencies
```

### Pull Request Guidelines
- Title: Clear, concise description of the change.
- Description: What changed? Why? Any risks?
- Self-review: Show you reviewed your own code.
- Link issues: Reference related GitHub issues.

### Code Review Ethos
- Be respectful. Critique code, not people.
- Explain your reasoning. "Don't do this" is unhelpful.
- Approve with trust, request changes with clarity.

---

## Quick Mantra

> **"Correct, Secure, Readable, Then Fast."**

Build it right. Build it safe. Build it clear. Optimize only when proven necessary.
