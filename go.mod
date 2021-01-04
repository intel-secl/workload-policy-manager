module intel/isecl/wpm/v3

require (
	github.com/google/uuid v1.1.1
	github.com/intel-secl/intel-secl/v3 v3.3.1
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.4.2
	github.com/stretchr/testify v1.3.0
	golang.org/x/sys v0.0.0-20200602225109-6fdc65e7d980 // indirect
	gopkg.in/yaml.v2 v2.3.0
	intel/isecl/lib/common/v3 v3.3.1
	intel/isecl/lib/flavor/v3 v3.3.1
)

replace intel/isecl/lib/flavor/v3 => gitlab.devtools.intel.com/sst/isecl/lib/flavor.git/v3 v3.3.1/develop

replace intel/isecl/lib/common/v3 => gitlab.devtools.intel.com/sst/isecl/lib/common.git/v3 v3.3.1/develop

replace github.com/vmware/govmomi => github.com/arijit8972/govmomi fix-tpm-attestation-output

replace github.com/intel-secl/intel-secl/v3 => gitlab.devtools.intel.com/sst/isecl/intel-secl.git/v3 v3.3.1/develop
