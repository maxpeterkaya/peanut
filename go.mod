module peanut

go 1.24

require (
	github.com/BurntSushi/toml v1.5.0
	github.com/go-chi/chi/v5 v5.2.2
	github.com/google/go-github/v67 v67.0.0
	github.com/joho/godotenv v1.5.1
	github.com/maxpeterkaya/peanut/common v0.0.0-00010101000000-000000000000
	github.com/rs/zerolog v1.34.0
)

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.4 // indirect
	github.com/prometheus/client_golang v1.16.0 // indirect
	github.com/prometheus/client_model v0.3.0 // indirect
	github.com/prometheus/common v0.42.0 // indirect
	github.com/prometheus/procfs v0.10.1 // indirect
	github.com/toshi0607/chi-prometheus v0.1.4 // indirect
	golang.org/x/sys v0.12.0 // indirect
	google.golang.org/protobuf v1.30.0 // indirect
)

replace github.com/maxpeterkaya/peanut/common => ./common
