package cmd

import (
	"fmt"
	"os"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	ddns "github.com/tnyeanderson/ddns/pkg"
)

type Config struct {
	viper  *viper.Viper
	Agent  ddns.Agent  `mapstructure:",squash"`
	Server ddns.Server `mapstructure:",squash"`
}

func (c *Config) Init() error {
	v := viper.NewWithOptions(
		viper.KeyDelimiter("::"),
	)
	v.SetEnvPrefix("DDNS")
	v.AutomaticEnv()
	c.viper = v
	// Read config if it exists
	if f := os.Getenv(EnvConfigFile); f != "" {
		v.SetConfigFile(f)

		if err := v.ReadInConfig(); err != nil {
			return err
		}

		if err := v.Unmarshal(&c, getDecoder()); err != nil {
			return err
		}

		fmt.Println(v.AllKeys())
	}
	return nil
}

func getDecoder() viper.DecoderConfigOption {
	return viper.DecodeHook(mapstructure.ComposeDecodeHookFunc(
		mapstructure.TextUnmarshallerHookFunc(), // added
		mapstructure.StringToTimeDurationHookFunc(),
		mapstructure.StringToSliceHookFunc(","),
	))

}

func init() {
}
