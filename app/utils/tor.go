package utils

import (
    "net/url"
    "golang.org/x/net/proxy"
    "net/http"
    "io/ioutil"
    "log"
)

func FetchViaTOR(requestUrl string) []byte {
    // Create a transport that uses Tor Browser's SocksPort.  If
    // talking to a system tor, this may be an AF_UNIX socket, or
    // 127.0.0.1:9050 instead.
    tbProxyURL, err := url.Parse("socks5://127.0.0.1:9050")
    if err != nil {
        log.Fatalf("Failed to parse proxy URL: %v\n", err)
    }

    // Get a proxy Dialer that will create the connection on our
    // behalf via the SOCKS5 proxy.  Specify the authentication
    // and re-create the dialer/transport/client if tor's
    // IsolateSOCKSAuth is needed.
    tbDialer, err := proxy.FromURL(tbProxyURL, proxy.Direct)
    if err != nil {
        log.Fatalf("Failed to obtain proxy dialer: %v\n", err)
    }

    // Make a http.Transport that uses the proxy dialer, and a
    // http.Client that uses the transport.
    tbTransport := &http.Transport{Dial: tbDialer.Dial}
    client := &http.Client{Transport: tbTransport}

    // Example: Fetch something.  Real code will probably want to use
    // client.Do() so they can change the User-Agent.
    resp, err := client.Get(requestUrl)
    if err != nil {
        log.Fatalf("Failed to issue GET request: %v\n", err)
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Fatalf("Failed to read the body: %v\n", err)
    }
    return body
}
