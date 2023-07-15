from setuptools import setup
setup(
    install_requires=[
        "grpcio", 
        "google-api-python-client", 
        "grpcio_tools[dev]", 
        "requests[dev]"
    ]
)