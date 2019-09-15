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

replace intel/isecl/lib/common => github.com/intel-secl/common v1.6-beta

replace intel/isecl/lib/flavor => github.com/intel-secl/flavor v1.6-beta

replace intel/isecl/lib/kms-client => github.com/intel-secl/kms-client v1.6-beta

replace intel/isecl/authservice => github.com/intel-secl/authservice v1.6-beta
