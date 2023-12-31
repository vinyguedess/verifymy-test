package main

import (
	"context"
	"net"
	"net/http"

	gorillaHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"

	"verifymy-golang-test/handlers"
	"verifymy-golang-test/middlewares"
	"verifymy-golang-test/repositories"
	"verifymy-golang-test/services"
)

func main() {
	godotenv.Load()

	fx.New(
		fx.Provide(
			func() (*zap.Logger, error) {
				zapConfig := zap.NewProductionConfig()
				zapConfig.EncoderConfig.MessageKey = "message"

				return zapConfig.Build()
			},
			NewHTTPServer,
			fx.Annotate(
				NewServeMux,
				fx.ParamTags(`*zap.Logger`, `group:"routes"`),
			),

			AsRoute(handlers.NewHealthCheckHandler),
			AsRoute(handlers.NewSignUpHandler),
			AsRoute(handlers.NewSignInHandler),
			AsRoute(handlers.NewShowProfileHandler),
			AsRoute(handlers.NewUpdateProfileHandler),
			AsRoute(handlers.NewListUsersHandler),
			AsRoute(handlers.NewDeleteUserByIdHandler),
		),
		fx.WithLogger(
			func(log *zap.Logger) fxevent.Logger {
				return &fxevent.ZapLogger{Logger: log}
			},
		),
		fx.Invoke(func(*http.Server) {}),
		repositories.Module,
		services.Module,
		handlers.Module,
	).Run()
}

func NewHTTPServer(lc fx.Lifecycle, mux *mux.Router, log *zap.Logger) *http.Server {
	corsMiddleware := gorillaHandlers.CORS(
		gorillaHandlers.AllowCredentials(),
		gorillaHandlers.AllowedMethods(
			[]string{
				http.MethodGet,
				http.MethodPost,
				http.MethodPatch,
				http.MethodDelete,
				http.MethodOptions,
			},
		),
		gorillaHandlers.AllowedHeaders(
			[]string{
				"Content-Type",
				"Origin",
				"Sec-fetch-site",
			},
		),
	)

	server := &http.Server{Addr: ":8080", Handler: corsMiddleware(mux)}
	lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				ln, err := net.Listen("tcp", server.Addr)
				if err != nil {
					return err
				}

				log.Info("Starting HTTP server", zap.String("addr", server.Addr))
				go server.Serve(ln)
				return nil
			},
			OnStop: func(ctx context.Context) error {
				return server.Shutdown(ctx)
			},
		},
	)
	return server
}

func NewServeMux(
	authService services.AuthService,
	handlers []handlers.Handler,
) *mux.Router {
	mux := mux.NewRouter()
	mux.Use(middlewares.AuthMiddleware(authService))

	mux.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))),
	)

	mux.PathPrefix("/swagger/").Handler(
		httpSwagger.Handler(
			httpSwagger.URL("/static/doc.json"),
			httpSwagger.DeepLinking(true),
			httpSwagger.DocExpansion("none"),
			httpSwagger.DomID("swagger-ui"),
		),
	).Methods(http.MethodGet)

	for _, h := range handlers {
		mux.Handle(h.Route(), h).Methods(h.Method()...)
	}

	return mux
}

func AsRoute(f interface{}) interface{} {
	return fx.Annotate(
		f,
		fx.As(new(handlers.Handler)),
		fx.ResultTags(`group:"routes"`),
	)
}
