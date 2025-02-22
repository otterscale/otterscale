<script lang="ts">
	import pb from '$lib/pb';

	import { Login } from './login';
	import { Favorite } from './favorite';
	import { Message } from './message';
	import { Menu } from './menu';

	let isValid = pb.authStore.isValid;
	pb.authStore.onChange(() => {
		isValid = pb.authStore.isValid;
	});
</script>

<div class="flex justify-end space-x-2">
	<Menu />
	{#if isValid}
		<Favorite />
		<Message />
	{:else}
		<Login />
	{/if}
</div>
