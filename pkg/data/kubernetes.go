// +buil ignore
// FIXME: wewrite against pure client-go

package data

import (
	"context"
	"os"
	"strconv"
	"strings"

	//"crypto/rsa"
	//"crypto/x509"
	//"encoding/pem"
	"fmt"
	"io/ioutil"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-logr/logr"
	k8sErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func init() {
	plugins = append(plugins, getK8sData)
}

func getK8sData(ctx context.Context, log logr.Logger, t *Trie) {
	config, err := rest.InClusterConfig()
	if err != nil {
		t.Insert(NotPresent{}, "Cloud", "Kubernetes")
		return
	}
	t.Insert(Some{"Present"}, "Cloud", "Kubernetes")

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Error(err, "Can't connect to k8s apiserver")
		return
	}

	version, err := clientSet.Discovery().ServerVersion()
	if err != nil {
		log.Error(err, "Can't get cluster version")
	} else {
		t.Insert(Some{fmt.Sprintf("%s %s", version.GitVersion, version.Platform)}, "Cloud", "Kubernetes", "Cluster", "Version")
	}

	saBytes, _ := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/token")
	saToken := string(saBytes)

	type k8sClaims struct {
		Namespace string `json:"kubernetes.io/serviceaccount/namespace"`
		Secret    string `json:"kubernetes.io/serviceaccount/secret.name"`
		Name      string `json:"kubernetes.io/serviceaccount/service-account.name"`
		Uid       string `json:"kubernetes.io/serviceaccount/service-account.uid"`

		jwt.StandardClaims
	}

	// TODO: unit test this with a ca.crt and satoek from a pod
	token, err := jwt.ParseWithClaims(saToken, &k8sClaims{}, nil)
	if err != nil {
		log.Error(err, "can't parse and verify token") // Will fail atm because we pass a nil keyFunc, but the token is still parsed, just not validated
		return
	}
	// This JWT is signed with the Service Account keypair.
	// This isn't the same as the apiserver CA keypair, so ca.crt on the disk can't validate it
	// As of Mar '20 there's no way to directly get the public part of the SA pair
	// Options are:
	//   * Call the TokenReview API: https://kubernetes.io/docs/reference/access-authn-authz/authentication/
	//   * Wait for OICD discovery to be implemented, and use that to get the public key: https://github.com/kubernetes/enhancements/blob/master/keps/sig-auth/20190730-oidc-discovery.md
	//     * maybe some of the work is there? https://kubernetes.io/docs/tasks/configure-pod-container/configure-service-account/

	/*func(token *jwt.Token) (interface{}, error) {
	can allegedly do this:
	      containers:
		        - name: envbin
		          image: mtinside/envbin:latest
		          volumeMounts:
		            - mountPath: /var/run/secrets/tokens
		              name: vault-token
		      volumes:
		        - name: vault-token
		          projected:
		            sources:
		            - serviceAccountToken:
		                path: vault-token
		                expirationSeconds: 7200
		                audience: vault

	but even with `minikube start --feature-gates='TokenRequest=true,TokenRequestProjection=true'`, Pods don't start


	    // verify rs256 alg

		asn1, err := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/ca.crt") - no!
		if err != nil {
			log.Fatalf("Can't read ca cert: %v", err)
		}

		block, _ := pem.Decode(asn1)
		if block == nil || block.Type != "CERTIFICATE" {
			log.Fatalf("Can't decode PEM: %v", err)
		}

		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			log.Fatalf("Can't parse ca cert: %v", err)
		}

		var pub *rsa.PublicKey
		var ok bool
		if pub, ok = cert.PublicKey.(*rsa.PublicKey); !ok {
			log.Fatalf("can't get public key from cert; %v", err)
		}

		return pub, nil
	}*/

	claims, ok := token.Claims.(*k8sClaims)
	if !ok /*|| !token.Valid*/ {
		log.Error(fmt.Errorf("ServiceAccount token invalid"), "Can't read k8s info")
		return
	}
	t.Insert(Some{claims.Namespace}, "Cloud", "Kubernetes", "Pod", "Namespace")
	t.Insert(Some{claims.Name}, "Cloud", "Kubernetes", "Pod", "ServiceAccount")

	hostname, _ := os.Hostname()
	pod, err := clientSet.CoreV1().Pods(claims.Namespace).Get(ctx, hostname, metav1.GetOptions{})
	if err != nil {
		if k8sErrors.IsForbidden(err) {
			log.Error(err, "Forbidden getting own Pod info; check RBAC")
			t.Insert(Some{"Forbidden"}, "Cloud", "Kubernetes", "Pod")
		} else if err == context.DeadlineExceeded {
			log.Error(err, "Timed out getting own Pod info")
			t.Insert(Some{"Timeout"}, "Cloud", "Kubernetes", "Pod")
		} else if k8sErrors.IsTimeout(err) { // client-go blew its own deadline? Is also an IsServerTimeout() to show the apiserver popped its deadline
			log.Error(err, "Timed out getting own Pod info")
			t.Insert(Some{"Timeout"}, "Cloud", "Kubernetes", "Pod")
		} else {
			log.Error(err, "Error getting own Pod info")
		}
	} else {
		t.Insert(Some{strconv.Itoa(len(pod.Spec.Containers))}, "Cloud", "Kubernetes", "Pod", "ContainersCount")

		images := []string{}
		for _, c := range pod.Spec.Containers {
			images = append(images, c.Image)
		}
		t.Insert(Some{strings.Join(images, ",")}, "Cloud", "Kubernetes", "Pod", "ContainersImages")

		if node, err := clientSet.CoreV1().Nodes().Get(ctx, pod.Spec.NodeName, metav1.GetOptions{}); err != nil {
			if k8sErrors.IsForbidden(err) {
				log.Error(err, "Forbidden getting own Node info; check RBAC")
				t.Insert(Some{"Forbidden"}, "Cloud", "Kubernetes", "Node")
			} else if err == context.DeadlineExceeded {
				log.Error(err, "Timed out getting own Node info")
				t.Insert(Some{"Timeout"}, "Cloud", "Kubernetes", "Node")
			} else if k8sErrors.IsTimeout(err) { // client-go blew its own deadline? Is also an IsServerTimeout() to show the apiserver popped its deadline
				log.Error(err, "Timed out getting own Node info")
				t.Insert(Some{"Timeout"}, "Cloud", "Kubernetes", "Node")
			} else {
				log.Error(err, "Error getting own Node info")
			}
		} else {
			t.Insert(Some{node.Status.Addresses[0].Address + " / " + node.Status.Addresses[1].Address}, "Cloud", "Kubernetes", "Node", "Address") // TODO loop
			t.Insert(Some{fmt.Sprintf("%s %s/%s", node.Status.NodeInfo.KubeletVersion, node.Status.NodeInfo.OperatingSystem, node.Status.NodeInfo.Architecture)}, "Cloud", "Kubernetes", "Node", "Version")
			t.Insert(Some{node.Status.NodeInfo.ContainerRuntimeVersion}, "Cloud", "Kubernetes", "Node", "ContainerRuntime")
			t.Insert(Some{node.Status.NodeInfo.OSImage}, "Cloud", "Kubernetes", "Node", "OS")
			t.Insert(Some{findSuffix(node.Labels, "node-role.kubernetes.io/")}, "Cloud", "Kubernetes", "Node", "Role")
			t.Insert(Some{node.Labels["node.kubernetes.io/instance-type"]}, "Cloud", "Kubernetes", "Node", "InstanceType")
			t.Insert(Some{node.Labels["topology.kubernetes.io/region"]}, "Cloud", "Kubernetes", "Node", "Region")
			t.Insert(Some{node.Labels["topology.kubernetes.io/zone"]}, "Cloud", "Kubernetes", "Node", "Zone")
		}
	}

	// TODO: get own namespace pods list
	// TODO: get nodes list: t.Insert(Some{strconv.Itoa(len(nodes.Items))}, "Cloud", "Kubernetes", "Cluster", "NodesCount")
	// TODO: search for parent ownerrefs, and build kubectl tree style location
}

func findSuffix(m map[string]string, pre string) string {
	for k := range m {
		if strings.HasPrefix(k, pre) {
			return strings.TrimPrefix(k, pre)
		}
	}
	return ""
}
