
# API package

The various controllers in API package make use of a few less-than-obvious mechanisms for reducing boilerplate.

## Struct tags

Each request is represented by a Go struct.  The fields on these structs contain `api:` struct tags that describe where in the request to find each piece of data.  Furthermore, there are `validate:` tags that define which values are permissible for each field.

The `validate:` tags are handled by the [gopkg.in/go-playground/validator.v9](https://github.com/go-playground/validator/tree/v9.24.0) package.  The only tags currently in use are:
- `gte=0`: a numeric value is expected, and it must be greater than or equal to 0
- `oneof=a b c`: the value must be one of the specified values (`a`, `b`, `c`).  Basically an enum.
- "Optional value" validators:
    - `omitempty`: use this for struct fields that are pointers.  This is only useful in combination with other validators (like `omitempty,gte=0`, which means "the value is optional, but must be >= 0 if it's present").
    - `isdefault`: use this for struct fields that are NOT pointers (like strings).  This is only useful in combination with other validators, and should be separated from them with a `|` pipe character (`isdefault|oneof=a b c`).

The `api:` struct tag takes two values separated by a comma.  The first value is the name of the field, and the second, prefixed by an `@`, describes where in the request that field is found.  This second field permits three possible values: `@query`, `@url_param`, and `@body`.

### `@query`

`@query` signals that the field should be obtained from the query string.  For example, the following request expects four query parameters, `offset`, `limit`, `ethereumAddress`, and `personId`:

```go
type MusicGroupListRequest struct {
    Offset       int            `api:"offset,@query" validate:"gte=0"`
    Limit        int            `api:"limit,@query"  validate:"gte=0"`
    EthereumAddr string         `api:"ethereumAddress,@query"`
    PersonID     *models.IDType `api:"personId,@query" validate:"omitempty,gte=0"`
}
```

The request URL will look something like `/api/musicgroup?offset=10&limit=5&ethereumAddress=0xdeadbeef&personId=6`

Note that the `ethereumAddress` and `personId` params can be omitted from the request URL entirely, and the request will still pass validation.


### `@url_param`

`@url_param` signals that the field should be obtained from one of the router-defined URL parameters.  In the router, these parameters are defined like so:

```go
r.Get("/cid/{cid}", rs.Cid)
```

The corresponding request struct would look like this:

```go
type MusicGroupGetByCIDRequest struct {
    CID string `api:"cid,@url_param" validate:"gte=0"`
}
```

Notice how the route param is called `cid` above, and the `api:` tag in the request struct also refers to `cid`.  These must match.


### `@body`

`@body` signals that the field should be decoded from the request body.  This field will very likely be a `struct`.  The name given to the field in the struct tag is ignored, but I tend to use `body` as in this example:

```go
type MusicGroupPostRequest struct {
    Body models.MusicGroup `api:"body,@body"`
}
```


# Authors/contributors

If the above is incomprehensible, reach out to one of the following folks:

- Bryn Bellomy (<bryn.bellomy@gmail.com>)


