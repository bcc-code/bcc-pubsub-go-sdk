# bcc-pubsub Go SDK

Go client for [bcc-pubsub](https://github.com/bcc-code/bcc-pubsub): publish events and manage subscriptions.

## Verifying pushes

Every push from bcc-pubsub carries an `X-PubSub-Token` header, signed by bcc-pubsub for your client id and the exact
body. Check it before you handle the message, so your endpoint only accepts pushes that bcc-pubsub sent to you.

`RequireDelivery` is `net/http` middleware. It answers 401 to a push without a valid token, and 503 when bcc-pubsub's
public keys cannot be fetched (bcc-pubsub then delivers the message again):

```go
verifier, err := pubsubsdk.NewDeliveryVerifier(pubsubsdk.DeliveryVerifierOptions{
    Issuer:   pubsubsdk.DeliveryIssuerProd, // or DeliveryIssuerSandbox, DeliveryIssuerDev
    Audience: os.Getenv("PUBSUB_CLIENT_ID"), // your own client id
})
if err != nil {
    log.Fatal(err)
}

http.Handle("/hooks/member-changed", pubsubsdk.RequireDelivery(verifier, 0, http.HandlerFunc(
    func(w http.ResponseWriter, r *http.Request) {
        claims, _ := pubsubsdk.DeliveryClaimsFrom(r.Context()) // subscription or dispatch, message id
        // r.Body is the verified body.
        w.WriteHeader(http.StatusOK)
    })))
```

Create one verifier and reuse it: it caches the public keys for an hour and fetches them again when bcc-pubsub starts
using a new key. With another router, call `verifier.Verify(ctx, r.Header.Get(pubsubsdk.DeliveryTokenHeader), body)`
on the raw body yourself. It returns a `*DeliveryVerificationError` (answer 401) or an error wrapping
`ErrDeliveryKeysUnavailable` (answer 5xx).

The token format and the key rotation are described in
[delivery verification](https://github.com/bcc-code/bcc-pubsub/blob/main/docs/delivery-verification.md). The tests run
the shared test vectors in `pubsubsdk/testdata/delivery-token`; `Scripts/update-test-vectors.sh` copies them from a
`bcc-pubsub` checkout.

## Code generation

The models in `pubsubsdk/model` are generated from the bcc-pubsub OpenAPI document with `Scripts/generate-sdks.sh`.
