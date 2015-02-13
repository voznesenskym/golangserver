package main

import (
	"net/http"
	// "strconv"
	"github.com/emicklei/go-restful"
	"gopkg.in/mgo.v2/bson"
	// "gopkg.in/mgo.v2"
	"fmt"
	// apns "github.com/anachronistic/apns"
)

type Message struct {
	Id      bson.ObjectId   `json:"id"bson:"_id"`
	Sender  bson.ObjectId   `json:"senderid"bson:"senderid"`
	Message string          `json:"message"bson:"message"`
	Targets []bson.ObjectId `json:"targets"bson:"targets"`
}

type MessageResource struct {
	// normally one would use DAO (data access object)
	messages  map[string]Message
	dataStore *DataStore
}

func (m MessageResource) Register(container *restful.Container) {
	ws := new(restful.WebService)
	ws.
		Path("/messages").
		Doc("Get Messages").
		Consumes(restful.MIME_XML, restful.MIME_JSON).
		Produces(restful.MIME_JSON, restful.MIME_XML)

	ws.Route(ws.GET("/{message-id}").To(m.findMessage).
		// docs
		Doc("get a message").
		Operation("findMessage").
		Param(ws.PathParameter("message-id", "identifier of the message").DataType("string")).
		Writes(Message{})) // on the response

	ws.Route(ws.POST("/send").To(m.sendMessage).
		// docs
		Doc("send a message").
		Operation("sendMessage").
		Reads(Message{})) // from the request

	container.Add(ws)
}

func (m MessageResource) findMessage(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("message-id")
	msg := m.messages[id]
	if len(msg.Id) == 0 {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusNotFound, "404: Message could not be found.")
		return
	}
	response.WriteEntity(msg)
}

func (m *MessageResource) sendMessage(request *restful.Request, response *restful.Response) {
	msg := new(Message)
	fmt.Println(msg)
	err := request.ReadEntity(msg)
	if err != nil {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError, err.Error())
		return
	}

	id, err := insertMessage(m.dataStore, msg)
	if err != nil {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError, err.Error())
		return
	}
	db := m.dataStore.session
	c := db.DB(m.dataStore.mongoDB).C("user")
	// var result []int
	// err := collection.Find(bson.M{"gender": "F"}).Distinct("age", &result)
	for index, element := range msg.Targets {
		// index is the index where we are
		// element is the element from someSlice for where we are
		// fmt.Println("Id")
		// fmt.Println(element)
		// query := c.FindId(element)
		// fmt.Println("query")
		// fmt.Println(query)
		// user :=

		var result []string
		err := c.FindId(element).Distinct("token", &result)
		if err != nil {
			fmt.Println("err")
			fmt.Println(err)
		}

		fmt.Println("token")
		fmt.Println(result)

		sendPush(result[index], id)
	}

	response.WriteHeader(http.StatusCreated)
	response.WriteEntity(msg)

}
