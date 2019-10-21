# Contributing guide

## Setting up development environment

The whole environment requires you to have:

1. `golang` 1.12+ for backend
2. `nodejs` 12+ (10+ should work as well) with `yarn` for frontend
3. [taskfile](https://taskfile.dev/#/installation) utility
4. [packr2](https://github.com/gobuffalo/packr/tree/master/v2)
5. [golangci-lint](https://github.com/golangci/golangci-lint) for linting

To start developing you should clone this repository in any directory you like:

```shell script
$ git clone git@github.com:ns3777k/mailcage.git
```

If you wanna work only on backend, you should build the frontend first with `task build:frontend` command, then
you can `task build:server` to rebuild the server. There's also `build:all` task to build all the binaries.

When you wanna change the frontend part, after building the server you should run it with
`--ui-assets-proxy-addr=127.0.0.1:3000` option. It means that instead of using pre-packed static files, server will
proxy it to `127.0.0.1:3000` which we're going set up next.

Frontend part is made with [create-react-app](https://github.com/facebook/create-react-app) and there's a command to
watch over files and rebuild them - `task watch:frontend`. By default, the watcher is available at `:3000`.

Finally, to see the ui, visit http://127.0.0.1:8025.

## Linting

Tasks `lint` and `lint:fix` let you check and fix the code.
