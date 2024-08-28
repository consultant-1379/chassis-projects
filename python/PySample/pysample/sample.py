# -*- coding: utf-8 -*-
from PySample.pysample import helpers


def get_response():
    """Get response."""
    return 'this is response'


def response():
    """ if available return response """
    if helpers.isResponseAvailable():
        print(get_response())
