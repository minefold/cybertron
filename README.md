# Cybertron

Cybertron is a versioned storage system for file archives that sits on S3.

## Heroku Setup

    git clone git://github.com/minefold/cybertron.git && cd cybertron
    heroku create
    heroku config:set AWS_ACCESS_KEY_ID=abcdef AWS_SECRET_ACCESS_KEY=abcdef S3_BUCKET=cybertron
    git push heroku master
    heroku open


## Example Usage

    # upload initial revision
    curl -iu $CYBERKEY cybertron.com/path/archive \
        -X POST
        -F file=@localfile.tar.gz

    # upload the 8th revision
    curl -iu $CYBERKEY cybertron.com/path/archive?rev=7 \
        -X PATCH
        -F file=@localfile.tar.gz
    # returns 409 if 7 is not head

    # download latest
    curl -iu $CYBERKEY cybertron.com/path/archive # 404

    curl -iu $CYBERKEY cybertron.com/path/archive.tar.gz

    # list revisions
    curl -iu $CYBERKEY cybertron.com/path/archive \
        -X HEAD

    # download revision
    curl -iu $CYBERKEY cybertron.com/path/archive?rev=47

    # delete file with revisions
    curl -iu $CYBERKEY cybertron.com/path/archive \
        -X DELETE

    # delete revision
    curl -iu $CYBERKEY cybertron.com/path/archive?rev=47 \
        -X DELETE

## TODO

Support content type. -H Content-Type=application/targz
Binary diff individual files in archive
Update individual files in archive
Support multiple download formats, through accepts header or extension