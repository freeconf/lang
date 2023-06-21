import sys

gold_update = False

def assert_equal(tester, actual_str, expected_file):
    global gold_update
    if gold_update:
        with open(expected_file, 'w') as f:
            f.write(actual_str)
    else:
        with open(expected_file, 'r') as f:
            tester.assertMultiLineEqual(f.read(), actual_str)


def parse_flags():
    global gold_update
    if len(sys.argv) > 1 and sys.argv[1] == "-update":
        sys.argv.pop()
        gold_update = True
