<script lang="ts">
	import { authClient } from '$lib/auth-client';
	import { Login } from './login';
	import { Favorite } from './favorite';
	import { Message } from './message';
	import { Preference } from './preference';

	const session = authClient.useSession();
</script>

<div class="flex justify-end space-x-2">
	<Preference />
	{#if $session.data}
		<Favorite />
		<Message />
	{:else}
		<Login />
	{/if}
</div>
