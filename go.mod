module db-operator

require (
	cloud.google.com/go v0.43.0 // indirect
	contrib.go.opencensus.io/exporter/ocagent v0.5.1 // indirect
	github.com/Azure/go-autorest/autorest v0.5.0 // indirect
	github.com/Azure/go-autorest/autorest/mocks v0.2.0 // indirect
	github.com/NYTimes/gziphandler v1.0.1 // indirect
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751 // indirect
	github.com/alecthomas/units v0.0.0-20190717042225-c3de453c63f4 // indirect
	github.com/appscode/jsonpatch v2.0.0+incompatible // indirect
	github.com/aws/aws-sdk-go v1.21.8
	github.com/census-instrumentation/opencensus-proto v0.2.1 // indirect
	github.com/coreos/prometheus-operator v0.31.1 // indirect
	github.com/emicklei/go-restful v2.9.6+incompatible // indirect
	github.com/go-logr/logr v0.1.0
	github.com/go-openapi/spec v0.19.2 // indirect
	github.com/go-openapi/swag v0.19.4 // indirect
	github.com/gobuffalo/envy v1.7.0 // indirect
	github.com/golang/groupcache v0.0.0-20190702054246-869f871628b6 // indirect
	github.com/googleapis/gnostic v0.3.0 // indirect
	github.com/gophercloud/gophercloud v0.2.0 // indirect
	github.com/gregjones/httpcache v0.0.0-20190611155906-901d90724c79 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.9.5 // indirect
	github.com/hashicorp/golang-lru v0.5.2 // indirect
	github.com/json-iterator/go v1.1.7 // indirect
	github.com/lib/pq v1.2.0
	github.com/magiconair/properties v1.8.1 // indirect
	github.com/mailru/easyjson v0.0.0-20190626092158-b2ccc519800e // indirect
	github.com/onsi/ginkgo v1.8.0 // indirect
	github.com/onsi/gomega v1.5.0 // indirect
	github.com/operator-framework/operator-sdk v0.9.0
	github.com/pelletier/go-toml v1.4.0 // indirect
	github.com/prometheus/common v0.6.0 // indirect
	github.com/prometheus/procfs v0.0.3 // indirect
	github.com/spf13/pflag v1.0.3
	github.com/spf13/viper v1.4.0
	golang.org/x/crypto v0.0.0-20190701094942-4def268fd1a4 // indirect
	golang.org/x/net v0.0.0-20190724013045-ca1201d0de80 // indirect
	golang.org/x/sys v0.0.0-20190726091711-fc99dfbffb4e // indirect
	golang.org/x/tools v0.0.0-20190725161231-2e34cfcb95cb // indirect
	gomodules.xyz/jsonpatch/v2 v2.0.0 // indirect
	google.golang.org/grpc v1.22.1 // indirect
	k8s.io/api v0.0.0-20190726022912-69e1bce1dad5
	k8s.io/apiextensions-apiserver v0.0.0-20190726024412-102230e288fd // indirect
	k8s.io/apimachinery v0.0.0-20190726022757-641a75999153
	k8s.io/client-go v11.0.0+incompatible
	k8s.io/klog v0.3.3 // indirect
	k8s.io/kube-openapi v0.0.0-20190722073852-5e22f3d471e6 // indirect
	k8s.io/kube-state-metrics v1.7.1 // indirect
	sigs.k8s.io/controller-runtime v0.1.12
	sigs.k8s.io/controller-tools v0.1.12
)

// Pinned to kubernetes-1.13.4
replace (
	k8s.io/api => k8s.io/api v0.0.0-20190222213804-5cb15d344471
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.0.0-20190228180357-d002e88f6236
	k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20190221213512-86fb29eff628
	k8s.io/client-go => k8s.io/client-go v0.0.0-20190228174230-b40b2a5939e4
)

replace (
	github.com/coreos/prometheus-operator => github.com/coreos/prometheus-operator v0.29.0
	k8s.io/kube-state-metrics => k8s.io/kube-state-metrics v1.6.0
	sigs.k8s.io/controller-runtime => sigs.k8s.io/controller-runtime v0.1.12
	sigs.k8s.io/controller-tools => sigs.k8s.io/controller-tools v0.1.11-0.20190411181648-9d55346c2bde
)

replace github.com/operator-framework/operator-sdk => github.com/operator-framework/operator-sdk v0.9.0

replace github.com/appscode/jsonpatch => gomodules.xyz/jsonpatch/v2 v2.0.0
