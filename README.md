mackerel-plugin-gearman [![Build Status](https://travis-ci.org/hfm/mackerel-plugin-gearman.svg?branch=master)](https://travis-ci.org/hfm/mackerel-plugin-gearman)
===

[Gearman](http://gearman.org/) custom metrics plugin for mackerel.io agent.

Synopsis
---

```sh
mackerel-plugin-gearman [-host=<host>] [-port=<port>] [-tempfile=<tempfile>] [-version]
```

```console
$ mackerel-plugin-gearman -h
Usage of mackerel-plugin-gearman:

  -H, -host=127.0.0.1                           Host of gearmand
  -p, -port=4730                                Port of gearmand
  -t, -tempfile=/tmp/mackerel-plugin-gearman    Temp file name
  -v, -version                                  Print version information and quit.
```

Example of mackerel-agent.conf

```toml
[plugin.metrics.gearman]
command = "/path/to/mackerel-plugin-gearman"
```
