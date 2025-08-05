#!/usr/bin/env python3
"""
Passenger WSGI wrapper for Go application
This file helps Passenger understand how to start the Go binary
"""

import os
import subprocess
import sys

def application(environ, start_response):
    """
    WSGI application that proxies to Go binary
    """
    # This is a placeholder - Passenger will use app_start_command instead
    start_response('200 OK', [('Content-Type', 'text/plain')])
    return [b'Go application should be running via app_start_command']

# Start the Go application if this script is run directly
if __name__ == '__main__':
    port = os.environ.get('PORT', '8080')
    os.environ['PORT'] = port
    os.environ['GIN_MODE'] = 'release'
    
    # Execute the Go binary
    subprocess.call(['./app'])