source common.sh

set -x
curl ${base_url}/handlers/cpu?value=0.1
time sleep 10
curl ${base_url}/handlers/cpu?value=1
time sleep 10
curl ${base_url}/handlers/cpu?value=16
time sleep 5
curl ${base_url}/handlers/cpu?value=1000
time sleep 5
curl ${base_url}/handlers/cpu?value=0
