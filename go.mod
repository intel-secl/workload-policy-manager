module intel/isecl/wpm

require (
	github.com/Gurpartap/logrus-stack v0.0.0-20170710170904-89c00d8a28f4 // indirect
	github.com/facebookgo/stack v0.0.0-20160209184415-751773369052 // indirect
	github.com/google/uuid v1.1.1
	github.com/sirupsen/logrus v1.4.0
	github.com/stretchr/testify v1.3.0
	gopkg.in/yaml.v2 v2.2.2
	intel/isecl/lib/common v0.0.0
	intel/isecl/lib/flavor v0.0.0
	intel/isecl/lib/kms-client v0.0.0
)

replace intel/isecl/lib/common => gitlab.devtools.intel.com/sst/isecl/lib/common.git v0.0.0-20190628222117-8960a9b04b5c

replace intel/isecl/lib/flavor => gitlab.devtools.intel.com/sst/isecl/lib/flavor.git v0.0.0-20190221164143-ac584a10db65

replace intel/isecl/lib/kms-client => gitlab.devtools.intel.com/sst/isecl/lib/kms-client.git v0.0.0-20190205015703-57d89d60bd9d
