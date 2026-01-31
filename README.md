
# Fleming

**Medical history, owned by the individual.  
Science, made auditable.**

Fleming is an open-source DeSci product and protocol that gives people full control over their medical history while enabling privacy-preserving, verifiable scientific research.

It is built in public, with a strong focus on correctness, UX, and long-term credibility.

---

## What is Fleming?

Fleming is a **user-owned medical timeline**.

Instead of fragmented records across hospitals, labs, and systems, Fleming provides a single, coherent view of a person’s medical history — organized over time and connected as a graph.

A timeline may include:
- Diagnoses
- Exams and lab results
- Prescriptions
- Medical documents
- Doctor interactions
- Consent and access history

The user decides what exists, who can see it, and for how long.

Nothing is shared by default.

---

## Core Principles

- **User sovereignty**  
  The patient owns the data, not the platform.

- **Explicit consent**  
  Every access is approved, scoped, and auditable.

- **Privacy by design**  
  All sensitive data is encrypted. The server cannot read it.

- **Honest architecture**  
  Storage is centralized by default, power is not.

- **Science without extraction**  
  Research participation is voluntary and verifiable.

---

## Doctors as Collaborators

Doctors do not own patient data.

They interact with a patient’s timeline by:
- Requesting access
- Uploading signed documents
- Proposing diagnoses or treatments
- Co-signing events

All actions are transparent and cryptographically attributable.

Control always remains with the patient.

---

Fleming does not:
- Recommend treatments
- Rank drugs
- Decide what is true
- Sell data

Instead, it supports science by providing:
- Clear data provenance
- Explicit consent
- Reproducible datasets
- Auditable participation

Users may voluntarily contribute parts of their timeline to research studies under limited, purpose-bound consent. Incentives (when present) reward **data quality and verification**, never medical outcomes.

---

## Security & Storage Model

- Medical data is encrypted using hybrid encryption (symmetric + asymmetric).
- The backend stores encrypted blobs and encrypted keys only.
- Private keys never leave the user’s device or wallet.
- A server compromise does not expose plaintext medical data.
- All access and consent changes are logged and auditable.

Fleming is honest about trust:
- The server is trusted for availability.
- It is **not** trusted with confidentiality.

---

## Identity & Authentication

Fleming uses **wallet-based authentication** instead of passwords.

Users authenticate by signing a challenge with a private key:
- No passwords
- No password resets
- No credential databases
- No third-party auth providers

From a UX perspective, this is presented as a secure, passwordless sign-in — not as a “crypto login”.

---

## Technology Overview

- **Frontend:** React
- **Backend:** Golang (modular monolith)
- **Storage:** Postgres + S3-compatible object storage
- **Encryption:** AES-GCM + public-key cryptography
- **Blockchain:** Solidity contracts for consent anchoring (no medical data on-chain)

Everything is configuration-driven and replaceable.

---

## Status

Fleming is in early development.

Expect:
- Breaking changes
- Incomplete features
- Design discussions
- Iteration in public

This is intentional.

---

## Philosophy

> Fleming does not try to remove trust.  
> It tries to make trust **explicit, minimal, and auditable**.

That is the foundation for both medical data ownership and credible decentralized science.

---

## License

Open source. License to be defined.
