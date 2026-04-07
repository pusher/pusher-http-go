# pusher-http-go v6 Modernization Plan

## Executive Summary

This document outlines the plan to modernize `pusher-http-go` to v6.0.0, targeting **Go 1.26+** with breaking changes to remove deprecated features, add context support, and eliminate code smells.

**Current State:**
- Go version: 1.14
- Last significant update: 2021 (v5.1.0)
- Contains deprecated `io/ioutil` usage
- Missing `context.Context` support
- Contains deprecated methods scheduled for removal

---

## 1. High Priority: Breaking Changes (v6.0.0)

### 1.1 Update Go Version

| File | Change |
|------|--------|
| `go.mod:3` | `go 1.14` → `go 1.26` |
| `.github/workflows/test.yml:14` | Test matrix `['1.18', '1.19']` → `['1.25', '1.26']` |
| `README.md:13` | Update "Supports Go 1.11 or greater" → "Supports Go 1.26 or greater" |

### 1.2 Add Context Support

Add `context.Context` parameter to all HTTP methods for proper cancellation/timeout support:

| File | Method | New Signature |
|------|--------|---------------|
| `client.go` | `Trigger` | `Trigger(ctx context.Context, channel, eventName string, data interface{}) error` |
| `client.go` | `TriggerWithParams` | Add ctx first param |
| `client.go` | `TriggerMulti` | Add ctx first param |
| `client.go` | `TriggerMultiWithParams` | Add ctx first param |
| `client.go` | `TriggerBatch` | Add ctx first param |
| `client.go` | `SendToUser` | Add ctx first param |
| `client.go` | `Channels` | Add ctx first param |
| `client.go` | `Channel` | Add ctx first param |
| `client.go` | `GetChannelUsers` | Add ctx first param |
| `request.go:24` | `http.NewRequest` → `http.NewRequestWithContext` |

**Example migration:**

```go
// Before (v5)
err := client.Trigger("my-channel", "my-event", data)

// After (v6)
err := client.Trigger(context.Background(), "my-channel", "my-event", data)

// With timeout
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
err := client.Trigger(ctx, "my-channel", "my-event", data)
```

### 1.3 Remove Deprecated Methods

| File | Method/Field | Status | Replacement |
|------|--------------|--------|-------------|
| `client.go:233-237` | `TriggerExclusive()` | **Remove** | `TriggerWithParams()` with `SocketID` |
| `client.go:248-252` | `TriggerMultiExclusive()` | **Remove** | `TriggerMultiWithParams()` with `SocketID` |
| `client.go:584-586` | `AuthenticatePrivateChannel()` | **Remove** | `AuthorizePrivateChannel()` |
| `client.go:624-626` | `AuthenticatePresenceChannel()` | **Remove** | `AuthorizePresenceChannel()` |
| `client.go:58` | `EncryptionMasterKey string` | **Remove** | `EncryptionMasterKeyBase64` |
| `client.go:730-742` | Logic for deprecated `EncryptionMasterKey` | **Remove** | N/A |

---

## 2. Medium Priority: Deprecated Package Updates

### 2.1 Replace `io/ioutil` (deprecated since Go 1.16)

| File | Line | Current | Replacement |
|------|------|---------|-------------|
| `request.go:7` | `import "io/ioutil"` | `import "io"` |
| `request.go:39` | `ioutil.ReadAll(response.Body)` | `io.ReadAll(response.Body)` |
| `client_test.go:6` | `import "io/ioutil"` | `import "io"` |
| `client_test.go:363,390,428,469,490` | `ioutil.ReadAll(req.Body)` | `io.ReadAll(req.Body)` |

### 2.2 Update Documentation Examples

| File | Lines | Update |
|------|-------|--------|
| `client.go:501,559,596,690` | Doc examples with `ioutil.ReadAll` | Use `io.ReadAll` |
| `README.md:420,463,536,738` | Examples with `ioutil.ReadAll` | Use `io.ReadAll` |

---

## 3. Medium Priority: Code Smells & Bug Fixes

### 3.1 Missing Error Check (Bug)

**File:** `request.go:24-28`

**Issue:** `http.NewRequest` error is not checked before using `req`.

**Current:**
```go
req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
for key, val := range headers {
    req.Header.Set(http.CanonicalHeaderKey(key), val)
}
```

**Fixed:**
```go
req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(body))
if err != nil {
    return nil, err
}
for key, val := range headers {
    req.Header.Set(http.CanonicalHeaderKey(key), val)
}
```

### 3.2 Replace Panics with Errors

Library code should never panic. These should return errors instead.

| File | Line | Function | Current | Fixed |
|------|------|----------|---------|-------|
| `crypto.go:67-78` | `formatMessage()` | `panic(err)` | Return `(string, error)` |
| `crypto.go:80-87` | `generateNonce()` | `panic(err)` | Return `([24]byte, error)` |

**Current:**
```go
func formatMessage(nonce string, cipherText string) string {
    // ...
    json, err := json.Marshal(encryptedMessage)
    if err != nil {
        panic(err)
    }
    return string(json)
}
```

**Fixed:**
```go
func formatMessage(nonce string, cipherText string) (string, error) {
    // ...
    json, err := json.Marshal(encryptedMessage)
    if err != nil {
        return "", fmt.Errorf("failed to marshal encrypted message: %w", err)
    }
    return string(json), nil
}
```

### 3.3 Replace `errors.New(fmt.Sprintf(...))` with `fmt.Errorf`

| File | Line | Current | Fixed |
|------|------|---------|-------|
| `encoder.go:50` | `errors.New(fmt.Sprintf("Event payload exceeded maximum size (%d bytes is too much)", len(payloadData)))` | `fmt.Errorf("event payload exceeded maximum size (%d bytes is too much)", len(payloadData))` |
| `encoder.go:59` | `errors.New(fmt.Sprintf("Paramater %s specified multiple times", k))` | `fmt.Errorf("parameter %s specified multiple times", k)` |

### 3.4 Fix Typo

| File | Line | Current | Fixed |
|------|------|---------|-------|
| `encoder.go:59` | `"Paramater"` | `"Parameter"` |

### 3.5 Simplify Boolean Returns

| File | Function | Current | Simplified |
|------|----------|---------|------------|
| `util.go:61-66` | `validChannel()` | `if ... { return false } return true` | `return len(channel) <= maxChannelNameSize && channelValidationRegex.MatchString(channel)` |
| `util.go:77-82` | `isEncryptedChannel()` | `if ... { return true } return false` | `return strings.HasPrefix(channel, "private-encrypted-")` |

---

## 4. Low Priority: Dependency Updates

### 4.1 Update Dependencies

| Dependency | Current | Update To |
|------------|---------|-----------|
| `golang.org/x/crypto` | `v0.0.0-20200709230013-948cd5f35899` | `v0.49.0` |
| `gopkg.in/stretchr/testify.v1` | `v1.2.2` | `github.com/stretchr/testify v1.11.1` |

**Note:** The testify import path changed from `gopkg.in/stretchr/testify.v1` to `github.com/stretchr/testify`. All test files need updating.

### 4.2 Update GitHub Actions

| File | Item | Current | Update To |
|------|------|---------|-----------|
| `.github/workflows/test.yml:10` | Runner | `ubuntu-20.04` | `ubuntu-latest` |
| `.github/workflows/test.yml:21` | Checkout | `actions/checkout@v2` | `actions/checkout@v4` |
| `.github/workflows/test.yml:24` | Setup Go | `actions/setup-go@v2` | `actions/setup-go@v5` |
| `.github/workflows/test.yml:40` | Runner (finish job) | `ubuntu-20.04` | `ubuntu-latest` |

---

## 5. File-by-File Change Summary

| File | Changes Required |
|------|------------------|
| `go.mod` | Go version 1.22, update dependencies |
| `client.go` | Add context params, remove deprecated methods/fields, update doc examples |
| `request.go` | Add context support, replace ioutil→io, add error check |
| `crypto.go` | Replace panic with error returns |
| `encoder.go` | Use fmt.Errorf, fix typo |
| `util.go` | Simplify boolean return statements |
| `client_test.go` | Replace ioutil→io, update tests for context params |
| `channel_authentication_test.go` | Update testify import |
| `crypto_test.go` | Update testify import |
| `request_url_test.go` | Update testify import |
| `response_parsing_test.go` | Update testify import |
| `util_test.go` | Update testify import |
| `webhook_test.go` | Update testify import |
| `README.md` | Update Go version requirement, update examples |
| `CHANGELOG.md` | Add v6.0.0 migration guide |
| `.github/workflows/test.yml` | Update Go versions, action versions |

---

## 6. Implementation Order

### Phase 1: Dependencies & Infrastructure
- [x] Update `go.mod` to Go 1.26
- [x] Update `golang.org/x/crypto` to v0.49.0
- [x] Update testify to `github.com/stretchr/testify v1.11.1`
- [x] Update all test file imports for testify
- [x] Update GitHub Actions workflow
- [x] Run `go mod tidy`

### Phase 2: Bug Fixes & Code Smells
- [x] Add missing error check in `request.go`
- [x] Replace panics with errors in `crypto.go`
- [x] Update callers of `formatMessage()` and `generateNonce()`
- [x] Fix typo and use fmt.Errorf in `encoder.go`
- [x] Simplify boolean returns in `util.go`

### Phase 3: ioutil Removal
- [x] Replace all `ioutil.ReadAll` with `io.ReadAll` in source files
- [x] Update documentation examples in `client.go`
- [x] Update README.md examples

### Phase 4: Breaking API Changes
- [x] Add `context.Context` parameter to `request()` function
- [x] Add `context.Context` to all public HTTP methods
- [x] Update all internal callers
- [x] Remove deprecated methods (`TriggerExclusive`, etc.)
- [x] Remove deprecated `EncryptionMasterKey` field and logic
- [x] Update all tests for new signatures

### Phase 5: Documentation
- [x] Update README.md with new Go version requirement
- [x] Update README.md examples with context usage
- [x] Write CHANGELOG.md entry for v6.0.0
- [x] Update code documentation/comments

### Phase 6: Verification
- [x] Run `go build ./...`
- [x] Run `go test ./...`
- [x] Run `go vet ./...`
- [ ] Run `staticcheck ./...` (if available)
- [ ] Verify GitHub Actions pass

---

## 7. Migration Guide (for CHANGELOG.md)

```markdown
## 6.0.0 - [DATE]

### ⚠️ Breaking Changes

#### Minimum Go Version
- **BREAKING**: Now requires Go 1.26 or later (previously Go 1.11+)

#### Context Support
- **BREAKING**: All HTTP methods now require `context.Context` as the first parameter:
  - `Trigger(ctx, channel, event, data)`
  - `TriggerWithParams(ctx, channel, event, data, params)`
  - `TriggerMulti(ctx, channels, event, data)`
  - `TriggerMultiWithParams(ctx, channels, event, data, params)`
  - `TriggerBatch(ctx, batch)`
  - `SendToUser(ctx, userId, event, data)`
  - `Channels(ctx, params)`
  - `Channel(ctx, name, params)`
  - `GetChannelUsers(ctx, name)`

#### Removed Deprecated APIs
- **BREAKING**: Removed `TriggerExclusive()` - use `TriggerWithParams()` with `SocketID` parameter
- **BREAKING**: Removed `TriggerMultiExclusive()` - use `TriggerMultiWithParams()` with `SocketID` parameter
- **BREAKING**: Removed `AuthenticatePrivateChannel()` - use `AuthorizePrivateChannel()`
- **BREAKING**: Removed `AuthenticatePresenceChannel()` - use `AuthorizePresenceChannel()`
- **BREAKING**: Removed `EncryptionMasterKey` field - use `EncryptionMasterKeyBase64`

### Changed
- Replaced deprecated `io/ioutil` with `io` package
- Updated `golang.org/x/crypto` to v0.49.0
- Updated `testify` to v1.11.1 with new import path `github.com/stretchr/testify`
- Improved error handling (no more panics in library code)
- Fixed typo: "Paramater" → "Parameter" in error messages

### Fixed
- Added missing error check after `http.NewRequest` in request handling
```

---

## 8. Risks & Considerations

### Breaking Change Impact
- All existing users will need to update their code to add `context.Context` parameters
- Users relying on deprecated methods must migrate before upgrading

### MD5 Usage
- `crypto.go` uses MD5 for body hashing (`md5Signature`)
- This is likely a **Pusher API requirement**, not a security issue
- Verify with Pusher API docs before considering any changes

### Testify Migration
- Import path change from `gopkg.in/stretchr/testify.v1` to `github.com/stretchr/testify`
- All test files need updating

---

## 9. Testing Checklist

- [ ] All existing tests pass with new changes
- [ ] Context cancellation works correctly
- [ ] Context timeout works correctly  
- [ ] Error handling in crypto functions works
- [ ] Encrypted channels still work
- [ ] Webhook validation still works
- [ ] All example code in README compiles
- [ ] GitHub Actions CI passes on Go 1.25 and 1.26
