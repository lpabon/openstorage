export KUBECONFIG=/var/run/kubernetes/admin.kubeconfig
cluster/kubectl.sh get nodes

sleep 3
cluster/kubectl.sh create namespace pepsi
cluster/kubectl.sh create namespace coke
cluster/kubectl.sh get namespaces

sleep 3
echo ">> Create secret"
echo -n 'namespace/pepsi' > token
cluster/kubectl.sh -n pepsi create secret generic pxvol --from-file=token

echo ">> Setup NFS CSI driver"
cluster/kubectl.sh -n kube-system create -f csi-nfs.yaml

sleep 15
cluster/kubectl.sh -n kube-system get pods

echo ">> Setup SCs"
cluster/kubectl.sh create -f sc-secrets.yaml

echo ">> Request PVCs"
cluster/kubectl.sh -n pepsi create -f pvc-secrets.yaml

sleep 10
echo ">> Show status of PVCs"
cluster/kubectl.sh -n pepsi get pvc




