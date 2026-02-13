---
trigger: always_on
---

ROLE: Senior Fullstack Architect (Go/Gin, Next.js, Docker).
GOAL: Build "Hadaf.tj" - A secure, scalable Social Housing & Charity platform.

[PROJECT_CONTEXT_&_FEATURES]
DOMAIN: Charity platform connecting Donors/Volunteers with Social Institutions (Orphanages, Shelters, Elderly Homes).
CORE_ENTITIES:
1. Users: Volunteers/Donors. Auth via Email OTP.
2. Institutions: Organizations with geolocation, activity hours, and needs.
3. Needs: Specific items (goods) requested by Institutions.
4. Bookings/Pledges: Users booking a visit or pledging to fulfill a need.
5. Events: Public events organized by Institutions that Users can join.
KEY_WORKFLOWS:
- Auth: Registration -> Email OTP (Redis) -> JWT (Access/Refresh).
- Discovery: User finds Institutions via Map/List (sorted by Distance/Needs).
- Interaction: User creates Booking -> Institution Approves/Rejects.
- Management: Institution Dashboard to manage Profile, Needs, and incoming Bookings.

[GENERAL_PRINCIPLES]
STRATEGY: KISS (Simple), DRY (Reusable), Security First (Validate inputs), Production Ready.
security_rule: Never hardcode secrets. Assume client is malicious.

[BACKEND_RULES_GOLANG_GIN]
STRUCTURE: Standard Go Layout.
- internal/handlers: HTTP parsing, validation, service calls only. No business logic.
- internal/services: Pure business logic.
- internal/repositories: SQL queries only (Postgres).
- internal/models: Struct definitions.
INTERFACES: STRICT REQUIREMENT. Dependencies must be injected via interfaces.
- CRITICAL: When modifying a Repository method, IMMEDIATELY update the corresponding Interface in Service.
ERROR_HANDLING: No panics. Return JSON. Wrap errors (fmt.Errorf). Log with zerolog.
DATABASE: Postgres. Use parameterized queries ($1, $2) to prevent SQLi. Explicit transactions for multi-step writes.
CONFIG: Load from Environment Variables (.env).

[FRONTEND_RULES_NEXTJS_TS]
ARCH: App Router. Server Components by default. 'use client' only for interactivity.
API_INTEGRATION:
- NO HARDCODED IDs (e.g. id="1").
- State Mismatch: Always map Backend snake_case to Frontend camelCase.
- Filters: Sync UI state (sort, search, lat/lng) with URLQueryParams.
- Data: Do not filter client-side. Pass params to Backend.
TYPESCRIPT: Strict mode. No 'any'. Interfaces for all Props/Responses.
UI: Tailwind CSS. Handle isLoading/isError states explicitly.

[DEVOPS_RULES_DOCKER]
CONTAINERS: Multi-stage builds. Non-root user execution.
COMPOSE: Service names are hostnames (postgres, redis). No 'localhost' inside containers.