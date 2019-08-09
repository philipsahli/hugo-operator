OPERATOR_NAME := hugo-operator
WATCH_NAMESPACE := "hugo-site"

crd:
		kubectl apply -f deploy/crds/hugo_v1alpha1_site_crd.yaml
		kubectl get crd --all-namespaces

deploy: crd

resource: 
		kubectl apply -f deploy/crds/hugo_v1alpha1_site_cr.yaml --namespace hugo-site

run:
		WATCH_NAMESPACE=hugo-site operator-sdk up local --verbose 