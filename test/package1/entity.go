package package1

// +import=shelf, Pkg=github.com/procyon-projects/shelf
import "time"

// +shelf:entity=User
// +shelf:table=users
type UserEntity struct {
	UserLoader
	// +shelf:id
	// +shelf:generated-value
	// +shelf:column:Unique=true, Length=500
	Id        int
	FirstName string
	LastName  string
	// +shelf:one-to-many
	Posts []UserEntity
}

// +shelf:entity=Post
// +shelf:table=posts
type PostEntity struct {
	// +shelf:id
	// +shelf:generated-value
	// +shelf:column=id
	Id int
	// +shelf:column=title
	Title string
	// +shelf:one-to-one
	// +shelf:join-column=post_detail_id
	PostDetails PostDetailsEntity
}

// +shelf:entity=PostDetails
// +shelf:table=post_details
type PostDetailsEntity struct {
	// +shelf:id
	// +shelf:generated-value
	// +shelf:column=id
	Id int
	// +shelf:column=created_on
	CreatedOn time.Time
	// +shelf:column=created_by
	CreatedBy string
}
