# Installing Prometheus on the Kubernetes Cluster using kube-prometheus

### Start minikube 
```
minikube delete && \
minikube start \
    --kubernetes-version=v1.25.2 \
    --memory=6g \
    --bootstrapper=kubeadm \
    --extra-config=kubelet.authentication-token-webhook=true \
    --extra-config=kubelet.authorization-mode=Webhook \
    --extra-config=scheduler.bind-address=0.0.0.0 \
    --extra-config=controller-manager.bind-address=0.0.0.0 &&
minikube addons disable metrics-server
```


### kube-prometheus Compatability with Kubernetes 
| kube-prometheus stack                                                                      | Kubernetes 1.21 | Kubernetes 1.22 | Kubernetes 1.23 | Kubernetes 1.24 | Kubernetes 1.25 | Kubernetes 1.26 | Kubernetes 1.27 |
|--------------------------------------------------------------------------------------------|-----------------|-----------------|-----------------|-----------------|-----------------|-----------------|-----------------|
| [`release-0.9`](https://github.com/prometheus-operator/kube-prometheus/tree/release-0.9)   | ✔               | ✔               | ✗               | ✗               | ✗               | x               | x               |
| [`release-0.10`](https://github.com/prometheus-operator/kube-prometheus/tree/release-0.10) | ✗               | ✔               | ✔               | ✗               | ✗               | x               | x               |
| [`release-0.11`](https://github.com/prometheus-operator/kube-prometheus/tree/release-0.11) | ✗               | ✗               | ✔               | ✔               | ✗               | x               | x               |
| [`release-0.12`](https://github.com/prometheus-operator/kube-prometheus/tree/release-0.12) | ✗               | ✗               | ✗               | ✔               | ✔               | x               | x               |
| [`main`](https://github.com/prometheus-operator/kube-prometheus/tree/main)                 | ✗               | ✗               | ✗               | ✗               | x               | ✔               | ✔               |

source - https://github.com/prometheus-operator/kube-prometheus/tree/main#compatibility


### Clone kube-prometheus Repository 
```
git clone https://github.com/prometheus-operator/kube-prometheus.git
cd kube-prometheus
git checkout release-0.12
``` 

### Create monitoring stack 

```
$ kubectl apply --server-side -f manifests/setup
$ kubectl wait \
	--for condition=Established \
	--all CustomResourceDefinition \
	--namespace=monitoring
$ kubectl apply -f manifests/
```

### Let's explore the stack
```
kubectl --namespace monitoring port-forward svc/prometheus-k8s 9090
kubectl --namespace monitoring port-forward svc/alertmanager-main 9093
kubectl --namespace monitoring port-forward svc/grafana 3000
```

### Simple usage of Prometheus and Grafana Dashboards
```
# Create a Deployment with 5 replicas

$ kubectl apply -f https://raw.githubusercontent.com/kubernetes/website/main/content/en/examples/service/load-balancer-example.yaml

# Create a StatefulSet

$ kubectl apply -f https://raw.githubusercontent.com/kubernetes/website/main/content/en/examples/application/web/web.yaml
```

