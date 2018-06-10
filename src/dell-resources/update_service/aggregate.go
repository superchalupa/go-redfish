package update_service

// this file should define the BMC Manager object golang data structures where
// we put all the data, plus the aggregate that pulls the data.  actual data
// population should happen in an impl class. ie. no dbus calls in this file

import (
	"context"

	"github.com/superchalupa/go-redfish/src/ocp/view"
	domain "github.com/superchalupa/go-redfish/src/redfishresource"

	eh "github.com/looplab/eventhorizon"
)

func EnhanceAggregate(ctx context.Context, v *view.View, baseUUID eh.UUID, ch eh.CommandHandler) {
	ch.HandleCommand(ctx,
		&domain.UpdateRedfishResourceProperties{
			ID: baseUUID,
			Properties: map[string]interface{}{
				"UpdateService": map[string]interface{}{"@odata.id": v.GetURI()},
			},
		})
}

func AddAggregate(ctx context.Context, v *view.View, ch eh.CommandHandler) *view.View {
	ch.HandleCommand(
		ctx,
		&domain.CreateRedfishResource{
			ID:          v.GetUUID(),
			Collection:  false,
			ResourceURI: v.GetURI(),
			Type:        "#UpdateService.v1_0_0.UpdateService",
			Context:     "/redfish/v1/$metadata#UpdateService.UpdateService",
			Privileges: map[string]interface{}{
				"GET":    []string{"Login"},
				"POST":   []string{}, // cannot create sub objects
				"PUT":    []string{},
				"PATCH":  []string{"ConfigureManager"},
				"DELETE": []string{}, // can't be deleted
			},
			Properties: map[string]interface{}{
				"ServiceEnabled": true,
				"Id":             "UpdateService",
				"Name":           "Update Service",
				"Description":    "Represents the properties for the Update Service",
				"Status": map[string]interface{}{
					"State":  "Enabled",
					"Health": "OK",
				},

				"FirmwareInventory": map[string]interface{}{
					"@odata.id": "/redfish/v1/UpdateService/FirmwareInventory",
				},
				"Actions": map[string]interface{}{
					"Oem": map[string]interface{}{
						"#DellUpdateService.v1_0_0.DellUpdateService.Reset": map[string]interface{}{
							"target": "/redfish/v1/UpdateService/Actions/Oem/DellUpdateService.Reset",
						},
						"UpdateService.v1_0_0#EID_674_UpdateService.Reset": map[string]interface{}{
							"target": "/redfish/v1/UpdateService/Actions/Oem/EID_674_UpdateService.Reset",
						},
						"#DellUpdateService.v1_0_0.DellUpdateService.Syncup": map[string]interface{}{
							"target": "/redfish/v1/UpdateService/Actions/Oem/DellUpdateService.Syncup",
						},
						"UpdateService.v1_0_0#EID_674_UpdateService.Syncup": map[string]interface{}{
							"target": "/redfish/v1/UpdateService/Actions/Oem/EID_674_UpdateService.Syncup",
						},
					},
				},
			}})

	return v
}