package main

import (
	"errors"
	"fmt"
	"log/slog"
	"sort"
	"strings"
	"sync"
)

type __dgi_Links struct {
	mu       sync.Mutex
	data     map[string]map[string]struct{}
	curModel string
}

type __dgi_Stack struct {
	data []string
}

func (s *__dgi_Stack) Pop() (string, error) {
	if s.IsEmpty() {
		return "", errors.New("stack is empty")
	}
	elem := s.data[len(s.data)-1]
	s.data = s.data[0 : len(s.data)-1]
	return elem, nil
}

func (s *__dgi_Stack) Push(elem string) {
	s.data = append(s.data, elem)
}

func (s *__dgi_Stack) IsEmpty() bool {
	return len(s.data) == 0
}

func (s *__dgi_Stack) Peek() (string, error) {
	if s.IsEmpty() {
		return "", errors.New("stack is empty")
	}
	return s.data[len(s.data)-1], nil
}

func (l *__dgi_Links) TopologicalSort() ([]string, error) {
	allModels := []string{}
	for key, _ := range l.data {
		if key != "" {
			allModels = append(allModels, key)
		}
	}

	sort.Strings(allModels)

	slog.Debug(fmt.Sprintf("performing topological sort for %d models", len(allModels)))
	finalStack := &__dgi_Stack{}
	curStack := &__dgi_Stack{}
	visited := map[string]struct{}{}
	for _, model := range allModels {
		if _, ok := visited[model]; ok {
			continue
		}

		curStack.Push(model)
		if err := l.dfs(curStack, visited); err != nil {
			return nil, fmt.Errorf("topological sort failed: %w", err)
		}

		for !curStack.IsEmpty() {
			elem, err := curStack.Pop()
			if err != nil {
				return nil, fmt.Errorf("popping from stack: %w", err)
			}

			finalStack.Push(elem)
		}
	}

	slog.Debug(fmt.Sprintf("topological sort completed: %v", finalStack.data))
	return finalStack.data, nil
}

func (l *__dgi_Links) dfs(curStack *__dgi_Stack, visited map[string]struct{}) error {
	if curStack.IsEmpty() {
		return nil
	}

	elem, err := curStack.Peek()
	if err != nil {
		return fmt.Errorf("error while creating links: %w", err)
	}

	visited[elem] = struct{}{}
	for key, _ := range l.data[elem] {
		if _, ok := visited[key]; ok {
			continue
		}
		curStack.Push(key)
		if err := l.dfs(curStack, visited); err != nil {
			return err
		}
	}
	return nil
}

func (l *__dgi_Links) StartGen(model string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	slog.Debug(fmt.Sprintf("starting generation tracking for %s", model))
	l.curModel = model
	if _, ok := l.data[l.curModel]; !ok {
		l.data[l.curModel] = map[string]struct{}{}
	}
}

func (l *__dgi_Links) EndGen(model string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	slog.Debug(fmt.Sprintf("ending generation tracking for %s", model))
	l.curModel = ""
}

func (l *__dgi_Links) AcceptSignal(model string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	slog.Debug(fmt.Sprintf("recording dependency from %s to %s", l.curModel, model))
	if _, ok := l.data[l.curModel]; !ok {
		l.data[l.curModel] = map[string]struct{}{}
	}

	v := l.data[l.curModel]
	v[model] = struct{}{}
}

func (l *__dgi_Links) PrettyPrint() {
	l.mu.Lock()
	defer l.mu.Unlock()

	var sb strings.Builder
	sb.WriteString("Links {\n")

	// Print current model
	sb.WriteString("  curModel: ")
	if l.curModel != "" {
		sb.WriteString(fmt.Sprintf("%q\n", l.curModel))
	} else {
		sb.WriteString("nil\n")
	}

	// Print data map
	sb.WriteString("  data: {\n")
	for src, dstMap := range l.data {
		srcStr := "<nil>"
		if src != "" {
			srcStr = fmt.Sprintf("%q", src)
		}
		sb.WriteString(fmt.Sprintf("    %s: {\n", srcStr))
		for dst := range dstMap {
			dstStr := "<nil>"
			if dst != "" {
				dstStr = fmt.Sprintf("%q", dst)
			}
			sb.WriteString(fmt.Sprintf("      %s\n", dstStr))
		}
		sb.WriteString("    }\n")
	}
	sb.WriteString("  }\n}\n")

	fmt.Println(sb.String())
}
