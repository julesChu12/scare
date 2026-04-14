#!/bin/bash
# =====================================================
# sCare C端业务流程测试脚本
# 覆盖：快速开通（未登录注册）、登录后提交服务请求、完整用户流程
# 用法: ./scripts/test_business_flow.sh [base_url]
# 默认: http://localhost:8080
# =====================================================

BASE_URL="${1:-http://localhost:8080}"
API_URL="$BASE_URL/api/v1"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

TOTAL=0; PASSED=0; FAILED=0

# ---------- 工具函数 ----------
info()  { echo -e "${BLUE}[INFO]${NC}  $1"; }
ok()    { echo -e "${GREEN}[OK]${NC}   $1"; }
fail()  { echo -e "${RED}[FAIL]${NC} $1"; }
warn()  { echo -e "${YELLOW}[WARN]${NC} $1"; }

# GET 请求，返回 "body\nhttp_code"
get() {
    local endpoint="$1"; local token="${2:-}"
    local h="Content-Type: application/json"
    [ -n "$token" ] && h="$h"$'\n'"Authorization: Bearer $token"
    curl -s -X GET "$API_URL$endpoint" -H "$h" -w '\nHTTP_CODE:%{http_code}'
}

# POST 请求，返回 "body\nhttp_code"
post() {
    local endpoint="$1"; local data="$2"; local token="${3:-}"
    local h="Content-Type: application/json"
    [ -n "$token" ] && h="$h"$'\n'"Authorization: Bearer $token"
    curl -s -X POST "$API_URL$endpoint" -H "$h" -d "$data" -w '\nHTTP_CODE:%{http_code}'
}

# PUT 请求，返回 "body\nhttp_code"
put() {
    local endpoint="$1"; local data="$2"; local token="${3:-}"
    local h="Content-Type: application/json"
    [ -n "$token" ] && h="$h"$'\n'"Authorization: Bearer $token"
    curl -s -X PUT "$API_URL$endpoint" -H "$h" -d "$data" -w '\nHTTP_CODE:%{http_code}'
}

# 验证响应 msg
assert_msg() {
    local name="$1"; local response="$2"; local expected_msg="$3"
    TOTAL=$((TOTAL + 1))
    local body=$(echo "$response" | sed '/HTTP_CODE:/d')
    local http_code=$(echo "$response" | grep "HTTP_CODE:" | sed 's/HTTP_CODE://')
    local msg=$(echo "$body" | grep -o "\"msg\":\s*\"[^\"]*\"" | head -1 | sed 's/"msg":\s*"\([^"]*\)"/\1/')
    if echo "$body" | grep -q "\"msg\":\s*\"$expected_msg\""; then
        ok "$name"
        PASSED=$((PASSED + 1)); return 0
    else
        fail "$name (expected msg: $expected_msg, got: $msg, code: $http_code)"
        FAILED=$((FAILED + 1)); return 1
    fi
}

# 验证响应 HTTP code
assert_code() {
    local name="$1"; local response="$2"; local expected="$3"
    TOTAL=$((TOTAL + 1))
    local http_code=$(echo "$response" | grep "HTTP_CODE:" | sed 's/HTTP_CODE://')
    if [ "$http_code" = "$expected" ]; then
        ok "$name (HTTP $http_code)"
        PASSED=$((PASSED + 1)); return 0
    else
        fail "$name (expected HTTP $expected, got $http_code)"
        FAILED=$((FAILED + 1)); return 1
    fi
}

# 验证字段非空
assert_not_empty() {
    local name="$1"; local value="$2"
    TOTAL=$((TOTAL + 1))
    if [ -n "$value" ] && [ "$value" != "null" ] && [ "$value" != "empty" ]; then
        ok "$name (value=$value)"
        PASSED=$((PASSED + 1)); return 0
    else
        fail "$name (value is empty)"
        FAILED=$((FAILED + 1)); return 1
    fi
}

# =====================================================
# 流程 1：未登录用户 - 快速开通（注册+登录+提交服务）
# =====================================================
test_quick_start_flow() {
    echo ""
    echo "=========================================="
    echo -e "${BLUE}流程 1：快速开通（未登录注册+提交服务）${NC}"
    echo "=========================================="

    local TS=$(date +%s) && local NEW_PHONE="139${TS: -8}"
    info "使用新手机号: $NEW_PHONE"

    # 1.1 发送验证码
    info "1.1 发送验证码..."
    local SEND_RESP=$(post "/c/auth/send-code" "{\"phone\":\"$NEW_PHONE\"}")
    assert_msg "发送验证码" "$SEND_RESP" "验证码已发送" || true

    # 1.2 快速开通（开发环境万能码 000000）
    info "1.2 快速开通（注册+登录+提交服务）..."
    local QS_RESP=$(post "/c/auth/quick-start" \
        "{
            \"phone\": \"$NEW_PHONE\",
            \"code\": \"000000\",
            \"name\": \"快速测试用户\",
            \"address\": \"北京市昌平区霍营街道华龙苑北里小区\",
            \"latitude\": 40.07,
            \"longitude\": 116.37,
            \"service_type\": \"meal\",
            \"description\": \"午餐送餐服务测试\",
            \"contact_name\": \"快速测试用户\",
            \"contact_phone\": \"$NEW_PHONE\"
        }")

    local QS_TOKEN=$(echo "$QS_RESP" | grep -o '"token":\s*"[^"]*"' | sed 's/"token":\s*"\([^"]*\)"/\1/')
    local QS_REQ_ID=$(echo "$QS_RESP" | grep -o '"id":\s*[0-9]*' | tail -1 | grep -o '[0-9]*')

    if echo "$QS_RESP" | grep -q "\"token\":"; then
        ok "快速开通成功 - token: ${QS_TOKEN:0:20}..."
        PASSED=$((PASSED + 1)); TOTAL=$((TOTAL + 1))
    else
        local QS_MSG=$(echo "$QS_RESP" | grep -o "\"msg\":\s*\"[^\"]*\"" | head -1 | sed 's/"msg":\s*"\([^"]*\)"/\1/')
        fail "快速开通失败 (msg=$QS_MSG)"
        FAILED=$((FAILED + 1)); TOTAL=$((TOTAL + 1)); return 1
    fi

    # 1.3 验证登录状态
    info "1.3 验证登录状态..."
    local ME_RESP=$(get "/c/auth/me" "$QS_TOKEN")
    assert_msg "获取当前用户信息" "$ME_RESP" "ok"

    # 1.4 查看提交的服务请求
    if [ -n "$QS_REQ_ID" ] && [ "$QS_REQ_ID" != "null" ] && [ -n "${QS_REQ_ID//[0-9]/}" ]; then
        info "1.4 查看提交的服务请求 (ID=$QS_REQ_ID)..."
        local REQ_RESP=$(get "/c/requests/$QS_REQ_ID" "$QS_TOKEN")
        local REQ_STATUS=$(echo "$REQ_RESP" | grep -o '"status":\s*"[^"]*"' | sed 's/"status":\s*"\([^"]*\)"/\1/')
        assert_not_empty "请求状态" "$REQ_STATUS"
    else
        warn "未返回有效请求ID，跳过请求详情检查"
    fi

    # 1.5 查看通知
    info "1.5 查看通知列表..."
    local NOTIF_RESP=$(get "/c/notifications" "$QS_TOKEN")
    assert_msg "通知列表" "$NOTIF_RESP" "ok"

    # 1.6 登出
    info "1.6 登出..."
    local LOGOUT_RESP=$(post "/c/auth/logout" "" "$QS_TOKEN")
    assert_msg "登出" "$LOGOUT_RESP" "logout successful"

    # 1.7 登出后验证 token 失效
    info "1.7 验证 token 已失效..."
    local CHECK_RESP=$(get "/c/auth/check" "$QS_TOKEN")
    assert_code "Token 已失效" "$CHECK_RESP" "401"
}

# =====================================================
# 流程 2：已注册用户 - 登录后提交服务请求
# =====================================================
test_login_submit_flow() {
    echo ""
    echo "=========================================="
    echo -e "${BLUE}流程 2：已注册用户登录 + 提交服务请求${NC}"
    echo "=========================================="

    local C_PHONE="13800000008"
    local C_PASSWORD="Test@123"
    info "使用测试账号: $C_PHONE"

    # 2.1 登录
    info "2.1 密码登录..."
    local LOGIN_RESP=$(post "/c/auth/login" \
        "{\"phone\":\"$C_PHONE\",\"password\":\"$C_PASSWORD\",\"type\":\"password\"}")
    local LOGIN_MSG=$(echo "$LOGIN_RESP" | grep -o "\"msg\":\s*\"[^\"]*\"" | head -1 | sed 's/"msg":\s*"\([^"]*\)"/\1/')
    local C_TOKEN=$(echo "$LOGIN_RESP" | grep -o '"token":\s*"[^"]*"' | sed 's/"token":\s*"\([^"]*\)"/\1/')

    if echo "$LOGIN_RESP" | grep -q "\"token\":"; then
        ok "登录成功 - token: ${C_TOKEN:0:20}..."
        PASSED=$((PASSED + 1)); TOTAL=$((TOTAL + 1))
    else
        fail "登录失败 (msg=$LOGIN_MSG)"
        FAILED=$((FAILED + 1)); TOTAL=$((TOTAL + 1)); return 1
    fi

    # 2.2 验证登录状态
    info "2.2 验证登录状态..."
    local ME_RESP=$(get "/c/auth/me" "$C_TOKEN")
    assert_msg "获取当前用户" "$ME_RESP" "ok"

    # 2.3 访问公开接口
    info "2.3 访问公开接口（新闻列表）..."
    local NEWS_RESP=$(get "/c/news")
    assert_msg "新闻列表（公开）" "$NEWS_RESP" "ok"

    # 2.4 提交服务请求 - 送餐
    info "2.4 提交服务请求（送餐）..."
    local REQ_RESP=$(post "/c/requests" \
        "{
            \"service_type\": \"meal\",
            \"description\": \"午餐送餐，老年餐，一份\",
            \"contact_name\": \"张大爷\",
            \"contact_phone\": \"$C_PHONE\",
            \"address\": \"北京市昌平区霍营街道华龙苑北里小区1号楼3单元501\",
            \"submit_location_lat\": 40.07,
            \"submit_location_lng\": 116.37,
            \"urgency\": \"normal\"
        }" "$C_TOKEN")
    local REQ_MSG=$(echo "$REQ_RESP" | grep -o "\"msg\":\s*\"[^\"]*\"" | head -1 | sed 's/"msg":\s*"\([^"]*\)"/\1/')
    local REQ_ID=$(echo "$REQ_RESP" | grep -o '"id":\s*[0-9]*' | head -1 | grep -o '[0-9]*')

    if echo "$REQ_RESP" | grep -q '"id":'; then
        ok "提交服务请求成功 (ID=$REQ_ID, msg=$REQ_MSG)"
        PASSED=$((PASSED + 1)); TOTAL=$((TOTAL + 1))
    else
        fail "提交服务请求失败 (msg=$REQ_MSG)"
        FAILED=$((FAILED + 1)); TOTAL=$((TOTAL + 1))
    fi

    # 2.5 提交服务请求 - 清洁
    info "2.5 提交服务请求（清洁）..."
    local REQ_RESP2=$(post "/c/requests" \
        "{
            \"service_type\": \"cleaning\",
            \"description\": \"家庭清洁服务，两室一厅\",
            \"contact_name\": \"张大爷\",
            \"contact_phone\": \"$C_PHONE\",
            \"address\": \"北京市昌平区霍营街道华龙苑北里小区1号楼3单元501\",
            \"submit_location_lat\": 40.07,
            \"submit_location_lng\": 116.37,
            \"urgency\": \"normal\"
        }" "$C_TOKEN")
    assert_msg "提交清洁请求" "$REQ_RESP2" "ok"

    # 2.6 查看我的请求列表
    info "2.6 查看我的服务请求列表..."
    local LIST_RESP=$(get "/c/requests" "$C_TOKEN")
    assert_msg "服务请求列表" "$LIST_RESP" "ok"
    local LIST_TOTAL=$(echo "$LIST_RESP" | grep -o '"total":\s*[0-9]*' | grep -o '[0-9]*')
    assert_not_empty "请求总数" "$LIST_TOTAL"

    # 2.7 查看请求详情
    if [ -n "$REQ_ID" ] && [ -n "${REQ_ID//[0-9]/}" ]; then
        info "2.7 查看请求详情 (ID=$REQ_ID)..."
        local DETAIL_RESP=$(get "/c/requests/$REQ_ID" "$C_TOKEN")
        assert_msg "请求详情" "$DETAIL_RESP" "ok"
    fi

    # 2.8 B端任务处理 + 用户评价
    info "2.8 B端处理任务..."
    local ADMIN_RESP=$(post "/b/auth/login" '{"phone":"13800000001","password":"Test@123"}')
    local ADMIN_TOKEN=$(echo "$ADMIN_RESP" | grep -o '"token":\s*"[^"]*"' | sed 's/"token":\s*"\([^"]*\)"/\1/')
    local STAFF_RESP=$(post "/b/auth/login" '{"phone":"13800000004","password":"Test@123"}')
    local STAFF_TOKEN=$(echo "$STAFF_RESP" | grep -o '"token":\s*"[^"]*"' | sed 's/"token":\s*"\([^"]*\)"/\1/')

    if [ -n "$ADMIN_TOKEN" ] && [ -n "$STAFF_TOKEN" ] && [ -n "$REQ_ID" ] && [ -z "${REQ_ID//[0-9]/}" ]; then
        local TASKS_RESP=$(get "/b/tasks?page_size=100" "$ADMIN_TOKEN")
        # jq: 精确匹配 request_id -> task id
        local TASK_ID=$(echo "$TASKS_RESP" | jq --argjson rid "$REQ_ID" '.data.items[] | select(.request_id == $rid) | .id' 2>/dev/null)
        # fallback: 精确匹配 "request_id":N 的前一个字段
        if [ -z "$TASK_ID" ]; then
            TASK_ID=$(echo "$TASKS_RESP" | grep -oP "\"request_id\":$REQ_ID," | head -1 | grep -oP '.{0,20}"id":\K[0-9]+')
        fi
        if [ -n "$TASK_ID" ]; then
            info "  认领任务 (TaskID=$TASK_ID)..."
            local CLAIM_RESP=$(post "/b/tasks/$TASK_ID/claim" "" "$STAFF_TOKEN")
            assert_msg "认领任务" "$CLAIM_RESP" "ok"

            info "  完成任务..."
            local COMPLETE_RESP=$(post "/b/tasks/$TASK_ID/complete" '{"staff_notes":"服务完成"}' "$STAFF_TOKEN")
            assert_msg "完成任务" "$COMPLETE_RESP" "ok"

            info "  用户评价..."
            local RATE_RESP=$(post "/c/requests/$REQ_ID/rate" '{"rating":5,"feedback":"服务态度很好，准时送达"}' "$C_TOKEN")
            assert_msg "评价服务" "$RATE_RESP" "ok"
        else
            warn "未找到对应任务，跳过认领/完成/评价"
        fi
    fi

    # 2.9 取消服务请求
    info "2.9 提交并取消服务请求..."
    local CANCEL_REQ=$(post "/c/requests" \
        "{
            \"service_type\": \"company\",
            \"description\": \"陪伴服务\",
            \"contact_name\": \"张大爷\",
            \"contact_phone\": \"$C_PHONE\",
            \"address\": \"北京市昌平区霍营街道华龙苑北里小区\",
            \"submit_location_lat\": 40.07,
            \"submit_location_lng\": 116.37
        }" "$C_TOKEN")
    local CANCEL_ID=$(echo "$CANCEL_REQ" | grep -o '"id":\s*[0-9]*' | head -1 | grep -o '[0-9]*')
    if [ -n "$CANCEL_ID" ] && [ -z "${CANCEL_ID//[0-9]/}" ]; then
        info "  取消服务请求 (ID=$CANCEL_ID)..."
        local CANCEL_RESP=$(post "/c/requests/$CANCEL_ID/cancel" "" "$C_TOKEN")
        assert_msg "取消服务请求" "$CANCEL_RESP" "ok"
    fi

    # 2.10 更新个人资料
    info "2.10 更新个人资料..."
    local PROFILE_RESP=$(put "/c/profile" \
        '{"name":"张大爷","address":"北京市昌平区霍营街道华龙苑北里小区1号楼3单元501"}' "$C_TOKEN")
    assert_msg "更新个人资料" "$PROFILE_RESP" "ok"

    # 2.11 登出
    info "2.11 登出..."
    local LOGOUT_RESP=$(post "/c/auth/logout" "" "$C_TOKEN")
    assert_msg "登出" "$LOGOUT_RESP" "logout successful"
}

# =====================================================
# 流程 3：家属代提交服务请求
# =====================================================
test_family_submit_flow() {
    echo ""
    echo "=========================================="
    echo -e "${BLUE}流程 3：家属代老人提交服务请求${NC}"
    echo "=========================================="

    local FAMILY_PHONE="13800000011"
    info "使用家属账号: $FAMILY_PHONE"

    # 3.1 家属登录
    info "3.1 家属登录..."
    local LOGIN_RESP=$(post "/c/auth/login" \
        "{\"phone\":\"$FAMILY_PHONE\",\"password\":\"Test@123\",\"type\":\"password\"}")
    local FAMILY_TOKEN=$(echo "$LOGIN_RESP" | grep -o '"token":\s*"[^"]*"' | sed 's/"token":\s*"\([^"]*\)"/\1/')

    if echo "$LOGIN_RESP" | grep -q "\"token\":"; then
        ok "家属登录成功"
        PASSED=$((PASSED + 1)); TOTAL=$((TOTAL + 1))
    else
        fail "家属登录失败，跳过此流程"
        FAILED=$((FAILED + 1)); TOTAL=$((TOTAL + 1)); return 1
    fi

    # 3.2 获取当前用户信息
    info "3.2 获取当前用户信息..."
    local ME_RESP=$(get "/c/auth/me" "$FAMILY_TOKEN")
    assert_msg "获取家属信息" "$ME_RESP" "ok"

    # 3.3 代老人提交服务请求
    info "3.3 代老人提交服务请求（陪诊）..."
    local REQ_RESP=$(post "/c/requests" \
        "{
            \"service_type\": \"medical\",
            \"description\": \"代老人预约陪诊服务，需要帮忙挂号和陪同就医\",
            \"contact_name\": \"张小华\",
            \"contact_phone\": \"$FAMILY_PHONE\",
            \"address\": \"北京市昌平区霍营街道华龙苑北里小区2号楼\",
            \"submit_location_lat\": 40.07,
            \"submit_location_lng\": 116.37,
            \"urgency\": \"urgent\"
        }" "$FAMILY_TOKEN")
    assert_msg "代提交服务请求（陪诊）" "$REQ_RESP" "ok"

    # 3.4 代老人提交送餐请求
    info "3.4 代老人提交服务请求（送餐）..."
    local REQ_RESP2=$(post "/c/requests" \
        "{
            \"service_type\": \"meal\",
            \"description\": \"老年营养餐，低盐低糖\",
            \"contact_name\": \"张小华\",
            \"contact_phone\": \"$FAMILY_PHONE\",
            \"address\": \"北京市昌平区霍营街道华龙苑北里小区2号楼\",
            \"submit_location_lat\": 40.07,
            \"submit_location_lng\": 116.37,
            \"urgency\": \"normal\"
        }" "$FAMILY_TOKEN")
    assert_msg "代提交服务请求（送餐）" "$REQ_RESP2" "ok"

    # 3.5 查看请求列表
    info "3.5 查看服务请求列表..."
    local LIST_RESP=$(get "/c/requests" "$FAMILY_TOKEN")
    assert_msg "服务请求列表" "$LIST_RESP" "ok"

    # 3.6 家属取消请求
    local CANCEL_ID=$(echo "$REQ_RESP2" | grep -o '"id":\s*[0-9]*' | head -1 | grep -o '[0-9]*')
    if [ -n "$CANCEL_ID" ] && [ -z "${CANCEL_ID//[0-9]/}" ]; then
        info "3.6 取消服务请求 (ID=$CANCEL_ID)..."
        local CANCEL_RESP=$(post "/c/requests/$CANCEL_ID/cancel" "" "$FAMILY_TOKEN")
        assert_msg "取消服务请求" "$CANCEL_RESP" "ok"
    fi

    # 3.7 登出
    info "3.7 登出..."
    local LOGOUT_RESP=$(post "/c/auth/logout" "" "$FAMILY_TOKEN")
    assert_msg "登出" "$LOGOUT_RESP" "logout successful"
}

# =====================================================
# 流程 4：B端任务处理全流程
# =====================================================
test_b_backend_flow() {
    echo ""
    echo "=========================================="
    echo -e "${BLUE}流程 4：B端任务处理全流程${NC}"
    echo "=========================================="

    info "4.1 Admin 登录..."
    local ADMIN_RESP=$(post "/b/auth/login" '{"phone":"13800000001","password":"Test@123"}')
    local ADMIN_TOKEN=$(echo "$ADMIN_RESP" | grep -o '"token":\s*"[^"]*"' | sed 's/"token":\s*"\([^"]*\)"/\1/')
    assert_not_empty "Admin token" "$ADMIN_TOKEN"

    info "4.2 Station Manager 登录..."
    local SM_RESP=$(post "/b/auth/login" '{"phone":"13800000002","password":"Test@123"}')
    local SM_TOKEN=$(echo "$SM_RESP" | grep -o '"token":\s*"[^"]*"' | sed 's/"token":\s*"\([^"]*\)"/\1/')

    info "4.3 Staff 登录..."
    local STAFF_RESP=$(post "/b/auth/login" '{"phone":"13800000004","password":"Test@123"}')
    local STAFF_TOKEN=$(echo "$STAFF_RESP" | grep -o '"token":\s*"[^"]*"' | sed 's/"token":\s*"\([^"]*\)"/\1/')

    info "4.4 查看任务池..."
    local POOL_RESP=$(get "/b/tasks/pool?page=1&page_size=5&status=dispatched" "$ADMIN_TOKEN")
    assert_msg "任务池" "$POOL_RESP" "ok"

    info "4.5 查看任务列表..."
    local TASKS_RESP=$(get "/b/tasks?page=1&page_size=10" "$ADMIN_TOKEN")
    assert_msg "任务列表" "$TASKS_RESP" "ok"

    info "4.6 查看我的任务（Staff）..."
    local MY_RESP=$(get "/b/tasks/my" "$STAFF_TOKEN")
    assert_msg "我的任务" "$MY_RESP" "ok"

    info "4.7 查看服务请求列表..."
    local REQ_RESP=$(get "/b/requests?page=1&page_size=10" "$ADMIN_TOKEN")
    assert_msg "服务请求列表" "$REQ_RESP" "ok"

    # 4.8 创建新服务请求（C端）并走完完整流程
    info "4.8 C端创建服务请求用于完整流程测试..."
    local C_RESP=$(post "/c/auth/login" '{"phone":"13800000008","password":"Test@123","type":"password"}')
    local C_TOKEN=$(echo "$C_RESP" | grep -o '"token":\s*"[^"]*"' | sed 's/"token":\s*"\([^"]*\)"/\1/')

    local NEW_REQ=$(post "/c/requests" \
        "{
            \"service_type\": \"care\",
            \"description\": \"日常照护服务测试\",
            \"contact_name\": \"张大爷\",
            \"contact_phone\": \"13800000008\",
            \"address\": \"北京市昌平区霍营街道华龙苑北里小区\",
            \"submit_location_lat\": 40.07,
            \"submit_location_lng\": 116.37
        }" "$C_TOKEN")
    local NEW_REQ_ID=$(echo "$NEW_REQ" | grep -o '"id":\s*[0-9]*' | head -1 | grep -o '[0-9]*')

    if [ -n "$NEW_REQ_ID" ] && [ -z "${NEW_REQ_ID//[0-9]/}" ]; then
        ok "C端创建服务请求成功 (ID=$NEW_REQ_ID)"
        PASSED=$((PASSED + 1)); TOTAL=$((TOTAL + 1))

        # 4.9 Admin 取消此请求
        info "4.9 Admin 取消请求 (ID=$NEW_REQ_ID)..."
        local CANCEL_RESP=$(post "/b/requests/$NEW_REQ_ID/cancel" "" "$ADMIN_TOKEN")
        assert_msg "Admin取消请求" "$CANCEL_RESP" "ok"

        # 4.10 再创建一个用于转派测试
        local TRANSFER_REQ=$(post "/c/requests" \
            "{
                \"service_type\": \"repair\",
                \"description\": \"家电维修测试\",
                \"contact_name\": \"张大爷\",
                \"contact_phone\": \"13800000008\",
                \"address\": \"北京市昌平区霍营街道华龙苑北里小区\",
                \"submit_location_lat\": 40.07,
                \"submit_location_lng\": 116.37
            }" "$C_TOKEN")
        local TRANSFER_REQ_ID=$(echo "$TRANSFER_REQ" | grep -o '"id":\s*[0-9]*' | head -1 | grep -o '[0-9]*')

        if [ -n "$TRANSFER_REQ_ID" ] && [ -z "${TRANSFER_REQ_ID//[0-9]/}" ]; then
            local STAFF2_RESP=$(post "/b/auth/login" '{"phone":"13800000005","password":"Test@123"}')
            local STAFF2_TOKEN=$(echo "$STAFF2_RESP" | grep -o '"token":\s*"[^"]*"' | sed 's/"token":\s*"\([^"]*\)"/\1/')

            # 查找对应任务
            local ALL_TASKS=$(get "/b/tasks?page_size=200" "$ADMIN_TOKEN")
            # jq: 精确匹配 request_id -> task id
            local TRANSFER_TASK_ID=$(echo "$ALL_TASKS" | jq --argjson rid "$TRANSFER_REQ_ID" '.data.items[] | select(.request_id == $rid) | .id' 2>/dev/null)
            # fallback: 精确匹配 "request_id":N 的前一个字段
            if [ -z "$TRANSFER_TASK_ID" ]; then
                TRANSFER_TASK_ID=$(echo "$ALL_TASKS" | grep -oP "\"request_id\":$TRANSFER_REQ_ID," | head -1 | grep -oP '.{0,20}"id":\K[0-9]+')
            fi

            if [ -n "$TRANSFER_TASK_ID" ]; then
                info "4.10 认领任务 (TaskID=$TRANSFER_TASK_ID)..."
                local CLAIM=$(post "/b/tasks/$TRANSFER_TASK_ID/claim" "" "$STAFF_TOKEN")
                assert_msg "认领任务" "$CLAIM" "ok"

                info "4.11 Admin 转派任务给 Staff2..."
                local TRANSFER=$(post "/b/tasks/$TRANSFER_TASK_ID/transfer" '{"staff_id":5}' "$ADMIN_TOKEN")
                assert_msg "转派任务" "$TRANSFER" "ok"

                info "4.12 Staff2 完成任务..."
                local COMPLETE=$(post "/b/tasks/$TRANSFER_TASK_ID/complete" '{"staff_notes":"家电维修完成"}' "$STAFF2_TOKEN")
                assert_msg "完成任务" "$COMPLETE" "ok"

                info "4.13 用户评价..."
                local RATE=$(post "/c/requests/$TRANSFER_REQ_ID/rate" '{"rating":5,"feedback":"维修师傅很专业"}' "$C_TOKEN")
                assert_msg "用户评价" "$RATE" "ok"
            else
                warn "未找到对应任务，跳过转派测试"
            fi
        fi
    else
        warn "C端创建请求失败，跳过后续测试"
    fi

    # 4.14 B端登出
    info "4.14 B端登出..."
    local B_LOGOUT=$(post "/b/auth/logout" "" "$ADMIN_TOKEN")
    assert_msg "Admin登出" "$B_LOGOUT" "logout successful"
}

# =====================================================
# 流程 5：Token 隔离与权限验证
# =====================================================
test_security_flow() {
    echo ""
    echo "=========================================="
    echo -e "${BLUE}流程 5：Token 隔离与权限验证${NC}"
    echo "=========================================="

    # 未登录访问需认证接口
    info "5.1 未登录访问需认证接口（应401）..."
    local NO_AUTH=$(get "/c/auth/me")
    assert_code "未登录 - C端 me" "$NO_AUTH" "401"

    # B端 token 访问 C端
    info "5.2 B端Token访问C端（应403）..."
    local B_LOGIN=$(post "/b/auth/login" '{"phone":"13800000001","password":"Test@123"}')
    local B_TOKEN=$(echo "$B_LOGIN" | grep -o '"token":\s*"[^"]*"' | sed 's/"token":\s*"\([^"]*\)"/\1/')
    local B_TO_C=$(get "/c/auth/me" "$B_TOKEN")
    assert_code "B端Token访问C端" "$B_TO_C" "403"

    # C端 token 访问 B端
    info "5.3 C端Token访问B端（应403）..."
    local C_LOGIN=$(post "/c/auth/login" '{"phone":"13800000008","password":"Test@123","type":"password"}')
    local C_TOKEN=$(echo "$C_LOGIN" | grep -o '"token":\s*"[^"]*"' | sed 's/"token":\s*"\([^"]*\)"/\1/')
    local C_TO_B=$(get "/b/auth/me" "$C_TOKEN")
    assert_code "C端Token访问B端" "$C_TO_B" "403"

    # 公开接口无需认证
    info "5.4 公开接口无需认证（新闻列表）..."
    local NEWS=$(get "/c/news")
    assert_msg "新闻列表（公开）" "$NEWS" "ok"

    info "5.5 公开接口无需认证（Banner列表）..."
    local BANNERS=$(get "/c/banners")
    assert_msg "Banner列表（公开）" "$BANNERS" "ok"

    info "5.6 公开接口无需认证（站点匹配）..."
    local MATCH=$(get "/c/stations/match?lng=116.38&lat=40.05")
    assert_msg "站点匹配（公开）" "$MATCH" "ok"
}

# =====================================================
# 主流程
# =====================================================
echo ""
echo "╔══════════════════════════════════════════════════╗"
echo "║     sCare C端业务流程测试  v1.0                  ║"
echo "║     Base URL: $API_URL                           ║"
echo "╚══════════════════════════════════════════════════╝"

# 健康检查
info "健康检查..."
HEALTH=$(curl -s "$API_URL/health")
if echo "$HEALTH" | grep -q '"msg":"ok"'; then
    ok "服务健康"
else
    fail "服务未启动，请先运行后端"
    echo "HEALTH: $HEALTH"
    exit 1
fi

test_security_flow
test_quick_start_flow
test_login_submit_flow
test_family_submit_flow
test_b_backend_flow

# =====================================================
# 测试结果汇总
# =====================================================
echo ""
echo "=========================================="
echo "测试结果汇总"
echo "=========================================="
echo -e "总计: $TOTAL"
echo -e "${GREEN}通过: $PASSED${NC}"
echo -e "${RED}失败: $FAILED${NC}"

if [ $FAILED -eq 0 ]; then
    echo -e "\n${GREEN}✓ 所有业务流程测试通过！${NC}"
    exit 0
else
    echo -e "\n${RED}✗ 有 $FAILED 个测试失败${NC}"
    exit 1
fi
