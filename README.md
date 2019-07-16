# :ledger: logbook

Logbook is a real-time log viewer for Kubernetes.

![Screenshot](screenshot.gif)

## Install

Logbook requires Go compiler:

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
- <kbd>f</kbd>: Enable and disable follow mode
- <kbd>Ctrl</kbd>+<kbd>D</kbd>: Scroll half-page down
- <kbd>Ctrl</kbd>+<kbd>U</kbd>: Scroll half-page up
- <kbd>Ctrl</kbd>+<kbd>F</kbd>: Scroll page down
- <kbd>Ctrl</kbd>+<kbd>B</kbd>: Scroll page up
- <kbd>G</kbd>: Scroll to bottom
- <kbd>g</kbd>: Scroll to top
- <kbd>Tab</kbd>: Switch containers
- <kbd>q</kbd>: Quit

## License

MIT
