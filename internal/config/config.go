package config

import (
	"fmt"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/go-logr/logr"
	"github.com/jessevdk/go-flags"
	"github.com/spf13/viper"

	"github.com/mt-inside/go-usvc"
)

func LoadConfig(log logr.Logger) error {
	var opts struct {
		Verbose          []bool `short:"v" long:"verbose"`
		OutputListenAddr string `long:"output-addr"`
		ApiListenAddr    string `long:"api-addr"`
		ConfigFileDir    string `short:"c" long:"config" description:"Path to the config file's directory. File must be called badpod.yaml"`
		FEDir            string `short:"w" long:"web" description:"Path to the root directory for the web interface"`

		BodyLen         int `short:"l" long:"body-lenght"` // TODO go-flags aNNOTATIONs
		Rate            int64
		DelayMS         int64
		ErrorRate       float64 `short:"e" long:"error-rate"`
		ErrorType       string  `short:"E" long:"error-type"`
		CrashRate       float64 `short:"x" long:"crash-rate" description:"Liklihood of crashing in a given time period."`
		Healthy         bool
		Ready           bool
		ReadyDelayMS    int64
		ShutdownDelayMS int64 `short:"s" long:"shutdown-delay"`
	}

	/* Flags and config */
	_, err := flags.Parse(&opts)
	if err != nil {
		return err
	}

	/* Defaults */
	viper.SetDefault("Verbosity", 0)
	viper.SetDefault("OutputListenAddr", ":8080")
	viper.SetDefault("ApiListenAddr", ":8081")

	viper.SetDefault("Rate", 1000)
	viper.SetDefault("DelayMS", 0)
	viper.SetDefault("ErrorRate", 0.0)
	viper.SetDefault("ErrorType", "http")

	viper.SetDefault("CrashRate", 0.0)
	viper.SetDefault("Healthy", true)
	viper.SetDefault("Ready", false)
	viper.SetDefault("ReadyDelayMS", 5000)

	viper.SetDefault("ShutdownDelayMS", 1000)

	/* Database */
	// TODO

	/* Config files */
	viper.SetConfigName("config") // config file name (minus extension)
	viper.SetConfigType("yaml")   // config file type, if `$Name` is found without an extension
	if viper.IsSet("ConfigFileDir") {
		viper.AddConfigPath(viper.GetString("ConfigFileDir"))
	}
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.config/envbin") // a directory

	// Read
	err = viper.ReadInConfig()
	if err != nil {
		log.Info("Can't read config file", "error", err)
		// continue loading...
	}

	// Watch
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config changed in file", e.Name)
		ConfigChanged()
	})

	/* We don't support the environment */

	/* Command-line flags using go-flags */
	// TODO move to pflags for this project
	if len(opts.Verbose) != 0 {
		viper.Set("Verbosity", len(opts.Verbose))
	}
	if opts.ApiListenAddr != "" {
		viper.Set("ApiListenAddr", opts.ApiListenAddr)
	}
	if opts.OutputListenAddr != "" {
		viper.Set("OutputListenAddr", opts.OutputListenAddr)
	}
	if opts.ErrorRate != 0.0 {
		viper.Set("ErrorRate", opts.ErrorRate)
	}
	if opts.ErrorType != "" {
		viper.Set("ErrorType", opts.ErrorType)
	}
	if opts.CrashRate != 0.0 {
		viper.Set("CrashRate", opts.CrashRate)
	}
	if opts.ShutdownDelayMS != 0 {
		viper.Set("ShutdownDelayMS", opts.ShutdownDelayMS)
	}
	if opts.FEDir != "" {
		viper.Set("FEDir", opts.FEDir)
	}
	// TODO: rest of the settings

	/* Kick lazily-instantiated stuff */
	ConfigChanged() // TODO: frp

	return nil
}

func ConfigChanged() {
	usvc.SetLevel(viper.GetInt("Verbosity"))
}

func ReadyTimer() {
	go func() {
		time.Sleep(time.Duration(viper.GetInt64("ReadyDelayMS")) * time.Millisecond)
		//TODO ideally cancel this if readiness is changed manually
		//TODO ideally save the programme start time, recalcualte absolute deadline for ready if ReadyDelayMS is changed and wait for that new context
		viper.Set("Ready", true)
		ConfigChanged()
	}()
}
