Getting - A hypermedia client for golang
============================================

Introduction
------------

The Getting library is an attempt at creating a 'generic' hypermedia client, it
supports an opinionated set of modern features REST services might have.

It is a port of the amazing [Ketting][1] library for javascript.

The library supports [HAL][hal], [JSON:API][jsonapi],
[Web Linking (HTTP Link Header)][2] and HTML5 links.

### Example

Coming soon.

Docs
----

Coming soon.

Notable Features
----------------

Getting provides a RESTful interface and make it easier to follow REST best
practices more strictly.

It provides some useful abstractions that make it easier to work with true
hypermedia / HATEAOS servers. It currently parses [HAL][hal] and has a deep
understanding of links and embedded resources. There's also support for parsing
and following links from HTML documents, and it understands the HTTP `Link:`
header.

Using this library it becomes very easy to follow links from a single bookmark,
and discover resources and features on the server.

Supported formats:

* [HAL][hal]
* [HTTP Link header][1] - automatically registers as links regardless of format.
* [JSON:API][jsonapi] - Understands the `links` object and registers collection
  members as `item` relationships.
* [application/problem+json][problem] - Will extract useful information from
  the standard problem object and embed them in exception objects.


### Following links

One core tenet of building a good REST service, is that URIs should be
discovered, not hardcoded in an application. It's for this reason that the
emphasis in this library is _not_ on URIs (like most libraries) but on
relation-types (the `rel`) and links.

Generally when interacting with a REST service, you'll want to only hardcode
a single URI (a bookmark) and discover all the other APIs from there on.

For example, consider that there is a some API at `https://api.example.org/`.
This API has a link to an API for news articles (`rel="articleCollection"`),
which has a link for creating a new article (`rel="new"`). When `POST`ing on
that URI, the API returns `201 Created` along with a `Location` header pointing
to the new article. On this location, a new `rel="author"` appears
automatically, pointing to the person that created the article.

This is how that interaction might look like:

Coming soon.

### Embedded resources

Embedded resources are a HAL feature. In situations when you are modeling a
'collection' of resources, in HAL you should generally just create links to
all the items in the collection. However, if a client wants to fetch all these
items, this can result in a lot of HTTP requests. HAL uses `_embedded` to work
around this. Using `_embedded` a user can effectively tell the HAL client about
the links in the collection and immediately send along the contents of those
resources, thus avoiding the overhead.

Getting understands `_embedded` and completely abstracts them away. If you use
Getting with a HAL server, you can therefore completely ignore them.

For example, given a collection resource with many resources that hal the
relationshiptype `item`, you might use the following API:

Coming soon.

Given the last example, if the server did _not_ use embedding, it will result
in a HTTP GET request for every item in the collection.

If the server _did_ use embedding, there will only be 1 GET request.

A major advantage of this, is that it allows a server to be upgradable. Hot
paths might be optimized using embedding, and the client seamlessly adjusts
to the new information.

Further reading:

* [Further reading](https://evertpot.com/rest-embedding-hal-http2/).
* [Hypertext Cache Pattern in HAL spec](https://tools.ietf.org/html/draft-kelly-json-hal-08#section-8.3).


Automatically parsing problem+json
----------------------------------

If your server emits application/problem+json documents ([RFC7807][problem])
on HTTP errors, the library will automatically extract the information from
that object, and also provide a better exception message (if the title
property is provided).

[1]: https://github.com/badgateway/ketting
[2]: https://tools.ietf.org/html/rfc8288 "Web Linking"

[hal]: http://stateless.co/hal_specification.html "HAL - Hypertext Application Language"
[jsonapi]: https://jsonapi.org/
[problem]: https://tools.ietf.org/html/rfc7807
