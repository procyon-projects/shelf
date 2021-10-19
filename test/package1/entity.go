package package1

// +import=shelf, Pkg=github.com/procyon-projects/shelf
import (
	"github.com/procyon-projects/shelf/test/package2"
	"time"
)

// +shelf:embeddable
type IdEntity struct {
	// +shelf:id
	// +shelf:generated-value
	Id int
}

// +shelf:embeddable
type BaseEntity struct {
	IdEntity
	// +shelf:column=created_on
	// +shelf:created-date
	CreatedOn time.Time
}

// +shelf:entity
// +shelf:table
type User struct {
	BaseEntity
	// +shelf:column=Unique=true, Length=500
	Email string

	FirstName string
	LastName  string
	// +shelf:enumerated=STRING
	Status package2.UserStatus

	// +shelf:embedded
	// +shelf:attribute-override=City, ColumnName="address_city"
	// +shelf:attribute-override="Country.Name", ColumnName="address_country"
	// +shelf:attribute-override=PostCode, ColumnName="address_post_code"
	Address *Address

	// +shelf:one-to-one:FetchType=LAZY,MappedBy=User
	CreditCard *CreditCard

	// +shelf:one-to-many
	Posts []Post
}

// +shelf:entity
type CreditCard struct {
	// +shelf:id
	// +shelf:generated-value
	Id     int
	Number string

	// +shelf:one-to-one:FetchType=LAZY
	User *User
}

// +shelf:embeddable
type Address struct {
	// +shelf:column=postCity
	City *[]string
	// +shelf:embedded
	// +shelf:attribute-override="Name", ColumnName="address_post_code"
	Country *Country
	// +shelf:column=postCode
	PostCode string
}

// +shelf:embeddable
type Country struct {
	// +shelf:column=Name
	Name string
}

// +shelf:entity=Post
// +shelf:table=posts
type Post struct {
	// +shelf:id
	// +shelf:generated-value
	// +shelf:column=id
	Id int
	// +shelf:column=title
	Title string
	// +shelf:one-to-one
	// +shelf:join-column=post_detail_id
	PostDetails *PostDetails
}

// +shelf:entity=PostDetails
// +shelf:table
type PostDetails struct {
	// +shelf:id
	// +shelf:generated-value
	// +shelf:column=id
	Id int
	// +shelf:column=created_on
	// +shelf:created-date
	CreatedOn time.Time
	// +shelf:column=created_by
	CreatedBy string
}
