package shelf

import (
	"errors"
	"fmt"
	"strings"
)

const (
	MarkerEntity = "shelf:entity"
	MarkerTable  = "shelf:table"

	MarkerId             = "shelf:id"
	MarkerGeneratedValue = "shelf:generated-value"

	MarkerColumn    = "shelf:column"
	MarkerTransient = "shelf:transient"
	MarkerLob       = "shelf:lob"

	MarkerEnumerated = "shelf:enumerated"

	MarkerRepository = "shelf:repository"
	MarkerQuery      = "shelf:query"

	MarkerEmbeddable        = "shelf:embeddable"
	MarkerEmbedded          = "shelf:embedded"
	MarkerAttributeOverride = "shelf:attribute-override"

	MarkerMapsId     = "shelf:maps-id"
	MarkerOneToOne   = "shelf:one-to-one"
	MarkerOneToMany  = "shelf:one-to-many"
	MarkerManyToOne  = "shelf:many-to-one"
	MarkerManyToMany = "shelf:many-to-many"

	MarkerTemporal = "shelf:temporal"
)

// +marker="shelf:entity", UseValueSyntax=true, Description="Specifies that the class is an entity."
type EntityMarker struct {
	// +marker:argument="Value", Optional=true, Description="The entity name."
	Name string `marker:"Value,useValueSyntax,optional"`
}

// +marker="shelf:table", UseValueSyntax=true, Description="Specifies the primary table for the marked entity."
type TableMarker struct {
	// +marker:argument="Value", Optional=true, Description="The name of the table."
	Name string `marker:"Value,useValueSyntax,optional"`
}

// +marker="shelf:id", Description="Specifies the primary key of an entity."
type IdMarker struct{}

// +marker="shelf:generated-value", Description="Provides for the specification of generation strategies \
//			for the values of primary keys."
type GeneratedValueMarker struct{}

// +marker="shelf:column", UseValueSyntax=true, Description="Specifies the mapped column for a persistent field."
type ColumnMarker struct {
	// +marker:argument="Name", Optional=true, Description="The name of the column."
	Name string `marker:"Value,useValueSyntax,optional"`
	// +marker:argument="Unique", Optional=true, Description="Whether the column is a unique key."
	Unique bool `marker:"Unique,optional"`
	// +marker:argument="Length", Optional=true, Description="The column length."
	Length int `marker:"Length,optional"`
}

// +marker="shelf:id", Description="Specifies the primary key of an entity."
type TransientMarker struct{}

// +marker="shelf-lob", Description="Specifies that a persistent field should be persisted as a large object \
//	to a database-supported large object type.
type LobMarker struct{}

// +marker="shelf:enumerated", UseValueSyntax=true, \
//			Description="Specifies that a persistent field should be persisted as a enumerated type."
type EnumeratedMarker struct {
	// +marker:argument="Value", \
	//	Options={STRING, ORDINAL}, \
	//	Description="The type used in mapping an enum type."
	Value string `marker:"Value,useValueSyntax"`
}

func (e EnumeratedMarker) Validate() error {
	matched := false
	enumeratedOptions := []string{"STRING", "ORDINAL"}

	for _, option := range enumeratedOptions {
		if strings.TrimSpace(e.Value) == option {
			matched = true
		}
	}

	if !matched {
		return fmt.Errorf("invalid Enumerated option. Here is the list of valid options %s", strings.Join(enumeratedOptions, ", "))
	}

	return nil
}

// +marker="shelf:repository", UseValueSyntax=true, Description="Specifies a repository."
type RepositoryMarker struct {
	// +marker:argument="Value", Description="The repository name."
	Name string `marker:"Value,useValueSyntax"`
	// +marker:argument="EntityStruct", Description="The entity name"
	Entity string `marker:"Entity"`
}

func (r RepositoryMarker) Validate() error {
	if strings.TrimSpace(r.Name) == "" {
		return errors.New("'Value' cannot be empty or nil")
	}

	if strings.TrimSpace(r.Entity) == "" {
		return errors.New("'Entity' cannot be empty or nil")
	}

	return nil
}

// +marker="shelf:query", UseValueSyntax=true, Description="Specifies a query."
type QueryMarker struct {
	// +marker:argument="Value", Description="The query string."
	Value string `marker:"Value,useValueSyntax"`
	// +marker:argument="Unique", Optional=true, Description="Whether the query is a native queery."
	NativeQuery bool `marker:"NativeQuery,optional"`
}

func (q QueryMarker) Validate() error {
	if strings.TrimSpace(q.Value) == "" {
		return errors.New("'Value' cannot be empty or nil")
	}

	return nil
}

// +marker="shelf:embeddable", Description="Specifies that a struct will be embedded by other entities."
type EmbeddableMarker struct{}

// +marker="shelf:embedded", Description="Specifies that an entity embed a struct"
type EmbeddedMarker struct{}

// +marker="shelf:attribute-override", Description="Specifies that the column property of embedded type will be overridden."
type AttributeOverrideMarker struct {
	// +marker:argument="Value", Description=" The association override mappings that are to be applied to the relationship field."
	Name string `marker:"Value,useValueSyntax"`
	// +marker:argument="ColumnName", Optional=true, Description="The name of the column."
	ColumnName string `marker:"ColumnName"`
	// +marker:argument="ColumnUnique", Optional=true, Description="Whether the column is a unique key."
	ColumnUnique bool `marker:"ColumnUnique,optional"`
	// +marker:argument="ColumnLength", Optional=true, Description="The column length."
	ColumnLength int `marker:"ColumnLength,optional"`
}

func (a AttributeOverrideMarker) Validate() error {
	if strings.TrimSpace(a.Name) == "" {
		return errors.New("'Value' cannot be empty or nil")
	}

	if strings.TrimSpace(a.ColumnName) == "" {
		return errors.New("'ColumnName' cannot be empty or nil")
	}

	return nil
}

// +marker="shelf:maps-id", Description=""
type MapsIdMarker struct{}

// +marker="shelf:one-to-one", Description="Specifies a single-valued association to another entity."
type OneToOneMarker struct {
	// +marker:argument="Cascade", Description="The operations that must be cascaded to the target of the association."
	Cascade []string `marker:"Cascade,optional"`
	// +marker:argument="FetchType", Description="Whether the association should be lazily loaded or must be eagerly fetched."
	FetchType string `marker:"FetchType,optional"`
	// +marker:argument="MappedBy", Description="The field that owns the relationship."
	MappedBy string `marker:"MappedBy,optional"`
}

func (o OneToOneMarker) Validate() error {

	if o.FetchType != "" {
		matched := false
		fetchTypeOptions := []string{"LAZY", "EAGER"}

		for _, option := range fetchTypeOptions {
			if strings.TrimSpace(o.FetchType) == option {
				matched = true
			}
		}

		if !matched {
			return fmt.Errorf("invalid FetchType option. Here is the list of valid options %s", strings.Join(fetchTypeOptions, ", "))
		}
	}

	if o.Cascade != nil && len(o.Cascade) != 0 {
		matched := false
		cascadeOptions := []string{"ALL", "PERSIST", "SAVE_UPDATE", "REMOVE"}

		for _, option := range cascadeOptions {
			if strings.TrimSpace(o.FetchType) == option {
				matched = true
			}
		}

		if !matched {
			return fmt.Errorf("invalid FetchType option. Here is the list of valid options %s", strings.Join(cascadeOptions, ", "))
		}
	}

	return nil
}

// +marker="shelf:one-to-many", Description="Specifies a many-valued association."
type OneToManyMarker struct {
	// +marker:argument="Cascade", Description="The operations that must be cascaded to the target of the association."
	Cascade []string `marker:"Cascade,optional"`
	// +marker:argument="FetchType", Description="Whether the association should be lazily loaded or must be eagerly fetched."
	FetchType string `marker:"FetchType,optional"`
	// +marker:argument="MappedBy", Description="The field that owns the relationship."
	MappedBy string `marker:"MappedBy,optional"`
}

func (o OneToManyMarker) Validate() error {

	if o.FetchType != "" {
		matched := false
		fetchTypeOptions := []string{"LAZY", "EAGER"}

		for _, option := range fetchTypeOptions {
			if strings.TrimSpace(o.FetchType) == option {
				matched = true
			}
		}

		if !matched {
			return fmt.Errorf("invalid FetchType option. Here is the list of valid options %s", strings.Join(fetchTypeOptions, ", "))
		}
	}

	if o.Cascade != nil && len(o.Cascade) != 0 {
		matched := false
		cascadeOptions := []string{"ALL", "PERSIST", "SAVE_UPDATE", "REMOVE"}

		for _, option := range cascadeOptions {
			if strings.TrimSpace(o.FetchType) == option {
				matched = true
			}
		}

		if !matched {
			return fmt.Errorf("invalid FetchType option. Here is the list of valid options %s", strings.Join(cascadeOptions, ", "))
		}
	}

	return nil
}

// +marker="shelf:many-to-one", Description="Specifies a single-valued association to another entity class that has many-to-one multiplicity"
type ManyToOneMarker struct {
	// +marker:argument="Cascade", Description="The operations that must be cascaded to the target of the association."
	Cascade []string `marker:"Cascade,optional"`
	// +marker:argument="FetchType", Description="Whether the association should be lazily loaded or must be eagerly fetched."
	FetchType string `marker:"FetchType,optional"`
	// +marker:argument="MappedBy", Description="The field that owns the relationship."
	MappedBy string `marker:"MappedBy,optional"`
}

func (o ManyToOneMarker) Validate() error {

	if o.FetchType != "" {
		matched := false
		fetchTypeOptions := []string{"LAZY", "EAGER"}

		for _, option := range fetchTypeOptions {
			if strings.TrimSpace(o.FetchType) == option {
				matched = true
			}
		}

		if !matched {
			return fmt.Errorf("invalid FetchType option. Here is the list of valid options %s", strings.Join(fetchTypeOptions, ", "))
		}
	}

	if o.Cascade != nil && len(o.Cascade) != 0 {
		matched := false
		cascadeOptions := []string{"ALL", "PERSIST", "SAVE_UPDATE", "REMOVE"}

		for _, option := range cascadeOptions {
			if strings.TrimSpace(o.FetchType) == option {
				matched = true
			}
		}

		if !matched {
			return fmt.Errorf("invalid FetchType option. Here is the list of valid options %s", strings.Join(cascadeOptions, ", "))
		}
	}

	return nil
}

// +marker="shelf:many-to-many", Description="Specifies a many-valued association with many-to-many multiplicity."
type ManyToManyMarker struct {
	// +marker:argument="Cascade", Description="The operations that must be cascaded to the target of the association."
	Cascade []string `marker:"Cascade,optional"`
	// +marker:argument="FetchType", Description="Whether the association should be lazily loaded or must be eagerly fetched."
	FetchType string `marker:"FetchType,optional"`
	// +marker:argument="MappedBy", Description="The field that owns the relationship."
	MappedBy string `marker:"MappedBy,optional"`
}

func (o ManyToManyMarker) Validate() error {

	if o.FetchType != "" {
		matched := false
		fetchTypeOptions := []string{"LAZY", "EAGER"}

		for _, option := range fetchTypeOptions {
			if strings.TrimSpace(o.FetchType) == option {
				matched = true
			}
		}

		if !matched {
			return fmt.Errorf("invalid FetchType option. Here is the list of valid options %s", strings.Join(fetchTypeOptions, ", "))
		}
	}

	if o.Cascade != nil && len(o.Cascade) != 0 {
		matched := false
		cascadeOptions := []string{"ALL", "PERSIST", "SAVE_UPDATE", "REMOVE"}

		for _, option := range cascadeOptions {
			if strings.TrimSpace(o.FetchType) == option {
				matched = true
			}
		}

		if !matched {
			return fmt.Errorf("invalid FetchType option. Here is the list of valid options %s", strings.Join(cascadeOptions, ", "))
		}
	}

	return nil
}

// +marker="shelf:temporal", UseValueSyntax=true
type TemporalMarker struct {
	// +marker:argument="Value", \
	//	Options={DATE, TIME, TIMESTAMP}, \
	//	Description="The type used in mapping an enum type."
	Value string `marker:"Value,useValueSyntax"`
}

func (t TemporalMarker) Validate() error {
	matched := false
	temporalOptions := []string{"DATE", "TIME", "TIMESTAMP"}

	for _, option := range temporalOptions {
		if strings.TrimSpace(t.Value) == option {
			matched = true
		}
	}

	if !matched {
		return fmt.Errorf("invalid Temporal option. Here is the list of valid options %s", strings.Join(temporalOptions, ", "))
	}

	return nil
}
