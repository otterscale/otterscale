<script lang="ts">
	import { PUBLIC_API_URL } from '$env/static/public';
	import { homePath, setupPath } from '$lib/path';
	import { breadcrumb } from '$lib/stores';
	import * as Code from '$lib/components/custom/code';
	import * as Terminal from '$lib/components/custom/terminal';
	import { m } from '$lib/paraglide/messages';
	import Icon from '@iconify/svelte';
	import { Button } from '$lib/components/ui/button';

	// Set breadcrumb navigation
	breadcrumb.set({ parents: [homePath], current: setupPath });

	const code = `sh -c "$(curl -fsSL https://raw.githubusercontent.com/otterscale/otterscale/main/scripts/install.sh" -- url=${PUBLIC_API_URL})`;

	let started = $state(true);
</script>

<div class="flex flex-col items-center justify-center">
	<h2 class="text-center text-3xl font-bold tracking-tight sm:text-4xl">{m.setup_environment()}</h2>

	{#if started}
		<Button
			class="text-muted-foreground mt-4 text-center text-lg"
			variant="ghost"
			size="lg"
			disabled
		>
			<Icon icon="ph:spinner-gap" class="size-6 animate-spin" />
			{m.setup_environment_installing()}
		</Button>

		<Terminal.Root class="m-6 max-w-6xl font-sans text-sm font-normal" delay={250}>
			<Terminal.TypingAnimation>&gt; jsrepo add ui/terminal</Terminal.TypingAnimation>
			<Terminal.Loading delay={1500}>
				{#snippet loadingMessage()}
					Fetching manifest
				{/snippet}
				{#snippet completeMessage()}
					<span class="text-green-500">
						<Icon icon="ph:check-bold" class="size-5" />
						✔ Retrieved blocks from github/ieedan/shadcn-svelte-extras
					</span>
				{/snippet}
			</Terminal.Loading>
			<Terminal.Loading delay={2750}>
				{#snippet loadingMessage()}
					Adding ui/terminal
				{/snippet}
				{#snippet completeMessage()}
					<span class="text-green-500">✔ Added ui/terminal</span>
				{/snippet}
			</Terminal.Loading>
			<Terminal.Loading delay={4000}>
				{#snippet loadingMessage()}
					Installing dependencies
				{/snippet}
				{#snippet completeMessage()}
					<span class="text-green-500">✔ Installed runed@^0.23.4</span>
				{/snippet}
			</Terminal.Loading>
			<Terminal.AnimatedSpan delay={5250} class="text-green-500">
				<span>✔ All done.</span>
			</Terminal.AnimatedSpan>
		</Terminal.Root>
	{:else}
		<Button
			class="text-muted-foreground mt-4 text-center text-lg"
			variant="ghost"
			size="lg"
			disabled
		>
			<Icon icon="ph:spinner-gap" class="size-6 animate-spin" />
			{m.setup_environment_waiting()}
		</Button>

		<div class="relative mx-auto flex w-full max-w-6xl flex-grow flex-col sm:mt-48">
			<p class="text-muted-foreground mt-4 text-center text-lg">
				{m.setup_environment_curl_description()}
			</p>

			<div class="w-full p-6">
				<Code.Root lang="bash" class="w-full" variant="secondary" {code} hideLines>
					<Code.CopyButton />
				</Code.Root>
			</div>

			<!-- Orange arrow decoration -->
			<div class="absolute end-0 top-0 hidden translate-x-20 -translate-y-12 md:block">
				<svg
					class="h-auto w-16 text-orange-500"
					width={121}
					height={135}
					viewBox="0 0 121 135"
					fill="none"
				>
					<path
						d="M5 16.4754C11.7688 27.4499 21.2452 57.3224 5 89.0164"
						stroke="currentColor"
						stroke-width={10}
						stroke-linecap="round"
					/>
					<path
						d="M33.6761 112.104C44.6984 98.1239 74.2618 57.6776 83.4821 5"
						stroke="currentColor"
						stroke-width={10}
						stroke-linecap="round"
					/>
					<path
						d="M50.5525 130C68.2064 127.495 110.731 117.541 116 78.0874"
						stroke="currentColor"
						stroke-width={10}
						stroke-linecap="round"
					/>
				</svg>
			</div>

			<!-- Cyan wavy decoration -->
			<div class="absolute start-0 bottom-0 hidden -translate-x-32 translate-y-10 md:block">
				<svg
					class="h-auto w-40 text-cyan-500"
					width={347}
					height={188}
					viewBox="0 0 347 188"
					fill="none"
				>
					<path
						d="M4 82.4591C54.7956 92.8751 30.9771 162.782 68.2065 181.385C112.642 203.59 127.943 78.57 122.161 25.5053C120.504 2.2376 93.4028 -8.11128 89.7468 25.5053C85.8633 61.2125 130.186 199.678 180.982 146.248L214.898 107.02C224.322 95.4118 242.9 79.2851 258.6 107.02C274.299 134.754 299.315 125.589 309.861 117.539L343 93.4426"
						stroke="currentColor"
						stroke-width={7}
						stroke-linecap="round"
					/>
				</svg>
			</div>
		</div>
	{/if}
</div>
