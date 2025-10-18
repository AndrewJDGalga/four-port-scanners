param (
    [array]$addrsStr=("127.0.0.1"),
    [string]$portsStr=("1,65535")
)

enum AddressOperationTypes {
    AddressSingle
    AddressRange
    AddressList
}

enum PortOperationTypes{
    PortSingle
    PortRange
    PortList
}

$targetPorts = $portsStr.Split(",")
$portCount = $targetPorts.Length


'''
for($i = 0; $i -le $operationType; $i++){
    Write-Host $targetPorts[$i]
}
'''