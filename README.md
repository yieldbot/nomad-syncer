## nomad-syncer

[![Build Status][travis-image]][travis-url] [![GoDoc][godoc-image]][godoc-url] [![Release][release-image]][release-url]

An opinionated CLI for Nomad

### Installation

#### Binary releases

| Linux | OSX |
|:---:|:---:|
| [64bit][download-linux-amd64-url] | [64bit][download-osx-amd64-url] |

See all [releases](https://github.com/yieldbot/nomad-syncer/releases)

#### Building from source
```
go get github.com/yieldbot/nomad-syncer
cd $GOPATH/src/github.com/yieldbot/nomad-syncer
go build
```

### Usage

#### Help

```bash
./nomad-syncer -h
```
```
Usage: nomad-syncer [OPTIONS] COMMAND [arg...]

An opinionated CLI for Nomad

Options:
  --docker-pull : Pull Docker images before sync
  --nomad       : Nomad url (default "http://localhost:4646")
  --proxy       : Proxy url
  -h, --help    : Display usage
  -pp           : Pretty print for JSON output
  -v, --version : Display version information
  -vv           : Display extended version information

Commands:
  add           : Add a job
  del           : Delete a job
  get           : Get a job information
  jobs          : Retrieve jobs
  sync          : Sync jobs via a file or directory
```

#### Setting Nomad Url

Default Nomad url is `http://localhost:4646`. But also you can use `--nomad` argument on each
command or set ENV variable with following command

```bash
export NOMAD_URL=http://localhost:4646
```

#### Setting Proxy Url

You can use `--proxy` argument on each command or set ENV variable with following command

```bash
export NOMAD_SYNCER_PROXY_URL=http://localhost:8888
```


#### Getting jobs

```bash
./nomad-syncer jobs
```

#### Syncing jobs

Syncing a file
```bash
./nomad-syncer sync examples/job-1.json
```

Syncing a directory
```bash
./nomad-syncer sync examples/
```

Syncing with `--docker-pull`

If you have private Docker repository and don't want to add credentials into 
jobs (see https://www.nomadproject.io/docs/drivers/docker.html) then you can use `--docker-pull`.
It allows you to pull Docker images before syncing job files.

```bash
HOME=/apps/nomad /apps/nomad/nomad-syncer --docker-pull sync examples/
```

`HOME=/apps/nomad` points to `.dockercfg` 's directory. It can be a user or `/etc` directory.

#### Adding a job

```bash
./nomad-syncer add "$(cat examples/job-1.json)"
```

#### Getting a job

```bash
./nomad-syncer get job-1
```

#### Deleting a job

```bash
./nomad-syncer del job-1
```

### TODO

- [ ] Auto binary release
- [ ] Add tests

### License

Licensed under The MIT License (MIT)  
For the full copyright and license information, please view the LICENSE.txt file.

[travis-url]: https://travis-ci.org/yieldbot/nomad-syncer
[travis-image]: https://travis-ci.org/yieldbot/nomad-syncer.svg?branch=master

[godoc-url]: https://godoc.org/github.com/yieldbot/nomad-syncer
[godoc-image]: https://godoc.org/github.com/yieldbot/nomad-syncer?status.svg

[release-url]: https://github.com/yieldbot/nomad-syncer/releases/tag/v1.0.0
[release-image]: https://img.shields.io/badge/release-v1.0.0-blue.svg

[coverage-url]: https://coveralls.io/github/yieldbot/nomad-syncer?branch=master
[coverage-image]: https://coveralls.io/repos/yieldbot/nomad-syncer/badge.svg?branch=master&service=github)

[download-linux-amd64-url]: https://github.com/yieldbot/nomad-syncer/releases/download/v1.0.0/nomad-syncer-linux-amd64.zip
[download-osx-amd64-url]: https://github.com/yieldbot/nomad-syncer/releases/download/v1.0.0/nomad-syncer-osx-amd64.zip