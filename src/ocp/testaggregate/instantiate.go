package testaggregate

import (
	"context"
	"errors"
	"sync"

	"github.com/Knetic/govaluate"
	eh "github.com/looplab/eventhorizon"
	"github.com/spf13/viper"

	"github.com/superchalupa/sailfish/src/log"
	"github.com/superchalupa/sailfish/src/ocp/model"
	"github.com/superchalupa/sailfish/src/ocp/view"
	domain "github.com/superchalupa/sailfish/src/redfishresource"
)

type viewFunc func(ctx context.Context, logger log.Logger, cfgMgr *viper.Viper, vw *view.View, cfg interface{}, parameters map[string]interface{}) error
type controllerFunc func(ctx context.Context, logger log.Logger, cfgMgr *viper.Viper, vw *view.View, cfg interface{}, parameters map[string]interface{}) error
type aggregateFunc func(ctx context.Context, logger log.Logger, cfgMgr *viper.Viper, vw *view.View, cfg interface{}, parameters map[string]interface{}) ([]eh.Command, error)

type Service struct {
	logger log.Logger
	sync.RWMutex
	ch                          eh.CommandHandler
	viewFunctionsRegistry       map[string]viewFunc
	controllerFunctionsRegistry map[string]controllerFunc
	aggregateFunctionsRegistry  map[string]aggregateFunc
}

func New(logger log.Logger, ch eh.CommandHandler) *Service {
	return &Service{
		logger:                      logger,
		ch:                          ch,
		viewFunctionsRegistry:       map[string]viewFunc{},
		controllerFunctionsRegistry: map[string]controllerFunc{},
		aggregateFunctionsRegistry:  map[string]aggregateFunc{},
	}
}

func (s *Service) RegisterViewFunction(name string, fn viewFunc) {
	s.Lock()
	defer s.Unlock()
	s.viewFunctionsRegistry[name] = fn
}

func (s *Service) GetViewFunction(name string) viewFunc {
	s.RLock()
	defer s.RUnlock()
	return s.viewFunctionsRegistry[name]
}

func (s *Service) RegisterControllerFunction(name string, fn controllerFunc) {
	s.Lock()
	defer s.Unlock()
	s.controllerFunctionsRegistry[name] = fn
}

func (s *Service) GetControllerFunction(name string) controllerFunc {
	s.RLock()
	defer s.RUnlock()
	return s.controllerFunctionsRegistry[name]
}

func (s *Service) RegisterAggregateFunction(name string, fn aggregateFunc) {
	s.Lock()
	defer s.Unlock()
	s.aggregateFunctionsRegistry[name] = fn
}

func (s *Service) GetAggregateFunction(name string) aggregateFunc {
	s.RLock()
	defer s.RUnlock()
	return s.aggregateFunctionsRegistry[name]
}

type config struct {
	Logger      []interface{}
	Models      map[string]map[string]interface{}
	View        []map[string]interface{}
	Controllers []map[string]interface{}
	Aggregate   string
}

func (s *Service) InstantiateFromCfg(ctx context.Context, cfgMgr *viper.Viper, name string, parameters map[string]interface{}) (log.Logger, *view.View, error) {

	newParams := map[string]interface{}{}
	for k, v := range parameters {
		newParams[k] = v
	}

	subCfg := cfgMgr.Sub("views")
	if subCfg == nil {
		s.logger.Crit("missing config file section: 'views'")
		return nil, nil, errors.New("invalid config section 'views'")
	}

	config := config{}

	err := subCfg.UnmarshalKey(name, &config)
	if err != nil {
		s.logger.Crit("unamrshal failed", "err", err, "name", name)
		return nil, nil, errors.New("unmarshal failed")
	}

	// Instantiate logger
	subLogger := s.logger.New(config.Logger...)
	subLogger.Debug("Instantiated new logger")

	// Instantiate view
	vw := view.New(view.WithDeferRegister())
	newParams["uuid"] = vw.GetUUID()

	// Instantiate Models
	for modelName, modelProperties := range config.Models {
		subLogger.Debug("creating model", "name", modelName)
		m := vw.GetModel(modelName)
		if m == nil {
			m = model.New()
		}
		for propName, propValue := range modelProperties {
			propValueStr, ok := propValue.(string)
			if !ok {
				continue
			}
			expr, err := govaluate.NewEvaluableExpressionWithFunctions(propValueStr, functions)
			if err != nil {
				subLogger.Crit("Failed to create evaluable expression", "propValueStr", propValueStr, "err", err)
				continue
			}
			propValue, err := expr.Evaluate(newParams)
			if err != nil {
				subLogger.Crit("expression evaluation failed", "expr", expr, "err", err)
				continue
			}

			subLogger.Debug("setting model property", "propname", propName, "propValue", propValue)
			m.UpdateProperty(propName, propValue)
		}
		vw.ApplyOption(view.WithModel(modelName, m))
	}

	// Run view functions
	for _, viewFn := range config.View {
		viewFnName, ok := viewFn["fn"]
		if !ok {
			subLogger.Crit("Missing function name")
			continue
		}
		viewFnNameStr := viewFnName.(string)
		if !ok {
			subLogger.Crit("fn name isnt a string")
			continue
		}
		viewFnParams, ok := viewFn["params"]
		if !ok {
			subLogger.Crit("Missing function parameters")
			continue
		}
		fn := s.GetViewFunction(viewFnNameStr)
		if fn == nil {
			subLogger.Crit("Could not find registered view function", "function", viewFnNameStr)
			continue
		}
		fn(ctx, subLogger, cfgMgr, vw, viewFnParams, newParams)
	}

	// Instantiate controllers
	for _, controllerFn := range config.Controllers {
		controllerFnName, ok := controllerFn["fn"]
		if !ok {
			subLogger.Crit("Missing function name")
			continue
		}
		controllerFnNameStr := controllerFnName.(string)
		if !ok {
			subLogger.Crit("fn name isnt a string")
			continue
		}
		controllerFnParams, ok := controllerFn["params"]
		if !ok {
			subLogger.Crit("Missing function parameters")
			continue
		}
		fn := s.GetControllerFunction(controllerFnNameStr)
		if fn == nil {
			subLogger.Crit("Could not find registered controller function", "function", controllerFnNameStr)
			continue
		}
		fn(ctx, subLogger, cfgMgr, vw, controllerFnParams, newParams)
	}

	// Instantiate aggregate
	func() {
		if len(config.Aggregate) == 0 {
			subLogger.Debug("no aggregate specified in config file to instantiate.")
			return
		}
		fn := s.GetAggregateFunction(config.Aggregate)
		if fn == nil {
			subLogger.Crit("invalid aggregate function")
			return
		}
		cmds, err := fn(ctx, subLogger, cfgMgr, vw, nil, newParams)
		if err != nil {
			subLogger.Crit("aggregate function returned nil")
			return
		}
		// We can get one or more commands back, handle them
		for _, cmd := range cmds {
			// if it's a resource create command, use the view ID for that
			if c, ok := cmd.(*domain.CreateRedfishResource); ok {
				c.ID = vw.GetUUID()
			}
			s.ch.HandleCommand(ctx, cmd)
		}
	}()

	// register the plugin
	domain.RegisterPlugin(func() domain.Plugin { return vw })
	vw.ApplyOption(view.AtClose(func() { domain.UnregisterPlugin(vw.PluginType()) }))

	return subLogger, vw, nil
}
