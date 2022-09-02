// Package cmd define the root command with the Cobra CLI library
package cmd

import (
	"casperParser/rpc"
	"casperParser/types/config"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/markbates/pkger"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/pflag"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/pkger"
)

const (
	envPrefix = "CASPER_PARSER"
)

var cfgFile string
var redis string
var cluster []string
var sentinel []string
var master string
var rpcEndpoint string
var databaseConnectionString string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "casperParser",
	Short: "casperParser help you parse and store off-chain data from the Casper Blockchain",
	Long:  ``,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return initializeConfig(cmd)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// init define all persistent flags
func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.casperParser or ./casperParser )")
	rootCmd.PersistentFlags().StringVarP(&redis, "redis", "r", "127.0.0.1:6379", "Redis single instance address. Lowest priority over cluster & sentinel flag.")
	rootCmd.PersistentFlags().StringSliceVarP(&cluster, "cluster", "c", []string{}, "Redis cluster instance addresses. Priority over redis flag but not the sentinel flag.")
	rootCmd.PersistentFlags().StringSliceVarP(&sentinel, "sentinel", "s", []string{}, "Redis sentinel addresses. Highest priority over redis & cluster flag.")
	rootCmd.PersistentFlags().StringVarP(&master, "master", "m", "mymaster", "Redis sentinel master name")
	rootCmd.PersistentFlags().StringVar(&rpcEndpoint, "rpc", "http://127.0.0.1:7777/rpc", "Casper RPC endpoint")
	rootCmd.PersistentFlags().StringVarP(&databaseConnectionString, "database", "d", "postgres://postgres:mypassword@localhost:5432/casper", "Postgres connection string. Prefer ENV variable to setup a secure connection")
}

// initializeConfig init Viper and bind flags to viper
func initializeConfig(cmd *cobra.Command) error {
	err := initViper()
	if err != nil {
		return err
	}

	// Bind the current command's flags to viper
	bindFlags(cmd)
	migrateDB()

	return nil
}

// migrateDB apply the database migrations
func migrateDB() {
	pkger.Include("../sql")
	m, err := migrate.New(
		"pkger://../sql",
		databaseConnectionString)
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Up(); err != nil && err.Error() != "no change" {
		log.Fatal(err)
	}
}

// initViper init the Viper singleton
func initViper() error {
	// Set the base name of the config file, without the file extension.
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home & current directory with name ".casperParser" (without extension).
		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName(".casperParser")
	}

	// Attempt to read the config file, gracefully ignoring errors
	// caused by a config file not being found. Return an error
	// if we cannot parse the config file.
	if err := viper.ReadInConfig(); err != nil {
		// It's okay if there isn't a config file
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}
	dt := viper.Get("config")
	err := mapstructure.Decode(dt, &config.ConfigParsed)
	if err != nil {
		log.Fatalf("Can't decode deploys config : %s\n", err)
	}
	log.Printf("Using config file: %s\n", viper.ConfigFileUsed())

	// When we bind flags to environment variables expect that the
	// environment variables are prefixed, e.g. a flag like --number
	// binds to an environment variable STING_NUMBER. This helps
	// avoid conflicts.
	viper.SetEnvPrefix(envPrefix)

	// Bind to environment variables
	// Works great for simple config names, but needs help for names
	// like --favorite-color which we fix in the bindFlags function
	viper.AutomaticEnv()
	return nil
}

// Bind each cobra flag to its associated viper configuration (config file and environment variable)
func bindFlags(cmd *cobra.Command) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		// Environment variables can't have dashes in them, so bind them to their equivalent
		// keys with underscores, e.g. --favorite-color to STING_FAVORITE_COLOR
		if strings.Contains(f.Name, "-") {
			envVarSuffix := strings.ToUpper(strings.ReplaceAll(f.Name, "-", "_"))
			viper.BindEnv(f.Name, fmt.Sprintf("%s_%s", envPrefix, envVarSuffix))
		}

		// Apply the viper config value to the flag when the flag is not set and viper has a value
		if !f.Changed && viper.IsSet(f.Name) {
			val := viper.Get(f.Name)
			cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
		}
	})
}

// getRedisConf from the flags of the command
func getRedisConf(cmd *cobra.Command) asynq.RedisConnOpt {
	var redisConf asynq.RedisConnOpt
	redisConf = asynq.RedisClientOpt{
		Addr: redis,
	}
	if cmd.Flags().Lookup("cluster").Changed {
		redisConf = asynq.RedisClusterClientOpt{
			Addrs: cluster,
		}
	}
	if cmd.Flags().Lookup("sentinel").Changed {
		log.Printf("Redis sentinel master name : %s\n", master)
		redisConf = asynq.RedisFailoverClientOpt{
			MasterName:    master,
			SentinelAddrs: sentinel,
		}
	}
	return redisConf
}

// getRpcClient from the flag of the command
func getRpcClient() *rpc.Client {
	return rpc.NewRpcClient(rpcEndpoint)
}

// getDatabaseConnectionString from the flag of the command
func getDatabaseConnectionString() string {
	return databaseConnectionString
}
