
#include <freeconf.h>

fc_err* fc_err_new(char *msg) {
    fc_err* err = (fc_err*) malloc(sizeof(fc_err));
	strcpy(&err->message[0], msg);
	return err;
}