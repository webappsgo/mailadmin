# Integrations

mailadmin administers mail software installed on its host, including Postfix, Dovecot, spam filtering, virus scanning, and supporting services. The Server Admin panel can inspect and control detected services through the host's available service manager, subject to operating-system privileges. OIDC, SAML, LDAP, and external webhook-delivery integrations are not currently exposed as working flows.

## API documentation and protocols

- Interactive OpenAPI explorer: `/server/docs/swagger`
- OpenAPI JSON: `/api/{api_version}/server/swagger`; `/api/swagger` is the direct current-version alias.
- GraphQL explorer: `/server/docs/graphql`
- Query-only GraphQL endpoint: `/api/{api_version}/server/graphql` (GET and POST); `/api/graphql` is its direct alias.
- Public discovery and PWA files: `/robots.txt`, `/sitemap.xml`, `/manifest.json`, and `/sw.js`.

GraphQL exposes `serverAbout` and `serverVersion`, plus generated query fields for listed, versioned REST `GET` routes. Those fields return JSON scalar data through the same route handlers, including their Bearer and tenant-ownership checks. It has no mutation root; write operations stay on REST.

OIDC provider metadata, webhooks, and native-app association endpoints are not currently enabled.

Mail-domain DNS guidance and external DNS updates require operator review. mailadmin must not assume authority to modify a third-party DNS zone.
