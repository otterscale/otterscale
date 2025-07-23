import { PUBLIC_BASE_URL } from "$env/static/public";
import { createAuthClient } from "better-auth/svelte";

export const authClient = createAuthClient({
    baseURL: PUBLIC_BASE_URL,
});
