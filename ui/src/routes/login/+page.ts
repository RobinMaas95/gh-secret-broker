import type { PageLoad } from "./$types";
import { error } from "@sveltejs/kit";

export const load: PageLoad = async ({ fetch, url }) => {
    try {
        const res = await fetch("/api/providers");
        if (!res.ok) {
            throw error(res.status, "Failed to load providers");
        }
        const providers = await res.json();
        const unauthorized = url.searchParams.has("unauthorized");
        return { providers, unauthorized };
    } catch (e) {
        if (e instanceof Response) throw e;
        throw error(500, "Failed to load providers");
    }
};
