// TODO: Read from env variables by placing a correct config inside the eventual docker container
let host = location.host
host = host.split(":")[0]
export default {
    apiHost: host.concat(':8079')
}