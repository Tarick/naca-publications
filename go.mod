module github.com/Tarick/naca-publications

go 1.15

replace github.com/Tarick/naca-publications => ./

// replace github.com/Tarick/naca-rss-feeds => ../naca-rss-feeds

require (
	github.com/Tarick/naca-rss-feeds v0.0.2-0.20201231093705-87e789c07d59
	github.com/asaskevich/govalidator v0.0.0-20200907205600-7a23bdc65eef
	github.com/go-chi/chi v1.5.1
	github.com/go-chi/cors v1.1.1
	github.com/go-chi/render v1.0.1
	github.com/go-chi/stampede v0.4.4
	github.com/go-ozzo/ozzo-validation/v4 v4.3.0
	github.com/gofrs/uuid v4.0.0+incompatible
	github.com/jackc/pgx/v4 v4.10.1
	github.com/prometheus/client_golang v1.9.0
	github.com/spf13/cobra v1.1.1
	github.com/spf13/viper v1.7.1
	go.uber.org/zap v1.16.0
)
