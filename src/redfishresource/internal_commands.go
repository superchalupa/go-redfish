package domain

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"

	eh "github.com/looplab/eventhorizon"
	"github.com/mitchellh/mapstructure"
	log "github.com/superchalupa/sailfish/src/log"
	"github.com/superchalupa/sailfish/src/looplab/event"
)

type syncEvent interface {
	Add(int)
	Wait()
}

func init() {
	eh.RegisterCommand(func() eh.Command { return &UpdateMetricRedfishResource{} })
	eh.RegisterCommand(func() eh.Command { return &CreateRedfishResource{} })
	eh.RegisterCommand(func() eh.Command { return &RemoveRedfishResource{} })
	eh.RegisterCommand(func() eh.Command { return &UpdateRedfishResourceProperties{} })
	eh.RegisterCommand(func() eh.Command { return &UpdateRedfishResourceProperties2{} })
	eh.RegisterCommand(func() eh.Command { return &RemoveRedfishResourceProperty{} })
	eh.RegisterCommand(func() eh.Command { return &InjectEvent{} })
}

const (
	CreateRedfishResourceCommand                 = eh.CommandType("internal:RedfishResource:Create")
	RemoveRedfishResourceCommand                 = eh.CommandType("internal:RedfishResource:Remove")
	UpdateMetricRedfishResourcePropertiesCommand = eh.CommandType("internal:RedfishResourceProperties:UpdateMetric")
	UpdateRedfishResourcePropertiesCommand       = eh.CommandType("internal:RedfishResourceProperties:Update")
	UpdateRedfishResourcePropertiesCommand2      = eh.CommandType("internal:RedfishResourceProperties:Update:2")
	RemoveRedfishResourcePropertyCommand         = eh.CommandType("internal:RedfishResourceProperties:Remove")
	InjectEventCommand                           = eh.CommandType("internal:Event:Inject")
)

// Static type checking for commands to prevent runtime errors due to typos
var _ = eh.Command(&CreateRedfishResource{})
var _ = eh.Command(&RemoveRedfishResource{})
var _ = eh.Command(&UpdateRedfishResourceProperties{})
var _ = eh.Command(&UpdateRedfishResourceProperties2{})
var _ = eh.Command(&RemoveRedfishResourceProperty{})
var _ = eh.Command(&InjectEvent{})

var immutableProperties = []string{"@odata.id", "@odata.type", "@odata.context"}

// CreateRedfishResource Command
type CreateRedfishResource struct {
	ID          eh.UUID `json:"id"`
	ResourceURI string
	Type        string
	Context     string
	Privileges  map[string]interface{}

	// optional stuff
	Headers       map[string]string      `eh:"optional"`
	Plugin        string                 `eh:"optional"`
	DefaultFilter string                 `eh:"optional"`
	Properties    map[string]interface{} `eh:"optional"`
	Meta          map[string]interface{} `eh:"optional"`
	Private       map[string]interface{} `eh:"optional"`
}

// AggregateType satisfies base Aggregate interface
func (c *CreateRedfishResource) AggregateType() eh.AggregateType { return AggregateType }

// AggregateID satisfies base Aggregate interface
func (c *CreateRedfishResource) AggregateID() eh.UUID { return c.ID }

// CommandType satisfies base Command interface
func (c *CreateRedfishResource) CommandType() eh.CommandType { return CreateRedfishResourceCommand }

func (c *CreateRedfishResource) Handle(ctx context.Context, a *RedfishResourceAggregate) error {

	requestLogger := ContextLogger(ctx, "internal_commands")
	requestLogger.Info("CreateRedfishResource", "META", a.Properties.Meta)

	if a.ID != eh.UUID("") {
		requestLogger.Error("Aggregate already exists!", "command", "CreateRedfishResource", "UUID", a.ID, "URI", a.ResourceURI, "request_URI", c.ResourceURI)
		return errors.New("Already created!")
	}
	a.ID = c.ID
	a.ResourceURI = c.ResourceURI
	a.DefaultFilter = c.DefaultFilter
	a.Plugin = c.Plugin
	a.Headers = make(map[string]string, len(c.Headers))
	for k, v := range c.Headers {
		a.Headers[k] = v
	}

	a.PrivilegeMap = make(map[HTTPReqType]interface{}, len(c.Privileges))
	for k, v := range c.Privileges {
		a.PrivilegeMap[MapStringToHTTPReq(k)] = v
	}

	// ensure no collisions
	for _, p := range immutableProperties {
		delete(c.Properties, p)
	}

	d := &RedfishResourcePropertiesUpdatedData{
		ID:            c.ID,
		ResourceURI:   a.ResourceURI,
		PropertyNames: []string{},
	}
	e := &RedfishResourcePropertyMetaUpdatedData{
		ID:          c.ID,
		ResourceURI: a.ResourceURI,
		Meta:        map[string]interface{}{},
	}

	v := map[string]interface{}{}
	a.Properties.Value = v
	a.Properties.Parse(c.Properties)
	a.Properties.Meta = c.Meta

	var resourceURI []string
	// preserve slashes
	for _, x := range strings.Split(c.ResourceURI, "/") {
		resourceURI = append(resourceURI, url.PathEscape(x))
	}

	v["@odata.id"] = strings.Join(resourceURI, "/")
	v["@odata.type"] = c.Type
	v["@odata.context"] = c.Context

	// send out event that it's created first
	a.PublishEvent(eh.NewEvent(RedfishResourceCreated, &RedfishResourceCreatedData{
		ID:          c.ID,
		ResourceURI: c.ResourceURI,
	}, time.Now()))

	// then send out possible notifications about changes in the properties or meta
	if len(d.PropertyNames) > 0 {
		a.PublishEvent(eh.NewEvent(RedfishResourcePropertiesUpdated, d, time.Now()))
	}
	if len(e.Meta) > 0 {
		a.PublishEvent(eh.NewEvent(RedfishResourcePropertyMetaUpdated, e, time.Now()))
	}

	return nil
}

// RemoveRedfishResource Command
type RemoveRedfishResource struct {
	ID          eh.UUID `json:"id"`
	ResourceURI string  `eh:"optional"`
}

// AggregateType satisfies base Aggregate interface
func (c *RemoveRedfishResource) AggregateType() eh.AggregateType { return AggregateType }

// AggregateID satisfies base Aggregate interface
func (c *RemoveRedfishResource) AggregateID() eh.UUID { return c.ID }

// CommandType satisfies base Command interface
func (c *RemoveRedfishResource) CommandType() eh.CommandType { return RemoveRedfishResourceCommand }

func (c *RemoveRedfishResource) Handle(ctx context.Context, a *RedfishResourceAggregate) error {
	a.PublishEvent(eh.NewEvent(RedfishResourceRemoved, &RedfishResourceRemovedData{
		ID:          c.ID,
		ResourceURI: a.ResourceURI,
	}, time.Now()))
	return nil
}

type RemoveRedfishResourceProperty struct {
	ID       eh.UUID `json:"id"`
	Property string  `eh:"optional"`
}

// AggregateType satisfies base Aggregate interface
func (c *RemoveRedfishResourceProperty) AggregateType() eh.AggregateType { return AggregateType }

// AggregateID satisfies base Aggregate interface
func (c *RemoveRedfishResourceProperty) AggregateID() eh.UUID { return c.ID }

// CommandType satisfies base Command interface
func (c *RemoveRedfishResourceProperty) CommandType() eh.CommandType {
	return RemoveRedfishResourcePropertyCommand
}
func (c *RemoveRedfishResourceProperty) Handle(ctx context.Context, a *RedfishResourceAggregate) error {
	properties := a.Properties.Value.(map[string]interface{})
	for key, _ := range properties {
		if key == c.Property {
			delete(properties, key)
		}
	}
	return nil
}

// toUpdate	{path2key : value}
type UpdateRedfishResourceProperties2 struct {
	ID         eh.UUID `json:"id"`
	Properties map[string]interface{}
}

// AggregateType satisfies base Aggregate interface
func (c *UpdateRedfishResourceProperties2) AggregateType() eh.AggregateType { return AggregateType }

// AggregateID satisfies base Aggregate interface
func (c *UpdateRedfishResourceProperties2) AggregateID() eh.UUID { return c.ID }

// CommandType satisfies base Command interface
func (c *UpdateRedfishResourceProperties2) CommandType() eh.CommandType {
	return UpdateRedfishResourcePropertiesCommand2
}

// aggregate is a.Properties.(RedfishresourceProperty)
// going through the aggregate it is [map]*RedfishResourceProperty...
// Updated to append to list.  TODO need a way to clean lists and prevent duplicates
func UpdateAgg(a *RedfishResourceAggregate, pathSlice []string, v interface{}, appendLimit int) error {
	loc, ok := a.Properties.Value.(map[string]interface{})
	if !ok {
		return errors.New(fmt.Sprintf("Updateagg: aggregate wis wrong type %T", a.Properties.Value))
	}

	plen := len(pathSlice) - 1
	for i, p := range pathSlice {
		k, ok := loc[p]
		if !ok {
			return fmt.Errorf("UpdateAgg Failed can not find %s in %+v", p, loc)
		}
		switch k.(type) {
		case *RedfishResourceProperty:
			k2, ok := k.(*RedfishResourceProperty)
			if !ok {
				return fmt.Errorf("UpdateAgg Failed, RedfishResourcePropertyFailed")
			}
			// metric events have the data appended
			switch v.(type) {
			case []interface{}, []map[string]interface{}:
				j := k2.Value.([]interface{})
				aggSLen := len(j)
				v2 := v.([]interface{})
				if aggSLen >= appendLimit {
					continue
				}
				if appendLimit < aggSLen+len(v2) {
					k2.Parse(v2[appendLimit-aggSLen:])
				} else {
					k2.Parse(v2)
				}
				k2.Parse(j)

				return nil
			default:
				if (plen == i) && (k2.Value != v) {
					k2.Value = v
				} else if plen == i {
					return nil
				} else {
					tmp := k2.Value
					loc, ok = tmp.(map[string]interface{})
					if !ok {
						return fmt.Errorf("UpdateAgg Failed %s type cast to map[string]interface{} for %+v  errored for %+v", a.ResourceURI, p, pathSlice)
					}
				}
			}

		default:
			return fmt.Errorf("agg update for slice %+v, received type %T instead of *RedfishResourceProperty", pathSlice, k)
		}
	}
	return nil

}

func GetValueinAgg(a *RedfishResourceAggregate, pathSlice []string) interface{} {
	a.Properties.Lock()
	defer a.Properties.Unlock()
	loc, ok := a.Properties.Value.(map[string]interface{})
	if !ok {
		return errors.New(fmt.Sprintf("GetValueinAgg: aggregate value is not a map[string]interface{}, but %T", a.Properties.Value))
	}

	plen := len(pathSlice) - 1
	for i, p := range pathSlice {
		k, ok := loc[p]
		if !ok {
			return fmt.Errorf("UpdateAgg Failed can not find %s in %+v", p, loc)
		}
		switch k.(type) {
		case *RedfishResourceProperty:
			k2, ok := k.(*RedfishResourceProperty)
			if !ok {
				return fmt.Errorf("UpdateAgg Failed, RedfishResourcePropertyFailed")
			}
			// metric events have the data appended
			if plen == i {
				return k2.Value
			} else if plen == i {
				return nil
			} else {
				tmp := k2.Value
				loc, ok = tmp.(map[string]interface{})
				if !ok {
					return fmt.Errorf("UpdateAgg Failed %s type cast to map[string]interface{} for %+v  errored for %+v", a.ResourceURI, p, pathSlice)
				}
			}

		default:
			return fmt.Errorf("agg update for slice %+v, received type %T instead of *RedfishResourceProperty", pathSlice, k)
		}
	}

	return nil

}

func validateValue(val interface{}) error {
	switch val.(type) {
	case []interface{}, map[string]interface{}:
		return fmt.Errorf("Update Agg does not support type %T", val)
	default:
		return nil
	}
}

//  This is handled by eventhorizon code.
//  When a CommandHandler "Handle" is called it will retrieve the aggregate from the DB.  and call this Handle. Then save the aggregate 'a' back to the db.  no locking is required..
// provide error when no change made..
func (c *UpdateRedfishResourceProperties2) Handle(ctx context.Context, a *RedfishResourceAggregate) error {

	if a.ID == eh.UUID("") {
		requestLogger := ContextLogger(ctx, "internal_commands")
		requestLogger.Error("Aggregate does not exist!", "UUID", a.ID, "URI", a.ResourceURI)
		return errors.New("non existent aggregate")
	}

	var err error = nil

	d := &RedfishResourcePropertiesUpdatedData2{
		ID:            c.ID,
		ResourceURI:   a.ResourceURI,
		PropertyNames: make(map[string]interface{}),
	}

	// update properties in aggregate
	for k, v := range c.Properties {
		pathSlice := strings.Split(k, "/")

		err := UpdateAgg(a, pathSlice, v, 0)

		if err == nil {
			d.PropertyNames[k] = v
		}
	}

	if len(d.PropertyNames) > 0 {
		a.PublishEvent(eh.NewEvent(RedfishResourcePropertiesUpdated2, d, time.Now()))
	}
	return err
}

type UpdateMetricRedfishResource struct {
	ID               eh.UUID                `json:"id"`
	Properties       map[string]interface{} `eh:"optional"`
	AppendLimit      int
	ReportUpdateType string
}

// AggregateType satisfies base Aggregate interface
func (c *UpdateMetricRedfishResource) AggregateType() eh.AggregateType { return AggregateType }

// AggregateID satisfies base Aggregate interface
func (c *UpdateMetricRedfishResource) AggregateID() eh.UUID { return c.ID }

// CommandType satisfies base Command interface
func (c *UpdateMetricRedfishResource) CommandType() eh.CommandType {
	return UpdateMetricRedfishResourcePropertiesCommand
}

//reportUpdateType int // 0-AppendStopsWhenFull, 1-AppendWrapsWhenFull, 3- NewReport, 4-Overwrite

// assume AppendStopsWhenFull
func (c *UpdateMetricRedfishResource) Handle(ctx context.Context, a *RedfishResourceAggregate) error {
	for k, v := range c.Properties {
		pathSlice := strings.Split(k, "/")
		if err := UpdateAgg(a, pathSlice, v, int(c.AppendLimit)); err != nil {
			fmt.Println("failed to updated agg")
			return err
		}
	}

	return nil
}

type UpdateRedfishResourceProperties struct {
	ID         eh.UUID                `json:"id"`
	Properties map[string]interface{} `eh:"optional"`
}

// AggregateType satisfies base Aggregate interface
func (c *UpdateRedfishResourceProperties) AggregateType() eh.AggregateType { return AggregateType }

// AggregateID satisfies base Aggregate interface
func (c *UpdateRedfishResourceProperties) AggregateID() eh.UUID { return c.ID }

// CommandType satisfies base Command interface
func (c *UpdateRedfishResourceProperties) CommandType() eh.CommandType {
	return UpdateRedfishResourcePropertiesCommand
}
func (c *UpdateRedfishResourceProperties) Handle(ctx context.Context, a *RedfishResourceAggregate) error {
	// ensure no collisions with immutable properties
	for _, p := range immutableProperties {
		delete(c.Properties, p)
	}

	d := &RedfishResourcePropertiesUpdatedData{
		ID:            c.ID,
		ResourceURI:   a.ResourceURI,
		PropertyNames: []string{},
	}
	e := &RedfishResourcePropertyMetaUpdatedData{
		ID:          c.ID,
		ResourceURI: a.ResourceURI,
		Meta:        map[string]interface{}{},
	}

	a.Properties.Parse(c.Properties)

	if len(d.PropertyNames) > 0 {
		a.PublishEvent(eh.NewEvent(RedfishResourcePropertiesUpdated, d, time.Now()))
	}
	if len(e.Meta) > 0 {
		a.PublishEvent(eh.NewEvent(RedfishResourcePropertyMetaUpdated, e, time.Now()))
	}

	return nil
}

type InjectEvent struct {
	ID          eh.UUID                  `json:"id" eh:"optional"`
	Name        eh.EventType             `json:"name"`
	Synchronous bool                     `eh:"optional"`
	Encoding    string                   `eh:"optional" json:"encoding"`
	EventData   map[string]interface{}   `json:"data" eh:"optional"`
	EventArray  []map[string]interface{} `json:"event_array" eh:"optional"`
	EventSeq    int64                    `json:"event_seq" eh:"optional"`

	ctx context.Context
}

// AggregateType satisfies base Aggregate interface
func (c *InjectEvent) AggregateType() eh.AggregateType { return AggregateType }

// AggregateID satisfies base Aggregate interface
func (c *InjectEvent) AggregateID() eh.UUID { return c.ID }

// CommandType satisfies base Command interface
func (c *InjectEvent) CommandType() eh.CommandType {
	return InjectEventCommand
}

var injectChanSlice chan *InjectEvent
var injectChan chan eh.Event

// inject event timeout
var IETIMEOUT time.Duration = 250 * time.Millisecond

func StartInjectService(logger log.Logger, d *DomainObjects) {
	injectChanSlice = make(chan *InjectEvent, 100)
	injectChan = make(chan eh.Event, 10)
	logger = logger.New("module", "injectservice")
	eb := d.EventBus
	ew := d.EventWaiter

	var s closeNotifier
	s, err := NewSdnotify()
	if err != nil {
		fmt.Printf("Error setting up SD_NOTIFY, using simulation instead: %s\n", err)
		s = SimulateSdnotify()
	}

	if interval := s.GetIntervalUsec(); interval == 0 {
		fmt.Printf("Watchdog interval is not set, so skipping watchdog setup. Set WATCHDOG_USEC to set.\n")
	} else {
		fmt.Printf("Setting up watchdog\n")

		// send watchdogs 3x per interval
		interval = interval / 3

		// set up listener for the watchdog events
		listener, err := ew.Listen(context.Background(), func(event eh.Event) bool {
			if event.EventType() == WatchdogEvent {
				return true
			}
			return false
		})

		if err != nil {
			panic("Could not start listener")
		}

		// goroutine to run sd_notify whenever we see a watchdog event
		go func() {
			defer s.Close()
			for {
				_, err := listener.Wait(context.Background())
				if err != nil {
					fmt.Printf("Watchdog wait exited\n")
					break
				}

				s.Notify("WATCHDOG=1")
			}
		}()

		// goroutine to inject watchdog events
		go func() {
			// inject a watchdog event every 10s. It will be processed by a listener elsewhere.
			for {
				time.Sleep(time.Duration(interval) * time.Microsecond)
				data, err := eh.CreateEventData("WatchdogEvent")
				if err != nil {
					injectChan <- eh.NewEvent(WatchdogEvent, data, time.Now())
				}

			}
		}()
	}

	// goroutine to synchronously handle the event inject queue
	go func() {
		queued := []*InjectEvent{}
		internalSeq := 0
		// find better way to initialize sequenceTimer
		sequenceTimer := time.NewTimer(IETIMEOUT)
		sequenceTimer.Stop()
		missingEvent := false
		tries := 0
		for {
			select {
			case event := <-injectChanSlice:
				//logger.Crit("InjectService Event received", "Sequence", event.EventSeq, "Name", event.Name)
				queued = append(queued, event)

				// ordered events are processed.
				for len(queued) != 0 && missingEvent == false {
					for _, evtPtr := range queued {
						// reset sailfish event seq counter
						if evtPtr.EventSeq == -1 {
							evtPtr.EventSeq = 0
							internalSeq = 0
						}
						eventSeq := int(evtPtr.EventSeq)
						//logger.Crit("InjectService: Event start", "Event Name", evtPtr.Name, "Sequence Number", evtPtr.EventSeq, "expected", internalSeq+1)

						if (eventSeq == internalSeq+1) || (internalSeq == 0 && eventSeq == 0) {
							// process event
							evtPtr.sendToChn(evtPtr.ctx)
							internalSeq = eventSeq
							queued[0] = nil
							queued = queued[1:]
						} else if internalSeq >= eventSeq {
							// drop all old events
							dropped_event := &DroppedEventData{
								Name:     evtPtr.Name,
								EventSeq: evtPtr.EventSeq,
							}

							queued[0] = nil
							queued = queued[1:]
							logger.Crit("InjectService: Event dropped", "Event Name", evtPtr.Name, "Sequence Number", evtPtr.EventSeq, "expected", internalSeq+1)
							eb.PublishEvent(evtPtr.ctx, eh.NewEvent(DroppedEvent, dropped_event, time.Now()))
						} else {
							tries += 1
							// missing event found, break and stop for loop
							missingEvent = true
							sequenceTimer = time.NewTimer(IETIMEOUT)
							break
						}

					}
				}

			// IETIMEOUT triggered here I will change the current sequence number and let the rest be handled above.
			case <-sequenceTimer.C:
				if len(queued) == 0 {
					sequenceTimer = time.NewTimer(IETIMEOUT)
					continue
				}

				sort.SliceStable(queued, func(i, j int) bool {
					return queued[i].EventSeq < queued[j].EventSeq
				})

				eventSeq := int(queued[0].EventSeq)
				// event sequence jumped!
				if eventSeq > internalSeq+1 {
					if tries < 40 {
						tries += 1
						sequenceTimer = time.NewTimer(IETIMEOUT)
						continue
					}
					logger.Crit("InjectService: Changing Internal Event Sequence", "# events in queue", len(queued), "before", internalSeq, "after", eventSeq-1)
				}

				tries = 0
				missingEvent = false
				if eventSeq > internalSeq {
					internalSeq = eventSeq - 1
				}
			}
		}
	}()

	go func() {
		for {
			event := <-injectChan

			eb.PublishEvent(context.Background(), event)

			ev, ok := event.(syncEvent)
			if event.EventType() == "LogEvent" {
				// do nothing
			} else if ok {
				ev.Wait()
			}

		}
	}()

}

const MAX_CONSOLIDATED_EVENTS = 5
const injectUUID = eh.UUID("49467bb4-5c1f-473b-0000-00000000000f")

type Decoder interface {
	Decode(d map[string]interface{}) error
}

func (c *InjectEvent) Handle(ctx context.Context, a *RedfishResourceAggregate) error {
	//testLogger := ContextLogger(ctx, "internal_commands").New("module", "inject_event")
	//testLogger.Crit("Event handle", "Sequence", c.EventSeq, "Name", c.Name)
	a.ID = injectUUID
	c.ctx = ctx

	if c.Synchronous {
		c.sendToChn(c.ctx)

	} else {
		injectChanSlice <- c
	}

	return nil
}

func (c *InjectEvent) sendToChn(ctx context.Context) error {

	requestLogger := ContextLogger(ctx, "internal_commands").New("module", "inject_event")
	//requestLogger.Crit("InjectService: Event sent", "Sequence", c.EventSeq, "Name", c.Name)

	eventList := make([]map[string]interface{}, 0, len(c.EventArray)+1)
	if len(c.EventData) > 0 {
		// comment out debug prints in the hot path, uncomment for debugging
		requestLogger.Debug("InjectEvent - ONE", "events", c.EventData)
		eventList = append(eventList, c.EventData) // preallocated
	}
	if len(c.EventArray) > 0 {
		// comment out debug prints in the hot path, uncomment for debugging
		requestLogger.Debug("InjectEvent - ARRAY", "events", c.EventArray)
		eventList = append(eventList, c.EventArray...) // preallocated
	}

	trainload := make([]eh.EventData, 0, MAX_CONSOLIDATED_EVENTS)
	for _, eventData := range eventList {
		// limit number of consolidated events to 30 to prevent overflowing queues and deadlocking
		if len(trainload) >= MAX_CONSOLIDATED_EVENTS {
			//fmt.Printf("Train (%s) leaving early: %d\n", c.Name, len(trainload))
			e := event.NewSyncEvent(c.Name, trainload, time.Now())
			e.Add(1)
			if c.Synchronous {
				defer e.Wait()
			}
			trainload = make([]eh.EventData, 0, MAX_CONSOLIDATED_EVENTS)
			injectChan <- e
		}

		// prefer to deserialize directly to a named type
		data, err := eh.CreateEventData(c.Name)

		// if the named type is not available, publish raw map[string]interface{} as eventData
		if err != nil {
			// this debug statement probably not hit too often, leave enabled for now
			// This is not the preferred path. Consider creating event if we hit this for specific events.
			requestLogger.Info("InjectEvent - event type not registered: injecting raw event.", "event name", c.Name, "error", err)
			trainload = append(trainload, eventData) //preallocated
			continue
		}

		// check if event wants to deserialize itself with a custom decoder
		if ds, ok := data.(Decoder); ok {
			err = ds.Decode(eventData)
			if err != nil {
				fmt.Printf("binary decode fail: %s\n", err)
				continue
			}
			trainload = append(trainload, data) //preallocated
			continue
		}

		// otherwise use default
		if c.Encoding == "binary" {
			structdata, err := base64.StdEncoding.DecodeString(eventData["data"].(string))
			if err != nil {
				fmt.Printf("ERROR decoding base64 event data: %s", err)
				continue
			}

			buf := bytes.NewReader(structdata)
			err = binary.Read(buf, binary.LittleEndian, data)
			if err != nil {
				fmt.Printf("binary decode fail: %s\n", err)
				continue
			}

			trainload = append(trainload, buf) //preallocated
		} else if c.Encoding == "json" || c.Encoding == "" {
			err = mapstructure.Decode(eventData, &data)
			if err != nil {
				trainload = append(trainload, eventData) //preallocated
				requestLogger.Warn("InjectEvent - could not decode event data, skipping event", "error", err, "raw-eventdata", eventData, "dest-event", data)
			} else {
				trainload = append(trainload, data) //preallocated
			}
		}

		// limit number of consolidated events to 30 to prevent overflowing queues and deadlocking
		if len(trainload) >= MAX_CONSOLIDATED_EVENTS {
			e := event.NewSyncEvent(c.Name, trainload, time.Now())
			trainload = make([]eh.EventData, 0, MAX_CONSOLIDATED_EVENTS)
			e.Add(1)
			if c.Synchronous {
				defer e.Wait()
			}

			injectChan <- e

		}

	}

	if len(trainload) > 0 {
		e := event.NewSyncEvent(c.Name, trainload, time.Now())
		if c.Synchronous {
			defer e.Wait()
		}
		e.Add(1)
		injectChan <- e

	}

	return nil
}
