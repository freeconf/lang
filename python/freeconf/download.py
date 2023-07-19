#!/usr/bin/env python3

name = "FreeCONF lang binary downloader"

desc = "Downloads the fc-lang binary from releases site for use with the FreeCONF applications"

usage = """
 This downloads the correct version of the fc-lang binary matching your installed freeconf python library.

 You do not run this script if you have already downloaded the library but it is harmless to run.

 If you store binary somewhere in the PATH, then freeconf python library will be able to find it. You can also use the FC_EXEC variable to explicitly set were binary should be found.

"""

# Only use packages that would have been installed when installing freeconf
# python package.  requests is safe.
import os
import requests
import argparse
import stat
import freeconf.driver

def run():
    parser = argparse.ArgumentParser(
        prog=name,
        description=desc,
        epilog=usage)

    parser.add_argument('-v', '--verbose', action='store_true')  # on/off flag
    parser.add_argument('-f', '--force', action='store_true')  # on/off flag
    parser.add_argument('-t', '--test', action='store_true', help='find for fc-lang or return non zero exit code')
    parser.add_argument('-d', '--dir', default=freeconf.driver.home_bin_dir(), help="where fc-lang binary should be located")  
    # This script should be designed to work on all OSes including Windows.
    args = parser.parse_args()

    def p(msg):
        if args.verbose:
            print(msg)

    if args.test:
        try:
            print(freeconf.driver.path_to_exe(verbose=args.verbose))
            exit(0)
        except freeconf.driver.ExecNotFoundException as e:
            p(e)
            exit(1)

    if not os.path.exists(args.dir):
        p(f"making {args.dir}")
        os.makedirs(args.dir)

    fname = freeconf.driver.exe_fname()
    file_path = os.path.join(args.dir, fname)
    if not os.path.isfile(file_path) or args.force:
        url = f"https://github.com/freeconf/lang/releases/download/v{freeconf.__version__}-alpha/{fname}"
        p(f"downloading {url} to {file_path}")
        resp = requests.get(url, stream=True)
        if resp.status_code >= 300:
            raise Exception("{resp._status_code} error downloading {url}")
        with open(file_path, 'wb') as fd:
            for chunk in resp.iter_content(chunk_size=2048):
                fd.write(chunk)
        p(f"downloaded successfully")
    else:
        p(f"{file_path} exists")

    if not os.access(file_path, os.X_OK):
        p("marking executable")
        st = os.stat(file_path)
        os.chmod(file_path, st.st_mode | stat.S_IEXEC)
    else:
        p("already executable")

if __name__ == '__main__':
    run()