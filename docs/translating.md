# Translating mailadmin

mailadmin embeds six complete locale catalogs in its binaries: English (`en`),
Spanish (`es`), French (`fr`), Arabic (`ar`), German (`de`), and Italian
(`it`). English is the source catalog. Locale files live in
`src/common/i18n/locales/` and use UTF-8 JSON.

## Translation rules

- Keep every key and its nesting unchanged. Translate values only.
- Keep named placeholders exactly as written, including braces. For example,
  `{count}` and `{total}` must not be renamed or removed.
- Keep product names, protocol names, routes, config keys, filenames, code,
  URLs, and values users type unchanged.
- Write complete, natural sentences in a consistent, respectful register.
- Follow the language's punctuation and typography conventions. Use logical
  CSS properties so translated text can expand and right-to-left text can flow
  correctly.
- Provide every plural category listed for the locale in the registry. Arabic
  uses `zero`, `one`, `two`, `few`, `many`, and `other`; Italian uses `one`,
  `many`, and `other`; English, Spanish, French, and German use `one` and
  `other`. An optional `zero` form is allowed for other locales.
- Keep `meta.language`, `meta.name`, `meta.native_name`, `meta.direction`,
  `meta.version`, and `meta.plural_categories` in sync with the registry.

## Adding or updating a locale

1. Update the English source catalog first when adding a message.
2. Translate or update the same key in every supported locale in the same
   change. Remove a retired message from all locale files.
3. For a new language, add a complete locale JSON file and one registry entry
   in `src/common/i18n/i18n.go`. The registry owns the language code, native
   name, text direction, and plural categories.
4. Check the layout at 320, 768, and 1280 pixels. Check Arabic with right-to-left
   direction and verify that long words and longer translated strings do not
   overflow.
5. Run the repository test suite. The i18n tests check embedded locale files,
   key parity, placeholders, and plural categories.

The supported-language list is derived from the package registry. Do not add a
language to documentation or configuration unless its complete catalog is
embedded and registered. All binaries use the same catalogs through
`src/common/i18n`; no locale files are loaded from the runtime filesystem.

## Glossary

| English | Translation guidance |
|---------|----------------------|
| account | A user's or administrator's identity; keep distinct from organization |
| administrator | A person who manages the server; do not translate route names |
| backup | A verified archive used to restore application data |
| domain | An Internet domain name, not a database domain model |
| email | Use the locale's standard term for electronic mail |
| organization | A tenant that groups users and owns scoped data |
| password | A secret used to authenticate; do not confuse with an API token |
| session | An authenticated browser or API interaction period |
| token | An opaque credential; do not translate token prefixes or values |
| rotate | Replace a secret or key while honoring its configured overlap period |
| revoke | Invalidate a credential so it can no longer be used |

Translation corrections are welcome through the normal project contribution
process. The English catalog remains authoritative for key names and meaning.
