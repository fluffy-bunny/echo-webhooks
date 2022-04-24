// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/cheekybits/genny

package claimsprovider

import (
	"reflect"
	"strings"

	di "github.com/fluffy-bunny/sarulabsdi"
	"github.com/rs/zerolog/log"
)

// ReflectTypeIClaimsProvider used when your service claims to implement IClaimsProvider
var ReflectTypeIClaimsProvider = di.GetInterfaceReflectType((*IClaimsProvider)(nil))

// AddSingletonIClaimsProvider adds a type that implements IClaimsProvider
func AddSingletonIClaimsProvider(builder *di.Builder, implType reflect.Type, implementedTypes ...reflect.Type) {
	implementedTypes = append(implementedTypes, ReflectTypeIClaimsProvider)
	_logAddIClaimsProvider("SINGLETON", implType, _getImplementedIClaimsProviderNames(implementedTypes...),
		_logIClaimsProviderExtra{
			Name:  "DI-BY",
			Value: "type",
		})
	di.AddSingleton(builder, implType, implementedTypes...)
}

// AddSingletonIClaimsProviderWithMetadata adds a type that implements IClaimsProvider
func AddSingletonIClaimsProviderWithMetadata(builder *di.Builder, implType reflect.Type, metaData map[string]interface{}, implementedTypes ...reflect.Type) {
	implementedTypes = append(implementedTypes, ReflectTypeIClaimsProvider)
	_logAddIClaimsProvider("SINGLETON", implType, _getImplementedIClaimsProviderNames(implementedTypes...),
		_logIClaimsProviderExtra{
			Name:  "DI-BY",
			Value: "type",
		},
		_logIClaimsProviderExtra{
			Name:  "DI-M",
			Value: metaData,
		})
	di.AddSingletonWithMetadata(builder, implType, metaData, implementedTypes...)
}

// AddSingletonIClaimsProviderByObj adds a prebuilt obj
func AddSingletonIClaimsProviderByObj(builder *di.Builder, obj interface{}, implementedTypes ...reflect.Type) {
	implementedTypes = append(implementedTypes, ReflectTypeIClaimsProvider)
	_logAddIClaimsProvider("SINGLETON", reflect.TypeOf(obj), _getImplementedIClaimsProviderNames(implementedTypes...),
		_logIClaimsProviderExtra{
			Name:  "DI-BY",
			Value: "obj",
		})
	di.AddSingletonWithImplementedTypesByObj(builder, obj, implementedTypes...)
}

// AddSingletonIClaimsProviderByObjWithMetadata adds a prebuilt obj
func AddSingletonIClaimsProviderByObjWithMetadata(builder *di.Builder, obj interface{}, metaData map[string]interface{}, implementedTypes ...reflect.Type) {
	implementedTypes = append(implementedTypes, ReflectTypeIClaimsProvider)
	_logAddIClaimsProvider("SINGLETON", reflect.TypeOf(obj), _getImplementedIClaimsProviderNames(implementedTypes...),
		_logIClaimsProviderExtra{
			Name:  "DI-BY",
			Value: "obj",
		},
		_logIClaimsProviderExtra{
			Name:  "DI-M",
			Value: metaData,
		})

	di.AddSingletonWithImplementedTypesByObjWithMetadata(builder, obj, metaData, implementedTypes...)
}

// AddSingletonIClaimsProviderByFunc adds a type by a custom func
func AddSingletonIClaimsProviderByFunc(builder *di.Builder, implType reflect.Type, build func(ctn di.Container) (interface{}, error), implementedTypes ...reflect.Type) {
	implementedTypes = append(implementedTypes, ReflectTypeIClaimsProvider)
	_logAddIClaimsProvider("SINGLETON", implType, _getImplementedIClaimsProviderNames(implementedTypes...),
		_logIClaimsProviderExtra{
			Name:  "DI-BY",
			Value: "func",
		})
	di.AddSingletonWithImplementedTypesByFunc(builder, implType, build, implementedTypes...)
}

// AddSingletonIClaimsProviderByFuncWithMetadata adds a type by a custom func
func AddSingletonIClaimsProviderByFuncWithMetadata(builder *di.Builder, implType reflect.Type, build func(ctn di.Container) (interface{}, error), metaData map[string]interface{}, implementedTypes ...reflect.Type) {
	implementedTypes = append(implementedTypes, ReflectTypeIClaimsProvider)
	_logAddIClaimsProvider("SINGLETON", implType, _getImplementedIClaimsProviderNames(implementedTypes...),
		_logIClaimsProviderExtra{
			Name:  "DI-BY",
			Value: "func",
		},
		_logIClaimsProviderExtra{
			Name:  "DI-M",
			Value: metaData,
		})

	di.AddSingletonWithImplementedTypesByFuncWithMetadata(builder, implType, build, metaData, implementedTypes...)
}

// AddTransientIClaimsProvider adds a type that implements IClaimsProvider
func AddTransientIClaimsProvider(builder *di.Builder, implType reflect.Type, implementedTypes ...reflect.Type) {
	implementedTypes = append(implementedTypes, ReflectTypeIClaimsProvider)
	_logAddIClaimsProvider("TRANSIENT", implType, _getImplementedIClaimsProviderNames(implementedTypes...),
		_logIClaimsProviderExtra{
			Name:  "DI-BY",
			Value: "type",
		})

	di.AddTransientWithImplementedTypes(builder, implType, implementedTypes...)
}

// AddTransientIClaimsProviderWithMetadata adds a type that implements IClaimsProvider
func AddTransientIClaimsProviderWithMetadata(builder *di.Builder, implType reflect.Type, metaData map[string]interface{}, implementedTypes ...reflect.Type) {
	implementedTypes = append(implementedTypes, ReflectTypeIClaimsProvider)
	_logAddIClaimsProvider("TRANSIENT", implType, _getImplementedIClaimsProviderNames(implementedTypes...),
		_logIClaimsProviderExtra{
			Name:  "DI-BY",
			Value: "type",
		},
		_logIClaimsProviderExtra{
			Name:  "DI-M",
			Value: metaData,
		})

	di.AddTransientWithImplementedTypesWithMetadata(builder, implType, metaData, implementedTypes...)
}

// AddTransientIClaimsProviderByFunc adds a type by a custom func
func AddTransientIClaimsProviderByFunc(builder *di.Builder, implType reflect.Type, build func(ctn di.Container) (interface{}, error), implementedTypes ...reflect.Type) {
	implementedTypes = append(implementedTypes, ReflectTypeIClaimsProvider)
	_logAddIClaimsProvider("TRANSIENT", implType, _getImplementedIClaimsProviderNames(implementedTypes...),
		_logIClaimsProviderExtra{
			Name:  "DI-BY",
			Value: "func",
		})

	di.AddTransientWithImplementedTypesByFunc(builder, implType, build, implementedTypes...)
}

// AddTransientIClaimsProviderByFuncWithMetadata adds a type by a custom func
func AddTransientIClaimsProviderByFuncWithMetadata(builder *di.Builder, implType reflect.Type, build func(ctn di.Container) (interface{}, error), metaData map[string]interface{}, implementedTypes ...reflect.Type) {
	implementedTypes = append(implementedTypes, ReflectTypeIClaimsProvider)
	_logAddIClaimsProvider("TRANSIENT", implType, _getImplementedIClaimsProviderNames(implementedTypes...),
		_logIClaimsProviderExtra{
			Name:  "DI-BY",
			Value: "func",
		},
		_logIClaimsProviderExtra{
			Name:  "DI-M",
			Value: metaData,
		})

	di.AddTransientWithImplementedTypesByFuncWithMetadata(builder, implType, build, metaData, implementedTypes...)
}

// AddScopedIClaimsProvider adds a type that implements IClaimsProvider
func AddScopedIClaimsProvider(builder *di.Builder, implType reflect.Type, implementedTypes ...reflect.Type) {
	implementedTypes = append(implementedTypes, ReflectTypeIClaimsProvider)
	_logAddIClaimsProvider("SCOPED", implType, _getImplementedIClaimsProviderNames(implementedTypes...),
		_logIClaimsProviderExtra{
			Name:  "DI-BY",
			Value: "type",
		})
	di.AddScopedWithImplementedTypes(builder, implType, implementedTypes...)
}

// AddScopedIClaimsProviderWithMetadata adds a type that implements IClaimsProvider
func AddScopedIClaimsProviderWithMetadata(builder *di.Builder, implType reflect.Type, metaData map[string]interface{}, implementedTypes ...reflect.Type) {
	implementedTypes = append(implementedTypes, ReflectTypeIClaimsProvider)
	_logAddIClaimsProvider("SCOPED", implType, _getImplementedIClaimsProviderNames(implementedTypes...),
		_logIClaimsProviderExtra{
			Name:  "DI-BY",
			Value: "type",
		},
		_logIClaimsProviderExtra{
			Name:  "DI-M",
			Value: metaData,
		})
	di.AddScopedWithImplementedTypesWithMetadata(builder, implType, metaData, implementedTypes...)
}

// AddScopedIClaimsProviderByFunc adds a type by a custom func
func AddScopedIClaimsProviderByFunc(builder *di.Builder, implType reflect.Type, build func(ctn di.Container) (interface{}, error), implementedTypes ...reflect.Type) {
	implementedTypes = append(implementedTypes, ReflectTypeIClaimsProvider)
	_logAddIClaimsProvider("SCOPED", implType, _getImplementedIClaimsProviderNames(implementedTypes...),
		_logIClaimsProviderExtra{
			Name:  "DI-BY",
			Value: "func",
		})
	di.AddScopedWithImplementedTypesByFunc(builder, implType, build, implementedTypes...)
}

// AddScopedIClaimsProviderByFuncWithMetadata adds a type by a custom func
func AddScopedIClaimsProviderByFuncWithMetadata(builder *di.Builder, implType reflect.Type, build func(ctn di.Container) (interface{}, error), metaData map[string]interface{}, implementedTypes ...reflect.Type) {
	implementedTypes = append(implementedTypes, ReflectTypeIClaimsProvider)
	_logAddIClaimsProvider("SCOPED", implType, _getImplementedIClaimsProviderNames(implementedTypes...),
		_logIClaimsProviderExtra{
			Name:  "DI-BY",
			Value: "func",
		},
		_logIClaimsProviderExtra{
			Name:  "DI-M",
			Value: metaData,
		})

	di.AddScopedWithImplementedTypesByFuncWithMetadata(builder, implType, build, metaData, implementedTypes...)
}

// RemoveAllIClaimsProvider removes all IClaimsProvider from the DI
func RemoveAllIClaimsProvider(builder *di.Builder) {
	builder.RemoveAllByType(ReflectTypeIClaimsProvider)
}

// GetIClaimsProviderFromContainer alternative to SafeGetIClaimsProviderFromContainer but panics of object is not present
func GetIClaimsProviderFromContainer(ctn di.Container) IClaimsProvider {
	return ctn.GetByType(ReflectTypeIClaimsProvider).(IClaimsProvider)
}

// GetManyIClaimsProviderFromContainer alternative to SafeGetManyIClaimsProviderFromContainer but panics of object is not present
func GetManyIClaimsProviderFromContainer(ctn di.Container) []IClaimsProvider {
	objs := ctn.GetManyByType(ReflectTypeIClaimsProvider)
	var results []IClaimsProvider
	for _, obj := range objs {
		results = append(results, obj.(IClaimsProvider))
	}
	return results
}

// SafeGetIClaimsProviderFromContainer trys to get the object by type, will not panic, returns nil and error
func SafeGetIClaimsProviderFromContainer(ctn di.Container) (IClaimsProvider, error) {
	obj, err := ctn.SafeGetByType(ReflectTypeIClaimsProvider)
	if err != nil {
		return nil, err
	}
	return obj.(IClaimsProvider), nil
}

// GetIClaimsProviderDefinition returns that last definition registered that this container can provide
func GetIClaimsProviderDefinition(ctn di.Container) *di.Def {
	def := ctn.GetDefinitionByType(ReflectTypeIClaimsProvider)
	return def
}

// GetIClaimsProviderDefinitions returns all definitions that this container can provide
func GetIClaimsProviderDefinitions(ctn di.Container) []*di.Def {
	defs := ctn.GetDefinitionsByType(ReflectTypeIClaimsProvider)
	return defs
}

// SafeGetManyIClaimsProviderFromContainer trys to get the object by type, will not panic, returns nil and error
func SafeGetManyIClaimsProviderFromContainer(ctn di.Container) ([]IClaimsProvider, error) {
	objs, err := ctn.SafeGetManyByType(ReflectTypeIClaimsProvider)
	if err != nil {
		return nil, err
	}
	var results []IClaimsProvider
	for _, obj := range objs {
		results = append(results, obj.(IClaimsProvider))
	}
	return results, nil
}

type _logIClaimsProviderExtra struct {
	Name  string
	Value interface{}
}

func _logAddIClaimsProvider(scopeType string, implType reflect.Type, interfaces string, extra ..._logIClaimsProviderExtra) {
	infoEvent := log.Info().
		Str("DI", scopeType).
		Str("DI-I", interfaces).
		Str("DI-B", implType.Elem().String())

	for _, extra := range extra {
		infoEvent = infoEvent.Interface(extra.Name, extra.Value)
	}

	infoEvent.Send()

}
func _getImplementedIClaimsProviderNames(implementedTypes ...reflect.Type) string {
	builder := strings.Builder{}
	for idx, implementedType := range implementedTypes {
		builder.WriteString(implementedType.Name())
		if idx < len(implementedTypes)-1 {
			builder.WriteString(", ")
		}
	}
	return builder.String()
}
