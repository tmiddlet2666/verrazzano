#!/bin/bash

compartment_id=$1
bastion_name=$2
KUBECONFIG=$3
public_key_file=$4
private_key_file=$5
port=$6
target_subnet_id=$7
echo "CREATE KUBECONFIG at $KUBECONFIG"

kubectl --kubeconfig $KUBECONFIG get vz -o yaml
bastion_id=$(oci bastion bastion list -c $compartment_id --name $bastion_name --bastion-lifecycle-state ACTIVE --all | jq '.data[].id' | sed -e 's/^"//' -e 's/"$//')
private_ip=$(kubectl --kubeconfig $KUBECONFIG get vz -o yaml | yq e ".items[0].status.instance.prometheusUrl" - | cut -d "." -f5- | cut -d "n" -f1 | rev | cut -c2- | rev)
session_id=$(oci bastion session create-port-forwarding --bastion-id $bastion_id --target-private-ip $private_ip --session-ttl 10800 --target-port 443 --ssh-public-key-file $public_key_file --wait-for-state SUCCEEDED | jq -r '.data.resources[].identifier')
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

echo $tunnel_command

# Run SSH command
eval $tunnel_command

touch /tmp/hosts
list=$(kubectl --kubeconfig $KUBECONFIG get vz -o yaml | yq e ".items[0].status.instance" -)
while IFS= read -r line; do
    echo "$private_ip $(echo $line | cut -d "/" -f3-)" >> /tmp/hosts
done <<< "$list"
cat /tmp/hosts
cat /tmp/hosts | sudo tee -a /etc/hosts >/dev/null
cat /etc/hosts
rm -rf /tmp/hosts
sudo iptables -t nat -A OUTPUT -p tcp --dport 443 -d $private_ip -j DNAT --to-destination 127.0.0.1:$port

sleep 5

while :; do curl -k -v https://$private_ip | echo "failed ping";sleep 30;done > "${port}_ping.out" 2>&1 &