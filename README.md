# GoRetro

![gofmt](https://github.com/justinawrey/goretro/actions/workflows/gofmt.yml/badge.svg)
![eslint](https://github.com/justinawrey/goretro/actions/workflows/eslint.yml/badge.svg)
![prettier](https://github.com/justinawrey/goretro/actions/workflows/prettier.yml/badge.svg)
![tsc](https://github.com/justinawrey/goretro/actions/workflows/tsc.yml/badge.svg)

GoRetro is a [NES](https://en.wikipedia.org/wiki/Nintendo_Entertainment_System) emulator written in Go. Using [Wails](https://wails.io/docs/introduction), I/O such as controller input, rendering, and audio are handled using cross-platform web apis accessed through the system default [WebView](https://en.wikipedia.org/wiki/WebView).

## Development

Run the emulator in development mode:

> [!IMPORTANT]
> Requirements:
>
> - [Go](https://go.dev)
> - [NPM](https://www.npmjs.com)
> - [Wails CLI](https://wails.io/docs/gettingstarted/installation#installing-wails) (make sure its in your `$PATH`)
> - [Platform specific dependencies](https://wails.io/docs/gettingstarted/installation#platform-specific-dependencies)
>
> You probably also need to `chmod` the `dev.sh` executable first.

```bash
cd internal/standalone
./dev.sh
```
