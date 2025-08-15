#!/bin/bash

echo "🚀 Testing SisuPass Email Client Service"
echo "========================================"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

BASE_URL="http://localhost:4000"

# Test health endpoint
echo -e "${YELLOW}1. Testing health endpoint...${NC}"
health_response=$(curl -s "$BASE_URL/health")
if [[ $? -eq 0 ]]; then
    echo -e "${GREEN}✅ Health check passed${NC}"
    echo "$health_response" | jq .
else
    echo -e "${RED}❌ Health check failed${NC}"
fi

echo ""

# Test templates endpoint
echo -e "${YELLOW}2. Testing templates endpoint...${NC}"
templates_response=$(curl -s "$BASE_URL/api/v1/email/templates")
if [[ $? -eq 0 ]]; then
    echo -e "${GREEN}✅ Templates endpoint works${NC}"
    echo "$templates_response" | jq .
else
    echo -e "${RED}❌ Templates endpoint failed${NC}"
fi

echo ""

# Test stats endpoint
echo -e "${YELLOW}3. Testing stats endpoint...${NC}"
stats_response=$(curl -s "$BASE_URL/api/v1/email/stats")
if [[ $? -eq 0 ]]; then
    echo -e "${GREEN}✅ Stats endpoint works${NC}"
    echo "$stats_response" | jq .
else
    echo -e "${RED}❌ Stats endpoint failed${NC}"
fi

echo ""

# Test send email endpoint
echo -e "${YELLOW}4. Testing send email endpoint...${NC}"
email_response=$(curl -s -X POST "$BASE_URL/api/v1/email/send" \
  -H "Content-Type: application/json" \
  -d '{
    "to": "test@example.com",
    "template": "user_welcome",
    "data": {
      "name": "Test User",
      "activation_token": "test123",
      "activation_url": "http://localhost:3000/activate?token=test123"
    },
    "priority": "normal"
  }')

if [[ $? -eq 0 ]]; then
    echo -e "${GREEN}✅ Send email endpoint works${NC}"
    echo "$email_response" | jq .
    
    # Extract email ID for status check
    email_id=$(echo "$email_response" | jq -r '.email.id')
    
    if [[ "$email_id" != "null" && "$email_id" != "" ]]; then
        echo ""
        echo -e "${YELLOW}5. Testing email status endpoint...${NC}"
        status_response=$(curl -s "$BASE_URL/api/v1/email/$email_id/status")
        if [[ $? -eq 0 ]]; then
            echo -e "${GREEN}✅ Email status endpoint works${NC}"
            echo "$status_response" | jq .
        else
            echo -e "${RED}❌ Email status endpoint failed${NC}"
        fi
    fi
else
    echo -e "${RED}❌ Send email endpoint failed${NC}"
fi

echo ""
echo -e "${BLUE}📖 Access Swagger UI at: $BASE_URL/swagger/${NC}"
echo -e "${GREEN}🎉 Test completed!${NC}"
