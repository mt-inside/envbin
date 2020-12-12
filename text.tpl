                 _     _
  ___ _ ____   _| |__ (_)_ __
 / _ \ '_ \ \ / / '_ \| | '_ \
|  __/ | | \ V /| |_) | | | | |_
 \___|_| |_|\_/ |_.__/|_|_| |_(_)

Version {{.Version}}, git {{.GitCommit}}
Built at {{.BuildTime}} with {{.GoVersion}}

REQUEST
Apparent source: {{.RequestIp}} ({{.RequestIpEnrich}})
User Agent: {{.UserAgent}}

HARDWARE
Virtualisation: {{.Virt}}
Apparent hardware: {{.Arch}}, {{.CpuName}}, {{.PhysCores}}/{{.VirtCores}} cores, {{.MemTotal}} RAM

OPERATING ENVIRONMENT
OS: {{.OsType}} {{.OsVersion}}, uptime {{.OsUptime}}
PID: {{.Pid}}, parent: {{.Ppid}} (#others: {{.OtherProcsCount}})
UID: {{.Uid}} (effective: {{.Euid}})
Primary GID: {{.Gid}} (effective: {{.Egid}}) (others: {{.Groups}})

NETWORK
Hostname: {{.Hostname}}, Primary IP: {{.HostIp}}
Apparent external IP: {{.ExternalIp}} ({{.ExternalIpEnrich}})

KUBERNETES
Present: {{.K8s}}
Version: {{.K8sVersion}}
Running in namespace {{.K8sNamespace}}
As ServiceAccount: {{.K8sServiceAccount}}
This Pod:
  Containers: {{.K8sThisPodContainersCount}}
  Images: {{.K8sThisPodContainersImages}}
Nodes: {{.K8sNodeCount}}
This Node:
  Address: {{.K8sNodeAddress}}
  Version: {{.K8sNodeVersion}}
  OS: {{.K8sNodeOS}}
  Container runtime: {{.K8sNodeRuntime}}
  Node role: {{.K8sNodeRole}}
  Cloud Instance: {{.K8sNodeCloudInstance}}
  Cloud Zone: {{.K8sNodeCloudZone}}
