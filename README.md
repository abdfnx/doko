<p align="center">
  <img src="https://raw.githubusercontent.com/abdfnx/doko/main/.github/assets/logo.svg" height="120px" />
</p>

> 🐳 The docker you know but with TUI.

![preview](https://user-images.githubusercontent.com/64256993/148515590-dccda7c1-73ea-45c6-80b6-901633861fde.gif)

this app is inspired from [lazydocker](https://github.com/jesseduffield/lazydocker)

## Installation

### Using script

- Shell

```
curl -sL https://git.io/doko | bash
```

- PowerShell

```
iwr -useb https://git.io/doko-win | iex
```

#### or with [**resto**](https://github.com/abdfnx/resto)

```sh
# shell
resto install https://git.io/doko

# powershell
resto install https://git.io/doko-win
```

**then close and open your**

### Go package manager

```sh
go install github.com/abdfnx/doko@latest
```

### Via Docker

```bash
docker run -itv /var/run/docker.sock:/var/run/docker.sock dokocli/doko
docker run -itv /var/run/docker.sock:/var/run/docker.sock dokocli/doko <FLAGS>
docker run -itv /var/run/docker.sock:/var/run/docker.sock dokocli/doko <CMD>
```

> full container:

```bash
docker run -itv /var/run/docker.sock:/var/run/docker.sock dokocli/doko-full
```

## Usage

- Open Doko UI

```sh
doko
```

- With specific endpoint

```sh
doko --endpoint <DOCKER_ENDPOINT>
```

- Use another docker engine version

```sh
doko --engine "1.40"
```

- Log file path

```sh
doko --log-file /home/doko/my-log.log
```

## Flags

```
    --ca string          The path to the TLS CA (ca.pem)
-c, --cert string        The path to the TLS certificate (cert.pem)
-e, --endpoint string    The docker endpoint to use (default "unix:///var/run/docker.sock")
-g, --engine string      The docker engine version (default "1.41")
    --help               Help for doko
-k, --key string         The path to the TLS key (key.pem)
-l, --log-file string    The path to the log file
-o, --log-level string   The log level (default "info")
```

## Keybindings (Shortcuts)

| name             | mission                | key(s)                                              |
| ---------------- | ---------------------- | --------------------------------------------------- |
| all              | quit                   | <kbd>q</kbd>                                        |
| all              | change panel           | <kbd>Tab</kbd> or <kbd>Shift</kbd> + <kbd>Tab</kbd> |
| list panels      | next entry             | <kbd>j</kbd> or <kbd>↓</kbd>                        |
| list panels      | next page              | <kbd>Ctrl</kbd> or <kbd>f</kbd>                     |
| list panels      | previous entry         | <kbd>k</kbd> or <kbd>↑</kbd>                        |
| list panels      | previous page          | <kbd>Ctrl</kbd> or <kbd>b</kbd>                     |
| list panels      | scroll to top          | <kbd>g</kbd>                                        |
| list panels      | scroll to bottom       | <kbd>G</kbd>                                        |
| image list       | pull image             | <kbd>p</kbd>                                        |
| image list       | import image           | <kbd>i</kbd>                                        |
| image list       | save image             | <kbd>s</kbd>                                        |
| image list       | load image             | <kbd>Ctrl</kbd> + <kbd>l</kbd>                      |
| image list       | find images            | <kbd>f</kbd>                                        |
| image list       | delete image           | <kbd>d</kbd>                                        |
| image list       | filter image           | <kbd>/</kbd>                                        |
| image list       | create container       | <kbd>c</kbd>                                        |
| image list       | inspect image          | <kbd>Enter</kbd>                                    |
| image list       | refresh image list     | <kbd>Ctrl</kbd> + <kbd>r</kbd>                      |
| container list   | export container       | <kbd>e</kbd>                                        |
| container list   | commit container       | <kbd>c</kbd>                                        |
| container list   | filter image           | <kbd>/</kbd>                                        |
| container list   | exec container cmd     | <kbd>Ctrl</kbd> + <kbd>e</kbd>                      |
| container list   | start container        | <kbd>t</kbd>                                        |
| container list   | stop container         | <kbd>s</kbd>                                        |
| container list   | kill container         | <kbd>Ctrl</kbd> + <kbd>k</kbd>                      |
| container list   | delete container       | <kbd>d</kbd>                                        |
| container list   | inspect container      | <kbd>Enter</kbd>                                    |
| container list   | rename container       | <kbd>r</kbd>                                        |
| container list   | refresh container list | <kbd>Ctrl</kbd> + <kbd>r</kbd>                      |
| container logs   | show container logs    | <kbd>Ctrl</kbd> + <kbd>l</kbd>                      |
| volume list      | create volume          | <kbd>c</kbd>                                        |
| volume list      | delete volume          | <kbd>d</kbd>                                        |
| volume list      | filter volume          | <kbd>/</kbd>                                        |
| volume list      | inspect volume         | <kbd>Enter</kbd>                                    |
| volume list      | refresh volume list    | <kbd>Ctrl</kbd> + <kbd>r</kbd>                      |
| network list     | delete network         | <kbd>d</kbd>                                        |
| network list     | inspect network        | <kbd>Enter</kbd>                                    |
| network list     | filter network         | <kbd>/</kbd>                                        |
| pull image       | pull image             | <kbd>Enter</kbd>                                    |
| pull image       | close panel            | <kbd>Esc</kbd>                                      |
| create container | next input box         | <kbd>Tab</kbd>                                      |
| create container | previous input box     | <kbd>Shift</kbd> + <kbd>Tab</kbd>                   |
| detail           | cursor dwon            | <kbd>j</kbd>                                        |
| detail           | cursor up              | <kbd>k</kbd>                                        |
| detail           | next page              | <kbd>Ctrl</kbd> or <kbd>f</kbd>                     |
| detail           | previous page          | <kbd>Ctrl</kbd> or <kbd>b</kbd>                     |
| search images    | search image           | <kbd>Enter</kbd>                                    |
| search images    | close panel            | <kbd>Esc</kbd>                                      |
| search result    | next image             | <kbd>j</kbd>                                        |
| search result    | previous image         | <kbd>k</kbd>                                        |
| search result    | pull image             | <kbd>Enter</kbd>                                    |
| search result    | close panel            | <kbd>q</kbd>                                        |
| create volume    | close panel            | <kbd>Esc</kbd>                                      |
| create volume    | next input box         | <kbd>Tab</kbd>                                      |
| create volume    | previous input box     | <kbd>Shift</kbd> + <kbd>Tab</kbd>                   |
