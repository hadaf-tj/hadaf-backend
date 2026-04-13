# Hadaf - Backend

## Description
This repository contains the backend service for the Hadaf platform. It provides RESTful APIs and core business logic implemented in Go, designed to facilitate targeted charitable contributions without financial transactions.

## Requirements
* Go 1.25.0
* PostgreSQL 15+ (recommended)
* Redis 7+ 
* MinIO

## Local Setup Instructions

1. **Clone the repository:**
   ```bash
   git clone <repository_url>
   cd hadaf-backend
   ```

2. **Configure Environment Variables:**
   * Create a `.env` file in the root directory based on `.env.example` (if provided).
   * Ensure `DB_URL`, `REDIS_URL`, `MINIO_ENDPOINT`, and `JWTSecretKey` are correctly set for local execution.

3. **Install Dependencies:**
   ```bash
   go mod download
   go mod tidy
   ```

4. **Run the Server:**
   ```bash
   go run main.go
   ```

## Contributing
Please note that direct commits to the `main` branch are restricted. To contribute, please fork the repository, follow our branching strategy, and submit a Pull Request.

For detailed guidelines on branch naming, commit messages, and the PR process, refer to the [Organization CONTRIBUTING.md](https://github.com/social-housing/.github/blob/main/profile/CONTRIBUTING.md).
