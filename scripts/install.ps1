param(
  [string]$Version = "",
  [string]$Repo = "Senasphy/monkeytui",
  [string]$BinName = "monkeytui",
  [string]$InstallDir = "$env:LOCALAPPDATA\Programs\monkeytui\bin"
)

$ErrorActionPreference = "Stop"

function Resolve-Arch {
  $arch = [System.Runtime.InteropServices.RuntimeInformation]::OSArchitecture.ToString().ToLowerInvariant()
  switch ($arch) {
    "x64" { return "amd64" }
    "arm64" { return "arm64" }
    default { throw "Unsupported architecture: $arch" }
  }
}

function Resolve-LatestVersion {
  $latest = Invoke-RestMethod -Uri "https://api.github.com/repos/$Repo/releases/latest" `
    -Headers @{ "User-Agent" = "monkeytui-installer" }
  return $latest.tag_name
}

if ([string]::IsNullOrWhiteSpace($Version)) {
  $Version = $env:MONKEYTUI_VERSION
}
if ([string]::IsNullOrWhiteSpace($Version)) {
  $Version = Resolve-LatestVersion
}
if ([string]::IsNullOrWhiteSpace($Version)) {
  throw "Failed to resolve release version"
}

# normalize version
$Version = $Version.TrimStart("v")

$arch = Resolve-Arch
$os = "windows"

$artifact = "${BinName}_${Version}_${os}_${arch}.zip"
$url = "https://github.com/$Repo/releases/download/v$Version/$artifact"

$tmpDir = Join-Path $env:TEMP ("monkeytui-" + [guid]::NewGuid().ToString("N"))
New-Item -ItemType Directory -Path $tmpDir | Out-Null

try {
  $archivePath = Join-Path $tmpDir $artifact

  Write-Host "Downloading $url ..."
  Invoke-WebRequest -Uri $url -OutFile $archivePath -ErrorAction Stop

  Expand-Archive -Path $archivePath -DestinationPath $tmpDir -Force

  $sourceExe = Join-Path $tmpDir "$BinName.exe"
  if (-not (Test-Path $sourceExe)) {
    throw "Binary not found in archive: $sourceExe"
  }

  New-Item -ItemType Directory -Force -Path $InstallDir | Out-Null
  $targetExe = Join-Path $InstallDir "$BinName.exe"
  Copy-Item -Force $sourceExe $targetExe

  # Add to PATH if missing
  $userPath = [Environment]::GetEnvironmentVariable("Path", "User")
  if ($userPath -notlike "*$InstallDir*") {
    $newPath = if ([string]::IsNullOrWhiteSpace($userPath)) {
      $InstallDir
    } else {
      "$userPath;$InstallDir"
    }
    [Environment]::SetEnvironmentVariable("Path", $newPath, "User")
    Write-Host "Added $InstallDir to PATH. Restart your terminal."
  }

  Write-Host "Installed $BinName $Version → $targetExe"
}
finally {
  if (Test-Path $tmpDir) {
    Remove-Item -Recurse -Force $tmpDir
  }
}
