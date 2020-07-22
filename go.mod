module intel/isecl/wpm/v2

go 1.14

require (
	github.com/google/uuid v1.1.1
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.4.2
	github.com/stretchr/testify v1.3.0
	golang.org/x/sys v0.0.0-20200602225109-6fdc65e7d980 // indirect
	gopkg.in/yaml.v2 v2.2.2
	intel/isecl/lib/clients/v2 v2.2.0
	intel/isecl/lib/common/v2 v2.2.0
	intel/isecl/lib/flavor/v2 v2.2.0
)

replace intel/isecl/lib/flavor/v2 => github.com/intel-secl/flavor/v2 v2.2.0

replace intel/isecl/lib/common/v2 => github.com/intel-secl/common/v2 v2.2.0

replace intel/isecl/lib/clients/v2 => github.com/intel-secl/clients/v2 v2.2.0
