from setuptools import setup

setup( 
    install_requires=[
        "grpcio", 
        "google-api-python-client", 
        "grpcio_tools[dev]", 
        "requests[dev]"
    ],

    # Need a python build master to figure out how to fake, cross-compile python
    # packages to build {os}-{arch} specific packages from single machine. Until
    # then I'm adding 17MB binary for each distro and pick right one at runtime.
    # at the expense of creating a larger dist than needed.
    data_files=[('', [
        'freeconf/data/fc-lang-0.1.0-darwin-amd64',
        'freeconf/data/fc-lang-0.1.0-darwin-arm64',
        'freeconf/data/fc-lang-0.1.0-linux-amd64',
        'freeconf/data/fc-lang-0.1.0-windows-amd64.exe'
    ])],
)
