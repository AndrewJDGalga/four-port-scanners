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

function Scan-AddressPortsAsync {
    $successAddress = "127.0.0.1"
    #$failAddress = "172.16.10.22"
    $ports = 1..65535

    $ports | Foreach-Object -Parallel {
        param($port, $successAddress)
        try {
            $socket = New-Object System.Net.Sockets.TcpClient($successAddress, $port)
            if($socket.Connected){
                "Port open: $port"
                $socket.Close()
            }
        } catch [System.Net.Sockets.SocketException] { "Skipping $port"}
    } -ArgumentList $_, $successAddress -ThrottleLimit 5
}

function Scan-AddressPorts {
    $successAddress = "127.0.0.1"
    $failAddress = "172.16.10.22"
    #port0 appears typically ignored?
    for($port = 1; $port -le 65535; $port++){
        #if (Test-Connection -Quiet -Count 1 -ComputerName $successAddress){
            try {
                $socket = New-Object System.Net.Sockets.TcpClient($successAddress, $port)
                if($socket.Connected){
                    "Port open: $port"
                    $socket.Close()
                }
            } catch [System.Net.Sockets.SocketException] { "Skipping $port"}
        #}
    }
}

Write-Host "Scanning..."
Scan-AddressPortsAsync