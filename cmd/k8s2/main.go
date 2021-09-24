package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/go-logr/logr"
	"github.com/kr/pretty"
	"github.com/mt-inside/go-usvc"
	"github.com/urfave/cli/v2"
	corev1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type Scale2 int32

const (
	//Kibi Scale2 = 10
	//Mebi Scale2 = 20
	Gibi Scale2 = 30
	//Tebi Scale2 = 40
)

func init() {
	spew.Config.DisableMethods = true
	spew.Config.DisablePointerMethods = true
}

func main() {
	log := usvc.GetLogger(true)

	kubeConfigPath := ""
	if home := homedir.HomeDir(); home != "" {
		kubeConfigPath = filepath.Join(home, ".kube", "config")
		if _, err := os.Stat(kubeConfigPath); err != nil {
			kubeConfigPath = ""
		}
	}

	app := cli.App{
		Name:  "kubeinspect",
		Usage: "Dump info about a k8s cluster",

		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "kubeconfig",
				Usage: "Absolute path to kubeconfig file",
				Value: kubeConfigPath,
			},
			&cli.StringFlag{
				Name:  "api",
				Usage: "URL of the Kubernetes API. Overrides any value in the kubeconfig.",
			},
		},

		Action: appMain,

		Metadata: map[string]interface{}{
			"log": log,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Error(err, "Error during exeuction")
	}
}

func appMain(c *cli.Context) error {
	log := c.App.Metadata["log"].(logr.Logger)

	k8s, err := getClientSet(log, c.String("kubeconfig"), c.String("url"))
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	svcRange(ctx, k8s)

	nodes, err := k8s.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return err
	}
	// _, err = clientset.CoreV1().Pods("default").Get(context.TODO(), "example-xxxxx", metav1.GetOptions{})
	// if errors.IsNotFound(err) {
	// 	fmt.Printf("Pod example-xxxxx not found in default namespace\n")
	// } else if statusError, isStatus := err.(*errors.StatusError); isStatus {
	// 	fmt.Printf("Error getting pod %v\n", statusError.ErrStatus.Message)
	// } else if err != nil {
	// 	panic(err.Error())
	// } else {
	// 	fmt.Printf("Found example-xxxxx pod in default namespace\n")
	// }

	version, err := k8s.Discovery().ServerVersion()
	if err != nil {
		return fmt.Errorf("Can't get cluster version: %w", err)
	}
	fmt.Printf("Kubernetes %s %s", version.GitVersion, version.Platform)
	fmt.Println()

	cores := &resource.Quantity{}
	ram := &resource.Quantity{}
	storage := &resource.Quantity{}
	ephemeral := &resource.Quantity{}
	for _, n := range nodes.Items {
		cores.Add(*n.Status.Capacity.Cpu())
		ram.Add(*n.Status.Capacity.Memory())
		storage.Add(*n.Status.Capacity.Storage())
		ephemeral.Add(*n.Status.Capacity.StorageEphemeral())
	}
	pow2Truncate(ram, Gibi)
	pow2Truncate(storage, Gibi)
	pow2Truncate(ephemeral, Gibi)
	fmt.Printf("%d Nodes (%v cores, %v ram, %v storage, %v ephemeral)", len(nodes.Items), cores, ram, storage, ephemeral)
	fmt.Println()

	// node-role.kubernetes.io/control-plane":""
	// "kubernetes.io/os":"linux"
	// "node.kubernetes.io/exclude-from-external-load-balancers":""
	// "kubernetes.io/arch":"amd64"
	// "node-role.kubernetes.io/master":""

	fmt.Print("Roles: ")
	pretty.Print(groupByRole(nodes))
	fmt.Println()

	// https://kubernetes.io/docs/reference/labels-annotations-taints/
	// All nodes
	render(nodes, "OSen", "kubernetes.io/os")
	render(nodes, "Arches", "kubernetes.io/arch")

	// Auto-taints
	render(nodes, "Not Ready", "node.kubernetes.io/not-ready")
	render(nodes, "Unreachable", "node.kubernetes.io/unreachable")
	render(nodes, "Unschedulable", "node.kubernetes.io/unschedulable")
	render(nodes, "Memory Pressure", "node.kubernetes.io/memory-pressure")
	render(nodes, "Disk Pressure", "node.kubernetes.io/disk-pressure")
	render(nodes, "Network Unavailable", "node.kubernetes.io/network-unavailable")
	render(nodes, "PID Pressure", "node.kubernetes.io/pid-pressure")

	// Cloud instances
	render(nodes, "Instance Types", "node.kubernetes.io/instance-type")
	render(nodes, "Regions", "topology.kubernetes.io/region")
	render(nodes, "Zones", "topology.kubernetes.io/zone")

	return nil
}

func svcRange(ctx context.Context, k8s kubernetes.Interface) {
	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "fake",
			Namespace: "kube-system",
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				corev1.ServicePort{
					Port: 443,
				},
			},
			ClusterIP: "1.1.1.1",
		},
	}

	_, err := k8s.CoreV1().Services("kube-system").Create(ctx, svc, metav1.CreateOptions{})
	var se *kerrors.StatusError
	if err != nil {
		if errors.As(err, &se) {
			cause := se.ErrStatus.Details.Causes[0]
			if cause.Type == metav1.CauseTypeFieldValueInvalid && cause.Field == "spec.clusterIPs" {
				fmt.Println(cause.Message[strings.LastIndex(cause.Message, " ")+1:])
				return
			}
		}
	}
	panic("Unexpectedly no error, wrong error, etc. Some error with the error.")
}

// k8s.io/apimachinery/pkg/api/resource.Quantity only knows how to round to decimal powers (despite nicely printing round numbers of binary powers)
func pow2Truncate(q *resource.Quantity, s Scale2) {
	val := q.Value()
	val = val &^ ((1 << s) - 1)
	q.Set(val)
}

func render(nodes *corev1.NodeList, title string, label string) {
	m := groupByLabel(nodes, label)

	if len(m) != 0 {
		fmt.Printf("%s: ", title)
		pretty.Print(m)
		fmt.Println()
	}
}

func groupByRole(ns *corev1.NodeList) (ret map[string]int) {
	ret = make(map[string]int)

	for _, n := range ns.Items {
		roles := getSufficies(n.Labels, "node-role.kubernetes.io/")
		for _, role := range roles {
			ret[role] = ret[role] + 1
		}
	}

	return
}

func getSufficies(ss map[string]string, pre string) []string {
	ret := []string{}

	for k := range ss {
		if strings.HasPrefix(k, pre) {
			ret = append(ret, strings.TrimPrefix(k, pre))
		}
	}
	return ret
}

func groupByLabel(ns *corev1.NodeList, key string) (ret map[string]int) {
	ret = make(map[string]int)

	for _, n := range ns.Items {
		val := n.Labels[key]
		if val != "" {
			ret[val] = ret[val] + 1
		}
	}

	return
}

func getClientSet(log logr.Logger, kubeConfigPath string, masterURL string) (kubernetes.Interface, error) {
	cfg, err := clientcmd.BuildConfigFromFlags(masterURL, kubeConfigPath)
	if err != nil {
		return nil, fmt.Errorf("Error getting kubeconfig: %w", err)
	}

	kubeClientSet, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("Error building kubernetes clientset: %w", err)
	}

	return kubeClientSet, nil
}
