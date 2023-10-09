### Minikub 
minikube start

#### Service endpoints
minikube service --all -n hotel-res1

### Build and push images
docker buildx build --push -t eirn/dsbpp_hotel_reserv:test .

### HeLM Build hotel reservation
helm upgrade --install test . --create-namespace --wait --set image.repository=eirn/dsbpp_hotel_reserv --set image.tag=test --set features.gcPercent=1000 --set features.memcTimeout=10 -n hotel-res1 --debug

helm upgrade --install shm . --create-namespace --wait --set image.repository=eirn/dsbpp_hotel_reserv --set image.tag=shm --set features.gcPercent=1000 --set features.memcTimeout=10 -n hotel-res1 --debug

helm upgrade --install baseline . --create-namespace --wait --set image.repository=eirn/dsbpp_hotel_reserv --set image.tag=baseline --set features.gcPercent=1000 --set features.memcTimeout=10 -n hotel-res1 --debug

helm upgrade --install updatedbaseline . --create-namespace --wait --set image.repository=eirn/dsbpp_hotel_reserv --set image.tag=updated_baseline --set features.gcPercent=1000 --set features.memcTimeout=10 -n hotel-res1 --debug

### Pyroscope
helm repo add pyroscope-io https://pyroscope-io.github.io/helm-chart
helm install pyroscope pyroscope-io/pyroscope --set service.type=NodePort

 kubectl port-forward (podname)  -n hotel-res1 4040:4040 --address='0.0.0.0'
