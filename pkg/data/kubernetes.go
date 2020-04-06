package data

import (
	"context"
	"log"
	"os"
	"strconv"

	//"crypto/rsa"
	//"crypto/x509"
	//"encoding/pem"
	"fmt"
	"github.com/ericchiang/k8s"
	corev1 "github.com/ericchiang/k8s/apis/core/v1"
	"io/ioutil"
	"github.com/dgrijalva/jwt-go"
)

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
			Secret string `json:"kubernetes.io/serviceaccount/secret.name"`
			Name string `json:"kubernetes.io/serviceaccount/service-account.name"`
			Uid string `json:"kubernetes.io/serviceaccount/service-account.uid"`

			jwt.StandardClaims
		}

		// TODO: unit test this with a ca.crt and satoek from a pod
		token, err := jwt.ParseWithClaims(saToken, &k8sClaims{}, nil, /*func(token *jwt.Token) (interface{}, error) {
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
		}*/)
		if err != nil {
			// log.Fatalf("can't parse and verify token: %v", err) // Will fail atm because we pass a nil keyFunc, but the token is still parsed, just not validated
		}

		if claims, ok := token.Claims.(*k8sClaims); ok /*&& token.Valid*/ {
			data["K8sNamespace"] = claims.Namespace
			data["K8sServiceAccount"] = claims.Name
		}

		data["K8sVersion"] = fmt.Sprintf("%s %s",version.GitVersion, version.Platform)

		var nodes corev1.NodeList
		if err := client.List(context.Background(), "", &nodes); err != nil {
			log.Fatalf("Can't list Nodes: %v", err)
		}
		data["K8sNodeCount"] = strconv.Itoa(len(nodes.Items))
		var self *corev1.Node
		hostname, _ := os.Hostname()
		for _, self = range nodes.Items {
			if *self.Metadata.Name == hostname { break }
		}
		log.Print(*self.Status.NodeInfo.ContainerRuntimeVersion)
		log.Print(*self.Status.NodeInfo.KubeletVersion)
		log.Print(*self.Status.NodeInfo.KubeProxyVersion)
		log.Print(*self.Status.NodeInfo.OsImage)
		log.Print(self.Metadata.Labels["node.kubernetes.io/instance-type"])
		log.Print(self.Metadata.Labels["topology.kubernetes.io/region"])
		log.Print(self.Metadata.Labels["topology.kubernetes.io/zone"])

		// TODO: get own pod container list
		// TODO: get own namespace pods list
		// TODO: search for parent ownerrefs, and build kubectl tree style location
	}

	return data
}
