source common.sh

set -x
curl -X POST ${base_url}/api/cpu?value=0.1
time sleep 10
curl -X POST ${base_url}/api/cpu?value=1
time sleep 10
curl -X POST ${base_url}/api/cpu?value=16
time sleep 5
curl -X POST ${base_url}/api/cpu?value=1000
time sleep 5
curl -X POST ${base_url}/api/cpu?value=0
