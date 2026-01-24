# Fleming Project Rules – General Principles

> **Fleming is a real, open-source, production-grade project. Treat every contribution as if it ships to thousands of users tomorrow.**

---

## 0) Non-Negotiables (TL;DR)

- **Correctness > speed**
- **Security > convenience**
- **Readable > clever**
- **Protocol-first architecture**
- **No silent failures**
- **No PII in logs**
- **Tests ship with logic**

> **Mantra:** *Correct. Secure. Readable. Then fast.*

---

## 1) Project Philosophy

### This is NOT a Toy Project
Fleming handles health history and medical context. That means:

- “Just make it work” is **not acceptable**
- Every change must be **production-grade**
- Bugs are **real harm**, not just “oops”

**We optimize for:**
1. **Correctness**
2. **Security**
3. **Maintainability**
4. **Performance (only when measured and needed)**

---

### Data Sovereignty & Privacy (Default: Minimize)
- Users **own** their data.
- We never store **unencrypted PII** on our servers.
- Every data flow must survive a security audit mindset.

**Rule of thumb:**  
> If you don’t need to store it, **don’t store it**.

**Never:**
- log PII, tokens, secrets, raw documents
- store plaintext medical content server-side unless explicitly designed and encrypted

---

### Open Source Ethos
Code must be understandable by strangers.

- Write for a developer who joins **2 years later**
- Document **why**, not only **what**
- Prefer **boring, battle-tested** solutions over hacks

---

## 2) The Fleming Protocol (Mental Model)

> **The Protocol is the foundation. Applications are built on top of it.**

```text
┌──────────────────────────────────────────────────────────────┐
│                      Applications (apps/)                    │
│                                                              │
│   ┌───────────────────┐   ┌───────────────────┐              │
│   │ Backend (Go API)   │   │ Web (React)        │              │
│   │ apps/backend/      │   │ apps/web/          │              │
│   └─────────┬─────────┘   └─────────┬─────────┘              │
│             │                       │                        │
│             ▼                       ▼                        │
│   ┌───────────────────────────────────────────────────────┐  │
│   │               Fleming Protocol (pkg/protocol/)         │  │
│   │                                                       │  │
│   │  • Canonical data models + schemas                     │  │
│   │  • Encryption primitives + key semantics               │  │
│   │  • Timeline event definitions                          │  │
│   │  • Consent + access semantics                          │  │
│   │  • Versioning + migrations strategy                    │  │
│   └───────────────────────────────────────────────────────┘  │
└──────────────────────────────────────────────────────────────┘
```

### Layers & Responsibilities

| Layer        | Location        | Responsibility                                      |
| :----------- | :-------------- | :-------------------------------------------------- |
| **Protocol** | `pkg/protocol/` | Canonical truth for types, schemas, crypto, consent |
| **Backend**  | `apps/backend/` | HTTP API that implements the protocol               |
| **Web**      | `apps/web/`     | UI that consumes the API and enforces UX rules      |

### Protocol Implications (Hard Rules)

#### Protocol types are canonical
Go types in `pkg/protocol/` define reality.
TypeScript types mirror, they don’t invent.

#### Apps depend on protocol — never the reverse
Never import from `apps/` into `pkg/`.

#### Third-party compatibility is intentional
Future clients (mobile, CLI, external tools) must implement the same protocol.

#### Schema evolution requires a migration story
If protocol changes: you must consider existing data + existing clients.

---

## 3) Critical Thinking Requirements

### Before Writing Any Code
1. **Read the existing code**
   - What already exists?
   - Why was it implemented this way?
2. **Clarify the requirement**
   - If ambiguous: ask.
   - If risky: call it out.
3. **Enumerate edge cases**
   - empty input
   - invalid data
   - partial failures
   - retries / timeouts
   - concurrency issues
4. **Threat-model the change**
   - Could this leak data?
   - Could this bypass consent?
   - Could this break encryption guarantees?
5. **Choose the simplest correct solution**
   - Avoid adding dependencies unless justified.

### While Writing Code
1. **Make decisions explainable**
2. **Prefer explicitness over magic**
3. **Handle failures intentionally**
4. **Add tests for business logic**
5. **Avoid breaking changes without a plan**

### Red Flags (Stop Signs)
- “Works on my machine”
- “I’ll add tests later”
- “This is just a quick fix”
- “Nobody will do that”

If you catch yourself thinking these: **pause and redesign.**

---

## 4) Engineering Principles (Non-Negotiable)

### SOLID (Practical Interpretation)
| Principle | Fleming interpretation                                  |
| :-------- | :------------------------------------------------------ |
| **S**     | Each module has one job and one reason to change        |
| **O**     | Add features by extension, not rewriting stable code    |
| **L**     | Swappable implementations must behave consistently      |
| **I**     | Small, focused interfaces beat “god interfaces”         |
| **D**     | Depend on abstractions; inject concrete implementations |

### DRY (But Don’t Over-Abstract)
- **Duplicate logic** → extract
- **Similar logic** → only extract if it stays the same for a reason

### KISS
- Simple is safer.
- Complexity is a permanent cost.

### YAGNI
- Don’t build “future features”.
- Build what is required **now**, with clean extension points only if cheap.

### Clean Code Standards
- **Names must be meaningful** (`userRepository`, not `ur`)
- **Functions should be small and focused**
- **Prefer refactoring over comment excuses**
- **Delete dead code** (git is the archive)

### Design Patterns (Use Only When They Reduce Complexity)
Allowed patterns when justified:
- Repository
- Service
- Factory
- Strategy
- Observer

Avoid “pattern cosplay”.

---

## 5) Agentic Development Rules (AI-Assisted Safety)

> These rules exist to prevent “fast wrong” changes.

### Plan Before Executing
Write a short plan before coding:
- what will change
- why it changes
- risks
- tests

Complex changes require explicit review before implementation.

### Verify Before Completing
1. **Run tests**
2. **Check for regressions**
3. **Validate behavior in runtime** (browser/API)

### Avoid Common AI Failure Modes
- **Don’t guess APIs** — verify them
- **Don’t refactor unrelated code**
- **Don’t ignore error messages**
- **Always read surrounding context before editing**

### Communication Requirements
- **State assumptions explicitly**
- **Explain tradeoffs briefly**
- **Ask when blocked instead of guessing**

---

## 6) Security & Quality Mindset

### Every Change Could Be The One
- the one that leaks data
- the one that breaks consent
- the one that crashes production

**Act accordingly.**

### Security Rules (Hard Requirements)
- **Validate and sanitize all input**
- **Never log secrets / tokens / PII**
- **Use parameterized queries** (no SQL string concat)
- **Secrets must come from env/secret manager** — never hardcode

### Quality Gates
- **80%+ coverage on business logic**
- **Zero lint errors** at commit time
- **TypeScript:** no `any`
- **Go:** no `//nolint` without justification

---

## 7) Collaboration & Communication

### Commit Messages (Conventional Commits)
```
feat: add user profile endpoint
fix: handle null avatar in timeline card
docs: update API documentation
refactor: extract auth middleware
test: add integration tests for timeline
chore: upgrade dependencies
```

### Pull Requests
- Clear title
- Explain **what** + **why**
- Mention risks and migration notes (if any)
- Self-review before requesting review
- Link related issues

### Review Ethos
- Critique code, not people
- Explain reasoning
- Be strict on security and correctness
- Be kind in communication

---

## Quick Mantra
> **Correct, Secure, Readable, Then Fast.**