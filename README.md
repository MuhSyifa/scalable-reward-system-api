# Scalable Reward System 🚀

A production-ready, highly scalable reward and loyalty management backend system built with Golang, Gin, PostgreSQL, and Redis. This project is designed to handle high-concurrency transactions like voucher redemptions while preventing race conditions (e.g., double spending) using Distributed Locks.

## 🌟 Key Features

*   **Robust Authentication**: JWT-based authentication with access and refresh tokens, plus Role-Based Access Control (Admin/User).
*   **Concurrent-Safe Redemption**: Utilizes Redis Distributed Locks (Redlock pattern) to ensure idempotency and prevent double-redemptions during high traffic.
*   **ACID Compliant Transactions**: Uses PostgreSQL database transactions to atomically deduct user points and decrease voucher stock.
*   **Idempotent APIs**: Middleware to cache and return identical responses for retried network requests.
*   **Rate Limiting**: Redis-backed rate limiter to prevent spamming and DDoS attacks.
*   **Background Workers**: Goroutine-based workers for processing expired points asynchronously.
*   **Clean Architecture**: Domain-driven design separating Delivery, Service, and Repository layers for high testability and maintainability.

## 🛠️ Tech Stack

*   **Language**: Golang (1.22+)
*   **Framework**: Gin Gonic
*   **Database**: PostgreSQL
*   **Cache & Queue**: Redis
*   **Query Builder**: sqlx
*   **Logging**: Zerolog

## 🏗️ Architecture Design

The project strictly follows **Clean Architecture** principles:
1.  **Domain Layer**: Core business entities and repository/service interfaces.
2.  **Repository Layer**: Data access layer communicating with PostgreSQL.
3.  **Service Layer**: Core business logic and transaction management.
4.  **Delivery Layer**: HTTP Handlers, Routing, and Middlewares.

## 🚀 Quick Setup (Docker)

1. Clone the repository
2. Spin up the containers
```bash
docker-compose up -d --build
```
3. The API will be available at `http://localhost:8080`. PostgreSQL on `5432` and Redis on `6379`.

## 📖 API Documentation

### Public Endpoints
*   `GET /api/v1/health` - Check system status
*   `POST /api/v1/register` - Register new user
*   `POST /api/v1/login` - Authenticate and receive JWT
*   `GET /api/v1/campaigns` - List active reward campaigns
*   `GET /api/v1/campaigns/:id/vouchers` - List vouchers for a campaign

### Protected Endpoints (Requires Bearer Token)
*   `POST /api/v1/redeem` - Redeem a voucher using points
    *   **Headers**: `Authorization: Bearer <token>`, `Idempotency-Key: <uuid>`
    *   **Body**: `{"voucher_id": "uuid"}`

## 🔮 Future Improvements
*   Implement `asynq` for robust persistent task queues.
*   Add comprehensive Unit and Integration testing for the Service layer.
*   Setup CI/CD pipeline using GitHub Actions.
