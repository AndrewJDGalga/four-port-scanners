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
    $batchSize = 10 #expecting load
    $currentJobs = @()

    for($i = 1; $i -le $totalJobs; $i++){
        $nextJob = Start-Job -ScriptBlock {
            param([int] $port)

            try{
                $socket = New-Object System.Net.Sockets.TcpClient($successAddress, $port)
                if($socket.Connected){
                    "Port open: $port"
                    $socket.Close()
                }
            } catch [System.Net.Sockets.SocketException] { }
        } -ArgumentList $i
        $currentJobs += $nextJob

        Write-Host "Added job $i of $totalJobs."

        if(($currentJobs.Count % $batchSize) -eq 0){
            Write-Host "Processing $($currentJobs.Count) jobs..."

            Wait-Job $currentJobs
            Receive-Job -Job $currentJobs | ForEach-Object { Write-Output $_ }
            $currentJobs | Remove-Job
            $currentJobs = @() #clear() has proven unreliable
        }
    }

    if($currentJobs.Count -gt 0){
        Wait-Job $currentJobs
        Receive-Job -Job $currentJobs | ForEach-Object { Write-Output $_ }
        $currentJobs | Remove-Job
        $currentJobs = @()
    }
}

function Scan-AddressPorts {
    $successAddress = "127.0.0.1"
    #$failAddress = "172.16.10.22"
    for($port = 1; $port -le 65535; $port++){
        try {
            $socket = New-Object System.Net.Sockets.TcpClient($successAddress, $port)
            if($socket.Connected){
                "Port open: $port"
                $socket.Close()
            }
        } catch [System.Net.Sockets.SocketException] { }#"Skipping $port"}
    }
}

Write-Host "Scanning..."
Scan-PortsJob