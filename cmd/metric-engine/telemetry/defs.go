package telemetry

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	eh "github.com/looplab/eventhorizon"
	"github.com/superchalupa/sailfish/cmd/metric-engine/metric"
	"github.com/superchalupa/sailfish/src/log"
)

// constants to refer to event types
const (
	// GET - get most redfish URIs.
	GenericGETCommandEvent  eh.EventType = "GenericGETCommandEvent"
	GenericGETResponseEvent eh.EventType = "GenericGETResponseEvent"

	// MRD - Metric Report Definition
	AddMRDCommandEvent     eh.EventType = "AddMetricReportDefinitionEvent"
	AddMRDResponseEvent    eh.EventType = "AddMetricReportDefinitionEventResponse"
	UpdateMRDCommandEvent  eh.EventType = "UpdateMetricReportDefinitionEvent"
	UpdateMRDResponseEvent eh.EventType = "UpdateMetricReportDefinitionEventResponse"
	DeleteMRDCommandEvent  eh.EventType = "DeleteMetricReportDefinitionEvent"
	DeleteMRDResponseEvent eh.EventType = "DeleteMetricReportDefinitionEventResponse"

	// MR - Metric Report
	DeleteMRCommandEvent  eh.EventType = "DeleteMetricReportEvent"
	DeleteMRResponseEvent eh.EventType = "DeleteMetricReportEventResponse"

	// MD - Metric Definition
	AddMDCommandEvent  eh.EventType = "AddMetricDefinitionEvent"
	AddMDResponseEvent eh.EventType = "AddMetricDefinitionEventResponse"

	// Trigger
	CreateTriggerCommandEvent  eh.EventType = "CreateTriggerEvent"
	CreateTriggerResponseEvent eh.EventType = "CreateTriggerEventResponse"
	AddTriggerCommandEvent     eh.EventType = "AddTriggerEvent"
	AddTriggerResponseEvent    eh.EventType = "AddTriggerEventResponse"
	UpdateTriggerCommandEvent  eh.EventType = "UpdateTriggerEvent"
	UpdateTriggerResponseEvent eh.EventType = "UpdateTriggerEventResponse"
	DeleteTriggerCommandEvent  eh.EventType = "DeleteTriggerEvent"
	DeleteTriggerResponseEvent eh.EventType = "DeleteTriggerEventResponse"

	// generic events
	DatabaseMaintenance eh.EventType = "DatabaseMaintenanceEvent"
	PublishClock        eh.EventType = "PublishClockEvent"
)

type GenericGETCommandData struct {
	metric.Command
	URI string
}

type GenericGETResponseData struct {
	metric.CommandResponse
}

func (u *GenericGETCommandData) UseVars(ctx context.Context, logger log.Logger, vars map[string]string) error {
	u.URI = vars["uri"]
	return nil
}

// MetricReportDefinitionData is the eh event data for adding a new report definition
type MetricReportDefinitionData struct {
	Name      string      `db:"Name"      json:"Id"`
	ShortDesc string      `db:"ShortDesc" json:"Name"`
	LongDesc  string      `db:"LongDesc"  json:"Description"`
	Type      string      `db:"Type"      json:"MetricReportDefinitionType"` // 'Periodic', 'OnChange', 'OnRequest'
	Updates   string      `db:"Updates"   json:"ReportUpdates"`              // 'AppendStopsWhenFull', 'AppendWrapsWhenFull', 'NewReport', 'Overwrite'
	Actions   StringArray `db:"Actions"   json:"ReportActions"`              // 	'LogToMetricReportsCollection', 'RedfishEvent'

	// Validation: It's assumed that TimeSpan is parsed on ingress. MRD Schema
	// specifies TimeSpan as a duration.
	// Represents number of seconds worth of metrics in a report. Metrics will be
	// reported from the Report generation as the "End" and metrics must have
	// timestamp > max(End-timespan, report start)
	TimeSpan RedfishDuration `db:"TimeSpan" json:"ReportTimespan"`

	// Validation: It's assumed that Period is parsed on ingress. Redfish
	// "Schedule" object is flexible, but we'll allow only period in seconds for
	// now When it gets to this struct, it needs to be expressed in Seconds.
	Period    RedfishDuration `db:"Period" json:"-"` // when type=periodic, it's a Redfish Duration
	Heartbeat RedfishDuration `db:"HeartbeatInterval" json:"MetricReportHeartbeatInterval"`
	Metrics   []RawMetricMeta `db:"Metrics"`

	Enabled      bool        `db:"Enabled" json:"MetricReportDefinitionEnabled"`
	SuppressDups bool        `db:"SuppressDups" json:"SuppressRepeatedMetricValue"`
	Hidden       bool        `db:"Hidden" json:"-"`
	TriggerList  StringArray `db:"-"         json:"-"` // Trigger [list]
}

// UnmarshalJSON special decoder for MetricReportDefinitionData to unmarshal the "period" specially
func (mrd *MetricReportDefinitionData) UnmarshalJSON(data []byte) error {
	type Alias MetricReportDefinitionData
	target := struct {
		*Alias
		Schedule *struct{ RecurrenceInterval RedfishDuration }
		Links    struct {
			Triggers []struct {
				OdataID string `json:"@odata.id"`
			}
		}
	}{
		Alias: (*Alias)(mrd),
	}

	if err := json.Unmarshal(data, &target); err != nil {
		return err
	}
	if target.Schedule != nil {
		mrd.Period = target.Schedule.RecurrenceInterval
	}
	// TriggerList
	for _, trg := range target.Links.Triggers {
		comps := strings.Split(trg.OdataID, "/")
		mrd.TriggerList = append(mrd.TriggerList, comps[len(comps)-1])
	}

	return nil
}

func (mrd MetricReportDefinitionData) GetTimeSpan() time.Duration {
	return time.Duration(mrd.TimeSpan)
}

// MetricDefifinitionData - is the eh event data for adding a new new definition (future)
type MetricDefinitionData struct {
	MetricID        string      `db:"MetricID"          json:"Id"`
	Name            string      `db:"Name"              json:"Name"`
	Description     string      `db:"Description"       json:"Description"`
	MetricType      string      `db:"MetricType"        json:"MetricType"`
	MetricDataType  string      `db:"MetricDataType"    json:"MetricDataType"`
	Units           string      `db:"Units"             json:"Units"`
	SensingInterval string      `db:"SensingInterval"   json:"SensingInterval"`
	Accuracy        float32     `db:"Accuracy"          json:"Accuracy"`
	DiscreteValues  StringArray `db:"DiscreteValues"   json:"DiscreteValues"`
}

// TriggerData - is the eh event data for adding a new new definition (future)
//   TODO: this definition is only for the current EEMID based event triggers, needs change/addition
//          for other types of trigger like numeric
type TriggerData struct {
	RedfishID      string      `db:"RedfishID"         json:"Id"`
	Name           string      `db:"Name"              json:"Name"`
	Description    string      `db:"Description"       json:"Description"`
	MetricType     string      `db:"MetricType"        json:"MetricType"`
	TriggerActions StringArray `db:"TriggerActions"    json:"TriggerActions"`
	MRDList        StringArray `db:"-"                 json:"-"`             // MRDID [list]
	EventTriggers  StringArray `db:"-"                 json:"EventTriggers"` // EEMIID [list]
}

// UnmarshalJSON special decoder for MetricReportDefinitionData to unmarshal "Links" and "DiscreteTriggers"
func (trg *TriggerData) UnmarshalJSON(data []byte) error {
	type Alias TriggerData
	target := struct {
		*Alias
		Links struct {
			MetricReportDefinitions []struct {
				OdataID string `json:"@odata.id"`
			}
		}
	}{
		Alias: (*Alias)(trg),
	}

	if err := json.Unmarshal(data, &target); err != nil {
		fmt.Println("error json.Unmarshal")
		return err
	}
	// MRD Name list array
	for _, mrd := range target.Links.MetricReportDefinitions {
		comps := strings.Split(mrd.OdataID, "/")
		trg.MRDList = append(trg.MRDList, comps[len(comps)-1])
	}

	return nil
}

// RawMetricMeta is a sub structure to help serialize stuff to db. it containst
// the stuff we are putting in or taking out of DB for Meta.
// Validation: It's assumed that Duration is parsed on ingress. The ingress
// format is (Redfish Duration): -?P(\d+D)?(T(\d+H)?(\d+M)?(\d+(.\d+)?S)?)?
// When it gets to this struct, it needs to be expressed in Seconds.
type RawMetricMeta struct {
	// Meta fields
	MetaID             int64           `db:"MetaID"`
	NamePattern        string          `db:"NamePattern" json:"MetricID"`
	CollectionFunction string          `db:"CollectionFunction" json:"CollectionFunction"`
	CollectionDuration RedfishDuration `db:"CollectionDuration" json:"CollectionDuration"`

	FQDDPattern     string      `db:"FQDDPattern" json:"-"`
	SourcePattern   string      `db:"SourcePattern" json:"-"`
	PropertyPattern string      `db:"PropertyPattern" json:"-"`
	Wildcards       StringArray `db:"Wildcards"  json:"-"`
}

// UnmarshalJSON special decoder for MetricReportDefinitionData to unmarshal the "period" specially
func (m *RawMetricMeta) UnmarshalJSON(data []byte) error {
	type Alias RawMetricMeta
	target := struct {
		*Alias
		Oem struct {
			Dell struct {
				FQDDPattern   string `json:"FQDD"`
				SourcePattern string `json:"Source"`
			} `json:"Dell"`
		} `json:"OEM"`
	}{
		Alias: (*Alias)(m),
	}

	if err := json.Unmarshal(data, &target); err != nil {
		return err
	}
	m.FQDDPattern = target.Oem.Dell.FQDDPattern
	m.SourcePattern = target.Oem.Dell.SourcePattern

	return nil
}

// MetricReportDefinition represents a DB record for a metric report
// definition. Basically adds ID and a static AppendLimit (for now, until we
// can figure out how to make this dynamic).
type MetricReportDefinition struct {
	MetricReportDefinitionData
	AppendLimit int   `db:"AppendLimit"`
	ID          int64 `db:"ID"`
}

type AddMRDCommandData struct {
	metric.Command
	MetricReportDefinitionData
}

func (a *AddMRDCommandData) UseInput(ctx context.Context, logger log.Logger, r io.Reader) error {
	decoder := json.NewDecoder(r)
	return decoder.Decode(&a.MetricReportDefinitionData)
}

type AddMRDResponseData struct {
	metric.CommandResponse
	metric.URIChanged
}

type UpdateMRDCommandData struct {
	metric.Command
	ReportDefinitionName string
	Patch                json.RawMessage
}

type UpdateMRDResponseData struct {
	metric.CommandResponse
	metric.URIChanged
	MetricReportDefinitionData
}

func (u *UpdateMRDCommandData) UseInput(ctx context.Context, logger log.Logger, r io.Reader) error {
	decoder := json.NewDecoder(r)
	return decoder.Decode(&u.Patch)
}

func (u *UpdateMRDCommandData) UseVars(ctx context.Context, logger log.Logger, vars map[string]string) error {
	u.ReportDefinitionName = vars["ID"]
	return nil
}

type DeleteMRDCommandData struct {
	metric.Command
	Name string
}

func (delMRD *DeleteMRDCommandData) UseVars(ctx context.Context, logger log.Logger, vars map[string]string) error {
	delMRD.Name = vars["ID"]
	return nil
}

type DeleteMRDResponseData struct {
	metric.CommandResponse
	metric.URIChanged
}

// MD defs
type MetricDefinition struct {
	MetricDefinitionData
}

type AddMDCommandData struct {
	metric.Command
	MetricDefinitionData
}

func (a *AddMDCommandData) UseInput(ctx context.Context, logger log.Logger, r io.Reader) error {
	decoder := json.NewDecoder(r)
	return decoder.Decode(&a.MetricDefinitionData)
}

type AddMDResponseData struct {
	metric.CommandResponse
}

// Trigger defs
type Trigger struct {
	*TriggerData
	ID int64 `db:"ID"`
}

type CreateTriggerCommandData struct {
	metric.Command
	TriggerData
}

func (a *CreateTriggerCommandData) UseInput(ctx context.Context, logger log.Logger, r io.Reader) error {
	decoder := json.NewDecoder(r)
	return decoder.Decode(a)
}

type CreateTriggerResponseData struct {
	metric.CommandResponse
}

type AddTriggerCommandData struct {
	metric.Command
	TriggerData
}

func (a *AddTriggerCommandData) UseInput(ctx context.Context, logger log.Logger, r io.Reader) error {
	decoder := json.NewDecoder(r)
	return decoder.Decode(&a.TriggerData)
}

type AddTriggerResponseData struct {
	metric.CommandResponse
	metric.URIChanged
}

type UpdateTriggerCommandData struct {
	metric.Command
	TriggerName string
	Patch       json.RawMessage
}

type UpdateTriggerResponseData struct {
	metric.CommandResponse
	metric.URIChanged
}

func (u *UpdateTriggerCommandData) UseInput(ctx context.Context, logger log.Logger, r io.Reader) error {
	decoder := json.NewDecoder(r)
	return decoder.Decode(&u.Patch)
}

func (u *UpdateTriggerCommandData) UseVars(ctx context.Context, logger log.Logger, vars map[string]string) error {
	u.TriggerName = vars["ID"]
	return nil
}

type DeleteTriggerCommandData struct {
	metric.Command
	Name string
}

func (delTrigger *DeleteTriggerCommandData) UseVars(ctx context.Context, logger log.Logger, vars map[string]string) error {
	delTrigger.Name = vars["ID"]
	return nil
}

type DeleteTriggerResponseData struct {
	metric.CommandResponse
	metric.URIChanged
}

// MR defs
type DeleteMRCommandData struct {
	metric.Command
	Name string
}

func (delMR *DeleteMRCommandData) UseVars(ctx context.Context, logger log.Logger, vars map[string]string) error {
	delMR.Name = vars["ID"]
	return nil
}

type DeleteMRResponseData struct {
	metric.CommandResponse
}

func Factory(et eh.EventType) func() (eh.Event, error) {
	return func() (eh.Event, error) {
		data, err := eh.CreateEventData(et)
		if err != nil {
			return nil, fmt.Errorf("could not create request report command: %w", err)
		}
		return eh.NewEvent(et, data, time.Now()), nil
	}
}

func RegisterEvents() {
	CMD := metric.NewCommand // save some verbosity

	// Full commands (request/response)
	// =========== GET - generically get any specific URI =======================
	eh.RegisterEventData(GenericGETCommandEvent, func() eh.EventData {
		return &GenericGETCommandData{Command: CMD(GenericGETResponseEvent)}
	})
	eh.RegisterEventData(GenericGETResponseEvent, func() eh.EventData {
		return &GenericGETResponseData{}
	})

	// =========== ADD MRD - AddMRDCommand ==========================
	eh.RegisterEventData(AddMRDCommandEvent, func() eh.EventData {
		return &AddMRDCommandData{Command: CMD(AddMRDResponseEvent)}
	})
	eh.RegisterEventData(AddMRDResponseEvent, func() eh.EventData { return &AddMRDResponseData{} })

	// =========== UPDATE MRD - UpdateMRDCommandEvent ====================
	eh.RegisterEventData(UpdateMRDCommandEvent, func() eh.EventData {
		return &UpdateMRDCommandData{Command: CMD(UpdateMRDResponseEvent)}
	})
	eh.RegisterEventData(UpdateMRDResponseEvent, func() eh.EventData { return &UpdateMRDResponseData{} })

	// =========== DEL MRD - DeleteMRDCommandEvent =======================
	eh.RegisterEventData(DeleteMRDCommandEvent, func() eh.EventData {
		return &DeleteMRDCommandData{Command: CMD(DeleteMRDResponseEvent)}
	})
	eh.RegisterEventData(DeleteMRDResponseEvent, func() eh.EventData { return &DeleteMRDResponseData{} })

	// =========== DEL MR - DeleteMRCommandEvent ==================================
	eh.RegisterEventData(DeleteMRCommandEvent, func() eh.EventData {
		return &DeleteMRCommandData{Command: CMD(DeleteMRResponseEvent)}
	})
	eh.RegisterEventData(DeleteMRResponseEvent, func() eh.EventData { return &DeleteMRResponseData{} })

	// =========== ADD MD - AddMDCommandEvent ==========================
	eh.RegisterEventData(AddMDCommandEvent, func() eh.EventData {
		return &AddMDCommandData{Command: CMD(AddMDResponseEvent)}
	})
	eh.RegisterEventData(AddMDResponseEvent, func() eh.EventData { return &AddMDResponseData{} })

	// =========== CREATE Trigger - CreateTriggerCommandEvent ==========================
	eh.RegisterEventData(CreateTriggerCommandEvent, func() eh.EventData {
		return &CreateTriggerCommandData{Command: metric.NewCommand(CreateTriggerResponseEvent)}
	})
	eh.RegisterEventData(CreateTriggerResponseEvent, func() eh.EventData { return &CreateTriggerResponseData{} })

	// =========== ADD Trigger - AddTriggerCommand ==========================
	eh.RegisterEventData(AddTriggerCommandEvent, func() eh.EventData {
		return &AddTriggerCommandData{Command: CMD(AddTriggerResponseEvent)}
	})
	eh.RegisterEventData(AddTriggerResponseEvent, func() eh.EventData { return &AddTriggerResponseData{} })

	// =========== UPDATE Trigger - UpdateTriggerCommandEvent ====================
	eh.RegisterEventData(UpdateTriggerCommandEvent, func() eh.EventData {
		return &UpdateTriggerCommandData{Command: CMD(UpdateTriggerResponseEvent)}
	})
	eh.RegisterEventData(UpdateTriggerResponseEvent, func() eh.EventData { return &UpdateTriggerResponseData{} })

	// =========== DEL Trigger - DeleteTriggerCommandEvent =======================
	eh.RegisterEventData(DeleteTriggerCommandEvent, func() eh.EventData {
		return &DeleteTriggerCommandData{Command: CMD(DeleteTriggerResponseEvent)}
	})
	eh.RegisterEventData(DeleteTriggerResponseEvent, func() eh.EventData { return &DeleteTriggerResponseData{} })

	// These aren't planned to ever be commands
	//   - no need for these to be callable from redfish or other interfaces
	eh.RegisterEventData(DatabaseMaintenance, func() eh.EventData { return "" })
}