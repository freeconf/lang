#include <jni.h>
#include "libfc.h"
#include "org_freeconf_fc_driver_Driver.h"


JNIEXPORT void JNICALL Java_org_freeconf_fc_driver_Driver_release
  (JNIEnv *env, jclass c, jlong poolId) {

    destruct(poolId);
}