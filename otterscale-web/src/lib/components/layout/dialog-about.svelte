<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import * as Dialog from '$lib/components/ui/dialog';
	import { siteConfig } from '$lib/config/site';
	import { m } from '$lib/paraglide/messages';
	import { staticPaths } from '$lib/path';
	import Icon from '@iconify/svelte';

	let { open = $bindable(false) }: { open: boolean } = $props();

	const links = [
		{ icon: 'ph:house-fill', text: m.homepage(), url: staticPaths.github.url },
		{ icon: 'ph:users-fill', text: m.contributors(), url: staticPaths.contributors.url },
		{ icon: 'ph:paper-plane-tilt-fill', text: m.feedback(), url: staticPaths.feedback.url }
	];

	const openLink = (url: string) => window.open(url, '_blank');
</script>

<Dialog.Root bind:open>
	<Dialog.Content class="p-12 sm:max-w-2xl">
		<div class="pointer-events-none absolute right-0 bottom-0 z-0 size-full">
			<svg
				aria-hidden="true"
				class="pointer-events-none absolute inset-0 h-full w-full fill-gray-400/30 stroke-gray-400/30"
				style="mask-image:radial-gradient(circle at 100% 100%, black 60%, transparent 100%);-webkit-mask-image:radial-gradient(circle at 100% 100%, black 60%, transparent 100%);opacity:0.4"
			>
				<defs>
					<pattern id=":S1:" width="40" height="40" patternUnits="userSpaceOnUse" x="-1" y="-1">
						<path d="M.5 40V.5H40" fill="none" stroke-dasharray="0"></path>
					</pattern>
				</defs>
				<rect width="100%" height="100%" stroke-width="0" fill="url(#:S1:)"></rect>
			</svg>
		</div>
		<Dialog.Header class="grid grid-cols-3 gap-4">
			<Dialog.Title class="flex flex-col items-center space-y-1">
				<Icon icon="fluent-emoji-flat:otter" class="size-16" />
				<span class="text-foreground text-xl font-semibold">{siteConfig.title}</span>
				<span class="text-muted-foreground text-sm">
					{import.meta.env.PACKAGE_VERSION}
				</span>
			</Dialog.Title>
			<Dialog.Description class="text-md col-span-2 flex flex-col justify-between gap-2 py-2">
				<span class="flex pt-2 font-light">{siteConfig.description}</span>
				<div class="flex justify-between">
					{#each links as { icon, text, url }}
						<Button variant="link" class="h-2 p-0 has-[>svg]:px-0" onclick={() => openLink(url)}>
							<Icon {icon} />
							{text}
						</Button>
					{/each}
				</div>
			</Dialog.Description>
		</Dialog.Header>
	</Dialog.Content>
</Dialog.Root>
