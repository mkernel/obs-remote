package main

import (
	"log"
	"net/http"

	obsws "github.com/christopher-dG/go-obs-websocket"
)

var client obsws.Client

func main() {
	client = obsws.Client{Host: "defiant.local", Port: 4444}
	if err := client.Connect(); err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect()
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("html/"))))
	http.HandleFunc("/api/", httpHandler)

	http.ListenAndServe(":8080", nil)

	/*req := obsws.NewSetCurrentSceneRequest("Dicecam")
	req.Send(client)
	req.Receive()*/
}

func httpHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/api/scene/gamecam" {
		req := obsws.NewSetCurrentSceneRequest("Gamecam")
		req.Send(client)
		req.Receive()
	}
	if r.URL.Path == "/api/scene/score" {
		req := obsws.NewSetCurrentSceneRequest("Fullscreen Scoreboard")
		req.Send(client)
		req.Receive()
	}
	if r.URL.Path == "/api/scene/dice" {
		req := obsws.NewSetCurrentSceneRequest("Dicecam")
		req.Send(client)
		req.Receive()
	}
	if r.URL.Path == "/api/scene/face" {
		req := obsws.NewSetCurrentSceneRequest("Facecam")
		req.Send(client)
		req.Receive()
	}
	if r.URL.Path == "/api/toggle/score" {
		toggle("Gamecam", "Scoreboard")
	}
	if r.URL.Path == "/api/toggle/grid" {
		toggle("Gamecam", "Grid")
	}
	if r.URL.Path == "/api/toggle/face" {
		toggle("Gamecam", "Facetime Camera")
	}
	if r.URL.Path == "/api/toggle/dice" {
		toggle("Gamecam", "dicecam")
	}
	http.Redirect(w, r, "/public/index.html", 302)
}

func toggle(scene string, item string) {
	req := obsws.NewGetSceneItemPropertiesRequest(scene, item)
	req.Send(client)
	data, _ := req.Receive()
	if data.Visible {
		nreq := obsws.NewSetSceneItemPropertiesRequest(
			scene, item,
			data.PositionX, data.PositionY,
			data.PositionAlignment, data.Rotation,
			data.ScaleX, data.ScaleY, data.CropTop,
			data.CropBottom, data.CropLeft, data.CropRight, false,
			data.Locked, data.BoundsType, data.BoundsAlignment, data.BoundsX, data.BoundsY)
		nreq.Send(client)
		nreq.Receive()

	} else {
		nreq := obsws.NewSetSceneItemPropertiesRequest(
			scene, item,
			data.PositionX, data.PositionY,
			data.PositionAlignment, data.Rotation,
			data.ScaleX, data.ScaleY, data.CropTop,
			data.CropBottom, data.CropLeft, data.CropRight, true,
			data.Locked, data.BoundsType, data.BoundsAlignment, data.BoundsX, data.BoundsY)
		nreq.Send(client)
		nreq.Receive()
	}
}
