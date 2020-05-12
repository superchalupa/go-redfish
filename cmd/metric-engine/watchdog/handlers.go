package watchdog

import (
	"context"
	"time"

	"golang.org/x/xerrors"

	eh "github.com/looplab/eventhorizon"
	"github.com/superchalupa/sailfish/src/log"
	"github.com/superchalupa/sailfish/src/looplab/event"
	"github.com/superchalupa/sailfish/src/ocp/am3"
)

const (
	WDEvent              = eh.EventType("Watchdog")
	watchdogsPerInterval = 3
)

type WDEventData struct {
	Seq int
}

type busComponents interface {
	GetBus() eh.EventBus
}

// StartWatchdogHandling will attach event handlers to ping systemd watchdog from AM3 for watchdog events
func StartWatchdogHandling(logger log.Logger, am3Svc am3.Service, d busComponents) error {
	s := 0
	eh.RegisterEventData(WDEvent, func() eh.EventData { s++; return &WDEventData{Seq: s} })

	var sd sdNotifier
	sd, err := NewSdnotify()
	if err != nil {
		logger.Warn("Running using simulation SD_NOTIFY", "err", err)
		sd = SimulateSdnotify()
	}

	// add the watchdog handling to the awesome mapper. meaning that the entire
	// event bus infra has to be working and functional for watchdog to be
	// pinged.
	err = am3Svc.AddEventHandler("Ping Systemd Watchdog", WDEvent, func(eh.Event) {
		err := sd.SDNotify("WATCHDOG=1")
		if err != nil {
			logger.Warn("sdnotify() api failed", "err", err)
		}
	})
	if err != nil {
		return xerrors.Errorf("error adding event handler: %w", err)
	}

	interval := sd.GetIntervalUsec()
	if interval == 0 {
		interval = 30000000
	}
	interval /= watchdogsPerInterval

	// set up separate thread that periodically publishes watchdog events on the bus
	go generateWatchdogEvents(logger, d.GetBus(), time.Duration(interval)*time.Microsecond)
	err = sd.SDNotify("READY=1")
	if err != nil {
		logger.Warn("sdnotify() api failed", "err", err)
	}

	return nil
}

func generateWatchdogEvents(logger log.Logger, bus eh.EventBus, interval time.Duration) {
	logger.Info("Setting up watchdog.", "interval-in-milliseconds", interval)
	watchdogTicker := time.NewTicker(interval)
	defer watchdogTicker.Stop()
	// this runs forever
	for range watchdogTicker.C {
		wdEventData, err := eh.CreateEventData(WDEvent)
		if err != nil {
			logger.Crit("event allocation failed", "err", err)
		}
		event.Publish(context.Background(), bus, WDEvent, wdEventData)
	}
}