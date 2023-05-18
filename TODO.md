yang
x  upsert* methods should return err, not selection
o  rename Rpc to Action

core
o  python weakrefs (gc)
o  full support of meta package
x  notifications
x  more val support incl. leaf-lists
?  val coersion
o  how to get things from context
o  choice

nodeutil
o  json writer
o  restconf server
o  gnmi

test suite
x  build driver compliance test harness

package
o  submit to Pypy
o  document py
o  document lang

bugs/tests
o python doesn't always exit after go stops in failed unit tests
o performance test results