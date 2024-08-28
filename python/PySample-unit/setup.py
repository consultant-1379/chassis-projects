#!/usr/bin/env python

from setuptools import setup

long_description = """
SAmple python project used to test ci pipelines for Mason2
"""

setup(
    name='PySample',
    version='0.0.1',
    description='Sample project for testing CI pipelines',
    long_description=long_description.strip(),
    packages=['pysample'],
    python_requires=">=3.4",
    install_requires=[
    ]
)
