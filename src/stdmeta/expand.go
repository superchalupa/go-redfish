package stdmeta

import (
	"context"
	"errors"
	"fmt"
	"path"
	"reflect"
	"sort"
	"strconv"
	"sync"

	"github.com/spf13/viper"
	"github.com/superchalupa/sailfish/src/dell-resources/attributes"
	"github.com/superchalupa/sailfish/src/log"
	"github.com/superchalupa/sailfish/src/ocp/model"
	"github.com/superchalupa/sailfish/src/ocp/testaggregate"
	"github.com/superchalupa/sailfish/src/ocp/view"
	domain "github.com/superchalupa/sailfish/src/redfishresource"
)

func RegisterFormatters(s *testaggregate.Service, d *domain.DomainObjects) {
	domain.RegisterPlugin(func() domain.Plugin { return &uriCollection{d: d} })

	expandOneFormatter := MakeExpandOneFormatter(d)
	s.RegisterViewFunction("withFormatter_expandone", func(ctx context.Context, logger log.Logger, cfgMgr *viper.Viper, cfgMgrMu *sync.RWMutex, vw *view.View, cfg interface{}, parameters map[string]interface{}) error {
		logger.Debug("Adding expandone formatter to view", "view", vw.GetURI())
		vw.ApplyOption(view.WithFormatter("expandone", expandOneFormatter))
		return nil
	})

	expandFormatter := MakeExpandListFormatter(d)
	s.RegisterViewFunction("withFormatter_expand", func(ctx context.Context, logger log.Logger, cfgMgr *viper.Viper, cfgMgrMu *sync.RWMutex, vw *view.View, cfg interface{}, parameters map[string]interface{}) error {
		logger.Debug("Adding expand formatter to view", "view", vw.GetURI())
		vw.ApplyOption(view.WithFormatter("expand", expandFormatter))
		return nil
	})

	fastExpandFormatter := MakeFastExpandListFormatter(d)
	s.RegisterViewFunction("withFormatter_fastexpand", func(ctx context.Context, logger log.Logger, cfgMgr *viper.Viper, cfgMgrMu *sync.RWMutex, vw *view.View, cfg interface{}, parameters map[string]interface{}) error {
		logger.Debug("Adding fast expand formatter to view", "view", vw.GetURI())
		vw.ApplyOption(view.WithFormatter("fastexpand", fastExpandFormatter))
		return nil
	})

	s.RegisterViewFunction("withFormatter_count", func(ctx context.Context, logger log.Logger, cfgMgr *viper.Viper, cfgMgrMu *sync.RWMutex, vw *view.View, cfg interface{}, parameters map[string]interface{}) error {
		logger.Debug("Adding count formatter to view", "view", vw.GetURI())
		vw.ApplyOption(view.WithFormatter("count", CountFormatter))
		return nil
	})

	s.RegisterViewFunction("withFormatter_attributeFormatter", func(ctx context.Context, logger log.Logger, cfgMgr *viper.Viper, cfgMgrMu *sync.RWMutex, vw *view.View, cfg interface{}, parameters map[string]interface{}) error {
		logger.Debug("Adding attributeFormatter formatter to view", "view", vw.GetURI())
		vw.ApplyOption(view.WithFormatter("attributeFormatter", attributes.FormatAttributeDump))
		return nil
	})

	s.RegisterViewFunction("withFormatter_formatOdataList", func(ctx context.Context, logger log.Logger, cfgMgr *viper.Viper, cfgMgrMu *sync.RWMutex, vw *view.View, cfg interface{}, parameters map[string]interface{}) error {
		logger.Debug("Adding FormatOdataList formatter to view", "view", vw.GetURI())
		vw.ApplyOption(view.WithFormatter("formatOdataList", FormatOdataList))
		return nil
	})

	s.RegisterViewFunction("stdFormatters", func(ctx context.Context, logger log.Logger, cfgMgr *viper.Viper, cfgMgrMu *sync.RWMutex, vw *view.View, cfg interface{}, parameters map[string]interface{}) error {
		logger.Debug("Adding standard formatters (expand, expandone, count, attributeFormatter, formatOdataList) to view", "view", vw.GetURI())
		vw.ApplyOption(view.WithFormatter("expandone", expandOneFormatter))
		vw.ApplyOption(view.WithFormatter("expand", expandFormatter))
		vw.ApplyOption(view.WithFormatter("fastexpand", fastExpandFormatter))
		vw.ApplyOption(view.WithFormatter("count", CountFormatter))
		vw.ApplyOption(view.WithFormatter("attributeFormatter", attributes.FormatAttributeDump))
		vw.ApplyOption(view.WithFormatter("formatOdataList", FormatOdataList))
		return nil
	})
}

type uriCollection struct {
	d *domain.DomainObjects
}

var (
	uriCollectionPlugin = domain.PluginType("uricollection")
)

func (t *uriCollection) PluginType() domain.PluginType { return uriCollectionPlugin }

func (t *uriCollection) PropertyGet(
	ctx context.Context,
	agg *domain.RedfishResourceAggregate,
	auth *domain.RedfishAuthorizationProperty,
	rrp *domain.RedfishResourceProperty,
	meta map[string]interface{},
) error {

	p, ok := meta["uribase"].(string)
	if !ok {
		return errors.New("uribase not specified as a string")
	}

	uriList := t.d.FindMatchingURIs(func(uri string) bool { return path.Dir(uri) == p })
	odata := make([]interface{}, len(uriList))

	sort.Slice(uriList, func(i, j int) bool {
		idx_i, _ := strconv.Atoi(path.Base(uriList[i]))
		idx_j, _ := strconv.Atoi(path.Base(uriList[j]))
		return idx_i > idx_j
	})

	used := 0
	for _, uri := range uriList {
		out, err := t.d.ExpandURI(ctx, uri)
		if err != nil {
			continue
		}
		odata[used] = out
		used = used + 1
	}

	rrp.Value = odata[:used]

	return nil
}

func MakeFastExpandListFormatter(d *domain.DomainObjects) func(context.Context, *view.View, *model.Model, *domain.RedfishResourceAggregate, *domain.RedfishResourceProperty, *domain.RedfishAuthorizationProperty, map[string]interface{}) error {
	return func(
		ctx context.Context,
		v *view.View,
		m *model.Model,
		agg *domain.RedfishResourceAggregate,
		rrp *domain.RedfishResourceProperty,
		auth *domain.RedfishAuthorizationProperty,
		meta map[string]interface{},
	) error {

		u := uriCollection{d: d}
		u.PropertyGet(ctx, agg, auth, rrp, meta)

		count, ok := meta["property"].(string)
		if ok {
			m.UpdateProperty(count, len(rrp.Value.([]interface{})))
		}

		return nil
	}
}

func MakeExpandListFormatter(d *domain.DomainObjects) func(context.Context, *view.View, *model.Model, *domain.RedfishResourceAggregate, *domain.RedfishResourceProperty, *domain.RedfishAuthorizationProperty, map[string]interface{}) error {
	return func(
		ctx context.Context,
		v *view.View,
		m *model.Model,
		agg *domain.RedfishResourceAggregate,
		rrp *domain.RedfishResourceProperty,
		auth *domain.RedfishAuthorizationProperty,
		meta map[string]interface{},
	) error {
		p, ok := meta["property"].(string)

		uris, ok := m.GetPropertyOk(p)
		if !ok {
			uris = []string{}
		}

		// possibly *slightly* larger than needed, but not likely. avoids garbage
		var odata []interface{}

		switch u := uris.(type) {
		case []string:
			odata = make([]interface{}, 0, len(u))
			for _, i := range u {
				out, err := d.ExpandURI(ctx, i)
				if err == nil {
					odata = append(odata, out) // preallocated
				}
			}
		case []interface{}:
			odata = make([]interface{}, 0, len(u))
			for _, i := range u {
				j, ok := i.(string)
				if !ok {
					continue
				}
				out, err := d.ExpandURI(ctx, j)
				if err == nil {
					odata = append(odata, out) // preallocated
				}
			}
		default:
			odata = make([]interface{}, 0)
		}

		rrp.Value = odata

		return nil
	}
}

func MakeExpandOneFormatter(d *domain.DomainObjects) func(context.Context, *view.View, *model.Model, *domain.RedfishResourceAggregate, *domain.RedfishResourceProperty, *domain.RedfishAuthorizationProperty, map[string]interface{}) error {
	return func(
		ctx context.Context,
		v *view.View,
		m *model.Model,
		agg *domain.RedfishResourceAggregate,
		rrp *domain.RedfishResourceProperty,
		auth *domain.RedfishAuthorizationProperty,
		meta map[string]interface{},
	) error {
		p, ok := meta["property"].(string)

		uri, ok := m.GetProperty(p).(string)
		if !ok {
			return errors.New("uri property not setup properly")
		}

		out, err := d.ExpandURI(ctx, uri)
		if err == nil {
			rrp.Value = out
		}

		return nil
	}
}

func CountFormatter(
	ctx context.Context,
	vw *view.View,
	m *model.Model,
	agg *domain.RedfishResourceAggregate,
	rrp *domain.RedfishResourceProperty,
	auth *domain.RedfishAuthorizationProperty,
	meta map[string]interface{},
) error {
	p, ok := meta["property"].(string)
	if !ok {
		return errors.New("property name to operate on not passed in meta")
	}

	arr := m.GetProperty(p)
	if arr == nil {
		rrp.Value = 0
		return errors.New("array property not setup properly, however setting count to 0")
	}

	v := reflect.ValueOf(arr)
	switch v.Kind() {
	case reflect.String:
		rrp.Value = v.Len()
	case reflect.Array:
		rrp.Value = v.Len()
	case reflect.Slice:
		rrp.Value = v.Len()
	case reflect.Map:
		rrp.Value = v.Len()
	case reflect.Chan:
		rrp.Value = v.Len()
	default:
		rrp.Value = 0
	}

	return nil
}

func FormatOdataList(ctx context.Context, v *view.View, m *model.Model, agg *domain.RedfishResourceAggregate, rrp *domain.RedfishResourceProperty, auth *domain.RedfishAuthorizationProperty, meta map[string]interface{}) error {
	p, ok := meta["property"].(string)
	if !ok {
		panic(fmt.Sprintf("Programming error: malformed aggregate. No property specified for agg: %+v", agg))
	}

	// called for missing property or wrong type
	sadpath := func() error {
		m.UnderLock(func() {
			uris, ok := m.GetPropertyOkUnlocked(p)
			if !ok {
				uris = []string{}
			}

			switch u := uris.(type) {
			case []string:
				// we get here if the array isn't sorted. need to sort
				sort.Strings(u)

			case []interface{}:
				// we get here if source is []interface{}, need to convert type to []string
				uriArr := make([]string, 0, len(u))
				for _, i := range u {
					if s, ok := i.(string); ok {
						uriArr = append(uriArr, s) //preallocated
					}
				}
				sort.Strings(uriArr)
				uris = uriArr
			default:
				// we get here if things are really wrong, just reset everything to an empty array
				uris = []string{}
			}

			m.UpdatePropertyUnlocked(p, uris)
			u := uris.([]string) // guaranteeed to be []string here
			odata := make([]interface{}, 0, len(u))
			for _, i := range u {
				odata = append(odata, &domain.RedfishResourceProperty{Value: map[string]interface{}{"@odata.id": i}}) // preallocated
			}
			rrp.Value = odata
		})
		return nil
	}

	// we may have an []interface{} or an []string, if the former, we have to "fix" it by converting from []interface{} to []string
	// we cant fix under read lock, need to take write lock. So, try happy path
	// first, we try to use the scalable, fast Read Lock, then if we have to
	// reformat, take the write lock and redo it

	// check property exists
	uris, ok := m.GetPropertyOk(p)
	if !ok {
		return sadpath()
	}

	// check if it's a []string
	u, ok := uris.([]string)
	if !ok {
		return sadpath()
	}

	// check if it's sorted
	if !sort.StringsAreSorted(u) {
		return sadpath()
	}

	rrp.Value = make([]interface{}, 0, len(u))
	for _, i := range u {
		rrp.Value = append(rrp.Value.([]interface{}), &domain.RedfishResourceProperty{Value: map[string]interface{}{"@odata.id": i}}) // preallocated
	}
	return nil
}
