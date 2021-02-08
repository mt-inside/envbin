<!doctype html>
<html class="no-js" lang="">

<head>
    <meta charset="utf-8">
    <title>Envbin</title>
    <meta name="description" content="">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/bulma@0.8.2/css/bulma.min.css">
    <!--<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/bulma-extensions@6.2.7/bulma-switch/dist/css/bulma-switch.min.css" integrity="sha256-hhNzSX9QCUNRpgKiGuOGzPtUdetKhSP4X/jQkkYgBzI=" crossorigin="anonymous">
    <script src="https://cdn.jsdelivr.net/npm/bulma-extensions@6.2.7/dist/js/bulma-extensions.min.js" integrity="sha256-q4zsxO0fpPm6VhtL/9QkCFE5ZkNa0yeUxhmt1VO1ev0=" crossorigin="anonymous"></script>-->
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/normalize/8.0.1/normalize.min.css">
    <meta name="theme-color" content="#fafafa">
</head>

<body>

<section class="section">
    <div class="container">
        <h1 class="title is-1">envbin</h1>

        <h1 class="title is-3">Version </h1>
        <p>
            Version {{.Version}}, git {{.GitCommit}}<br>
            Built at {{.BuildTime}} with {{.GoVersion}}<br>
        </p>

        <h1 class="title is-3">Request</h1>
        <p>
            Apparent source: {{.RequestIp}} ({{.RequestIpEnrich}})<br>
            Proxies: {{.ProxyChain}}<br>
            User agent: {{.UserAgent}}<br>
        </p>

        <h1 class="title is-3">Hardware</h1>
        <p>
            Virtualisation: {{.Virt}}<br>
            Firmware: {{.FirmwareType}}<br>
            Apparent hardware: {{.Arch}}, {{.CpuName}}, {{.PhysCores}}/{{.VirtCores}} cores, {{.MemTotal}} RAM<br>
        </p>

        <h1 class="title is-3">Operating Environment</h1>
        <p>
            OS: {{.OsType}} {{.OsVersion}}, uptime {{.OsUptime}}<br>
            PID: {{.Pid}}, parent: {{.Ppid}} (#others: {{.OtherProcsCount}})<br>
            UID: {{.Uid}} (effective: {{.Euid}})<br>
            Primary GID: {{.Gid}} (effective: {{.Egid}}) (others: {{.Groups}})<br>
        </p>

        <h1 class="title is-3">Network</h1>
        <p>
            Hostname: {{.Hostname}}, Primary IP: {{.HostIp}}<br>
            Apparent external IP: {{.ExternalIp}} ({{.ExternalIpEnrich}})<br>
        </p>

        <h1 class="title is-3">Kubernetes</h1>
        <p>
            Present: {{.K8s}}<br>
            Control plane version: {{.K8sVersion}}<br>
            Nodes: {{.K8sNodeCount}}<br>
        </p>
        <h1 class="title is-5">This Pod</h1>
        <p>
            Running in namespace: {{.K8sNamespace}}<br>
            As ServiceAccount: {{.K8sServiceAccount}}<br>
            Containers: {{.K8sThisPodContainersImages}}<br>
        </p>
        <h1 class="title is-5">This Node</h1>
        <p>
            Address: {{.K8sNodeAddress}}<br>
            Version: {{.K8sNodeVersion}}<br>
            OS: {{.K8sNodeOS}}<br>
            Container runtime: {{.K8sNodeRuntime}}<br>
            Node role(s): {{.K8sNodeRole}}<br>
            Cloud instance type: {{.K8sNodeCloudInstance}}<br>
            Cloud zone: {{.K8sNodeCloudZone}}<br>
        </p>

    </div>

</section>

<script src="https://use.fontawesome.com/releases/v5.3.1/js/all.js"></script>

</body>

</html>
