package httphandlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"nx_trainee_forum/forum/application/config"
	"nx_trainee_forum/forum/httphandlers/authorization"
	"nx_trainee_forum/forum/models"
	"regexp"
	"strconv"

	"gorm.io/gorm"
)

func CommentsHandler(cfg *config.Config, db *gorm.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		rPath := r.URL.Path
		reCommentsID := regexp.MustCompile(`^\/comments\/\d+(\/)??$`)
		reComments := regexp.MustCompile(`^\/comments(\/)??$`)

		switch {
		case reComments.Match([]byte(rPath)):
			switch r.Method {
			case http.MethodGet: // list comments with filters
				listCommentsHTTP(db, w, r)
			case http.MethodPost: // create comment in:json
				createCommentHTTP(cfg, db, w, r)
			case http.MethodPut: // update comment in:json
				updateCommentHTTP(cfg, db, w, r)
			default:
				ResponseError(w, http.StatusMethodNotAllowed, "")
				return
			}
		case reCommentsID.Match([]byte(rPath)):
			switch r.Method {
			case http.MethodGet: // get comments/{id}
				getCommentByIDHTTP(db, w, r)
			case http.MethodDelete: // delete comments/{id}
				deleteCommentHTTP(cfg, db, w, r)
			default:
				ResponseError(w, http.StatusMethodNotAllowed, "")
				return
			}
		default:
			ResponseError(w, http.StatusBadRequest, "")
		}
	})

}

//@Summary List comments
//@description list comments with filtering
//@Param postId query int false "ID of post"
//@Param xml query string false "show data like XML"
//@Success 200
//@Failure default
//@Router /comments/ [get]
//@Security ApiKeyAuth
func listCommentsHTTP(DB *gorm.DB, w http.ResponseWriter, r *http.Request) {

	var param map[string]interface{} = make(map[string]interface{})
	postId := r.FormValue("postId")
	if postId != "" {
		var err error
		param["postId"], err = strconv.Atoi(r.FormValue("postId"))
		if err != nil {
			ResponseError(w, http.StatusBadRequest, "")
			return
		}
	}
	var cc models.Comments
	result := cc.ListComments(DB, param)
	if result.Error != nil {
		ResponseError(w, http.StatusInternalServerError, "")
		return
	}
	if responseXML(r) {
		xmlWrite(w, cc)
	} else {
		jsonWrite(w, cc.Comments)
	}
}

//@Summary Show comment
//@description Get comment by ID
//@Param id path int true "ID of comment"
//@Param xml query string false "show data like XML"
//@Success 200
//@Failure default
//@Router /comments/{id} [get]
//@Security ApiKeyAuth
func getCommentByIDHTTP(DB *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var param map[string]interface{} = make(map[string]interface{})
	var err error
	param["id"], err = strconv.Atoi(reNum.FindString(r.URL.Path))
	if err != nil {
		ResponseError(w, http.StatusBadRequest, "")
		return
	}
	var cmnt models.Comment
	result := cmnt.GetComment(DB, param)
	if result.Error != nil {
		ResponseError(w, http.StatusNotFound, "")
		return
	}

	if responseXML(r) {
		xmlWrite(w, cmnt)
	} else {
		jsonWrite(w, cmnt)
	}
}

type createCommentStruct struct {
	PostID int `json:"postId"`
	Name   string
	Email  string
	Body   string
}
type updateCommentStruct struct {
	ID    int
	Name  string
	Email string
	Body  string
}

//@Summary Create comment
//@description create comment
//@Accept json
//@Produce json
//@Param RequestPost body createCommentStruct true "JSON structure for creating post"
//@Success 200,201
//@Failure 400
//@Failure default
//@Router /comments/ [post]
//@Security ApiKeyAuth
func createCommentHTTP(cfg *config.Config, DB *gorm.DB, w http.ResponseWriter, r *http.Request) {
	u := authorization.GetCurrentUser(cfg, DB, r)
	if u.ID == 0 {
		ResponseError(w, http.StatusInternalServerError, "")
		return
	}

	reqBody, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		ResponseError(w, http.StatusBadRequest, "")
		return
	}
	var c models.Comment
	err = json.Unmarshal(reqBody, &c)
	if err != nil {
		ResponseError(w, http.StatusInternalServerError, "")
		return
	}
	if c.Name == "" || c.Email == "" || c.Body == "" || c.PostID == 0 {
		ResponseError(w, http.StatusBadRequest, "")
		return
	}
	c.UserID = u.ID
	result := c.CreateComment(DB)
	if result.Error != nil {
		ResponseError(w, http.StatusBadRequest, "")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	jsonC, _ := json.MarshalIndent(c, "", "  ")
	fmt.Fprintln(w, string(jsonC))
}

//@Summary Update comment
//@description update comment
//@Accept json
//@Produce json
//@Param RequestPost body updateCommentStruct true "JSON structure for creating post"
//@Success 200
//@Failure 400
//@Failure default
//@Router /comments/ [put]
//@Security ApiKeyAuth
func updateCommentHTTP(cfg *config.Config, DB *gorm.DB, w http.ResponseWriter, r *http.Request) {
	u := authorization.GetCurrentUser(cfg, DB, r)
	if u.ID == 0 {
		ResponseError(w, http.StatusInternalServerError, "")
		return
	}
	reqBody, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		ResponseError(w, http.StatusBadRequest, "")
		return
	}
	var c models.Comment
	err = json.Unmarshal(reqBody, &c)
	if err != nil {
		ResponseError(w, http.StatusInternalServerError, "")
		return
	}
	var cUpd models.Comment
	result := cUpd.GetComment(DB, map[string]interface{}{"id": c.ID})
	if result.Error != nil || result.RowsAffected == 0 {
		ResponseError(w, http.StatusBadRequest, "")
		return
	}
	if cUpd.UserID != u.ID {
		ResponseError(w, http.StatusBadRequest, "")
		return
	}
	result = c.UpdateComment(DB)
	if result.Error != nil {
		ResponseError(w, http.StatusBadRequest, "")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	jsonC, _ := json.MarshalIndent(c, "", "  ")
	fmt.Fprintln(w, string(jsonC))
}

//@Summary Delete comment
//@descripton delete comment by ID
//@Param id path int true "ID of deleting comment"
//@Success 200
//@Failure default
//@Router /comments/{id} [delete]
//@Security ApiKeyAuth
func deleteCommentHTTP(cfg *config.Config, DB *gorm.DB, w http.ResponseWriter, r *http.Request) {
	u := authorization.GetCurrentUser(cfg, DB, r)
	if u.ID == 0 {
		ResponseError(w, http.StatusInternalServerError, "")
		return
	}

	cID, err := strconv.Atoi(reNum.FindString(r.URL.Path))
	if err != nil {
		ResponseError(w, http.StatusBadRequest, "")
		return
	}
	var cDel models.Comment
	result := cDel.GetComment(DB, map[string]interface{}{"id": cID})
	if result.Error != nil || result.RowsAffected == 0 {
		ResponseError(w, http.StatusBadRequest, "")
		return
	}
	if cDel.UserID != u.ID {
		ResponseError(w, http.StatusBadRequest, "")
		return
	}
	var c models.Comment = models.Comment{ID: cID, UserID: u.ID}
	result = c.DeleteComment(DB)
	if result.Error != nil {
		ResponseError(w, http.StatusInternalServerError, "")
		return
	}
	w.WriteHeader(http.StatusOK)
}
