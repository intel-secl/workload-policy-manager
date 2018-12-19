package config

import (
	"io/ioutil"
	"log"
	"strings"
)

/*
 *
 * @author srege
 *
 */
var Configuration struct {
	BaseURL             string
	Username            string
	Password            string
	KMSTlsCertSHA256    string
	EnvelopeKeyLocation string
}
