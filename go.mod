module intel/isecl/wpm

require (
	github.com/google/uuid v1.1.1
	github.com/pkg/errors v0.8.1
	github.com/sirupsen/logrus v1.4.2
	github.com/stretchr/testify v1.3.0
	gopkg.in/yaml.v2 v2.2.2
	intel/isecl/lib/clients v0.0.0
	intel/isecl/lib/common v1.0.0-Beta
	intel/isecl/lib/flavor v0.0.0
)

replace intel/isecl/lib/flavor => gitlab.devtools.intel.com/sst/isecl/lib/flavor.git v2.0/develop

replace intel/isecl/lib/common => gitlab.devtools.intel.com/sst/isecl/lib/common.git v2.0/develop

replace intel/isecl/lib/clients => gitlab.devtools.intel.com/sst/isecl/lib/clients.git v2.0/develop
