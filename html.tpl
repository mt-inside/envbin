<html>
<head></head>
<body>
<h1>{{.SessionName}}</h1>
<p>
Request: {{.RequestNumber}}<br>
envbin {{.Version}}, {{.GoVersion}}<br>
Started at {{.StartTime}}, up {{.RunTime}}<br>
</p>

<p>
Host: {{.OsType}} {{.OsVersion}}, uptime {{.OsUptime}}<br>
Virtualisation: {{.Virt}}<br>
Hardware: {{.Arch}}, {{.CpuName}}, {{.PhysCores}}/{{.VirtCores}} cores, {{.MemTotal}} RAM<br>
Procs: {{.ProcCount}} procs<br>
Evnironment: PID: {{.Pid}}, U/GID: {{.Uid}}/{{.Gid}}<br>
Network: {{.Hostname}}, {{.Ip}}<br>
</p>

<p>
Memory: {{.MemUseVirtual}} virtual, {{.MemUsePhysical}} physical<br>
GC Runs: {{.GcRuns}}<br>
CPU: time {{.CpuSelfTime}}<br>
</p>

<table>
<tr>
<td>Quit</td>
<td><form method="get" action="/api/exit">
<input type="text" name="code" value="0">
<button>exit</button>
</form></td>
</tr>
<tr>
<td>Liveness: {{.SettingLiveness}}</td>
<td><form method="get" action="/api/liveness">
<button name="value" value="true">true</button>
<button name="value" value="false">false</button>
</form></td>
</tr>
<tr>
<td>Readiness: {{.SettingReadiness}}</td>
<td><form method="get" action="/api/readiness">
<button name="value" value="true">true</button>
<button name="value" value="false">false</button>
</form></td>
</tr>
<tr>
<td>Latency: {{.SettingLatency}}</td>
<td><form method="get" action="/api/delay">
<input type="text" name="value" value="{{.SettingLatency}}">
<button>set</button>
</form></td>
</tr>
<tr>
<td>Bandwidth: {{.SettingBandwidth}}</td>
<td><form method="get" action="/api/bandwidth">
<input type="text" name="value" value="{{.SettingBandwidth}}">
<button>set</button>
</form></td>
</tr>
<tr>
<td>Error Rate: {{.SettingErrorRate}}</td>
<td><form method="get" action="/api/errorrate">
<input type="text" name="value" value="{{.SettingErrorRate}}">
<button>set</button>
</form></td>
</tr>
<tr>
<td>Allocate:</td>
<td><form method="get" action="/api/allocate">
<button name="value" value="1024">1kB</button>
<button name="value" value="1048576">1MB</button>
<button name="value" value="1073741824">1GB</button>
</form></td>
</tr>
<tr>
<td>CPU Use: {{.SettingCpuUse}}</td>
<td><form method="get" action="/api/cpu">
<input type="text" name="value" value="{{.SettingCpuUse}}">
<button>set</button>
</form></td>
</tr>
</table>
</body>
</html>
