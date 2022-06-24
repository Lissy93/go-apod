

<h1 align="center">Go APOD</h1>

<p align="center">
  <i>A CORS-enabled, no-auth wrapper to NASA's Astronomy Picture of the Day </i><br>
  <b>🌌 Public API: <a href="https://go-apod.herokuapp.com/">go-apod.herokuapp.com</a></b>
</p>


## API Usage


### `/apod`

> Returns full JSON info about todays picture

**Example**

```
GET https://go-apod.herokuapp.com/apod
```

**Response**

```json
{
  "date": "2022-06-20",
  "explanation": "There, just right of center, what is that? The surface of Mars keeps revealing new surprises with the recent discovery of finger-like rock spires. The small nearly-vertical rock outcrops were imaged last month by the robotic Curiosity rover on Mars. Although similar in size and shape to small snakes, the leading explanation for their origin is as conglomerations of small minerals left by water flowing through rock crevices. After these relatively dense minerals filled the crevices, they were left behind when the surrounding rock eroded away.  Famous rock outcrops on Earth with a similar origin are called hoodoos. NASA's Curiosity Rover continues to search for new signs of ancient water in Gale Crater on Mars, while also providing a geologic background important for future human exploration.   Explore Your Universe: Random APOD Generator",
  "hdurl": "https://apod.nasa.gov/apod/image/2206/MarsFingers_Curiosity_1338.jpg",
  "media_type": "image",
  "service_version": "v1",
  "title": "Rock Fingers on Mars",
  "url": "https://apod.nasa.gov/apod/image/2206/MarsFingers_Curiosity_960.jpg"
}
```

---

### `/image`

> Returns todays image

**Example**

```html
<img
  src="https://go-apod.herokuapp.com/image"
  alt="Astronomy Picture of the Day"
  width="350"
/>
```

**Response**

<img src="https://go-apod.herokuapp.com/image" alt="Astronomy Picture of the Day" width="350" />

---

## Deployment

> _Go-APOD can be self-hosted, either with Docker, via the 1-click Heroku deployment, or by running the executable directly._<br>
> A NASA API Key is required, which you can sign up for at [api.nasa.gov](https://api.nasa.gov/).

### Heroku

[![Deploy to Heroku](https://www.herokucdn.com/deploy/button.svg)](https://heroku.com/deploy?template=https://github.com/Lissy93/go-apod)

### Docker
A multi-arch container is available on DockerHub, under [`lissy93/apod`](https://hub.docker.com/r/lissy93/apod), or GHCR  as [`ghcr.io/lissy93/go-apod`](https://github.com/Lissy93/go-apod/pkgs/container/go-apod).<br> Or, use this [`docker-compose.yml`](https://github.com/Lissy93/go-apod/blob/master/docker-compose.yml) template, and just populate with your API key and run `docker compose up`.

```bash
docker run -p 8080:8080 -e NASA_API_KEY='XXX' -d lissy93/apod
```

### From Executable

Each release has pre-compiled binaries attached for Windows, Mac and Linux, which can be run directly.
From the [Releases Page](https://github.com/Lissy93/go-apod/releases), download and extract the version for your system, then execute it with: `NASA_API_KEY='XXX' ./go-apod`

### From Source

See the [Building Locally](#building-locally) section below

---


## Building Locally

> If you haven't already done so, you'll need to [install Go Lang](https://go.dev/doc/install).<br>
> Then clone the repo `git clone https://github.com/Lissy93/go-apod.git && cd go-apod`


### Commands
- Run Directly > `go run .`
- Compile App > `go build`
- Run Tests > `go test`

### Environmental Variables

- `NASA_API_KEY` (Required) - Your API Key, you can sign up for one at [api.nasa.gov](https://api.nasa.gov/)
- `PORT` (Optional) - The port to start the web server on, defaults to `8080`
- `CORS_ALLOWED_ORIGINS` (Optional) - List of origins which can use the API, defaults to `*` / all
- `NASA_BASE_URL` (Optional) - The base URL upstream GET requests, defaults to NASA's APOD API

---

## App

> A simple web app which displays today's image and info is also included. It can either be self-hosted or accessed below.

<p align="center">
  <a href="https://apod.as93.net">
    <img width="500" src="https://api.apiflash.com/v1/urltoimage?access_key=64850d88f6c645b3a144a493e725f769&url=https%3A%2F%2Fgo-apod.herokuapp.com%2F&format=webp&width=770&height=770&ttl=86400&response_type=image&wait_until=page_loaded&css=.link-wrapper%7Bdisplay%3Anone%3B%7D" />
  </a>
</p>

---

## License

Licensed under [MIT](https://github.com/Lissy93/go-apod/blob/master/LICENSE), © [Alicia Sykes](https://aliciasykes.com) 2022
