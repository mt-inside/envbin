// TODO: read namespace labels and annotations
// TODO: read pod labels and annotions
// TODO: can find the Service(s)? Find endpoint with same IP, that has a ref to the service?

package fetchers

import (
	"context"
	"errors"
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
	. "github.com/mt-inside/envbin/pkg/data/trie"
)

func init() {
	data.RegisterPlugin(getK8sDataInCluster)
}

// inCluster wraps generic getter
// cmdline to call generic getting with its own config (clientSet?)
// both to return a trie
// * cmdline thing to share a renderer with "dump"

func getK8sDataInCluster(ctx context.Context, log logr.Logger, vals chan<- InsertMsg) {
	config, err := rest.InClusterConfig()
	if err != nil {
		vals <- Insert(NotPresent(), "Cloud", "Kubernetes")
		return
	}

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		vals <- Insert(Error(err), "Cloud", "Kubernetes")
		log.Error(err, "Can't connect to k8s apiserver")
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

func getK8sData(ctx context.Context, log logr.Logger, clientSet *kubernetes.Clientset, vals chan<- InsertMsg) {
	// Control Plane
	version, err := clientSet.Discovery().ServerVersion()
	if err != nil {
		vals <- Insert(k8sFromError(log, err), "Cloud", "Kubernetes", "Masters")
		log.Error(err, "Can't get cluster version")
	} else {
		vals <- Insert(Some(version.GitVersion), "Cloud", "Kubernetes", "Masters", "Version")
		vals <- Insert(Some(version.Platform), "Cloud", "Kubernetes", "Masters", "Platform")
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
	//vals <- Insert(k8sFromError(log, err), "Cloud", "Kubernetes", "Identity")
	//return
	// }
	claims, ok := token.Claims.(*k8sClaims)
	if !ok {
		vals <- Insert(k8sFromError(log, err), "Cloud", "Kubernetes", "Identity")
		return
	}
	vals <- Insert(Some(claims.Kubernetes.Namespace), "Cloud", "Kubernetes", "Namespace", "Name")
	vals <- Insert(Some(claims.Kubernetes.ServiceAccount.Name), "Cloud", "Kubernetes", "Pod", "ServiceAccount")

	// Nodes

	nodes, err := clientSet.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		vals <- Insert(k8sFromError(log, err), "Cloud", "Kubernetes", "Nodes")
		// Keep going becuase we might have permission to read other stuff
	}
	for i, n := range nodes.Items {
		vals <- Insert(Some(n.Status.NodeInfo.Architecture), "Cloud", "Kubernetes", "Nodes", strconv.Itoa(i), "Architecture")
		vals <- Insert(Some(quantity2str(n.Status.Capacity.Cpu())), "Cloud", "Kubernetes", "Nodes", strconv.Itoa(i), "Cores")
		vals <- Insert(Some(quantity2str(n.Status.Capacity.Memory())), "Cloud", "Kubernetes", "Nodes", strconv.Itoa(i), "RAM")
		vals <- Insert(Some(quantity2str(n.Status.Capacity.Storage())), "Cloud", "Kubernetes", "Nodes", strconv.Itoa(i), "Storage")
		vals <- Insert(Some(quantity2str(n.Status.Capacity.StorageEphemeral())), "Cloud", "Kubernetes", "Nodes", strconv.Itoa(i), "Ephemeral")
		for k, v := range n.Labels {
			vals <- Insert(Some(v), "Cloud", "Kubernetes", "Nodes", strconv.Itoa(i), "Labels", k)
		}
	}

	// Pods in this Namespace

	pods, err := clientSet.CoreV1().Pods(claims.Kubernetes.Namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		vals <- Insert(k8sFromError(log, err), "Cloud", "Kubernetes", "Namespace", "Pods")
		return
	}
	for i, p := range pods.Items {
		vals <- Insert(Some(p.Name), "Cloud", "Kubernetes", "Namespace", "Pods", strconv.Itoa(i), "Name")
		for j, c := range p.Spec.Containers {
			vals <- Insert(Some(c.Image), "Cloud", "Kubernetes", "Namespace", "Pods", strconv.Itoa(i), "Containers", strconv.Itoa(j), "Image")
		}
	}

	// Containers in this Pod

	hostname, _ := os.Hostname()
	pod, err := clientSet.CoreV1().Pods(claims.Kubernetes.Namespace).Get(ctx, hostname, metav1.GetOptions{})
	if err != nil {
		vals <- Insert(k8sFromError(log, err), "Cloud", "Kubernetes", "Pod")
		return
	}
	for i, c := range pod.Spec.Containers {
		vals <- Insert(Some(c.Name), "Cloud", "Kubernetes", "Pod", "Containers", strconv.Itoa(i), "Name")
		vals <- Insert(Some(c.Image), "Cloud", "Kubernetes", "Pod", "Containers", strconv.Itoa(i), "Image")
	}

	vals <- Insert(Some(pod.Spec.NodeName), "Cloud", "Kubernetes", "Pod", "NodeName")

	// TODO: get own namespace pods list
	// TODO: get nodes list: vals <- Insert(Some(strconv.Itoa(len(nodes.Items))), "Cloud", "Kubernetes", "Cluster", "NodesCount")
	// TODO: search for parent ownerrefs, and build kubectl tree style location
}
func k8sFromError(log logr.Logger, err error) Value {
	if k8sErrors.IsForbidden(err) {
		log.Error(err, "Kubernetes error: Forbidden; check RBAC")
		return Forbidden()
	} else if err == context.DeadlineExceeded {
		log.Error(err, "Kubernetes error: Timed out")
		return Timeout(time.Second) // TODO: use the actual timeout!
	} else if k8sErrors.IsTimeout(err) { // client-go blew its own deadline?
		log.Error(err, "Kubernetes error: Timed out")
		return Timeout(time.Second) // TODO: use actual timeout!
		// Is also an IsServerTimeout() to show the apiserver popped its deadline
	} else {
		log.Error(err, "Kubernetes error: Other")
		return Error(errors.New("Unknown kubernetes error"))
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
