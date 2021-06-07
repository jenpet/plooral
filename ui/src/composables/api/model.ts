export class Organization {
    slug? : string
    name? : string
    description? : string
    hidden? : boolean
    protected? : boolean
    tags? : string[]
    user_credentials? : SecurityCredentials
    owner_credentials? : SecurityCredentials
}

export class SecurityCredentials {
    password? : string
}

export class OrganizationCreationRequestBody {
    slug? : string
    name? : string
    description? : string
    protected? : boolean
    hidden? : boolean
    password? : string
    password_confirmation? : string
}

export class Board {
    slug? : string
    name? : string
    content? : any
}