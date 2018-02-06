# Progress Bars As A Service

Keep track of the progress of, well, anything!

# HTTP Interface

HTTP defaults to `application/json` unless you set the `Accept` header to `application/text`.

## POST /

Create a new progress bar value.

    $ curl -vX POST localhost:9999/
    *   Trying ::1...
    * TCP_NODELAY set
    * Connected to localhost (::1) port 9999 (#0)
    > POST / HTTP/1.1
    > Host: localhost:9999
    > User-Agent: curl/7.54.0
    > Accept: */*
    >
    < HTTP/1.1 200 OK
    < Content-Type: application/json
    < Date: Tue, 06 Feb 2018 15:36:54 GMT
    < Content-Length: 81
    <
    {"id":"81a8a58e-85e7-425d-95f3-66542e4c4212","progress":0,"token":"CNuYC19tDJ4"}

## PUT /{id}

Update a progress bar value.

    $ curl -vX PUT -d "token=CNuYC19tDJ4&progress=10" localhost:9999/81a8a58e-85e7-425d-95f3-66542e4c4212
    *   Trying ::1...
    * TCP_NODELAY set
    * Connected to localhost (::1) port 9999 (#0)
    > PUT /81a8a58e-85e7-425d-95f3-66542e4c4212 HTTP/1.1
    > Host: localhost:9999
    > User-Agent: curl/7.54.0
    > Accept: */*
    > Content-Length: 29
    > Content-Type: application/x-www-form-urlencoded
    >
    * upload completely sent off: 29 out of 29 bytes
    < HTTP/1.1 200 OK
    < Content-Type: application/json
    < Date: Tue, 06 Feb 2018 15:37:11 GMT
    < Content-Length: 60
    <
    {"id":"81a8a58e-85e7-425d-95f3-66542e4c4212","progress":10}

## GET /{id}

Get a progress bar value.

    $ curl -v localhost:9999/81a8a58e-85e7-425d-95f3-66542e4c4212
    *   Trying ::1...
    * TCP_NODELAY set
    * Connected to localhost (::1) port 9999 (#0)
    > GET /81a8a58e-85e7-425d-95f3-66542e4c4212 HTTP/1.1
    > Host: localhost:9999
    > User-Agent: curl/7.54.0
    > Accept: */*
    >
    < HTTP/1.1 200 OK
    < Content-Type: application/json
    < Date: Tue, 06 Feb 2018 15:37:21 GMT
    < Content-Length: 60
    <
    {"id":"81a8a58e-85e7-425d-95f3-66542e4c4212","progress":10}

# gRPC

The gRPC client is in `github.com/jmhobbs/pbaas/pb`


## Example

    $ cat testClient.go
    package main

    import (
      "context"
      "log"

      "github.com/jmhobbs/pbaas/pb"
      "google.golang.org/grpc"
    )

    func main() {
      conn, err := grpc.Dial("127.0.0.1:12345", grpc.WithInsecure())
      if err != nil {
        log.Fatalf("did not connect: %v", err)
      }
      defer conn.Close()

      c := pb.NewProgressBarServiceClient(conn)

      r, err := c.NewProgressBar(context.Background(), &pb.NewProgressBarRequest{5})
      if err != nil {
        log.Fatalf("could not create progress bar: %v", err)
      }

      log.Println("New Progress Bar:", r)
    }
    $ go run testClient.go
    2018/02/06 09:47:26 New Progress Bar: id:"1530d60f-0550-491a-88ee-e63b9a33ac02" token:"FIqsf2sV1h8"
    $
