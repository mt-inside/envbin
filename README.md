### On builds, linking, and CGO

gousb needs libusb, so we need CGO to call it.
Because we use CGO, we don't get static linking by default.
Trying to statically link is basically a folly across Linux and Darwin, and using Go's native cross-compilation (GH Actions do this when building raw binaries).
Thus we don't try, and will get dynamically-linked binaries.
libusb will need to be present when we cross-compile? Or do we just need headers?
These binaries will need libusb, and the right version of the right libc at runtime, which isn't ideal.

TODO The `native` tag ??
