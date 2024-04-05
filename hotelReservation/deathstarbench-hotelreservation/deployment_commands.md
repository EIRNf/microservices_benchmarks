Minikub
Minikub
minikube start

Service endpoints
minikube service --all -n hotel-res1

Build and push images
docker buildx build --push -t eirn/dsbpp_hotel_reserv:test .

HeLM Build hotel reservation
helm upgrade --install test . --create-namespace --wait --set image.repository=eirn/dsbpp_hotel_reserv --set image.tag=test --set features.gcPercent=1000 --set features.memcTimeout=10 -n hotel-res1 --debug

helm upgrade --install shm . --create-namespace --wait --set image.repository=eirn/dsbpp_hotel_reserv --set image.tag=shm --set features.gcPercent=1000 --set features.memcTimeout=10 -n hotel-res1 --debug

helm upgrade --install baseline . --create-namespace --wait --set image.repository=eirn/dsbpp_hotel_reserv --set image.tag=baseline --set features.gcPercent=1000 --set features.memcTimeout=10 -n hotel-res1 --debug

helm upgrade --install updatedbaseline . --create-namespace --wait --set image.repository=eirn/dsbpp_hotel_reserv --set image.tag=updated_baseline --set features.gcPercent=1000 --set features.memcTimeout=10 -n hotel-res1 --debug

Pyroscope
helm repo add pyroscope-io https://pyroscope-io.github.io/helm-chart 
helm install pyroscope pyroscope-io/pyroscope --set service.type=NodePort

kubectl port-forward (podname) -n hotel-res1 4040:4040 --address='0.0.0.0'

kubectl port-forward (podname) --address 0.0.0.0 4040:(podPort) 


## Full Flow for testing

 ### Defaults
- REMOTE_RUNNER
    - IP ADDRESS
- Hotelres 
    - deployment/frontend-shm-hotelres - port: 5000
    - appnames: geo, search, etc...
- Local IP Address

### Inputs
- RELEASE_NAME 
    Release name and image to deploy
- REMOTE
    Indicates remote or local tests
- TEST
    - functional
        - endpoints?
    - wrk2
        - 
    - bench
        - endpoints?
        - num_instances
        - num_requests (per instance)
- CAPTURE_LOGS
    - app: 
- GRAPHS:


### Build
- docker buildx build --push -t eirn/dsbpp_hotel_reserv:(REALEASE_NAME) .
### Deploy
- minikube delete && minikube start
- cd deathstarbench-hotelreservation/helm_hotelReservation
- helm upgrade --install (REALEASE_NAME) . --wait --set image.repository=eirn/dsbpp_hotel_reserv --set image.tag=(REALEASE_NAME) --set features.gcPercent=1000 --set features.memcTimeout=10 --debug
helm upgrade --install shm . --wait --set image.repository=eirn/dsbpp_hotel_reserv --set image.tag=shm --set features.gcPercent=1000 --set features.memcTimeout=10 --debug


// doesnt work 
- helm uninstall (REALEASE_NAME)

### Port Forward
*primarily for frontend or jeager endpoint*

- kubectl port-forward (podname) --address 0.0.0.0 4040:(podPort) 
- kubectl port-forward deployment/frontend-(REALEASE_NAME)-(hotelres) 4040:5000

- kubectl port-forward deployment/frontend-shm-hotelres 4040:5000

### Run Tests

#### Functional Tests
#### wrk2
#### Golang Benchmarks

### Capture Extra Metrics
kubectl logs -l app-name=geo -f
kubectl logs -l app-name=search -f


### Collect outputs
GRAPHS?
- will depend on test run

