// +build ignore
// FIXME: wewrite against pure client-go

package data

import (
	"context"
	"errors"
	"log"
	"os"
	"strconv"
	"strings"

	//"crypto/rsa"
	//"crypto/x509"
	//"encoding/pem"
	"fmt"
	"io/ioutil"

	"github.com/dgrijalva/jwt-go"
	"github.com/ericchiang/k8s"
	corev1 "github.com/ericchiang/k8s/apis/core/v1"
)

func init() {
	plugins = append(plugins, getK8sData)
}

func getK8sData() map[string]string {
	data := map[string]string{}

	client, err := k8s.NewInClusterClient()

	data["K8s"] = "no"
	if err == nil {
		data["K8s"] = "yes"

		disco := k8s.NewDiscoveryClient(client)
		version, _ := disco.Version(context.Background()) // TODO: better context (timeout of 1s propagated from the root
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
		token, err := jwt.ParseWithClaims(saToken, &k8sClaims{}, nil) /*func(token *jwt.Token) (interface{}, error) {
		    // This JWT is signed with the Service Account keypair.
		    // This isn't the same as the apiserver CA keypair, so ca.crt on the disk can't validate it
		    // As of Mar '20 there's no way to directly get the public part of the SA pair
		    // Options are:
		    //   * Call the TokenReview API: https://kubernetes.io/docs/reference/access-authn-authz/authentication/
		    //   * Wait for OICD discovery to be implemented, and use that to get the public key: https://github.com/kubernetes/enhancements/blob/master/keps/sig-auth/20190730-oidc-discovery.md
		    //     * maybe some of the work is there? https://kubernetes.io/docs/tasks/configure-pod-container/configure-service-account/

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
		if err != nil {
			// log.Fatalf("can't parse and verify token: %v", err) // Will fail atm because we pass a nil keyFunc, but the token is still parsed, just not validated
		}

		claims, ok := token.Claims.(*k8sClaims)
		if !ok /*|| !token.Valid*/ {
			log.Fatal("ServiceAccount token invalid")
		}
		data["K8sNamespace"] = claims.Namespace
		data["K8sServiceAccount"] = claims.Name

		data["K8sVersion"] = fmt.Sprintf("%s %s", version.GitVersion, version.Platform)

		hostname, _ := os.Hostname()
		var thisPod corev1.Pod
		if err := client.Get(context.Background(), claims.Namespace, hostname, &thisPod); err != nil {
			var k8sErr *k8s.APIError
			if errors.As(err, &k8sErr) {
				log.Printf("Forbidden to get Pod, check RBAC: %v", k8sErr)
				data["K8sThisPodContainersCount"] = "forbidden"
			} else {
				log.Fatalf("Can't get Pod: %v", err)
			}
		} else {
			data["K8sThisPodContainersCount"] = strconv.Itoa(len(thisPod.Spec.Containers))

			data["K8sThisPodContainersImages"] = ""
			for _, c := range thisPod.Spec.Containers {
				data["K8sThisPodContainersImages"] += *c.Image + ", "
			}
		}

		var nodes corev1.NodeList
		if err := client.List(context.Background(), "", &nodes); err != nil {
			var k8sErr *k8s.APIError
			if errors.As(err, &k8sErr) {
				log.Printf("Forbidden to list nodes, check RBAC: %v", k8sErr)
				data["K8sNodeCount"] = "forbidden"
			} else {
				log.Fatalf("Can't list Nodes: %v", err)
			}
		} else {
			data["K8sNodeCount"] = strconv.Itoa(len(nodes.Items))
			var self *corev1.Node
			for _, self = range nodes.Items {
				if *self.Metadata.Name == *thisPod.Spec.NodeName {
					data["K8sNodeAddress"] = *self.Status.Addresses[0].Address + " / " + *self.Status.Addresses[1].Address // TODO loop
					data["K8sNodeVersion"] = fmt.Sprintf("%s %s/%s", *self.Status.NodeInfo.KubeletVersion, *self.Status.NodeInfo.OperatingSystem, *self.Status.NodeInfo.Architecture)
					data["K8sNodeRuntime"] = *self.Status.NodeInfo.ContainerRuntimeVersion
					data["K8sNodeOS"] = *self.Status.NodeInfo.OsImage
					data["K8sNodeRole"] = findSuffix(self.Metadata.Labels, "node-role.kubernetes.io/")
					data["K8sNodeCloudInstance"] = self.Metadata.Labels["node.kubernetes.io/instance-type"]
					data["K8sNodeCloudRegion"] = self.Metadata.Labels["topology.kubernetes.io/region"]
					data["K8sNodeCloudZone"] = self.Metadata.Labels["topology.kubernetes.io/zone"]
				}
			}
		}

		// TODO: get own pod container list
		// TODO: get own namespace pods list
		// TODO: search for parent ownerrefs, and build kubectl tree style location
	}

	return data
}

func findSuffix(m map[string]string, pre string) string {
	for k, _ := range m {
		if strings.HasPrefix(k, pre) {
			return strings.TrimPrefix(k, pre)
		}
	}
	return ""
}
