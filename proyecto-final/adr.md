# ADR-001: separar los préstamos de portátiles de los préstamos de libros

- **Estado:** aceptada
- **Fecha:** 18 de noviembre de 2026
- **Autor:** Carlos Antonio Ramírez Galvis

## Contexto

La Biblioteca API ya permite prestar libros. Ese diseño usa una cantidad de copias (`available_copies`) y una tabla `loans` que referencia directamente a `books`. El nuevo requisito pide registrar portátiles, prestarlos, devolverlos y consultar los préstamos del usuario sin dañar ninguna función existente.

Un portátil no es equivalente a una copia anónima de un libro. Cada equipo físico tiene un serial único, puede estar disponible, prestado o en mantenimiento, y su historial debe apuntar siempre al mismo equipo. Además, el proyecto debe conservar la arquitectura por capas, la autenticación JWT y los roles actuales.

## Fuerzas que influyen en la decisión

- Reducir el riesgo de regresión en los préstamos de libros.
- Representar correctamente un activo individual con serial.
- Mantener handlers, servicios y repositorios con responsabilidades claras.
- Evitar una migración grande de datos en un proyecto académico de alcance corto.
- Dejar una evolución razonable si aparecen otros tipos de recursos.
- Controlar solicitudes concurrentes para que un equipo no se preste dos veces.

## Opciones evaluadas

### 1. Extender `loans` con `laptop_id`

La tabla tendría dos relaciones opcionales: libro o portátil. Se reutilizarían algunas rutas y structs.

**A favor:** menos tablas nuevas y apariencia inicial de reutilización.  
**En contra:** aparecen nulos, validaciones por tipo y consultas condicionales; el modelo puede quedar en estados inválidos. La lógica de copias de libros se mezclaría con la de serial y estado de portátiles.

### 2. Reemplazar el diseño por un recurso prestable genérico

Libros y portátiles pasarían a un modelo común, posiblemente con `resources`, `items` y un único `loans`.

**A favor:** ofrece una base más flexible para muchos recursos.  
**En contra:** exige migrar el modelo ya entregado, aumenta el alcance y obliga a crear abstracciones sin conocer todavía las reglas de los próximos recursos. El riesgo de romper funcionalidades existentes es alto.

### 3. Crear un módulo paralelo para portátiles

Se agregan las tablas `laptops` y `laptop_loans` y las mismas capas utilizadas por el sistema actual. Se comparten JWT, roles, conexión a MySQL y convenciones HTTP.

**A favor:** el cambio queda aislado, el modelo expresa el serial y el estado del equipo, y no se migran préstamos existentes.  
**En contra:** hay clases y consultas parecidas a las del módulo de libros. Si crecen muchos tipos de activos, esta repetición deberá revisarse.

## Decisión

Se adopta la **opción 3: módulo paralelo para portátiles**.

Cada fila de `laptops` representa un equipo físico y su `serial_number` es único. Se manejan tres estados: `available`, `loaned` y `maintenance`. Los préstamos se guardan en `laptop_loans`, separados de `loans`.

Las reglas elegidas son:

- El catálogo de portátiles puede ser consultado por cualquier usuario autenticado.
- Solo un administrador puede registrar un portátil.
- Estudiantes y administradores autenticados pueden solicitar y devolver sus propios préstamos.
- La duración del préstamo es de **7 días**.
- Un equipo debe estar `available` para prestarse.
- Una devolución solo puede realizarla el usuario dueño del préstamo.
- Un préstamo devuelto no puede devolverse otra vez.

El préstamo y la devolución usan transacciones de MySQL. Durante el préstamo se bloquea la fila del portátil con `SELECT ... FOR UPDATE`, se crea el préstamo y se cambia el estado en una sola operación atómica. Así, dos peticiones simultáneas no pueden obtener el mismo equipo. La devolución también actualiza el préstamo y el equipo en una transacción.

El nuevo esquema se agrega en `migrations/002_laptops.sql`. La migración `001` conserva el esquema original; únicamente se corrigió el separador inválido que ya existía entre los dos últimos registros del seed y que impedía iniciar una base vacía.

## Consecuencias

### Positivas

- Los flujos de libros no necesitan cambiar.
- El modelo refleja que cada portátil es individual.
- Las reglas nuevas quedan concentradas en su propio servicio.
- Las transacciones evitan estados parciales y préstamos dobles por concurrencia.
- Los códigos HTTP distinguen errores de validación, no encontrado, conflicto y falta de permisos.

### Negativas y deuda aceptada

- Existen dos módulos de préstamo y cierta repetición estructural.
- Los estados están definidos como `ENUM`; agregar un estado requiere una migración.
- No se implementó todavía un módulo genérico para cualquier activo.

### Evolución prevista

Si se solicita un tercer recurso con identidad individual y reglas similares, se compararán los tres módulos y se podrá extraer un concepto común de `asset` o `loanable_item`. Esa refactorización deberá incluir una estrategia de migración y no se hará solo para eliminar unas pocas líneas repetidas.
