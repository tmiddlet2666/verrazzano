#!/bin/bash

compartment_id=$1
region=$2
cluster_name=$3
bastion_name=$4
public_key_file=$5
private_key_file=$6
KUBECONFIG=$7
port=$8
#vcn_id=$9
#sec_list_id="${10}"

#oci network vcn get --vcn-id=$vcn_id
#oci network security-list update --ingress-security-rules='[{"description": "east west","icmpOptions": null,"isStateless": false,"protocol": "all","source": "10.196.0.0/16","sourceType": "CIDR_BLOCK","tcpOptions": null,"udpOptions": null},{"description": null,"icmpOptions": {"code": null,"type": 3},"isStateless": false,"protocol": "1","source": "10.196.0.0/16","sourceType": "CIDR_BLOCK","tcpOptions": null,"udpOptions": null},{"description": null,"icmpOptions": {"code": 4,"type": 3},"isStateless": false,"protocol": "1","source": "0.0.0.0/0","sourceType": "CIDR_BLOCK","tcpOptions": null,"udpOptions": null},{"description": null,"icmpOptions": null,"isStateless": false,"protocol": "6","source": "0.0.0.0/0","sourceType": "CIDR_BLOCK","tcpOptions": {"destinationPortRange": {"max": 22,"min": 22},"sourcePortRange": null},"udpOptions": null},{"description": null,"icmpOptions": null,"isStateless": false,"protocol": "6","source": "0.0.0.0/0","sourceType": "CIDR_BLOCK","tcpOptions": {"destinationPortRange": {"max": 6443,"min": 6443},"sourcePortRange": null},"udpOptions": null},{"description": null,"icmpOptions": null,"isStateless": false,"protocol": "6","source": "0.0.0.0/0","sourceType": "CIDR_BLOCK","tcpOptions": {"destinationPortRange": {"max": 443,"min": 443},"sourcePortRange": null},"udpOptions": null}]' --force --security-list-id=$sec_list_id 


#oci bastion bastion create --bastion-type STANDARD --compartment-id $compartment_id --target-subnet-id $target_subnet_id --client-cidr-list '["0.0.0.0/0"]' --max-session-ttl 10800 --name $bastion_name
#exit 0

rm $KUBECONFIG
#oci ce cluster delete --force --cluster-id=ocid1.cluster.oc1.ap-tokyo-1.aaaaaaaaq2qg35j4h7tvzronojvfqrbg524o5be67uv4ncc64c4j3use3j7q

oci lb load-balancer delete --force --load-balancer-id=ocid1.loadbalancer.oc1.ap-tokyo-1.aaaaaaaa2b7djjw5w3s62l4f2ofzdzujywblt7l2tachhpmkta6njhnvma3a
oci lb load-balancer delete --force --load-balancer-id=ocid1.loadbalancer.oc1.ap-tokyo-1.aaaaaaaa2cmdt6iaumv4oeilwz6p2botkrehdc4phlyqqzvbxquubb26hwkq
oci lb load-balancer delete --force --load-balancer-id=ocid1.loadbalancer.oc1.ap-tokyo-1.aaaaaaaa32b7j7qqqwxxhsdbndtrt5p2ammrijim756n62ni4iqto4ynbtfq
oci lb load-balancer delete --force --load-balancer-id=ocid1.loadbalancer.oc1.ap-tokyo-1.aaaaaaaa3gszpoby6rzf5bzjggrljq6ubzool5jpli4ieavebml35clncvua
oci lb load-balancer delete --force --load-balancer-id=ocid1.loadbalancer.oc1.ap-tokyo-1.aaaaaaaa3iuq4epbxzbkjhbgo65ni7sr2yad4y6b2aewdjn5rzugqdlntkra
oci lb load-balancer delete --force --load-balancer-id=ocid1.loadbalancer.oc1.ap-tokyo-1.aaaaaaaa52ci6o7a546d3xbrn7olxe3s2wnkmzu3ajkzjcralbr7oqvvinrq
oci lb load-balancer delete --force --load-balancer-id=ocid1.loadbalancer.oc1.ap-tokyo-1.aaaaaaaa7zrpjma4bifyl4yw67xpbyklicsndrfcm262gzruwzvrkkbps6ha
oci lb load-balancer delete --force --load-balancer-id=ocid1.loadbalancer.oc1.ap-tokyo-1.aaaaaaaaa4cwcnhyuowo2mlan5sdbxodm7ewrzlp6icynb3wk2zhaufxbkxa
oci lb load-balancer delete --force --load-balancer-id=ocid1.loadbalancer.oc1.ap-tokyo-1.aaaaaaaaa5zveevat2ip3awnfhu7haghujeo7wdoxlldfc6r7kikxouojtda
oci lb load-balancer delete --force --load-balancer-id=ocid1.loadbalancer.oc1.ap-tokyo-1.aaaaaaaaahfkwk7rciq5cdl455rcnixeyhdjs2slzvzvoetxcevcw7f5d6bq
oci ce cluster list -c $compartment_id
#exit 0
oci lb load-balancer list --compartment-id=$compartment_id

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

tunnel_command="while :; do { while :; do echo echo ping; sleep 10; done } | ${tunnel_command};sleep 10;done > $port.out 2>&1 &"

cp $KUBECONFIG "${KUBECONFIG}_original"

# Substitute 127.0.0.1 into kubeconfig file
sed -i.bak "s/${api_private_endpoint}/127.0.0.1:$port/g" $KUBECONFIG

echo $tunnel_command

# Run SSH command
eval $tunnel_command

sleep 5

echo "KUBECTL READY TO USE"

while :; do kubectl --kubeconfig=$KUBECONFIG get nodes | echo "failed ping";sleep 10;done > $port_ping.out 2>&1 &