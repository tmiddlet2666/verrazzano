#!/bin/bash

compartment_id=$1
region=$2
cluster_name=$3
bastion_name=$4
public_key_file=$5
private_key_file=$6
KUBECONFIG=$7
port=$8
echo "CREATE KUBECONFIG at $KUBECONFIG"

#oci bastion bastion create --bastion-type STANDARD --compartment-id $compartment_id --target-subnet-id $target_subnet_id --client-cidr-list '["0.0.0.0/0"]' --max-session-ttl 10800 --name $bastion_name
#exit 0

rm $KUBECONFIG
oci ce cluster delete --force --cluster-id=ocid1.cluster.oc1.ap-tokyo-1.aaaaaaaaq2qg35j4h7tvzronojvfqrbg524o5be67uv4ncc64c4j3use3j7q
oci ce cluster delete --force --cluster-id=ocid1.cluster.oc1.ap-tokyo-1.aaaaaaaaepzv2u5zdingui5xkmbzvdegtnwaoybaxxrr6nixwcc6zhhx7d3a
oci ce cluster delete --force --cluster-id=ocid1.cluster.oc1.ap-tokyo-1.aaaaaaaaxtriif7aygae5ou6zfsmzdvmffufsxndlc7coensccyo6livb44q
oci ce cluster delete --force --cluster-id=ocid1.cluster.oc1.ap-tokyo-1.aaaaaaaau3ujoh2ipy52kn73ah4w6py4ofkrbd3yfxihcenhgcnbn6epa6zq
oci ce cluster delete --force --cluster-id=ocid1.cluster.oc1.ap-tokyo-1.aaaaaaaaj7nmey3meald2y6nczevcocjdt6pbcooj4nk52244ckktchkcjuq
oci ce cluster delete --force --cluster-id=ocid1.cluster.oc1.ap-tokyo-1.aaaaaaaaxx2d62qx5nl4zongga2ktdfbmsu4i2wnglurukfv3cbsj6dixhiq
oci ce cluster delete --force --cluster-id=ocid1.cluster.oc1.ap-tokyo-1.aaaaaaaabj6arlkwq2ytbgkbjiusvx5vdio7ghxq45oka4ukfcbir5vnpmpa
oci ce cluster delete --force --cluster-id=ocid1.cluster.oc1.ap-tokyo-1.aaaaaaaasqrr6ywgb2v7jruzvovlivdhh2emkr6kfwq73ayxucvjif7pjw2q
oci ce cluster delete --force --cluster-id=ocid1.cluster.oc1.ap-tokyo-1.aaaaaaaabj6arlkwq2ytbgkbjiusvx5vdio7ghxq45oka4ukfcbir5vnpmpa

oci ce cluster list -c $compartment_id
exit 0

cluster_id=$(oci ce cluster list -c $compartment_id --name $cluster_name --lifecycle-state ACTIVE | jq '.data[].id' | sed -e 's/^"//' -e 's/"$//')
oci ce cluster create-kubeconfig \
	--cluster-id $cluster_id \
	--file $KUBECONFIG \
	--region $region \
	--token-version 2.0.0 \
	--kube-endpoint PRIVATE_ENDPOINT

bastion_id=$(oci bastion bastion list -c $compartment_id --name $bastion_name --bastion-lifecycle-state ACTIVE --all | jq '.data[].id' | sed -e 's/^"//' -e 's/"$//')
api_private_endpoint=$(oci ce cluster get --cluster-id $cluster_id | jq '.data.endpoints["private-endpoint"]' | sed -e 's/^"//' -e 's/"$//')
private_ip=$(echo "$api_private_endpoint" | cut -d ':' -f1)
session_id=$(oci bastion session create-port-forwarding --bastion-id $bastion_id --target-private-ip $private_ip --session-ttl 10800 --target-port 6443 --ssh-public-key-file $public_key_file --wait-for-state SUCCEEDED | jq -r '.data.resources[].identifier')

echo "ACCESS KUBERNETES CLUSTER VIA PORT FORWARDING"
username=$(oci bastion session get --session-id $session_id | jq '.data["target-resource-details"]["target-resource-operating-system-user-name"]' | sed -e 's/^"//' -e 's/"$//')
bastion_ip=$(oci bastion session get --session-id $session_id | jq '.data["target-resource-details"]["target-resource-private-ip-address"]' | sed -e 's/^"//' -e 's/"$//')
tunnel_command=$(oci bastion session get --session-id $session_id | jq '.data["ssh-metadata"]["command"]' | sed -e 's/^"//' -e 's/"$//')

# Remove \ from command
tunnel_command=${tunnel_command//'\'/''}

# Substite the private key path for <privateKey> in the bastion SSH command
tunnel_command="${tunnel_command//<privateKey>/$private_key_file}"

# Add the k8s api forwarding port to the command, as well as necessary flags
tunnel_command="${tunnel_command/${username}@${bastion_ip}/-f ${username}@${bastion_ip} -L $port:${api_private_endpoint} -N}"

# Substite the localport in the bastion SSH command
tunnel_command="${tunnel_command//<localPort>/$port}"

# Disable host key verification
tunnel_command="${tunnel_command//ssh -i/ssh -4 -v -o StrictHostKeyChecking=no -o ServerAliveInterval=30 -o ServerAliveCountMax=5 -o ExitOnForwardFailure=yes -i}"

tunnel_command="while true; do { while true; do echo echo ping; sleep 10; done } | ${tunnel_command};sleep 10;done &"

cp $KUBECONFIG "${KUBECONFIG}_original"

# Substitute 127.0.0.1 into kubeconfig file
sed -i.bak "s/${api_private_endpoint}/127.0.0.1:$port/g" $KUBECONFIG

echo $tunnel_command

# Run SSH command
eval $tunnel_command

sleep 5

echo "KUBECTL READY TO USE"