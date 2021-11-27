package cont

type Role int

const (
	Reporter Role = iota + 1
	Developer
	Maintainer
	Owner
	Admin
)

func (r Role) String() string {
	return [...]string{"Reporter", "Developer", "Maintainer", "Owner", "Admin"}[r-1]
}

func (r Role) EnumIndex() int {
	return int(r)
}
