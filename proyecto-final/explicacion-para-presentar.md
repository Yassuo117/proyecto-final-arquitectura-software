# Explicación para presentar el proyecto

Este documento está escrito como apoyo para hablar. Lo mejor es entender las ideas y explicarlas con palabras propias, no memorizar cada frase.

## Presentación de diez minutos

### 1. Introducción

> Mi proyecto parte de una API de biblioteca que ya permitía registrar usuarios, administrar libros y prestar libros. El cambio solicitado fue agregar el préstamo de portátiles sin dañar lo que ya existía. Antes de programar revisé dónde estaba acoplado el préstamo al concepto de libro y encontré dependencias en la base de datos, los modelos, el repositorio, el servicio y los handlers.

La idea importante aquí es que no se empezó modificando archivos al azar. Primero se entendió el diseño existente y el impacto del cambio.

### 2. El problema arquitectónico

> El préstamo anterior tiene un `book_id` y controla disponibilidad con una cantidad de copias. Un portátil es diferente porque cada equipo físico tiene un serial único y un estado propio. Por ejemplo, dos ejemplares del mismo libro pueden tratarse como copias equivalentes, pero dos portátiles del mismo modelo siguen siendo equipos diferentes.

Esta diferencia explica casi toda la decisión. Los libros usan cantidades; los portátiles usan identidad individual.

### 3. Alternativas que evalué

> Consideré tres alternativas. La primera era agregar un `laptop_id` opcional a la tabla de préstamos existente. La descarté porque la tabla terminaría con campos nulos y muchas condiciones para saber si el préstamo era de libro o de portátil. La segunda era convertir todo a un modelo genérico de recursos prestables. Era más flexible, pero obligaba a migrar lo que ya funcionaba y aumentaba bastante el riesgo. La tercera era crear un módulo paralelo para portátiles, reutilizando autenticación, roles, MySQL y la arquitectura por capas. Elegí esta última porque representa correctamente los equipos y afecta lo mínimo posible el módulo de libros.

La desventaja aceptada es que existe cierta repetición entre los dos módulos. En este proyecto es preferible a crear una abstracción demasiado grande sin conocer todavía las reglas de otros recursos.

### 4. Cómo quedó organizada la solución

> La solución conserva las tres capas principales. El handler se encarga de HTTP, recibe el JSON y devuelve códigos de estado. El service representa los casos de uso y define reglas como los siete días de préstamo. El repository contiene el SQL y las transacciones. Los modelos representan portátiles y préstamos.

Archivos para mostrar:

```text
internal/handler/laptop_handler.go
internal/service/laptop_service.go
internal/repository/laptop_repository.go
internal/repository/laptop_loan_repository.go
internal/models/laptop.go
```

Una forma corta de explicarlo es:

```text
Petición HTTP → Handler → Service → Repository → MySQL
```

- **Handler:** entiende HTTP, JWT, parámetros y códigos de respuesta.
- **Service:** entiende el caso de uso y calcula el vencimiento.
- **Repository:** entiende SQL, bloqueos y transacciones.
- **MySQL:** conserva los datos.

### 5. Modelo de datos

> Agregué una migración nueva llamada `002_laptops.sql`; no mezclé el esquema nuevo dentro de la migración inicial. La tabla `laptops` guarda marca, modelo, serial único y estado. La tabla `laptop_loans` relaciona el usuario con el portátil y guarda fechas y estado del préstamo.

Estados del portátil:

- `available`: se puede prestar.
- `loaned`: tiene un préstamo activo.
- `maintenance`: está registrado, pero no se puede prestar.

Estados del préstamo:

- `active`: todavía no se ha devuelto.
- `returned`: ya terminó.

El serial tiene una restricción `UNIQUE`. Por eso no se pueden registrar dos equipos con el mismo identificador.

### 6. Rutas implementadas

| Método y ruta | Función | Permiso |
|---|---|---|
| `GET /api/laptops` | Listar portátiles | Usuario autenticado |
| `POST /api/laptops` | Registrar portátil | Solo administrador |
| `POST /api/laptop-loans` | Solicitar portátil | Usuario autenticado |
| `GET /api/laptop-loans/me` | Ver préstamos propios | Usuario autenticado |
| `PUT /api/laptop-loans/:id/return` | Devolver portátil | Dueño del préstamo |

> No se creó otro sistema de usuarios. Estas rutas usan el mismo JWT y los mismos roles que ya tenía la API.

### 7. Flujo de un préstamo

> Cuando llega la solicitud, el middleware valida el token y deja el `user_id` en el contexto. El handler valida el JSON. El servicio calcula una fecha de vencimiento de siete días y llama al repositorio. El repositorio inicia una transacción y bloquea la fila del portátil con `SELECT ... FOR UPDATE`. Si está disponible, crea el préstamo y cambia el equipo a `loaned`. Finalmente confirma la transacción y responde con HTTP 201.

La transacción es importante porque crear el préstamo y cambiar el estado deben ocurrir juntos. Si una operación falla, ninguna queda aplicada a medias.

El bloqueo `FOR UPDATE` resuelve este caso:

1. Dos usuarios solicitan el mismo portátil casi al mismo tiempo.
2. La primera petición bloquea la fila.
3. La segunda debe esperar.
4. La primera cambia el estado a `loaned` y confirma.
5. La segunda continúa, encuentra que ya no está disponible y recibe `409 Conflict`.

### 8. Devolución y permisos

> La devolución también usa una transacción. Primero se busca y bloquea el préstamo. Después se comprueba que exista, que pertenezca al usuario autenticado y que no haya sido devuelto antes. Si todo es correcto, el préstamo cambia a `returned` y el portátil vuelve a `available`.

Códigos útiles para explicar:

- `200`: operación correcta.
- `201`: recurso o préstamo creado.
- `400`: datos enviados incorrectamente.
- `401`: falta token o es inválido.
- `403`: el usuario está autenticado, pero no tiene permiso.
- `404`: portátil o préstamo inexistente.
- `409`: conflicto con el estado actual, como serial repetido, equipo ocupado o devolución repetida.

### 9. Diagramas

> Preparé cuatro diagramas. El de contexto muestra a estudiantes y administradores usando la Biblioteca API. El de contenedores muestra la API, MySQL y Adminer coordinados por Docker Compose. El de componentes baja al nivel de handlers, services y repositories. Finalmente, el diagrama de secuencia muestra todas las validaciones del préstamo de un portátil.

Al presentar, no es necesario leer cada caja. Se debe explicar qué nivel de detalle responde cada diagrama.

### 10. Verificación

> No me limité a comprobar que compilara. Reconstruí los contenedores, confirmé que MySQL quedara healthy y probé nuevamente el registro, login, libros y préstamo de libros. Después ejecuté los casos de portátiles. En total se aprobaron 22 casos HTTP, incluidos token ausente, estudiante intentando crear un equipo, serial duplicado, equipo ya ocupado, devolución de otro usuario, ID inexistente y devolución repetida. También probé las migraciones en una base completamente vacía.

Comandos de calidad ejecutados:

```text
go test ./...
go vet ./...
git diff --check
docker compose up -d --build
```

### 11. Cierre

> En conclusión, agregué el préstamo de portátiles sin cambiar el funcionamiento de los libros. La decisión principal fue separar los dos módulos porque manejan la disponibilidad de forma diferente. El resultado mantiene las capas existentes, reutiliza la seguridad y usa transacciones para conservar la consistencia. La solución no intenta resolver todos los posibles recursos futuros, pero deja claro cuándo tendría sentido extraer un modelo más genérico.

## Demostración corta en vivo

Si el profesor pide mostrar el sistema, este orden toma pocos minutos:

1. Ejecutar `docker compose ps` y señalar los tres contenedores y `healthy` en MySQL.
2. Abrir `http://localhost:8080/swagger/index.html`.
3. Mostrar `GET /health` o `GET /ping`.
4. Mostrar las rutas de `laptops` y `laptop-loans` en Swagger.
5. Abrir `migrations/002_laptops.sql` y señalar las dos tablas.
6. Abrir `laptop_loan_repository.go` y señalar `Beginx`, `FOR UPDATE`, `Commit` y `Rollback`.
7. Mostrar `proyecto-final/verificacion.md` y la evidencia de 22 casos.

No conviene improvisar cambios en la base durante la exposición. Para una demostración completa se puede usar `pruebas-proyecto-final.http` siguiendo el orden de sus peticiones.

## Preguntas difíciles y respuestas claras

### ¿Por qué no reutilizó la tabla `loans`?

Porque esa tabla depende directamente de `book_id` y el servicio modifica cantidades de copias. Agregar otro tipo habría mezclado dos reglas distintas y aumentado la posibilidad de registros inválidos.

### ¿Su solución viola DRY por repetir estructura?

Existe repetición estructural, pero todavía no hay suficiente información para diseñar una abstracción estable. La separación es intencional porque los modelos tienen reglas distintas. Si aparece un tercer activo individual, se puede refactorizar con evidencia real.

### ¿Por qué siete días?

Fue una decisión de negocio tomada y documentada en el ADR. La regla está centralizada en el servicio, por lo que se puede cambiar sin modificar las rutas ni el acceso a datos.

### ¿Qué pasa si la aplicación falla después de crear el préstamo?

La creación del préstamo y el cambio de estado están dentro de una transacción. Si ocurre un error antes del `Commit`, el `Rollback` evita que quede un cambio parcial.

### ¿Dos usuarios pueden prestar el mismo portátil al mismo tiempo?

No mediante la API. El repositorio bloquea la fila del equipo mientras decide y actualiza su estado. La segunda transacción ve el estado `loaned` y responde con conflicto.

### ¿Por qué `maintenance` si no se implementó una ruta para cambiarlo?

Es un estado válido del dominio y evita representar un equipo dañado como disponible. La administración de mantenimiento puede ser una evolución posterior; no era parte de los endpoints mínimos solicitados.

### ¿Cómo comprobó que no rompió los libros?

La batería final registra un usuario, lista libros, crea un préstamo de libro y lo devuelve. Esos casos se ejecutaron después de construir la imagen con el módulo nuevo.

### ¿Qué mejoraría con más tiempo?

Agregaría pruebas unitarias con interfaces para simular repositorios, una ruta administrativa para mantenimiento y una restricción adicional en base de datos para reforzar que solo exista un préstamo activo por portátil. También evaluaría un modelo genérico únicamente si aparecen más recursos con reglas parecidas.

## Lo mínimo que debes recordar

Si te pones nervioso, recuerda estas cinco ideas:

1. Libro por cantidad; portátil por serial.
2. Se compararon tres opciones y se eligió la de menor riesgo.
3. Se mantuvieron handler, service y repository.
4. Las transacciones y el bloqueo evitan estados parciales y préstamos dobles.
5. Se aprobaron 22 casos y el flujo anterior de libros sigue funcionando.
