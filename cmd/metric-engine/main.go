package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"

	eh "github.com/looplab/eventhorizon"
	eventpublisher "github.com/looplab/eventhorizon/publisher/local"
	"github.com/superchalupa/sailfish/src/looplab/eventbus"
	"github.com/superchalupa/sailfish/src/looplab/eventwaiter"

	"github.com/superchalupa/sailfish/src/log15adapter"
)

type busComponents struct {
	EventBus       eh.EventBus
	EventWaiter    *eventwaiter.EventWaiter
	EventPublisher eh.EventPublisher
}

func (d *busComponents) GetBus() eh.EventBus                 { return d.EventBus }
func (d *busComponents) GetWaiter() *eventwaiter.EventWaiter { return d.EventWaiter }
func (d *busComponents) GetPublisher() eh.EventPublisher     { return d.EventPublisher }

func main() {
	flag.StringSliceP("listen", "l", []string{}, "Listen address.  Formats: (http:[ip]:nn, https:[ip]:port)")

	cfgMgr := viper.New()
	if err := cfgMgr.BindPFlags(flag.CommandLine); err != nil {
		panic(fmt.Sprintf("Could not bind viper flags: %s", err))
	}
	cfgMgr.SetEnvPrefix("ME") // set up viper env variable mappings
	cfgMgr.AutomaticEnv()     // automatically pull in overrides from env

	// load main config file
	cfgMgr.SetConfigName("metric-engine")
	cfgMgr.AddConfigPath(".")
	cfgMgr.AddConfigPath("/etc/")
	if err := cfgMgr.ReadInConfig(); err != nil {
		panic(fmt.Sprintf("Could not read config file: %s", err))
	}

	// Local config for running from the build tree
	if fileExists("local-metric-engine.yaml") {
		fmt.Fprintf(os.Stderr, "Reading local-metric-engine.yaml config\n")
		cfgMgr.SetConfigFile("local-metric-engine.yaml")
		if err := cfgMgr.MergeInConfig(); err != nil {
			panic(fmt.Sprintf("Error reading local config file: %s", err))
		}
	}
	flag.Parse() // read command line flags after config files

	logger := log15adapter.Make()
	logger.SetupLogHandlersFromConfig(cfgMgr)
	ctx, cancel := context.WithCancel(context.Background())
	intr := make(chan os.Signal, 1)
	signal.Notify(intr, os.Interrupt)
	go func() {
		// wait until <CTRL>-C
		<-intr
		logger.Crit("INTERRUPTED, Cancelling...")
		cancel()
	}()

	d := &busComponents{
		EventBus:       eventbus.NewEventBus(),
		EventPublisher: eventpublisher.NewEventPublisher(),
		EventWaiter:    eventwaiter.NewEventWaiter(eventwaiter.SetName("Main"), eventwaiter.NoAutoRun),
	}

	d.EventBus.AddHandler(eh.MatchAny(), d.EventPublisher)
	d.EventPublisher.AddObserver(d.EventWaiter)
	go d.GetWaiter().Run()

	setup(ctx, logger, cfgMgr, d)     // set up idrac specific
	h := starthttp(logger, cfgMgr, d) // (optional) startup of http services

	<-ctx.Done() // wait until everything is done
	shutdown()   // shut down idrac stuff
	h.shutdown() // shut down http servers
}

func fileExists(fn string) bool {
	fd, err := os.Stat(fn)
	if os.IsNotExist(err) {
		return false
	}
	return !fd.IsDir()
}
