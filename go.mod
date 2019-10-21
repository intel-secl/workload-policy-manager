module intel/isecl/wpm

require (
	github.com/Gurpartap/logrus-stack v0.0.0-20170710170904-89c00d8a28f4 // indirect
	github.com/facebookgo/stack v0.0.0-20160209184415-751773369052 // indirect
	github.com/google/uuid v1.1.1
	github.com/sirupsen/logrus v1.4.2
	github.com/stretchr/testify v1.3.0
	gopkg.in/yaml.v2 v2.2.2
	intel/isecl/authservice v0.0.0
	intel/isecl/lib/common v1.0.0-Beta
	intel/isecl/lib/flavor v0.0.0
	intel/isecl/lib/kms-client v0.0.0
)

replace intel/isecl/lib/common => gitlab.devtools.intel.com/sst/isecl/lib/common.git v0.0.0-20191004192529-c4a1171a5ec2

replace intel/isecl/lib/flavor => gitlab.devtools.intel.com/sst/isecl/lib/flavor.git v0.0.0-20190915015315-7d9923b58ff3

replace intel/isecl/lib/kms-client => gitlab.devtools.intel.com/sst/isecl/lib/kms-client.git v0.0.0-20191015220541-4091cac1d743

replace intel/isecl/authservice => gitlab.devtools.intel.com/sst/isecl/authservice.git v0.0.0-20191006164949-962eb84a0210
