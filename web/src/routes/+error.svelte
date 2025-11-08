<script lang="ts">
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import { m } from '$lib/paraglide/messages';

	const errorMessages: Record<number, string> = {
		400: m.error_message_400(),
		401: m.error_message_401(),
		403: m.error_message_403(),
		404: m.error_message_404(),
		500: m.error_message_500(),
		501: m.error_message_501(),
		502: m.error_message_502(),
		503: m.error_message_503()
	};
</script>

{#if page.error}
	<div class="h-svh">
		<div class="m-auto flex h-full w-1/4 flex-col items-center justify-center gap-2">
			<h1 class="text-9xl leading-tight font-bold">{page.status}</h1>
			<span class="font-medium">{page.error.message}</span>
			<p class="line-clamp-2 text-center text-muted-foreground">
				{errorMessages[page.status]}
			</p>
			<p class="text-center text-muted-foreground"></p>
			<div class="mt-6 flex gap-4">
				<button
					class="[&amp;_svg]:pointer-events-none [&amp;_svg]:size-4 [&amp;_svg]:shrink-0 inline-flex h-9 items-center justify-center gap-2 rounded-md border border-input bg-background px-4 py-2 text-sm font-medium whitespace-nowrap shadow-xs transition-colors hover:bg-accent hover:text-accent-foreground focus-visible:ring-1 focus-visible:ring-ring focus-visible:outline-hidden disabled:pointer-events-none disabled:opacity-50"
					onclick={() => history.back()}>{m.go_back()}</button
				>
				<a
					class="[&amp;_svg]:pointer-events-none [&amp;_svg]:size-4 [&amp;_svg]:shrink-0 inline-flex h-9 items-center justify-center gap-2 rounded-md bg-primary px-4 py-2 text-sm font-medium whitespace-nowrap text-primary-foreground shadow-sm transition-colors hover:bg-primary/90 focus-visible:ring-1 focus-visible:ring-ring focus-visible:outline-hidden disabled:pointer-events-none disabled:opacity-50"
					href={resolve('/')}>{m.back_to_home()}</a
				>
			</div>
		</div>
	</div>
{/if}
