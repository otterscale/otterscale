import { redirect } from "@sveltejs/kit";
import { auth } from "$lib/auth";
import { dashboardPath, loginPath } from "$lib/path";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ request, url }) => {
	const session = await auth.api.getSession({
		headers: request.headers,
	});

	if (session) {
		redirect(302, dashboardPath);
	}

	redirect(302, `${loginPath}${url.search}`);
};
