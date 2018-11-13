# new-backend

## Why dependency inject into server?
For one it makes mock testing things much easier, because you can just create your own server instance and test that with more granularity.  See [How I write Go HTTP services after seven years](https://medium.com/statuscode/how-i-write-go-http-services-after-seven-years-37c208122831) for more info on that.

## Current structure of backend
- One route ("/") sends the entire react app (~1.5 MB!!)
- The domains you see in the browser are not always requests sent to our backend
	- They're just react telling you lies, and hitting the api
	- e.g.:
		- "/login" is not registered at all in our current backend
		- But, that page is in the js from the react app
		- When pressing login, react sends a POST request with form info to this URL:
			- `api.boilermake.org/v1/users/login`
- There's some pros and cons things about this
	- Good:
		- More uniform request handling
			- Everything goes through the api, all render is done client side
		- Only make one http request for static content
		- Requests after will generally be small
		- Moves complexity from backend to frontend (maybe neutral?)
	- Bad:
		- Huge initial load times
		- Overall will typically send more over the wire
			- Depends on number of pages hit, but on average this is probably true
		- Puts more work on the client by making them render and create the DOM
			- Can be especially slow on mobile/lower power devices
			- and drain batter quicker
		- Worse SEO
			- Google can't index dynamic pages as well as static
			- Kinda... It may be better now than it was before
		- Requires writing a lot more js
			- Again may be neutral or even a positive, but I'd rather write go

## Thoughts on new structure
- Static pages should be static
	- Made with go templates or some server side render thing
	- Make sure these are being cached
- API should be used in a few select spots
	- Announcements
	- Some exec pages (For live updates)
	- Maybe application (Again, for live updates)
- Each of those places can have it's own front end app/framework thing
	- ofc with lots of sharing between them
	- test

Or maybe not... This is something we should probs discuss as a team

## Ideal way to structure handlers:
- Validate input
- Query/update some data
- Return a response
