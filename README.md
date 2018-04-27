# gocheck-sitemap
Monitors the URLs in sitemap continuously and exposes the HTTP statuses as a service

# Run in a docker container
```
docker run -p 3000:3000 -e SITEMAP=https://www.sitemaps.org/sitemap.xml redhatraptor/gocheck-sitemap
```

# See the http statuses
curl http://localhost:3000/

```
{
  "https://www.sitemaps.org/": 200,
  "https://www.sitemaps.org/faq.html": 200,
  "https://www.sitemaps.org/protocol.html": 200
}
```

# Future
- Page load time
- Page size
- Redirects
- First byte 
- more
