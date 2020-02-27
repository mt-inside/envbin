source common.sh

set -x
curl -w "%{http_code}\n" ${base_url}/ready
curl ${base_url}/api/readiness?value=false
curl -w "%{http_code}\n" ${base_url}/ready
curl ${base_url}/api/readiness?value=true
curl -w "%{http_code}\n" ${base_url}/ready
