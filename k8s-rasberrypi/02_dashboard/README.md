Launch Kubernetes Dashboard
===========================

管理アカウントを作成する
--------------------

```shell
$ kubectl apply -f service-account-admin.yaml
serviceaccount/admin-user created

$ kubectl apply -f cluster-role-admin.yaml
clusterrolebinding.rbac.authorization.k8s.io/admin-user created

$ kubectl -n kube-system describe secret $(kubectl -n kube-system get secret | grep admin-user | awk '{print $1}')
Name:         admin-user-token-m42k7
Namespace:    kube-system
Labels:       <none>
Annotations:  kubernetes.io/service-account.name: admin-user
              kubernetes.io/service-account.uid: 418f493a-011d-11e9-b480-b827eb24fca9

Type:  kubernetes.io/service-account-token

Data
====
namespace:  11 bytes
token:     <jwt>
ca.crt:     1025 bytes
```

Dashboard をリソース作成
----------------------

ARM 用のイメージを使用する。

```shell
$ kubectl apply -f kubernetes-dashboard-arm.yaml
```