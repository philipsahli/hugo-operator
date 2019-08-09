# commands

    operator-sdk new --repo github.com/philipsahli/hugo-operator hugo-operator

    operator-sdk add api --api-version=hugo.sahli.net/v1alpha1 --kind=Site

    operator-sdk add controller --api-version=hugo.sahli.net/v1alpha1 --kind=Site