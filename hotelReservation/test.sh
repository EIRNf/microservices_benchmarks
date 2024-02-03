#!/bin/bash -x

## PROJECT VARIABLES
DOCKER_PROJECT=eirn/dsbpp_hotel_reserv
RELEASE_NAME=shm

## REMOTE VARIABLES
REMOTE=0
DEPLOY_EXTERNAL_IP="$(gcloud compute instances describe instance-1 \
  --format='get(networkInterfaces[0].accessConfigs[0].natIP)')"
TEST_EXTERNAL_IP="$(gcloud compute instances describe instance-3 \
  --format='get(networkInterfaces[0].accessConfigs[0].natIP)')"
ENDPOINT_PORT=4040
DEPLOY_ADDRESS=$DEPLOY_EXTERNAL_IP:$ENDPOINT_PORT

## LOCAL VARIABLES
cd "$(dirname "$0")"
SCRIPT_DIR="$(pwd)"
DEPLOY_HELM_DIR="$("$SCRIPT_DIR"/deathstarbench-hotelreservation/helm_hotelReservation)"

NODE_PORT=$(kubectl get --namespace default -o jsonpath="{.spec.ports[0].nodePort}" services frontend-shm-hotelres)
NODE_IP=$(kubectl get nodes --namespace default -o jsonpath="{.items[0].status.addresses[0].address}")
echo http://$NODE_IP:$NODE_PORT


ENDPOINT_TESTING=$SCRIPT_DIR/endpoint_testing

RUN_DIR=$SCRIPT_DIR/endpoint_testing/runs

## TESTING

### VERIFY FUNCTIONALITY
cd ENDPOINT_TESTING

### REMOTE
if $REMOTE
then

else
### Local 
go test --url=http://127.0.0.1:4040
fi

### RUN BENCHMARKS
#### ENDPOINT TESTING
cd ENDPOINT_TESTING
./endpoint_testing
#### WRK2
cd $SCRIPT_DIR
../wrk2/wrk -D exp -t 8 -c 8 -d 60 -L -s ./wrk2/scripts/hotels/hotels.lua $DEPLOY_ADDRESS -r  -R 1085

