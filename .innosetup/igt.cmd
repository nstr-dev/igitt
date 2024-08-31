@ECHO OFF
if exist "" (
	"" %*
) else (
	igitt mkalias
	igitt %*
)