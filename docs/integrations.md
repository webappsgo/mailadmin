# Integrations

mailadmin is designed to administer mail software installed on its host, including Postfix, Dovecot, spam filtering, virus scanning, and supporting services. The currently implemented server-rendered panel and API do not expose complete service-control or third-party identity integration flows yet; do not treat configuration fields or roadmap scope as proof of a working integration.

## API documentation and protocols

- Interactive OpenAPI explorer: `/server/docs/swagger`
- OpenAPI JSON: `/api/{api_version}/server/swagger`; `/api/swagger` is the direct current-version alias.
- GraphQL explorer: `/server/docs/graphql`
- Query-only GraphQL endpoint: `/api/{api_version}/server/graphql` (GET and POST); `/api/graphql` is its direct alias.
- Public discovery files: `/robots.txt`, `/sitemap.xml`, and `/manifest.webmanifest`.

GraphQL exposes `serverAbout` and `serverVersion`, plus generated query fields for each listed, versioned REST `GET` route. Those fields return JSON scalar data and call the same route handlers, including their Bearer and tenant-ownership checks. It has no mutation root; write operations stay on REST.

OIDC provider metadata, webhooks, and native-app association endpoints are not currently enabled.

Mail-domain DNS guidance and external DNS updates require operator review. mailadmin must not assume authority to modify a third-party DNS zone.
