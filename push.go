package main

import (
	"fmt"
	apns "github.com/anachronistic/apns"
	// "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	// "gopkg.in/mgo.v2/bson"
	"strings"
)

func sendPush(token string, id bson.ObjectId) {
	deviceToken := strings.Replace(token, "-", "", -1)
	fmt.Println(deviceToken)
	payload := apns.NewPayload()
	payload.Alert = "Hello, world!"
	payload.Badge = 42
	payload.Sound = "bingbong.aiff"

	pn := apns.NewPushNotification()
	pn.DeviceToken = deviceToken
	pn.AddPayload(payload)

	client := apns.NewClient("gateway.sandbox.push.apple.com:2195", "cert.pem", "key.pem")
	resp := client.Send(pn)

	alert, _ := pn.PayloadString()
	fmt.Println("  Alert:", alert)
	fmt.Println("  Success:", resp.Success)
	fmt.Println("  Error:", resp.Error)
}
