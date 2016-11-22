package cmd

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/dainis/gpdownload/download"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
}

var cfgFile string

var RootCmd = &cobra.Command{
	Use:   "gpdownload [flags] <package>",
	Short: "Download apk from playstore",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := &download.Config{
			Email:    viper.GetString("email"),
			Password: viper.GetString("password"),
			DeviceId: viper.GetString("device_id"),
		}

		cfg.OutputFile, _ = cmd.Flags().GetString("output")

		download.Download(cfg, args[0])
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		fields := []string{"email", "password", "device_id"}

		for _, f := range fields {
			v := viper.GetString(f)

			if v == "" {
				log.Fatalf("%s can't be empty", f)
			}
		}

		if len(args) == 0 || args[0] == "" {
			log.Fatal("Specify package to download")
		}
	},
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gpdownload.yaml)")

	RootCmd.Flags().StringP("output", "o", "", "Save download apk as")
	RootCmd.Flags().StringP("email", "e", "", "Google account email")
	RootCmd.Flags().StringP("password", "p", "", "Google account password")
	RootCmd.Flags().StringP("device_id", "d", "", "Device id")

	viper.BindPFlag("email", RootCmd.Flags().Lookup("email"))
	viper.BindPFlag("password", RootCmd.Flags().Lookup("password"))
	viper.BindPFlag("device_id", RootCmd.Flags().Lookup("device_id"))
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName(".gpdownload")
	viper.AddConfigPath("$HOME")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.WithError(err).Fatalf("Failed to read config file %s", viper.ConfigFileUsed())
	} else {
		log.Infof("Using config file: %s", viper.ConfigFileUsed())
	}
}
