package application

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net/http"
	"nx_trainee_forum/forum/application/config"
	"nx_trainee_forum/forum/httphandlers"
	"nx_trainee_forum/forum/httphandlers/middleware"
	"nx_trainee_forum/forum/models"
	"os"
	"time"

	httpSwagger "github.com/swaggo/http-swagger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Application struct {
	DB     *gorm.DB
	Config *config.Config
	Router *http.ServeMux
	server *http.Server
	ctx    context.Context
	cancel context.CancelFunc
}

func New() *Application {
	app := Application{}
	//init configuration
	app.Config = config.New()
	//init Router
	app.Router = http.NewServeMux()
	//init Server
	app.server = &http.Server{
		Handler: app.Router,
		Addr:    app.Config.HostAddr,
	}
	app.ctx, app.cancel = context.WithCancel(context.Background())
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
	initRouters(&app)
	return &app
}

func (app *Application) Start() {
	//start Server
	fmt.Println("App start")
	go app.server.ListenAndServe()
	var command string
	for {
		fmt.Print(">>: ")
		myscan := bufio.NewScanner(os.Stdin)
		myscan.Scan()
		command = myscan.Text()
		switch command {
		case "server shutdown":
			fmt.Println(command)
			app.Close()
			return
		default:
			fmt.Println("invalid command: " + command)
		}
	}
}

func (app *Application) Close() {
	sql, _ := app.DB.DB()
	sql.Close()
	ctxsd, cancel := context.WithTimeout(app.ctx, 5*time.Second)
	defer cancel()
	app.server.Shutdown(ctxsd)
	app = nil
}

func initRouters(app *Application) {
	router := app.Router
	router.Handle("/", httphandlers.MainHandler(app.DB, app.Config))
	router.Handle("/public", http.NotFoundHandler())
	router.Handle("/public/", httphandlers.PublicHandler())
	router.Handle("/logout/", httphandlers.LogoutHandler(app.Config, app.DB))
	router.Handle("/auth/", httphandlers.Authentication(app.Config, app.DB))
	router.Handle("/getapikey", middleware.Authorization(app.Config, app.DB, httphandlers.GetAPIKeyHandler(app.DB, app.Config)))
	router.Handle("/posts", middleware.Authorization(app.Config, app.DB, httphandlers.PostsHandler(app.Config, app.DB)))
	router.Handle("/posts/", middleware.Authorization(app.Config, app.DB, httphandlers.PostsHandler(app.Config, app.DB)))
	router.Handle("/comments", middleware.Authorization(app.Config, app.DB, httphandlers.CommentsHandler(app.Config, app.DB)))
	router.Handle("/comments/", middleware.Authorization(app.Config, app.DB, httphandlers.CommentsHandler(app.Config, app.DB)))
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
