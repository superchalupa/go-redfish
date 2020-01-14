package attributes

import (
	"context"
	"errors"
	"sync"

	"github.com/spf13/viper"

	a "github.com/superchalupa/sailfish/src/dell-resources/attributedef"
	"github.com/superchalupa/sailfish/src/log"
	"github.com/superchalupa/sailfish/src/ocp/model"
	"github.com/superchalupa/sailfish/src/ocp/testaggregate"
	"github.com/superchalupa/sailfish/src/ocp/view"
	domain "github.com/superchalupa/sailfish/src/redfishresource"

	eh "github.com/looplab/eventhorizon"
)

func FormatAttributeDump(
	ctx context.Context,
	v *view.View,
	m *model.Model,
	agg *domain.RedfishResourceAggregate,
	rrp *domain.RedfishResourceProperty,
	auth *domain.RedfishAuthorizationProperty,
	meta map[string]interface{},
) error {
	p, ok := meta["property"]
	if !ok {
		return errors.New("fallback")
	}

	prop, ok := p.(string)
	if !ok {
		return errors.New("fallback")
	}

	m.RLock()
	defer m.RUnlock()
	attributes, ok := m.GetPropertyUnlocked(prop).(map[string]map[string]map[string]interface{})
	if !ok {
		rrp.Value = map[string]interface{}{}
		return nil
	}

	var ad a.AttributeData
	res := map[string]interface{}{}
	for group, v := range attributes {
		for index, v2 := range v {
			for name, value := range v2 {
				if ad.ReadAllowed(value, auth) {
					res[group+"."+index+"."+name] = ad.Value
				}
			}
		}
	}
	rrp.Value = res

	return nil
}

func RegisterAggregate(s *testaggregate.Service) {
	s.RegisterAggregateFunction("attributes_uri",
		func(ctx context.Context, subLogger log.Logger, cfgMgr *viper.Viper, cfgMgrMu *sync.RWMutex, vw *view.View, extra interface{}, params map[string]interface{}) ([]eh.Command, error) {
			return []eh.Command{
				&domain.CreateRedfishResource{
					ResourceURI: vw.GetURI(),
					Type:        "#OemAttributes.v1_0_0.OemAttributes",
					Context:     params["rooturi"].(string) + "/$metadata#OemAttributes.OemAttributes",

					Privileges: map[string]interface{}{
						"GET":   []string{"Login"},
						"PATCH": []string{"ConfigureManager"},
					},
					Properties: map[string]interface{}{
						"Id@meta":           vw.Meta(view.GETProperty("unique_name_attr"), view.GETModel("default")),
						"Name":              "Oem Attributes",
						"Description":       "This is the manufacturer/provider specific list of attributes.",
						"AttributeRegistry": "ManagerAttributeRegistry.v1_0_0",
						"Attributes@meta":   vw.Meta(view.GETProperty("attributes"), view.GETFormatter("attributeFormatter"), view.GETModel("default"), view.PropPATCH("attributes", "ar_dump")),
					}},
			}, nil
		})
}
