# Security Audit Report

This report captures the findings from the API security review and the actions applied to close the highest-risk issues.

## 1) Password hashes exposed in API responses
- Severity: High
- Status: Fixed Now
- Finding: If the user model is serialized directly, password hashes or internal auth state could be returned to clients.
- Fix: Added PasswordHash with `json:"-"` and kept auth data behind isolated DTOs.

## 2) Missing authentication on staff actions
- Severity: High
- Status: Fixed Now
- Finding: Book create/update routes were not protected by token validation, allowing unauthenticated writes.
- Fix: Added `AuthMiddleware` and restricted staff routes to `librarian` and `admin` roles.

## 3) Missing RBAC enforcement
- Severity: High
- Status: Fixed Now
- Finding: Even if a user had a valid session, the app would not limit actions based on role.
- Fix: Added `RequireRole` to enforce `reader`, `librarian`, and `admin` permissions before route execution.

## 4) User enumeration via login responses
- Severity: Medium
- Status: Fixed Now
- Finding: Different error messages for unknown email or wrong password can reveal account existence.
- Fix: Both cases now return the same HTTP 401 payload: `{"error":{"code":"unauthorized","message":"invalid credentials"}}`.

## 5) JWT secret management risk
- Severity: Medium
- Status: Fixed Now
- Finding: JWT signing secret must never be hard-coded in code or committed to source control.
- Fix: Token generation and validation read the signing key strictly from `JWT_SECRET`.

## 6) Sensitive request body leakage in logs
- Severity: Low
- Status: Tracked for Later
- Finding: If future logging code prints entire request bodies, passwords will appear in logs.
- Fix: Keep request logging limited to metadata and avoid body serialization; future logging should redact or omit `password` fields.

## 7) Mass assignment through direct model binding
- Severity: Medium
- Status: Fixed Now
- Finding: Binding the raw JSON to a DB model allows client-controlled fields such as role or password hash to be updated.
- Fix: Introduced dedicated DTOs (`RegisterRequest`, `LoginRequest`, `UpdateProfileRequest`) with restricted fields.

## Summary
The project now closes the main auth and RBAC issues in a production-safe way: password hashing, JWT validation, RBAC enforcement, and generic login failures. Remaining logging hardening is tracked for later review and can be tightened as the codebase grows.
