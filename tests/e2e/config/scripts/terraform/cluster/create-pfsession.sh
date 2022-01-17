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
oci ce cluster delete --force --cluster-id=ocid1.cluster.oc1.ap-tokyo-1.aaaaaaaa2f6jwfar5ypknb2mtp2ctpseh3dz5w3lyalacptlacyrzxfzyjwq
oci ce cluster delete --force --cluster-id=ocid1.cluster.oc1.ap-tokyo-1.aaaaaaaapxlgdrstk4dxhbqcda25mqgyawcv7exs3jdtq2aj3chhnrwobvhq
oci ce cluster delete --force --cluster-id=ocid1.cluster.oc1.ap-tokyo-1.aaaaaaaaal77ysxulxce5jrysy5esqaougrig7qgjaa7gwqeecqsu7sv7uqq
oci ce cluster delete --force --cluster-id=ocid1.cluster.oc1.ap-tokyo-1.aaaaaaaam566vmtipgcwufq4oqwnav6he4cafkus7hfcr6t4uco6dlqqb6uq
oci ce cluster delete --force --cluster-id=ocid1.cluster.oc1.ap-tokyo-1.aaaaaaaagng7fk22lwrcxuntno7zrugeopfo2dmx7sknfnygwcg53rkqdtpq
oci ce cluster delete --force --cluster-id=ocid1.cluster.oc1.ap-tokyo-1.aaaaaaaa76cxgu5xztvl3ai2cuywqui5qldur5g3sz7kknpahcsafgqntj2a
oci ce cluster delete --force --cluster-id=ocid1.cluster.oc1.ap-tokyo-1.aaaaaaaatbvtkqjmaaewdsw7w74tz6m2zgogvjcbpcflimn75crqjsxnxtcq
oci ce cluster delete --force --cluster-id=ocid1.cluster.oc1.ap-tokyo-1.aaaaaaaab3cwwjtyzmgpfgocvfm5tybjx7442salc4ui3prvncis2jd2iz7q
oci ce cluster delete --force --cluster-id=ocid1.cluster.oc1.ap-tokyo-1.aaaaaaaal2ghmfmtdoatyoclmgq2knthcysqhfrqjbzrlm635c74l4subkaq
oci ce cluster delete --force --cluster-id=ocid1.cluster.oc1.ap-tokyo-1.aaaaaaaa34dujoybznggbn53cksbvduv22psv3g6izybavv6tcc7fn3re4ya
oci ce cluster delete --force --cluster-id=ocid1.cluster.oc1.ap-tokyo-1.aaaaaaaarlweqkh7qopm5xon4beycqg7frscwo4ydtvfowr34ctmxv234gra
oci ce cluster delete --force --cluster-id=ocid1.cluster.oc1.ap-tokyo-1.aaaaaaaakntkeplcddtvr4dbtw57xapzou2ja6nflyr3ic2ohcaqvsltxeda
oci ce cluster delete --force --cluster-id=ocid1.cluster.oc1.ap-tokyo-1.aaaaaaaarlerbw2oqu25gf6nx6xkn7amu6cxd4je2hkvzivrbcrnleggev7q
oci ce cluster delete --force --cluster-id=ocid1.cluster.oc1.ap-tokyo-1.aaaaaaaasudt5kgnypsc2o4habgd3ldzypixehidddt4ffdxdc4kecgeooea
oci ce cluster delete --force --cluster-id=ocid1.cluster.oc1.ap-tokyo-1.aaaaaaaas6whpry6jqwo3kj5oah5j6yre7r6oygi25u7elukfcxck3lojrjq
oci ce cluster delete --force --cluster-id=ocid1.cluster.oc1.ap-tokyo-1.aaaaaaaav44hmczh2oglsoua6bvjx2qyzahtp6lwxoshomeamcecotasu2qq
oci ce cluster delete --force --cluster-id=ocid1.cluster.oc1.ap-tokyo-1.aaaaaaaamjg27cvjjnrv6peeydidikr7lin65lhppqyetxovtckkmykujfga
oci ce cluster list -c $compartment_id
exit 0

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