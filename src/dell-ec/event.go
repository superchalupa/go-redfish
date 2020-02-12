package dell_ec

import (
	eh "github.com/looplab/eventhorizon"
)

const (
	ComponentEvent    eh.EventType = "ComponentEvent"
	LogEvent          eh.EventType = "LogEvent"
	FaultEntryAdd     eh.EventType = "FaultEntryAdd"
	FaultEntryRemove  eh.EventType = "FaultEntryRemove"
	FaultEntriesClear eh.EventType = "FaultEntriesClear"
)

func init() {
	eh.RegisterEventData(ComponentEvent, func() eh.EventData { return &ComponentEventData{} })
	eh.RegisterEventData(LogEvent, func() eh.EventData { return &LogEventData{} })
	eh.RegisterEventData(FaultEntryAdd, func() eh.EventData { return &FaultEntryAddData{} })
	eh.RegisterEventData(FaultEntryRemove, func() eh.EventData { return &FaultEntryRmData{} })
	eh.RegisterEventData(FaultEntriesClear, func() eh.EventData { return &struct{}{} })
}


type ComponentEventData struct {
	Id         string
	Type       string
	FQDD       string
	ParentFQDD string
}

type LogEventData struct {
	Description string
	EntryType   string
	Id          int
	Created     string
	Message     string
	MessageArgs []string
	MessageID   string
	Name        string
	Severity    string
	Category    string
	Action      string
	FQDD        string
	LogAlert    string
	EventId     string
}

type FaultEntryRmData struct {
	Id   int
	Name string
}

type FaultEntryAddData struct {
	Description   string
	EntryType     string
	Id            int
	Created       string
	Message       string
	MessageArgs   []string
	MessageID     string
	Name          string
	Severity      string
	Category      string
	Action        string
	CompSubSystem string
	SubSystem     string
	FQDD          string
}
