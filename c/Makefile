TEST_SRC = $(wildcard ./test/test_*.c)
TESTS = $(foreach T,$(TEST_SRC),$(basename $(notdir $(T))))
TEST_OBJ = $(foreach T,$(TESTS), out/$(T).out)
export LD_LIBRARY_PATH = ./lib
export YANGPATH = $(abspath ../yang)

test: $(TESTS)

debug:
	echo $(TEST_SRC)
	echo $(TESTS)
	echo $(TEST_OBJ)

$(TESTS) : % : out/%.out
	./out/$@.out

out/%.out : test/%.c
	gcc -Wall -I./lib -L./lib ./test/$*.c -o $@ -lfc
