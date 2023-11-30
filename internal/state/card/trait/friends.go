package trait

const (
	FriendsTrait = "Friends"
)

type FriendsArgs struct {
	Location string
	Trait    struct {
		Type string
		Args interface{}
	}
}
