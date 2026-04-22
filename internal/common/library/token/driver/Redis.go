package driver

type Redis struct {
}

func (t *Redis) Construct(option map[string]any) *Redis {
	return t
}

func (t *Redis) Set(s string, i int, i2 int) {
	//TODO implement me
	panic("implement me")
}

func (t *Redis) Get(s string) any {
	//TODO implement me
	panic("implement me")
}

func (t *Redis) Check(i int, i2 int) {
	//TODO implement me
	panic("implement me")
}

func (t *Redis) Delete(s string) {
	//TODO implement me
	panic("implement me")
}

func (t *Redis) Clear(s string) {
	//TODO implement me
	panic("implement me")
}

func (t *Redis) Handler() {
	//TODO implement me
	panic("implement me")
}

func (t *Redis) GetEncryptedToken(s string) {
	//TODO implement me
	panic("implement me")
}

func (t *Redis) GetExpiredIn(i int) {
	//TODO implement me
	panic("implement me")
}
