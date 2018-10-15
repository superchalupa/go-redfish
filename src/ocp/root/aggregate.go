package root

import (
	"context"

	domain "github.com/superchalupa/sailfish/src/redfishresource"

	eh "github.com/looplab/eventhorizon"
)

type view interface {
	GetUUID() eh.UUID
	GetURI() string
}

func AddAggregate(ctx context.Context, v view, ch eh.CommandHandler, eb eh.EventBus) {
	ch.HandleCommand(
		ctx,
		&domain.CreateRedfishResource{
			ID:          v.GetUUID(),
			ResourceURI: v.GetURI(),
			Type:        "#ServiceRoot.v1_0_2.ServiceRoot",
			Context:     "/redfish/v1/$metadata#ServiceRoot.ServiceRoot",

			Privileges: map[string]interface{}{
				"GET": []string{"Unauthenticated"},
			},
			Properties: map[string]interface{}{
				"Id":             "RootService",
				"Name":           "Root Service",
				"Description":    "Root Service",
				"RedfishVersion": "1.0.2",
				"@odata.etag":    `W/"abc123"`,
			}})

	return
}
