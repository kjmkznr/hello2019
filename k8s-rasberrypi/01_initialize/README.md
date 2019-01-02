Kubernetes Cluster on Raspberry Pi 3 Model B+
=============================================

1. Setup Raspberry Pi SD Card Image
------------------------------------

```shell
$ fdisk -lu 2018-11-13-raspbian-stretch-lite.img
Disk 2018-11-13-raspbian-stretch-lite.img: 1.8 GiB, 1866465280 bytes, 3645440 sectors
Units: sectors of 1 * 512 = 512 bytes
Sector size (logical/physical): 512 bytes / 512 bytes
I/O size (minimum/optimal): 512 bytes / 512 bytes
Disklabel type: dos
Disk identifier: 0x7ee80803

Device                                Boot Start     End Sectors  Size Id Type
2018-11-13-raspbian-stretch-lite.img1       8192   98045   89854 43.9M  c W95 FAT32 (LBA)
2018-11-13-raspbian-stretch-lite.img2      98304 3645439 3547136  1.7G 83 Linux
$ sudo mount -o loop,offset=$(( 512 * 8192 )) 2018-11-13-raspbian-stretch-lite.img /mnt/boot
$ sudo touch /mnt/boot/ssh
$ sudo nvim /mnt/boot/cmdline.txt
( cgroup_enable=cpuset cgroup_enable=memory cgroup_memory=1 を行の最後に追加 )
```

2. Upgrade System
------------------

SDカードを焼いたあと、RaspberryPi を起動する。

```shell
$ sudo apt update
$ sudo apt upgrade
```

3. Disable Swap
---------------

Swap を有効にしていると Kubernetes 起動時に警告がでるため。

```shell
$ sudo dphys-swapfile swapoff
$ sudo dphys-swapfile uninstall
$ sudo update-rc.d dphys-swapfile remove
```

4. Install Docker
-----------------

ここで入れるのは stable 版

```shell
$ sudo apt install gnupg2 software-properties-common
$ curl -fsSL https://download.docker.com/linux/$(. /etc/os-release; echo "$ID")/gpg | sudo apt-key add -
$ echo "deb [arch=armhf] https://download.docker.com/linux/$(. /etc/os-release; echo "$ID") \
    $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list
$ sudo apt update
$ sudo apt install docke-ce
```

5. Install kubeadm
------------------

```shell
$ curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | sudo apt-key add -
$ cat <<EOF | sudo tee /etc/apt/sources.list.d/kubernetes.list
deb https://apt.kubernetes.io/ kubernetes-xenial main
EOF
$ sudo apt update
$ sudo apt install -y kubelet kubeadm kubectl
$ sudo apt-mark hold kubelet kubeadm kubectl
```

6. Setup Kubernetes Cluster
----------------------------

### On master node

```shell
$ sudo kubeadm init --pod-network-cidr=10.244.0.0/16
[init] Using Kubernetes version: v1.13.1
[preflight] Running pre-flight checks
        [WARNING SystemVerification]: this Docker version is not on the list of validated versions: 18.09.0. Latest validated version: 18.06
[preflight] Pulling images required for setting up a Kubernetes cluster
[preflight] This might take a minute or two, depending on the speed of your internet connection
[preflight] You can also perform this action in beforehand using 'kubeadm config images pull'
[kubelet-start] Writing kubelet environment file with flags to file "/var/lib/kubelet/kubeadm-flags.env"
[kubelet-start] Writing kubelet configuration to file "/var/lib/kubelet/config.yaml"
[kubelet-start] Activating the kubelet service
[certs] Using certificateDir folder "/etc/kubernetes/pki"
[certs] Generating "ca" certificate and key
[certs] Generating "apiserver-kubelet-client" certificate and key
[certs] Generating "apiserver" certificate and key
[certs] apiserver serving cert is signed for DNS names [k8s kubernetes kubernetes.default kubernetes.default.svc kubernetes.default.svc.cluster.local] and IPs [10.96.0.1 192.168.5.119]
[certs] Generating "front-proxy-ca" certificate and key
[certs] Generating "front-proxy-client" certificate and key
[certs] Generating "etcd/ca" certificate and key
[certs] Generating "etcd/peer" certificate and key
[certs] etcd/peer serving cert is signed for DNS names [k8s localhost] and IPs [192.168.5.119 127.0.0.1 ::1]
[certs] Generating "apiserver-etcd-client" certificate and key
[certs] Generating "etcd/server" certificate and key
[certs] etcd/server serving cert is signed for DNS names [k8s localhost] and IPs [192.168.5.119 127.0.0.1 ::1]
[certs] Generating "etcd/healthcheck-client" certificate and key
[certs] Generating "sa" key and public key
[kubeconfig] Using kubeconfig folder "/etc/kubernetes"
[kubeconfig] Writing "admin.conf" kubeconfig file
[kubeconfig] Writing "kubelet.conf" kubeconfig file
[kubeconfig] Writing "controller-manager.conf" kubeconfig file
[kubeconfig] Writing "scheduler.conf" kubeconfig file
[control-plane] Using manifest folder "/etc/kubernetes/manifests"
[control-plane] Creating static Pod manifest for "kube-apiserver"
[control-plane] Creating static Pod manifest for "kube-controller-manager"
[control-plane] Creating static Pod manifest for "kube-scheduler"
[etcd] Creating static Pod manifest for local etcd in "/etc/kubernetes/manifests"
[wait-control-plane] Waiting for the kubelet to boot up the control plane as static Pods from directory "/etc/kubernetes/manifests". This can take up to 4m0s
[kubelet-check] Initial timeout of 40s passed.
[apiclient] All control plane components are healthy after 81.508702 seconds
[uploadconfig] storing the configuration used in ConfigMap "kubeadm-config" in the "kube-system" Namespace
[kubelet] Creating a ConfigMap "kubelet-config-1.13" in namespace kube-system with the configuration for the kubelets in the cluster
[patchnode] Uploading the CRI Socket information "/var/run/dockershim.sock" to the Node API object "k8s" as an annotation
[mark-control-plane] Marking the node k8s as control-plane by adding the label "node-role.kubernetes.io/master=''"
[mark-control-plane] Marking the node k8s as control-plane by adding the taints [node-role.kubernetes.io/master:NoSchedule]
[bootstrap-token] Using token: kz4kug.ns803n2fjcs1kfzy
[bootstrap-token] Configuring bootstrap tokens, cluster-info ConfigMap, RBAC Roles
[bootstraptoken] configured RBAC rules to allow Node Bootstrap tokens to post CSRs in order for nodes to get long term certificate credentials
[bootstraptoken] configured RBAC rules to allow the csrapprover controller automatically approve CSRs from a Node Bootstrap Token
[bootstraptoken] configured RBAC rules to allow certificate rotation for all node client certificates in the cluster
[bootstraptoken] creating the "cluster-info" ConfigMap in the "kube-public" namespace
[addons] Applied essential addon: CoreDNS
[addons] Applied essential addon: kube-proxy

Your Kubernetes master has initialized successfully!

To start using your cluster, you need to run the following as a regular user:

  mkdir -p $HOME/.kube
  sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
  sudo chown $(id -u):$(id -g) $HOME/.kube/config

You should now deploy a pod network to the cluster.
Run "kubectl apply -f [podnetwork].yaml" with one of the options listed at:
  https://kubernetes.io/docs/concepts/cluster-administration/addons/

You can now join any number of machines by running the following on each node
as root:

  kubeadm join 192.168.5.119:6443 --token xxxx.xxxx --discovery-token-ca-cert-hash sha256:97b11801a4e40632de54e28e5389be8c63ebb7f1a4be01b9f5bebc645bbc64b4

$ mkdir -p $HOME/.kube
$ sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
$ sudo chown $(id -u):$(id -g) $HOME/.kube/config
```

7. Setup flannel
-----------------

### On master node

```shell
$ kubectl apply -f kube-flannel.yml
```

8. Setup kubernetes node
-------------------------

```shell
$ kubeadm join 192.168.5.119:6443 --token kz4kug.ns803n2fjcs1kfzy --discovery-token-ca-cert-hash sha256:97b11801a4e40632de54e28e5389be8c63ebb7f1a4be01b9f5bebc645bbc64b4
```

```shell
$ kubectl get nodes
NAME         STATUS     ROLES    AGE     VERSION
k8s          NotReady   master   19m     v1.13.1
k8s-node02   NotReady   <none>   9m52s   v1.13.1
k8s-node03   NotReady   <none>   9m31s   v1.13.1
```