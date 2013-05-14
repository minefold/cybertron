# Cybertron

Cybertron is a versioned storage system for file archives that sits on S3.

Cybertron uses DynamoDB to prevent you from uploading an earlier or conflicting revision of your archive.

## Local Setup

    psql -c "create database cybertron_development"
    go build && POSTGRES_URL=postgres://localhost/cybertron_development ./cybertron

## Heroku Setup

    git clone git://github.com/minefold/cybertron.git && cd cybertron
    heroku create
    heroku config:set AWS_ACCESS_KEY=abcdef AWS_SECRET_KEY=abcdef S3_URL=https://s3.amazonaws.com/cybertron
    git push heroku master
    heroku open


## Example Usage

    # upload initial archive
    curl -iu $CYBERKEY cybertron.com/path/archive?rev=10 \
        -X POST
        -F file=@localfile.tar.gz

    # update rev from 10 to 15
    curl -iu $CYBERKEY cybertron.com/path/archive?rev=10 \
        -X PATCH
        -H 'X-Rev:15'
        -F file=@localfile.tar.gz
    # returns 409 if 7 is not head

    # List revisions
    curl -iu $CYBERKEY cybertron.com/path/archive.json
    curl -iu $CYBERKEY cybertron.com/path/archive.json?params

    # download revision
    curl -u $CYBERKEY cybertron.com/path/archive # 404
    curl -u $CYBERKEY cybertron.com/path/archive.tar.gz
    curl -u $CYBERKEY cybertron.com/path/archive.tar.gz?rev=47

    # delete revisions
    curl -iu $CYBERKEY cybertron.com/path/archive \  # all revisions
        -X DELETE
    curl -iu $CYBERKEY cybertron.com/path/archive?rev=47 \  # one revision
        -X DELETE

## TODO

authentication
dlog
Support tar gz
Support tar lzo
Support content type. -H Content-Type=application/targz
Binary diff individual files in archive
Update individual files in archive
Support multiple download formats, through accepts header or extension
graceful restarts
