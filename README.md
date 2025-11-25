# 🚀 Multi-Service Monorepo with Full CI/CD (Node.js, Python, Go)

This repository is a complete **multi-service monorepo** containing three microservices written in **Node.js**, **Python (FastAPI)**, and **Go**, along with a fully automated **CI/CD pipeline in Jenkins**.

The system includes:
✔ Change detection per service
✔ Automated linting
✔ Automated testing
✔ Security scanning
✔ Docker image builds
✔ Isolation per service
✔ Standardized, extensible workflow for real-world DevOps setups

---

## 📁 Project Structure

```
.
├── user-service/                 # Node.js + Express microservice
│   ├── src/
│   ├── tests/
│   ├── Dockerfile
│   ├── package.json
│   ├── package-lock.json
│   └── eslint.config.js
│
├── transaction-service/          # Python + FastAPI microservice
│   ├── main.py
│   ├── test_main.py
│   ├── requirements.txt
│   └── Dockerfile
│
├── notification-service/         # Go (net/http) microservice
│   ├── main.go
│   ├── main_test.go
│   ├── go.mod
│   └── Dockerfile
│
├── detect_changes.sh             # Detects which services changed between commits
├── lint.sh                       # Unified linting script for all languages
├── test.sh                       # Unified testing script for all services
├── scan.sh                       # Security scanning for all services
├── trufflehog_exclude.txt        # Exclusions for secret scanning
└── Jenkinsfile                   # Complete Jenkins CI/CD pipeline
```

---

# 🧩 Microservices Overview

## 1️⃣ user-service (Node.js + Express)

A simple REST service that performs user validation.

Features:

* Request validation
* Express routing
* Unit tests using Jest + Supertest
* ESLint configuration
* Dockerfile for containerization

Key files:

* `package.json`
* `eslint.config.js`
* `src/index.js` *(if applicable)*

---

## 2️⃣ transaction-service (Python + FastAPI)

A microservice that calculates the total sum of **positive** transactions.

Features:

* FastAPI endpoints
* Input validation
* Unit tests with pytest
* Coverage reports
* Bandit security scanning
* Dockerfile for deployment

Key files:

* `main.py`
* `test_main.py`
* `requirements.txt`

---

## 3️⃣ notification-service (Go)

A Go microservice that validates emails and handles notification requests.

Features:

* Email validation with regex
* HTTP handlers
* Request size limiting (`MaxBytesReader`)
* Unit tests for email validation and handler logic
* gosec security scanning
* Dockerfile included

Key files:

* `main.go`
* `main_test.go`
* `go.mod`

---

# 🔧 CI/CD Scripts

## 🟦 detect_changes.sh

Detects which microservices have changed by comparing the last two commits.

This ensures:

* Only modified services are linted
* Only modified services are tested
* Only modified services are scanned or built

---

## 🟩 lint.sh

Automatically selects the right linting tool per language:

| Language | Tool          |
| -------- | ------------- |
| Node.js  | ESLint        |
| Python   | Flake8        |
| Go       | golangci-lint |

---

## 🟧 test.sh

Standardized test runner:

| Language | Framework                      |
| -------- | ------------------------------ |
| Node.js  | Jest + coverage                |
| Python   | pytest + coverage (XML + HTML) |
| Go       | `go test` + coverage profile   |

---

## 🟥 scan.sh

Security scanning tailored to each ecosystem:

| Language | Tool      |
| -------- | --------- |
| Node.js  | npm audit |
| Python   | Bandit    |
| Go       | gosec     |

---

# ⚙️ Jenkins CI/CD Pipeline

The `Jenkinsfile` implements a complete multi-stage CI/CD workflow:

### Pipeline Stages:

1. **Checkout**
2. **Detect Changes** (using `detect_changes.sh`)
3. **Lint** (based on language)
4. **Security Scan**
5. **Unit Tests**
6. **Docker Build** (only for changed services)

The pipeline is modular, easy to extend, and follows DevOps best practices.

---

# 🧪 Running Tests

### Run tests for all services:

```
./test.sh
```

### Run tests for a specific service:

```
SERVICES_CHANGED="transaction-service" ./test.sh
```

---

# ▶️ Running Services Locally

## user-service (Node.js)

```
cd user-service
npm install
npm start
```

## transaction-service (Python)

```
cd transaction-service
pip install -r requirements.txt
uvicorn main:app --reload
```

## notification-service (Go)

```
cd notification-service
go run main.go
```

---

# 🐳 Docker

### Build a Docker image for any service:

```
docker build -t user-service ./user-service
```

You can repeat the same command for:

* `transaction-service`
* `notification-service`

---

# 🧰 Technologies Used

* Node.js / Express
* Python / FastAPI
* Go / net/http
* Jest, pytest, go test
* ESLint, Flake8, golangci-lint
* npm audit, Bandit, gosec
* Docker
* Jenkins CI/CD
* Git-based change detection

---

# ✅ Summary

This monorepo demonstrates a clean, scalable DevOps architecture:

✔ Multi-language microservices
✔ Unified tooling across services
✔ Automated CI/CD with selective execution
✔ Real-world compliant workflow
✔ Modular & extendable design
