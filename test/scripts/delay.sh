source common.sh

set -x
curl ${base_url}/handlers/delay
curl ${base_url}/handlers/delay?nou
curl ${base_url}/handlers/delay?value
curl ${base_url}/handlers/delay?value=nou
curl -X POST ${base_url}/handlers/delay?value=2

curl ${base_url}/handlers/delay?value=2
time curl ${base_url}/
curl ${base_url}/handlers/delay?value=0
