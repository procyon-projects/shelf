package package1

// +import=shelf, Pkg=github.com/procyon-projects/shelf
import (
	"context"
)

// +shelf:repository="user-repository", Entity=User
type UserRepository interface {
	LoadPosts() UserRepository

	Count(ctx context.Context) int
	ExistsById(ctx context.Context, id int) bool

	Delete(ctx context.Context, user *User)
	DeleteById(ctx context.Context, id int)
	DeleteAll(ctx context.Context, user []*User)
	DeleteAllById(ctx context.Context, ids []int)

	Save(ctx context.Context, user *User)
	SaveAll(ctx context.Context, user []*User)

	FindById(ctx context.Context, id int) *User
	FindAll(ctx context.Context) []*User
	FindAllById(ctx context.Context, ids []int) []*User

	FindByFirstNameAndLastName(ctx context.Context, firstName, lastName string) *User
	// +shelf:query="FROM User WHERE FirstName = %1 AND LastName = %2"
	CustomQuery(ctx context.Context, firstName string) *User
}
