envbin {{.Version}}, {{.GoVersion}}
Started at {{.StartTime}}, up {{.RunTime}}

Session: {{.SessionName}}
Request: {{.RequestNumber}}

Host: {{.OsType}} {{.OsVersion}}, uptime {{.OsUptime}}
Virtualisation: {{.Virt}}
Kubernetes: {{.K8s}} {{.K8sVersion}} in {{.K8sNamespace}} running as {{.K8sServiceAccount}}
Hardware: {{.Arch}}, {{.CpuName}}, {{.PhysCores}}/{{.VirtCores}} cores, {{.MemTotal}} RAM
Procs: {{.ProcCount}} procs
Evnironment: PID: {{.Pid}}, U/GID: {{.Uid}}/{{.Gid}}
Network: {{.Hostname}}, {{.Ip}}

Memory: {{.MemUseVirtual}} virtual, {{.MemUsePhysical}} physical
GC Runs: {{.GcRuns}}
CPU: time {{.CpuSelfTime}}

Liveness: {{.SettingLiveness}}
Readiness: {{.SettingReadiness}}
Latency: {{.SettingLatency}}
Bandwidth: {{.SettingBandwidth}}
Error rate: {{.SettingErrorRate}}
Cpu Use: {{.SettingCpuUse}}
