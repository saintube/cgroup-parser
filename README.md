# Cgroup Parser

cgroup-parser is a commandline tool writing for parsing linux cgroup files.

e.g. parse `cpuacct.stat` to get realtime cpu usages.

```bash
$ # get realtime cpu usage of k8s pod xxxx with collecting interval 1000 milli-seconds
$ cgroup-parser cpuacct -p /sys/fs/cgroup/cpuacct/kubepods-slices/podxxxx/cpuacct.stat --interval 1000
```
