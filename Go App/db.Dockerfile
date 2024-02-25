# comenzar con la imagen base
FROM mysql:8.0.23

# importar datos al contenedor
# Todos los scripts en docker-entrypoint-initdb.d/ se ejecutan automáticamente durante el inicio del contenedor
COPY ./database/*.sql /docker-entrypoint-initdb.d/