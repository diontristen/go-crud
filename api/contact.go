package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/diontristen/go-crud/postgres"
	"github.com/diontristen/go-crud/util"
)

// ContactFavorites is a field that contains a contact's favorites
type ContactFavorites struct {
	Colors []string `json:"colors"`
}

// Contact represents a Contact model in the database
type getContact struct {
	ID        int              `db:"id"`
	Name      string           `db:"name"`
	Address   string           `db:"address"`
	Phone     string           `db:"phone"`
	Favorites ContactFavorites `db:"favorites"`
	CreatedAt string           `db:"created_at"`
	UpdatedAt string           `db:"updated_at"`
}

func (a *ContactFavorites) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &a)
}

type getContactResponse []getContact

type postContact struct {
	Name      string           `json:"name" validate:"required"`
	Address   string           `json:"address" validate:"required"`
	Phone     string           `json:"phone" validate:"required"`
	Favorites ContactFavorites `json:"favorites"  validate:"required"`
}

type deleteContact struct {
	ID int `db:"id"`
}

type updateContact struct {
	ID     int               `json:"id"`
	Update updateContactJSON `json:"update"`
}

type updateContactJSON map[string]interface{}

func GetContact(w http.ResponseWriter, r *http.Request, ac *util.AppContext) {
	q := "SELECT * FROM contacts"
	args := map[string]interface{}{}
	var contact getContact
	lst := getContactResponse{}
	err := postgres.ForRows(q, args, ac.DB, &contact, func() (bool, error) {
		lst = append(lst, contact)
		return true, nil
	})

	if err != nil {
		util.WriteError(w, err)
	}
	fmt.Println(lst)

	util.WriteJSON(w, lst)
}

func PostContact(w http.ResponseWriter, r *http.Request, ac *util.AppContext) {
	data := postContact{}

	err := json.NewDecoder(r.Body).Decode(&data)

	if err != nil {
		util.ReturnBadRequest(w, "Could not parse JSON. Probably a malformed JSON.")
		return
	}

	validationErrors, err := util.ValidateData(data)

	if err != nil {
		util.ReturnBadRequest(w, err.Error())
		return
	}

	if validationErrors != nil {
		util.ReturnBadRequestJSON(w, validationErrors)
		return
	}

	isOkay := util.ValidateExisting(w, "name", data.Name, ac, "post")

	if !isOkay {
		return
	}

	var id int64

	datJSON, _ := json.Marshal(data.Favorites)

	insertSQL := `
		insert into "contacts" (
			name,
			address,
			phone,
			favorites
		)
		values ($1, $2, $3, $4)
		returning id;
	`

	sqlError := ac.DB.QueryRow(insertSQL, data.Name, data.Address, data.Phone, string(datJSON)).Scan(&id)
	if sqlError != nil {
		util.ReturnBadRequest(w, fmt.Sprint(sqlError))
		return
	}

	json.NewEncoder(w).Encode(id)
}

func UpdateContact(w http.ResponseWriter, r *http.Request, ac *util.AppContext) {

	data := updateContact{}
	err := json.NewDecoder(r.Body).Decode(&data)

	if err != nil {
		util.ReturnBadRequest(w, "Could not parse JSON. Probably a malformed JSON.")
		return
	}

	isOkay := util.ValidateExisting(w, "id", data.ID, ac, "update")

	if !isOkay {
		return
	}

	set := []string{}

	for k := range data.Update {
		set = append(set, fmt.Sprintf("%s = '%s'", k, data.Update[k]))
	}

	setStr := ""
	if len(set) > 0 {
		setStr = strings.Join(set, " , ")
	}

	q := util.TemplatePrintf(`
		UPDATE "contacts"
			SET {{.set}}
		WHERE
			id = {{.id}}
	`,
		map[string]interface{}{
			"set": setStr,
			"id":  data.ID,
		})

	fmt.Println(q)
	sqlError := ac.DB.QueryRow(q).Err()

	if sqlError != nil {
		util.ReturnBadRequest(w, fmt.Sprint(sqlError))
		return
	}

	w.WriteHeader(http.StatusNoContent)

}

func RemoveContact(w http.ResponseWriter, r *http.Request, ac *util.AppContext) {
	data := deleteContact{}
	err := json.NewDecoder(r.Body).Decode(&data)

	if err != nil {
		util.ReturnBadRequest(w, "Could not parse JSON. Probably a malformed JSON.")
		return
	}

	isOkay := util.ValidateExisting(w, "id", data.ID, ac, "update")

	if !isOkay {
		return
	}

	removeSQL := `
		DELETE FROM "contacts"
		WHERE
			id = $1
	`

	sqlError := ac.DB.QueryRow(removeSQL, data.ID).Err()

	if sqlError != nil {
		util.ReturnBadRequest(w, fmt.Sprint(sqlError))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
