# Proyecto final de Arquitectura de Software

**Estudiante:** Carlos Antonio Ramírez Galvis  
**Modalidad:** individual  
**Tema:** ampliación de la Biblioteca API para préstamos de portátiles

## Contenido de la entrega

Los documentos solicitados se encuentran en esta misma carpeta:

1. [`analisis.md`](analisis.md): revisión del sistema existente, acoplamientos, alternativas y comparación.
2. [`adr.md`](adr.md): decisión arquitectónica aceptada, reglas de negocio y consecuencias.
3. [`diagramas.md`](diagramas.md): C4 de contexto, C4 de contenedores, componentes y secuencia de préstamo.
4. [`verificacion.md`](verificacion.md): ambiente usado, resultados y checklist de pruebas.

Material de apoyo:

- [`evidencias/resultado-http.txt`](evidencias/resultado-http.txt): salida de los 22 casos ejecutados.
- [`evidencias/resultado-migraciones.txt`](evidencias/resultado-migraciones.txt): comprobación desde una base vacía.
- [`guion-defensa.md`](guion-defensa.md): resumen para preparar la sustentación individual.
- [`explicacion-para-presentar.md`](explicacion-para-presentar.md): explicación oral completa, demostración y preguntas posibles.
- [`checklist-entrega.md`](checklist-entrega.md): última revisión antes de compartir el enlace.

## Ubicación de la implementación

El código se conserva en la estructura original del repositorio:

```text
migrations/002_laptops.sql
internal/models/laptop.go
internal/repository/laptop_repository.go
internal/repository/laptop_loan_repository.go
internal/service/laptop_service.go
internal/handler/laptop_handler.go
internal/handler/routes.go
cmd/api/main.go
```

Las peticiones manuales están en `pruebas-proyecto-final.http` y la verificación completa se puede repetir con `scripts/verificar-proyecto.ps1`.

## Decisión en una frase

Se agregó un módulo paralelo para portátiles porque cada equipo tiene serial y estado propio, mientras los libros se administran por cantidad de copias. Con esto se evita modificar o migrar el préstamo de libros que ya funcionaba.
