# infra

## cluster scalability

1. start from 1 node
2. setup control-plane with complete default config from kubeKey
3. add 2 more nodes(control-plane,worker,etcd) in yaml
4. `kk add nodes -f ./ai-server/infra/dev.yaml`

## image registry

1. use tencent cloud personal plan for early development
2. later: [self-hosted harbor ha](https://github.com/kubesphere/kubekey/blob/master/docs/harbor-ha.md)

## user traffic scalability

1. reduce cloud module using like CLB/CVM(dedicated LB server)
2. expose 443 ports on each node ready for ingress
3. access any one of servers via nodePubIP:443
4. use dns loadbalancer to route traffic to any one of servers
