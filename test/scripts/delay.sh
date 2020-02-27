source common.sh

set -x
curl ${base_url}/api/delay
curl ${base_url}/api/delay?nou
curl ${base_url}/api/delay?value
curl ${base_url}/api/delay?value=nou
curl -X POST ${base_url}/api/delay?value=2

curl ${base_url}/api/delay?value=2
time curl ${base_url}/
curl ${base_url}/api/delay?value=0
