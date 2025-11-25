package main

import (
	"sync"
)

type __dgi_DataGenGenerators struct {
	minimal                func() *__datagen_minimalGenerator
	multiple_types         func() *__datagen_multiple_typesGenerator
	nested                 func() *__datagen_nestedGenerator
	simple                 func() *__datagen_simpleGenerator
	with_builtin_functions func() *__datagen_with_builtin_functionsGenerator
	with_conditionals      func() *__datagen_with_conditionalsGenerator
	with_maps              func() *__datagen_with_mapsGenerator
	with_metadata          func() *__datagen_with_metadataGenerator
	with_misc              func() *__datagen_with_miscGenerator
	with_slices            func() *__datagen_with_slicesGenerator

	__links *__dgi_Links
}

type __dgi_Metadata struct {
	Count int
	Tags  map[string]string
}

func minimalFunc(model *__datagen_minimalGenerator, tail string) func() *__datagen_minimalGenerator {
	return func() *__datagen_minimalGenerator {
		model.datagen.__links.AcceptSignal(tail)
		return model
	}
}
func multiple_typesFunc(model *__datagen_multiple_typesGenerator, tail string) func() *__datagen_multiple_typesGenerator {
	return func() *__datagen_multiple_typesGenerator {
		model.datagen.__links.AcceptSignal(tail)
		return model
	}
}
func nestedFunc(model *__datagen_nestedGenerator, tail string) func() *__datagen_nestedGenerator {
	return func() *__datagen_nestedGenerator {
		model.datagen.__links.AcceptSignal(tail)
		return model
	}
}
func simpleFunc(model *__datagen_simpleGenerator, tail string) func() *__datagen_simpleGenerator {
	return func() *__datagen_simpleGenerator {
		model.datagen.__links.AcceptSignal(tail)
		return model
	}
}
func with_builtin_functionsFunc(model *__datagen_with_builtin_functionsGenerator, tail string) func() *__datagen_with_builtin_functionsGenerator {
	return func() *__datagen_with_builtin_functionsGenerator {
		model.datagen.__links.AcceptSignal(tail)
		return model
	}
}
func with_conditionalsFunc(model *__datagen_with_conditionalsGenerator, tail string) func() *__datagen_with_conditionalsGenerator {
	return func() *__datagen_with_conditionalsGenerator {
		model.datagen.__links.AcceptSignal(tail)
		return model
	}
}
func with_mapsFunc(model *__datagen_with_mapsGenerator, tail string) func() *__datagen_with_mapsGenerator {
	return func() *__datagen_with_mapsGenerator {
		model.datagen.__links.AcceptSignal(tail)
		return model
	}
}
func with_metadataFunc(model *__datagen_with_metadataGenerator, tail string) func() *__datagen_with_metadataGenerator {
	return func() *__datagen_with_metadataGenerator {
		model.datagen.__links.AcceptSignal(tail)
		return model
	}
}
func with_miscFunc(model *__datagen_with_miscGenerator, tail string) func() *__datagen_with_miscGenerator {
	return func() *__datagen_with_miscGenerator {
		model.datagen.__links.AcceptSignal(tail)
		return model
	}
}
func with_slicesFunc(model *__datagen_with_slicesGenerator, tail string) func() *__datagen_with_slicesGenerator {
	return func() *__datagen_with_slicesGenerator {
		model.datagen.__links.AcceptSignal(tail)
		return model
	}
}

func __dgi_initGeneratorsAndModels() (*__dgi_DataGenGenerators, map[string]__dgi_RecordGenerator) {
	minimalGenerator := __init___datagen_minimalGenerator()
	multiple_typesGenerator := __init___datagen_multiple_typesGenerator()
	nestedGenerator := __init___datagen_nestedGenerator()
	simpleGenerator := __init___datagen_simpleGenerator()
	with_builtin_functionsGenerator := __init___datagen_with_builtin_functionsGenerator()
	with_conditionalsGenerator := __init___datagen_with_conditionalsGenerator()
	with_mapsGenerator := __init___datagen_with_mapsGenerator()
	with_metadataGenerator := __init___datagen_with_metadataGenerator()
	with_miscGenerator := __init___datagen_with_miscGenerator()
	with_slicesGenerator := __init___datagen_with_slicesGenerator()

	// Construct directory instances bottom-up so children are available

	datagen := &__dgi_DataGenGenerators{
		minimal:                minimalFunc(minimalGenerator, "minimal"),
		multiple_types:         multiple_typesFunc(multiple_typesGenerator, "multiple_types"),
		nested:                 nestedFunc(nestedGenerator, "nested"),
		simple:                 simpleFunc(simpleGenerator, "simple"),
		with_builtin_functions: with_builtin_functionsFunc(with_builtin_functionsGenerator, "with_builtin_functions"),
		with_conditionals:      with_conditionalsFunc(with_conditionalsGenerator, "with_conditionals"),
		with_maps:              with_mapsFunc(with_mapsGenerator, "with_maps"),
		with_metadata:          with_metadataFunc(with_metadataGenerator, "with_metadata"),
		with_misc:              with_miscFunc(with_miscGenerator, "with_misc"),
		with_slices:            with_slicesFunc(with_slicesGenerator, "with_slices"),

		__links: &__dgi_Links{
			mu:   sync.Mutex{},
			data: map[string]map[string]struct{}{},
		},
	}
	minimalGenerator.datagen = datagen
	multiple_typesGenerator.datagen = datagen
	nestedGenerator.datagen = datagen
	simpleGenerator.datagen = datagen
	with_builtin_functionsGenerator.datagen = datagen
	with_conditionalsGenerator.datagen = datagen
	with_mapsGenerator.datagen = datagen
	with_metadataGenerator.datagen = datagen
	with_miscGenerator.datagen = datagen
	with_slicesGenerator.datagen = datagen

	// model registry
	models := map[string]__dgi_RecordGenerator{
		"minimal":                minimalGenerator.Gen,
		"multiple_types":         multiple_typesGenerator.Gen,
		"nested":                 nestedGenerator.Gen,
		"simple":                 simpleGenerator.Gen,
		"with_builtin_functions": with_builtin_functionsGenerator.Gen,
		"with_conditionals":      with_conditionalsGenerator.Gen,
		"with_maps":              with_mapsGenerator.Gen,
		"with_metadata":          with_metadataGenerator.Gen,
		"with_misc":              with_miscGenerator.Gen,
		"with_slices":            with_slicesGenerator.Gen,
	}

	return datagen, models
}
