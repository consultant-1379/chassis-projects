module nef-nrfintegration

go 1.12

replace k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20190404173353-6a84e37a896d

replace cloud.google.com/go => github.com/GoogleCloudPlatform/google-cloud-go v0.39.0

require (
	gerrit.ericsson.se/nef/nef-golangcommon latest
	github.com/evanphx/json-patch v4.2.0+incompatible // indirect
	github.com/google/gofuzz v1.0.0 // indirect
	github.com/googleapis/gnostic v0.2.0 // indirect
	github.com/json-iterator/go v1.1.6 // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	golang.org/x/net v0.0.0-20190514140710-3ec191127204
	golang.org/x/oauth2 v0.0.0-20190517181255-950ef44c6e07 // indirect
	golang.org/x/time v0.0.0-20190308202827-9d24e82272b4 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	k8s.io/api v0.0.0-20190515023547-db5a9d1c40eb
	k8s.io/apimachinery v0.0.0-20190515023456-b74e4c97951f
	k8s.io/client-go v11.0.0+incompatible
	k8s.io/kube-openapi v0.0.0-20190510232812-a01b7d5d6c22 // indirect
	k8s.io/utils v0.0.0-20190506122338-8fab8cb257d5 // indirect
	sigs.k8s.io/yaml v1.1.0 // indirect
)
