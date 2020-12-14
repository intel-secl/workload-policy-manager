/*
 * Copyright (C) 2019 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package devtest

import (
	config "intel/isecl/wpm/v3/config"
	ci "intel/isecl/wpm/v3/pkg/containerimageflavor"
	i "intel/isecl/wpm/v3/pkg/imageflavor"
	"testing"

	"github.com/stretchr/testify/assert"
)

// assign working KMS URL to "config.Configuration.Kms.APIURL" variable
func TestCreateImageFlavor(t *testing.T) {
	config.Configuration.Kms.APIURL = "https://url.to.kms.instance:port/v1/"
	imageFlavor, err := i.CreateImageFlavor("label", "", "cirros-x86.qcow2", "cirros-x86.qcow2_enc", "", false)
	assert.Nil(t, err)
	assert.NotNil(t, imageFlavor)
}

func TestCreateImageFlavorToFile(t *testing.T) {
	config.Configuration.Kms.APIURL = "https://url.to.kms.instance:port/v1/"
	imageFlavor, err := i.CreateImageFlavor("label", "image_flavor.txt", "cirros-x86.qcow2", "cirros-x86.qcow2_enc", "", false)
	assert.Nil(t, err)
	assert.NotNil(t, imageFlavor)
}

func TestCreateContainerImageFlavor(t *testing.T) {
	config.Configuration.Kms.APIURL = "https://url.to.kms.instance:port/v1/"
	imageFlavor, err := ci.CreateContainerImageFlavor("hello-world", "latest", "", "", "", false, false, "", "")
	assert.Nil(t, err)
	assert.NotNil(t, imageFlavor)
}

func TestCreateContainerImageFlavorToFile(t *testing.T) {
	config.Configuration.Kms.APIURL = "https://url.to.kms.instance:port/v1/"
	imageFlavor, err := ci.CreateContainerImageFlavor("hello-world", "latest", "", "", "", false, false, "", "container_flavor.txt")
	assert.Nil(t, err)
	assert.NotNil(t, imageFlavor)
}
