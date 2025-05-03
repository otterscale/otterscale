import type { LayoutServerLoad } from './$types';

export const load = (async ({ locals, depends }) => {
    depends('app:user');
    return {
        user: locals.user
    };
}) satisfies LayoutServerLoad;