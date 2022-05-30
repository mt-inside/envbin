// TODO: read namespace labels and annotations
// TODO: read pod labels and annotions
// TODO: can find the Service(s)? Find endpoint with same IP, that has a ref to the service?

package fetchers

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"github.com/go-logr/logr"
	"github.com/golang-jwt/jwt/v4"
	k8sErrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/mt-inside/envbin/pkg/data"
	"github.com/mt-inside/envbin/pkg/data/trie"
)

func init() {
	data.RegisterPlugin(getK8sDataInCluster)
}

// inCluster wraps generic getter
// cmdline to call generic getting with its own config (clientSet?)
// both to return a trie
// * cmdline thing to share a renderer with "dump"

func getK8sDataInCluster(ctx context.Context, log logr.Logger, vals chan<- trie.InsertMsg) {
	config, err := rest.InClusterConfig()
	if err != nil {
		vals <- trie.Insert(trie.NotPresent(), "Cloud", "Kubernetes")
		return
	}

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		vals <- trie.Insert(trie.Error(fmt.Errorf("can't connect to apiserver: %w", err)), "Cloud", "Kubernetes")
		return
	}

	getK8sData(ctx, log, clientSet, vals)
}

type k8sClaims struct {
	jwt.RegisteredClaims

	Kubernetes struct {
		Namespace string `json:"namespace"`
		Pod       struct {
			Name string `json:"name"`
			Uid  string `json:"uid"`
		} `json:"pod"`
		ServiceAccount struct {
			Name string `json:"name"`
			Uid  string `json:"uid"`
		} `json:"serviceaccount"`
	} `json:"kubernetes.io"`
}

func getK8sData(ctx context.Context, log logr.Logger, clientSet *kubernetes.Clientset, vals chan<- trie.InsertMsg) {
	// Control Plane
	version, err := clientSet.Discovery().ServerVersion()
	if err != nil {
		vals <- trie.Insert(k8sFromError(log, fmt.Errorf("can't get cluster version: %w", err)), "Cloud", "Kubernetes", "Masters")
	} else {
		vals <- trie.Insert(trie.Some(version.GitVersion), "Cloud", "Kubernetes", "Masters", "Version")
		vals <- trie.Insert(trie.Some(version.Platform), "Cloud", "Kubernetes", "Masters", "Platform")
	}

	// IP Ranges

	//svcRange(ctx, clientSet) // TODO

	// Namespace & Identity

	saBytes, _ := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/token")
	saToken := string(saBytes)
	// TODO: unit test this with a ca.crt and satoek from a pod
	token, err := jwt.ParseWithClaims(saToken, &k8sClaims{}, nil)
	// if err != nil {
	// Will fail atm because we pass a nil keyFunc, but the token is still parsed, just not validated
	//vals <- trie.Insert(k8sFromError(log, err), "Cloud", "Kubernetes", "Identity")
	//return
	// }
	claims, ok := token.Claims.(*k8sClaims)
	if !ok {
		vals <- trie.Insert(k8sFromError(log, err), "Cloud", "Kubernetes", "Identity")
		return
	}
	vals <- trie.Insert(trie.Some(claims.Kubernetes.Namespace), "Cloud", "Kubernetes", "Namespace", "Name")
	vals <- trie.Insert(trie.Some(claims.Kubernetes.ServiceAccount.Name), "Cloud", "Kubernetes", "Pod", "ServiceAccount")

	// Nodes

	nodes, err := clientSet.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		vals <- trie.Insert(k8sFromError(log, err), "Cloud", "Kubernetes", "Nodes")
		// Keep going becuase we might have permission to read other stuff
	}
	for i, n := range nodes.Items {
		vals <- trie.Insert(trie.Some(n.Status.NodeInfo.Architecture), "Cloud", "Kubernetes", "Nodes", strconv.Itoa(i), "Architecture")
		vals <- trie.Insert(trie.Some(quantity2str(n.Status.Capacity.Cpu())), "Cloud", "Kubernetes", "Nodes", strconv.Itoa(i), "Cores")
		vals <- trie.Insert(trie.Some(quantity2str(n.Status.Capacity.Memory())), "Cloud", "Kubernetes", "Nodes", strconv.Itoa(i), "RAM")
		vals <- trie.Insert(trie.Some(quantity2str(n.Status.Capacity.Storage())), "Cloud", "Kubernetes", "Nodes", strconv.Itoa(i), "Storage")
		vals <- trie.Insert(trie.Some(quantity2str(n.Status.Capacity.StorageEphemeral())), "Cloud", "Kubernetes", "Nodes", strconv.Itoa(i), "Ephemeral")
		for k, v := range n.Labels {
			vals <- trie.Insert(trie.Some(v), "Cloud", "Kubernetes", "Nodes", strconv.Itoa(i), "Labels", k)
		}
	}

	// Pods in this Namespace

	pods, err := clientSet.CoreV1().Pods(claims.Kubernetes.Namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		vals <- trie.Insert(k8sFromError(log, err), "Cloud", "Kubernetes", "Namespace", "Pods")
		return
	}
	for i, p := range pods.Items {
		vals <- trie.Insert(trie.Some(p.Name), "Cloud", "Kubernetes", "Namespace", "Pods", strconv.Itoa(i), "Name")
		for j, c := range p.Spec.Containers {
			vals <- trie.Insert(trie.Some(c.Image), "Cloud", "Kubernetes", "Namespace", "Pods", strconv.Itoa(i), "Containers", strconv.Itoa(j), "Image")
		}
	}

	// Containers in this Pod

	hostname, _ := os.Hostname()
	pod, err := clientSet.CoreV1().Pods(claims.Kubernetes.Namespace).Get(ctx, hostname, metav1.GetOptions{})
	if err != nil {
		vals <- trie.Insert(k8sFromError(log, err), "Cloud", "Kubernetes", "Pod")
		return
	}
	for i, c := range pod.Spec.Containers {
		vals <- trie.Insert(trie.Some(c.Name), "Cloud", "Kubernetes", "Pod", "Containers", strconv.Itoa(i), "Name")
		vals <- trie.Insert(trie.Some(c.Image), "Cloud", "Kubernetes", "Pod", "Containers", strconv.Itoa(i), "Image")
	}

	vals <- trie.Insert(trie.Some(pod.Spec.NodeName), "Cloud", "Kubernetes", "Pod", "NodeName")

	// TODO: get own namespace pods list
	// TODO: get nodes list: vals <- trie.Insert(trie.Some(strconv.Itoa(len(nodes.Items))), "Cloud", "Kubernetes", "Cluster", "NodesCount")
	// TODO: search for parent ownerrefs, and build kubectl tree style location
}
func k8sFromError(log logr.Logger, err error) trie.Value {
	if k8sErrors.IsForbidden(err) {
		return trie.Forbidden()
	} else if err == context.DeadlineExceeded {
		return trie.Timeout(time.Second) // TODO: use the actual timeout!
	} else if k8sErrors.IsTimeout(err) { // client-go blew its own deadline?
		return trie.Timeout(time.Second) // TODO: use actual timeout!
		// Is also an IsServerTimeout() to show the apiserver popped its deadline
	} else {
		return trie.Error(fmt.Errorf("unknown kubernetes error: %w", err))
	}
}

// func svcRange(ctx context.Context, k8s kubernetes.Interface) {
// 	svc := &corev1.Service{
// 		ObjectMeta: metav1.ObjectMeta{
// 			Name:      "fake",
// 			Namespace: "kube-system",
// 		},
// 		Spec: corev1.ServiceSpec{
// 			Ports: []corev1.ServicePort{
// 				corev1.ServicePort{
// 					Port: 443,
// 				},
// 			},
// 			ClusterIP: "1.1.1.1",
// 		},
// 	}

// 	_, err := k8s.CoreV1().Services("kube-system").Create(ctx, svc, metav1.CreateOptions{})
// 	var se *k8sErrors.StatusError
// 	if err != nil {
// 		if errors.As(err, &se) {
// 			cause := se.ErrStatus.Details.Causes[0]
// 			if cause.Type == metav1.CauseTypeFieldValueInvalid && cause.Field == "spec.clusterIPs" {
// 				fmt.Println(cause.Message[strings.LastIndex(cause.Message, " ")+1:])
// 				return
// 			}
// 		}
// 	}
// 	panic("Unexpectedly no error, wrong error, etc. Some error with the error.")
// }

func quantity2int64(q *resource.Quantity) int64 {
	x, ok := q.AsInt64()
	if !ok {
		panic("Failed to unwrap k8s resource.Quantity - overflows int64")
	}
	return x
}
func quantity2str(q *resource.Quantity) string {
	n := quantity2int64(q)
	return strconv.FormatInt(n, 10)
}
