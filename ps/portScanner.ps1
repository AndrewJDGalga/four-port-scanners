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

function Scan-PortsJob {
    $totalJobs = 65535
    $batchSize = 10
    $currentJobs = @()

    for($i = 1; $i -le $totalJobs; $i++){
        $nextJob = Start-Job -ScriptBlock {
            param([int] $port)

            $socket = New-Object System.Net.Sockets.TcpClient($successAddress, $port)
            if($socket.Connected){
                "Port open: $port"
                $socket.Close()
            }
        } -ArgumentList $i

    }
}

function Scan-AddressPorts {
    $successAddress = "127.0.0.1"
    $failAddress = "172.16.10.22"
    #port0 appears typically ignored?
    for($port = 1; $port -le 65535; $port++){
        try {
            $socket = New-Object System.Net.Sockets.TcpClient($successAddress, $port)
            if($socket.Connected){
                "Port open: $port"
                $socket.Close()
            }
        } catch [System.Net.Sockets.SocketException] { "Skipping $port"}
    }
}

Write-Host "Scanning..."
