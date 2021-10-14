package main

type MarkerFlag int

const (
	Column MarkerFlag = iota
	Transient
	Id
	GeneratedValue
	Temporal
	CreatedDate
	LastModifiedDate
	Enumerated
	Lob
	Embedded
	AttributeOverride
	OneToOne
	OneToMany
	ManyToOne
	ManyToMany
)
