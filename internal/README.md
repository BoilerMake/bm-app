# `internal/`
This directory should have packages that we don't want others to use.  This is stuff like "business logic", so http handlers and all that jazz. This doesn't me that other will not be able to see this directory, just that they probably shouldn't.

## Structure
### `app/`
Holds application logic
- HTTP handlers
- HTTP server
- Routing

### `pkg/`
Holds internal libraries
- None right now, there may never be any, idk
- But if there were, they should go in there
