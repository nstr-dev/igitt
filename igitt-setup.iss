; Script generated by the Inno Setup Script Wizard.
; SEE THE DOCUMENTATION FOR DETAILS ON CREATING INNO SETUP SCRIPT FILES!

#define MyAppName "Igitt"
#define MyAppPublisher "nstr.dev"
#define MyAppURL "nstr.dev"
#define MyAppExeName "igitt.exe"

#ifndef MyAppVersion
  #define MyAppVersion "1.0.0"
#endif

[Setup]
; NOTE: The value of AppId uniquely identifies this application. Do not use the same AppId value in installers for other applications.
; (To generate a new GUID, click Tools | Generate GUID inside the IDE.)
AppId=239D7644-5A09-407F-A751-13577B4845BC
AppName={#MyAppName}
AppVersion={#MyAppVersion}
;AppVerName={#MyAppName} {#MyAppVersion}
AppPublisher={#MyAppPublisher}
AppPublisherURL={#MyAppURL}
AppSupportURL={#MyAppURL}
AppUpdatesURL={#MyAppURL}
DefaultDirName={autopf}\{#MyAppName}
; "ArchitecturesAllowed=x64compatible" specifies that Setup cannot run
; on anything but x64 and Windows 11 on Arm.
ArchitecturesAllowed=x64compatible
; "ArchitecturesInstallIn64BitMode=x64compatible" requests that the
; install be done in "64-bit mode" on x64 or Windows 11 on Arm,
; meaning it should use the native 64-bit Program Files directory and
; the 64-bit view of the registry.
ArchitecturesInstallIn64BitMode=x64compatible
DefaultGroupName={#MyAppName}
DisableProgramGroupPage=yes
; Remove the following line to run in administrative install mode (install for all users.)
PrivilegesRequired=lowest
; PrivilegesRequiredOverridesAllowed=dialog
OutputBaseFilename=igitt-setup
OutputDir=.innosetup
Compression=lzma
SolidCompression=yes
WizardStyle=modern
ChangesEnvironment=yes
UninstallDisplayName={#MyAppName}
UninstallDisplayIcon={app}\{#MyAppExeName}

[UninstallDelete]
Type: filesandordirs; Name: "{app}"

[Registry]
; Add program directory to PATH for current user if the checkbox is selected
Root: HKCU; Subkey: "Environment"; ValueType: expandsz; ValueName: "Path"; \
  ValueData: "{olddata};{app}"; Check: ShouldAddPathCheck

; Register application for "Add/Remove Programs"
Root: HKCU; Subkey: "Software\Microsoft\Windows\CurrentVersion\Uninstall\{#MyAppName}"; \
  ValueType: string; ValueName: "DisplayName"; ValueData: "{#MyAppName}"

Root: HKCU; Subkey: "Software\Microsoft\Windows\CurrentVersion\Uninstall\{#MyAppName}"; \
  ValueType: string; ValueName: "UninstallString"; ValueData: """{uninstallexe}"""

Root: HKCU; Subkey: "Software\Microsoft\Windows\CurrentVersion\Uninstall\{#MyAppName}"; \
  ValueType: string; ValueName: "DisplayVersion"; ValueData: "{#MyAppVersion}"

Root: HKCU; Subkey: "Software\Microsoft\Windows\CurrentVersion\Uninstall\{#MyAppName}"; \
  ValueType: string; ValueName: "Publisher"; ValueData: "{#MyAppPublisher}"

[Languages]
Name: "english"; MessagesFile: "compiler:Default.isl"

[Files]
Source: "{#MyAppExeName}"; DestDir: "{app}";
Source: ".innosetup\igitt.exe.log"; DestDir: "{app}";
Source: ".innosetup\igt.cmd"; DestDir: "{app}";
Source: ".innosetup\igittconfig.yaml"; DestDir: "{app}";
; NOTE: Don't use "Flags: ignoreversion" on any shared system files

[Icons]
Name: "{group}\{#MyAppName}"; Filename: "{app}\{#MyAppExeName}"

[Code]
var
  AddToPathPage: TInputOptionWizardPage;

function NeedsAddPath(Param: string): Boolean;
var
  OrigPath: string;
begin
  if not RegQueryStringValue(HKCU, 'Environment', 'Path', OrigPath) then
  begin
    Result := True;
    exit;
  end;
  { Ensure we search for ;Param; in a consistent case to prevent duplicates }
  Result := Pos(';' + UpperCase(Param) + ';', ';' + UpperCase(OrigPath) + ';') = 0;
end;

function ShouldAddPathCheck: Boolean;
begin
  { Return true only if the user has checked the option and the path isn’t already present }
  Result := AddToPathPage.Values[0] and NeedsAddPath('{app}');
end;

procedure RemoveFromPath();
var
  PathValue: string;
  AppPath: string;
begin
  if RegQueryStringValue(HKCU, 'Environment', 'Path', PathValue) then
  begin
    AppPath := ExpandConstant('{app}');
    { Remove only our application's path safely }
    StringChangeEx(PathValue, ';' + AppPath, '', True);
    StringChangeEx(PathValue, AppPath + ';', '', True);
    StringChangeEx(PathValue, AppPath, '', True);
    { Save the modified PATH value or delete the key if it's empty }
    if PathValue <> '' then
      RegWriteStringValue(HKCU, 'Environment', 'Path', PathValue)
    else
      RegDeleteValue(HKCU, 'Environment', 'Path');
  end;
end;

procedure CurUninstallStepChanged(CurUninstallStep: TUninstallStep);
begin
  if CurUninstallStep = usPostUninstall then  { Run after uninstallation }
  begin
    RemoveFromPath();
  end;
end;

procedure InitializeWizard;
begin
  { Create a custom page after the Select Destination Location page }
  AddToPathPage := CreateInputOptionPage(wpSelectDir,
    'Additional Options',
    'Select additional options:',
    'Choose whether to add the application folder to the PATH environment variable.',
    True, False);
  AddToPathPage.Add('Add application folder to PATH');
  AddToPathPage.Add('Do not add folder to PATH');
  { Set default state (checked) }
  AddToPathPage.Values[0] := True;
end;