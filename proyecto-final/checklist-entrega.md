# Checklist antes de entregar

## Archivos obligatorios

- [x] `proyecto-final/analisis.md`
- [x] `proyecto-final/adr.md`
- [x] `proyecto-final/diagramas.md`
- [x] `proyecto-final/verificacion.md`
- [x] Migración nueva `migrations/002_laptops.sql`
- [x] Implementación separada en handler, service y repository

## Funcionalidad

- [x] Administrador registra un portátil.
- [x] Usuario autenticado lista portátiles.
- [x] Usuario solicita un portátil disponible.
- [x] Usuario consulta sus préstamos.
- [x] Usuario devuelve su propio portátil.
- [x] Un equipo ocupado no se presta otra vez.
- [x] Se controlan serial repetido, token ausente, rol incorrecto y devoluciones inválidas.
- [x] Registro, login, libros y préstamos de libros siguen funcionando.

## Calidad y evidencia

- [x] `go test ./...`
- [x] `go vet ./...`
- [x] `git diff --check`
- [x] Docker Compose construye la API.
- [x] API, MySQL y Adminer levantan.
- [x] MySQL aparece como healthy.
- [x] 22 casos HTTP aprobados.
- [x] Migraciones verificadas desde una base vacía.
- [x] Swagger actualizado.
- [x] Historial dividido por hitos.
- [x] `.env` y respaldo local excluidos de Git.

## Pendiente únicamente al publicar

- [ ] Crear o confirmar el fork en la cuenta del estudiante.
- [ ] Subir la rama del proyecto al fork.
- [ ] Revisar en GitHub que los diagramas Mermaid se rendericen.
- [ ] Copiar el enlace del repositorio para el profesor.

No se debe subir `.env`, contraseñas, tokens ni el respaldo local de MySQL.
