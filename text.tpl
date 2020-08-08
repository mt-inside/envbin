                 _     _
  ___ _ ____   _| |__ (_)_ __
 / _ \ '_ \ \ / / '_ \| | '_ \
|  __/ | | \ V /| |_) | | | | |_
 \___|_| |_|\_/ |_.__/|_|_| |_(_)

Version {{.Version}}, git {{.GitCommit}}
Built at {{.BuildTime}} with {{.GoVersion}}
Started at {{.StartTime}}; up {{.RunTime}}

SESSION
Name: {{.SessionName}}
Request: {{.RequestNumber}}

YOU
IP: {{.RequestIP}}
User Agent: {{.UserAgent}}

HOST
OS: {{.OsType}} {{.OsVersion}}, uptime {{.OsUptime}}
Virtualisation: {{.Virt}}
Hardware: {{.Arch}}, {{.CpuName}}, {{.PhysCores}}/{{.VirtCores}} cores, {{.MemTotal}} RAM
Procs: {{.ProcCount}} procs
Evnironment: PID: {{.Pid}}, U/GID: {{.Uid}}/{{.Gid}}
Hostname: {{.Hostname}}, primary IP: {{.Ip}}

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

RESOURCES
Memory: {{.MemUseVirtual}} virtual, {{.MemUsePhysical}} physical
GC Runs: {{.GcRuns}}
CPU Time: {{.CpuSelfTime}}

SETTINGS
Liveness: {{.SettingLiveness}}
Readiness: {{.SettingReadiness}}
Latency: {{.SettingLatency}}
Bandwidth: {{.SettingBandwidth}}
Error rate: {{.SettingErrorRate}}
Cpu Use: {{.SettingCpuUse}}
