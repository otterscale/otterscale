<script lang="ts">
	import type { User } from 'better-auth';
	import { Login } from './login';
	import { Favorite } from './favorite';
	import { Message } from './message';
	import { Preference } from './preference';

	export let user: User;
</script>

<div class="flex justify-end space-x-2">
	<Preference />
	{#if user}
		<Favorite />
		<Message />
	{:else}
		<Login />
	{/if}
</div>
