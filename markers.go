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
//						to a database-supported large object type.
type LobMarker struct{}

// +marker="shelf:enumerated", UseValueSyntax=true, \
//			Description="Specifies that a persistent field should be persisted as a enumerated type."
type EnumeratedMarker struct {
	// +marker:argument="Value", \
	//					Options={STRING, ORDINAL},
	//					Description="The type used in mapping an enum type."
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
}

func (r RepositoryMarker) Validate() error {
	if strings.TrimSpace(r.Name) == "" {
		return errors.New("'Value' cannot be empty or nil")
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
