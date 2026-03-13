package main

import (
	"flag"
	"fmt"
	"net/http"

	"wallet/service/wallet/api/internal/config"
	"wallet/service/wallet/api/internal/handler"
	"wallet/service/wallet/api/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/wallet-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	// Swagger UI: serve embedded static files
	server.AddRoute(rest.Route{
		Method: http.MethodGet,
		Path:   "/swagger/:file",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			file := r.URL.Path[len("/swagger/"):]
			swaggerHandler(file)(w, r)
		},
	})
	server.AddRoute(rest.Route{
		Method: http.MethodGet,
		Path:   "/swagger",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "/swagger/index.html", http.StatusFound)
		},
	})

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	fmt.Printf("Swagger UI: http://localhost:%d/swagger/index.html\n", c.Port)
	server.Start()
}
