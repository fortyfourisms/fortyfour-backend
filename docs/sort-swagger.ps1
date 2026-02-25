# sort-swagger.ps1
$json = Get-Content swagger.json -Raw | ConvertFrom-Json

# buat dictionary ordered
$sorted = [ordered]@{}

# salin key-value ke sorted
$sorted['swagger']   = $json.swagger
$sorted['info']      = $json.info
$sorted['host']      = $json.host
$sorted['basePath']  = $json.basePath
$sorted['schemes']   = $json.schemes

# sort paths & definitions lalu masukkan
$sortedPaths      = $json.paths.psobject.properties | Sort-Object Name
$sortedDefinitions= $json.definitions.psobject.properties | Sort-Object Name

$pathsOrdered      = [ordered]@{}
foreach ($p in $sortedPaths)      { $pathsOrdered[$p.Name]      = $p.Value }

$defsOrdered       = [ordered]@{}
foreach ($d in $sortedDefinitions) { $defsOrdered[$d.Name]   = $d.Value }

$sorted['paths']      = $pathsOrdered
$sorted['definitions']= $defsOrdered

# tulis kembali
$sorted | ConvertTo-Json -Depth 100 | Set-Content swagger.json