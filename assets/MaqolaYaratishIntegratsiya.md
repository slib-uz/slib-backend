## Overview

External integration endpoint for creating articles in journals via API. Requires BasicAuth.

The journal is determined from the authenticated client credentials (each integration client is bound to a journal in the system). If the client has no journal configured, the API returns 401 Unauthorized.

Authors (affiliations) are referenced by their ScienceID: the backend looks the author up in the `authors` table and uses its internal ID. If the author is not yet stored locally, it is fetched from ScienceID and saved automatically. You can also pre-resolve an author via the [Find author by ScienceID](#find-author-by-scienceid) endpoint.

Study field IDs for `study_fields_ids` can be obtained from the [List study fields](#list-study-fields) endpoint. Pass the returned `id` values in the array when creating an article.

Language ID for `language_id` can be obtained from the [List languages](#list-languages) endpoint. Pass the returned `id` value when creating an article.

Organizations (affiliations) are resolved by TIN: if an organization with the given TIN already exists in the system, its existing record is used; otherwise a new organization is created with the provided TIN and name. The article, organizations, and author affiliations are all created atomically in a single database transaction.

## Endpoint

```txt
POST https://api-service.slib.uz/api/integration/articles
```

## Authentication

\- Type: BasicAuth

\- Required: Yes

## Request Body

```json
{
    "name": {
        "uz": "Maqola nomi",
        "ru": "Название статьи",
        "en": "Article name"
    },
    "co_authors_count": 3,
    "study_fields_ids": [10, 20],
    "language_id": 1,
    "content_file": "base64_encoded_content",
    "tags": ["research", "science"],
    "doi": "10.1234/example.doi",
    "expert_conclusion_file": "base64_encoded_expert_conclusion",
    "annotation": {
        "uz": "Annotatsiya",
        "ru": "Аннотация",
        "en": "Annotation"
    },
    "affiliations": [
        {
            "science_id": "AAA-1111-2024",
            "organization_tin": "302345678",
            "organization_name": "Toshkent Davlat Universiteti",
            "position_name": "Professor"
        },
        {
            "science_id": "BBB-2222-2024",
            "organization_tin": "301122334",
            "organization_name": "Samarqand Davlat Universiteti",
            "position_name": "Dotsent"
        }
    ],
    "references": ["Reference 1", "Reference 2"],
    "published_date": "2026-01-15"
}
```

## Field Descriptions

<table><tr><td>Field</td><td>Type</td><td>Required</td><td>Description</td></tr><tr><td>name</td><td>map[string]string</td><td>Yes</td><td>Article name by language (locale: title)</td></tr><tr><td>co_authors_count</td><td>uint</td><td>Yes</td><td>Number of co-authors</td></tr><tr><td>study_fields_ids</td><td>[]uint</td><td>Yes</td><td>Study field IDs. Obtain valid IDs from <a href="#list-study-fields">List study fields</a></td></tr><tr><td>language_id</td><td>uint</td><td>Yes</td><td>Language identifier. Obtain a valid ID from <a href="#list-languages">List languages</a></td></tr><tr><td>content_file</td><td>string</td><td>Yes</td><td>Base64 encoded article content</td></tr><tr><td>tags</td><td>[]string</td><td>Yes</td><td>Article tags</td></tr><tr><td>doi</td><td>string</td><td>No</td><td>Digital Object Identifier</td></tr><tr><td>expert_conclusion_file</td><td>string</td><td>No</td><td>Base64 encoded expert conclusion</td></tr><tr><td>annotation</td><td>map[string]string</td><td>Yes</td><td>Annotation by language</td></tr><tr><td>affiliations</td><td>[]object</td><td>Yes</td><td>Author affiliations, created inline (see below)</td></tr><tr><td>references</td><td>[]string</td><td>Yes</td><td>Article references</td></tr><tr><td>published_date</td><td>string</td><td>No</td><td>Publication date, format YYYY-MM-DD (defaults to now)</td></tr></table>

### Affiliation object (`affiliations[]`)

Each affiliation binds an author to an organization for this article. The author is resolved by `science_id` (found locally, or fetched from ScienceID and saved). The organization is looked up by `organization_tin`; if not found, it is created with `organization_tin` + `organization_name`.

> The article's authors (co-authors) are taken from these affiliations - there is no separate co-author list. Each unique resolved author is attached to the article as a co-author.

<table><tr><td>Field</td><td>Type</td><td>Required</td><td>Description</td></tr><tr><td>science_id</td><td>string</td><td>Yes</td><td>Author ScienceID. Backend resolves it to the internal author (creating it locally from ScienceID if needed)</td></tr><tr><td>organization_tin</td><td>string</td><td>Yes</td><td>Organization TIN. Used to find or create the organization</td></tr><tr><td>organization_name</td><td>string</td><td>Yes</td><td>Organization name. Used only when a new organization is created</td></tr><tr><td>position_name</td><td>string</td><td>No</td><td>Author's position at the organization</td></tr></table>

## Response

Success (200)

```json
{
    "status": 200,
    "message": "Article created successfully"
}
```

## Error Responses

400 Bad Request

```json
{
    "status": 400,
    "message": "Invalid StudyField IDs: [99]"
}
```

```json
{
    "status": 400,
    "message": "affiliations[0]: science_id is required"
}
```

```json
{
    "status": 400,
    "message": "affiliations[0]: organization_tin is required"
}
```

```json
{
    "status": 400,
    "message": "affiliations[0]: organization_name is required"
}
```

## 401 Unauthorized

```json
{
    "status": 401,
    "message": "unauthorized"
}
```

Returned when Basic Auth fails, or when the authenticated client has no journal configured (`journal_id` is null).

## 404 Not Found

```json
{
    "status": 404,
    "message": "Resource not found"
}
```

## Example Usage

## cURL

```txt
curl -X POST https://api-service.slib.uz/api/integration/articles \
-H "Authorization: Basic base64(username:password)" \
-H "Content-Type: application/json" \
-d '{
    "name": {"en": "Sample Article"},
    "co_authors_count": 1,
    "study_fields_ids": [789],
    "language_id": 1,
    "content_file": "base64content...",
    "tags": ["research"],
    "annotation": {"en": "Sample annotation"},
    "affiliations": [
        {
            "science_id": "AAA-1111-2024",
            "organization_tin": "302345678",
            "organization_name": "Tashkent State University",
            "position_name": "Professor"
        }
    ],
    "references": ["Ref 1"],
    "published_date": "2026-01-15"
}'
```

## Go example

```go
import (
    "bytes"
    "encoding/json"
    "net/http"
)

payload := map[string]interface{}{
    "name":             map[string]string{"en": "Article Title"},
    "co_authors_count": 1,
    "study_fields_ids": []int{1},
    "language_id":      1,
    "content_file":     "base64content",
    "tags":             []string{"tag1"},
    "annotation":       map[string]string{"en": "Abstract"},
    "affiliations": []map[string]interface{}{
        {
            "science_id":        "AAA-1111-2024",
            "organization_tin":  "302345678",
            "organization_name": "Tashkent State University",
            "position_name":     "Professor",
        },
    },
    "references":     []string{"Ref"},
    "published_date": "2026-01-15",
}

body, _ := json.Marshal(payload)
req, _ := http.NewRequest("POST", "https://api-service.slib.uz/api/integration/articles", bytes.NewBuffer(body))
req.SetBasicAuth("username", "password")
req.Header.Set("Content-Type", "application/json")

client := &http.Client{}
resp, _ := client.Do(req)
defer resp.Body.Close()
```

---

## Find author by ScienceID

Helper endpoint to resolve (and, if necessary, create) an author from a ScienceID. Use it to look up an author before sending an article, or simply pass the `science_id` directly inside `affiliations` (the create endpoint resolves it automatically). If the author does not yet exist locally, it is fetched from ScienceID and saved, then returned.

### Endpoint

```txt
GET https://api-service.slib.uz/api/integration/authors/find-by-id/{science_id}
```

### Authentication

\- Type: BasicAuth

\- Required: Yes

### Path Parameters

<table><tr><td>Parameter</td><td>Type</td><td>Required</td><td>Description</td></tr><tr><td>science_id</td><td>string</td><td>Yes</td><td>Author ScienceID</td></tr></table>

### Response

Success (200)

```json
{
    "id": 123,
    "full_name": "Ism Familiya",
    "science_id": "AAA-1111-2024",
    "photo": null,
    "academic_title": null,
    "academic_degree": null,
    "orcid_id": null
}
```

The `id` field is the internal author ID that the affiliation is bound to.

### Error Responses

401 Unauthorized

```json
{
    "status": 401,
    "message": "Unauthorized"
}
```

404 Not Found

```json
{
    "status": 404,
    "message": "Resource not found"
}
```

### cURL

```txt
curl -X GET https://api-service.slib.uz/api/integration/authors/find-by-id/AAA-1111-2024 \
-H "Authorization: Basic base64(username:password)"
```

---

## List study fields

Helper endpoint to list available study fields and obtain their `id` values for `study_fields_ids` when creating an article. Filter by `journal_id` to get only the fields linked to the target journal.

### Endpoint

```txt
GET https://api-service.slib.uz/api/studyfield/list
```

### Authentication

\- Type: None

\- Required: No

### Query Parameters

<table><tr><td>Parameter</td><td>Type</td><td>Required</td><td>Description</td></tr><tr><td>journal_id</td><td>int</td><td>No</td><td>Filter study fields by journal ID</td></tr><tr><td>search</td><td>string</td><td>No</td><td>Search by name (uz, ru, en). Used when <code>journal_id</code> is not provided</td></tr></table>

### Response

Success (200)

```json
[
    {
        "id": 10,
        "name": {
            "uz": "Matematika",
            "ru": "Математика",
            "en": "Mathematics"
        },
        "parent_id": null,
        "code": 101,
        "parent": null,
        "children": []
    },
    {
        "id": 20,
        "name": {
            "uz": "Fizika",
            "ru": "Физика",
            "en": "Physics"
        },
        "parent_id": null,
        "code": 102,
        "parent": null,
        "children": []
    }
]
```

Use the `id` field from each item in the `study_fields_ids` array when creating an article.

### cURL

```txt
curl -X GET "https://api-service.slib.uz/api/studyfield/list?journal_id=123"
```

```txt
curl -X GET "https://api-service.slib.uz/api/studyfield/list?search=matematika"
```

---

## List languages

Helper endpoint to list available languages and obtain the `id` value for `language_id` when creating an article.

### Endpoint

```txt
GET https://api-service.slib.uz/api/language/list
```

### Authentication

\- Type: None

\- Required: No

### Response

Success (200)

```json
[
    {
        "id": 1,
        "name": {
            "uz": "O'zbek",
            "ru": "Узбекский",
            "en": "Uzbek"
        },
        "code": "uz"
    },
    {
        "id": 2,
        "name": {
            "uz": "Rus",
            "ru": "Русский",
            "en": "Russian"
        },
        "code": "ru"
    },
    {
        "id": 3,
        "name": {
            "uz": "Ingliz",
            "ru": "Английский",
            "en": "English"
        },
        "code": "en"
    }
]
```

Use the `id` field from the chosen language in the `language_id` field when creating an article.

### cURL

```txt
curl -X GET https://api-service.slib.uz/api/language/list
```
