# CSR signer sketch

**For informational purposes only**

This repository is just a sketch to outline how to

- Generate an ED25519 based keypair
- Create a Certificate Signing Request (CSR)
- Send the CSR to a server
- Sign the certificate in the CSR
- Return signed certificate to the client

## Building

This demo consists of a `client`and a `server` binary.  You build both with make:

```shell
make
```

Now you should have `bin/client` and `bin/server`.

## Running

First start the server

```shell
bin/server
```

...then run the client

```shell
bin/client
```

On the server side you should see something like this when you start the server and run the client once:

```text
021/09/23 16:25:13 server up, listening to :8881
2021/09/23 16:25:16 Got CSR from CN=sample client certificate,1.2.840.113549.1.9.1=#0c1075736572406578616d706c652e636f6d with signature f2dd572564530e1050161fb1ecfbf6b95ad74ededada3d85ce18540cfe6a3a143ce1ce375506723c49dadcec7f9c00fa4a09df36cf8e95c27a11bc22db943205 
2021/09/23 16:25:16 signature ok
2021/09/23 16:25:16 created certificate:
-----BEGIN CERTIFICATE-----
MIIBNDCB56ADAgECAgECMAUGAytlcDAaMRgwFgYDVQQKEw9CbGluZCBGYWl0aCBJ
bmMwHhcNMjEwOTIzMTQyNTE2WhcNMjEwOTI0MTQyNTE2WjAkMSIwIAYDVQQDExlz
YW1wbGUgY2xpZW50IGNlcnRpZmljYXRlMCowBQYDK2VwAyEAOG9oDOl+3dY0SSMi
eaDJqVRYzbmBbBivlW8lmsONhWSjSDBGMA4GA1UdDwEB/wQEAwIHgDATBgNVHSUE
DDAKBggrBgEFBQcDAjAfBgNVHSMEGDAWgBTn9fa2sejzG9W44lc+DswCFSBHcDAF
BgMrZXADQQAM9QQmD3AGtmbtJ2a75XrXzwaUKMIiV8DLTIEpUQgS7J5Gqlw9FCrP
ktUpifD94RcNEHaxj4GWwnh+vKZjZXkL
-----END CERTIFICATE-----
```

On the client side you should see something like this:

```text
-----BEGIN PUBLIC KEY-----
MCowBQYDK2VwAyEAOG9oDOl+3dY0SSMieaDJqVRYzbmBbBivlW8lmsONhWQ=
-----END PUBLIC KEY-----

-----BEGIN PRIVATE KEY-----
MC4CAQAwBQYDK2VwBCIEICAWyEHgiJFnq0XhhaTtOmJrKFxL697SwpWivuzkdHV3
-----END PRIVATE KEY-----

Client certificate signed by server:
-----BEGIN CERTIFICATE-----
MIIBNDCB56ADAgECAgECMAUGAytlcDAaMRgwFgYDVQQKEw9CbGluZCBGYWl0aCBJ
bmMwHhcNMjEwOTIzMTQyNTE2WhcNMjEwOTI0MTQyNTE2WjAkMSIwIAYDVQQDExlz
YW1wbGUgY2xpZW50IGNlcnRpZmljYXRlMCowBQYDK2VwAyEAOG9oDOl+3dY0SSMi
eaDJqVRYzbmBbBivlW8lmsONhWSjSDBGMA4GA1UdDwEB/wQEAwIHgDATBgNVHSUE
DDAKBggrBgEFBQcDAjAfBgNVHSMEGDAWgBTn9fa2sejzG9W44lc+DswCFSBHcDAF
BgMrZXADQQAM9QQmD3AGtmbtJ2a75XrXzwaUKMIiV8DLTIEpUQgS7J5Gqlw9FCrP
ktUpifD94RcNEHaxj4GWwnh+vKZjZXkL
-----END CERTIFICATE-----

Issuer: O=Blind Faith Inc
Authority Key ID: e7f5f6b6b1e8f31bd5b8e2573e0ecc0215204770
Public key algorithm: Ed25519
```

## `csrparse`

This is a utility to parse and dump CSR.  This can be used to verify that you can parse the CSR from other systems if you need to.
