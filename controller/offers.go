package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/maciekmm/HackYeah/middleware"
	"github.com/maciekmm/HackYeah/model"
	"github.com/maciekmm/HackYeah/utils"
)

type Offers struct {
	Database *gorm.DB
}

func (o *Offers) Register(router *mux.Router) {
	router.Handle("/", http.HandlerFunc(o.HandleGetAll)).Methods(http.MethodGet)
	router.Handle("/", middleware.RequiresAuth(model.RoleUser, http.HandlerFunc(o.HandleAdd))).Methods(http.MethodPost)
	router.Handle("/{id:[0-9]+}/", middleware.RequiresAuth(model.RoleUser, http.HandlerFunc(o.HandleGetSingle))).Methods(http.MethodGet)
	router.Handle("/{id:[0-9]+}/", middleware.RequiresAuth(model.RoleUser, http.HandlerFunc(o.HandlePatchSingle))).Methods(http.MethodPatch)
	router.Handle("/{id:[0-9]+}/", middleware.RequiresAuth(model.RoleUser, http.HandlerFunc(o.HandlePutSingle))).Methods(http.MethodPut)
	router.Handle("/{id:[0-9]+}/", middleware.RequiresAuth(model.RoleUser, http.HandlerFunc(o.HandleDelete))).Methods(http.MethodDelete)
	router.Handle("/user/{id:[0-9]+}/", middleware.RequiresAuth(model.RoleUser, http.HandlerFunc(o.HandleGetUser))).Methods(http.MethodGet)
}

func (o *Offers) HandleAdd(rw http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(middleware.ContextUserKey).(*model.User)
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	offer := model.Offer{}
	if err := decoder.Decode(&offer); err != nil {
		(&utils.ErrorResponse{
			Errors:      []string{model.ErrOfferInternal.Error()},
			DebugErrors: []string{err.Error()},
		}).Write(http.StatusBadRequest, rw)
		return
	}
	offer.UserID = user.ID

	if err := offer.Add(o.Database); err != nil {
		if res, ok := err.(*utils.ErrorResponse); ok {
			res.Write(http.StatusBadRequest, rw)
			return
		}
		(&utils.ErrorResponse{
			Errors:      []string{model.ErrOfferInternal.Error()},
			DebugErrors: []string{err.Error()},
		}).Write(http.StatusInternalServerError, rw)
		return
	}

	rw.WriteHeader(http.StatusOK)
	if byt, err := json.Marshal(offer); err == nil {
		rw.Write(byt)
	}
}

func (o *Offers) HandleDelete(rw http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(middleware.ContextUserKey).(*model.User)
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		(&utils.ErrorResponse{
			Errors:      []string{model.ErrOfferIDInvalid.Error()},
			DebugErrors: []string{err.Error()},
		}).Write(http.StatusBadRequest, rw)
		return
	}

	db := o.Database
	if user.Role < model.RoleAdmin {
		db = db.Where("user_id = ?", user.ID)
	}

	if res := db.Where("id = ?", id).Delete(&model.Offer{}); res.Error != nil {
		(&utils.ErrorResponse{
			Errors:      []string{model.ErrOfferInternal.Error()},
			DebugErrors: []string{res.Error.Error()},
		}).Write(http.StatusInternalServerError, rw)
		return
	}
	rw.WriteHeader(http.StatusOK)
}

func (o *Offers) HandleGetAll(rw http.ResponseWriter, r *http.Request) {
	offers := []model.Offer{}
	res := o.Database
	if res := res.Find(&offers); res.Error != nil {
		(&utils.ErrorResponse{
			Errors:      []string{model.ErrOfferInternal.Error()},
			DebugErrors: []string{res.Error.Error()},
		}).Write(http.StatusInternalServerError, rw)
		return
	}

	byt, err := json.Marshal(&offers)
	if err != nil {
		(&utils.ErrorResponse{
			Errors:      []string{model.ErrOfferInternal.Error()},
			DebugErrors: []string{err.Error()},
		}).Write(http.StatusInternalServerError, rw)
		return
	}
	rw.WriteHeader(http.StatusOK)
	rw.Write(byt)
}

func (o *Offers) HandleGetUser(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		(&utils.ErrorResponse{
			Errors:      []string{model.ErrOfferUserIDInvalid.Error()},
			DebugErrors: []string{err.Error()},
		}).Write(http.StatusBadRequest, rw)
		return
	}

	offers := []model.Offer{}
	res := o.Database
	if res := res.Where("user_id = ?", id).Find(&offers); res.Error != nil {
		(&utils.ErrorResponse{
			Errors:      []string{model.ErrOfferInternal.Error()},
			DebugErrors: []string{res.Error.Error()},
		}).Write(http.StatusInternalServerError, rw)
		return
	}

	byt, err := json.Marshal(&offers)
	if err != nil {
		(&utils.ErrorResponse{
			Errors:      []string{model.ErrOfferInternal.Error()},
			DebugErrors: []string{err.Error()},
		}).Write(http.StatusInternalServerError, rw)
		return
	}
	rw.WriteHeader(http.StatusOK)
	rw.Write(byt)
}

func (o *Offers) HandleGetSingle(rw http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(middleware.ContextUserKey).(*model.User)
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		(&utils.ErrorResponse{
			Errors:      []string{model.ErrOfferIDInvalid.Error()},
			DebugErrors: []string{err.Error()},
		}).Write(http.StatusBadRequest, rw)
		return
	}

	offer := model.Offer{}
	res := o.Database
	if res := res.First(&offer, uint(id)); res.Error != nil {
		(&utils.ErrorResponse{
			Errors:      []string{model.ErrOfferInternal.Error()},
			DebugErrors: []string{res.Error.Error()},
		}).Write(http.StatusInternalServerError, rw)
		return
	}
	offer.UserID = user.ID

	byt, err := json.Marshal(&offer)
	if err != nil {
		(&utils.ErrorResponse{
			Errors:      []string{model.ErrOfferInternal.Error()},
			DebugErrors: []string{err.Error()},
		}).Write(http.StatusInternalServerError, rw)
		return
	}
	rw.WriteHeader(http.StatusOK)
	rw.Write(byt)
}

func (o *Offers) HandlePatchSingle(rw http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(middleware.ContextUserKey).(*model.User)
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		(&utils.ErrorResponse{
			Errors:      []string{model.ErrOfferIDInvalid.Error()},
			DebugErrors: []string{err.Error()},
		}).Write(http.StatusBadRequest, rw)
		return
	}

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	offer := model.Offer{}
	if err := decoder.Decode(&offer); err != nil {
		(&utils.ErrorResponse{
			Errors:      []string{model.ErrOfferInternal.Error()},
			DebugErrors: []string{err.Error()},
		}).Write(http.StatusBadRequest, rw)
		return
	}

	offer.UserID = user.ID
	mod := model.Offer{}
	mod.ID = uint(id)

	db := o.Database
	if user.Role != model.RoleAdmin {
		mod.UserID = user.ID
		db = db.Where("user_id = ?", user.ID)
	}

	if res := db.Model(&mod).Updates(&offer); res.Error != nil {
		(&utils.ErrorResponse{
			Errors:      []string{model.ErrOfferInternal.Error()},
			DebugErrors: []string{res.Error.Error()},
		}).Write(http.StatusInternalServerError, rw)
		return
	}

	rw.WriteHeader(http.StatusOK)
}

func (o *Offers) HandlePutSingle(rw http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(middleware.ContextUserKey).(*model.User)
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		(&utils.ErrorResponse{
			Errors:      []string{model.ErrOfferIDInvalid.Error()},
			DebugErrors: []string{err.Error()},
		}).Write(http.StatusBadRequest, rw)
		return
	}

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	offer := model.Offer{}
	if err := decoder.Decode(&offer); err != nil {
		(&utils.ErrorResponse{
			Errors:      []string{model.ErrOfferInternal.Error()},
			DebugErrors: []string{err.Error()},
		}).Write(http.StatusBadRequest, rw)
		return
	}

	offer.ID = uint(id)
	db := o.Database
	if user.Role < model.RoleAdmin {
		db.Where("user_id = ?", user.ID)
	}
	if res := db.Save(&offer); res.Error != nil {
		(&utils.ErrorResponse{
			Errors:      []string{model.ErrOfferInternal.Error()},
			DebugErrors: []string{res.Error.Error()},
		}).Write(http.StatusInternalServerError, rw)
		return
	}
	rw.WriteHeader(http.StatusOK)
}
