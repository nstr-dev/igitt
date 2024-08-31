@ECHO OFF
if exist "" (
    "" %*
) else (
		echo.
		echo ================ igitt ================
		echo.
    echo Run following command to [32mcreate[0m an alias script [34m^(igitt =^> igt^)[0m:
		echo.
		echo [36migitt mkalias[0m
		echo.
		echo.
		echo If you don't want this alias, run this to [31mremove[0m this script:
		echo.
		echo [36mdel %~dpnx0[0m
)