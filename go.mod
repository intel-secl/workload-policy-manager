module intel/isecl/wpm/v2

require (
	github.com/google/uuid v1.1.1
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.4.2
	github.com/stretchr/testify v1.3.0
	gopkg.in/yaml.v2 v2.2.2
	intel/isecl/lib/clients/v2 v2.0.0
	intel/isecl/lib/common/v2 v2.0.0
	intel/isecl/lib/flavor/v2 v2.0.0
)

replace intel/isecl/lib/flavor/v2 => gitlab.devtools.intel.com/sst/isecl/lib/flavor.git/v2 v2.1/develop

replace intel/isecl/lib/common/v2 => gitlab.devtools.intel.com/sst/isecl/lib/common.git/v2 v2.1/develop

replace intel/isecl/lib/clients/v2 => gitlab.devtools.intel.com/sst/isecl/lib/clients.git/v2 v2.1/develop
