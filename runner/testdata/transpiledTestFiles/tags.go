package main

import (
	"fmt"
	"strings"
)

func parseTags(tagsStr string) (map[string]string, error) {
	tags := make(map[string]string)

	if tagsStr == "" {
		return tags, nil
	}

	pairs := strings.Split(tagsStr, ",")

	for _, pair := range pairs {
		pair = strings.TrimSpace(pair)

		parts := strings.SplitN(pair, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid tag format: %s (expected key=value)", pair)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		if key == "" {
			return nil, fmt.Errorf("empty key in tag: %s", pair)
		}

		tags[key] = value
	}

	return tags, nil
}

func getModelsMetadata(datagen *DataGenGenerators) map[string]Metadata {
	out := make(map[string]Metadata)
	if datagen.with_conditionals != nil {
		out["with_conditionals"] = datagen.with_conditionals().Metadata()
	}
	if datagen.with_slices != nil {
		out["with_slices"] = datagen.with_slices().Metadata()
	}
	if datagen.multiple_types != nil {
		out["multiple_types"] = datagen.multiple_types().Metadata()
	}
	if datagen.nested != nil {
		out["nested"] = datagen.nested().Metadata()
	}
	if datagen.simple != nil {
		out["simple"] = datagen.simple().Metadata()
	}
	if datagen.with_maps != nil {
		out["with_maps"] = datagen.with_maps().Metadata()
	}
	if datagen.with_metadata != nil {
		out["with_metadata"] = datagen.with_metadata().Metadata()
	}
	if datagen.with_misc != nil {
		out["with_misc"] = datagen.with_misc().Metadata()
	}
	if datagen.minimal != nil {
		out["minimal"] = datagen.minimal().Metadata()
	}
	if datagen.with_builtin_functions != nil {
		out["with_builtin_functions"] = datagen.with_builtin_functions().Metadata()
	}
	return out
}

func getMatchingModels(modelsMetadata map[string]Metadata, need map[string]string) []string {
	matchedModels := make([]string, 0)
	if len(need) == 0 {
		return matchedModels
	}

	for name, md := range modelsMetadata {
		match := true
		for k, v := range need {
			if md.Tags[k] != v {
				match = false
				break
			}
		}
		if match {
			matchedModels = append(matchedModels, name)
		}
	}
	return matchedModels
}
