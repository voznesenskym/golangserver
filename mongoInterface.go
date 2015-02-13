package main

import (
	"fmt"
	// "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func insertUser(datastore *DataStore, user *User) error {

	fmt.Printf("Inserting User")
	db := datastore.session
	c := db.DB(datastore.mongoDB).C("user")

	usr := *user
	usr.Id = bson.NewObjectId()
	fmt.Printf("setting User id")
	fmt.Printf("Created new NewObjectId %s", usr.Id)
	// return c.Update(&usr)
	return c.Insert(&usr)
}

func insertMessage(datastore *DataStore, message *Message) (id bson.ObjectId, err error) {
	fmt.Println("Inserting Message")
	db := datastore.session
	c := db.DB(datastore.mongoDB).C("message")

	msg := *message
	msg.Id = bson.NewObjectId()
	fmt.Println("setting Message id")
	fmt.Println("Created new NewObjectId %s", msg.Id)
	// return c.Update(&msg)
	err = c.Insert(&msg)
	return msg.Id, err
}
