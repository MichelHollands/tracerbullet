# tracerbullet

This module enables debug logging for 1 single HTTP request by adding a user defined header when using the [slog](golang.org/x/exp/slog) package. For an example see [example/main.go](example/main.go).

## Usage

It works by creating the following:
- a type that implements the AccessChecker interface
- a middleware that takes a handler, a logger and an AccessChecker as parameters
- a [RoundTripper](https://pkg.go.dev/net/http#RoundTripper) that takes a Transport, a logger and an AccessChecker as parameters
- a slog [Logger](https://pkg.go.dev/golang.org/x/exp/slog#Logger) that uses the [Handler](handler.go) defined in this package

The middleware should be wrapped around an existing [http.Handler](https://pkg.go.dev/net/http#Handler). It checks the HTTP Header specified in the AccessChecker. If that exists and has the correct value then the X-SlogTracer-Set header is set in the context. The slog handler will check for the existince of that header and enable debug logging if so.

The RoundTripper has to be added to any Transport that connects to other HTTP service. If needed it sets the HTTP header specified in the AccessChecker in the outgoing HTTP requests.

[NewStaticAccessChecker](checker.go) creates an AccessChecker that uses the hard-coded header and value specified in it's parameters. This should not be used in production.

[example/main.go](example/main.go) shows all of these used in a test program.

## Docker compose example

To run an example follow these steps:

```
cd example
docker compose up --build --remove-orphans
```

This runs 3 services. Server1 calls server2 which in turn calls server3. 

Call the `/action` endpoint on the first service:

```
curl http://localhost:3000/action
```

Note that only log lines with level=INFO are shown. Now call the `/action` endpoint on the first service with the header specified in the AccessChecker:

```
curl -H "X-TracerBullet: abcdef" http://localhost:3000/action
````

Note that log lines with level=DEBUG are shown as well.
