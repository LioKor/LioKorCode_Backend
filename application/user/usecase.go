package user

type UseCase interface {
	GetId(token string) (uint64, error)
	StoreSession(uid uint64) (string, error)
}
