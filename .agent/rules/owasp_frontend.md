---
trigger: glob
description: "*.ts", "*.tsx"
globs: "*.ts", "*.tsx"
---

# OWASP Frontend Security Rules

> **Trigger**: Use these rules when building, reviewing, or modifying frontend code in `apps/web/`.

> [!IMPORTANT]
> For comprehensive security reviews, invoke the `/cc-skill-security-review` skill which provides deeper vulnerability analysis, OWASP Top 10 mapping, and remediation guidance.

---

## 1. Input Validation & Sanitization

| Rule                                               | Implementation                            |
| -------------------------------------------------- | ----------------------------------------- |
| Never use `innerHTML` or `dangerouslySetInnerHTML` | Use `textContent` or React's JSX escaping |
| Validate all user inputs client-side               | Use Zod schemas before API calls          |
| Sanitize any data displayed from external sources  | Use DOMPurify if HTML rendering required  |

```tsx
// ❌ Dangerous
element.innerHTML = userInput;

// ✅ Safe
element.textContent = userInput;
```

---

## 2. XSS Prevention

| Rule                                                   | Implementation                              |
| ------------------------------------------------------ | ------------------------------------------- |
| Let React handle escaping                              | Never bypass with `dangerouslySetInnerHTML` |
| Use CSP headers                                        | Configure in backend/nginx                  |
| Avoid `eval()`, `new Function()`, `setTimeout(string)` | Always use function references              |

**Content-Security-Policy Header** (configure in backend):
```
Content-Security-Policy: 
  default-src 'self';
  script-src 'self';
  style-src 'self' 'unsafe-inline';
  img-src 'self' data: https:;
  connect-src 'self' https://api.example.com;
```

---

## 3. Authentication Token Storage

| Method              | When to Use           | Risk Level                      |
| ------------------- | --------------------- | ------------------------------- |
| **HttpOnly Cookie** | Preferred for JWTs    | Low (immune to XSS)             |
| **sessionStorage**  | Short-lived tokens    | Medium (cleared on tab close)   |
| **localStorage**    | Avoid for auth tokens | High (persists, XSS-accessible) |

```typescript
// ✅ Preferred: HttpOnly cookie set by backend
// Frontend doesn't handle token storage

// 🟡 Acceptable: sessionStorage for wallet-derived keys
sessionStorage.setItem('kek', derivedKey);
window.addEventListener('beforeunload', () => sessionStorage.clear());
```

---

## 4. Cryptography in Browser

| Rule                          | Implementation                              |
| ----------------------------- | ------------------------------------------- |
| Use WebCrypto API exclusively | Never use npm crypto packages               |
| Never expose keys in logs     | Use `console.log` sparingly, never log keys |
| Clear keys on logout          | Explicitly zero memory where possible       |

```typescript
// ✅ WebCrypto only
const key = await crypto.subtle.generateKey(
  { name: 'AES-GCM', length: 256 },
  true,
  ['encrypt', 'decrypt']
);

// ❌ Never use
import CryptoJS from 'crypto-js'; // npm package
```

---

## 5. Key Management

| Rule                               | Implementation                    |
| ---------------------------------- | --------------------------------- |
| Derive keys from wallet signatures | Use HKDF with WebCrypto           |
| Never persist KEK to storage       | Keep in memory only               |
| Clear on logout                    | Explicitly clear `sessionStorage` |

```typescript
// Key derivation from wallet signature
const signature = await wallet.signMessage('Fleming Key Derivation v1');
const keyMaterial = await crypto.subtle.importKey(
  'raw',
  new TextEncoder().encode(signature),
  'HKDF',
  false,
  ['deriveKey']
);
const kek = await crypto.subtle.deriveKey(
  { name: 'HKDF', salt, info, hash: 'SHA-256' },
  keyMaterial,
  { name: 'AES-GCM', length: 256 },
  false,
  ['wrapKey', 'unwrapKey']
);
```

---

## 6. CORS & API Security

| Rule                                 | Implementation                       |
| ------------------------------------ | ------------------------------------ |
| Strict origin validation             | Backend validates `Origin` header    |
| Include credentials only when needed | `credentials: 'same-origin'` default |
| Use CSRF tokens for mutations        | Or rely on SameSite cookies          |

```typescript
// API client configuration
const api = {
  fetch: (url: string, options: RequestInit = {}) =>
    fetch(url, {
      ...options,
      credentials: 'same-origin',
      headers: {
        'Content-Type': 'application/json',
        ...options.headers,
      },
    }),
};
```

---

## 7. Dependency Security

| Rule                       | Implementation                       |
| -------------------------- | ------------------------------------ |
| Run `pnpm audit` regularly | Add to CI pipeline                   |
| Pin dependency versions    | Use exact versions in `package.json` |
| Review new dependencies    | Check npm advisories before adding   |
| Minimal dependencies       | Prefer stdlib over npm packages      |

```bash
# Regular security audit
pnpm audit

# Update vulnerable packages
pnpm audit --fix
```

---

## 8. Error Handling

| Rule                               | Implementation                         |
| ---------------------------------- | -------------------------------------- |
| Never expose stack traces to users | Catch errors, show generic messages    |
| Log errors to structured logging   | Not to console in production           |
| Validate API error responses       | Don't trust error messages from server |

```typescript
// ❌ Dangerous
catch (error) {
  alert(error.stack);
}

// ✅ Safe
catch (error) {
  console.error('Operation failed:', error); // Internal logging
  toast.error('Something went wrong. Please try again.');
}
```

---

## 9. Transport Security

| Rule                           | Implementation                     |
| ------------------------------ | ---------------------------------- |
| Enforce HTTPS everywhere       | `Strict-Transport-Security` header |
| No mixed content               | All resources over HTTPS           |
| Certificate pinning (optional) | For mobile apps                    |

---

## 10. Secure Defaults

| Setting                | Value                             |
| ---------------------- | --------------------------------- |
| `SameSite` cookie      | `Strict` or `Lax`                 |
| `Secure` cookie flag   | `true`                            |
| `HttpOnly` cookie flag | `true` (for auth tokens)          |
| Referrer-Policy        | `strict-origin-when-cross-origin` |
| X-Content-Type-Options | `nosniff`                         |
| X-Frame-Options        | `DENY`                            |

---

## Security Review Skill

For comprehensive security analysis beyond these rules, use:

```
/cc-skill-security-review
```

This skill provides:
- OWASP Top 10 vulnerability mapping
- Automated code scanning recommendations
- Threat modeling guidance
- Remediation priorities
- Compliance checking (HIPAA, GDPR considerations)

---

## Checklist for PRs

Before merging frontend changes, verify:

- [ ] No `innerHTML` or `dangerouslySetInnerHTML`
- [ ] No hardcoded secrets or API keys
- [ ] All user inputs validated
- [ ] Crypto uses WebCrypto API only
- [ ] Keys cleared on logout
- [ ] No console.log with sensitive data
- [ ] `pnpm audit` passes
- [ ] CSP headers configured
