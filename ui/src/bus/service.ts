import {APIAdapter} from "@/bus/api";

export function init(orgSlug : string, boardSlug : string) {
    // TODO: Init setup is still ugly
    // instantiate APIAdapter
    new APIAdapter(orgSlug, boardSlug)
}