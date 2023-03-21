package golaravel

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/CloudyKit/jet/v6"
	"github.com/PaulDong/golaravel/render"
	"github.com/PaulDong/golaravel/session"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

const version = "1.0.0"

type Golaravel struct {
	AppName  string
	Debug    bool
	Version  string
	ErrorLog *log.Logger
	InfoLog  *log.Logger
	RootPath string
	Routes   *chi.Mux
	Render   *render.Render
	JetViews *jet.Set
	Session  *scs.SessionManager
	DB       Database
	config   config
}

type config struct {
	port        string
	renderer    string
	cookie      cookieConfig
	sessionType string
	database    databaseConfig
}

func (g *Golaravel) New(rootPath string) error {
	pathConfig := initPaths{
		rootPath:    rootPath,
		folderNames: []string{"handlers", "migrations", "views", "data", "public", "tmp", "logs", "middleware"},
	}
	err := g.Init(pathConfig)
	if err != nil {
		return err
	}
	err = g.checkDotEnv(rootPath)
	if err != nil {
		return err
	}
	// read .env
	err = godotenv.Load(rootPath + "/.env")
	if err != nil {
		return err
	}

	//create loggers
	infoLog, errorLog := g.startLoggers()

	// connect to database
	if os.Getenv("DATABASE_TYPE") == "" {
		db, err := g.OpenDB(os.Getenv("DATABASE_TYPE"), g.BuildDSN())
		if err != nil {
			errorLog.Println(err)
			os.Exit(1)
		}
		g.DB = Database{
			DataType: os.Getenv("DATABASE_TYPE"),
			Pool:     db,
		}
	}
	g.InfoLog = infoLog
	g.ErrorLog = errorLog

	g.Debug, _ = strconv.ParseBool(os.Getenv("DEBUG"))
	g.Version = version
	g.RootPath = rootPath
	g.Routes = g.routes().(*chi.Mux)

	g.config = config{
		port:     os.Getenv("PORT"),
		renderer: os.Getenv("RENDERER"),
		cookie: cookieConfig{
			name:     os.Getenv("COOKIE_NAME"),
			lifetime: os.Getenv("COOKIE_LIFETIME"),
			persist:  os.Getenv("COOKIE_PERSIST"),
			secure:   os.Getenv("COOKIE_SECURE"),
			domain:   os.Getenv("COOKIE_DOMAIN"),
		},
		sessionType: os.Getenv("SESSION_TYPE"),
		database: databaseConfig{
			database: os.Getenv("DATABASE_TYPE"),
			dsn:      g.BuildDSN(),
		},
	}

	// create session
	sess := session.Session{
		CookieLifetime: g.config.cookie.lifetime,
		CookiePersist:  g.config.cookie.persist,
		CookieName:     g.config.cookie.name,
		SessionType:    g.config.sessionType,
		CookieDomain:   g.config.cookie.domain,
		DBPool:         g.DB.Pool,
	}

	g.Session = sess.InitSession()

	var views = jet.NewSet(
		jet.NewOSFileSystemLoader(fmt.Sprintf("%s/views", rootPath)),
		jet.InDevelopmentMode(),
	)

	g.JetViews = views

	g.createRenderer()

	return nil
}

func (g *Golaravel) Init(p initPaths) error {
	root := p.rootPath
	for _, path := range p.folderNames {
		// create folder if it doesn't exist
		err := g.CreateDirIfNotExist(root + "/" + path)
		if err != nil {
			return err
		}
	}
	return nil
}

func (g *Golaravel) ListenAndServe() {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", os.Getenv("PORT")),
		ErrorLog:     g.ErrorLog,
		Handler:      g.Routes,
		IdleTimeout:  30 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 600 * time.Second,
	}

	defer g.DB.Pool.Close()

	g.InfoLog.Printf("Listening on port %s", os.Getenv("PORT"))
	err := srv.ListenAndServe()
	if err != nil {
		g.ErrorLog.Fatal(err)
	}
}

func (g *Golaravel) checkDotEnv(path string) error {
	err := g.CreateFileIfNotExist(fmt.Sprintf("%s/.env", path))
	if err != nil {
		return err
	}
	return nil
}

func (g *Golaravel) startLoggers() (*log.Logger, *log.Logger) {
	var infoLog *log.Logger
	var errorLog *log.Logger
	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	return infoLog, errorLog
}

func (g *Golaravel) createRenderer() {
	myRenderer := render.Render{
		Renderer: g.config.renderer,
		RootPath: g.RootPath,
		Port:     g.config.port,
		JetViews: g.JetViews,
		Session:  g.Session,
	}

	g.Render = &myRenderer
}

func (g *Golaravel) BuildDSN() string {
	var dsn string
	switch os.Getenv("DATABASE_TYPE") {
	case "postgres", "postgresql":
		dsn = fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s timezone=UTC connect_timeout=5",
			os.Getenv("DATABASE_HOST"),
			os.Getenv("DATABASE_PORT"),
			os.Getenv("DATABASE_USER"),
			os.Getenv("DATABASE_NAME"),
			os.Getenv("DATABASE_SSL_MODE"))
		if os.Getenv("DATABASE_PASS") != "" {
			dsn = fmt.Sprintf("%s password=%s", dsn, os.Getenv("DATABASE_PASS"))
		}
	case "mysql", "mariadb":

	default:

	}
	return dsn
}
