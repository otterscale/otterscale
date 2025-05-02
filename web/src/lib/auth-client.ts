import { createAuthClient } from "better-auth/svelte"
import { adminClient, organizationClient } from "better-auth/client/plugins"

export const authClient = createAuthClient({
    baseURL: import.meta.env.PUBLIC_URL,
    plugins: [
        adminClient(),
        organizationClient({
            teams: {
                enabled: true
            }
        })
    ]
})