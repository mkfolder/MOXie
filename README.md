<p align="center">
  <img src="https://img.shields.io/badge/status-beta-yellow?style=flat-square" alt="Status: Beta"/>
  <img src="https://img.shields.io/badge/network-Solana-9945FF?style=flat-square&logo=solana" alt="Network: Solana"/>
  <img src="https://img.shields.io/badge/go-%3E%3D1.25-blue?style=flat-square&logo=go" alt="Go"/>
  <img src="https://img.shields.io/badge/next.js-16-black?style=flat-square&logo=next.js" alt="Next.js 16"/>
  <img src="https://img.shields.io/badge/license-AGPLv3-green?style=flat-square" alt="License: GNU AGPLv3"/>
</p>

<p align="center">
  <h1 align="center">MOXie</h1>
  <p align="center"><em>Solana-native payment gateway for the modern merchant.</em></p>
</p>

---

MOXie is a **payment gateway** that lets merchants accept **SOL** (Solana's native token) directly into their own wallet — no middleman, no custodial risk. Payments are tracked on-chain using the **Solana Memo program**, eliminating the need for smart contracts.

> **Note:** This project is in early development. The **backend core is functional**, but the **UI is still a scaffold**.
> Currently **only Solana native token (SOL)** is supported.

---

## How It Works

### Registration (one-time)

```text
Merchant                  MOXie API                       Helius
   │                         │                              │
   │  POST /auth/register    │                              │
   │────────────────────────►│  CreateWebhook(addresses)    │
   │                         │─────────────────────────────►│
   │                         │                              │
```

### Payment flow

```text
Customer                   MOXie API                         Helius                    Merchant
   │                          │                                │                          │
   │  POST /orders/create     │                                │                          │
   │─────────────────────────►│                                │                          │
   │   { amount, metadata }   │                                │                          │
   │◄─────────────────────────│                                │                          │
   │    { address, memo }     │                                │                          │
   │                          │                                │                          │
   │   Send SOL + memo        │                                │                          │
   │   (on-chain)             │                                │                          │
   │──────────────────────────────────────────────────────────►│   (Helius detects tx)    │
   │                          │                                │                          │
   │                          │◄─── POST /helius-webhook ───── │                          │
   │                          │                                │                          │
   │                          │  Decode memo → find order      │                          │
   │                          │  Match account → extract SOL   │                          │
   │                          │  Accumulate payment            │                          │
   │                          │                                │                          │
   │                          │  POST /webhook (if fulfilled)  │                          │
   │                          │──────────────────────────────────────────────────────────►│
```

**The payment flow description:**

1. **Create an order** — The merchant's backend calls `POST /orders/create` with the requested amount and optional metadata. MOXie returns a destination Solana address, a unique base58-encoded memo, and a `qrcode_data` field containing a `solana:` URL for wallet scanning.

2. **Customer sends SOL + memo** — The customer sends SOL to the merchant's address with the memo attached (via the [Solana Memo Program](https://spl.solana.com/memo)). The memo acts as a payment reference, linking the on-chain transfer to the order.

3. **Helius detects the transfer** — A Helius webhook (registered during merchant onboarding) detects the incoming transaction and sends it to MOXie's backend.

4. **Match & accumulate** — MOXie decodes the memo to find the corresponding order, checks the native balance change on the merchant's account, and accumulates the paid amount. If the total meets or exceeds the requested amount, the order is marked as paid.

5. **Notify the merchant** — MOXie sends an HTTP POST to the merchant's registered webhook URL with the order details and transaction data.

### QR code scanning (Solana Pay)

The `POST /orders/create` response includes a `qrcode_data` field with a `solana:` URL. Merchants can render this as a QR code. When scanned with a Solana wallet (e.g. Phantom, Backpack), the wallet automatically populates the recipient, amount, and memo — the customer just taps confirm.

The Solana Pay endpoints (`GET /solpay/:id`, `POST /solpay/:id`) serve the metadata and transaction building required by the [Solana Pay spec](https://docs.solanapay.com).

---

## Project Structure

```text
moxie/
├── api/                         # Go backend
│   ├── cmd/
│   │   ├── server/main.go       # HTTP server entrypoint
│   │   └── migration/migrate.go # DB auto-migration
│   ├── config/default.yaml      # App configuration
│   └── internal/
│       ├── handler/             # HTTP handlers (auth, orders, webhook)
│       ├── service/             # Business logic (payment processing)
│       ├── constants/           # Shared constants (domain, subdomain)
│       ├── helius/              # Helius RPC client & types
│       ├── models/              # GORM models (Merchant, Order)
│       ├── routes/              # Route registration
│       ├── workers/             # Background workers (order cleanup)
│       └── db/                  # Generic repository layer
│
├── ui/                          # Next.js frontend (WIP)
│   ├── pages/                   # Page routes
│   ├── components/              # Reusable UI components
│   └── layouts/                 # Page layouts
│
├── caddy/                        # Caddy reverse proxy
│   ├── Caddyfile                 # Routes: /api/* → API, everything else → UI
│   └── Dockerfile
│
├── docker-compose.yaml          # API + UI + PostgreSQL + pgAdmin + Caddy
└── Makefile                     # Docker manipulation
```

---

## API Endpoints

| Method | Path                     | Description                               |
|--------|--------------------------|-------------------------------------------|
| `POST` | `/auth/register`         | Register a new merchant                   |
| `POST` | `/auth/login`            | Authenticate a merchant                   |
| `GET`  | `/health`                | Health check                              |
| `GET`  | `/orders/find-all`       | List all orders                           |
| `GET`  | `/orders/find/:id`       | Get order by ID                           |
| `POST` | `/orders/create`         | Create a new order                        |
| `GET`  | `/solpay/:id`            | Solana Pay metadata for wallet to display |
| `POST` | `/solpay/:id`            | Build transaction for Solana Pay flow     |
| `POST` | `/helius-webhook/handle` | Incoming Helius webhook                   |

---

## Tech Stack

| Layer           | Technology                                                     |
|-----------------|----------------------------------------------------------------|
| **Backend**     | Go 1.25, Fiber v3, GORM, PostgreSQL 17, Zap logging            |
| **Frontend**    | Next.js 16 (Pages Router), React 19, HeroUI v3, Tailwind v4    |
| **Blockchain**  | Solana, Solana Pay, Memo Program, Helius RPC + Webhooks        |
| **Infra**       | Docker, Docker Compose, pgAdmin                                |

---

## Getting Started

### Prerequisites

- [Docker](https://docs.docker.com/get-docker/) & [Docker Compose](https://docs.docker.com/compose/)
- A [Helius](https://helius.xyz) API key (for Solana transaction monitoring)

### Quick Start

```bash
make compose
```

This spins up five services:

| Service      | Port        | Description                          |
|--------------|-------------|--------------------------------------|
| **caddy**    | 80          | Reverse proxy (single entry point)   |
| **api**      | 7654        | Go/Fiber backend (internal)          |
| **ui**       | 4321        | Next.js frontend (internal)          |
| **postgres** | 5432        | PostgreSQL database                  |
| **pgadmin**  | 1444        | Database admin panel                 |

Caddy acts as the single entry point:

- [`http://localhost`](http://localhost) → UI
- [`http://localhost/api`](http://localhost/api) → API

### Configuration

Copy the example environment file and edit as needed:

```bash
cp api/.env.example api/.env
```

Key configuration values are in `api/config/default.yaml`:

```yaml
server:
  port: 7654

postgres:
  dsn: "host=postgres user=moxie password=moxie dbname=moxie port=5432 sslmode=disable"

helius:
  api_url: "https://api-devnet.helius-rpc.com"

workers:
  cleaner:
    interval: 30m
    order_expiration: 24h
```

---

## What's Implemented vs. What's Missing

### ✅ Backend (Functional)

- Merchant registration & authentication (bcrypt + base58 API keys)
- Order CRUD with unique memo generation (base58-encoded UUID)
- Helius webhook integration for real-time payment monitoring
- On-chain memo decoding and order matching
- Payment accumulation with partial payment support
- Merchant notification via outbound webhook
- Solana Pay QR code generation (`solana:` URL scheme)
- Solana Pay metadata & transaction endpoints for wallet scanning
- Background worker for expired order cleanup

### 🚧 In Progress

- **JWT token issuing** — currently using base58 UUID as session token
- **Transaction storage** — payments are accumulated on the order row; a dedicated transaction history table is planned
- **Frontend dashboard** — the UI is currently the HeroUI starter template, awaiting MOXie-specific branding, merchant portal, and analytics

### 🔜 Planned

- Transaction history and detailed analytics
- SPL token support
- Multi-chain support (Base, Sui)
- SDK/client libraries for popular languages
- Dashboard with real-time payment monitoring

---

## Architecture Decisions

**Why Memo Program instead of smart contracts?**  
The Solana Memo program allows attaching arbitrary data to SOL transfers. This means merchants can receive payments directly into their own wallet (no custodial risk) while MOXie handles the order matching off-chain. No smart contract audits, no additional attack surface.

**Why Helius?**  
Helius provides enhanced transaction webhooks that include account-level balance changes, making it straightforward to detect incoming SOL payments and associate them with orders without polling the chain.

---

## Contributing

This is a personal project — I'm not accepting code contributions. However, bug reports, issues, and suggestions are very welcome.

---

<p align="center">
  <strong>MOXie</strong> — payments, but faster.
</p>
