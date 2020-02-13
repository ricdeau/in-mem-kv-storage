# in-mem-kv-storage

To launch server:

    $ git clone github.com/ricdeau/in-mem-kv-storage
 
    $ cd ./in-mem-kv-storage
 
    $ go build -o <out-file-name> github.com/ricdeau/in-mem-kv-storage
 
    $ ./<out-file-name> 
    to launch on port 5339 with default options or
    $ ./<out-file-name> --help
    to see other options
    
Store value example:
    
    curl -X PUT -H 'Content-Type: <content-type>' -d 'some-data' <base-address>/api/storage/<key>
    
    If new value returns 201 if updated 200.
    
Get value example:

    curl -X GET <base-address>/api/storage/<key>
    
    If value with this key doen't exist returns 404 otherwise 200 and value.
    
Delete value example:
    
    curl -X DELETE <base-address>/api/storage/<key>
    
    Returns 204, removes value if it exists.
    

To launch client:

    In project's directory 
    
    $ cd ./client
    
    $ go build -o client
    
    $ ./client --help
    to see options
    
Client usage example:
    
    ./client -c 100 -n 100000 -ct image/jpg -m PUT -t http://localhost:5339/api/storage/key -f ~/tQXloGfkwMg.jpg
    Total time: 4.477003028s
    Average time per request: 44.77µs
    Average rps: 22336.37

    
 