Setup MetalLB
==============

## Create Resources

元は https://raw.githubusercontent.com/google/metallb/v0.7.3/manifests/metallb.yaml であるが、
ARM 環境では起動しないので、イメージを書き換えたものを利用する。

```shell
$ kubectl apply -f metallb.yaml
namespace/metallb-system created
serviceaccount/controller created
serviceaccount/speaker created
clusterrole.rbac.authorization.k8s.io/metallb-system:controller created
clusterrole.rbac.authorization.k8s.io/metallb-system:speaker created
role.rbac.authorization.k8s.io/config-watcher created
clusterrolebinding.rbac.authorization.k8s.io/metallb-system:controller created
clusterrolebinding.rbac.authorization.k8s.io/metallb-system:speaker created
rolebinding.rbac.authorization.k8s.io/config-watcher created
daemonset.apps/speaker created
deployment.apps/controller created
```

```shell
$ kubectl get pods -n metallb-system
NAME                          READY   STATUS    RESTARTS   AGE
controller-58f794d964-ggtjq   1/1     Running   0          69s
speaker-7dsh7                 1/1     Running   0          69s
speaker-t6lw6                 1/1     Running   0          69s
```

## Configure MetalLB for BGP

```shell
$ kubectl apply -f metallb-bgp.yaml
```


## Try

```shell
$ kubectl apply -f lb-nginx.yaml
deployment.apps/nginx created
service/nginx created
$ kubectl get svc
NAME         TYPE           CLUSTER-IP      EXTERNAL-IP   PORT(S)        AGE
kubernetes   ClusterIP      10.96.0.1       <none>        443/TCP        6d
nginx        LoadBalancer   10.100.116.53   192.51.100.0   80:31655/TCP   34s

$ curl 192.51.100.0
<!DOCTYPE html>
<html>
<head>
<title>Welcome to nginx!</title>
<style>
    body {
        width: 35em;
        margin: 0 auto;
        font-family: Tahoma, Verdana, Arial, sans-serif;
    }
</style>
</head>
<body>
<h1>Welcome to nginx!</h1>
<p>If you see this page, the nginx web server is successfully installed and
working. Further configuration is required.</p>

<p>For online documentation and support please refer to
<a href="http://nginx.org/">nginx.org</a>.<br/>
Commercial support is available at
<a href="http://nginx.com/">nginx.com</a>.</p>

<p><em>Thank you for using nginx.</em></p>
</body>
</html>
```