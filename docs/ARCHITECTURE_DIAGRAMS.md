# Architecture Diagrams

## 1. Request Flow with Middleware Stack

```
┌─────────────────────────────────────────────────────────────────┐
│ HTTP Request: GET /api/v1/admin/dashboard/auth-trends           │
│ Header: Authorization: Bearer <jwt_token>                       │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
              ┌───────────────────────────────┐
              │  Route Matching: /api/v1/admin│
              └───────────────────────────────┘
                              │
                              ▼
        ╔═══════════════════════════════════════════╗
        ║    MIDDLEWARE LAYER 1: RequireAuth        ║
        ║                                           ║
        ║  1. Extract JWT from Authorization header ║
        ║  2. Verify EdDSA signature                ║
        ║  3. Check token not expired               ║
        ║  4. GET session from Redis using JTI      ║
        ║  5. Verify session.is_active = true       ║
        ║  6. Verify session.mfa_status = verified  ║
        ║  7. Set context: user_id, jti, role       ║
        ║                                           ║
        ║  ✓ PASS: Continue  ✗ FAIL: 401 Unauthed  ║
        ╚═══════════════════════════════════════════╝
                              │
                              ▼
        ╔═══════════════════════════════════════════╗
        ║    MIDDLEWARE LAYER 2: RequireRole        ║
        ║                                           ║
        ║  1. Get role from context (set by Layer 1)│
        ║  2. Check if role in allowed_roles        ║
        ║  3. For admin: check role == "Admin"      ║
        ║                                           ║
        ║  ✓ PASS: Continue  ✗ FAIL: 403 Forbidden ║
        ╚═══════════════════════════════════════════╝
                              │
                              ▼
        ╔═══════════════════════════════════════════╗
        ║ MIDDLEWARE LAYER 3: ValidateConsistency   ║
        ║                                           ║
        ║  1. Re-extract JWT from header            ║
        ║  2. Verify EdDSA signature                ║
        ║  3. GET session from Redis                ║
        ║  4. Compare JWT.role == session.role      ║
        ║  5. If mismatch: possible tampering       ║
        ║                                           ║
        ║  ✓ PASS: Continue  ✗ FAIL: 403 Forbidden ║
        ╚═══════════════════════════════════════════╝
                              │
                              ▼
              ┌───────────────────────────────┐
              │     HTTP Handler Execution    │
              │                               │
              │ DashboardHandler.             │
              │ GetAuthenticationTrends()     │
              └───────────────────────────────┘
                              │
                              ▼
              ┌───────────────────────────────┐
              │    Business Logic Layer       │
              │                               │
              │ DashboardService.             │
              │ GetAuthenticationTrends()     │
              │ - Validate interval param     │
              │ - Call repository             │
              └───────────────────────────────┘
                              │
                              ▼
              ┌───────────────────────────────┐
              │    Data Access Layer          │
              │                               │
              │ DashboardRepository.          │
              │ GetAuthTrendsByInterval()     │
              │ - Query ActivityLog table     │
              │ - Aggregate by DATE_TRUNC    │
              │ - Fill missing intervals      │
              └───────────────────────────────┘
                              │
                              ▼
                    ┌──────────────────┐
                    │  PostgreSQL DB   │
                    │                  │
                    │  ActivityLog     │
                    │  - id            │
                    │  - user_id       │
                    │  - type          │
                    │  - created_at    │
                    └──────────────────┘
                              │
                              ▼
              ┌───────────────────────────────┐
              │     Data Processing Layer     │
              │                               │
              │ - Sort by timestamp           │
              │ - Aggregate success/failure   │
              │ - Fill missing time buckets   │
              └───────────────────────────────┘
                              │
                              ▼
              ┌───────────────────────────────┐
              │      Response Building        │
              │                               │
              │ {                             │
              │   "interval": "hour",         │
              │   "data": [                   │
              │     {                         │
              │       "timestamp": "...",     │
              │       "successCount": 120,    │
              │       "failureCount": 5       │
              │     }                         │
              │   ]                           │
              │ }                             │
              └───────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│ HTTP Response: 200 OK                                           │
│ Body: Time-series authentication trends data                   │
└─────────────────────────────────────────────────────────────────┘
```

## 2. Dependency Injection Graph

```
┌────────────────────────────────────────────────────────┐
│                   External Resources                   │
├────────────────────────────────────────────────────────┤
│  PostgreSQL DB Connection          Redis Connection   │
└────────────────┬───────────────────────────┬───────────┘
                 │                           │
                 ▼                           ▼
        ┌─────────────────┐        ┌──────────────────┐
        │  config.Database│        │  config.Redis    │
        └────────┬────────┘        └────────┬─────────┘
                 │                         │
                 ▼                         ▼
    ┌──────────────────────────┐  ┌─────────────────┐
    │    UserRepository        │  │ SessionRepository
    │    WebAuthnCredential    │  │ Used for JWT DI  │
    │    NotificationRepository│  └────────┬────────┘
    │    DashboardRepository ◄─────────────┤ Uses Redis
    │    (NEW)               │             │
    └────────┬───────────────┘             │
             │                              │
             ▼                              │
    ┌──────────────────────────┐           │
    │   DashboardService       │◄──────────┘
    │   (NEW)                  │
    │   - GetAuthenticationTrends
    │   - GetDashboardMetrics  │
    └────────┬───────────────┘
             │
             ▼
    ┌──────────────────────────┐
    │  DashboardHandler        │
    │  (NEW)                   │
    │  - GetAuthenticationTrends
    │  - GetDashboardMetrics   │
    └────────┬───────────────┘
             │
             ▼
         Routes
       (main.go)
```

## 3. Role Authorization Decision Tree

```
                        ┌─ Admin Access Request ─┐
                        │ (Get auth-trends)      │
                        └───────────┬────────────┘
                                    │
                              ┌─────▼────────┐
                              │ Decode JWT   │
                              └─────┬────────┘
                                    │
                          ┌─────────▼──────────┐
                          │ Signature Valid?   │
                          └─┬──────────────┬──┘
                     Valid │              │ Invalid
                          ▼              ▼
                  ┌──────────────┐  ┌─────────────┐
                  │ Get Session  │  │ 401 Error   │
                  │ from Redis   │  └─────────────┘
                  └─┬────────────┘
                    │
           ┌────────▼─────────┐
           │ Session Found?   │
           └┬──────────────┬──┘
      Yes  │              │ No
           ▼              ▼
    ┌──────────────┐  ┌─────────────┐
    │ Check Status │  │ 401 Error   │
    ├──────────────┤  └─────────────┘
    │ Active?      │
    │ MFA Verified?│
    └┬────────────┬┘
   Yes│           │ No
      ▼           ▼
  ┌─────────┐ ┌──────────────┐
  │Get role │ │ 401 Error    │
  └┬────────┘ └──────────────┘
   │
   ▼
  ┌──────────────────────┐
  │ Role in allowed list?│
  │ (e.g., "Admin")      │
  └┬────────────────┬───┘
 Yes│               │ No
   ▼               ▼
┌──────────────┐ ┌──────────────────┐
│Verify role   │ │ 403 Forbidden    │
│consistency   │ │ (Client tried    │
│with Redis    │ │  access Admin)   │
└┬──────────┬──┘ └──────────────────┘
 │          │
 │    Match │
 │          │
 ▼          ▼
┌────────┐ ┌──────────────────────┐
│ GRANT  │ │ 403 Forbidden        │
│ ACCESS │ │ (Role mismatch -     │
│        │ │  possible tampering) │
└────────┘ └──────────────────────┘
```

## 4. Data Flow: ActivityLog → Response

```
PostgreSQL
ActivityLog Table
    │
    │ WHERE:
    │   - created_at >= NOW() - 24 hours
    │   - type IN ('Login_Success', 'Failed_Login')
    │
    ▼
┌─────────────────────────────────┐
│ Raw Activity Logs (Example)     │
├─────────────────────────────────┤
│ id │ created_at        │ type   │
│ 1  │ 2026-03-29 00:05  │ Success│
│ 2  │ 2026-03-29 00:12  │ Failed │
│ 3  │ 2026-03-29 00:58  │ Success│
│ 4  │ 2026-03-29 01:03  │ Success│
│ 5  │ 2026-03-29 01:15  │ Success│
│ 6  │ 2026-03-29 02:00  │ Failed │
└─────────────────────────────────┘
    │
    │ DATE_TRUNC('hour', created_at)
    │ GROUP BY time_bucket, type
    │ COUNT(*)
    │
    ▼
┌─────────────────────────────────┐
│ Aggregated by Hour              │
├─────────────────────────────────┤
│ hour        │ type    │ count   │
│ 00:00:00    │ Success │ 2       │
│ 00:00:00    │ Failed  │ 1       │
│ 01:00:00    │ Success │ 2       │
│ 02:00:00    │ Failed  │ 1       │
└─────────────────────────────────┘
    │
    │ Combine success/failure per hour
    │ Fill missing hours with zeros
    │
    ▼
┌─────────────────────────────────┐
│ Time-Series Data Points         │
├─────────────────────────────────┤
│ timestamp  │ success │ failure  │
│ 00:00:00   │ 2       │ 1        │
│ 01:00:00   │ 2       │ 0        │
│ 02:00:00   │ 0       │ 1        │
│ 03:00:00   │ 0       │ 0        │
│ ...        │ ...     │ ...      │
└─────────────────────────────────┘
    │
    │ JSON serialization
    │
    ▼
┌─────────────────────────────────┐
│ API Response                    │
├─────────────────────────────────┤
│ {                               │
│   "interval": "hour",           │
│   "data": [                     │
│     {                           │
│       "timestamp": "2026-...",   │
│       "successCount": 2,        │
│       "failureCount": 1         │
│     },                          │
│     ...                         │
│   ]                             │
│ }                               │
└─────────────────────────────────┘
    │
    │ HTTP 200 OK
    │
    ▼
Browser/Client
│
▼ Visualizes as:
├─ Line Chart (success trend)
├─ Area Chart (success vs failure)
├─ Bar Chart (comparisons)
└─ Data Table (drill-down)
```

## 5. Component Interaction Diagram

```
┌────────────────────────────────────────────────────────────┐
│                        Gin Router                          │
│                                                            │
│    /api/v1/admin/dashboard/auth-trends [GET]             │
└────────────────┬─────────────────────────────────────────┘
                 │
    ┌────────────▼───────────────┐
    │   Middleware Chain         │
    ├────────────────────────────┤
    │ 1. RequireAuth             │
    │    ├─ JWTService           │
    │    └─ SessionRepository    │
    │ 2. RequireRole("Admin")    │
    │ 3. ValidateRoleConsistency │
    │    ├─ JWTService           │
    │    └─ SessionRepository    │
    └────────────┬───────────────┘
                 │
    ┌────────────▼───────────────────────┐
    │     DashboardHandler               │
    │  (HTTP Layer)                      │
    ├────────────────────────────────────┤
    │ + GetAuthenticationTrends()        │
    │ + GetDashboardMetrics()            │
    │                                    │
    │ Responsibilities:                  │
    │ - Parse query parameters           │
    │ - Call service                     │
    │ - Format HTTP response             │
    │ - Handle HTTP errors               │
    └────────────┬───────────────────────┘
                 │
    ┌────────────▼───────────────────────┐
    │     DashboardService               │
    │  (Business Logic Layer)            │
    ├────────────────────────────────────┤
    │ + GetAuthenticationTrends()        │
    │ + GetDashboardMetrics()            │
    │                                    │
    │ Responsibilities:                  │
    │ - Validate parameters              │
    │ - Business logic                   │
    │ - Orchestrate repository calls     │
    │ - Data transformation              │
    └────────────┬───────────────────────┘
                 │
    ┌────────────▼─────────────────────────┐
    │    DashboardRepository               │
    │  (Data Access Layer)                 │
    ├──────────────────────────────────────┤
    │ + GetAuthTrendsByInterval()          │
    │ + GetAuthTrendsMetrics()             │
    │ + fillMissingIntervals()             │
    │                                      │
    │ Responsibilities:                    │
    │ - Database queries                   │
    │ - Data aggregation                   │
    │ - Time-series processing             │
    │ - Gap filling                        │
    └────────────┬──────────────────────────┘
                 │
    ┌────────────▼────────────────────┐
    │  PostgreSQL Database            │
    │                                 │
    │  Tables:                        │
    │  ├─ activity_logs (queries)      │
    │  ├─ sessions (via FK)            │
    │  └─ users (via FK)              │
    │                                 │
    │  Indexes:                       │
    │  ├─ created_at                  │
    │  ├─ type                        │
    │  └─ user_id                     │
    └─────────────────────────────────┘
```

## 6. Security Validation Points

```
Request with JWT Token
         │
         ▼
    ┌────────────────────────┐
    │ VALIDATION POINT 1:    │
    │ JWT Signature Check    │
    │ (EdDSA Cryptography)   │
    └────┬───────────────────┘
         │
    ┌────▼────────────────────┐
    │ VALIDATION POINT 2:     │
    │ Token Expiration Check  │
    │ (exp claim)             │
    └────┬───────────────────┘
         │
    ┌────▼──────────────────────┐
    │ VALIDATION POINT 3:       │
    │ Session Exists in Redis   │
    │ (Revocation check)        │
    └────┬──────────────────────┘
         │
    ┌────▼──────────────────────┐
    │ VALIDATION POINT 4:       │
    │ Session Status Check      │
    │ (is_active, MFA verified) │
    └────┬──────────────────────┘
         │
    ┌────▼──────────────────────┐
    │ VALIDATION POINT 5:       │
    │ Role Authorization Check  │
    │ (role in allowed list)    │
    └────┬──────────────────────┘
         │
    ┌────▼──────────────────────┐
    │ VALIDATION POINT 6:       │
    │ Role Consistency Check    │
    │ (JWT.role == Cache.role)  │
    └────┬──────────────────────┘
         │
    ┌────▼──────────────────────┐
    │ All Checks Passed ✓       │
    │ → Grant Access            │
    └──────────────────────────┘
```

## 7. Request Journey Timeline

```
T+0ms    - HTTP Request arrived
         - GET /api/v1/admin/dashboard/auth-trends?interval=hour

T+1ms    - Route matched
         - Middleware chain starts

T+2ms    - RequireAuth: Extract JWT from header
T+3ms    - RequireAuth: Verify EdDSA signature
T+5ms    - RequireAuth: Redis GET session:<jti>
T+15ms   - RequireAuth: Validate session state
T+16ms   - RequireAuth: Set context variables
         - RequireRole: Compare role with "Admin"
T+17ms   - ValidateRoleConsistency: Re-verify JWT
T+19ms   - ValidateRoleConsistency: Redis GET session
T+29ms   - ValidateRoleConsistency: Compare roles
T+30ms   - Handler execution begins

T+31ms   - DashboardHandler: Parse query params
T+32ms   - Call DashboardService

T+33ms   - DashboardService: Validate interval
T+34ms   - Call DashboardRepository

T+35ms   - DashboardRepository: Build SQL query
T+40ms   - PostgreSQL: Execute DATE_TRUNC query
T+80ms   - PostgreSQL: Return results (e.g., 24 rows)
T+81ms   - Repository: Aggregate success/failure
T+82ms   - Repository: Fill missing intervals
T+83ms   - Repository: Sort by timestamp
T+84ms   - Return to Service

T+85ms   - Service: Return GetAuthTrendsResponse
T+86ms   - Return to Handler

T+87ms   - Handler: JSON serialize response
T+88ms   - HTTP 200 OK response sent

Total: ~88ms for full request cycle
```
