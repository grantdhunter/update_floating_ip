# UFI Update Floating IP

A small utility that automatically assigns a Kubernetes DigitalOcean droplet to a floating IP.

# Background
I am runnign a single node kubernetes cluster on DO and do not want to pay for a load balancer. This works find and dandy until DO runs an upgrade which recycles my kube node that is assigned to my floating IP. I build this utility to automatically reassing the droplet to the floating IP any time the node starts up.

# Design
UFI runs as a daemonset so it's garaunteed to have a single pod per node. This works very nicely for a single node cluster and I am guessing that it will work ok for a multi node cluster, the last node to start will get the assignment to the floating IP.

The actual application runs does it's business then sleeps until the end of the universe (or life time of the node). This
 is kinda a hack to get around the fact that kubernetes doens't have the ability to run a job once during node start.

# Running
``` update_floating_ip [--wait] your.domain
```
