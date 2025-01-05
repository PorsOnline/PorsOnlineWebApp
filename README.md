# Online Survey System

This project is a Golang-based online survey system that uses Docker for containerization. The application connects to two PostgreSQL databases and serves on port 8080.

For installation follow the instructions below:

### 1.Clone the Repository
```bash
git clone <repository-url>
cd <repository-name>
```
### 3. set the env variables in .env file
```bash
POSTGRES_HOST=localhost
POSTGRES_USER=you-username
POSTGRES_PASSWORD=your-password
POSTGRES_SecretDB_NAME=your-secrets-database
POSTGRES_DB_NAME=your-db-database
POSTGRES_PORT=5432
POSTGRES_SecretDB_PORT=5433
CONFIG_PATH=./config.json
GOOGLE_SMTP_PASSWORD=""
```
### 2.Run the program
```bash
docker compose up -d --build
```