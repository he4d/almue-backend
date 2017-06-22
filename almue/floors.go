package almue

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/render"
	"github.com/he4d/almue/model"
)

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
		if shutters, _ := a.store.GetShutterList(); shutters != nil {
			resp.Shutters = shutterListPayload{}
			for _, shutter := range shutters {
				resp.Shutters = append(resp.Shutters, a.newShutterPayloadResponse(shutter))
			}
		}
	}

	if resp.Lightings == nil {
		if lightings, _ := a.store.GetLightingList(); lightings != nil {
			resp.Lightings = lightingListPayload{}
			for _, lighting := range lightings {
				resp.Lightings = append(resp.Lightings, a.newLightingPayloadResponse(lighting))
			}
		}
	}

	return resp
}

func (a *Almue) getAllFloors(w http.ResponseWriter, r *http.Request) {
	floors, err := a.store.GetFloorList()
	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
	if err := render.RenderList(w, r, a.newFloorListPayloadResponse(floors)); err != nil {
		render.Render(w, r, ErrRender(err))
	}
}

func (a *Almue) createFloor(w http.ResponseWriter, r *http.Request) {
	f := &floorPayload{}
	if err := render.Bind(r, f); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	id, err := a.store.CreateFloor(f.Floor)
	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	floor := f.Floor
	floor.ID = id

	render.Status(r, http.StatusCreated)
	render.Render(w, r, a.newFloorPayloadResponse(floor))
}

func (a *Almue) getFloor(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	floor := ctx.Value(contextKeyFloor).(*model.Floor)

	if err := render.Render(w, r, a.newFloorPayloadResponse(floor)); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
}

//TODO: hier muss weiter gearbeitet werden!!
func (a *Almue) updateFloor(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	oldFloor, ok := ctx.Value(contextKeyFloor).(*model.Floor)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	f := new(model.Floor)
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(f); err != nil {
		respondWithError(w, 400)
		return
	}
	defer r.Body.Close()
	f.ID = oldFloor.ID

	if err := a.store.UpdateFloor(f); err != nil {
		respondWithError(w, 500)
		return
	}

	floor, err := a.store.GetFloor(oldFloor.ID)
	if err != nil {
		respondWithError(w, 500)
		return
	}

	respondWithJSON(w, http.StatusOK, floor)
}

func (a *Almue) deleteFloor(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	floor, ok := ctx.Value(contextKeyFloor).(*model.Floor)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	//TODO: update a.devices (delete related devices)

	if err := a.store.DeleteFloor(floor.ID); err != nil {
		respondWithError(w, 500)
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}
