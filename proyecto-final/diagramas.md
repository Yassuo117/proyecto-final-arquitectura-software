# Diagramas de arquitectura

**Autor:** Carlos Antonio Ramírez Galvis

## C4 nivel 1: contexto del sistema

```mermaid
flowchart LR
    student[Estudiante] -->|Consulta catálogo, solicita y devuelve recursos| system[Biblioteca API]
    admin[Administrador] -->|Administra libros y portátiles| system
    system -->|Lee y guarda usuarios, catálogos y préstamos| db[(MySQL)]
```

La API es el sistema central. Los dos tipos de usuario acceden mediante HTTP y la aplicación guarda su información en MySQL.

## C4 nivel 2: contenedores

```mermaid
flowchart LR
    user[Usuario] -->|HTTP/JSON + JWT| api[API Go / Gin<br/>Puerto 8080]
    api -->|SQL mediante sqlx| mysql[(MySQL 8<br/>Puerto 3306)]
    admin[Administrador] -->|Consulta técnica opcional| adminer[Adminer<br/>Puerto 8081]
    adminer -->|SQL| mysql
    compose[Docker Compose] -.->|Construye y ejecuta| api
    compose -.->|Ejecuta y verifica salud| mysql
    compose -.->|Ejecuta| adminer
```

Docker Compose coordina los tres contenedores. La API no expone directamente la base de datos a sus clientes; toda regla de negocio pasa por Go.

## Diagrama de componentes de la API

```mermaid
flowchart TB
    routes[Router Gin] --> authmw[Middleware JWT y roles]
    authmw --> ah[AuthHandler]
    authmw --> bh[BookHandler]
    authmw --> blh[LoanHandler]
    authmw --> lh[LaptopHandler]

    ah --> as[AuthService]
    bh --> bs[BookService]
    blh --> bls[LoanService]
    lh --> ls[LaptopService]

    as --> ur[UserRepository]
    bs --> br[BookRepository]
    bls --> lr[LoanRepository]
    bls --> br
    ls --> lapr[LaptopRepository]
    ls --> llr[LaptopLoanRepository]

    ur --> db[(MySQL)]
    br --> db
    lr --> db
    lapr --> db
    llr --> db
```

Los handlers traducen HTTP, los servicios expresan los casos de uso y los repositorios contienen SQL. El módulo nuevo respeta la estructura que ya tenía el proyecto.

## Secuencia: solicitar un portátil

```mermaid
sequenceDiagram
    actor U as Usuario
    participant M as Middleware JWT
    participant H as LaptopHandler
    participant S as LaptopService
    participant R as LaptopLoanRepository
    participant DB as MySQL

    U->>M: POST /api/laptop-loans + JWT + laptop_id
    M->>M: Validar firma y vencimiento
    alt Token inválido o ausente
        M-->>U: 401 Unauthorized
    else Token válido
        M->>H: user_id y petición
        H->>H: Validar JSON e ID
        alt Datos inválidos
            H-->>U: 400 Bad Request
        else Datos válidos
            H->>S: Borrow(user_id, laptop_id)
            S->>S: Calcular vencimiento a 7 días
            S->>R: Borrow(user_id, laptop_id, due_at)
            R->>DB: BEGIN
            R->>DB: SELECT status FOR UPDATE
            alt Equipo inexistente
                DB-->>R: Sin filas
                R-->>H: sql.ErrNoRows
                H-->>U: 404 Not Found
            else Equipo no disponible
                R->>DB: ROLLBACK
                R-->>H: ErrLaptopUnavailable
                H-->>U: 409 Conflict
            else Equipo disponible
                R->>DB: INSERT laptop_loans
                R->>DB: UPDATE laptops SET status='loaned'
                R->>DB: COMMIT
                R-->>S: Préstamo creado
                S-->>H: Préstamo creado
                H-->>U: 201 Created
            end
        end
    end
```
