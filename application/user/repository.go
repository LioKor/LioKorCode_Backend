package user

type Repository interface {
	GetId(token string) (uint64, error)
	StoreSession(token string, uid uint64) error
}
