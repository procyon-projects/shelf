package main

import "github.com/procyon-projects/marker"

type TypeMetadata struct {
	IsPointer   bool
	IsAnonymous bool
}

func TypeMetadataFromType(typ marker.Type) *TypeMetadata {
	typeMetadata := &TypeMetadata{}
	return typeMetadata
}
