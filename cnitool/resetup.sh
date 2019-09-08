#!/usr/bin/env bash
# https://github.com/YaoZengzeng/KubernetesResearch/blob/master/%E6%B7%B1%E5%85%A5%E7%90%86%E8%A7%A3CNI.md

echo "remove interface: cni0"
ip link delete cni0

ROOT_VETH=$(ip link | grep -E "veth\\w+"  | awk '{print $2}' | cut -d '@' -f1)
echo "remove interface: $ROOT_VETH"
ip link delete $ROOT_VETH

sudo rm -rf /var/lib/cni/*

echo "remove netns ns"
ip netns delete ns

echo "add netns ns"
ip netns add ns