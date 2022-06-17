#!/bin/bash

function delete_k8s_resources() {
  ( [ -z "$5" ] && kubectl get $1 --no-headers -o custom-columns="$2" || kubectl get $1 --no-headers -o custom-columns="$2" -n "$5" ) \
    | ( [ -z "$4" ] && cat || awk "$4" ) \
    | ( [ -z "$5" ] && xargsr kubectl delete $1 || xargsr kubectl delete $1 -n "$5" ) \
    || err_return $? "$3" || return $? # return on pipefail
}

k delete ns cattle-system
delete_k8s_resources mutatingwebhookconfigurations ":metadata.name,:metadata.labels" "Could not delete MutatingWebhookConfigurations from Rancher" '/cattle.io|app:rancher/ {print $1}' \
  || return $?
delete_k8s_resources validatingwebhookconfigurations ":metadata.name,:metadata.labels" "Could not delete ValidatingWebhookConfigurations from Rancher" '/cattle.io|app:rancher/ {print $1}' \
  || return $?
kubectl delete mutatingwebhookconfigurations.admissionregistration.k8s.io rancher.cattle.io --ignore-not-found \
  || err_return $?

helm ls -n fleet-system | awk '/fleet/ {print $1}' | xargsr helm uninstall -n fleet-system \
  || err_return $? "Could not delete fleet-system charts from helm" || return $? # return on pipefail
helm ls -n cattle-fleet-system | awk '/fleet/ {print $1}' | xargsr helm uninstall -n cattle-fleet-system \
  || err_return $? "Could not delete cattle-fleet-system charts from helm" || return $? # return on pipefail
helm ls -n cattle-fleet-local-system | awk '/fleet/ {print $1}' | xargsr helm uninstall -n cattle-fleet-local-system \
  || err_return $? "Could not delete cattle-fleet-local-system charts from helm" || return $? # return on pipefail
helm ls -n fleet-system | awk '/fleet/ {print $1}' | xargsr helm -n fleet-system uninstall \
  || err_return $? "Could not delete fleet-system charts from helm" || return $? # return on pipefail
helm ls -n cattle-system | awk '/rancher/ {print $1}' | xargsr helm uninstall -n cattle-system \
  || err_return $? "Could not delete cattle-system from helm" || return $? # return on pipefail

kubectl api-resources --api-group=management.cattle.io --verbs=delete -o name \
  | xargsr -n 1 kubectl get --all-namespaces --ignore-not-found -o custom-columns=":kind,:metadata.name,:metadata.namespace" \
  | awk '{res="";if ($1 != "") res=tolower($1)".management.cattle.io "tolower($2); if ($3 != "<none>" && res != "") res=res" -n "$3; if (res != "") cmd="kubectl patch "res" -p \x027{\"metadata\":{\"finalizers\":null}}\x027 --type=merge;kubectl delete --ignore-not-found "res; print cmd}' \
  | sh \
  || err_return $? "There were errors deleting rancher CRs"  # Continue if fai

crd_content=$(kubectl get crds --no-headers -o custom-columns=":metadata.name,:spec.group" | awk '/coreos.com|cattle.io/')
while [ "$crd_content" ]
do
  # remove finalizers from crds
  # Ignore patch failures and attempt to delete the resources anyway.
  patch_k8s_resources crds ":metadata.name,:spec.group" "Could not remove finalizers from CustomResourceDefinitions in Rancher" '/coreos.com|cattle.io/ {print $1}' '{"metadata":{"finalizers":null}}' \
    || true

  # delete crds
  # This process is backgrounded in order to timeout due to finalizers hanging
  delete_k8s_resources crds ":metadata.name,:spec.group" "Could not delete CustomResourceDefinitions from Rancher" '/coreos.com|management.cattle.io|cattle.io|fleet/ {print $1}' \
    || return $? &# return on pipefail
  sleep 30
  kill $! || true
  crd_content=$(kubectl get crds --no-headers -o custom-columns=":metadata.name,:spec.group" | awk '/coreos.com|cattle.io/')
done

delete_k8s_resources clusterrolebinding ":metadata.name,:metadata.labels" "Could not delete ClusterRoleBindings from Rancher" '/cattle.io|app:rancher|rancher-webhook|fleetworkspace-|fleet-|gitjob/ {print $1}' \
  || return $? # return on pipefail

delete_k8s_resources clusterrole ":metadata.name,:metadata.labels" "Could not delete ClusterRoles from Rancher" '/cattle.io|app:rancher|fleetworkspace-|fleet-|gitjob/ {print $1}' \
  || return $? # return on pipefail

default_names=("default" "kube-node-lease" "kube-public" "kube-system")
for namespace in "${default_names[@]}"
do
  delete_k8s_resources rolebinding ":metadata.name" "Could not delete RoleBindings from Rancher in namespace ${namespace}" '/clusterrolebinding-/' "${namespace}" \
    || return $? # return on pipefail
  delete_k8s_resources rolebinding ":metadata.name" "Could not delete RoleBindings from Rancher in namespace ${namespace}" '/^rb-/' "${namespace}" \
    || return $? # return on pipefail
done

kubectl delete configmap cattle-controllers -n kube-system  --ignore-not-found=true || err_return $? "Could not delete ConfigMap from Rancher in namespace kube-system" || return $?
kubectl delete configmap rancher-controller-lock -n kube-system --ignore-not-found=true || err_return $? "Could not delete ConfigMap rancher-controller-lock in namespace kube-system" || return $?

patch_k8s_resources namespaces ":metadata.name" "Could not remove finalizers from namespaces in Rancher" '/^cattle-|^local|^p-|^user-|^fleet|^rancher/ {print $1}' '{"metadata":{"finalizers":null}}' \
  || return $? # return on pipefail

if kubectl get serviceaccount -n cattle-system rancher > /dev/null 2>&1 ; then
  if ! kubectl delete serviceaccount -n cattle-system rancher ; then
    error "Failed to delete the service account rancher in namespace cattle-system."
  fi
fi

if kubectl get serviceaccount -n cattle-system default > /dev/null 2>&1 ; then
  if ! kubectl delete serviceaccount -n cattle-system default ; then
    error "Failed to delete the service account default in namespace cattle-system."
  fi
fi

rancher_namespaces=("cattle-fleet-clusters-system" "cattle-fleet-local-system" "cattle-fleet-system" "cattle-global-data" "cattle-global-nt" "cattle-impersonation-system" "fleet-default" "fleet-local")
for namespace in "${rancher_namespaces[@]}"
do
  if ! kubectl delete namespace ${namespace} --ignore-not-found=true ; then
    error "Failed to delete the namespace ${namespace}"
  fi
done

for namespace in "${default_names[@]}"
do
  kubectl get secret -n "${namespace}" --no-headers -o custom-columns=":metadata.name,:metadata.annotations" \
    | awk '/field.cattle.io\/projectId:/ {print $1}' \
    | xargsr -I resource kubectl annotate secret resource -n "${namespace}" field.cattle.io/projectId- \
    || err_return $? "Could not delete Annotations from Rancher" || return $? # return on pipefail
done

kubectl get namespaces --no-headers -o custom-columns=":metadata.name,:metadata.finalizers" \
  | awk '/controller.cattle.io/ {print $1}' \
  | xargsr kubectl patch namespace -p '{"metadata":{"finalizers":null}}' --type=merge \
  || err_return $? "Could not remove Rancher finalizers from all namespaces" || return $? # return on pipefail
