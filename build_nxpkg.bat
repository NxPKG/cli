cd ..\crates\nxpkgrepo
cargo build --target x86_64-pc-windows-gnu
if %errorlevel% neq 0 exit /b %errorlevel%

copy ..\..\target\x86_64-pc-windows-gnu\debug\nxpkg.exe ..\..\target\debug\nxpkg.exe
cd ..\..\cli