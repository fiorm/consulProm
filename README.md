# consulProm
e2e test architecture for consul service discovery mechanism 

## steps to follow : 

1. Make sure you have minikube  and helm is installed. 

2. Run the below command in your terminal 

helm install -f consul-helm/helm-consul-values.yml hashicorp ./consul-helm

3. To check the consul ui in browser 

minikube service hashicorp-consul-ui



