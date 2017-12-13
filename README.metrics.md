= Accessing cluster metrics

make sure to have a kubernetes cluster version >= 1.8. [minikube](https://github.com/kubernetes/minikube) is good for experimentation.

== Accessing core metrics

Those are the same metrics that tools like `kubectl top` or the autoscaler use.

1. make sure the metrics server is deployed and running. On minikube, is not active by default.
   To deploy metrics-server, follow [those instructions](https://github.com/kubernetes-incubator/metrics-server/blob/master/README.md)
2. metrics API is discoverable and accessibile like the standard k8s APIs. For example, using the kube-proxy:
```
curl http://127.0.0.1:8080/apis/metrics.k8s.io/v1beta1/
{
  "kind": "APIResourceList",
  "apiVersion": "v1",
  "groupVersion": "metrics.k8s.io/v1beta1",
  "resources": [
    {
      "name": "nodes",
      "singularName": "",
      "namespaced": false,
      "kind": "NodeMetrics",
      "verbs": [
        "get",
        "list"
      ]
    },
    {
      "name": "pods",
      "singularName": "",
      "namespaced": true,
      "kind": "PodMetrics",
      "verbs": [
        "get",
        "list"
      ]
    }
  ]
}
```
