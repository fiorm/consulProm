.PHONY: update
update:  
    @eval $$(minikube docker-env) ;\
    docker image build -t message-api:v1 -f Dockerfile .
    kubectl set image deployment/message-api *=message-api:v1

##  work is pending here 