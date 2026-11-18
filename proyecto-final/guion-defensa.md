# Guion corto para la defensa

Este archivo no reemplaza los documentos de entrega. Sirve para preparar la explicación de aproximadamente diez minutos.

## Orden sugerido

### 1. Problema (1 minuto)

La API solo prestaba libros y todo el módulo estaba acoplado a `book_id`, copias disponibles, título y autor. El nuevo requisito era prestar portátiles sin romper ese flujo.

### 2. Diferencia de dominio (1 minuto)

Los libros se manejan por cantidades; un portátil es un activo físico individual con serial y estado. Por eso no era correcto tratarlo como otra copia de libro.

### 3. Alternativas (2 minutos)

Se compararon tres opciones: poner `laptop_id` en `loans`, rehacer todo con un recurso genérico o crear un módulo paralelo. La primera mezcla reglas y permite datos inválidos. La segunda es flexible, pero obliga a migrar lo existente y aumenta mucho el riesgo. Se eligió el módulo paralelo porque tiene menor impacto y representa mejor el equipo individual.

### 4. Decisión e implementación (2 minutos)

La solución agrega `laptops` y `laptop_loans` mediante la migración `002`. Mantiene las capas handler, service y repository; reutiliza JWT, roles, MySQL y Docker Compose. El administrador registra equipos, y los usuarios autenticados pueden listar, solicitar y devolver los propios.

### 5. Consistencia (1 minuto)

Al prestar se inicia una transacción y se bloquea el portátil con `SELECT ... FOR UPDATE`. Luego se crea el préstamo y se cambia el estado a `loaned`. Si algo falla, se revierte todo. La devolución también es transaccional. Esto evita el préstamo doble por peticiones concurrentes.

### 6. Diagramas (1 minuto)

Mostrar en orden: contexto, contenedores, componentes y secuencia. En el diagrama de secuencia se ven las validaciones del JWT, JSON, existencia y disponibilidad antes del `201 Created`.

### 7. Evidencia (2 minutos)

Se probaron las migraciones desde una base vacía, los tres contenedores, el flujo anterior de libros y los casos nuevos. La batería automatizada aprobó 22 casos, incluidos sin token, rol incorrecto, serial duplicado, equipo ocupado, devolución ajena, inexistente y repetida.

## Preguntas que pueden hacer

**¿Por qué no hizo un modelo totalmente genérico?**  
Porque solo existen dos recursos con reglas distintas. Generalizar ahora aumentaba el alcance y el riesgo. Si aparece otro activo individual con reglas parecidas, ya habría evidencia para extraer esa abstracción.

**¿Cuál es la principal desventaja de su decisión?**  
Hay repetición estructural entre préstamos de libros y portátiles. Se aceptó a cambio de mantener separados los modelos y no migrar la funcionalidad existente.

**¿Cómo evita prestar dos veces el mismo portátil?**  
La fila del portátil se bloquea dentro de una transacción. Solo la primera petición encuentra el estado `available`; la siguiente espera y luego encuentra `loaned`, por lo que recibe 409.

**¿Por qué el préstamo dura siete días?**  
Es una regla explícita elegida para el proyecto y documentada en el ADR. Está centralizada en el servicio, de modo que se puede cambiar sin modificar el handler ni el repositorio.

**¿Por qué existen los estados `available`, `loaned` y `maintenance`?**  
`available` y `loaned` resuelven el flujo solicitado. `maintenance` permite sacar un equipo del catálogo de préstamos sin borrarlo ni inventar un préstamo falso.

**¿Qué prueba que los libros no se dañaron?**  
La verificación vuelve a listar, prestar y devolver un libro después de construir la nueva imagen. Esos casos aprobaron junto con los del módulo nuevo.
