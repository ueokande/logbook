# :ledger: logbook

Logbook is a real-time log viewer for Kubernetes.

![Screenshot](screenshot.gif)

## Install

Download a latest version of the binary from releases, or use `go get` as follows:

```console
$ go get -u github.com/ueokande/logbook
```

## Usage

```console
$ logbook [--kubeconfig KUBECONFIG] [--namespace NAMESPACE]

Flags:
  --kubeconfig  Path to kubeconfig file
  --namespace   Kubernetes namespace
```

- <kbd>Ctrl</kbd>+<kbd>n</kbd>: Select next pod
- <kbd>Ctrl</kbd>+<kbd>p</kbd>: Scroll previous pod
- <kbd>j</kbd>: Scroll down
- <kbd>k</kbd>: Scroll up
- <kbd>h</kbd>: Scroll left
- <kbd>l</kbd>: Scroll right
- <kbd>f</kbd>: Enable and disable follow mode
- <kbd>Ctrl</kbd>+<kbd>D</kbd>: Scroll half-page down
- <kbd>Ctrl</kbd>+<kbd>U</kbd>: Scroll half-page up
- <kbd>Ctrl</kbd>+<kbd>F</kbd>: Scroll page down
- <kbd>Ctrl</kbd>+<kbd>B</kbd>: Scroll page up
- <kbd>G</kbd>: Scroll to bottom
- <kbd>g</kbd>: Scroll to top
- <kbd>Tab</kbd>: Switch containers
- <kbd>/</kbd>: Search forward for matching line.
- <kbd>n</kbd>: Repeat previous search.
- <kbd>N</kbd>: Repeat previous search in reverse direction.
- <kbd>q</kbd>: Quit

## License

MIT

[releases]: https://github.com/ueokande/logbook/releases
