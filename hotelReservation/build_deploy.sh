#!/bin/bash -x

## PROJECT VARIABLES
DOCKER_PROJECT=eirn/dsbpp_hotel_reserv
RELEASE_NAME=shm

## LOCAL VARIABLES
cd "$(dirname "$0")"
SCRIPT_DIR="$(pwd)"
DEPLOY_HELM_DIR="$SCRIPT_DIR/deathstarbench-hotelreservation/helm_hotelReservation"

## CLEAN KUBERNETES 
minikube delete && minikube start

## BUILD
docker buildx build --push -t $DOCKER_PROJECT:$RELEASE_NAME .

## DEPLOY
cd $DEPLOY_HELM_DIR
helm upgrade --install $RELEASE_NAME . --wait --set image.repository=$DOCKER_PROJECT --set image.tag=$RELEASE_NAME --set features.gcPercent=1000 --set features.memcTimeout=10 --debug

## PORT FORWARD
cd $SCRIPT_DIR
kubectl port-forward deployment/frontend-$RELEASE_NAME-hotelres --address 0.0.0.0 4040:5000

## TESTING
### Call Test script
# ./test.sh