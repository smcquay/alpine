Simple cli tool for serving files in a directory.

Installation:

    go get mcquay.me/servedir

use:

    servedir  # default port

    # specify port
    servedir -port 6666

    # allow serving hidden files/dirs
    servedir -hidden

    # serve https, with http redirect
    TLS_CERT=/path/to/cert.pem TLS_KEY=/path/to/key.pem servedir

    # or see -help
