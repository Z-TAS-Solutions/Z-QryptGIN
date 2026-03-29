#!/bin/bash

# Dashboard Authentication Trends - Implementation Verification Checklist
# Run this after building the project to verify everything is in place

echo "=========================================="
echo "Z-QryptGIN Dashboard Implementation Check"
echo "=========================================="
echo ""

PROJECT_ROOT="c:\Users\Cyberlowspecs\Documents\Coding\Z-TAS\Z-QryptGIN"

# Color codes
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Counter
PASSED=0
FAILED=0

check_file() {
    local filepath=$1
    local description=$2
    
    if [ -f "$filepath" ]; then
        echo -e "${GREEN}✓${NC} PASS: $description"
        ((PASSED++))
    else
        echo -e "${RED}✗${NC} FAIL: $description (File not found: $filepath)"
        ((FAILED++))
    fi
}

check_content() {
    local filepath=$1
    local search_term=$2
    local description=$3
    
    if grep -q "$search_term" "$filepath" 2>/dev/null; then
        echo -e "${GREEN}✓${NC} PASS: $description"
        ((PASSED++))
    else
        echo -e "${RED}✗${NC} FAIL: $description (Term not found in $filepath)"
        ((FAILED++))
    fi
}

echo "1. Checking if new files were created..."
echo "=========================================="
check_file "$PROJECT_ROOT/internal/app/dto/dashboard_dto.go" "Dashboard DTO file created"
check_file "$PROJECT_ROOT/internal/app/repository/dashboard_repo.go" "Dashboard Repository created"
check_file "$PROJECT_ROOT/internal/app/service/dashboard_service.go" "Dashboard Service created"
check_file "$PROJECT_ROOT/api/handlers/dashboard_handler.go" "Dashboard Handler created"
check_file "$PROJECT_ROOT/docs/DASHBOARD_AUTH_TRENDS_API.md" "Dashboard API documentation created"
check_file "$PROJECT_ROOT/docs/RBAC_ARCHITECTURE.md" "RBAC Architecture documentation created"
check_file "$PROJECT_ROOT/docs/IMPLEMENTATION_SUMMARY.md" "Implementation Summary created"
check_file "$PROJECT_ROOT/docs/ARCHITECTURE_DIAGRAMS.md" "Architecture Diagrams created"
echo ""

echo "2. Checking if new middleware was added..."
echo "=========================================="
check_content "$PROJECT_ROOT/api/server/auth.go" "func RequireRole" "RequireRole middleware implemented"
check_content "$PROJECT_ROOT/api/server/auth.go" "func ValidateRoleConsistency" "ValidateRoleConsistency middleware implemented"
echo ""

echo "3. Checking if Session DTO was updated with Role field..."
echo "=========================================="
check_content "$PROJECT_ROOT/internal/app/dto/session.go" "Role.*string" "Session DTO includes Role field"
echo ""

echo "4. Checking if main.go was updated..."
echo "=========================================="
check_content "$PROJECT_ROOT/cmd/api/main.go" "dashboardRepo := repository.NewDashboardRepository" "DashboardRepository initialization added"
check_content "$PROJECT_ROOT/cmd/api/main.go" "dashboardSvc := service.NewDashboardService" "DashboardService initialization added"
check_content "$PROJECT_ROOT/cmd/api/main.go" "dashboardHandler := handlers.NewDashboardHandler" "DashboardHandler initialization added"
check_content "$PROJECT_ROOT/cmd/api/main.go" "dashboard.GET.*auth-trends.*dashboardHandler.GetAuthenticationTrends" "Dashboard routes added"
check_content "$PROJECT_ROOT/cmd/api/main.go" "RequireRole.*Admin" "Role-based middleware added to routes"
echo ""

echo "5. Checking if DTOs have required fields..."
echo "=========================================="
check_content "$PROJECT_ROOT/internal/app/dto/dashboard_dto.go" "AuthTrendDataPoint struct" "AuthTrendDataPoint DTO defined"
check_content "$PROJECT_ROOT/internal/app/dto/dashboard_dto.go" "GetAuthTrendsResponse struct" "GetAuthTrendsResponse DTO defined"
check_content "$PROJECT_ROOT/internal/app/dto/dashboard_dto.go" "timestamp.*json" "AuthTrendDataPoint has timestamp field"
check_content "$PROJECT_ROOT/internal/app/dto/dashboard_dto.go" "successCount" "AuthTrendDataPoint has successCount field"
check_content "$PROJECT_ROOT/internal/app/dto/dashboard_dto.go" "failureCount" "AuthTrendDataPoint has failureCount field"
echo ""

echo "6. Checking Repository interface implementation..."
echo "=========================================="
check_content "$PROJECT_ROOT/internal/app/repository/dashboard_repo.go" "GetAuthTrendsByInterval" "GetAuthTrendsByInterval method implemented"
check_content "$PROJECT_ROOT/internal/app/repository/dashboard_repo.go" "GetAuthTrendsMetrics" "GetAuthTrendsMetrics method implemented"
check_content "$PROJECT_ROOT/internal/app/repository/dashboard_repo.go" "fillMissingIntervals" "fillMissingIntervals helper implemented"
check_content "$PROJECT_ROOT/internal/app/repository/dashboard_repo.go" "DATE_TRUNC" "PostgreSQL DATE_TRUNC used for aggregation"
echo ""

echo "7. Checking Service implementation..."
echo "=========================================="
check_content "$PROJECT_ROOT/internal/app/service/dashboard_service.go" "DashboardService interface" "DashboardService interface defined"
check_content "$PROJECT_ROOT/internal/app/service/dashboard_service.go" "func.*GetAuthenticationTrends" "GetAuthenticationTrends method implemented"
check_content "$PROJECT_ROOT/internal/app/service/dashboard_service.go" "Supported values are 'minute' or 'hour'" "Interval validation in service"
echo ""

echo "8. Checking Handler implementation..."
echo "=========================================="
check_content "$PROJECT_ROOT/api/handlers/dashboard_handler.go" "func.*GetAuthenticationTrends.*c \*gin.Context" "GetAuthenticationTrends HTTP handler implemented"
check_content "$PROJECT_ROOT/api/handlers/dashboard_handler.go" "DefaultQuery.*interval.*hour" "Interval parameter with default value"
check_content "$PROJECT_ROOT/api/handlers/dashboard_handler.go" "StatusBadRequest" "Returns 400 for invalid interval"
check_content "$PROJECT_ROOT/api/handlers/dashboard_handler.go" "StatusInternalServerError" "Returns 500 for server errors"
echo ""

echo "9. Checking RBAC implementation..."
echo "=========================================="
check_content "$PROJECT_ROOT/api/server/auth.go" "func RequireRole.*allowedRoles.*string.*gin.HandlerFunc" "RequireRole signature correct"
check_content "$PROJECT_ROOT/api/server/auth.go" "role == allowedRole" "Role comparison logic in RequireRole"
check_content "$PROJECT_ROOT/api/server/auth.go" "jwt_\.role ==.*session_\.role" "Role consistency check in ValidateRoleConsistency"
check_content "$PROJECT_ROOT/api/server/auth.go" "403.*Forbidden" "403 Forbidden returned on role failure"
echo ""

echo "10. Checking documentation completeness..."
echo "=========================================="
check_content "$PROJECT_ROOT/docs/DASHBOARD_AUTH_TRENDS_API.md" "GET /api/v1/admin/dashboard/auth-trends" "API documentation includes endpoint"
check_content "$PROJECT_ROOT/docs/DASHBOARD_AUTH_TRENDS_API.md" "interval.*query parameter" "API docs include query parameters"
check_content "$PROJECT_ROOT/docs/RBAC_ARCHITECTURE.md" "Role-Based Access Control" "RBAC documentation present"
check_content "$PROJECT_ROOT/docs/RBAC_ARCHITECTURE.md" "RequireRole" "RBAC docs explain RequireRole"
check_content "$PROJECT_ROOT/docs/RBAC_ARCHITECTURE.md" "ValidateRoleConsistency" "RBAC docs explain ValidateRoleConsistency"
echo ""

echo "=========================================="
echo "VERIFICATION SUMMARY"
echo "=========================================="
echo -e "Passed: ${GREEN}$PASSED${NC}"
echo -e "Failed: ${RED}$FAILED${NC}"

if [ $FAILED -eq 0 ]; then
    echo -e "\n${GREEN}✓ All checks passed! Implementation is complete.${NC}"
    exit 0
else
    echo -e "\n${YELLOW}⚠ Some checks failed. Please review the failures above.${NC}"
    exit 1
fi
