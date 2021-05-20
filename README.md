# PLOORAL
... is a multi-tenant collaboration whiteboard located in your browser **proof of concept**. 
I'd suggest to not look at the code since it is a shining counterexample of clean code which will hopefully be refactored someday.

# Development
Everything runs in docker (20.10+) containers during development. 

### 1. Install robo
[robo](https://github.com/tj/robo) is a simple YAML-based task runner written in Go.

```sh
$ curl -sf https://gobinaries.com/tj/robo | sh
```

### 2. Startup
```
$ robo kickstart
```
Spins up all necessary containers for local development with hot reload.

## Utilized Technologies
- [Vue 3](https://vuejs.org/) (including a sad first attempt using the new [Composition API](https://v3.vuejs.org/api/composition-api.html#composition-api))
- [Bootstrap](https://getbootstrap.com/)
- [Konva.js](https://konvajs.org/) builds the basic heart of the board itself
- [Golang](https://golang.org/)

# License
"PLOORAL" is Open Source software released under the MIT license.