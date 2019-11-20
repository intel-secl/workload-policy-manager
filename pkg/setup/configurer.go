/*
 * Copyright (C) 2019 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package setup

import (
	csetup "intel/isecl/lib/common/setup"
	"intel/isecl/wpm/config"

	"github.com/pkg/errors"
)

type Configurer struct {
}

var ErrMessageSetupIncomplete = errors.New("Configuration is not complete - setup tasks can be completed only after configuration")

func (cnfr Configurer) Run(c csetup.Context) error {
	log.Trace("setup/configurer:Run() Entering")
	defer log.Trace("setup/configurer:Run() Leaving")
	// save configuration from config.yml
	if cnfr.Validate(c) == nil {
		log.Debug("setup/configurer:Run() Configurer setup task is complete")
		return nil
	}

	err := config.SaveConfiguration(c)
	if err != nil {
		return err
	}

	return nil
}

func (cnfr Configurer) Validate(c csetup.Context) error {
	log.Trace("setup/configurer:Validate() Entering")
	defer log.Trace("setup/configurer:Validate() Leaving")

	if config.Configuration.ConfigComplete != true {
		return errors.New("setup/configurer:Validate() Configuration is not complete")
	}
	return nil
}
