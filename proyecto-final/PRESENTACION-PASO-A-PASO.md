# Presentación paso a paso (10 minutos)

Este es el orden exacto. Durante la exposición avanza de izquierda a derecha por las pestañas que quedaron abiertas en Cursor.

## Antes de hablar

Comprueba que Docker Desktop diga que el motor está funcionando y que Swagger abra en:

`http://localhost:8080/swagger/index.html`

En Cursor deben quedar abiertas estas pestañas, en este orden:

1. `diagramas.md` en vista previa.
2. `adr.md` en vista previa.
3. `002_laptops.sql`.
4. `routes.go`.
5. `laptop_loan_repository.go`.
6. `verificacion.md` en vista previa.
7. `PRESENTACION-PASO-A-PASO.md` en vista previa, para estudiar antes de comenzar.

## Paso 1 — Presentarme y explicar el problema (0:00 a 1:00)

**Qué mostrar:** la primera parte de `diagramas.md`.

**Qué decir:**

> Buenos días. Mi proyecto consiste en ampliar una API de biblioteca que ya manejaba usuarios, libros y préstamos. Lo que agregué fue el préstamo de portátiles. El problema es que un libro se controla por cantidad de copias, pero un portátil es un equipo individual: tiene un serial y puede estar disponible, prestado o en mantenimiento. Por eso no era buena idea tratarlo exactamente igual que un libro.

## Paso 2 — Explicar las opciones que consideré (1:00 a 2:30)

**Qué hacer:** cambia a la pestaña `adr.md`.

**Qué decir:**

> Antes de programar comparé tres opciones. La primera era meter un campo de portátil dentro de la tabla de préstamos de libros, pero eso mezclaba dos cosas distintas y dejaba campos vacíos. La segunda era cambiar todo a un sistema genérico de recursos prestables. Esa opción podía servir a futuro, pero obligaba a modificar bastante código que ya funcionaba. La tercera era hacer un módulo separado para portátiles y reutilizar la autenticación, la base de datos y la arquitectura que ya existían. Escogí esa opción porque era la más clara y la que tenía menos riesgo de dañar el módulo de libros.

## Paso 3 — Mostrar cómo quedó la arquitectura (2:30 a 3:45)

**Qué hacer:** vuelve a `diagramas.md` y muestra el diagrama de componentes.

**Qué decir:**

> Mantuve la arquitectura por capas del proyecto. El handler recibe la petición y responde por HTTP. El service contiene las reglas del negocio. El repository se comunica con MySQL. También reutilicé el JWT y los roles. Esto permite agregar la nueva función sin mezclar la lógica en un solo archivo.

Si te preguntan por Docker:

> Docker Compose levanta tres servicios: la API, MySQL y Adminer. Así el proyecto se puede ejecutar de la misma forma en otro computador.

## Paso 4 — Mostrar la base de datos (3:45 a 4:45)

**Qué hacer:** abre la pestaña `002_laptops.sql`.

**Qué señalar:** las tablas `laptops` y `laptop_loans`.

**Qué decir:**

> En lugar de cambiar la migración inicial, agregué una segunda migración. La tabla laptops guarda cada equipo con marca, modelo, serial único y estado. La tabla laptop_loans guarda quién pidió el equipo, las fechas y el estado del préstamo. La relación se hace con llaves foráneas para mantener la integridad de los datos.

## Paso 5 — Mostrar las rutas (4:45 a 5:45)

**Qué hacer:** cambia a `routes.go`.

**Qué señalar:** `/laptops`, `/laptop-loans` y el middleware de administrador.

**Qué decir:**

> Estas son las rutas nuevas. Un usuario autenticado puede ver los equipos, solicitar uno, consultar sus préstamos y devolver el suyo. Para registrar un portátil se necesita el rol de administrador. Si no hay token, la API responde 401; si el rol o el propietario no es correcto, responde 403.

## Paso 6 — Explicar cómo se evita prestar dos veces el mismo equipo (5:45 a 7:00)

**Qué hacer:** cambia a `laptop_loan_repository.go`.

**Qué señalar:** `Beginx`, `FOR UPDATE`, `Commit` y `Rollback`.

**Qué decir:**

> Esta es la parte más importante para la consistencia. Cuando alguien pide un portátil, se inicia una transacción y se bloquea la fila del equipo con FOR UPDATE. Después se revisa que siga disponible, se crea el préstamo y se cambia el estado a loaned. Si algo falla se hace rollback y no queda la información a medias. Esto también evita que dos usuarios puedan llevarse el mismo portátil al mismo tiempo.

No expliques cada línea. Solo señala esas cuatro partes.

## Paso 7 — Hacer la demostración en Swagger (7:00 a 8:30)

**Qué hacer:** cambia al navegador y muestra Swagger.

**Qué decir mientras bajas por la página:**

> En Swagger quedaron documentadas las rutas. Aquí está el inicio de sesión, el catálogo de portátiles, la creación para administradores, la solicitud, la consulta de préstamos propios y la devolución. Para la demostración no necesito crear todo de nuevo, porque la batería automática ya ejecuta el flujo completo; aquí muestro que las operaciones están disponibles y documentadas.

Muestra, en este orden:

1. `POST /auth/login`.
2. `GET /laptops`.
3. `POST /laptops`.
4. `POST /laptop-loans`.
5. `GET /laptop-loans/me`.
6. `PUT /laptop-loans/{id}/return`.

No presiones `Try it out` a menos que el profesor te pida ejecutar una petición.

## Paso 8 — Mostrar las pruebas (8:30 a 9:30)

**Qué hacer:** vuelve a Cursor y abre `verificacion.md`.

**Qué decir:**

> Para comprobar que la ampliación no dañó el sistema ejecuté las pruebas de Go, la revisión con go vet, la construcción de Docker y 22 casos HTTP. Los casos incluyen el préstamo y la devolución anterior de libros, además del flujo de portátiles. También probé errores como token ausente, rol incorrecto, serial repetido, equipo ocupado, devolución de otro usuario, préstamo inexistente y devolución repetida. Los 22 casos aprobaron.

## Paso 9 — Cerrar la presentación (9:30 a 10:00)

**Qué mostrar:** vuelve al diagrama de componentes.

**Qué decir:**

> En conclusión, agregué los préstamos de portátiles sin reemplazar lo que ya funcionaba. Separé las reglas de libros y portátiles, mantuve la arquitectura por capas y utilicé transacciones para proteger la consistencia. La solución tiene algo de repetición entre módulos, pero reduce el riesgo y deja claro cómo funciona cada tipo de préstamo. Gracias.

## Si el profesor hace preguntas

**¿Por qué no hiciste todo genérico?**

> Porque para este alcance solo había dos recursos con reglas distintas. Hacerlo genérico obligaba a migrar el módulo anterior y aumentaba el riesgo. Si después aparecen más activos parecidos, ahí sí tendría sentido extraer una solución común.

**¿Cómo evitas el préstamo doble?**

> Uso una transacción y bloqueo la fila del portátil con SELECT FOR UPDATE. La segunda petición encuentra el equipo como prestado y recibe un conflicto 409.

**¿Por qué separaste las tablas?**

> Porque los libros se administran por cantidad, mientras cada portátil tiene identidad e historial propios.

**¿Qué significa rollback?**

> Significa deshacer los cambios de la operación si uno de los pasos falla, para no dejar la base de datos a medias.

**¿Qué significa JWT?**

> Es el token que identifica al usuario después de iniciar sesión. La API lo revisa antes de permitir las rutas protegidas.

**¿Qué mejorarías después?**

> Agregaría pruebas unitarias más específicas, manejo del estado de mantenimiento desde la API y una interfaz web para los usuarios.

## Frase para recordar si te quedas en blanco

> Separé los portátiles porque cada equipo tiene serial y estado propio; mantuve las capas existentes y protegí los préstamos con transacciones.

