<!doctype html>
<html class="no-js" lang="">

<head>
    <meta charset="utf-8">
    <title>Envbin</title>
    <meta name="description" content="">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/bulma@0.8.1/css/bulma.min.css">
    <!--<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/bulma-extensions@6.2.7/bulma-switch/dist/css/bulma-switch.min.css" integrity="sha256-hhNzSX9QCUNRpgKiGuOGzPtUdetKhSP4X/jQkkYgBzI=" crossorigin="anonymous">
    <script src="https://cdn.jsdelivr.net/npm/bulma-extensions@6.2.7/dist/js/bulma-extensions.min.js" integrity="sha256-q4zsxO0fpPm6VhtL/9QkCFE5ZkNa0yeUxhmt1VO1ev0=" crossorigin="anonymous"></script>-->
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/normalize/8.0.1/normalize.min.css">
    <meta name="theme-color" content="#fafafa">
</head>

<body>

<section class="section">
    <div class="container">
        <h1 class="title">envbin</h1>
    </div>
</section>

<section class="section">

    <nav class="panel is-primary">
        <p class="panel-heading">
            Version
        </p>
        <div class="panel-block">
            Version {{.Version}}, git {{.GitCommit}}<br>
            Built at {{.BuildTime}} with {{.GoVersion}}<br>
            Started at {{.StartTime}}; up {{.RunTime}}<br>
        </div>
    </nav>

    <nav class="panel is-primary">
        <p class="panel-heading">
            Session
        </p>
        <div class="panel-block">
            Name: {{.SessionName}}<br>
            Request: {{.RequestNumber}}<br>
        </div>
    </nav>

    <nav class="panel is-primary">
        <p class="panel-heading">
            Host
        </p>
        <div class="panel-block">
            OS: {{.OsType}} {{.OsVersion}}, uptime {{.OsUptime}}<br>
            Virtualisation: {{.Virt}}<br>
            Hardware: {{.Arch}}, {{.CpuName}}, {{.PhysCores}}/{{.VirtCores}} cores, {{.MemTotal}} RAM<br>
            Procs: {{.ProcCount}} procs<br>
            Evnironment: PID: {{.Pid}}, U/GID: {{.Uid}}/{{.Gid}}<br>
            Hostname: {{.Hostname}}, Primary IP: {{.Ip}}<br>
        </div>
    </nav>

    <nav class="panel is-primary">
        <p class="panel-heading">
            Kubernetes
        </p>
        <div class="panel-block">
            Present: {{.K8s}}<br>
            Version: {{.K8sVersion}}<br>
            Running in namespace: {{.K8sNamespace}}<br>
            As ServiceAccount: {{.K8sServiceAccount}}<br>
        </div>
    </nav>

    <nav class="panel is-primary">
        <p class="panel-heading">
            Resources
        </p>
        <div class="panel-block">
            Memory: {{.MemUseVirtual}} virtual, {{.MemUsePhysical}} physical<br>
            GC Runs: {{.GcRuns}}<br>
            CPU Time: {{.CpuSelfTime}}<br>
        </div>
    </nav>

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
