#!/bin/bash

compartment_id=$1
region=$2
cluster_name=$3
bastion_name=$4
public_key_file=$5
private_key_file=$6
KUBECONFIG=$7
vcn_id=$8
sec_list_id=$9

oci bastion bastion list --compartment-id=$compartment_id
oci network vcn get --vcn-id=$vcn_id
oci network security-list get --security-list-id=$sec_list_id

echo "CREATE KUBECONFIG at $KUBECONFIG"

rm $KUBECONFIG

cluster_id=$(oci ce cluster list -c $compartment_id --name $cluster_name --lifecycle-state ACTIVE | jq '.data[].id' | sed -e 's/^"//' -e 's/"$//')
echo "cluster_id is $cluster_id"
oci ce cluster get --cluster-id $cluster_id
oci ce cluster create-kubeconfig \
	--cluster-id $cluster_id \
	--file $KUBECONFIG \
	--region $region \
	--token-version 2.0.0 \
	--kube-endpoint PRIVATE_ENDPOINT

bastion_id=$(oci bastion bastion list -c $compartment_id --name $bastion_name --bastion-lifecycle-state ACTIVE --all | jq '.data[].id' | sed -e 's/^"//' -e 's/"$//')
echo "bastion_id is $bastion_id"
api_private_endpoint=$(oci ce cluster get --cluster-id $cluster_id | jq '.data.endpoints["private-endpoint"]' | sed -e 's/^"//' -e 's/"$//')
echo "api_private_endpoint is $api_private_endpoint"
private_ip=$(echo "$api_private_endpoint" | cut -d ':' -f1)
session_id=$(oci bastion session create-port-forwarding --bastion-id $bastion_id --target-private-ip $private_ip --target-port 6443 --ssh-public-key-file $public_key_file --wait-for-state SUCCEEDED | jq -r '.data.resources[].identifier')

echo "ACCESS KUBERNETES CLUSTER VIA PORT FORWARDING"
username=$(oci bastion session get --session-id $session_id | jq '.data["target-resource-details"]["target-resource-operating-system-user-name"]' | sed -e 's/^"//' -e 's/"$//')
bastion_ip=$(oci bastion session get --session-id $session_id | jq '.data["target-resource-details"]["target-resource-private-ip-address"]' | sed -e 's/^"//' -e 's/"$//')
tunnel_command=$(oci bastion session get --session-id $session_id | jq '.data["ssh-metadata"]["command"]' | sed -e 's/^"//' -e 's/"$//')

# Remove \ from command
tunnel_command=${tunnel_command//'\'/''}

# Substite the private key path for <privateKey> in the bastion SSH command
tunnel_command="${tunnel_command//<privateKey>/$private_key_file}"

# Add the k8s api forwarding port to the command, as well as necessary flags
tunnel_command="${tunnel_command/${username}@${bastion_ip}/-f ${username}@${bastion_ip} -L 6443:${api_private_endpoint} -N}"

# Substite the localport in the bastion SSH command
tunnel_command="${tunnel_command//<localPort>/6443}"

# Disable host key verification
tunnel_command="${tunnel_command//ssh -i/ssh -4 -v -o StrictHostKeyChecking=no -i}"

#tunnel_command="${tunnel_command} &"

# Substitute 127.0.0.1 into kubeconfig file
sed -i.bak "s/${api_private_endpoint}/127.0.0.1:6443/g" $KUBECONFIG

echo $tunnel_command

# Run SSH command
eval $tunnel_command

sleep 5

echo "KUBECTL READY TO USE"