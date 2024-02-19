#!/bin/bash -xe 

build=true

## PROJECT VARIABLES
DOCKER_PROJECT=eirn/dsbpp_hotel_reserv
RELEASE_NAME=shm

## LOCAL VARIABLES
cd "$(dirname "$0")"
SCRIPT_DIR="$(pwd)"
DEPLOY_HELM_DIR="$SCRIPT_DIR/deathstarbench-hotelreservation/helm_hotelReservation"

## CLEAN KUBERNETES 
minikube delete
# minikube start
 minikube start --extra-config=kubelet.feature-gates="CPUManager=true"
# minikube start --extra-config=kubelet.feature-gates="CPUManager=true" --extra-config=kubelet.config=$SCRIPT_DIR/deathstarbench-hotelreservation/kubelet.yaml


# minikube start --feature-gates=CPUManager=true --v=5 --force-systemd=true  --extra-config=kubeadm.ignore-preflight-errors=SystemVerification --extra-config=kubelet.config=kubelet.config=$SCRIPT_DIR/deathstarbench-hotelreservation/kubelet.yaml


# --extra-config=kubelet.config=$SCRIPT_DIR/deathstarbench-hotelreservation/kubelet.yaml


## BUILD
if [ "$build" = true ] ; then
    docker buildx build --push -t $DOCKER_PROJECT:$RELEASE_NAME .
fi

## DEPLOY
cd $DEPLOY_HELM_DIR
helm upgrade --install $RELEASE_NAME . --wait --set image.repository=$DOCKER_PROJECT --set image.tag=$RELEASE_NAME --set features.gcPercent=1000 --set features.memcTimeout=10 --debug

## PORT FORWARD
cd $SCRIPT_DIR
kubectl port-forward deployment/frontend-$RELEASE_NAME-hotelres --address 0.0.0.0 4040:5000

## TESTING
### Call Test script
# ./test.sh
