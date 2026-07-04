package models

import "fmt"

// Catalog stores all available tests by ID.
type Catalog struct {
	tests map[string]TestDefinition
	order []string
}

func NewCatalog(definitions ...TestDefinition) (*Catalog, error) {
	catalog := &Catalog{
		tests: make(map[string]TestDefinition, len(definitions)),
		order: make([]string, 0, len(definitions)),
	}

	for _, definition := range definitions {
		if err := definition.Validate(); err != nil {
			return nil, err
		}
		if _, exists := catalog.tests[definition.ID]; exists {
			return nil, fmt.Errorf("duplicate test id %q", definition.ID)
		}

		catalog.tests[definition.ID] = definition
		catalog.order = append(catalog.order, definition.ID)
	}

	return catalog, nil
}

func DefaultCatalog() (*Catalog, error) {
	return NewCatalog(GardnerTest(), CliftonTest())
}

func (c *Catalog) All() []TestDefinition {
	definitions := make([]TestDefinition, 0, len(c.order))
	for _, id := range c.order {
		definitions = append(definitions, c.tests[id])
	}
	return definitions
}

func (c *Catalog) Get(id string) (TestDefinition, bool) {
	definition, ok := c.tests[id]
	return definition, ok
}
