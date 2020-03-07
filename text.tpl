== {{.SessionName}} ==
Request: {{.RequestNumber}}
envbin {{.Version}}, {{.GoVersion}}
Started at {{.StartTime}}, up {{.RunTime}}

Host: {{.OsType}} {{.OsVersion}}, uptime {{.OsUptime}}
Virtualisation: {{.Virt}}
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
