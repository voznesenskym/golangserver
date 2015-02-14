package main

import (
	"fmt"
	"github.com/emicklei/go-restful"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
)

type User struct {
	Id          bson.ObjectId `json:"id"bson:"_id"`
	Name        string        `json:"name" bson:"name"`
	DeviceToken string        `json:"token" bson:"token"`
}

type UserResource struct {
	// normally one would use DAO (data access object)
	users     map[string]User
	dataStore *DataStore
}

type DataStore struct {
	mongoDB string
	session *mgo.Session
}

func (u UserResource) Register(container *restful.Container) {
	ws := new(restful.WebService)
	ws.
		Path("/users").
		Doc("Manage Users").
		Consumes(restful.MIME_XML, restful.MIME_JSON, "application/octet-stream", "application/x-msgpack", "multipart/form-data", "application/x-www-form-urlencoded").
		Produces(restful.MIME_JSON, restful.MIME_XML, "application/octet-stream", "application/x-msgpack", "multipart/form-data", "application/x-www-form-urlencoded") // you can specify this per route as well

	ws.Route(ws.GET("/{user-id}").To(u.findUser).
		// docs
		Doc("get a user").
		Operation("findUser").
		Param(ws.PathParameter("user-id", "identifier of the user").DataType("string")).
		Writes(User{})) // on the response

	ws.Route(ws.PUT("/{user-id}").To(u.updateUser).
		// docs
		Doc("update a user").
		Operation("updateUser").
		Param(ws.PathParameter("user-id", "identifier of the user").DataType("string")).
		ReturnsError(409, "duplicate user-id", nil).
		Reads(User{})) // from the request

	ws.Route(ws.POST("").To(u.createUser).
		// docs
		Doc("create a user").
		Operation("createUser").
		Reads(User{})) // from the request

	ws.Route(ws.DELETE("/{user-id}").To(u.removeUser).
		// docs
		Doc("delete a user").
		Operation("removeUser").
		Param(ws.PathParameter("user-id", "identifier of the user").DataType("string")))

	container.Add(ws)
}

func (u UserResource) findUser(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("user-id")
	usr := u.users[id]
	if len(usr.Id) == 0 {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusNotFound, "404: User could not be found.")
		return
	}
	response.WriteEntity(usr)
}

func (u *UserResource) createUser(request *restful.Request, response *restful.Response) {
	usr := new(User)
	fmt.Println(usr)
	err := request.ReadEntity(usr)
	if err != nil {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError, err.Error())
		return
	}

	err = insertUser(u.dataStore, usr)
	if err != nil {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError, err.Error())
		return
	}
	response.WriteHeader(http.StatusCreated)
	response.WriteEntity(usr)
}

func (u *UserResource) updateUser(request *restful.Request, response *restful.Response) {
	usr := new(User)
	err := request.ReadEntity(&usr)
	if err != nil {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError, err.Error())
		return
	}

	// u.users[usr.Id] = *usr
	response.WriteEntity(usr)
}

func (u *UserResource) removeUser(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("user-id")
	delete(u.users, id)
}

func main() {

	mongoDB := "gotest"
	// session, err := mgo.Dial("localhost:27017")
	session, err := mgo.Dial("46.28.202.117:27017")
	if err != nil {
		panic(err)
	}

	dataStore := DataStore{mongoDB, session}
	defer session.Close()

	wsContainer := restful.NewContainer()
	u := UserResource{map[string]User{}, &dataStore}
	m := MessageResource{map[string]Message{}, &dataStore}
	// m..Register(wsContainer)
	m.Register(wsContainer)
	u.Register(wsContainer)

	log.Printf("start listening on localhost:5000")
	server := &http.Server{Addr: ":5000", Handler: wsContainer}
	log.Fatal(server.ListenAndServe())
}
