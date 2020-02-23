source common.sh

set -x
curl -X POST ${base_url}/api/delay
curl -X POST ${base_url}/api/delay?nou
curl -X POST ${base_url}/api/delay?value
curl -X POST ${base_url}/api/delay?value=nou
curl -X GET ${base_url}/api/delay?value=2

curl -X POST ${base_url}/api/delay?value=2
time curl ${base_url}/
curl -X POST ${base_url}/api/delay?value=0
