source common.sh

set -x
curl -w "%{http_code}\n" ${base_url}/live
curl http://localhost:8080/api/liveness?false
curl -w "%{http_code}\n" ${base_url}/live
curl http://localhost:8080/api/liveness?true
curl -w "%{http_code}\n" ${base_url}/live
