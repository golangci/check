package a

type A struct {
	name string
}

type B struct {
	A
	age int
}

func F(b B) string {
	_ = b.age
	return b.name
}
