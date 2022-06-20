

<h1 align="center">Go APOD</h1>

<p align="center">
  <i>A CORS-enabled, no-auth wrapper to NASA's Astronomy Picture of the Day </i><br>
  <b>Public API: <a href="https://go-apod.herokuapp.com/">go-apod.herokuapp.com</a></b>
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

## Building and Running

#### Setup
If you haven't already done so, you'll need to [install Go Lang](https://go.dev/doc/install).

Then clone the repo `git clone https://github.com/Lissy93/go-apod.git && cd go-apod`


#### Build / Run
- Development - `go run main.go`
- Production - `go build -o bin/apod main.go`
- Testing - `go test`

#### Environmental Variables

- `NASA_API_KEY` (Required) - Your API Key, you can sign up for one at [api.nasa.gov](https://api.nasa.gov/)
- `PORT` (Optional) - The port to start the web server on, defaults to `8080`


## Deployment

[![Deploy to Heroku](https://www.herokucdn.com/deploy/button.svg)](https://heroku.com/deploy?template=https://github.com/Lissy93/go-apod)

## License

Licensed under [MIT](https://github.com/Lissy93/go-apod/blob/master/LICENSE), © [Alicia Sykes](https://aliciasykes.com) 2022
