# NASA-APOD

A go app that surfaces NASA's Astronomy Picture of the Day.

The web server exposes three routes:

	/       - Serves up static docs site as homepage
	/apod   - Fetches and returns JSON from APOD API
	/image  - Returns raw image from today's APOD img URL

