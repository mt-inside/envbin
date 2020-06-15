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
            Started at {{.StartTime}}; up {{.RunTime}}<br>
        </p>

        <h1 class="title is-3">Session</h1>
        <p>
            Name: {{.SessionName}}<br>
            Request: {{.RequestNumber}}<br>
        </p>

        <h1 class="title is-3">Host</h1>
        <p>
            OS: {{.OsType}} {{.OsVersion}}, uptime {{.OsUptime}}<br>
            Virtualisation: {{.Virt}}<br>
            Hardware: {{.Arch}}, {{.CpuName}}, {{.PhysCores}}/{{.VirtCores}} cores, {{.MemTotal}} RAM<br>
            Procs: {{.ProcCount}} procs<br>
            Evnironment: PID: {{.Pid}}, U/GID: {{.Uid}}/{{.Gid}}<br>
            Hostname: {{.Hostname}}, Primary IP: {{.Ip}}<br>
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

        <h1 class="title is-3">Resources</h1>
        <p>
            Memory: {{.MemUseVirtual}} virtual, {{.MemUsePhysical}} physical<br>
            GC Runs: {{.GcRuns}}<br>
            CPU Time: {{.CpuSelfTime}}<br>
        </p>

    </div>

</section>

<section class="section">

    <div class="field is-horizontal">
        <div class="field-label is-normal">
            <label class="label">Quit</label>
        </div>
        <div class="field-body">
            <form method="get" action="handlers/exit">
                <div class="field has-addons">
                    <p class="control has-icons-left">
                        <input class="input" type="text" name="code" value="0" placeholder="Exit Code">
                        <span class="icon is-small is-left"><i class="fas fa-skull"></i></span>
                    </p>
                    <p class="control">
                        <button class="button">Exit</button>
                    </p>
                </div>
            </form>
        </div>
    </div>

    <div class="field is-horizontal">
        <div class="field-label is-normal">
            <label class="label">Liveness</label>
        </div>
        <div class="field-body">
            <div class="field has-addons">
                <p class="control">
                    <form method="get" action="handlers/liveness">
                        <button class="button {{if eq .SettingLiveness "true"}}is-primary{{end}}" name="value" value="true">true</button>
                        <button class="button {{if eq .SettingLiveness "false"}}is-primary{{end}}" name="value" value="false">false</button>
                    </form>
                </p>
            </div>
        </div>
    </div>

    <div class="field is-horizontal">
        <div class="field-label is-normal">
            <label class="label">Readiness</label>
        </div>
        <div class="field-body">
            <div class="field has-addons">
                <div class="control">
                    <form method="get" action="handlers/readiness">
                        <button class="button {{if eq .SettingReadiness "true"}}is-primary{{end}}" name="value" value="true">true</button>
                        <button class="button {{if eq .SettingReadiness "false"}}is-primary{{end}}" name="value" value="false">false</button>
                    </form>
                </div>
            </div>
        </div>
    </div>

    <div class="field is-horizontal">
        <div class="field-label is-normal">
            <label class="label">Latency</label>
        </div>
        <div class="field-body">
            <form method="get" action="handlers/delay">
                <div class="field has-addons">
                    <p class="control has-icons-left">
                        <input class="input" type="text" name="value" value="{{.SettingLatency}}" placeholder="s">
                        <span class="icon is-small is-left"><i class="fas fa-user-clock"></i></span>
                    </p>
                    <p class="control">
                        <button class="button">Set</button>
                    </p>
                </div>
            </form>
        </div>
    </div>

    <div class="field is-horizontal">
        <div class="field-label is-normal">
            <label class="label">Bandwidth</label>
        </div>
        <div class="field-body">
            <form method="get" action="handlers/bandwidth">
                <div class="field has-addons">
                    <p class="control has-icons-left">
                        <input class="input" type="text" name="value" value="{{.SettingBandwidth}}" placeholder="bytes/s">
                        <span class="icon is-small is-left"><i class="fas fa-tachometer-alt"></i></span>
                    </p>
                    <p class="control">
                        <button class="button">Set</button>
                    </p>
                </div>
            </form>
        </div>
    </div>

    <div class="field is-horizontal">
        <div class="field-label is-normal">
            <label class="label">Error Rate</label>
        </div>
        <div class="field-body">
            <form method="get" action="handlers/errorrate">
                <div class="field has-addons">
                    <p class="control has-icons-left">
                        <input class="input" type="text" name="value" value="{{.SettingErrorRate}}" placeholder="rate [0-1]">
                        <span class="icon is-small is-left"><i class="fas fa-bomb"></i></span>
                    </p>
                    <p class="control">
                        <button class="button">Set</button>
                    </p>
                </div>
            </form>
        </div>
    </div>

    <div class="field is-horizontal">
        <div class="field-label is-normal">
            <label class="label">Allocate</label>
        </div>
        <div class="field-body">
            <div class="field has-addons">
                <div class="control">
                    <form method="get" action="handlers/allocate">
                        <button class="button" name="value" value="1024">1kB</button>
                        <button class="button" name="value" value="1048576">1MB</button>
                        <button class="button" name="value" value="1073741824">1GB</button>
                    </form>
                </div>
            </div>
        </div>
    </div>

    <div class="field is-horizontal">
        <div class="field-label is-normal">
            <label class="label">CPU Use</label>
        </div>
        <div class="field-body">
            <form method="get" action="handlers/cpu">
                <div class="field has-addons">
                    <p class="control has-icons-left">
                        <input class="input" type="text" name="value" value="{{.SettingCpuUse}}" placeholder="cores [0-n]">
                        <span class="icon is-small is-left"><i class="fas fa-microchip"></i></span>
                    </p>
                    <p class="control">
                        <button class="button">Set</button>
                    </p>
                </div>
            </form>
        </div>
    </div>

</section>

<script src="https://use.fontawesome.com/releases/v5.3.1/js/all.js"></script>

</body>

</html>
