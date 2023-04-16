if [ "$1" == "off" ]; then
unset FC_LANG_EXEC
unset FC_LANG_DBG_ADDR
else
export FC_LANG_EXEC="${HOME}/work/lang/bin/fc-lang-dbg"
export FC_LANG_DBG_ADDR=":9999"
fi