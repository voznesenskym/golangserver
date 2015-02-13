package message

import (
	"net/http"
	"github.com/emicklei/go-restful"
	)

type Message struct {
	Sender, Id, Message string
	Targets []string
}

type MessageResource struct {
	// normally one would use DAO (data access object)
	messages map[string]Message
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