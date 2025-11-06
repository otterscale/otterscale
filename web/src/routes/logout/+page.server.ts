import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ locals, cookies }) => {
	cookies.delete('OS_TOKEN', { path: '/' });
	locals.user = null;
};
