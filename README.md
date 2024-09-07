# GoRetro

![gofmt](https://github.com/justinawrey/goretro/actions/workflows/gofmt.yml/badge.svg)
![eslint](https://github.com/justinawrey/goretro/actions/workflows/eslint.yml/badge.svg)
![prettier](https://github.com/justinawrey/goretro/actions/workflows/prettier.yml/badge.svg)
![tsc](https://github.com/justinawrey/goretro/actions/workflows/tsc.yml/badge.svg)

GoRetro is a [NES](https://en.wikipedia.org/wiki/Nintendo_Entertainment_System) emulator written in Go. Using [Wails](https://wails.io/docs/introduction), I/O such as controller input, rendering, and audio are handled using cross-platform web apis accessed through the system default [WebView](https://en.wikipedia.org/wiki/WebView).

## Development

Run the emulator in development mode:

> [!IMPORTANT]
> You probably need to `chmod` the executable first.

```bash
cd internal/standalone
./dev.sh
```
