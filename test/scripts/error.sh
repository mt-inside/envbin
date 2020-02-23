source common.sh

set -x
curl -X POST ${base_url}/api/errorrate?value=1
curl -o /dev/null -s -w "%{http_code}\n" ${base_url}
curl -X POST ${base_url}/api/errorrate?value=0.5
curl -o /dev/null -s -w "%{http_code}\n" ${base_url}
curl -o /dev/null -s -w "%{http_code}\n" ${base_url}
curl -o /dev/null -s -w "%{http_code}\n" ${base_url}
curl -o /dev/null -s -w "%{http_code}\n" ${base_url}
curl -X POST ${base_url}/api/errorrate?value=0
