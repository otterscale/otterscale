// External libraries
import type { Handle } from '@sveltejs/kit';
import { sequence } from '@sveltejs/kit/hooks';
import { svelteKitHandler } from "better-auth/svelte-kit";

// Internal modules
import { auth } from "$lib/auth";
import { i18n } from '$lib/i18n';

// Define individual handlers
const handleParaglide: Handle = i18n.handle();
const handleAuth: Handle = async ({ event, resolve }) => {
    const { request } = event;
    const session = await auth.api.getSession({
        headers: request.headers
    });

    if (session && session.user) {
        event.locals.user = session.user;
    }

    return svelteKitHandler({ event, resolve, auth });
};

// Combine handlers using sequence
export const handle = sequence(handleParaglide, handleAuth);