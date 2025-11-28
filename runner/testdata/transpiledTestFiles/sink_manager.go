package main

import (
	"fmt"
	"log/slog"
	"slices"
)

func __dgi_orchestrateSinks(topologicallySorted []string, allData map[string][]__dgi_Record, cfg *__dgi_Config) error {
	if cfg.ClearData {
		slog.Info("clearing existing data from sinks")
		if err := __dgi_clearAllData(topologicallySorted, allData, cfg); err != nil {
			return fmt.Errorf("error in clearing data: %w", err)
		}
	}

	slog.Info("loading data into sinks")
	return __dgi_loadAllData(topologicallySorted, allData, cfg)
}

func __dgi_clearAllData(topologicallySorted []string, allData map[string][]__dgi_Record, cfg *__dgi_Config) error {
	reversedTopologicallySorted := slices.Clone(topologicallySorted)
	slices.Reverse(reversedTopologicallySorted)
	slog.Debug(fmt.Sprintf("clearing data in reverse topological order: %v", reversedTopologicallySorted))
	for _, name := range reversedTopologicallySorted {
		if _, ok := allData[name]; !ok {
			continue
		}

		if err := __dgi_clearModelSinks(name, allData[name], cfg); err != nil {
			return fmt.Errorf("error clearing sinks for model %s: %w", name, err)
		}
	}
	slog.Info("data deletion completed successfully")
	return nil
}

func __dgi_loadAllData(topologicallySorted []string, allData map[string][]__dgi_Record, cfg *__dgi_Config) error {
	slog.Debug(fmt.Sprintf("loading data in topological order: %v", topologicallySorted))
	for _, name := range topologicallySorted {
		records, ok := allData[name]
		if !ok {
			continue
		}

		if err := __dgi_loadModelSinks(name, records, cfg); err != nil {
			return fmt.Errorf("%q, skipping further models", err)
		}
	}
	slog.Info("data loading completed successfully")
	return nil
}

// clearModelSinks routes model records to configured sinks per config.json
func __dgi_clearModelSinks(modelName string, records []__dgi_Record, cfg *__dgi_Config) error {
	sinks, err := cfg.SinkSpecsForModel(modelName)
	if err != nil {
		return fmt.Errorf("error while getting sink specs for model %s: %w", modelName, err)
	}

	slog.Debug(fmt.Sprintf("clearing %s from %d sinks", modelName, len(sinks)))
	for _, s := range sinks {
		switch s.SinkType {
		case __dgi_SinkTypeMySQL:
			err := __dgi_clearMysqlSink(s, modelName)
			if err != nil {
				return fmt.Errorf("error while clearing MySQL sink %s: %w", s.SinkName, err)
			}
		case __dgi_SinkTypePostgres:
			err := __dgi_clearPostgresSink(s, modelName)
			if err != nil {
				return fmt.Errorf("error while clearing Postgres sink %s: %w", s.SinkName, err)
			}
		default:
			return fmt.Errorf("unsupported sink_type %q for model %q", s.SinkType, modelName)
		}
	}
	return nil
}

// loadModelSinks routes model records to configured sinks per config.json
func __dgi_loadModelSinks(modelName string, records []__dgi_Record, cfg *__dgi_Config) error {
	sinks, err := cfg.SinkSpecsForModel(modelName)
	if err != nil {
		return fmt.Errorf("error while getting sink specs for model %s: %w", modelName, err)
	}

	slog.Debug(fmt.Sprintf("loading %s to %d sinks with %d records", modelName, len(sinks), len(records)))
	for _, s := range sinks {
		switch s.SinkType {
		case __dgi_SinkTypeMySQL:
			err := __dgi_loadMysqlSink(s, modelName, records)
			if err != nil {
				return fmt.Errorf("error in loading MySQL sink %s: %w", s.SinkName, err)
			}
		case __dgi_SinkTypePostgres:
			err := __dgi_loadPostgresSink(s, modelName, records)
			if err != nil {
				return fmt.Errorf("error in loading Postgres sink %s: %w", s.SinkName, err)
			}
		default:
			return fmt.Errorf("unsupported sink_type %q for model %q", s.SinkType, modelName)
		}
	}
	return nil
}

func __dgi_loadMysqlSink(sinkSpec *__dgi_SinkSpec, modelName string, records []__dgi_Record) error {
	var sc __dgi_MySQLConfig
	if err := sinkSpec.ConfigInto(&sc); err != nil {
		return fmt.Errorf("mysql sink %q config: %w", sinkSpec.SinkName, err)
	}

	if sc.BatchSize <= 0 {
		sc.BatchSize = len(records)
	}

	switch modelName {
	case "minimal":
		typed := make([]*__datagen_minimal, 0, len(records))
		for _, r := range records {
			typed = append(typed, r.(*__datagen_minimal))
		}

		return Sink_mysql___datagen_minimal_data(modelName, typed, &sc)
	case "multiple_types":
		typed := make([]*__datagen_multiple_types, 0, len(records))
		for _, r := range records {
			typed = append(typed, r.(*__datagen_multiple_types))
		}

		return Sink_mysql___datagen_multiple_types_data(modelName, typed, &sc)
	case "nested":
		typed := make([]*__datagen_nested, 0, len(records))
		for _, r := range records {
			typed = append(typed, r.(*__datagen_nested))
		}

		return Sink_mysql___datagen_nested_data(modelName, typed, &sc)
	case "simple":
		typed := make([]*__datagen_simple, 0, len(records))
		for _, r := range records {
			typed = append(typed, r.(*__datagen_simple))
		}

		return Sink_mysql___datagen_simple_data(modelName, typed, &sc)
	case "with_builtin_functions":
		typed := make([]*__datagen_with_builtin_functions, 0, len(records))
		for _, r := range records {
			typed = append(typed, r.(*__datagen_with_builtin_functions))
		}

		return Sink_mysql___datagen_with_builtin_functions_data(modelName, typed, &sc)
	case "with_conditionals":
		typed := make([]*__datagen_with_conditionals, 0, len(records))
		for _, r := range records {
			typed = append(typed, r.(*__datagen_with_conditionals))
		}

		return Sink_mysql___datagen_with_conditionals_data(modelName, typed, &sc)
	case "with_maps":
		typed := make([]*__datagen_with_maps, 0, len(records))
		for _, r := range records {
			typed = append(typed, r.(*__datagen_with_maps))
		}

		return Sink_mysql___datagen_with_maps_data(modelName, typed, &sc)
	case "with_metadata":
		typed := make([]*__datagen_with_metadata, 0, len(records))
		for _, r := range records {
			typed = append(typed, r.(*__datagen_with_metadata))
		}

		return Sink_mysql___datagen_with_metadata_data(modelName, typed, &sc)
	case "with_misc":
		typed := make([]*__datagen_with_misc, 0, len(records))
		for _, r := range records {
			typed = append(typed, r.(*__datagen_with_misc))
		}

		return Sink_mysql___datagen_with_misc_data(modelName, typed, &sc)
	case "with_slices":
		typed := make([]*__datagen_with_slices, 0, len(records))
		for _, r := range records {
			typed = append(typed, r.(*__datagen_with_slices))
		}

		return Sink_mysql___datagen_with_slices_data(modelName, typed, &sc)
	default:
		return fmt.Errorf("mysql sink not implemented for model %q", modelName)
	}
}

func __dgi_clearMysqlSink(sinkSpec *__dgi_SinkSpec, modelName string) error {
	var sc __dgi_MySQLConfig
	if err := sinkSpec.ConfigInto(&sc); err != nil {
		return fmt.Errorf("mysql sink %q config: %w", sinkSpec.SinkName, err)
	}

	switch modelName {
	case "minimal":
		return Clear_mysql___datagen_minimal_data(modelName, &sc)
	case "multiple_types":
		return Clear_mysql___datagen_multiple_types_data(modelName, &sc)
	case "nested":
		return Clear_mysql___datagen_nested_data(modelName, &sc)
	case "simple":
		return Clear_mysql___datagen_simple_data(modelName, &sc)
	case "with_builtin_functions":
		return Clear_mysql___datagen_with_builtin_functions_data(modelName, &sc)
	case "with_conditionals":
		return Clear_mysql___datagen_with_conditionals_data(modelName, &sc)
	case "with_maps":
		return Clear_mysql___datagen_with_maps_data(modelName, &sc)
	case "with_metadata":
		return Clear_mysql___datagen_with_metadata_data(modelName, &sc)
	case "with_misc":
		return Clear_mysql___datagen_with_misc_data(modelName, &sc)
	case "with_slices":
		return Clear_mysql___datagen_with_slices_data(modelName, &sc)
	default:
		return fmt.Errorf("mysql sink not implemented for model %q", modelName)
	}
}

func __dgi_loadPostgresSink(sinkSpec *__dgi_SinkSpec, modelName string, records []__dgi_Record) error {
	var sc __dgi_PostgresConfig
	if err := sinkSpec.ConfigInto(&sc); err != nil {
		return fmt.Errorf("postgres sink %q config: %w", sinkSpec.SinkName, err)
	}

	if sc.BatchSize <= 0 {
		sc.BatchSize = len(records)
	}

	switch modelName {
	case "minimal":
		typed := make([]*__datagen_minimal, 0, len(records))
		for _, r := range records {
			typed = append(typed, r.(*__datagen_minimal))
		}

		return Sink_postgres___datagen_minimal_data(modelName, typed, &sc)
	case "multiple_types":
		typed := make([]*__datagen_multiple_types, 0, len(records))
		for _, r := range records {
			typed = append(typed, r.(*__datagen_multiple_types))
		}

		return Sink_postgres___datagen_multiple_types_data(modelName, typed, &sc)
	case "nested":
		typed := make([]*__datagen_nested, 0, len(records))
		for _, r := range records {
			typed = append(typed, r.(*__datagen_nested))
		}

		return Sink_postgres___datagen_nested_data(modelName, typed, &sc)
	case "simple":
		typed := make([]*__datagen_simple, 0, len(records))
		for _, r := range records {
			typed = append(typed, r.(*__datagen_simple))
		}

		return Sink_postgres___datagen_simple_data(modelName, typed, &sc)
	case "with_builtin_functions":
		typed := make([]*__datagen_with_builtin_functions, 0, len(records))
		for _, r := range records {
			typed = append(typed, r.(*__datagen_with_builtin_functions))
		}

		return Sink_postgres___datagen_with_builtin_functions_data(modelName, typed, &sc)
	case "with_conditionals":
		typed := make([]*__datagen_with_conditionals, 0, len(records))
		for _, r := range records {
			typed = append(typed, r.(*__datagen_with_conditionals))
		}

		return Sink_postgres___datagen_with_conditionals_data(modelName, typed, &sc)
	case "with_maps":
		typed := make([]*__datagen_with_maps, 0, len(records))
		for _, r := range records {
			typed = append(typed, r.(*__datagen_with_maps))
		}

		return Sink_postgres___datagen_with_maps_data(modelName, typed, &sc)
	case "with_metadata":
		typed := make([]*__datagen_with_metadata, 0, len(records))
		for _, r := range records {
			typed = append(typed, r.(*__datagen_with_metadata))
		}

		return Sink_postgres___datagen_with_metadata_data(modelName, typed, &sc)
	case "with_misc":
		typed := make([]*__datagen_with_misc, 0, len(records))
		for _, r := range records {
			typed = append(typed, r.(*__datagen_with_misc))
		}

		return Sink_postgres___datagen_with_misc_data(modelName, typed, &sc)
	case "with_slices":
		typed := make([]*__datagen_with_slices, 0, len(records))
		for _, r := range records {
			typed = append(typed, r.(*__datagen_with_slices))
		}

		return Sink_postgres___datagen_with_slices_data(modelName, typed, &sc)
	default:
		return fmt.Errorf("postgres sink not implemented for model %q", modelName)
	}
}

func __dgi_clearPostgresSink(sinkSpec *__dgi_SinkSpec, modelName string) error {
	var sc __dgi_PostgresConfig
	if err := sinkSpec.ConfigInto(&sc); err != nil {
		return fmt.Errorf("postgres sink %q config: %w", sinkSpec.SinkName, err)
	}

	switch modelName {
	case "minimal":
		return Clear_postgres___datagen_minimal_data(modelName, &sc)
	case "multiple_types":
		return Clear_postgres___datagen_multiple_types_data(modelName, &sc)
	case "nested":
		return Clear_postgres___datagen_nested_data(modelName, &sc)
	case "simple":
		return Clear_postgres___datagen_simple_data(modelName, &sc)
	case "with_builtin_functions":
		return Clear_postgres___datagen_with_builtin_functions_data(modelName, &sc)
	case "with_conditionals":
		return Clear_postgres___datagen_with_conditionals_data(modelName, &sc)
	case "with_maps":
		return Clear_postgres___datagen_with_maps_data(modelName, &sc)
	case "with_metadata":
		return Clear_postgres___datagen_with_metadata_data(modelName, &sc)
	case "with_misc":
		return Clear_postgres___datagen_with_misc_data(modelName, &sc)
	case "with_slices":
		return Clear_postgres___datagen_with_slices_data(modelName, &sc)
	default:
		return fmt.Errorf("postgres sink not implemented for model %q", modelName)
	}
}

func __dgi_getRecordCount(cfg *__dgi_Config, modelName string, metadata __dgi_Metadata) int {
	for _, m := range cfg.Models {
		if m.ModelName == modelName {
			if m.Count != nil {
				return *m.Count
			}
			return metadata.Count
		}
	}

	return metadata.Count
}
