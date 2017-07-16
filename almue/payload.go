package almue

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/he4d/almue/model"
)

//-- FLOOR PAYLOAD --//

type floorPayload struct {
	*model.Floor
	Shutters  shutterListPayload  `json:"shutters,omitempty"`
	Lightings lightingListPayload `json:"lightings,omitempty"`
}

type floorListPayload []*floorPayload

func (f *floorPayload) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (f *floorPayload) Bind(r *http.Request) error {
	return nil
}

func (a *Almue) newFloorListPayloadResponse(floors []*model.Floor) []render.Renderer {
	list := []render.Renderer{}
	for _, floor := range floors {
		list = append(list, a.newFloorPayloadResponse(floor))
	}
	return list
}

func (a *Almue) newFloorPayloadResponse(floor *model.Floor) *floorPayload {
	resp := &floorPayload{Floor: floor}

	if resp.Shutters == nil {
		if shutters, _ := a.store.GetShutterListOfFloor(floor.ID); shutters != nil {
			resp.Shutters = shutterListPayload{}
			for _, shutter := range shutters {
				resp.Shutters = append(resp.Shutters, a.newShutterPayloadResponse(shutter))
			}
		}
	}

	if resp.Lightings == nil {
		if lightings, _ := a.store.GetLightingListOfFloor(floor.ID); lightings != nil {
			resp.Lightings = lightingListPayload{}
			for _, lighting := range lightings {
				resp.Lightings = append(resp.Lightings, a.newLightingPayloadResponse(lighting))
			}
		}
	}

	return resp
}

//-- SHUTTER PAYLOAD --//
type shutterPayload struct {
	*model.Shutter
}

type shutterListPayload []*shutterPayload

func (s *shutterPayload) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *shutterPayload) Bind(r *http.Request) error {
	return nil
}

func (a *Almue) newShutterListPayloadResponse(shutters []*model.Shutter) []render.Renderer {
	list := []render.Renderer{}
	for _, shutter := range shutters {
		list = append(list, a.newShutterPayloadResponse(shutter))
	}
	return list
}

func (a *Almue) newShutterPayloadResponse(shutter *model.Shutter) *shutterPayload {
	resp := &shutterPayload{Shutter: shutter}

	return resp
}

//-- LIGHTING PAYLOAD --//
type lightingPayload struct {
	*model.Lighting
}

type lightingListPayload []*lightingPayload

func (l *lightingPayload) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (l *lightingPayload) Bind(r *http.Request) error {
	return nil
}

func (a *Almue) newLightingListPayloadResponse(lightings []*model.Lighting) []render.Renderer {
	list := []render.Renderer{}
	for _, lighting := range lightings {
		list = append(list, a.newLightingPayloadResponse(lighting))
	}
	return list
}

func (a *Almue) newLightingPayloadResponse(lighting *model.Lighting) *lightingPayload {
	resp := &lightingPayload{Lighting: lighting}

	return resp
}
