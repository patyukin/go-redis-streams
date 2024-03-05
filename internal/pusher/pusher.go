package pusher

type PushStorage interface {
	GetIdsForPush() error
}

type Pusher struct {
	storage PushStorage
}

func NewPusher(storage PushStorage) *Pusher {
	return &Pusher{
		storage: storage,
	}
}

func (p *Pusher) Push() {

}
