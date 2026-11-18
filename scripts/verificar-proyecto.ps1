$ErrorActionPreference = 'Stop'
$baseUrl = 'http://localhost:8080'
$results = [System.Collections.Generic.List[object]]::new()

function Invoke-Api {
    param(
        [string]$Method,
        [string]$Path,
        [string]$Token = '',
        [object]$Body = $null
    )

    $params = @{
        Uri                = "$baseUrl$Path"
        Method             = $Method
        SkipHttpErrorCheck = $true
    }
    if ($Token) {
        $params.Headers = @{ Authorization = "Bearer $Token" }
    }
    if ($null -ne $Body) {
        $params.ContentType = 'application/json'
        $params.Body = $Body | ConvertTo-Json -Compress
    }
    return Invoke-WebRequest @params
}

function Confirm-Status {
    param([string]$Case, [int]$Expected, $Response)

    $actual = [int]$Response.StatusCode
    $results.Add([pscustomobject]@{
        Caso      = $Case
        Esperado  = $Expected
        Obtenido  = $actual
        Resultado = if ($actual -eq $Expected) { 'OK' } else { 'FALLÓ' }
    })
    if ($actual -ne $Expected) {
        throw "$Case esperaba HTTP $Expected y obtuvo $actual. Respuesta: $($Response.Content)"
    }
}

$stamp = Get-Date -Format 'yyyyMMddHHmmss'
$studentEmail = "final.student.$stamp@example.com"
$otherEmail = "final.other.$stamp@example.com"
$adminEmail = "final.admin.$stamp@example.com"
$password = 'pass@123'

$health = Invoke-Api GET '/health'
Confirm-Status 'Health' 200 $health
$ping = Invoke-Api GET '/ping'
Confirm-Status 'Ping' 200 $ping

$studentRegister = Invoke-Api POST '/api/auth/register' -Body @{
    name = 'Estudiante prueba final'; email = $studentEmail; password = $password
}
Confirm-Status 'Registrar estudiante' 201 $studentRegister
$studentToken = ($studentRegister.Content | ConvertFrom-Json).token

$otherRegister = Invoke-Api POST '/api/auth/register' -Body @{
    name = 'Segundo estudiante'; email = $otherEmail; password = $password
}
Confirm-Status 'Registrar segundo estudiante' 201 $otherRegister
$otherToken = ($otherRegister.Content | ConvertFrom-Json).token

$adminRegister = Invoke-Api POST '/api/auth/register' -Body @{
    name = 'Administrador prueba final'; email = $adminEmail; password = $password
}
Confirm-Status 'Registrar usuario para admin' 201 $adminRegister

& docker compose exec -T mysql mysql -uroot -proot_pass biblioteca -e "UPDATE users SET role='admin' WHERE email='$adminEmail';" 2>$null
if ($LASTEXITCODE -ne 0) { throw 'No se pudo asignar el rol admin para la prueba.' }

$adminLogin = Invoke-Api POST '/api/auth/login' -Body @{ email = $adminEmail; password = $password }
Confirm-Status 'Login de administrador' 200 $adminLogin
$adminToken = ($adminLogin.Content | ConvertFrom-Json).token

$me = Invoke-Api GET '/api/auth/me' -Token $studentToken
Confirm-Status 'Perfil autenticado' 200 $me
$books = Invoke-Api GET '/api/books' -Token $studentToken
Confirm-Status 'Listar libros existentes' 200 $books
$bookID = ($books.Content | ConvertFrom-Json)[0].id
$bookLoan = Invoke-Api POST '/api/loans' -Token $studentToken -Body @{ book_id = $bookID }
Confirm-Status 'Préstamo de libro existente' 201 $bookLoan
$bookLoanID = ($bookLoan.Content | ConvertFrom-Json).id
$bookReturn = Invoke-Api PUT "/api/loans/$bookLoanID/return" -Token $studentToken
Confirm-Status 'Devolución de libro existente' 200 $bookReturn

$anonymousList = Invoke-Api GET '/api/laptops'
Confirm-Status 'Catálogo sin token' 401 $anonymousList
$studentCreate = Invoke-Api POST '/api/laptops' -Token $studentToken -Body @{
    brand = 'Lenovo'; model = 'ThinkPad E14'; serial_number = "UTB-DENIED-$stamp"
}
Confirm-Status 'Estudiante no puede crear portátil' 403 $studentCreate

$serial = "UTB-LAP-$stamp"
$created = Invoke-Api POST '/api/laptops' -Token $adminToken -Body @{
    brand = 'Lenovo'; model = 'ThinkPad E14'; serial_number = $serial
}
Confirm-Status 'Admin crea portátil' 201 $created
$laptop = $created.Content | ConvertFrom-Json

$duplicate = Invoke-Api POST '/api/laptops' -Token $adminToken -Body @{
    brand = 'Lenovo'; model = 'ThinkPad E14'; serial_number = $serial
}
Confirm-Status 'Serial duplicado' 409 $duplicate

$catalog = Invoke-Api GET '/api/laptops' -Token $studentToken
Confirm-Status 'Usuario lista portátiles' 200 $catalog

$borrow = Invoke-Api POST '/api/laptop-loans' -Token $studentToken -Body @{ laptop_id = $laptop.id }
Confirm-Status 'Solicitar portátil disponible' 201 $borrow
$loan = $borrow.Content | ConvertFrom-Json

$duplicateLoan = Invoke-Api POST '/api/laptop-loans' -Token $otherToken -Body @{ laptop_id = $laptop.id }
Confirm-Status 'Impedir préstamo de equipo ocupado' 409 $duplicateLoan

$myLoans = Invoke-Api GET '/api/laptop-loans/me' -Token $studentToken
Confirm-Status 'Listar mis préstamos de portátil' 200 $myLoans

$foreignReturn = Invoke-Api PUT "/api/laptop-loans/$($loan.id)/return" -Token $otherToken
Confirm-Status 'Impedir devolución ajena' 403 $foreignReturn
$invalidReturn = Invoke-Api PUT '/api/laptop-loans/999999999/return' -Token $studentToken
Confirm-Status 'Devolución inexistente' 404 $invalidReturn
$return = Invoke-Api PUT "/api/laptop-loans/$($loan.id)/return" -Token $studentToken
Confirm-Status 'Devolver portátil propio' 200 $return
$repeatedReturn = Invoke-Api PUT "/api/laptop-loans/$($loan.id)/return" -Token $studentToken
Confirm-Status 'Impedir devolución repetida' 409 $repeatedReturn

$results | Format-Table -AutoSize
Write-Host "`nVERIFICACIÓN COMPLETA: $($results.Count) casos aprobados." -ForegroundColor Green
