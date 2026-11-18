# Verificación del proyecto

**Fecha de ejecución:** 18 de noviembre de 2026  
**Responsable:** Carlos Antonio Ramírez Galvis

## Entorno utilizado

- Docker Engine 29.8.2
- Docker Compose 5.5.1
- MySQL 8.0
- API construida desde el `Dockerfile` del repositorio
- Pruebas realizadas contra `http://localhost:8080`

## Estado de contenedores

Después de ejecutar `docker compose up -d --build`, se obtuvo:

| Contenedor | Estado | Puerto |
|---|---|---|
| `biblioteca_mysql` | Up (healthy) | 3306 |
| `biblioteca_api` | Up | 8080 |
| `biblioteca_adminer` | Up | 8081 |

La base existente se respaldó antes de reconstruir los contenedores. No fue necesario borrar su volumen.

## Verificación desde una base vacía

También se inició un contenedor temporal de MySQL con un volumen vacío y la carpeta `migrations` montada en `/docker-entrypoint-initdb.d`. MySQL ejecutó correctamente `001_init.sql` y después `002_laptops.sql`.

El resultado incluyó las tablas `users`, `books`, `loans`, `laptops` y `laptop_loans`, junto con los seis libros de ejemplo. Después de consultar el esquema, el contenedor y el volumen temporal fueron eliminados. El resultado está en [`evidencias/resultado-migraciones.txt`](evidencias/resultado-migraciones.txt).

## Pruebas automatizadas de la API

El script [`../scripts/verificar-proyecto.ps1`](../scripts/verificar-proyecto.ps1) crea usuarios únicos, asigna un administrador de prueba, ejecuta los casos y detiene la ejecución si un código HTTP no coincide.

| Área | Caso comprobado | Esperado | Resultado |
|---|---|---:|---|
| Infraestructura | Health y ping | 200 | OK |
| Autenticación | Registro, login y perfil JWT | 200/201 | OK |
| Regresión | Listar libros | 200 | OK |
| Regresión | Prestar y devolver un libro | 201/200 | OK |
| Seguridad | Consultar portátiles sin token | 401 | OK |
| Roles | Estudiante intenta crear portátil | 403 | OK |
| Catálogo | Administrador crea portátil | 201 | OK |
| Catálogo | Serial repetido | 409 | OK |
| Catálogo | Usuario autenticado lista equipos | 200 | OK |
| Préstamo | Solicitar equipo disponible | 201 | OK |
| Préstamo | Segundo usuario solicita equipo ocupado | 409 | OK |
| Préstamo | Consultar mis préstamos | 200 | OK |
| Devolución | Otro usuario intenta devolver | 403 | OK |
| Devolución | ID inexistente | 404 | OK |
| Devolución | Propietario devuelve el equipo | 200 | OK |
| Devolución | Repetir la devolución | 409 | OK |

En total se aprobaron **22 casos**. La salida completa está en [`evidencias/resultado-http.txt`](evidencias/resultado-http.txt).

## Calidad del código

Se ejecutaron los siguientes comandos sin errores:

```text
gofmt -w ...
go test ./...
go vet ./...
git diff --check
```

También se regeneró Swagger con las rutas de `laptops` y `laptop-loans`.

## Lista final del enunciado

- [x] La solución construye con Docker Compose.
- [x] Los tres contenedores levantan y MySQL queda healthy.
- [x] Health, ping, registro, login y libros funcionan.
- [x] El préstamo anterior de libros sigue funcionando.
- [x] Un administrador puede crear portátiles.
- [x] Un usuario autenticado puede listar portátiles.
- [x] Un usuario puede solicitar, consultar y devolver sus portátiles.
- [x] Un portátil ocupado no se puede prestar nuevamente.
- [x] Se verificaron token ausente, rol incorrecto, serial duplicado, devolución ajena, inexistente y repetida.
- [x] La ampliación usa una migración `002`.
- [x] Se mantienen las capas handler, service y repository.
- [x] Se documentan análisis, ADR y cuatro diagramas Mermaid.
