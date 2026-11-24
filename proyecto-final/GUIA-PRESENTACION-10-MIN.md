# Guía para presentar el proyecto en 10 minutos

## Antes de comenzar

Deja abiertas estas pestañas, en este orden:

1. `proyecto-final/diagramas.md` en vista previa de Markdown.
2. `migrations/002_laptops.sql`.
3. `internal/handler/routes.go`.
4. `internal/repository/laptop_loan_repository.go`.
5. Swagger en `http://localhost:8080/swagger/index.html`.
6. Una terminal en la raíz del proyecto para mostrar `docker compose ps` y, si te lo piden, ejecutar `scripts/verificar-proyecto.ps1`.

No intentes explicar todos los archivos. La idea es demostrar el problema, la decisión arquitectónica, la implementación y las pruebas.

## Guion cronometrado

### 0:00–1:00 — Presentación del problema

> Este proyecto parte de una API de biblioteca que ya manejaba usuarios, libros y préstamos. El nuevo requisito fue permitir el préstamo de portátiles sin dañar lo que ya funcionaba. El reto principal era que un libro se controla por cantidad de copias, mientras que cada portátil es un equipo físico individual con serial y estado propio.

Muestra el primer diagrama de `diagramas.md`.

### 1:00–2:30 — Alternativas analizadas

> Analicé tres opciones. La primera era agregar un `laptop_id` opcional a la tabla de préstamos de libros, pero eso mezclaba dos reglas distintas y producía muchos campos nulos. La segunda era convertir todo a un modelo genérico de recursos prestables; era flexible, pero implicaba migrar y arriesgar el módulo existente. La tercera era crear un módulo paralelo para portátiles y reutilizar la autenticación, los roles, MySQL y la arquitectura por capas. Elegí la tercera por su menor impacto y porque representa mejor el dominio.

Si te preguntan dónde está documentado, abre `proyecto-final/adr.md`.

### 2:30–4:00 — Arquitectura de la solución

Muestra los diagramas de contenedores y componentes.

> La API conserva la separación handler, service y repository. El handler recibe HTTP y valida la entrada; el service concentra las reglas de negocio; el repository accede a MySQL. Las rutas protegidas pasan por JWT y algunas además exigen el rol de administrador. Docker Compose levanta la API, MySQL y Adminer.

Abre `internal/handler/routes.go` y señala las rutas de `laptops` y `laptop-loans`.

### 4:00–5:00 — Modelo de datos

Abre `migrations/002_laptops.sql`.

> La migración 002 agrega dos tablas sin modificar el esquema original. `laptops` guarda marca, modelo, serial único y estado. `laptop_loans` relaciona usuario y portátil, y conserva las fechas y el estado del préstamo. Los estados del equipo son `available`, `loaned` y `maintenance`.

### 5:00–6:30 — Regla crítica de consistencia

Abre `internal/repository/laptop_loan_repository.go` y señala `Beginx`, `FOR UPDATE`, `Commit` y `Rollback`.

> La operación de préstamo es transaccional. Primero bloquea la fila del portátil con `SELECT FOR UPDATE`, comprueba que esté disponible, crea el préstamo y cambia el equipo a `loaned`. Si falla cualquier paso, hace rollback. Así se evita que dos peticiones simultáneas presten el mismo equipo. La devolución usa la misma idea: actualiza el préstamo y libera el portátil en una sola transacción.

### 6:30–8:00 — Demostración funcional

Muestra Swagger y recorre estas rutas, sin necesidad de ejecutar todas:

| Ruta | Qué demuestra |
|---|---|
| `POST /api/auth/login` | Autenticación y obtención de JWT |
| `GET /api/laptops` | Catálogo protegido |
| `POST /api/laptops` | Creación restringida a administrador |
| `POST /api/laptop-loans` | Préstamo de un equipo disponible |
| `GET /api/laptop-loans/me` | Consulta de los préstamos propios |
| `PUT /api/laptop-loans/{id}/return` | Devolución del equipo propio |

> La API devuelve códigos HTTP coherentes: 201 al crear, 401 sin token, 403 por rol o propietario incorrecto, 404 si no existe y 409 para conflictos como serial repetido, equipo ocupado o devolución repetida.

### 8:00–9:15 — Verificación

En la terminal muestra:

```powershell
docker compose ps
```

Luego explica:

> La solución se verificó con `go test`, `go vet`, construcción de la imagen y 22 casos HTTP automatizados. Las pruebas incluyen el flujo anterior de libros y los casos nuevos de portátiles, además de errores de autenticación, autorización, duplicados y conflictos. También se probaron las migraciones desde una base vacía.

Si el profesor pide evidencia, abre `proyecto-final/verificacion.md` o ejecuta:

```powershell
.\scripts\verificar-proyecto.ps1
```

El resultado esperado termina con `VERIFICACIÓN COMPLETA: 22 casos aprobados`.

### 9:15–10:00 — Cierre

> En conclusión, la decisión priorizó extender el sistema sin romperlo. El módulo paralelo conserva las reglas propias de cada recurso, reutiliza la infraestructura existente y protege la consistencia con transacciones. La desventaja es cierta repetición estructural, pero se aceptó porque reduce el riesgo y mantiene claro el dominio.

## Respuestas cortas a preguntas probables

**¿Por qué no hizo una tabla genérica?**  
Porque solo había dos recursos con reglas diferentes. Generalizar ahora aumentaba el alcance y obligaba a migrar lo que ya funcionaba.

**¿Cómo evita el préstamo doble?**  
Con una transacción y un bloqueo `SELECT ... FOR UPDATE` sobre la fila del portátil.

**¿Quién puede registrar equipos?**  
Solo un usuario autenticado con rol `admin`.

**¿Puede un usuario devolver el préstamo de otro?**  
No. Se valida el propietario y se responde 403.

**¿Qué quedó reutilizado?**  
JWT, roles, conexión MySQL, Docker Compose, manejo HTTP y la estructura por capas.

**¿Qué desventaja tiene la decisión?**  
Existe repetición entre el módulo de préstamos de libros y el de portátiles, pero se mantiene el dominio separado y se evita una migración riesgosa.

## Plan de emergencia si algo no abre

- Si Swagger no carga, muestra `docs/swagger.yaml` y `pruebas-proyecto-final.http`.
- Si Docker tarda, explica con `proyecto-final/evidencias/resultado-http.txt` y `resultado-migraciones.txt`.
- Si no recuerdas una respuesta, vuelve a la idea central: *portátil individual por serial, módulo paralelo, arquitectura por capas y transacciones para consistencia*.

