package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/go-logr/logr"
	"github.com/kr/pretty"
	"github.com/mt-inside/envbin/pkg/utils"
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

type NodeBom struct {
	arch  string
	cores int64
	ram   int64
}
type CountedNodeBom struct {
	bom   NodeBom
	count int
}
type NodeDetails struct {
	name string
	bom  NodeBom
}

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

func unwrapQuantity(q *resource.Quantity) int64 {
	x, ok := q.AsInt64()
	if !ok {
		panic("Failed to unwrap k8s resource.Quantity - overflows int64")
	}
	return x
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

	k8sNodes, err := k8s.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
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

	nodes := []NodeDetails{}
	boms := map[NodeBom]int{}

	totalCores := &resource.Quantity{}
	totalRam := &resource.Quantity{}
	totalStorage := &resource.Quantity{}
	totalEphemeral := &resource.Quantity{}
	for _, n := range k8sNodes.Items {
		// TODO: groupByRole & groupByLabel upfront here (array of labels to pre-calc)
		// Rounded becuase nodes often report very similar, but not idential, amounts of RAM. 10 binary places seems about right from some finger-in-the-air
		bom := NodeBom{n.Status.NodeInfo.Architecture, unwrapQuantity(n.Status.Capacity.Cpu()), utils.Round(unwrapQuantity(n.Status.Capacity.Memory()), 2, 10)}
		boms[bom]++

		node := NodeDetails{n.ObjectMeta.Name, bom}
		nodes = append(nodes, node)

		totalCores.Add(*n.Status.Capacity.Cpu())
		totalRam.Add(*n.Status.Capacity.Memory())
		totalStorage.Add(*n.Status.Capacity.Storage())
		totalEphemeral.Add(*n.Status.Capacity.StorageEphemeral())
	}

	fmt.Println("Hardware: ")
	fmt.Printf("  Total %d nodes; %v cores, %v ram, %v storage, %v ephemeral",
		len(nodes),
		totalCores,
		utils.FormatIEC(unwrapQuantity(totalRam), 2),
		utils.FormatIEC(unwrapQuantity(totalStorage), 2),
		utils.FormatIEC(unwrapQuantity(totalEphemeral), 2),
	)
	fmt.Println()

	if false { // all
		for _, node := range nodes {
			fmt.Printf("  %s  %s %3d cores %s ram\n",
				node.name,
				node.bom.arch,
				node.bom.cores,
				utils.FormatIEC(node.bom.ram, 2),
			)
			// TODO show role, any taints/auto-taints, instance/zone/region
		}
	} else {
		boms_sorted := []CountedNodeBom{}
		for node, count := range boms {
			boms_sorted = append(boms_sorted, CountedNodeBom{node, count})
		}
		sort.Slice(boms_sorted, func(i, j int) bool {
			if boms_sorted[i].bom.cores == boms_sorted[j].bom.cores {
				return boms_sorted[i].bom.ram > boms_sorted[j].bom.ram
			}
			return boms_sorted[i].bom.cores > boms_sorted[j].bom.cores
		})
		for _, cn := range boms_sorted {
			fmt.Printf("  %4d  %s %3d cores %s ram\n", cn.count, cn.bom.arch, cn.bom.cores, utils.FormatIEC(cn.bom.ram, 2))
		}

		fmt.Print("Roles: ")
		pretty.Print(groupByRole(k8sNodes))
		fmt.Println()

		// https://kubernetes.io/docs/reference/labels-annotations-taints/
		// All nodes
		render(k8sNodes, "OSen", "kubernetes.io/os")
		render(k8sNodes, "Arches", "kubernetes.io/arch")

		// Auto-taints
		render(k8sNodes, "Not Ready", "node.kubernetes.io/not-ready")
		render(k8sNodes, "Unreachable", "node.kubernetes.io/unreachable")
		render(k8sNodes, "Unschedulable", "node.kubernetes.io/unschedulable")
		render(k8sNodes, "Memory Pressure", "node.kubernetes.io/memory-pressure")
		render(k8sNodes, "Disk Pressure", "node.kubernetes.io/disk-pressure")
		render(k8sNodes, "Network Unavailable", "node.kubernetes.io/network-unavailable")
		render(k8sNodes, "PID Pressure", "node.kubernetes.io/pid-pressure")

		// Cloud instances
		render(k8sNodes, "Instance Types", "node.kubernetes.io/instance-type")
		render(k8sNodes, "Regions", "topology.kubernetes.io/region")
		render(k8sNodes, "Zones", "topology.kubernetes.io/zone")
	}

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
