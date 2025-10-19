param (
    [array]$addrsStr=("127.0.0.1"),
    [string]$portsStr=("1,65535")
)

enum AddressOperationTypes {
    AddressSingle
    AddressRange
    AddressList
    AddressErr
}
$addrType = [AddressOperationTypes]::AddressErr

enum PortOperationTypes{
    PortSingle
    PortRange
    PortList
    PortErr
}
$portType = [PortOperationTypes]::PortErr

$targetPorts = $portsStr.Split(",")
$targetAddrs = $addrsStr.Split(",")

if (Test-Connection -Quiet -Count 1 -ComputerName 1.1.1.1){
    Write-Host "Computer up"
}

function Scan-AddressPorts(){
    $testAddress = 1.1.1.1
    for($i = 0; $i -le 1..65535; $i++){
        
    }
}
