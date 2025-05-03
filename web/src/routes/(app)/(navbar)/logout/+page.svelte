<script lang="ts">
	import { goto, invalidate } from '$app/navigation';
	import { i18n } from '$lib/i18n';
	import { authClient } from '$lib/auth-client';

	import { onMount } from 'svelte';

	let countdown = 3;

	onMount(() => {
		authClient.signOut({
			fetchOptions: {
				async onSuccess() {
					await invalidate('app:user');
				}
			}
		});

		const timer = setInterval(() => {
			countdown--;
			if (countdown <= 0) {
				clearInterval(timer);
				goto(i18n.resolveRoute('/'));
			}
		}, 1000);

		return () => clearInterval(timer);
	});

	var color = [
		'from-red-500 to-orange-300/80',
		'from-blue-500 to-purple-300/80',
		'from-green-500 to-teal-300/80',
		'from-pink-500 to-rose-300/80',
		'from-yellow-500 to-amber-300/80'
	][Math.floor(Math.random() * 5)];
</script>

<div class="flex min-h-[calc(100vh_-_theme(spacing.16))] w-full items-center justify-center">
	<div class="flex flex-col items-center gap-4">
		<span
			class="pointer-events-none flex justify-center whitespace-pre-wrap bg-gradient-to-b bg-clip-text text-center text-4xl font-semibold leading-none text-transparent {color}"
		>
			Good bye!
		</span>
		<span class="text-sm text-muted-foreground">Redirecting in {countdown} seconds...</span>
	</div>
</div>
