# 🥚 Egg Carton

> **The problem:** AI agents need API keys. Storing them in `.env` files or shell history is a security nightmare.  
> **The solution:** Zero-trust secret injection - your keys stay encrypted in the cloud, injected only at runtime.

**Egg Carton** is a serverless secret manager built for developers who work with AI agents. It uses AWS KMS envelope encryption to protect your secrets, OAuth for authentication, and a dead-simple CLI for secret injection.

```bash
# Store secrets once
egg lay ANTHROPIC_API_KEY "sk-ant-..."
egg lay OPENAI_API_KEY "sk-..."

# Use them anywhere, without ever exposing plaintext
egg hatch -- python ai_agent.py
egg hatch -- node autonomous_workflow.js

# Secrets appear as env vars, never touch disk
```

### Why this matters for AI development
- **AI agents scan everything** - .env files, shell history, logs. Egg Carton keeps secrets out of reach.
- **Seamless injection** - No manual `export` commands. Just prefix with `egg hatch --`.
- **Zero-trust by default** - Secrets encrypted with unique keys, scoped to your identity.
- **OAuth in 1 click** - No AWS credentials to manage. Just `egg login` and you're done.

<div align="center">
  <img src="public/demo.gif" alt="Egg Carton Demo" width="800"/>
</div>

---


### Core Commands
| Command | What it does |
|---------|-------------|
| `egg login` | Authenticate via OAuth (opens browser) |
| `egg lay KEY "value"` | Store a secret (or use alias: `egg add`) |
| `egg get KEY` | Retrieve a secret |
| `egg get` | List all your secret keys |
| `egg hatch -- <cmd>` | Run command with secrets injected (or use alias: `egg run`) |
| `egg break KEY` | Delete a secret |

---

## 🆚 Why Not Just Use AWS Secrets Manager?

Secrets Manager works, but it's **built for enterprises, not individual use**:

| | AWS Secrets Manager | Egg Carton |
|---|---|---|
| **Setup** | IAM policies, resource policies, cross-account access | `egg login` |
| **CLI** | `aws secretsmanager get-secret-value --secret-id MySecret --query SecretString --output text` | `egg get MySecret` |
| **Secret Injection** | Write bash scripts: `export KEY=$(aws secretsmanager...)` | `egg hatch -- command` |
| **Auth** | AWS IAM credentials | OAuth (browser login) |

---

## 🔐 How It Works

**Architecture:** Serverless backend (Lambda + DynamoDB + KMS) + CLI with OAuth authentication

```
┌─────────────┐
│  egg CLI    │  ← You interact here
└──────┬──────┘
       │ HTTPS (OAuth token in header)
       ▼
┌─────────────┐
│ API Gateway │  ← Validates JWT
└──────┬──────┘
       │
       ▼
┌─────────────┐       ┌──────────┐       ┌─────────┐
│   Lambda    │ ◄───► │ DynamoDB │       │   KMS   │
│  Functions  │       │ (secrets)│ ◄───► │(crypto) │
└─────────────┘       └──────────┘       └─────────┘
```

### Security: Envelope Encryption

Each secret gets its own unique encryption key:

1. **Generate** a 256-bit Data Encryption Key (DEK) via KMS
2. **Encrypt** your secret with the DEK using AES-256-GCM
3. **Encrypt** the DEK itself with KMS master key
4. **Store** both encrypted secret + encrypted DEK in DynamoDB
5. **Decrypt** on retrieval (KMS unwraps DEK → DEK decrypts secret)


### Authentication: OAuth PKCE Flow

No AWS credentials needed. Just browser-based OAuth:

1. `egg login` → Opens browser to Cognito
2. You authenticate (Google, etc.)
3. Cognito redirects to `localhost:8080/callback` with auth code
4. CLI exchanges code for JWT tokens (access + refresh)
5. Tokens stored in `~/.eggcarton/credentials.json` (0600 permissions)
6. Access token (1hr expiry) used for API calls, auto-refreshed when needed

---

## 🛠️ For Developers

<details>
<summary><b>Infrastructure Setup (Terraform)</b></summary>

```bash
# Prerequisites: AWS CLI configured, Terraform installed
cd egg-carton

# Configure
cp terraform.tfvars.example terraform.tfvars
# Edit terraform.tfvars with your Google OAuth credentials

# Deploy
terraform init
terraform plan -out=tfplan
terraform apply tfplan

# Get outputs (API endpoint, Cognito config)
terraform output
```

**Stack:**
- Lambda (provided.al2023 custom runtime)
- DynamoDB (PAY_PER_REQUEST)
- KMS (with automatic rotation)
- API Gateway HTTP API
- Cognito User Pool with Google IdP
- CloudWatch Dashboard

</details>

<details>
<summary><b>CLI Development</b></summary>

```bash
cd cli

# Development
go run main.go login
go run main.go lay TEST_KEY "test value"

# Build
make build

# Install system-wide
make install  # installs to /usr/local/bin/egg

# Test
make test
```

**Tech Stack:**
- Go 1.21+
- Cobra (CLI framework)
- OAuth 2.0 PKCE implementation
- JWT parsing for token management

</details>

<details>
<summary><b>Project Structure</b></summary>

```
egg-carton/
├── cli/                       # CLI tool
│   ├── commands/              # login, lay, get, break, hatch
│   ├── auth/                  # OAuth PKCE + token refresh
│   ├── api/                   # HTTP client for Lambda API
│   └── config/                # Config + token storage
├── cmd/lambda/                # Lambda functions
│   ├── put_egg/               # Store secret
│   ├── get_egg/               # Retrieve secrets
│   └── break_egg/             # Delete secret
├── pkg/crypto/                # AES-256-GCM encryption
├── main.tf                    # Infrastructure
├── cognito.tf                 # OAuth setup
└── cloudwatch_dashboard.tf    # Monitoring
```

</details>

---

## Monitoring

CloudWatch dashboard includes:
- Lambda invocations, errors, duration
- DynamoDB read/write capacity, throttles
- API Gateway requests, 4xx/5xx errors
- KMS API call counts
- Cognito authentication metrics

Access via AWS Console or `terraform output dashboard_url`

---

## 🤝 Contributing

Contributions welcome! Please feel free to submit a Pull Request.

**Development priorities:**
- [ ] Add secret rotation policies
- [ ] Support for secret versioning
- [ ] CLI autocomplete (bash/zsh)
- [ ] Export secrets to .env format
- [ ] Multi-region replication

---

## 📝 License

MIT License - see [LICENSE](LICENSE) file
