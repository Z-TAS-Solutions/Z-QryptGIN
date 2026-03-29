# Role-Based Access Control (RBAC) Architecture

## Overview

This document describes the role-based access control system implemented in Z-QryptGIN, specifically designed to prevent unauthorized access to admin endpoints by client token holders.

## Problem Statement

Previous architecture lacked proper role validation between JWT claims and cached session state. This could allow:
- A client JWT token to potentially access admin endpoints (if endpoints didn't validate role)
- Role tampering or cache inconsistency issues
- No cross-validation between JWT and session cache role

## Solution: Multi-Layer RBAC

### Layer 1: JWT + Session Validation (RequireAuth)
**Location**: `api/server/auth.go`

Middleware that:
1. Extracts and validates JWT signature (EdDSA cryptographic verification)
2. Checks JWT has not expired
3. Retrieves session from Redis cache using JTI (JWT ID)
4. Validates session is active and MFA-verified
5. Extracts role into Gin context

```go
func RequireAuth(jwtSvc service.JWTService, sessionRepo repository.SessionRepository) gin.HandlerFunc
```

**Cross-Validation Points**:
- ✓ JWT signature valid?
- ✓ Token not expired?
- ✓ Session exists in Redis?
- ✓ Session is active?
- ✓ MFA is verified?

**Output**: Role stored in context for next layer

### Layer 2: Role Validation (RequireRole)
**Location**: `api/server/auth.go`

Middleware that:
1. Retrieves role from context (set by RequireAuth)
2. Checks if user's role matches one of the allowed roles
3. Rejects with 403 Forbidden if no match

```go
func RequireRole(allowedRoles ...string) gin.HandlerFunc
```

**Usage Examples**:
```go
// Admin-only endpoint
protected.Use(server.RequireRole("Admin"))

// Multiple allowed roles
protected.Use(server.RequireRole("Admin", "Supervisor"))

// Client-only endpoint (though not implemented yet)
// protected.Use(server.RequireRole("Client"))
```

**Validation**: Is user's role in the allowed list?

### Layer 3: Role Consistency Validation (ValidateRoleConsistency)
**Location**: `api/server/auth.go`

Middleware that performs strict cross-validation:
1. Re-extracts JWT from Authorization header
2. Verifies JWT signature
3. Retrieves session from Redis cache
4. **Compares JWT role WITH session cache role**
5. Rejects with 403 Forbidden if mismatch

```go
func ValidateRoleConsistency(jwtSvc service.JWTService, sessionRepo repository.SessionRepository) gin.HandlerFunc
```

**Security Benefit**: Detects if:
- Token was tampered with after creation
- Redis cache was corrupted
- Role was inconsistently stored
- Session was migrated between different role contexts

**Validation**: Does JWT role == cached session role?

## Data Flow Diagram

```
HTTP Request
    ↓
URL Route Match
    ↓
RequireAuth Middleware
├─ Extract JWT from header
├─ Verify signature (EdDSA)
├─ Get session from Redis
├─ Check session active & MFA verified
└─ Set role in context
    ↓
RequireRole("Admin") Middleware
├─ Get role from context
├─ Check if role in allowed list
└─ Return 403 if not allowed
    ↓
ValidateRoleConsistency Middleware
├─ Re-extract JWT from header
├─ Verify signature (EdDSA)
├─ Get session from Redis
├─ Compare JWT.role == cachedsession.role
└─ Return 403 if mismatch
    ↓
Handler
├─ Get user_id, jti, role from context
└─ Process request
    ↓
Response
```

## Role Definition

Current roles defined in `internal/app/database/model_validation.go`:

```go
type UserRole string

const (
    RoleAdmin  UserRole = "Admin"
    RoleClient UserRole = "Client"
)
```

Can be extended with:
- "Supervisor"
- "Auditor"
- "Support"
- etc.

## Session DTO Enhancement

**File**: `internal/app/dto/session.go`

Added `Role` field to enable cross-validation:

```go
type Session struct {
    ID           string    `json:"id"`
    UserID       uint      `json:"user_id"`
    JTI          string    `json:"jti"`
    Role         string    `json:"role"`           // ← NEW
    DeviceName   string    `json:"device_name"`
    DeviceID     string    `json:"device_id"`
    IPAddress    string    `json:"ip_address"`
    UserAgent    string    `json:"user_agent"`
    IsActive     bool      `json:"is_active"`
    MfaStatus    MfaStatus `json:"mfa_status"`
    LastActiveAt time.Time `json:"last_active_at"`
    ExpiresAt    time.Time `json:"expires_at"`
    Location     string    `json:"location"`
}
```

When session is created in JWT service, role is stored in Redis:
```json
{
  "id": "session-123",
  "user_id": 45,
  "jti": "uuid-v7-string",
  "role": "Admin",
  "is_active": true,
  "mfa_status": "verified",
  ...
}
```

## Route Configuration Example

### Admin Dashboard Routes
**File**: `cmd/api/main.go`

```go
admin := router.Group("/api/v1/admin")
{
    // Registration (no auth required)
    admin.POST("/users/register/new", userRegistrationHandler.Register)
    
    // Protected admin routes
    protected := admin.Group("")
    protected.Use(server.RequireAuth(jwtService, sessionRepo))        // Layer 1
    protected.Use(server.RequireRole("Admin"))                         // Layer 2
    protected.Use(server.ValidateRoleConsistency(jwtService, sessionRepo)) // Layer 3
    {
        dashboard := protected.Group("/dashboard")
        {
            dashboard.GET("/auth-trends", dashboardHandler.GetAuthenticationTrends)
            dashboard.GET("/metrics", dashboardHandler.GetDashboardMetrics)
        }
    }
}
```

### Client User Routes
**File**: `cmd/api/main.go`

```go
user := router.Group("/api/v1/user")
{
    // Protected client routes (no admin access)
    protected := user.Group("")
    protected.Use(server.RequireAuth(jwtService, sessionRepo))
    protected.Use(server.RequireRole("Client"))  // ← Clients only
    {
        protected.GET("/notifications", userHandler.GetNotifications)
        protected.GET("/sessions", sessionHandler.GetActiveSessions)
    }
}
```

## Accessing Role in Handlers

Handlers can access the user's role from Gin context:

```go
func (h *DashboardHandler) GetAuthenticationTrends(c *gin.Context) {
    // Already validated to be Admin by middleware
    
    userIDInterface, _ := c.Get("user_id")
    userID, _ := userIDInterface.(uint)
    
    roleInterface, _ := c.Get("role")
    role, _ := roleInterface.(string)  // Always "Admin" at this point
    
    // Process request knowing user is definitely an admin
}
```

## Security Assumptions & Guarantees

### What We Guarantee
- ✓ User is authentic (JWT signature verified)
- ✓ Session was not revoked (Redis check)
- ✓ User has required role
- ✓ User's role hasn't been tampered with
- ✓ MFA has been completed
- ✓ JWT hasn't expired
- ✓ Admin endpoints are protected from client users

### What We Don't Guarantee
- ✗ User still exists in database (no DB check)
- ✗ User's status is still "Active" (no DB check)
- ✗ Permissions beyond role (no fine-grained ACL)

For stronger guarantees, add additional checks in handlers:
```go
user, err := h.userRepo.FindByID(userID)
if err != nil || user.Status != "Active" {
    return errors.New("user inactive or deleted")
}
```

## Testing Role-Based Access

### Test 1: Admin Can Access Admin Endpoint
```bash
# Create admin JWT
admin_token="..." # Admin role in JWT

curl -H "Authorization: Bearer $admin_token" \
  http://localhost:8080/api/v1/admin/dashboard/auth-trends
# Expected: 200 OK with data
```

### Test 2: Client Cannot Access Admin Endpoint
```bash
# Create client JWT
client_token="..." # Client role in JWT

curl -H "Authorization: Bearer $client_token" \
  http://localhost:8080/api/v1/admin/dashboard/auth-trends
# Expected: 403 Forbidden
```

### Test 3: Tampered Token Detected
```bash
# Modify JWT manually (change role Admin → Client)
tampered_token="..."

curl -H "Authorization: Bearer $tampered_token" \
  http://localhost:8080/api/v1/admin/dashboard/auth-trends
# Expected: 403 Forbidden (ValidateRoleConsistency detects mismatch)
```

### Test 4: Revoked Session Blocked
```bash
# Revoke session in Redis manually
# redis-cli DEL session:jti-value

curl -H "Authorization: Bearer $token" \
  http://localhost:8080/api/v1/admin/dashboard/auth-trends
# Expected: 401 Unauthorized (session not found)
```

## Future Enhancements

1. **Fine-Grained Permissions**: Add permission-level checks (e.g., "admin.read.auth-trends")
2. **Role Inheritance**: Define role hierarchies (e.g., Admin inherits Client permissions)
3. **Custom Middleware per Handler**: Some handlers might need different role checks
4. **Audit Logging**: Log all RBAC violations for security monitoring
5. **Dynamic Role Loading**: Load roles/permissions from database instead of hardcoding
6. **Rate Limiting per Role**: Different rate limits for admin vs client
7. **Least Privilege SPA Tokens**: Shorter expiry for admin actions

## Troubleshooting

### Getting 403 Forbidden with valid token?
1. Check JWT role matches session cache role (verify both set correctly)
2. Ensure MFA status is "verified" in session cache
3. Verify session hasn't been manually revoked in Redis
4. Check endpoint middleware order (RequireAuth before RequireRole)

### Getting 401 Unauthorized unexpectedly?
1. Check JWT hasn't expired (exp claim)
2. Verify session still exists in Redis (didn't time out)
3. Check session.is_active is true
4. Verify MfaStatus == "verified"

### Admin token works but client token doesn't work for client route?
1. Might have "Admin" role but endpoint expects "Client"
2. Some endpoints might allow both roles (check RequireRole call)
3. Verify client route actually has RequireRole middleware
