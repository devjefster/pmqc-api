`# 🚀 PMQC API

### A Go-based REST API for fetching and storing **PMQC (Programa de Monitoramento da Qualidade dos Combustíveis)** data from the Brazilian government.

## 📌 Features
a
- ✅ Fetches PMQC JSON data from the **ANP API** (Agência Nacional do Petróleo)
- ✅ Stores PMQC data in **PostgreSQL**
- ✅ Uses **Gin** for REST API routing
- ✅ Implements **Swagger (OpenAPI) documentation**
- ✅ Supports **parallel fetching & batch inserts**
- ✅ Includes **logging & retries for robustness**
- ✅ Uses **Docker & Docker Compose** for easy deployment

---
## 🛠️ **Installation & Setup**

### **1️⃣ Clone the Repository**
```sh
git clone https://github.com/devjefster/pmqc-api.git
cd pmqc-api` 
```
### **2️⃣ Install Go Dependencies**
```sh
go mod tidy
```
### **3️⃣ Setup Environment Variables**

Create a `.env` file with:
`DATABASE_URL=postgres://user:password@db:5432/pmqc  LOG_LEVEL=info`

----------

## 🚀 **Running the API**

### **1️⃣ Running Locally**
Ensure **PostgreSQL** is running, then:
`go run main.go`

### **2️⃣ Running with Docker**

Build and start services:
`docker-compose up --build -d`

Verify that the API is running:
`curl http://localhost:8080/ping`

----------

## 🔥 **API Endpoints**

### **📌 Fetcher API**


Fetch PMQC data for a specific year/month
`GET /fetch?year=2024&month=10`

Fetch all available PMQC data
`GET /fetch/all`
### **📌 Storage API**

Get all stored amostras (supports pagination & filters)
`GET /amostras`

Get a specific amostra by ID
`GET /amostras/{id}`

----------

## 📖 **Swagger API Documentation**

API documentation is available at:
`http://localhost:8080/swagger/index.html`

To regenerate Swagger docs:
`swag init --parseDependency --parseInternal`

----------
## 🏗 **Deployment**

### **Deploying with Docker Compose**
`docker-compose up -d --build`

### **Deploying to Production**

1.  Set up a **PostgreSQL database**
2.  Use **Docker or Kubernetes**
3.  Configure **environment variables** (`DATABASE_URL`, `LOG_LEVEL`)
4.  Run:
    `go build -o pmqc-api ./main.go ./pmqc-api`
----------

## 🤝 **Contributing**

1.  **Fork** the project
2.  **Create a feature branch** (`git checkout -b feature-xyz`)
3.  **Commit changes** (`git commit -m "Added new feature XYZ"`)
4.  **Push to GitHub** (`git push origin feature-xyz`)
5.  **Create a Pull Request**

----------

## 📜 **License**

This project is licensed under the **MIT License**.

----------