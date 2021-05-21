package application

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"nx_trainee_forum/forum/application/config"
	"nx_trainee_forum/forum/httphandlers"
	"nx_trainee_forum/forum/httphandlers/middleware"
	"nx_trainee_forum/forum/models"

	httpSwagger "github.com/swaggo/http-swagger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Application struct {
	DB     *gorm.DB
	Config *config.Config
	Router *http.ServeMux
	Server *http.Server
}

func New() *Application {
	app := Application{}
	//init configuration
	app.Config = config.New()
	//init Router
	app.Router = http.NewServeMux()
	//init Server
	app.Server = &http.Server{
		Handler: app.Router,
		Addr:    app.Config.HostAddr,
	}
	return &app
}

func (app *Application) Start() {
	gormDialector := mysql.New(mysql.Config{
		DSN: fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", app.Config.DB.UserDB, app.Config.DB.PassDB, app.Config.DB.HostDB, app.Config.DB.PortDB, app.Config.DB.NameDB),
	})
	var err error
	app.DB, err = gorm.Open(gormDialector, &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	//init DB tables
	initDBTables(app.DB)
	//init Routers
	initRouters(app)
	//start Server
	fmt.Println("App start")
	app.Server.ListenAndServe()
}

func (app *Application) Close() {
	fmt.Println("App stop")
	sql, _ := app.DB.DB()
	sql.Close()
	app.Server.Shutdown(context.Background())
	app = nil

}

func initRouters(app *Application) {
	router := app.Router
	router.Handle("/", httphandlers.MainHandler(app.DB))
	router.Handle("/public", http.NotFoundHandler())
	router.Handle("/public/", httphandlers.PublicHandler())
	router.Handle("/logout/", httphandlers.LogoutHandler(app.DB))
	router.Handle("/auth/", httphandlers.Authentification(app.Config, app.DB))
	router.Handle("/posts", middleware.Authorization(app.DB, httphandlers.PostsHandler(app.DB)))
	router.Handle("/posts/", middleware.Authorization(app.DB, httphandlers.PostsHandler(app.DB)))
	router.Handle("/comments", middleware.Authorization(app.DB, httphandlers.CommentsHandler(app.DB)))
	router.Handle("/comments/", middleware.Authorization(app.DB, httphandlers.CommentsHandler(app.DB)))
	router.HandleFunc("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("localhost/swagger/doc.json"),
	))
}

func initDBTables(db *gorm.DB) {
	if !db.Migrator().HasTable(&models.Comment{}) {
		db.Migrator().CreateTable(&models.Comment{})
	}
	if !db.Migrator().HasTable(&models.Post{}) {
		db.Migrator().CreateTable(&models.Post{})
		db.Migrator().CreateConstraint(&models.Post{}, "Comments")
	} else {
		if !db.Migrator().HasConstraint(&models.Post{}, "Comments") {
			db.Migrator().CreateConstraint(&models.Post{}, "Comments")
		}
	}

	if !db.Migrator().HasTable(&models.User{}) {
		db.Migrator().CreateTable(&models.User{})
		db.Migrator().CreateConstraint(&models.User{}, "Posts")
		db.Migrator().CreateConstraint(&models.User{}, "Comments")
	} else {
		if !db.Migrator().HasConstraint(&models.User{}, "Posts") {
			db.Migrator().CreateConstraint(&models.User{}, "Posts")
		}
		if !db.Migrator().HasConstraint(&models.User{}, "Comments") {
			db.Migrator().CreateConstraint(&models.User{}, "Comments")
		}
	}
}
