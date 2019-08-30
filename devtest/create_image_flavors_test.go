/*
 * Copyright (C) 2019 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package imageflavor

import (
	config "intel/isecl/wpm/config"
	ci "intel/isecl/wpm/pkg/containerimageflavor"
	i "intel/isecl/wpm/pkg/imageflavor"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateImageFlavor(t *testing.T) {
	config.Configuration.Kms.APIURL = "https://10.105.168.214:443/v1/"
	config.Configuration.Kms.APIUsername = "kms-admin"
	config.Configuration.Kms.APIPassword = "password"
	config.Configuration.Kms.TLSSha384 = "313f4798df8605b37bf89d68bef596e0a7ce338088a48dd389553d80bb512b76"
	imageFlavor, err := i.CreateImageFlavor("label", "", "cirros-x86.qcow2", "cirros-x86.qcow2_enc", "", false)
	assert.Nil(t, err)
	assert.NotNil(t, imageFlavor)
}

func TestCreateImageFlavorToFile(t *testing.T) {
	config.Configuration.Kms.APIURL = "https://10.105.168.214:443/v1/"
	config.Configuration.Kms.APIUsername = "kms-admin"
	config.Configuration.Kms.APIPassword = "password"
	config.Configuration.Kms.TLSSha384 = "313f4798df8605b37bf89d68bef596e0a7ce338088a48dd389553d80bb512b76"
	imageFlavor, err := i.CreateImageFlavor("label", "image_flavor.txt", "cirros-x86.qcow2", "cirros-x86.qcow2_enc", "", false)
	assert.Nil(t, err)
	assert.NotNil(t, imageFlavor)
}

func TestCreateContainerImageFlavor(t *testing.T) {
	config.Configuration.Kms.APIURL = "https://10.105.168.214:443/v1/"
	config.Configuration.Kms.APIUsername = "kms-admin"
	config.Configuration.Kms.APIPassword = "password"
	config.Configuration.Kms.TLSSha384 = "313f4798df8605b37bf89d68bef596e0a7ce338088a48dd389553d80bb512b76"
	imageFlavor, err := ci.CreateContainerImageFlavor("hello-world", "latest", "", "", "", false, false, "", "")
	assert.Nil(t, err)
	assert.NotNil(t, imageFlavor)
}

func TestCreateContainerImageFlavorToFile(t *testing.T) {
	config.Configuration.Kms.APIURL = "https://10.105.168.214:443/v1/"
	config.Configuration.Kms.APIUsername = "kms-admin"
	config.Configuration.Kms.APIPassword = "password"
	config.Configuration.Kms.TLSSha384 = "313f4798df8605b37bf89d68bef596e0a7ce338088a48dd389553d80bb512b76"
	imageFlavor, err := ci.CreateContainerImageFlavor("hello-world", "latest", "", "", "", false, false, "", "container_flavor.txt")
	assert.Nil(t, err)
	assert.NotNil(t, imageFlavor)
}
