param (
    [array]$addrsStr=("127.0.0.1"),
    [string]$portsStr=("1,65535")
)

$targetPorts = $portsStr.Split(",")
for($i = 0; $i -le $targetPorts.Length; $i++){
    Write-Host $targetPorts[$i]
}
#Write-Host $targetPorts.Length