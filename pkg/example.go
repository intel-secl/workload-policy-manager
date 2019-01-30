package pkg

import (
	"fmt"
	csetup "intel/isecl/lib/common/setup"
)

type SaloneeInfo struct {
}

func (s SaloneeInfo) Run(c csetup.Context) error {
	fmt.Println("This is Salonee")
	return nil
}
func (s SaloneeInfo) Validate(c csetup.Context) error {

	return nil
}
