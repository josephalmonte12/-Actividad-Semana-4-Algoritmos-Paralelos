wait-for "${DATABASE_HOST}:${DATABASE_PORT}" -- "$@"

# Observe sus archivos .go e invoque go build si los archivos cambiaron.
CompileDaemon --build="go build -o main main.go"  --command=./main