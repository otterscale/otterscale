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
		503: m.error_message_503(),
	};
</script>

{#if page.error}
	<div class="h-svh">
		<div class="m-auto flex h-full w-1/4 flex-col items-center justify-center gap-2">
			<h1 class="text-9xl leading-tight font-bold">{page.status}</h1>
			<span class="font-medium">{page.error.message}</span>
			<p class="text-muted-foreground line-clamp-2 text-center">
				{errorMessages[page.status]}
			</p>
			<p class="text-muted-foreground text-center"></p>
			<div class="mt-6 flex gap-4">
				<button
					class="focus-visible:ring-ring [&amp;_svg]:pointer-events-none [&amp;_svg]:size-4 [&amp;_svg]:shrink-0 border-input bg-background hover:bg-accent hover:text-accent-foreground inline-flex h-9 items-center justify-center gap-2 rounded-md border px-4 py-2 text-sm font-medium whitespace-nowrap shadow-xs transition-colors focus-visible:ring-1 focus-visible:outline-hidden disabled:pointer-events-none disabled:opacity-50"
					onclick={() => history.back()}>{m.go_back()}</button
				>
				<a
					class="focus-visible:ring-ring [&amp;_svg]:pointer-events-none [&amp;_svg]:size-4 [&amp;_svg]:shrink-0 bg-primary text-primary-foreground hover:bg-primary/90 inline-flex h-9 items-center justify-center gap-2 rounded-md px-4 py-2 text-sm font-medium whitespace-nowrap shadow-sm transition-colors focus-visible:ring-1 focus-visible:outline-hidden disabled:pointer-events-none disabled:opacity-50"
					href={resolve('/')}>{m.back_to_home()}</a
				>
			</div>
		</div>
	</div>
{/if}
