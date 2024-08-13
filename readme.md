
# run ren rennen 🏃

rennen (Dutch for running) runs a predefined set of commands simultaneously and has the option to trigger individual commands. it's designed to be simple, clean, and easy to use.

<img alt="Welcome to VHS" src="https://github.com/thejawker/rennen/blob/main/ren.gif" width="800" />

## installation

right now i really have only tested it on mac so dunno if it works on windows or linux

### mac
```bash
brew install thejawker/tappen/rennen
```

## usage

once you've downloaded the binary, you can run rennen with the following command `ren`:

```bash
ren init # to create a ren.json file
ren # to start rennen
```

## configuration
rennen requires a configuration file named `ren.json` in the same directory where the binary is run. this file should contain a `processes` array where each object represents a process that rennen should manage. each process object should have a `shortname`, `command`, and `description`.

here's an example of what the `ren.json` file could look like:

```json
{
  "processes": [
    {
      "shortname": "frontend",
      "command": "yarn start",
      "description": "starts the frontend server"
    },
    {
      "shortname": "server",
      "command": "php artisan serve",
      "description": "starts the laravel server"
    }
    //... add so many more processes
  ]
}
```

## development setup

if you want to contribute to rennen or run it in a development environment, you'll need to set up your environment first. here's how you can do it:

1. clone the repository:

```bash
git clone https://github.com/thejawker/rennen.git
```

2. navigate to the project directory:

```bash
cd rennen
```

3. install the dependencies:

```bash
go mod download
```

4. build the project:

```bash
make build
```

5. run the project:

```bash
make run
```

6. check the logs
since the only thing visible is the output of the commands, you can open a new terminal and run the following command:
```bash
make logs
```

## roadmap
read these things that i should do but probably won't [here](todo.md)

## creds
it's built with the [bubbletea](https://github.com/charmbracelet/bubbletea) framework, and is released using goreleaser.

## contributing
well, you really don't have to

hope you enjoy less pain in your life now you can ren! 


if you have any questions or run into any issues, please know that i might never respond.
