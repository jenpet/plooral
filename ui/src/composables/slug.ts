import slug from "slugify";

const slugLength = 60

export default (s:string) : string => {
    return slug(s.substring(0, slugLength-1), {
        strict: true
    })
}