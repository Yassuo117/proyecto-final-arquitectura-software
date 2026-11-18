# Análisis del cambio solicitado

**Estudiante:** Carlos Antonio Ramírez Galvis  
**Modalidad:** individual

## 1. ¿Cómo funciona actualmente el préstamo de libros?

La API está organizada por capas. El `LoanHandler` recibe la petición HTTP y obtiene el usuario desde el JWT; `LoanService` aplica las reglas del préstamo; y `LoanRepository` guarda o consulta los datos en MySQL. Además, el servicio de préstamos depende de `BookRepository` para verificar y modificar la cantidad de copias disponibles.

El flujo actual es:

1. El estudiante envía un `book_id` a `POST /api/loans`.
2. El servicio busca el libro y revisa que `available_copies` sea mayor que cero.
3. Se valida que ese usuario no tenga un préstamo activo del mismo libro.
4. Se crea el préstamo por 14 días.
5. Se descuenta una copia disponible.
6. Al devolverlo, se marca el préstamo como devuelto y se suma la copia.

## 2. Puntos acoplados al concepto de libro

| Lugar | Acoplamiento encontrado |
|---|---|
| `migrations/001_init.sql` | La tabla `loans` exige un `book_id` y tiene una llave foránea directa a `books`. |
| `internal/models/loan.go` | `Loan` y `LoanRequest` contienen `BookID`; el detalle contiene título y autor. |
| `internal/repository/loan_repository.go` | Las consultas hacen `JOIN books`, crean préstamos con `book_id` y buscan duplicados por usuario y libro. |
| `internal/service/loan_service.go` | La disponibilidad se calcula con copias del libro y se modifican esas copias al prestar/devolver. |
| `internal/handler/loan_handler.go` | Los mensajes y la documentación HTTP hablan específicamente de libros. |
| `internal/handler/routes.go` | Las rutas `/loans` están conectadas al handler anterior. |

Por lo anterior, cambiar simplemente `book_id` por algo genérico afectaría una funcionalidad que ya existe y obligaría a migrar datos. También mezclaría dos formas diferentes de inventario: los libros se controlan por cantidad de copias, pero cada portátil debe identificarse por su serial.

## 3. Pregunta de evolución futura

Si más adelante la biblioteca presta calculadoras, tabletas o cámaras, la solución debería permitir agregarlas sin duplicar toda la aplicación. Sin embargo, tampoco conviene generalizar antes de conocer las reglas de esos recursos. Por ejemplo, una calculadora podría manejarse por cantidad, mientras una cámara podría necesitar serial, accesorios y estado de mantenimiento.

La decisión de este proyecto busca un equilibrio: no romper lo que funciona, mantener las mismas capas y dejar un límite claro para el nuevo dominio. Si aparece un tercer recurso con reglas parecidas a las de los portátiles, ahí sí existiría evidencia suficiente para extraer un módulo genérico de activos prestables.

## 4. Opciones consideradas

### Opción A: agregar `laptop_id` a la tabla `loans`

Se podría permitir que un préstamo tenga `book_id` o `laptop_id`. Reutiliza parte del código, pero introduce columnas nulas y obliga a preguntar el tipo de recurso en muchas consultas. También es fácil terminar con registros inválidos que tengan los dos campos o ninguno.

### Opción B: crear un modelo genérico de recursos prestables

Consiste en crear entidades como `resources`, `resource_items` y `loans`, y migrar libros y portátiles a ese modelo. Es la opción más flexible a largo plazo, pero implica rehacer una parte grande del sistema, migrar el préstamo actual y asumir reglas de recursos que todavía no se conocen.

### Opción C: módulo paralelo de portátiles

Se crean `laptops` y `laptop_loans`, con sus modelos, repositorios, servicio y handler. La autenticación, el middleware y la infraestructura se reutilizan. La diferencia de negocio queda explícita: cada portátil tiene serial y estado propio. Esta opción agrega algo de código, pero reduce el riesgo sobre los libros y cabe en el alcance del proyecto.

## 5. Comparación

| Criterio | A. Una tabla con dos FK | B. Modelo genérico | C. Módulo paralelo |
|---|---:|---:|---:|
| Riesgo de romper libros | Medio | Alto | Bajo |
| Complejidad de implementación | Media | Alta | Media |
| Claridad del modelo | Baja | Alta | Alta |
| Extensibilidad futura | Media | Alta | Media |
| Migración de datos existentes | Sí | Sí | No |
| Ajuste al tiempo y alcance | Medio | Bajo | Alto |

La opción seleccionada es la **C**. Su principal costo es cierta repetición, pero en este caso la repetición refleja reglas diferentes y evita una abstracción prematura.
