# Agnos Backend Assignment

## Getting Started

Clone the repository:

```bash
git clone <repository-url>
cd agnos-assignment
```

Start all services:

```bash
docker compose up -d
```

The application will be available at:

* **API:** `http://localhost:8080/api`
* **pgAdmin:** `http://localhost:8081`

## Services

* **API** — Go + Gin
* **Database** — PostgreSQL
* **Reverse Proxy** — Nginx
* **Database Management** — pgAdmin

## Stop the Application

```bash
docker compose down
```
