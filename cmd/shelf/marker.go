package main

type MarkerFlag int

const (
	Column MarkerFlag = 1 << iota
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

const (
	Associations = OneToOne | OneToMany | ManyToOne | ManyToMany
	DateAndTime  = Temporal | CreatedDate | LastModifiedDate
)
