module github.com/Tarick/naca-publications

go 1.15

replace github.com/Tarick/naca-publications => ./

require (
	github.com/Tarick/naca-rss-feeds v0.0.1
	github.com/go-chi/chi v4.1.2+incompatible
	github.com/go-chi/cors v1.1.1
	github.com/go-chi/render v1.0.1
	github.com/go-chi/stampede v0.4.4
	github.com/go-ozzo/ozzo-validation/v4 v4.2.2
	github.com/gofrs/uuid v3.3.0+incompatible
	github.com/jackc/pgx/v4 v4.9.0
	github.com/spf13/cobra v1.1.0
	github.com/spf13/viper v1.7.1
	go.uber.org/zap v1.16.0
	golang.org/x/crypto v0.0.0-20201012173705-84dcc777aaee // indirect
	golang.org/x/sys v0.0.0-20201016160150-f659759dc4ca // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
)
