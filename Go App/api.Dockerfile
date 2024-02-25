FROM golang:1.16-alpine
WORKDIR /app

# agregar algunos paquetes necesarios
RUN apk update && \
    apk add libc-dev && \
    apk add gcc && \
    apk add make

# evitar la reinstalación de proveedores en cada cambio en el código fuente
COPY ./go.mod go.sum ./
RUN go mod download && go mod verify

# Instale Compile Daemon para comenzar. Lo usaremos para observar los cambios en los archivos go.
RUN go get github.com/githubnemo/CompileDaemon

# Copie y cree la aplicación
COPY . .
COPY ./entrypoint.sh /entrypoint.sh

# esperarlo requiere bash, que Alpine no incluye de forma predeterminada. Utilice esperar en su lugar
ADD https://raw.githubusercontent.com/eficode/wait-for/v2.1.0/wait-for /usr/local/bin/wait-for
RUN chmod +rx /usr/local/bin/wait-for /entrypoint.sh

ENTRYPOINT [ "sh", "/entrypoint.sh" ]