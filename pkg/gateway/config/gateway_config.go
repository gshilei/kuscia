// Copyright 2023 Ant Group Co., Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"

	"github.com/secretflow/kuscia/pkg/gateway/utils"
	"github.com/secretflow/kuscia/pkg/utils/kusciaconfig"
	"github.com/secretflow/kuscia/pkg/utils/nlog"
)

type GatewayConfig struct {
	RootDir       string `yaml:"rootdir,omitempty"`
	Namespace     string `yaml:"namespace,omitempty"`
	ConfBasedir   string `yaml:"confBasedir,omitempty"`
	DomainKeyFile string `yaml:"domainKeyFile,omitempty"`
	WhiteListFile string `yaml:"whiteListFile,omitempty"`

	ExternalPort   uint32 `yaml:"externalPort,omitempty"`
	HandshakePort  uint32 `yaml:"handshakePort,omitempty"`
	XDSPort        uint32 `yaml:"xdsPort,omitempty"`
	EnvoyAdminPort uint32 `yaml:"envoyAdminPort,omitempty"`

	IdleTimeout  int `yaml:"idleTimeout,omitempty"`
	ResyncPeriod int `yaml:"resyncPeriod,omitempty"`

	MasterConfig   *kusciaconfig.MasterConfig `yaml:"master,omitempty"`
	ExternalTLS    *kusciaconfig.TLSConfig    `yaml:"externalTLS,omitempty"`
	InnerServerTLS *kusciaconfig.TLSConfig    `yaml:"InnerServerTLS,omitempty"`
	InnerClientTLS *kusciaconfig.TLSConfig    `yaml:"InnerClientTLS,omitempty"`
}

func DefaultStaticGatewayConfig() *GatewayConfig {
	g := &GatewayConfig{
		Namespace:     "default",
		ConfBasedir:   "./conf",
		DomainKeyFile: "",
		WhiteListFile: "",

		ExternalPort:   1080,
		HandshakePort:  1054,
		XDSPort:        10001,
		EnvoyAdminPort: 10000,
		IdleTimeout:    60,
		ResyncPeriod:   600,
		MasterConfig:   &kusciaconfig.MasterConfig{},
	}
	return g
}

func LoadOverrideConfig(config *GatewayConfig, configPath string) (*GatewayConfig, error) {
	if configPath == "" {
		return config, nil // no need to load config file
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal(data, config)
	if err != nil {
		return config, err
	}

	nlog.Infof("Gateway config: %+v", config)

	return config, config.CheckConfig()
}

func (config *GatewayConfig) CheckConfig() error {
	var err error
	err = kusciaconfig.CheckTLSConfig(config.InnerServerTLS, "innerServerTLS")
	if err != nil {
		return err
	}

	err = kusciaconfig.CheckTLSConfig(config.ExternalTLS, "externalTLS")
	if err != nil {
		return err
	}

	err = kusciaconfig.CheckTLSConfig(config.InnerClientTLS, "innerClientTLS")
	if err != nil {
		return err
	}

	return kusciaconfig.CheckMasterConfig(config.MasterConfig)
}

func (config *GatewayConfig) GetEnvoyNodeID() string {
	hostname := utils.GetHostname()
	envoyNodeCluster := fmt.Sprintf("kuscia-gateway-%s", config.Namespace)
	return fmt.Sprintf("%s-%s", envoyNodeCluster, hostname)
}