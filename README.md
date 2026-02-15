# 🥚 Egg Carton

**"Crack the code, not the security."**

Egg Carton is a serverless secret management proxy built for the era of autonomous AI agents. It eliminates the need for plain-text `.env` files by providing a hardened, encrypted "shell" for API keys, injecting them into agentic workflows (like Claude Code) only at the moment of execution.

Modern AI agents require high-level access to various APIs. Storing these keys locally or in shell history is a massive security risk. Egg Carton solves this by acting as a **Zero-Trust Middleman** between your local machine and the cloud.

## Security Model

### Envelope Encryption Flow

```
┌──────────────────────────────────────────────────────────────┐
│                     ENVELOPE NCRYPTION                       │
├──────────────────────────────────────────────────────────────┤
│                                                              │
│  1. Master Key (KMS)                                         │
│     ┌────────────────────────┐                               │
│     │  KMS Master Key         │  ◀── Never leaves KMS        │
│     │  (256-bit)             │      AWS managed              │
│      └──────────┬─────────────┘                              │
│                │                                             │
│                │ generates                                   │
│                ▼                                             │
│  2. Data Encryption Key (DEK)                               │
│     ┌────────────────────────┐                              │
│     │  Plaintext DEK         │  ◀── Used once, then deleted │
│     │  (256-bit AES key)     │                              │
│     └──────────┬─────────────┘                              │
│                │                                             │
│                │ encrypts                                    │
│                ▼                                             │
│  3. Your Secret                                             │
│     ┌────────────────────────┐                              │
│     │  "password123"         │  ◀── Your actual data        │
│     └──────────┬─────────────┘                              │
│                │                                             │
│                │ becomes                                     │
│                ▼                                             │
│  4. Ciphertext                                              │
│     ┌────────────────────────┐                              │
│     │  [encrypted bytes]     │  ◀── Stored in DynamoDB      │
│     └────────────────────────┘                              │
│                                                              │
│  5. Encrypted DEK                                           │
│     ┌────────────────────────┐                              │
│     │  [KMS-wrapped DEK]     │  ◀── Also stored in DynamoDB │
│     └────────────────────────┘                              │
│                                                              │
└──────────────────────────────────────────────────────────────┘
```

### Why This Is Secure

1. **Master Key Protection**
   - Never leaves AWS KMS hardware security modules (HSMs)
   - FIPS 140-2 validated
   - Automatic key rotation available

2. **Unique DEK Per Secret**
   - Each secret has its own encryption key
   - If one DEK is compromised, only one secret is affected
   - No single point of failure

3. **Encrypted at Rest**
   - DynamoDB table encrypted with KMS
   - Ciphertext is encrypted
   - Encrypted DEK is double-encrypted (once by master key, once by table encryption)

4. **Encrypted in Transit**
   - HTTPS/TLS for all API calls
   - TLS for all AWS service calls

5. **No Plaintext Storage**
   - Plaintext secret never touches disk
   - Plaintext DEK exists only in Lambda memory briefly
   - Lambda containers are ephemeral

6. **Audit Trail**
   - All KMS operations logged to CloudTrail
   - All API calls logged to CloudWatch
   - DynamoDB streams available for change tracking
