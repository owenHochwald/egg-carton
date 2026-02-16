# 🥚 EggCarton CLI

A secure CLI tool for managing secrets using AWS Lambda and Cognito authentication.

## Features

- **EC-12**: `egg login` - OAuth authentication via Cognito Hosted UI with PKCE flow
- **EC-13**: CRUD operations (`add`, `get`, `break`) - Manage secrets via authenticated API calls
- **EC-14**: `egg run` - Inject secrets as environment variables and execute commands

## Project Structure

```
cli/
├── main.go                 # CLI entry point with cobra commands
├── config/
│   └── config.go          # Configuration management (API endpoint, credentials, etc.)
├── auth/
│   ├── login.go           # OAuth flow implementation (PKCE)
│   ├── server.go          # Local callback server for OAuth redirect
│   └── token.go           # Token storage and refresh logic
├── api/
│   └── client.go          # HTTP client for Lambda API calls (with JWT auth)
├── commands/
│   ├── login.go           # `egg login` command
│   ├── add.go             # `egg add [key] [value]` command
│   ├── get.go             # `egg get [key]` command
│   ├── break.go           # `egg break [key]` command
│   └── run.go             # `egg run -- [command]` command
└── README.md              # This file
```

## Development Roadmap

### Phase 1: Setup & Configuration (EC-12 Foundation)
**Files to create:** `config/config.go`

**Your task:**
1. Create a config struct to store:
   - API endpoint from terraform output
   - Cognito configuration (user pool ID, client ID, domain)
   - Token storage location (`~/.eggcarton/credentials.json`)
2. Implement config loading from environment variables or config file
3. Add helper functions to save/load tokens

**Key considerations:**
- Store tokens securely (file permissions 0600)
- Support multiple profiles (optional)
- Validate config on load

**Example structure:**
```go
type Config struct {
    APIEndpoint    string
    CognitoConfig  CognitoConfig
    TokenPath      string
}

type CognitoConfig struct {
    UserPoolID string
    ClientID   string
    Domain     string
    Region     string
}
```

### Phase 2: OAuth Login Flow (EC-12)
**Files to create:** `auth/login.go`, `auth/server.go`, `auth/token.go`

**Your task:**

#### Step 2a: PKCE Generation (`auth/login.go`)
1. Generate code verifier (43-128 random characters)
2. Generate code challenge (SHA256 hash of verifier, base64url encoded)
3. Build authorization URL with PKCE parameters

**OAuth URL format:**
```
https://{cognito-domain}/oauth2/authorize?
  client_id={client-id}
  &response_type=code
  &scope=openid+email+profile
  &redirect_uri=http://localhost:8080/callback
  &code_challenge={challenge}
  &code_challenge_method=S256
```

#### Step 2b: Local Callback Server (`auth/server.go`)
1. Start HTTP server on `localhost:8080`
2. Handle `/callback` route to capture authorization code
3. Exchange code for JWT tokens using PKCE verifier
4. Gracefully shutdown server after receiving tokens

**Token exchange endpoint:**
```
POST https://{cognito-domain}/oauth2/token
Content-Type: application/x-www-form-urlencoded

grant_type=authorization_code
&client_id={client-id}
&code={authorization-code}
&redirect_uri=http://localhost:8080/callback
&code_verifier={original-verifier}
```

#### Step 2c: Token Management (`auth/token.go`)
1. Parse JWT response (access_token, id_token, refresh_token, expires_in)
2. Save tokens to `~/.eggcarton/credentials.json`
3. Implement token refresh when expired
4. Add helper to get current valid access token

**Hints:**
- Use `crypto/rand` for secure random generation
- Use `crypto/sha256` for code challenge
- Use `encoding/base64` with RawURLEncoding
- Use `net/http` for both opening browser and local server
- Consider using `github.com/pkg/browser` to open system browser

### Phase 3: API Client (EC-13 Foundation)
**Files to create:** `api/client.go`

**Your task:**
1. Create HTTP client that adds JWT bearer token to requests
2. Implement methods for each Lambda endpoint:
   - `PutEgg(key, value string) error`
   - `GetEgg(key string) (string, error)`
   - `BreakEgg(key string) error`
3. Handle API errors and token expiration
4. Parse Lambda responses

**API Endpoints (based on your Lambda functions):**
```
PUT    {api-endpoint}/eggs
GET    {api-endpoint}/eggs/{owner}
DELETE {api-endpoint}/eggs/{owner}
```

**Request format example (PUT):**
```json
{
  "owner": "user-sub-from-jwt",
  "secret_id": "key",
  "plaintext": "value"
}
```

**Key considerations:**
- Extract `sub` claim from JWT for owner field
- Set `Authorization: Bearer {access-token}` header
- Handle 401 responses by refreshing token
- Return meaningful error messages

### Phase 4: CLI Commands (EC-13)
**Files to create:** `commands/login.go`, `commands/add.go`, `commands/get.go`, `commands/break.go`

**Your task:**
1. Use `github.com/spf13/cobra` to define commands
2. Wire up each command to call the appropriate API client method

**Command implementations:**

#### `egg login`
```go
// Pseudocode:
- Check if already logged in (valid token exists)
- If yes, ask to re-login or skip
- Initiate OAuth flow (call auth.StartLogin())
- Print success message with user info
```

#### `egg add [key] [value]`
```go
// Pseudocode:
- Validate arguments
- Get valid access token (refresh if needed)
- Call api.PutEgg(key, value)
- Print success message
```

#### `egg get [key]`
```go
// Pseudocode:
- Validate arguments
- Get valid access token
- Call api.GetEgg(key)
- Print decrypted value to stdout
```

#### `egg break [key]`
```go
// Pseudocode:
- Validate arguments
- Get valid access token
- Call api.BreakEgg(key)
- Print confirmation message
```

### Phase 5: Secret Injection (EC-14)
**Files to create:** `commands/run.go`

**Your task:**
1. Fetch all user secrets from API
2. Parse secrets into key-value pairs
3. Set as environment variables
4. Execute user command with `os/exec`

**Implementation steps:**

#### Step 5a: Fetch All Secrets
- Extend API client to support listing all secrets
- You may need to update Lambda to support scanning user's eggs

#### Step 5b: Environment Injection
```go
// Pseudocode:
- Get all secrets for authenticated user
- Build environment variable map
- Merge with current os.Environ()
- Create exec.Command with custom env
- Wire up stdin/stdout/stderr
- Run command and wait for completion
```

**Command format:**
```bash
egg run -- go run main.go
egg run -- npm start
egg run -- ./my-script.sh
```

**Key considerations:**
- Use `os/exec.Command` to run subprocess
- Pass through stdin/stdout/stderr for interactivity
- Preserve existing environment variables
- Handle command parsing (split on `--`)
- Exit with same code as subprocess

### Phase 6: Main Entry Point
**Files to create:** `cli/main.go`

**Your task:**
1. Initialize cobra root command
2. Add all subcommands
3. Set up global flags (--profile, --verbose, etc.)
4. Handle graceful errors

**Example structure:**
```go
package main

import (
    "github.com/spf13/cobra"
    "github.com/owenHochwald/egg-carton/cli/commands"
)

func main() {
    rootCmd := &cobra.Command{
        Use:   "egg",
        Short: "EggCarton - Secure secret management CLI",
    }

    rootCmd.AddCommand(commands.LoginCmd)
    rootCmd.AddCommand(commands.AddCmd)
    rootCmd.AddCommand(commands.GetCmd)
    rootCmd.AddCommand(commands.BreakCmd)
    rootCmd.AddCommand(commands.RunCmd)

    if err := rootCmd.Execute(); err != nil {
        os.Exit(1)
    }
}
```

## Dependencies

Add these to your `go.mod`:

```bash
go get github.com/spf13/cobra@latest
go get github.com/pkg/browser@latest
go get golang.org/x/oauth2@latest  # Optional, for OAuth helpers
```

## Build & Install

```bash
cd cli
go build -o egg main.go
sudo mv egg /usr/local/bin/
```

## Testing Checklist

- [ ] `egg login` opens browser and saves tokens
- [ ] `egg add TEST_KEY test_value` stores secret
- [ ] `egg get TEST_KEY` retrieves and decrypts secret
- [ ] `egg break TEST_KEY` deletes secret
- [ ] `egg run -- echo $TEST_KEY` injects secret into command
- [ ] Token refresh works when access token expires
- [ ] Error messages are helpful and clear

## Security Notes

- ✅ Uses PKCE flow (no client secret needed)
- ✅ Tokens stored with 0600 permissions
- ✅ Access tokens short-lived (1 hour)
- ✅ Refresh tokens allow re-authentication without browser
- ✅ All API calls use HTTPS with bearer token
- ✅ Secrets encrypted at rest in DynamoDB with KMS

## Next Steps

1. Start with Phase 1 (config management)
2. Move to Phase 2 (OAuth login)
3. Test login flow before moving to API client
4. Implement CRUD operations
5. Finally, build the `run` command

Ask me questions as you build each component! I'll help guide you through any tricky parts.
