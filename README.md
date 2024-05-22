## sqlx + squirrel

Repositorio para la charla sobre cómo usar sqlx y squirrel.

Para probar los ejemplos descritos en este repositorio, clona el repositorio:

```sh
git clone --recurse-submodules https://github.com/danvergara/sqlx_squirrel
```

Ve al directorio ` pagila` y ejecuta:

```sh
docker compose up -d
```

Tendrás la base de datos de Pagila (un port a Postgres de la famosa base de datos [Sakila](https://dev.mysql.com/doc/sakila/en/)) corriendo en un contenedor de postgres.

Para correr la presentación, ejecuta [present](https://pkg.go.dev/golang.org/x/tools/present):

```sh
present
```
