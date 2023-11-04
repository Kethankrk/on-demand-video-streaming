package database

type Video struct {
	ID     string `bson:"_id"`
	Title  string `bson:"title"`
	IsDash bool   `bson:"isdash"`
}
