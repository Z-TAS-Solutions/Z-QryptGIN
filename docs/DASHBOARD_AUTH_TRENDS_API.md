## Dashboard Authentication Trends - Usage Guide

### Endpoint
```
GET /api/v1/admin/dashboard/auth-trends
```

### Authentication
- **Required**: Bearer JWT token in Authorization header
- **Role Required**: "Admin"
- **Session Requirement**: MFA must be verified

### Request

#### Headers
```
Authorization: Bearer <jwt_token>
Content-Type: application/json
```

#### Query Parameters
| Parameter | Type | Default | Options | Description |
|-----------|------|---------|---------|-------------|
| `interval` | string | `"hour"` | `"hour"`, `"minute"` | Time-series aggregation interval |

### Examples

#### Using cURL (hour interval - 24 hourly data points)
```bash
curl -X GET "http://localhost:8080/api/v1/admin/dashboard/auth-trends?interval=hour" \
  -H "Authorization: Bearer eyJhbGciOiJFZERTQSIsInR5cCI6IkpXVCJ9..." \
  -H "Content-Type: application/json"
```

#### Using cURL (minute interval - 1440 minute data points)
```bash
curl -X GET "http://localhost:8080/api/v1/admin/dashboard/auth-trends?interval=minute" \
  -H "Authorization: Bearer <token>"
```

#### Using JavaScript/Fetch
```javascript
const token = "your-jwt-token-here";

fetch("/api/v1/admin/dashboard/auth-trends?interval=hour", {
  method: "GET",
  headers: {
    "Authorization": `Bearer ${token}`,
    "Content-Type": "application/json"
  }
})
.then(res => res.json())
.then(data => console.log(data))
.catch(err => console.error(err));
```

### Success Response (Status: 200 OK)

#### Hour Interval Response
```json
{
  "interval": "hour",
  "data": [
    {
      "timestamp": "2026-03-29T00:00:00Z",
      "successCount": 120,
      "failureCount": 5
    },
    {
      "timestamp": "2026-03-29T01:00:00Z",
      "successCount": 98,
      "failureCount": 8
    },
    {
      "timestamp": "2026-03-29T02:00:00Z",
      "successCount": 150,
      "failureCount": 3
    }
    // ... up to 24 hours
  ]
}
```

#### Minute Interval Response
```json
{
  "interval": "minute",
  "data": [
    {
      "timestamp": "2026-03-28T00:00:00Z",
      "successCount": 2,
      "failureCount": 0
    },
    {
      "timestamp": "2026-03-28T00:01:00Z",
      "successCount": 1,
      "failureCount": 0
    },
    // ... 1440 data points for 24 hours
  ]
}
```

### Error Responses

#### 400 Bad Request - Invalid Interval
```json
{
  "error": "BadRequest",
  "message": "Invalid interval. Supported values are 'minute' or 'hour'"
}
```

#### 401 Unauthorized - Missing Token
```json
{
  "error": "Unauthorized",
  "message": "missing authorization header"
}
```

#### 401 Unauthorized - Invalid/Expired Token
```json
{
  "error": "Unauthorized",
  "message": "invalid or expired token"
}
```

#### 401 Unauthorized - Session Revoked
```json
{
  "error": "Unauthorized",
  "message": "session expired or revoked"
}
```

#### 403 Forbidden - Not Admin
```json
{
  "error": "Forbidden",
  "message": "You do not have access to this resource"
}
```

#### 403 Forbidden - Role Mismatch
```json
{
  "error": "Forbidden",
  "message": "role mismatch between token and session - possible security violation"
}
```

#### 500 Internal Server Error
```json
{
  "error": "InternalServerError",
  "message": "Failed to fetch authentication trends"
}
```

### Response Fields Explanation

#### Response Root
- `interval` (string): The aggregation interval used ("minute" or "hour")
- `data` (array): Array of time-series data points, sorted chronologically

#### Data Point
- `timestamp` (string, ISO 8601): Start of the time bucket (UTC)
- `successCount` (integer): Number of successful authentications in this interval
- `failureCount` (integer): Number of failed authentications in this interval

### Security Notes

1. **Role-Based Access**: Only `Admin` role users can access this endpoint
2. **Cross-Validation**: JWT role is validated against Redis session cache role
3. **Session Freshness**: Session must be active and MFA-verified
4. **Token Revocation**: If session is revoked in Redis, request fails with 401

### Implementation Details

The endpoint:
- Queries the last 24 hours of authentication activity
- Aggregates Activity Logs by the specified interval
- Fills missing time buckets with zero counts (prevents visualization gaps)
- Supports minute-level or hour-level granularity
- Returns empty buckets (count=0) for periods with no activity

### Performance Notes

- **Hour interval**: ~24 data points returned
- **Minute interval**: ~1440 data points returned
- Uses efficient PostgreSQL DATE_TRUNC for aggregation
- Database queries are indexed on created_at and type columns
- Suitable for real-time dashboards and time-series visualizations
