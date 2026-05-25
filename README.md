# Makesweet server in golang

## What is this?

A server implementation of [paulfitz/makesweet](https://github.com/paulfitz/makesweet) meme gif generator using golang builded to be a microservice for a Discord bot.

This version uses a Python port of the original tool.

## Tools

This project uses the following tools

<div>
  <table>
    <tr>
      <th style="text-align:center">Golang</th>
      <th style="text-align:center">Docker</th>
      <th style="text-align:center">Gin</th>
      <th style="text-align:center">Swaggo</th>
    </tr>
    <tr>
      <td style="text-align: center"><a href="https://go.dev"><img src="https://go.dev/blog/go-brand/Go-Logo/SVG/Go-Logo_Blue.svg" height="90" alt="Golang" /></a></td>
      <td style="text-align: center"><a href="https://www.docker.com"><img src="https://uxwing.com/wp-content/themes/uxwing/download/brands-and-social-media/docker-icon.svg" height="90" alt="Docker" /></a></td>
      <td style="text-align: center"><a href="https://gin-gonic.com"><img src="https://raw.githubusercontent.com/gin-gonic/logo/master/color.png" height="90" alt="Gin" /></a></td>
      <td style="text-align: center"><a href="https://github.com/swaggo/swag"><img src="https://raw.githubusercontent.com/swaggo/swag/master/assets/swaggo.png" height="90" alt="Swaggo" /></a></td>
    </tr>
  </table>
</div>

## How to use?

1. Clone this repo to your desired location.
2. Go to the `makesweet-server` dir.
3. Copy the `.env.example` file to a file named `.env` and modify the env value (if you want).
4. Initialize the submodule with `git submodule update --init`
5. Use the `docker compose up` command to run the compose.
6. The server will run in http://localhost:8080/api.
7. You can use Swagger UI to test the API or check the endpoints in http://localhost:8080/api/docs/index.html
