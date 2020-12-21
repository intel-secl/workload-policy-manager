module intel/isecl/wpm/v3

require (
	github.com/google/uuid v1.1.1
	github.com/intel-secl/intel-secl/v3 v3.3.0
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.4.2
	github.com/stretchr/testify v1.3.0
	golang.org/x/sys v0.0.0-20200602225109-6fdc65e7d980 // indirect
	gopkg.in/yaml.v2 v2.3.0
	intel/isecl/lib/common/v3 v3.3.0
	intel/isecl/lib/flavor/v3 v3.3.0
)

replace intel/isecl/lib/flavor/v3 => github.com/intel-secl/flavor/v3 v3.3.0

replace intel/isecl/lib/common/v3 => github.com/intel-secl/common/v3 v3.3.0

replace github.com/vmware/govmomi => github.com/arijit8972/govmomi fix-tpm-attestation-output
