# Rearch

TODO find old notes on this (EN? Notion? GH issues? GH project!)

I think (all in this repo, just packages & cmds)
* envbin (libenv?) as a library, API call that returns the tree
* client binary that makes that call, pretty-prints it. Also with a dump command like the "daemon" currently has (use the http-log indenting stringbuilder?)
* envbin server binary that makes taht call, serves it as JSON etc over an "API"
* badpod server binary that has the :8081 config API / config file/args for setting error-rate, bitrate, etc, and an :8080 API that returns parts of the libenv tree according to a template (forget lorumIpsum, just have a --repeat option to keep repeating the template if we need long output). Should also be able to give tree prefix and get that whole part pretty-printed. No whole-tree as API on this one.
  * find out what container tag the LI course wants, and have that be a special build that replicates that output (or just make that the default template)
  * controlplane for this. Just read config file from http? (means we can drop the config API which I remember being a PITA). So can use bucket / etc? Also support reading config file from k8s configMap

## Code hygene
Tetratelabs group.Run for the daemons & their config
Tetratelabs telemetry + my zap logger + logging best-practice
Pretty-print client to use indenting stringbuilder & http-log's styler?

# Details
* move like, everything, to: https://github.com/osquery/osquery-go
* bus topology (intra-usb, and pci->usb controllers, and intra-pci (domain->bus->device))
  * Build one unified topology tree (each entry has a bus type / device type). This lets clients render it as-is, or walk the tree and filter for certain types to give eg a list of USB devices. They can also filter for type and only print leaf nodes to avoid all those "root hubs"
  * don't think thunderbolt controllers are multiplexers, think they just a) are a device on a bus (lane) and b) have lanes (busses) going through them transparently, which will magically have devices on them if anything's plugged into that TB port
