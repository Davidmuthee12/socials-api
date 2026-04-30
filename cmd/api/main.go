package main

import (
	"expvar"
	"runtime"
	"time"

	"github.com/Davidmuthee12/socials/internal/auth"
	"github.com/Davidmuthee12/socials/internal/db"
	"github.com/Davidmuthee12/socials/internal/env"
	"github.com/Davidmuthee12/socials/internal/mailer"
	ratelimiter "github.com/Davidmuthee12/socials/internal/rateLimiter"
	"github.com/Davidmuthee12/socials/internal/store"
	cache "github.com/Davidmuthee12/socials/internal/store/cache"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

const version = "1.1.4"

//	@title	Socials APP API

//	@description	This API for the social app.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@BasePath					/v1
//
//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						Authorization
//	@description

func main() {
	cfg := config{
		addr:        env.GetString("ADDR", ":8000"),
		apiURL:      env.GetString("EXTERNAL_URL", "localhost:8080"),
		frontendURL: env.GetString("FRONTEND_URL", "http://localhost:5173"),
		db: dbConfig{
			addr:          env.GetString("DB_ADDR", "postgres://user:adminpassword@localhost/social?sslmode=disable"),
			maxOpenConns:  env.GetInt("DB_MAX_OPEN_CONNS", 25),
			maxIddleConns: env.GetInt("DB_MAX_IDLE_CONNS", 25),
			maxIddleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		redisCfg: redisConfig{
			addr:    env.GetString("REDIS_ADDR", "localhost:6379"),
			pw:      env.GetString("REDIS_PS", ""),
			db:      env.GetInt("REDIS_DB", 0),
			enabled: env.GetBool("REDIS_ENABLED", false),
		},
		env: env.GetString("ENV", "development"),
		mail: mailConfig{
			exp:       time.Hour * 24 * 3, // user has 3 days to accept invitation
			fromEmail: env.GetString("FROM_EMAIL", ""),
			sendGrid: sendGridConfig{
				apiKey: env.GetString("SENDGRID_API_KEY", ""),
			},
			mailTrap: mailTrapConfig{
				apiKey: env.GetString("MAIL_TRAP_API_KEY", ""),
			},
		},
		auth: authConfig{
			basic: basicConfig{
				user: env.GetString("AUTH_BASIC_USER", "admin"),
				pass: env.GetString("AUTH_BASIC_PASS", "admin"),
			},
			token: tokenConfig{
				secret: env.GetString("AUTH_TOKEN_SECRET", "example"),
				exp:    time.Hour * 24 * 3, //3 days
				iss:    "socials",
			},
		},
		rateLimiter: ratelimiter.Config{
			RequestPerTimeFrame: env.GetInt("RATELIMITER_REQUEST_COUNT", 20),
			TimeFrame:           time.Second * 5,
			Enabled:             env.GetBool("RATE_LIMITER_ENABLED", true),
		},
	}

	// Logger
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	// DB
	db, err := db.New(
		cfg.db.addr,
		cfg.db.maxOpenConns,
		cfg.db.maxIddleConns,
		cfg.db.maxIddleTime,
	)

	if err != nil {
		logger.Fatal(err)
	}

	defer db.Close()
	logger.Info("Database connection pool established")

	// Cache
	var rdb *redis.Client
	if cfg.redisCfg.enabled {
		rdb = cache.NewRedisClient(cfg.redisCfg.addr, cfg.redisCfg.pw, cfg.redisCfg.db)
		logger.Info("Redis cache connection pool established")

		defer rdb.Close()
	}

	ratelimiter := ratelimiter.NewFixedWindowLimiter(
		cfg.rateLimiter.RequestPerTimeFrame,
		cfg.rateLimiter.TimeFrame,
	)

	store := store.NewStorage(db)
	cacheStorage := cache.NewRedisStorage(rdb)

	// mailer -- im using mailTrap since sendGrid api encountered issues
	mailer := mailer.NewSendgrid(cfg.mail.sendGrid.apiKey, cfg.mail.fromEmail)
	// mailtrap, err := mailer.NewMailTrapClient(cfg.mail.mailTrap.apiKey, cfg.mail.fromEmail)
	// if err != nil {
	// 	logger.Fatal(err)
	// }

	jwtAuthenticator := auth.NewJWTAuthenticator(cfg.auth.token.secret, cfg.auth.token.iss, cfg.auth.token.iss)

	app := &application{
		config:       cfg,
		store:        store,
		cacheStorage: cacheStorage,
		logger:       logger,
		mailer:       mailer,
		// mailer: mailtrap,
		authenticator: jwtAuthenticator,
		rateLimiter:   ratelimiter,
	}

	// Metrics collections
	expvar.NewString("version").Set(version)
	expvar.Publish("database", expvar.Func(func() any {
		return db.Stats()
	}))

	expvar.Publish("goroutines", expvar.Func(func() any {
		return runtime.NumGoroutine()
	}))

	mux := app.mount()
	logger.Fatal(app.run(mux))
}
