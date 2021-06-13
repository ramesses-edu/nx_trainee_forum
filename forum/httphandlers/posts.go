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

func PostsHandler(cfg *config.Config, db *gorm.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		rPath := r.URL.Path
		rePostsComments := regexp.MustCompile(`^\/posts\/\d+\/comments(\/)??$`)
		rePostsID := regexp.MustCompile(`^\/posts\/\d+(\/)??$`)
		rePosts := regexp.MustCompile(`^\/posts(\/)??$`)

		switch {
		case rePosts.Match([]byte(rPath)):
			switch r.Method {
			case http.MethodGet: //list posts with filters
				listPostsHTTP(db, w, r)
			case http.MethodPost: //create post in:json
				createPostHTTP(cfg, db, w, r)
			case http.MethodPut: //update post  in:json
				updatePostHTTP(cfg, db, w, r)
			default:
				ResponseError(w, http.StatusMethodNotAllowed, "")
				return
			}
		case rePostsID.Match([]byte(rPath)):
			switch r.Method {
			case http.MethodGet: // get posts/{id}
				getPostByIDHTTP(db, w, r)
			case http.MethodDelete: // delete posts/{id}
				deletePostHTTP(cfg, db, w, r)
			default:
				ResponseError(w, http.StatusMethodNotAllowed, "")
				return
			}
		case rePostsComments.Match([]byte(rPath)):
			switch r.Method {
			case http.MethodGet: // list comments like->/comments?postId={id}
				listPostCommentsHTTP(db, w, r)
			default:
				ResponseError(w, http.StatusMethodNotAllowed, "")
				return
			}
		default:
			ResponseError(w, http.StatusBadRequest, "")
		}
	})

}

//@Summary List posts
//@Description get posts
//@Produce json
//@Param userId query integer false "posts filter by user"
//@Param xml query string false "show data like XML"
//@success 200
//@Failure 400,404
//@Failure 500
//@Failure default
//@Router /posts/ [get]
//@Security ApiKeyAuth
func listPostsHTTP(DB *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var param map[string]interface{} = make(map[string]interface{})
	userId := r.FormValue("userId")
	if userId != "" {
		var err error
		param["userId"], err = strconv.Atoi(r.FormValue("userId"))
		if err != nil {
			ResponseError(w, http.StatusBadRequest, "")
			return
		}
	}
	var pp models.Posts
	result := pp.ListPosts(DB, param)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			ResponseError(w, http.StatusNotFound, "")
			return
		}
		ResponseError(w, http.StatusInternalServerError, "")
		return
	}
	if responseXML(r) {
		xmlWrite(w, pp)
	} else {
		jsonWrite(w, pp.Posts)
	}
}

//@Summary Show a posts
//@Description get post by ID
//@Produce json
//@Param id path integer true "Post ID"
//@Param xml query string false "show data like XML"
//@success 200
//@Failure 400,404
//@Failure 500
//@Failure default
//@Router /posts/{id} [get]
//@Security ApiKeyAuth
func getPostByIDHTTP(DB *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var param map[string]interface{} = make(map[string]interface{})
	var err error
	param["id"], err = strconv.Atoi(reNum.FindString(r.URL.Path))
	if err != nil {
		ResponseError(w, http.StatusBadRequest, "")
		return
	}
	var p models.Post
	result := p.GetPost(DB, param)
	if result.Error != nil {
		ResponseError(w, http.StatusNotFound, "")
		return
	}
	if responseXML(r) {
		xmlWrite(w, p)
	} else {
		jsonWrite(w, p)
	}
}

type createPostStruct struct {
	Title string
	Body  string
}
type updatePostStruct struct {
	ID    int
	Title string
	Body  string
}

//@Summary Create post
//@Description create post
//@Accept json
//@Produce json
//@Param RequestPost body createPostStruct true "JSON structure for creating post"
//@Success 200,201
//@Failure 400
//@Failure default
//@Router /posts/ [POST]
//@Security ApiKeyAuth
func createPostHTTP(cfg *config.Config, DB *gorm.DB, w http.ResponseWriter, r *http.Request) {
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
	var p models.Post
	err = json.Unmarshal(reqBody, &p)
	if err != nil {
		ResponseError(w, http.StatusInternalServerError, "")
		return
	}
	if p.Title == "" || p.Body == "" {
		ResponseError(w, http.StatusBadRequest, "")
		return
	}
	p.UserID = u.ID
	result := p.CreatePost(DB)
	if result.Error != nil {
		ResponseError(w, http.StatusBadRequest, "")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	jsonP, _ := json.MarshalIndent(p, "", "  ")
	fmt.Fprintln(w, string(jsonP))
}

//@Summary Update post
//@Description update post
//@Accept json
//@Produce json
//@Param RequestPost body updatePostStruct true "JSON structure for updating post"
//@Success 200
//@Failure 400
//@Failure default
//@Router /posts/ [put]
//@Security ApiKeyAuth
func updatePostHTTP(cfg *config.Config, DB *gorm.DB, w http.ResponseWriter, r *http.Request) {
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
	var p models.Post
	err = json.Unmarshal(reqBody, &p)
	if err != nil {
		ResponseError(w, http.StatusInternalServerError, "")
		return
	}
	var pUpd models.Post
	result := pUpd.GetPost(DB, map[string]interface{}{"id": p.ID})
	if result.Error != nil || result.RowsAffected == 0 {
		ResponseError(w, http.StatusBadRequest, "")
		return
	}
	if pUpd.UserID != u.ID {
		ResponseError(w, http.StatusBadRequest, "")
		return
	}
	result = p.UpdatePost(DB)
	if result.Error != nil {
		ResponseError(w, http.StatusBadRequest, "")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	jsonP, _ := json.MarshalIndent(p, "", "  ")
	fmt.Fprintln(w, string(jsonP))
}

//@Summary Delete post
//@Description delete post by ID
//@Param id path int true "ID of deleting post"
//@Success 200
//@Failure default
//@Router /posts/{id} [delete]
//@Security ApiKeyAuth
func deletePostHTTP(cfg *config.Config, DB *gorm.DB, w http.ResponseWriter, r *http.Request) {
	u := authorization.GetCurrentUser(cfg, DB, r)
	if u.ID == 0 {
		ResponseError(w, http.StatusInternalServerError, "")
		return
	}

	pID, err := strconv.Atoi(reNum.FindString(r.URL.Path))
	if err != nil {
		ResponseError(w, http.StatusInternalServerError, "")
		return
	}
	var pDel models.Post
	result := pDel.GetPost(DB, map[string]interface{}{"id": pID})
	if result.Error != nil || result.RowsAffected == 0 {
		ResponseError(w, http.StatusBadRequest, "")
		return
	}
	if pDel.UserID != u.ID {
		ResponseError(w, http.StatusBadRequest, "")
		return
	}
	var p models.Post = models.Post{ID: pID, UserID: u.ID}
	result = p.DeletePost(DB)
	if result.Error != nil {
		ResponseError(w, http.StatusInternalServerError, "")
		return
	}
	w.WriteHeader(http.StatusOK)
}

//@Summary List comments of post
//@Description List comments like request /comments?postId={id}
//@Param id path int true "ID of post"
//@Param xml query string false "show data like XML"
//@Router /posts/{id}/comments [get]
//@Success 200
//@Failure default
//@Security ApiKeyAuth
func listPostCommentsHTTP(DB *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var param map[string]interface{} = make(map[string]interface{})
	var err error
	param["postId"], err = strconv.Atoi(reNum.FindString(r.URL.Path))
	if err != nil {
		ResponseError(w, http.StatusBadRequest, "")
		return
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
