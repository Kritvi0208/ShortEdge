# 🔗 ShortEdge: URL Shortener with Insights

**ShortEdge** is a powerful, production-ready URL shortener built with the [GoFr v1.42.0](https://gofr.dev) framework. It supports branded links, real-time analytics, and full REST API access — making it a great choice for teams that need full control, customization, and data privacy without relying on third-party platforms.

---

## 📈 Features

| ✅ Feature                  | 📌 Description                                                                 |
|----------------------------|--------------------------------------------------------------------------------|
| 🔗 Branded Short URLs       | Create short links with custom codes or aliases                               |
| 📊 Real-Time Analytics      | Track IP-based location, browser, and device info on every visit              |
| 🔐 Public/Private Toggle    | Control visibility of links and analytics data                                |
| 🧾 Full REST API            | Clean, CRUD-complete API for developers                                       |
| ⏳ Link Expiry              | Set optional expiration for time-bound links                                  |
| 🧬 Swagger UI               | Interactive API documentation via Swagger                                     |
| 📈 Prometheus Integration   | Monitor performance with built-in `/metrics` endpoint                         |
| 🌐 Basic Frontend           | Minimal UI to shorten and manage links                                        |
| ❤️ Health Check            | `/health` endpoint for liveness and deployment readiness                       |


---


## 📂 Project Structure

```
ShortEdge/
├── cmd/
│   └── main.go              # Application entrypoint, initializes GoFr app
│
├── handler/                 # API route handlers (HTTP layer)
│   ├── url.go               # Handles URL shortening, update, delete, redirect
│   └── visit.go             # Handles visit analytics endpoints
│
├── service/                 # Business logic (middle layer)
│   ├── url.go               # Short link creation, update, expiry logic
│   └── visit.go             # Processes visitor tracking and analytics
│
├── store/                   # Data access layer (PostgreSQL queries)
│   ├── url.go               # DB methods for creating, updating, fetching URLs
│   └── visit.go             # DB methods for storing/retrieving visit logs
│
├── model/                   # Domain models (data representations)
│   └── visit.go             # Visit model struct (IP, country, device, etc.)
│
├── factory/                 # Dependency injection setup
│   └── store.go             # Initializes DB and returns store interfaces
│
├── middleware/              # Optional middleware (auth, logging, etc.)
│   └── middleware.go        # Sample middleware handler
│
├── static/                  # Frontend static files (HTML + JS + CSS)
│   ├── script.js            # Shared JS functions
│   ├── style.css            # Basic styles for frontend
│   └── scripts/
│       ├── dashboard.js     # Handles dashboard UI logic
│       ├── shorten.js       # Shorten form behavior
│       └── token.js         # Token/localStorage helper
│
├── docs/                    # Swagger documentation
│   ├── swagger.yaml         # OpenAPI spec in YAML
│   ├── swagger.json         # Auto-generated Swagger JSON
│   └── docs.go              # Swagger annotations (via swaggo)
│
├── .env                     # Environment variables (DB config)
├── .gitignore               # Git ignored files list
├── go.mod                   # Go module file
├── go.sum                   # Go dependencies checksum
└── README.md                # Project documentation
```
---

## 🔧 Setup & Run

### 1. Clone Repo

```
git clone https://github.com/Kritvi0208/ShortEdge.git
cd ShortEdge
```

### 2. Setup .env
Create a .env file:
```
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=yourpassword
DB_NAME=shortedge
```
### 3. Run Server
```
go run ./cmd/main.go
```

### 4. Access

_Swagger UI:_ http://localhost:8000/swagger/index.html

_Frontend:_ http://localhost:8000

_Metrics:_ http://localhost:8000/metrics

_Health:_ 	http://localhost:8000/health

---
## 🌐 API Endpoints

| 🧭 Method | 🌍 Endpoint              | 📌 Purpose                                 |
|----------|--------------------------|--------------------------------------------|
| `POST`   | `/shorten`               | Create a short/branded URL                 |
| `GET`    | `/{code}`                | Redirect to original URL                   |
| `GET`    | `/analytics/{code}`      | View visit analytics for a short URL       |
| `GET`    | `/all`                   | List all shortened URLs                    |
| `PUT`    | `/update/{code}`         | Edit long URL or toggle visibility         |
| `DELETE`| `/delete/{code}`          | Delete a short URL                         |
| `GET`    | `/health`                | Health check for deployment                |
| `GET`    | `/metrics`               | Prometheus metrics for observability       |
| `GET`    | `/swagger/index.html`    | Interactive Swagger API documentation      |

---
## 🧱 Architecture
```
                  ┌────────────────────────┐                         ┌─────────────────────────────────────────────┐
                  │     External Clients   │                         │  🌐 Users send HTTP requests via frontend,  
                  │  (Browser / Postman)   │                         │     curl, or tools like Postman               
                  └───────────┬────────────┘                         └─────────────────────────────────────────────┘
                              │
                       [HTTP Requests]
                              ▼
                ┌────────────────────────────┐                       ┌─────────────────────────────────────────────┐
                │      Handler Layer         │                       │  🎯 Handles routes, binds JSON, performs   
                │  Accepts + Validates input │                       │     validations, and calls service layer    
                │  Calls service functions   │                       └─────────────────────────────────────────────┘
                └────────────┬───────────────┘
                             ▼
                ┌────────────────────────────┐                       ┌─────────────────────────────────────────────┐
                │      Service Layer         │                       │  🧠 Core business logic lives here: expiry 
                │  Business rules + logic    │                       │     handling, privacy, branding, analytics  
                │  Prepares models, logic    │                       └─────────────────────────────────────────────┘
                └────────────┬───────────────┘
                             ▼
                ┌────────────────────────────┐                       ┌─────────────────────────────────────────────┐
                │       Store Layer          │                       │  🗄️ Executes raw SQL queries against        
                │  Executes SQL/Postgres     │                       │     PostgreSQL using GoFr DB utilities      
                │  Interfaces for mocking    │                       │     Returns typed models to the service     
                └────────────────────────────┘                       └─────────────────────────────────────────────┘
```
---

## 📚 SOLID Design Principles

| 🧠 Principle                 | ✅ How It’s Used in ShortEdge                                        |
|-----------------------------|----------------------------------------------------------------------|
| **S - Single Responsibility** | Handler, Service, and Store layers are cleanly separated             |
| **O - Open/Closed**           | Extend features via interfaces without changing core logic          |
| **L - Liskov Substitution**   | Interfaces allow mock substitution during testing                    |
| **I - Interface Segregation** | Each service/store uses only what it needs — no bloated interfaces  |
| **D - Dependency Inversion**  | Factory injects DB abstraction, not hardcoded Postgres code         |

---

## ⚙️ Tech Stack

| 🛠️ Tool/Technology | 📌 Role in Project                                         |
|--------------------|-----------------------------------------------------------|
| **Go**             | Primary backend programming language                      |
| **GoFr**           | Scalable backend framework with clean architecture        |
| **PostgreSQL**     | Relational database to persist short links and analytics  |
| **Swagger (swaggo)**| API documentation generator (served at `/swagger/...`)   |
| **Prometheus**     | Built-in monitoring at `/metrics` for observability       |
| **ipwho.is**       | IP geolocation API (used to resolve country from IP)      |
| **uasurfer**       | Parses `User-Agent` to extract browser & device details   |
| **HTML/CSS/JS**    | Basic static frontend for shortening & managing URLs      |

---
**GoFr Summer of Code 2025 — Assignment 7: URL Shortener with Insights**
---
Licensed under the MIT License — free to use, modify, and distribute.
