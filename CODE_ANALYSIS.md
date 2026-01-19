# Tesla Fetch - Code Analysis

## Folder Structure Assessment

### Current Structure
```
tesla-fetch/
├── api/
│   ├── data/
│   ├── types/
│   ├── vehicle-status/
│   └── wake/
├── auth/
├── command/
├── dependencies/
├── draw-status/
├── post-command/
├── services/
│   ├── get-data-handler/
│   └── wake-vehiclel-handler/   <- typo: "vehiclel"
└── ui/
```

### Issues with Current Structure

1. **Typo in folder name**: `wake-vehiclel-handler` should be `wake-vehicle-handler`

2. **Inconsistent naming conventions**:
   - Hyphenated folders: `draw-status`, `vehicle-status`, `post-command`, `get-data-handler`
   - Non-hyphenated: `auth`, `command`, `ui`, `services`
   - Go convention prefers no hyphens in package names

3. **Deeply nested single-file packages**:
   - `services/get-data-handler/` contains only one file
   - `services/wake-vehiclel-handler/` is empty
   - `api/data/`, `api/wake/`, `api/vehicle-status/` each have one file

4. **Unused/dead code**:
   - `dependencies/dependencies.go` imports non-existent `tfetch/vehicle-state`
   - `post-command/post-command.go` calls `auth.CheckLogin(status)` incorrectly (wrong signature)
   - `services/wake-vehiclel-handler/wake-vehicle-handle.go` is empty

5. **Types scattered across multiple packages**:
   - `api/types/` has API response types
   - `auth/auth-types.go` has auth types
   - `command/command-types.go` has command types
   - `draw-status/draw-types.go` has draw types

### Suggested Structure
```
tesla-fetch/
├── cmd/
│   └── tfetch/
│       └── main.go
├── internal/
│   ├── api/
│   │   ├── client.go          # HTTP client wrapper
│   │   ├── vehicle.go         # Vehicle data/state APIs
│   │   └── wake.go
│   ├── auth/
│   │   ├── auth.go
│   │   ├── oauth.go
│   │   └── token.go
│   ├── command/
│   │   └── command.go
│   └── ui/
│       ├── display.go
│       └── spinner.go
├── pkg/
│   └── types/
│       └── types.go           # Shared types
├── go.mod
└── go.sum
```

---

## Code Patterns Analysis

### What Looks Good

1. **Interface-based design** (`auth/auth-types.go:10-15`, `ui/status-logger.go:3-6`):
   - `AuthMethods`, `StatusLoggerMethods`, `VehicleMethods`, `WakeMethods` interfaces enable testing and loose coupling

2. **Encrypted token storage** (`auth/token-store.go`):
   - AES-GCM encryption with salt
   - Proper key derivation
   - Secure file permissions (0600/0700)

3. **Clean dependency injection in main** (`main.go:57-66`):
   - `AppDependencies` struct cleanly wires services

4. **Good use of channels for async UI** (`ui/loading-spinner.go`, `ui/status-logger.go`):
   - Non-blocking spinner updates via channels
   - `NoopLogger` for silent mode

5. **Context with timeout for commands** (`command/handle-command.go:19`):
   - 30-second timeout prevents hanging

6. **Proper defer for cleanup** (multiple files):
   - `defer res.Body.Close()`, `defer close(status)`, `defer listener.Close()`

---

### What Needs Improvement

#### 1. Global state / Package-level variables

**Problem** (`auth/call-auth.go:75-76`, `auth/auth-types.go:51-53`):
```go
var tokenChan = make(chan *Token)
var errChan = make(chan error)
var StateStore string
```
These create race conditions and make testing difficult.

**Fix**: Move these into the `AuthService` struct.

---

#### 2. Environment variable access scattered everywhere

**Problem**: `os.Getenv()` called directly in multiple places:
- `api/data/get-vehicle-data-api.go:15-16`
- `api/vehicle-status/vehicle-state-api.go:17-18`
- `api/wake/wake-api.go:19-20`
- `command/handle-command.go` (implicitly via VIN)

**Fix**: Create a `Config` struct loaded once and passed through dependency injection.

---

#### 3. Inconsistent error handling

**Problem** (`auth/call-auth.go:21-25`):
```go
baseUrl, err := url.Parse("...")
if err != nil {
    log.Fatalf("Malformed auth url: %s", err)  // Exits program
    return nil, err                             // Dead code
}
```
Mixing `log.Fatalf()` (exits) with returning errors creates confusion.

**Problem** (`auth/call-auth.go:117-125`):
```go
if state != StateStore || state == "" {
    http.Error(w, "Internal auth error", http.StatusBadRequest)
    errChan <- fmt.Errorf("...")
    // Missing return! Continues execution
}
```

**Fix**: Choose one pattern - either return errors or use fatal, not both. Add missing `return` statements.

---

#### 4. Missing error handling for `rand.Read`

**Problem** (`auth/call-auth.go:197-199`):
```go
func generateState() string {
    b := make([]byte, 16)
    rand.Read(b)  // Error ignored!
    return base64.URLEncoding.EncodeToString(b)
}
```

**Fix**: Check the error return value.

---

#### 5. HTTP client not reused

**Problem**: Every API call creates a new `http.Client{}`:
```go
client := &http.Client{}  // In vehicle-state-api.go:22
client := &http.Client{}  // In get-vehicle-data-api.go:35
client := &http.Client{}  // In wake-api.go:24
```

**Fix**: Create one shared client with proper timeouts:
```go
var httpClient = &http.Client{Timeout: 30 * time.Second}
```

---

#### 6. Hardcoded magic values

**Problem**:
- `localhost:8080` in `auth/call-auth.go:79`
- `5 * time.Minute` auth timeout in `auth/call-auth.go:102`
- `30 * time.Second` wake timeout in `api/wake/wake-api.go:63`
- `5 * time.Second` poll interval in `api/wake/wake-api.go:64`

**Fix**: Move to constants or config.

---

#### 7. Redundant code in error function

**Problem** (`api/vehicle-status/vehicle-state-api.go:45`):
```go
err = errors.New(fmt.Sprintf("..."))  // Redundant
```

**Fix**: Use `fmt.Errorf()` directly.

---

#### 8. `openBrowser` returns nil even on error

**Problem** (`auth/call-auth.go:202-220`):
```go
func openBrowser(url string) error {
    // ...
    if err != nil {
        fmt.Println("Failed to open browser:", err)
    }
    return nil  // Always returns nil!
}
```

**Fix**: Return the actual error.

---

#### 9. Custom `max` function when builtin exists

**Problem** (`draw-status/draw-status.go:104-109`):
```go
func max(a int, b int) int { ... }
```
Go 1.21+ has builtin `max()`. Since you're on Go 1.22.5, this is unnecessary.

---

#### 10. No tests

**Problem**: Zero test files found.

**Fix**: Add tests, especially for:
- Token encryption/decryption
- API response parsing
- Flag validation

---

#### 11. Broken/dead code

- `dependencies/dependencies.go` imports `tfetch/vehicle-state` which doesn't exist
- `post-command/post-command.go` won't compile - calls `auth.CheckLogin(status)` but signature requires `ui.StatusLoggerMethods`
- `services/wake-vehiclel-handler/wake-vehicle-handle.go` is empty
- Commented-out code in `main.go:79-95`

---

#### 12. No HTTP response status handling for token exchange

**Problem** (`auth/call-auth.go:181-193`):
```go
resp, err := http.Post(tokenUrl, "application/json", bytes.NewBuffer(payload))
// No check of resp.StatusCode!
```

**Fix**: Check for non-200 responses before decoding.

---

## Summary

| Category | Good | Needs Work |
|----------|------|------------|
| Architecture | Interface-based DI | Scattered config, global state |
| Security | AES-GCM encryption, file perms | Ignored rand error |
| Error handling | Some proper patterns | Inconsistent (fatal vs return) |
| Code organization | Logical separation | Too many tiny packages |
| Testing | - | No tests at all |
| Code hygiene | Proper defers | Dead code, typos, commented code |
| HTTP | - | No client reuse, no timeouts |

---

## Priority Fixes

1. Remove dead code (`dependencies/`, `post-command/`, empty handler, commented code)
2. Fix the typo in `wake-vehiclel-handler`
3. Create a shared config loaded once at startup
4. Add HTTP client timeouts and reuse
5. Fix error handling inconsistencies (the missing `return` in callback is a bug)
6. Add tests for critical paths (auth, token storage)
7. Consolidate tiny packages into fewer, larger ones
