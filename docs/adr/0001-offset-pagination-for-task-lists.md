# Offset pagination for the task list

Status: accepted

`GET /tasks` returns a user's tasks one page at a time. We paginate with `limit` and `offset` rather than a keyset cursor, so pagination state stays two readable integers that show up in logs and can be reproduced with curl.

## Context

The task list is filtered by user, optionally by an inclusive date range, and ordered newest first with `id` as a tiebreak. Tasks are append-only through the API: there is no delete, update, or patch endpoint, and `created_at` is set by the server.

## Considered options

- **Keyset cursor.** Stable under concurrent inserts and index-friendly, but state is an opaque token the client must carry and forward, malformed tokens need their own error path, and there is no page number or total.
- **`page` and `pageSize`.** Same SQL as offset, friendlier to a page-jump UI, but only meaningful with a total count, which means a filtered `COUNT(*)` on every request.
- **Hard-capped list.** Return at most N tasks with no way to fetch beyond them. Simple, but not pagination.
- **Offset and limit.** Chosen. Legible, ordinary SQL, no token lifecycle, and no count on the hot path.

## Consequences

- A task inserted between page requests shifts later rows down, so the last task of a page can reappear at the top of the next. This is accepted. A client that needs a frozen walk pins `endDate` to the first page's upper bound; no new parameter is required.
- There is no total count and no page number, so clients cannot build a "page 3 of 12" view.
- `OFFSET` costs more as the offset grows. Acceptable for a per-user list bounded by a date range.
- The choice rests on the list being append-only. If tasks become deletable or editable through the API, drift stops being insert-only and a keyset cursor deserves reconsideration.
