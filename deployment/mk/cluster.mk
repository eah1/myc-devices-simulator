
KIND_CLUSTER := myc-devices-simulator-cluster

kind-up:
	kind create cluster \
		--image kindest/node:v1.21.1@sha256:69860bda5563ac81e3c0057d654b5253219618a22ec3a346306239bba8cfa1a6 \
		--name $(KIND_CLUSTER) \
		--config deployment/k8s/kind/kind-config.yaml
	kubectl create namespace myc-devices-simulator-cluster
	kubectl config set-context --current --namespace=myc-devices-simulator-cluster-system
	kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/kind/deploy.yaml

kind-down:
	kind delete cluster --name $(KIND_CLUSTER)

kind-status:
	kubectl get nodes -o wide
	kubectl get svc -o wide
	kubectl get pods -o wide --watch --all-namespaces

kind-apply-db:
	kustomize build deployment/k8s/kind/myc-devices-simulator-database-pod | kubectl apply -f -
	kubectl wait --namespace=database-system --timeout=120s --for=condition=Available deployment/database-pod
kind-status-db:
	kubectl get pods -o wide --watch --namespace=database-system
kind-logs-db:
	kubectl logs -l app=database --namespace=database-system --all-containers=true -f --tail=100
kind-db-migration:
	GOOSE_DRIVER=postgres GOOSE_DBSTRING="$(POSTGRES_KIND_URI)" goose -dir "$(SCHEMA_DIR)" status
	GOOSE_DRIVER=postgres GOOSE_DBSTRING="$(POSTGRES_KIND_URI)" goose -dir "$(SCHEMA_DIR)" up

kind-load-simulator-api:
	kind load docker-image $(API_IMAGE_NAME):$(VERSION) --name $(KIND_CLUSTER)
kind-apply-simulator-api:
	kustomize build deployment/k8s/kind/myc-devices-simulator-api-pod | kubectl apply -f -
kind-logs-simulator-api:
	kubectl logs -l app=simulator-api --all-containers=true -f --tail=100 --namespace=device-simulator-system