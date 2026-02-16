# 🚀 Quick Start Guide - Building the EggCarton CLI

## ✅ Setup Complete!

You now have:
- ✅ Go module initialized
- ✅ Cobra CLI framework installed
- ✅ Browser package for opening system browser
- ✅ Complete folder structure
- ✅ Template files with TODO comments

## 📋 Your Implementation Order

### **Phase 1: Configuration (Start Here!)**

**File:** `cli/config/config.go`

**What to implement:**
1. `SaveTokens()` - Save tokens to `~/.eggcarton/credentials.json` with 0600 permissions
2. `LoadTokens()` - Read tokens from file
3. `IsTokenValid()` - Check if access token is expired

**Test it:**
```go
// Create a test file to verify your config works
tokens := &config.TokenData{
    AccessToken: "test",
    ExpiresIn: 3600,
    IssuedAt: time.Now().Unix(),
}
cfg, _ := config.LoadConfig()
cfg.SaveTokens(tokens)
loaded, _ := cfg.LoadTokens()
fmt.Println(loaded.IsTokenValid())
```

---

### **Phase 2: PKCE & OAuth Flow**

#### **Step 2.1:** `cli/auth/login.go`

**What to implement:**
1. `GeneratePKCEChallenge()` - Create code verifier and challenge

**Key algorithm:**
```
verifier = base64url(random_32_bytes)
challenge = base64url(sha256(verifier))
```

**Test it:**
```go
pkce, _ := auth.GeneratePKCEChallenge()
fmt.Println("Verifier:", pkce.Verifier)
fmt.Println("Challenge:", pkce.Challenge)
// Verifier should be 43 chars, Challenge should be 43 chars
```

2. `BuildAuthorizationURL()` - Construct the OAuth URL

**Expected output:**
```
https://eggcarton-auth-uqhqvdut.auth.us-west-1.amazoncognito.com/oauth2/authorize?client_id=1vccvf2hh5amna78lurbn9bjhi&response_type=code&scope=openid+email+profile&redirect_uri=http://localhost:8080/callback&code_challenge=<YOUR_CHALLENGE>&code_challenge_method=S256
```

#### **Step 2.2:** `cli/auth/server.go`

**What to implement:**
1. `StartCallbackServer()` - Local HTTP server to catch OAuth redirect

**Key points:**
- Listen on `:8080`
- Extract `code` from query params
- Return nice HTML page
- Shutdown after receiving code
- Handle timeout (5 minutes)

**Test it:**
```bash
# In terminal 1:
go run main.go login  # Will start server

# In terminal 2:
curl "http://localhost:8080/callback?code=TEST123"
# Should see success page and server should shutdown
```

#### **Step 2.3:** `cli/auth/token.go`

**What to implement:**
1. `ExchangeCodeForTokens()` - Trade authorization code for JWT
2. `RefreshAccessToken()` - Get new access token when expired

**Request format:**
```
POST https://eggcarton-auth-uqhqvdut.auth.us-west-1.amazoncognito.com/oauth2/token
Content-Type: application/x-www-form-urlencoded

grant_type=authorization_code&client_id=1vccvf2hh5amna78lurbn9bjhi&code=<CODE>&redirect_uri=http://localhost:8080/callback&code_verifier=<VERIFIER>
```

#### **Step 2.4:** Wire up `cli/commands/login.go`

**Flow:**
```go
1. cfg := config.LoadConfig()
2. pkce := auth.GeneratePKCEChallenge()
3. authURL := auth.BuildAuthorizationURL(...)
4. browser.OpenURL(authURL)  // Opens browser!
5. code := auth.StartCallbackServer()  // Waits for callback
6. tokens := auth.ExchangeCodeForTokens(code, pkce.Verifier)
7. cfg.SaveTokens(tokens)
8. fmt.Println("✅ Login successful!")
```

**Test the full login flow:**
```bash
cd cli
go run main.go login
# Browser should open, you authenticate, then CLI saves tokens
```

---

### **Phase 3: API Client**

**File:** `cli/api/client.go`

**What to implement:**
1. `ExtractOwnerFromToken()` - Decode JWT to get user's `sub` claim
2. `PutEgg()` - Store a secret
3. `GetEgg()` - Retrieve a secret
4. `BreakEgg()` - Delete a secret

**JWT Decoding Hint:**
```go
parts := strings.Split(accessToken, ".")
// parts[1] is the payload
payloadBytes, _ := base64.RawURLEncoding.DecodeString(parts[1])
var claims map[string]interface{}
json.Unmarshal(payloadBytes, &claims)
sub := claims["sub"].(string)  // This is the owner/user ID
```

**Test API client:**
```bash
# First login
go run main.go login

# Then test each operation
go run main.go add TEST_KEY test_value
go run main.go get TEST_KEY
go run main.go break TEST_KEY
```

---

### **Phase 4: Wire up Commands**

**Files:** `cli/commands/add.go`, `get.go`, `break.go`

**Pattern for each command:**
```go
func runAdd(cmd *cobra.Command, args []string) error {
    // 1. Load config
    cfg, _ := config.LoadConfig()
    
    // 2. Load tokens
    tokens, _ := cfg.LoadTokens()
    if tokens == nil {
        return fmt.Errorf("not logged in, run: egg login")
    }
    
    // 3. Check/refresh if needed
    if !tokens.IsTokenValid() {
        tokens, _ = auth.RefreshAccessToken(...)
        cfg.SaveTokens(tokens)
    }
    
    // 4. Extract owner from JWT
    owner, _ := api.ExtractOwnerFromToken(tokens.AccessToken)
    
    // 5. Create API client
    client := api.NewClient(cfg.APIEndpoint, tokens.AccessToken)
    
    // 6. Call appropriate method
    err := client.PutEgg(owner, key, value)
    
    // 7. Print result
    fmt.Printf("✅ Secret '%s' stored successfully!\n", key)
    return err
}
```

---

### **Phase 5: The Killer Feature - `egg run`**

**File:** `cli/commands/run.go`

**What to implement:**

**Step 5.1:** Add a Lambda to list all eggs (or enhance GetEgg)
```go
// You might need to update your Lambda to support:
// GET /eggs?owner={owner}
// Returns: [{"secret_id": "KEY1", "plaintext": "val1"}, ...]
```

**Step 5.2:** Implement the run command
```go
func runRun(cmd *cobra.Command, args []string) error {
    // 1. Find the "--" separator
    dashIndex := -1
    for i, arg := range args {
        if arg == "--" {
            dashIndex = i
            break
        }
    }
    if dashIndex == -1 {
        return fmt.Errorf("usage: egg run -- [command]")
    }
    
    // 2. Extract command after "--"
    command := args[dashIndex+1:]
    
    // 3. Fetch all secrets
    owner, _ := api.ExtractOwnerFromToken(tokens.AccessToken)
    client := api.NewClient(cfg.APIEndpoint, tokens.AccessToken)
    secrets, _ := client.ListEggs(owner)  // map[string]string
    
    // 4. Build environment
    env := os.Environ()  // Copy existing env vars
    for key, value := range secrets {
        env = append(env, fmt.Sprintf("%s=%s", key, value))
    }
    
    // 5. Execute command
    cmd := exec.Command(command[0], command[1:]...)
    cmd.Env = env
    cmd.Stdin = os.Stdin
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    
    err := cmd.Run()
    
    // 6. Exit with same code
    if exitErr, ok := err.(*exec.ExitError); ok {
        os.Exit(exitErr.ExitCode())
    }
    
    return err
}
```

**Test it:**
```bash
# Store some secrets
egg add DATABASE_URL postgres://localhost/mydb
egg add API_KEY secret123

# Test injection
egg run -- env | grep DATABASE_URL
# Should output: DATABASE_URL=postgres://localhost/mydb

# Run actual program with secrets
egg run -- go run main.go
egg run -- npm start
```

---

## 🧪 Full Testing Flow

```bash
# 1. Build the CLI
cd cli
go build -o egg main.go

# 2. Login
./egg login
# ✅ Browser opens, authenticate, tokens saved

# 3. Add secrets
./egg add DB_HOST localhost
./egg add DB_PORT 5432
./egg add API_SECRET supersecret

# 4. Get a secret
./egg get DB_HOST
# Output: localhost

# 5. List with run (create a test script)
echo '#!/bin/bash\necho "DB_HOST=$DB_HOST"' > test.sh
chmod +x test.sh
./egg run -- ./test.sh
# Output: DB_HOST=localhost

# 6. Delete a secret
./egg break DB_PORT

# 7. Install globally
sudo mv egg /usr/local/bin/
egg --help
```

## 🐛 Debugging Tips

1. **Login not working?**
   - Check Cognito domain and client ID in config.go
   - Make sure redirect URI matches (http://localhost:8080/callback)
   - Check browser console for errors

2. **API calls failing?**
   - Print the access token to verify format
   - Check API endpoint URL
   - Look at Lambda logs in CloudWatch
   - Verify Lambda has proper CORS headers

3. **Token expired?**
   - Implement token refresh in commands
   - Check `IsTokenValid()` logic

4. **Secrets not injecting?**
   - Print environment variables to debug
   - Make sure ListEggs returns correct format
   - Check Lambda query endpoint

## 📚 Resources

- [Cobra Docs](https://github.com/spf13/cobra)
- [AWS Cognito OAuth Endpoints](https://docs.aws.amazon.com/cognito/latest/developerguide/token-endpoint.html)
- [PKCE Spec](https://oauth.net/2/pkce/)
- [Go crypto/rand](https://pkg.go.dev/crypto/rand)
- [Go os/exec](https://pkg.go.dev/os/exec)

## 🎯 Success Criteria

- [ ] `egg login` authenticates via browser
- [ ] Tokens saved securely with 0600 permissions
- [ ] `egg add` stores encrypted secrets
- [ ] `egg get` retrieves decrypted secrets
- [ ] `egg break` deletes secrets
- [ ] `egg run` injects secrets into subprocess
- [ ] Token refresh works automatically
- [ ] Good error messages for users

---

## 🚀 Start Coding!

**Begin with Phase 1 - config.go**

Open `cli/config/config.go` and implement the `SaveTokens()` function. Once that works, move to `LoadTokens()`, then `IsTokenValid()`.

Ask me questions as you go! I'm here to help guide you through any tricky parts. 🥚

Happy coding! 🎉
