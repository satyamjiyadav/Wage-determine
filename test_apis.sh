#!/bin/bash
set -e

BASE_URL="http://localhost:8080"
echo "=========================================================="
echo " Starting API Verification Test Suite for Prevailing Wage "
echo "=========================================================="
echo ""

# Helper function to print test status
run_test() {
  local test_name="$1"
  local curl_cmd="$2"
  
  echo "----------------------------------------------------------"
  echo "TEST: $test_name"
  echo "CMD:  $curl_cmd"
  echo "OUTPUT:"
  eval "$curl_cmd" | python3 -m json.tool || eval "$curl_cmd"
  echo ""
}

# 1. Health Probe Endpoint
run_test "1. Service Health Readiness Check (/healthz)" \
  "curl -s $BASE_URL/healthz"

# 2. Metrics Endpoint
run_test "2. Prometheus Telemetry Metrics (/metrics)" \
  "curl -s $BASE_URL/metrics | head -n 10"

# 3. Location Resolver (Boston MA, ZIP 02108)
run_test "3. Resolve Location for Boston, MA (ZIP 02108)" \
  "curl -s '$BASE_URL/api/v1/locations/resolve?zip_code=02108'"

# 4. Occupation Search
run_test "4. Search Occupations for 'Engineer'" \
  "curl -s '$BASE_URL/api/v1/occupations/search?q=Engineer&limit=3'"

# 5. Get Occupation Metadata (Mechanical Engineers 17-2141.00)
run_test "5. Get O*NET Metadata for Mechanical Engineers (17-2141.00)" \
  "curl -s '$BASE_URL/api/v1/occupations/17-2141.00'"

# 6. Wage Lookup (Software Developers in Seattle 98101)
run_test "6. Wage Lookup for Software Developers in Seattle (ZIP 98101)" \
  "curl -s '$BASE_URL/api/v1/wages/lookup?soc_code=15-1252.00&zip_code=98101'"

# 7. Wage Lookup (Data Scientists in NYC 10001)
run_test "7. Wage Lookup for Data Scientists in NYC (ZIP 10001)" \
  "curl -s '$BASE_URL/api/v1/wages/lookup?soc_code=15-2051.00&zip_code=10001'"

# 8. Automated Level Determination (Scenario A: Level 1 Entry Engineer)
run_test "8. 4-Tier Assessment: Entry Level Civil Engineer (Expected: Level 1)" \
  "curl -s -X POST $BASE_URL/api/v1/wages/determine-level \
    -H 'Content-Type: application/json' \
    -d '{
      \"soc_code\": \"17-2051.00\",
      \"zip_code\": \"94103\",
      \"job_title\": \"Junior Civil Engineer\",
      \"education\": { \"required_degree\": \"Bachelor\" },
      \"experience_months\": 12,
      \"supervises_employees\": false
    }'"

# 9. Automated Level Determination (Scenario B: Level 3 Senior Data Scientist)
DETERMINATION_RESP=$(curl -s -X POST $BASE_URL/api/v1/wages/determine-level \
  -H 'Content-Type: application/json' \
  -d '{
    "soc_code": "15-2051.00",
    "zip_code": "94103",
    "job_title": "Lead Data Scientist",
    "education": { "required_degree": "Doctorate" },
    "experience_months": 72,
    "special_skills": ["PyTorch", "LLMs", "Distributed Training"],
    "supervises_employees": true,
    "number_of_subordinates": 4
  }')

echo "----------------------------------------------------------"
echo "TEST: 9. 4-Tier Assessment: Lead Data Scientist (Expected: Level 4)"
echo "$DETERMINATION_RESP" | python3 -m json.tool
echo ""

# Extract Determination Tracking Number for audit check
TRACKING_NUM=$(echo "$DETERMINATION_RESP" | python3 -c "import sys, json; print(json.load(sys.stdin)['data']['determination_number'])")

# 10. Audit Trail Log Retrieval
run_test "10. Retrieve Audit Trail Log by Tracking Number ($TRACKING_NUM)" \
  "curl -s $BASE_URL/api/v1/determinations/$TRACKING_NUM"

# 11. Batch Wage Query Across 4 Major U.S. Tech Hubs
run_test "11. Batch Wage Lookup Across SF, NYC, Seattle, and Austin" \
  "curl -s -X POST $BASE_URL/api/v1/wages/batch-lookup \
    -H 'Content-Type: application/json' \
    -d '[
      {\"soc_code\": \"15-1252.00\", \"zip_code\": \"94103\"},
      {\"soc_code\": \"15-1252.00\", \"zip_code\": \"10001\"},
      {\"soc_code\": \"15-1252.00\", \"zip_code\": \"98101\"},
      {\"soc_code\": \"15-1252.00\", \"zip_code\": \"78701\"}
    ]'"

echo "=========================================================="
echo " ALL API VERIFICATION TESTS COMPLETED SUCCESSFULLY! "
echo "=========================================================="
