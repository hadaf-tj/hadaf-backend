# Hadaf - Backend

## Description
This repository contains the backend service for the Hadaf platform. It provides RESTful APIs and core business logic implemented in Go, designed to facilitate targeted charitable contributions without financial transactions.

## Local Setup Instructions (100% Working Method)

The most reliable and recommended way to run the backend locally is using **Docker Compose**. The backend relies on multiple services (PostgreSQL, Redis, MinIO), and Docker automatically provisions all of them along with the Go application itself.

### Prerequisites
* Docker and Docker Compose installed.

### Steps
1. **Clone the repository:**
   ```bash
   git clone https://github.com/hadaf-tj/hadaf-backend
   cd hadaf-backend
   ```

2. **Configure Environment Variables:**
   Create a `.env` file in the root directory based on `.env_example`:
   ```bash
   cp .env_example .env
   ```
   *Note: Ensure all required keys (e.g., `GOOGLE_CLIENT_ID`, `GOOGLE_CLIENT_SECRET`) are filled in properly for OAuth to work.*

3. **Start the Application:**
   Run the following command to build the image and start all services in the background:
   ```bash
   docker compose up -d --build
   ```
   *The backend API will be available at `http://localhost:8000`.*
   *Database migrations are applied automatically on startup.*

4. **View Logs (Optional):**
   To monitor the backend logs:
   ```bash
   docker compose logs -f app
   ```

5. **Stop the Application:**
   ```bash
   docker compose down
   ```

---

## Architecture Overview

The Hadaf backend is built as a modular monolith in Go, following clean architecture principles. This structure is designed to be highly maintainable and easily extensible as new features are added.

### Tech Stack
- **Language:** Go (Golang)
- **Database:** PostgreSQL (Relational data, User profiles, Needs, Promises)
- **Caching & OTP:** Redis
- **Object Storage:** MinIO (S3-compatible storage for avatars and file uploads)
- **Migrations:** `golang-migrate`

### System Layers
The architecture is divided into distinct layers, allowing developers to add new features without disrupting existing logic:

1. **Routing & Handlers (`internal/handlers/`)**
   - Built with the `gin-gonic/gin` framework.
   - Responsible for receiving HTTP requests, extracting payloads, and returning JSON responses.
   - *How to extend:* When adding a new feature (e.g., "Comments"), create a new handler here and register its routes.

2. **Business Logic (`internal/services/`)**
   - Contains the core logic of the application (e.g., OAuth validation, token generation, mapping needs to institutions).
   - *How to extend:* Add new service methods here. Handlers will call these services rather than interacting with the database directly.

3. **Data Access (`internal/repository/`)**
   - Manages all interactions with external state (PostgreSQL, Redis, MinIO).
   - *How to extend:* Write new SQL queries or storage commands here. Services will call the repository layer to fetch or save data.

4. **Models (`internal/models/`)**
   - Shared data structures and domain entities used across all layers.

### Adding New Features
To add a new feature to the platform:
1. Define the data structure in `internal/models/`.
2. Create a new database migration in `migration/`.
3. Implement the database queries in `internal/repository/`.
4. Add the business logic in `internal/services/`.
5. Expose the REST API endpoint in `internal/handlers/`.

## Database Management & Seeding

During development, you may need to interact with the database directly or load test data.

### Accessing the Database Manually
To open an interactive PostgreSQL console (`psql`) inside the running Docker container:
```bash
docker compose exec postgres psql -U postgres -d shb
```
*(You can exit the console by typing `\q` and pressing Enter).*

### Loading Seed Data
If you need to populate the database with initial test data, you can execute the provided seed script (`migration/seed.sql`) against the running database:
```bash
cat migration/seed.sql | docker compose exec -T postgres psql -U postgres -d shb
```

---

## Contributing
Please note that direct commits to the `main` branch are restricted. To contribute, please fork the repository, follow our branching strategy, and submit a Pull Request.

For detailed guidelines on branch naming, commit messages, and the PR process, refer to the [Organization CONTRIBUTING.md](https://github.com/social-housing/.github/blob/main/profile/CONTRIBUTING.md).
