# 🧮 PackCalc — Automated Packing Calculator API

[![Go Version](https://img.shields.io/badge/Go-1.25-blue.svg)](https://golang.org)
![License](https://img.shields.io/badge/license-MIT-green.svg)
![AWS](https://img.shields.io/badge/AWS-EC2%20%7C%20RDS-orange.svg)
![Docker](https://img.shields.io/badge/Container-Docker-blue.svg)

---

## 🏗️ Overview

**PackCalc** is a backend API written in Go that calculates optimal pack distributions for product orders.  
It’s a production-ready, cloud-deployed system featuring **CI/CD automation**, **secure HTTPS**, **AWS infrastructure**, **custom ORM migrations**, and **multi-environment fallback** for reliability.

---

### 🌐 Frontend Demo

A public frontend to interact with the API is available here:  
👉 **NO MORE AVAILABLE**

This interface allows users to test requests visually, sending quantities and viewing pack results in real-time production.

---

### 🌐 Backend Production

👉 **NO MORE AVAILABLE**

This api results in real-time production development.

---

## 🚀 Features

- ⚙️ **Go 1.25 backend** with modular architecture
- 🧩 **Custom ORM & Migration Tool** — [`github.com/svvictorelias/go-migrate/pkg/migrate`](https://github.com/svvictorelias/go-migrate)
- 🐘 **AWS RDS (PostgreSQL)** with automatic **in-memory fallback**
- ☁️ **Deployed on AWS EC2** with Nginx reverse proxy
- 🔒 **Full HTTPS setup** using Let’s Encrypt (Certbot)
- 🔁 **GitHub Actions CI/CD** — build, test, migrate, and deploy automatically
- 🐳 **Dockerfile** ready for **Kubernetes** or **Terraform** pipelines
- 💓 Built-in health checks (`/health`)
- 🧠 Mock data layer for local/offline development

---

## 🧩 Architecture

```text
Frontend (Vercel)
       │
       ▼
HTTPS (Nginx reverse proxy)
       │
       ▼
  Go Backend API (PackCalc)
       ├── /health
       ├── /packs
       └── /calculate
       │
       ├── PostgreSQL (AWS RDS)
       └── In-memory fallback store
```

---

## ⚙️ Environment Configuration

Global variables on EC2 are defined in `/etc/environment`:

```bash
DATABASE_URL=postgres://postgres:postgres@localhost:5432/packcalc?sslmode=disable
PORT=8080
```

If the database becomes unavailable, the API logs:

```
DB not available: DATABASE_URL not set. Falling back to mock store (development).
```

and automatically switches to in-memory mode.

---

## 🧰 Local Development

### 🔧 Clone the repository

```bash
git clone git@github.com:svvictorelias/shipping-pack-backend.git
cd shipping-pack-backend
```

### ⚙️ Run locally

```bash
go mod tidy
docker compose up --build -d
make build
make migrations
make run
```

> The API starts on `http://localhost:8080` by default.

---

## 🧩 API Endpoints (cURL Examples)

### 0) Health Check

Checks if the service is up.

```bash
curl -i http://localhost:8080/health
```

**Expected:**

```json
{
  "status": "ok",
  "ts": "2025-10-26T10:55:45-03:00"
}
```

---

### 1) List Available Pack Sizes

Returns all pack sizes used by the optimizer.

```bash
curl -i http://localhost:8080/packs
```

**Example:**

```json
{ "packs": [250, 500, 1000, 2000, 5000] }
```

---

### 2) Update Available Pack Sizes (POST JSON)

```bash
curl -i -X POST http://localhost:8080/packs   -H "Content-Type: application/json"   -d '{"packs":[250, 500, 1000, 2000, 5000]}'
```

**Example response:**

```json
{
  "ok": true
}
```

### 3) Calculate Optimal Packs (POST JSON)

Computes the minimal oversupply first, then minimal number of packs.

```bash
curl -i -X POST http://localhost:8080/calculate   -H "Content-Type: application/json"   -d '{"items": 53}'
```

**Example response:**

```json
{
  "counts": {
    "54": 1
  },
  "pack_count": 1,
  "total_items": 54,
  "waste": 1
}
```

**Rules applied:**

1. Whole packs only (no splitting).
2. Minimize items sent above the requested amount.
3. Tie-breaker: minimize number of packs.

---


## 🧰 Makefile — Build, Run, Migrations, and Tests

### File: `Makefile`

### Usage

- **Build binary**

  ```bash
  make build
  ```

  Outputs `bin/packcalc`.

- **Run migrations**

  ```bash
  make migrations
  ```

  Executes migrations using the personal library `github.com/svvictorelias/go-migrate/pkg/migrate`.

- **Run API with local env**

  ```bash
  make run
  ```

  Loads variables from `.env.local` and starts the server.

- **Run tests with coverage**

  ```bash
  make test
  ```

  Produces `coverage.out` and prints a summary.

- **HTML coverage report**

  ```bash
  make cover-html
  ```

  Generates `coverage.html` for visual inspection.

- **Clean coverage artifacts**
  ```bash
  make test-clean
  ```

---

## 🔁 Continuous Deployment (GitHub Actions)

Each push to the `main` branch triggers a complete CI/CD pipeline:

1. ✅ Build & test the Go binary
2. 🧩 Run migrations against RDS
3. 🚀 Upload the binary to EC2
4. 🔄 Restart the app automatically via SSH

### 🔐 Secrets used:

- `EC2_SSH_KEY` → private SSH key for EC2 access
- `EC2_HOST` → EC2 public hostname or IP
- `DATABASE_URL` → PostgreSQL connection string (RDS)

---

## ☁️ AWS Infrastructure

### EC2 Instance

- OS: **Amazon Linux 2023**
- Reverse proxy: **Nginx**
- Managed HTTPS with **Certbot (Let’s Encrypt)**

### RDS Database

- Engine: PostgreSQL 15
- Connection: `sslmode=require`
- Security Group allows inbound connections from EC2

---

## 🔒 Nginx + HTTPS Configuration

### Enable HTTPS with Certbot

```bash
sudo dnf install certbot python3-certbot-nginx -y
sudo certbot --nginx -d shipping.fiianalise.com.br
sudo certbot renew --dry-run
```

---

## 🐳 Dockerfile (Kubernetes/Terraform Ready)

```dockerfile
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download && go build -o packcalc ./cmd/packcalc

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/packcalc .
EXPOSE 8080
CMD ["./packcalc"]
```

This container image can be used in:

- Kubernetes Deployments
- Terraform ECS modules
- Local Docker Compose setups

---

## 🧠 Extra Highlights

### 🧩 Custom ORM & Migration Engine

PackCalc uses a personal ORM and migration library:
[`github.com/svvictorelias/go-migrate/pkg/migrate`](https://github.com/svvictorelias/go-migrate)

Features:

- Transactional schema migrations
- Auto rollback on failure
- Lightweight ORM abstractions
- Compatible with any PostgreSQL database

---

### 🚀 Automated CI/CD Pipeline

- Builds, tests, migrates, and deploys on every push
- Uses GitHub Secrets for secure credential handling
- Zero-downtime redeployments
- Easily extendable for staging/production environments

---

### 🐳 Docker, K8s, and Terraform Integration

- Minimal container image (~12MB)
- Fully portable binary
- Supports K8s Deployments, Services, and Ingress
- Works with Terraform `aws_instance` or `aws_ecs_service` modules

---

### 🐘 AWS RDS with Fallback

- Production: PostgreSQL on AWS RDS
- Development: automatic in-memory fallback when the database is unreachable
- Ensures reliability and zero downtime for read-only routes

---

### ☁️ EC2 Hosting

- Hosted on **Amazon Linux 2023**
- Managed through **Nginx reverse proxy**
- Logs accessible via `/home/ec2-user/app.log`
- System environment variables stored in `/etc/environment`

---

### 🌐 DNS + HTTPS

| Component | Description                         |
| --------- | ----------------------------------- |
| DNS       | Managed via Registro.br             |
| Record    | `A` record → EC2 public IP          |
| HTTPS     | Enabled via Certbot (Let’s Encrypt) |
| Proxy     | Nginx (443 → 8080)                  |

---

## 💓 Health Check

Endpoint:

```bash
curl https://HOST/health
```

Response:

```json
{
  "status": "ok",
  "ts": "2025-10-26T10:55:45-03:00"
}
```

---

## 🧰 Development Mode

When `DATABASE_URL` is missing, PackCalc automatically switches to mock mode:

```
DB not available: DATABASE_URL not set. Falling back to mock store (development).
```

This mode uses an **in-memory data store**, perfect for local testing or CI pipelines.

Access:

```
http://localhost:8080/health
```

---

## 📄 License

MIT License © 2025 — Developed by **Victor Elias**

---

## ✨ Summary

PackCalc combines **Go performance**, **AWS scalability**, and **modern DevOps automation**.  
It’s production-ready, cloud-secured, and easy to extend across environments —  
from **local development** to **multi-cloud Kubernetes deployments**.
