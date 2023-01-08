#include <jni.h>
#include "libfc.h"
#include "org_freeconf_fc_parser_Parser.h"

JNIEXPORT jobject JNICALL Java_org_freeconf_fc_parser_Parser_parse
    (JNIEnv *env, jobject o, jstring ypath, jstring yname)
{
    const char *cstrYpath = (*env)->GetStringUTFChars(env, ypath, 0);
    const char *cstrYname = (*env)->GetStringUTFChars(env, yname, 0);
    Module m = parser((char *)cstrYpath, (char *)cstrYname);
    (*env)->ReleaseStringUTFChars(env, ypath, cstrYpath);
    (*env)->ReleaseStringUTFChars(env, yname, cstrYname);
    if (m.ident == NULL)
    {
        return NULL;
    }

    jclass cls = (*env)->FindClass(env, "org/freeconf/fc/meta/Module");
    jmethodID constructor = (*env)->GetMethodID(env, cls, "<init>", "(JLjava/lang/String;Ljava/lang/String;)V");
    if (constructor == NULL)
    {
        return NULL;
    }
    jstring ident = (*env)->NewStringUTF(env, m.ident);
    jstring desc = (*env)->NewStringUTF(env, m.desc);
    jobject object = (*env)->NewObject(env, cls, constructor, m.poolId, ident, desc);
    return object;
}