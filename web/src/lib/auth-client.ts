import { ssoClient } from '@better-auth/sso/client';
import { createAuthClient } from 'better-auth/svelte';

import { env } from '$env/dynamic/public';

export const authClient = createAuthClient({
	baseURL: env.PUBLIC_URL,
	plugins: [ssoClient()],
});
