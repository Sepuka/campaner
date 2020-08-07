package context

type Context struct {
	User *User
}

type User struct {
	Timezone int32
}
