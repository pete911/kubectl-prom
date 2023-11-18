# kubectl-prom

Kubectl plugin to query prometheus. `kubectl-prom query '<prom-query>'` will send query to prometheus server in the
cluster and returns response.

- no need to manually port forward
- response time can be easily tested by prepending `time` command
- response can be piped to `jq`

Prometheus server is selected by default with `app.kubernetes.io/name=prometheus,app.kubernetes.io/instance=prometheus`
labels (all namespaces are searched). This can be overridden if the prometheus server pods are using different labels
by supplying `-prom-label` and/or `-prom-namespace` flags.

## local testing

- run kind cluster `./e2e/e2e`
- `./kubectl-prom query 'rate(container_cpu_usage_seconds_total{namespace="kube-system",pod="kube-scheduler-prom-test-control-plane"}[2m])' | jq .`
```json
[
  {
    "metric": {
      "beta_kubernetes_io_arch": "arm64",
      "beta_kubernetes_io_os": "linux",
      "container": "kube-scheduler",
      "cpu": "total",
      "id": "/kubelet.slice/kubelet-kubepods.slice/abcdef.scope",
      "image": "registry.k8s.io/kube-scheduler:v1.27.3",
      "instance": "prom-test-control-plane",
      "job": "kubernetes-nodes-cadvisor",
      "kubernetes_io_arch": "arm64",
      "kubernetes_io_hostname": "prom-test-control-plane",
      "kubernetes_io_os": "linux",
      "name": "bf90d5ee9216f36b5be5006957d9434095748c66af572d84b449d82ef9c8c3df",
      "namespace": "kube-system",
      "pod": "kube-scheduler-prom-test-control-plane"
    },
    "value": [
      1700302921.619,
      "0.0030641861415124844"
    ]
  }
]
```
- when finished testing, delete cluster `kind delete cluster --name prom-test`
