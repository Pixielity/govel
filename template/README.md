<div align="center">

<img src="https://www.govel.com/logo.png?v=1.14.x" width="300" alt="Logo">

[![Doc](https://pkg.go.dev/badge/github.com/govel/framework)](https://pkg.go.dev/github.com/govel/framework)
[![Go](https://img.shields.io/github/go-mod/go-version/govel/framework)](https://go.dev/)
[![Release](https://img.shields.io/github/release/govel/framework.svg)](https://github.com/govel/framework/releases)
[![Test](https://github.com/govel/framework/actions/workflows/test.yml/badge.svg)](https://github.com/govel/framework/actions)
[![Report Card](https://goreportcard.com/badge/github.com/govel/framework)](https://goreportcard.com/report/github.com/govel/framework)
[![Codecov](https://codecov.io/gh/govel/framework/branch/master/graph/badge.svg)](https://codecov.io/gh/govel/framework)
![License](https://img.shields.io/github/license/govel/framework)

</div>

## About Govel

Govel is a web application framework with complete functions and good scalability. As a starting scaffolding to help
Gopher quickly build their own applications.

The framework style is consistent with [Laravel](https://github.com/laravel/laravel), let Php developer don't need to learn a
new framework, but also happy to play around Golang! In tribute to Laravel!

Welcome to star, PR and issues！

## Getting started

```
// Generate APP_KEY
go run . artisan key:generate

// Route
facades.Route().Get("/", userController.Show)

// ORM
facades.Orm().Query().With("Author").First(&user)

// Task Scheduling
facades.Schedule().Command("send:emails name").EveryMinute()

// Log
facades.Log().Debug(message)

// Cache
value := facades.Cache().Get("govel", "default")

// Queues
err := facades.Queue().Job(&jobs.Test{}, []queue.Arg{}).Dispatch()
```

## Documentation

Online documentation [https://www.govel.com](https://www.govel.com)

Example [https://github.com/govel/example](https://github.com/govel/example)

> To optimize the documentation, please submit a PR to the documentation
> repository [https://github.com/govel/docs](https://github.com/govel/docs)

## Main Function

|                                                                                        |                                                                 |                                                                          |                                                                       |                                                                                |
|----------------------------------------------------------------------------------------|-----------------------------------------------------------------|--------------------------------------------------------------------------|-----------------------------------------------------------------------|--------------------------------------------------------------------------------|
| [Config](https://www.govel.com/getting-started/configuration.html)                   | [Http](https://www.govel.com/the-basics/routing.html)         | [Authentication](https://www.govel.com/security/authentication.html)   | [Authorization](https://www.govel.com/security/authorization.html)  | [Orm](https://www.govel.com/orm/getting-started.html)                        |
| [Migrate](https://www.govel.com/database/migrations.html)                                 | [Logger](https://www.govel.com/the-basics/logging.html)       | [Cache](https://www.govel.com/digging-deeper/cache.html)               | [Grpc](https://www.govel.com/the-basics/grpc.html)                  | [Artisan Console](https://www.govel.com/digging-deeper/artisan-console.html) |
| [Task Scheduling](https://www.govel.com/digging-deeper/task-scheduling.html)         | [Queue](https://www.govel.com/digging-deeper/queues.html)     | [Event](https://www.govel.com/digging-deeper/event.html)               | [FileStorage](https://www.govel.com/digging-deeper/filesystem.html) | [Mail](https://www.govel.com/digging-deeper/mail.html)                       |
| [Validation](https://www.govel.com/the-basics/validation.html)                       | [Mock](https://www.govel.com/testing/mock.html)               | [Hash](https://www.govel.com/security/hashing.html)                    | [Crypt](https://www.govel.com/security/encryption.html)             | [Carbon](https://www.govel.com/digging-deeper/helpers.html)                  |
| [Package Development](https://www.govel.com/digging-deeper/package-development.html) | [Testing](https://www.govel.com/testing/getting-started.html) | [Localization](https://www.govel.com/digging-deeper/localization.html) | [Session](https://www.govel.com/the-basics/session.html)            |                                                                                |

## Roadmap

[For Detail](https://github.com/govel/govel/issues?q=is%3Aissue+is%3Aopen)

## Excellent Extend Packages

[For Detail](https://www.govel.com/getting-started/packages.html)

## Contributors

This project exists thanks to all the people who contribute, to participate in the contribution, please see [Contribution Guide](https://www.govel.com/getting-started/contributions.html).

<a href="https://github.com/hwbrzzl" target="_blank"><img src="https://avatars.githubusercontent.com/u/24771476?v=4" width="48" height="48"></a>
<a href="https://github.com/DevHaoZi" target="_blank"><img src="https://avatars.githubusercontent.com/u/115467771?v=4" width="48" height="48"></a>
<a href="https://github.com/kkumar-gcc" target="_blank"><img src="https://avatars.githubusercontent.com/u/84431594?v=4" width="48" height="48"></a>
<a href="https://github.com/almas-x" target="_blank"><img src="https://avatars.githubusercontent.com/u/9382335?v=4" width="48" height="48"></a>
<a href="https://github.com/merouanekhalili" target="_blank"><img src="https://avatars.githubusercontent.com/u/1122628?v=4" width="48" height="48"></a>
<a href="https://github.com/hongyukeji" target="_blank"><img src="https://avatars.githubusercontent.com/u/23145983?v=4" width="48" height="48"></a>
<a href="https://github.com/sidshrivastav" target="_blank"><img src="https://avatars.githubusercontent.com/u/28773690?v=4" width="48" height="48"></a>
<a href="https://github.com/Juneezee" target="_blank"><img src="https://avatars.githubusercontent.com/u/20135478?v=4" width="48" height="48"></a>
<a href="https://github.com/dragoonchang" target="_blank"><img src="https://avatars.githubusercontent.com/u/1432336?v=4" width="48" height="48"></a>
<a href="https://github.com/dhanusaputra" target="_blank"><img src="https://avatars.githubusercontent.com/u/35093673?v=4" width="48" height="48"></a>
<a href="https://github.com/mauri870" target="_blank"><img src="https://avatars.githubusercontent.com/u/10168637?v=4" width="48" height="48"></a>
<a href="https://github.com/Marian0" target="_blank"><img src="https://avatars.githubusercontent.com/u/624592?v=4" width="48" height="48"></a>
<a href="https://github.com/ahmed3mar" target="_blank"><img src="https://avatars.githubusercontent.com/u/12982325?v=4" width="48" height="48"></a>
<a href="https://github.com/flc1125" target="_blank"><img src="https://avatars.githubusercontent.com/u/14297703?v=4" width="48" height="48"></a>
<a href="https://github.com/zzpwestlife" target="_blank"><img src="https://avatars.githubusercontent.com/u/12382180?v=4" width="48" height="48"></a>
<a href="https://github.com/juantarrel" target="_blank"><img src="https://avatars.githubusercontent.com/u/7213379?v=4" width="48" height="48"></a>
<a href="https://github.com/Kamandlou" target="_blank"><img src="https://avatars.githubusercontent.com/u/77993374?v=4" width="48" height="48"></a>
<a href="https://github.com/livghit" target="_blank"><img src="https://avatars.githubusercontent.com/u/108449432?v=4" width="48" height="48"></a>
<a href="https://github.com/jeff87218" target="_blank"><img src="https://avatars.githubusercontent.com/u/29706585?v=4" width="48" height="48"></a>
<a href="https://github.com/shayan-yousefi" target="_blank"><img src="https://avatars.githubusercontent.com/u/19957980?v=4" width="48" height="48"></a>
<a href="https://github.com/zxdstyle" target="_blank"><img src="https://avatars.githubusercontent.com/u/38398954?v=4" width="48" height="48"></a>
<a href="https://github.com/milwad-dev" target="_blank"><img src="https://avatars.githubusercontent.com/u/98118400?v=4" width="48" height="48"></a>
<a href="https://github.com/mdanialr" target="_blank"><img src="https://avatars.githubusercontent.com/u/48054961?v=4" width="48" height="48"></a>
<a href="https://github.com/KlassnayaAfrodita" target="_blank"><img src="https://avatars.githubusercontent.com/u/113383200?v=4" width="48" height="48"></a>
<a href="https://github.com/YlanzinhoY" target="_blank"><img src="https://avatars.githubusercontent.com/u/102574758?v=4" width="48" height="48"></a>
<a href="https://github.com/gouguoyin" target="_blank"><img src="https://avatars.githubusercontent.com/u/13517412?v=4" width="48" height="48"></a>
<a href="https://github.com/dzham" target="_blank"><img src="https://avatars.githubusercontent.com/u/10853451?v=4" width="48" height="48"></a>
<a href="https://github.com/praem90" target="_blank"><img src="https://avatars.githubusercontent.com/u/6235720?v=4" width="48" height="48"></a>
<a href="https://github.com/vendion" target="_blank"><img src="https://avatars.githubusercontent.com/u/145018?v=4" width="48" height="48"></a>
<a href="https://github.com/tzsk" target="_blank"><img src="https://avatars.githubusercontent.com/u/13273787?v=4" width="48" height="48"></a>
<a href="https://github.com/ycb1986" target="_blank"><img src="https://avatars.githubusercontent.com/u/12908032?v=4" width="48" height="48"></a>
<a href="https://github.com/BadJacky" target="_blank"><img src="https://avatars.githubusercontent.com/u/113529280?v=4" width="48" height="48"></a>
<a href="https://github.com/NiteshSingh17" target="_blank"><img src="https://avatars.githubusercontent.com/u/79739154?v=4" width="48" height="48"></a>

## Sponsor

Better development of the project is inseparable from your support, reward us by [Open Collective](https://opencollective.com/govel).

<p align="left"><img src="https://www.govel.com/reward.png" width="200"></p>

## Group

Welcome more discussion in Discord.

[https://discord.gg/cFc5csczzS](https://discord.gg/cFc5csczzS)

## License

The Govel framework is open-sourced software licensed under the [MIT license](https://opensource.org/licenses/MIT).
