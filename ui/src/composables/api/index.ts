import axios, {AxiosResponse} from 'axios'
import {Board, Organization} from "@/composables/api/model";
import config from "@/config";

const apiBasePath = 'api/v1'
const apiBaseUrl = `http://${config.apiHost}/${apiBasePath}`

const client = axios.create({
    baseURL: apiBaseUrl
})

export const getOrganizations = async () : Promise<Organization[]>=> {
    return performAPICall(client.get('/orgs'))
        // TODO: No clue on how to get a decent generics typing here to get rid of casting
        .then((resp ) : Organization[] => {
            return <Organization[]> resp.data
        })
}

export const getBoards = async (orgSlug : string) : Promise<Board[]> => {
    return performAPICall(client.get(`/orgs/${orgSlug}/boards`))
        .then((resp) : Board[] => {
            return <Board[]> resp.data
        })
}

export const getBoard = async (orgSlug : string, boardSlug : string) : Promise<Board> => {
    return performAPICall(client.get(`/orgs/${orgSlug}/boards/${boardSlug}`))
        .then((resp) : Board => {
            return <Board> resp.data
        })
}

export const storeBoard = async (orgSlug : string, boardSlug : string, content : string) : Promise<Board> => {
    return performAPICall(client.post(`/orgs/${orgSlug}/boards/${boardSlug}`, content))
        .then((resp) : Board => {
            return <Board> resp.data
        })
}


async function performAPICall<T>(call: Promise<AxiosResponse>): Promise<APIResponse<T>> {
    return new Promise<APIResponse<T>>((resolve, reject) => {
        call
            .then(response => {
                if (response.status >= 200 && response.status < 400) {
                    resolve(APIResponse.fromResponse(response))
                } else {
                    const msg = `API communication error. HTTP status code: ${response.status}. Body: ${JSON.stringify(response.data)}`
                    reject(new Error(msg))
                }})
            .catch((err) => {reject(err)}
            )
    })
}

export class APIResponse<T> {
    errors: string[]
    data: T
    constructor(errors: string[], data: T) {
        this.errors = errors
        this.data = data
    }
    static fromResponse<T>(res: AxiosResponse): APIResponse<T> {
        const errs = res.data.errors !== undefined ? res.data.errors : []
        return new APIResponse<T>(errs, res.data.data)
    }
}

export * from './model'