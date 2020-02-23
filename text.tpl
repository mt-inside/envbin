envbin {{.Version}}, {{.GoVersion}}
Started at {{.StartTime}}

Host: {{.OsType}} {{.OsVersion}}, uptime {{.OsUptime}}
Evnironment: PID: {{.Pid}}, U/GID: {{.Uid}}/{{.Gid}}
Hardware: {{.Arch}}, {{.CpuName}}, {{.PhysCores}}/{{.VirtCores}} cores, {{.MemTotal}} RAM
Procs: {{.ProcCount}} procs
Network: {{.Hostname}}, {{.Ip}}

Memory: {{.MemUseVirtual}} virtual, {{.MemUsePhysical}} physical
CPU: time {{.CpuSelfTime}}

Health: {{.SettingHealth}}
Liveness: {{.SettingLiveness}}
Latency: {{.SettingLatency}}
Bandwidth: {{.SettingBandwidth}}
