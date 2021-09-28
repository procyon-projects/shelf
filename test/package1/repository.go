package package1

// +import=shelf, Pkg=github.com/procyon-projects/shelf
import "context"

// +shelf:loader=user-loader
type UserLoader interface {
	LoadPosts(ctx context.Context)
}

// +shelf:repository="user-repository"
type UserRepository interface {
	Count(ctx context.Context) int
	ExistsById(ctx context.Context, id int) bool

	Delete(ctx context.Context, user *UserEntity)
	DeleteById(ctx context.Context, id int)
	DeleteAll(ctx context.Context, user []*UserEntity)
	DeleteAllById(ctx context.Context, ids []int)

	Save(ctx context.Context, user *UserEntity)
	SaveAll(ctx context.Context, user []*UserEntity)

	FindById(ctx context.Context, id int) *UserEntity
	FindAll(ctx context.Context) []*UserEntity
	FindAllById(ctx context.Context, ids []int)

	FindByFirstNameAndLastName(ctx context.Context, firstName, lastName string) *UserEntity
	// +shelf:query="FROM User WHERE FirstName = %1 AND LastName = %2"
	CustomQuery(ctx context.Context, firstName string) *UserEntity
}
