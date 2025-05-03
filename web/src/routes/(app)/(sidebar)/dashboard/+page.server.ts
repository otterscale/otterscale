import { redirect } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';
import { setCallback } from '$lib/callback';

export const load: PageServerLoad = async ({ url, locals, depends }) => {
    depends('app:user');
    if (!locals.user) {
        redirect(302, setCallback('/login', url.pathname));
    }
    return {};
};