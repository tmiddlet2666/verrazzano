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

rm $KUBECONFIG
oci ce cluster delete --force --cluster-id=ocid1.cluster.oc1.ap-tokyo-1.aaaaaaaarkzyy7cfaxopanhis2lwacgxbiu3x3cctxzna3fbgcteytlth3gq
oci ce cluster delete --force --cluster-id=ocid1.cluster.oc1.ap-tokyo-1.aaaaaaaaxm6pbizy6ttbgr6vjmonwrmrqd7i4rog3wdzxfmp4cx2576qhn4q
oci ce cluster delete --force --cluster-id=ocid1.cluster.oc1.ap-tokyo-1.aaaaaaaafndedtyhqi62qqjmxnlmvhzxcwjhhjkd2qespli44ca6bw3izh3a
oci ce cluster delete --force --cluster-id=ocid1.cluster.oc1.ap-tokyo-1.aaaaaaaaeh36dn5457jjawztx6cncwnsepy22gkmwvpgk5coccxoitbuhmga
oci ce cluster delete --force --cluster-id=ocid1.cluster.oc1.ap-tokyo-1.aaaaaaaajr7n5jxfch2e4ycjk6ieuyxtyk5xsbrccix5r36nrcqgfoco4aza
oci ce cluster delete --force --cluster-id=ocid1.cluster.oc1.ap-tokyo-1.aaaaaaaapecrg7gjaxkoamfv3i57pskgrckc4xnuxzvztxzqfcdkpfkpec3a
oci ce cluster delete --force --cluster-id=ocid1.cluster.oc1.ap-tokyo-1.aaaaaaaaghwezenuhp2avduw3b2f4dssezcc7wy5ockwih7joc5azo7f5eoa
oci ce cluster delete --force --cluster-id=ocid1.cluster.oc1.ap-tokyo-1.aaaaaaaavyv63ynbuo475oyfxxkbbv62ftqhxdwlerdvlncp5cshl5ndvska
oci ce cluster delete --force --cluster-id=ocid1.cluster.oc1.ap-tokyo-1.aaaaaaaax5zuf3g754o4lbucrh37nfqvxjn4e5epsfqg6w2v5cr7agpfr67a
oci ce cluster delete --force --cluster-id=ocid1.cluster.oc1.ap-tokyo-1.aaaaaaaavyv63ynbuo475oyfxxkbbv62ftqhxdwlerdvlncp5cshl5ndvska
oci ce cluster delete --force --cluster-id=ocid1.cluster.oc1.ap-tokyo-1.aaaaaaaaeh36dn5457jjawztx6cncwnsepy22gkmwvpgk5coccxoitbuhmga
oci ce cluster list -c $compartment_id
#exit 0

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
tunnel_command="${tunnel_command//ssh -i/ssh -4 -v -o StrictHostKeyChecking=no -o ServerAliveInterval=5 -o ServerAliveCountMax=10000 -i}"

tunnel_command="${tunnel_command} &"

cp $KUBECONFIG "${KUBECONFIG}_original"

# Substitute 127.0.0.1 into kubeconfig file
sed -i.bak "s/${api_private_endpoint}/127.0.0.1:$port/g" $KUBECONFIG

echo $tunnel_command

# Run SSH command
eval $tunnel_command

sleep 5

echo "KUBECTL READY TO USE"