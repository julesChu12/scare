#!/bin/bash
# sCare 多身份架构迁移 API 测试脚本
# 使用方法: ./run_tests.sh [base_url]
# 测试完成后会在当前目录生成 test_responses_YYYYMMDD_HHMMSS.json

BASE_URL="${1:-http://localhost:8080/api/v1}"
PASSWORD="Test@123"

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 计数器
PASSED=0
FAILED=0

# 响应记录文件
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")
RESPONSE_FILE="${SCRIPT_DIR}/test_responses_${TIMESTAMP}.json"

# 初始化响应 JSON
echo "{" > "$RESPONSE_FILE"
echo '  "test_info": {' >> "$RESPONSE_FILE"
echo "    \"timestamp\": \"$(date -Iseconds)\"," >> "$RESPONSE_FILE"
echo "    \"base_url\": \"$BASE_URL\"" >> "$RESPONSE_FILE"
echo '  },' >> "$RESPONSE_FILE"
echo '  "responses": {' >> "$RESPONSE_FILE"
FIRST_RESPONSE=true

# 记录响应到文件
record_response() {
    local test_id="$1"
    local http_code="$2"
    local body="$3"

    if [ "$FIRST_RESPONSE" = true ]; then
        FIRST_RESPONSE=false
    else
        echo "," >> "$RESPONSE_FILE"
    fi

    # 转义 JSON 字符串中的特殊字符
    local escaped_body=$(echo "$body" | python3 -c "import sys,json; print(json.dumps(sys.stdin.read().strip()))")

    echo "    \"$test_id\": {" >> "$RESPONSE_FILE"
    echo "      \"http_code\": $http_code," >> "$RESPONSE_FILE"
    echo "      \"body\": $escaped_body" >> "$RESPONSE_FILE"
    echo -n "    }" >> "$RESPONSE_FILE"
}

# 测试函数
test_api() {
    local test_id="$1"
    local name="$2"
    local method="$3"
    local endpoint="$4"
    local expected_status="$5"
    local expected_msg="$6"
    local data="$7"
    local token="$8"

    local headers="-H 'Content-Type: application/json'"
    if [ -n "$token" ]; then
        headers="$headers -H 'Authorization: Bearer $token'"
    fi

    if [ -n "$data" ]; then
        response=$(eval "curl -s -w '\n%{http_code}' -X $method '$BASE_URL$endpoint' $headers -d '$data'")
    else
        response=$(eval "curl -s -w '\n%{http_code}' -X $method '$BASE_URL$endpoint' $headers")
    fi

    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | sed '$d')
    msg=$(echo "$body" | python3 -c "import sys,json; print(json.load(sys.stdin).get('msg',''))" 2>/dev/null || echo "")

    # 记录响应到文件
    record_response "$test_id" "$http_code" "$body"

    if [ "$http_code" = "$expected_status" ] && [ "$msg" = "$expected_msg" ]; then
        echo -e "${GREEN}✓${NC} $name"
        ((PASSED++))
        echo "$body"
    else
        echo -e "${RED}✗${NC} $name"
        echo "  Expected: status=$expected_status, msg=$expected_msg"
        echo "  Got: status=$http_code, msg=$msg"
        ((FAILED++))
    fi
    echo ""
}

echo "========================================"
echo "sCare 多身份架构迁移 API 测试"
echo "Base URL: $BASE_URL"
echo "========================================"
echo ""

# 1. 健康检查
echo "--- 基础测试 ---"
test_api "health_check" "健康检查" "GET" "/health" "200" "ok"

# 2. B端登录测试
echo "--- B端登录测试 ---"

# Admin 登录
ADMIN_RESP=$(curl -s -X POST "$BASE_URL/b/auth/login" -H "Content-Type: application/json" -d '{"phone":"13800000001","password":"Test@123"}')
ADMIN_HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$BASE_URL/b/auth/login" -H "Content-Type: application/json" -d '{"phone":"13800000001","password":"Test@123"}')
ADMIN_TOKEN=$(echo "$ADMIN_RESP" | python3 -c "import sys,json; print(json.load(sys.stdin)['data']['token'])" 2>/dev/null)
ADMIN_IDENTITIES=$(echo "$ADMIN_RESP" | python3 -c "import sys,json; print(json.load(sys.stdin)['data']['identities'])" 2>/dev/null)
record_response "b_auth_admin_login" "$ADMIN_HTTP_CODE" "$ADMIN_RESP"
if [ -n "$ADMIN_TOKEN" ]; then
    echo -e "${GREEN}✓${NC} Admin 登录成功"
    echo "  identities: $ADMIN_IDENTITIES"
    ((PASSED++))
else
    echo -e "${RED}✗${NC} Admin 登录失败"
    ((FAILED++))
fi
echo ""

# Station Manager 登录
SM_RESP=$(curl -s -X POST "$BASE_URL/b/auth/login" -H "Content-Type: application/json" -d '{"phone":"13800000002","password":"Test@123"}')
SM_HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$BASE_URL/b/auth/login" -H "Content-Type: application/json" -d '{"phone":"13800000002","password":"Test@123"}')
SM_TOKEN=$(echo "$SM_RESP" | python3 -c "import sys,json; print(json.load(sys.stdin)['data']['token'])" 2>/dev/null)
SM_IDENTITIES=$(echo "$SM_RESP" | python3 -c "import sys,json; print(json.load(sys.stdin)['data']['identities'])" 2>/dev/null)
record_response "b_auth_station_manager_login" "$SM_HTTP_CODE" "$SM_RESP"
if [ -n "$SM_TOKEN" ]; then
    echo -e "${GREEN}✓${NC} Station Manager 登录成功"
    echo "  identities: $SM_IDENTITIES"
    ((PASSED++))
else
    echo -e "${RED}✗${NC} Station Manager 登录失败"
    ((FAILED++))
fi
echo ""

# Staff 登录
STAFF_RESP=$(curl -s -X POST "$BASE_URL/b/auth/login" -H "Content-Type: application/json" -d '{"phone":"13800000004","password":"Test@123"}')
STAFF_HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$BASE_URL/b/auth/login" -H "Content-Type: application/json" -d '{"phone":"13800000004","password":"Test@123"}')
STAFF_TOKEN=$(echo "$STAFF_RESP" | python3 -c "import sys,json; print(json.load(sys.stdin)['data']['token'])" 2>/dev/null)
STAFF_IDENTITIES=$(echo "$STAFF_RESP" | python3 -c "import sys,json; print(json.load(sys.stdin)['data']['identities'])" 2>/dev/null)
record_response "b_auth_staff_login" "$STAFF_HTTP_CODE" "$STAFF_RESP"
if [ -n "$STAFF_TOKEN" ]; then
    echo -e "${GREEN}✓${NC} Staff 登录成功"
    echo "  identities: $STAFF_IDENTITIES"
    ((PASSED++))
else
    echo -e "${RED}✗${NC} Staff 登录失败"
    ((FAILED++))
fi
echo ""

# 3. C端登录测试
echo "--- C端登录测试 ---"

# Elderly 密码登录
ELDERLY_RESP=$(curl -s -X POST "$BASE_URL/c/auth/login" -H "Content-Type: application/json" -d '{"phone":"13800000008","password":"Test@123","type":"password"}')
ELDERLY_HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$BASE_URL/c/auth/login" -H "Content-Type: application/json" -d '{"phone":"13800000008","password":"Test@123","type":"password"}')
ELDERLY_TOKEN=$(echo "$ELDERLY_RESP" | python3 -c "import sys,json; print(json.load(sys.stdin)['data']['token'])" 2>/dev/null)
ELDERLY_TYPE=$(echo "$ELDERLY_RESP" | python3 -c "import sys,json; print(json.load(sys.stdin)['data']['customer_type'])" 2>/dev/null)
record_response "c_auth_elderly_login" "$ELDERLY_HTTP_CODE" "$ELDERLY_RESP"
if [ -n "$ELDERLY_TOKEN" ]; then
    echo -e "${GREEN}✓${NC} C端 elderly 登录成功"
    echo "  customer_type: $ELDERLY_TYPE"
    ((PASSED++))
else
    echo -e "${RED}✗${NC} C端 elderly 登录失败"
    ((FAILED++))
fi
echo ""

# 4. 多身份用户测试
echo "--- 多身份用户测试 (李奶奶 - elderly + staff) ---"

# B端登录
MULTI_B_RESP=$(curl -s -X POST "$BASE_URL/b/auth/login" -H "Content-Type: application/json" -d '{"phone":"13800000009","password":"Test@123"}')
MULTI_B_HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$BASE_URL/b/auth/login" -H "Content-Type: application/json" -d '{"phone":"13800000009","password":"Test@123"}')
MULTI_B_IDENTITIES=$(echo "$MULTI_B_RESP" | python3 -c "import sys,json; print(json.load(sys.stdin)['data']['identities'])" 2>/dev/null)
MULTI_B_PRIMARY=$(echo "$MULTI_B_RESP" | python3 -c "import sys,json; print(json.load(sys.stdin)['data']['primary'])" 2>/dev/null)
record_response "multi_identity_b_login" "$MULTI_B_HTTP_CODE" "$MULTI_B_RESP"
if [ "$MULTI_B_IDENTITIES" = "['staff']" ]; then
    echo -e "${GREEN}✓${NC} 多身份用户 B端登录成功"
    echo "  identities: $MULTI_B_IDENTITIES (B端身份)"
    echo "  primary: $MULTI_B_PRIMARY (全局主身份)"
    ((PASSED++))
else
    echo -e "${RED}✗${NC} 多身份用户 B端登录失败"
    echo "  Got identities: $MULTI_B_IDENTITIES"
    ((FAILED++))
fi
echo ""

# C端登录
MULTI_C_RESP=$(curl -s -X POST "$BASE_URL/c/auth/login" -H "Content-Type: application/json" -d '{"phone":"13800000009","password":"Test@123","type":"password"}')
MULTI_C_HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$BASE_URL/c/auth/login" -H "Content-Type: application/json" -d '{"phone":"13800000009","password":"Test@123","type":"password"}')
MULTI_C_TYPE=$(echo "$MULTI_C_RESP" | python3 -c "import sys,json; print(json.load(sys.stdin)['data']['customer_type'])" 2>/dev/null)
record_response "multi_identity_c_login" "$MULTI_C_HTTP_CODE" "$MULTI_C_RESP"
if [ "$MULTI_C_TYPE" = "elderly" ]; then
    echo -e "${GREEN}✓${NC} 多身份用户 C端登录成功"
    echo "  customer_type: $MULTI_C_TYPE"
    ((PASSED++))
else
    echo -e "${RED}✗${NC} 多身份用户 C端登录失败"
    ((FAILED++))
fi
echo ""

# 5. 权限测试
echo "--- 权限验证测试 ---"

# Admin 访问用户列表
ADMIN_ACCESS=$(curl -s "$BASE_URL/b/users?page=1&page_size=5" -H "Authorization: Bearer $ADMIN_TOKEN")
ADMIN_ACCESS_CODE=$(curl -s -o /dev/null -w "%{http_code}" "$BASE_URL/b/users?page=1&page_size=5" -H "Authorization: Bearer $ADMIN_TOKEN")
ADMIN_ACCESS_MSG=$(echo "$ADMIN_ACCESS" | python3 -c "import sys,json; print(json.load(sys.stdin)['msg'])" 2>/dev/null)
record_response "permission_admin_access" "$ADMIN_ACCESS_CODE" "$ADMIN_ACCESS"
if [ "$ADMIN_ACCESS_MSG" = "ok" ]; then
    echo -e "${GREEN}✓${NC} Admin 访问用户列表 - 权限通过"
    ((PASSED++))
else
    echo -e "${RED}✗${NC} Admin 访问用户列表 - 应该成功但失败"
    ((FAILED++))
fi
echo ""

# Staff 访问用户列表 (应被拒绝)
STAFF_ACCESS=$(curl -s "$BASE_URL/b/users?page=1&page_size=5" -H "Authorization: Bearer $STAFF_TOKEN")
STAFF_ACCESS_CODE=$(curl -s -o /dev/null -w "%{http_code}" "$BASE_URL/b/users?page=1&page_size=5" -H "Authorization: Bearer $STAFF_TOKEN")
STAFF_ACCESS_MSG=$(echo "$STAFF_ACCESS" | python3 -c "import sys,json; print(json.load(sys.stdin)['msg'])" 2>/dev/null)
record_response "permission_staff_forbidden" "$STAFF_ACCESS_CODE" "$STAFF_ACCESS"
if [ "$STAFF_ACCESS_MSG" = "forbidden" ]; then
    echo -e "${GREEN}✓${NC} Staff 访问用户列表 - RBAC 正确拒绝"
    ((PASSED++))
else
    echo -e "${RED}✗${NC} Staff 访问用户列表 - 应该被拒绝但通过了"
    ((FAILED++))
fi
echo ""

# C端用户访问 C端路由 (应成功 - 跳过 RBAC)
C_ACCESS=$(curl -s "$BASE_URL/c/requests" -H "Authorization: Bearer $ELDERLY_TOKEN")
C_ACCESS_CODE=$(curl -s -o /dev/null -w "%{http_code}" "$BASE_URL/c/requests" -H "Authorization: Bearer $ELDERLY_TOKEN")
C_ACCESS_MSG=$(echo "$C_ACCESS" | python3 -c "import sys,json; print(json.load(sys.stdin)['msg'])" 2>/dev/null)
record_response "permission_c_end_skip_rbac" "$C_ACCESS_CODE" "$C_ACCESS"
if [ "$C_ACCESS_MSG" = "ok" ]; then
    echo -e "${GREEN}✓${NC} C端用户访问服务请求 - 跳过 RBAC 成功"
    ((PASSED++))
else
    echo -e "${RED}✗${NC} C端用户访问服务请求 - 应该成功但失败"
    ((FAILED++))
fi
echo ""

# C端用户访问 B端路由 (应被拒绝 - 端类型不匹配)
C_B_ACCESS=$(curl -s "$BASE_URL/b/users" -H "Authorization: Bearer $ELDERLY_TOKEN")
C_B_ACCESS_CODE=$(curl -s -o /dev/null -w "%{http_code}" "$BASE_URL/b/users" -H "Authorization: Bearer $ELDERLY_TOKEN")
C_B_ACCESS_MSG=$(echo "$C_B_ACCESS" | python3 -c "import sys,json; print(json.load(sys.stdin)['msg'])" 2>/dev/null)
record_response "permission_c_end_type_mismatch" "$C_B_ACCESS_CODE" "$C_B_ACCESS"
if [ "$C_B_ACCESS_MSG" = "token end type mismatch" ]; then
    echo -e "${GREEN}✓${NC} C端访问B端路由 - 端类型检查正确拒绝"
    ((PASSED++))
else
    echo -e "${RED}✗${NC} C端访问B端路由 - 应该被拒绝"
    echo "  Got: $C_B_ACCESS_MSG"
    ((FAILED++))
fi
echo ""

# 6. 错误场景测试
echo "--- 错误场景测试 ---"

# 无B端身份用户尝试B端登录
NO_B_RESP=$(curl -s -X POST "$BASE_URL/b/auth/login" -H "Content-Type: application/json" -d '{"phone":"13800000008","password":"Test@123"}')
NO_B_CODE=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$BASE_URL/b/auth/login" -H "Content-Type: application/json" -d '{"phone":"13800000008","password":"Test@123"}')
NO_B_MSG=$(echo "$NO_B_RESP" | python3 -c "import sys,json; print(json.load(sys.stdin)['msg'])" 2>/dev/null)
record_response "b_auth_no_identity" "$NO_B_CODE" "$NO_B_RESP"
if [ "$NO_B_MSG" = "user has no B-end identity" ]; then
    echo -e "${GREEN}✓${NC} 无B端身份用户登录 - 正确拒绝"
    ((PASSED++))
else
    echo -e "${RED}✗${NC} 无B端身份用户登录 - 应该被拒绝"
    echo "  Got: $NO_B_MSG"
    ((FAILED++))
fi
echo ""

# 密码错误
WRONG_PWD=$(curl -s -X POST "$BASE_URL/b/auth/login" -H "Content-Type: application/json" -d '{"phone":"13800000001","password":"WrongPassword"}')
WRONG_PWD_CODE=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$BASE_URL/b/auth/login" -H "Content-Type: application/json" -d '{"phone":"13800000001","password":"WrongPassword"}')
WRONG_PWD_MSG=$(echo "$WRONG_PWD" | python3 -c "import sys,json; print(json.load(sys.stdin)['msg'])" 2>/dev/null)
record_response "b_auth_invalid_password" "$WRONG_PWD_CODE" "$WRONG_PWD"
if [ "$WRONG_PWD_MSG" = "invalid credentials" ]; then
    echo -e "${GREEN}✓${NC} 密码错误 - 正确拒绝"
    ((PASSED++))
else
    echo -e "${RED}✗${NC} 密码错误 - 应该返回 invalid credentials"
    echo "  Got: $WRONG_PWD_MSG"
    ((FAILED++))
fi
echo ""

# 关闭响应 JSON 文件
echo "" >> "$RESPONSE_FILE"
echo "  }," >> "$RESPONSE_FILE"
echo "  \"summary\": {" >> "$RESPONSE_FILE"
echo "    \"passed\": $PASSED," >> "$RESPONSE_FILE"
echo "    \"failed\": $FAILED," >> "$RESPONSE_FILE"
echo "    \"total\": $((PASSED + FAILED))" >> "$RESPONSE_FILE"
echo "  }" >> "$RESPONSE_FILE"
echo "}" >> "$RESPONSE_FILE"

# 总结
echo "========================================"
echo "测试结果汇总"
echo "========================================"
echo -e "通过: ${GREEN}$PASSED${NC}"
echo -e "失败: ${RED}$FAILED${NC}"
echo ""
echo -e "${YELLOW}响应记录已保存到: ${NC}$RESPONSE_FILE"
echo ""

if [ $FAILED -eq 0 ]; then
    echo -e "${GREEN}所有测试通过!${NC}"
    exit 0
else
    echo -e "${RED}有 $FAILED 个测试失败${NC}"
    exit 1
fi
