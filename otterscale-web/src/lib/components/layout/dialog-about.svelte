<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import * as Dialog from '$lib/components/ui/dialog';
	import { siteConfig } from '$lib/config/site';
	import { m } from '$lib/paraglide/messages';
	import { staticPaths } from '$lib/path';
	import Icon from '@iconify/svelte';

	let { open = $bindable(false) }: { open: boolean } = $props();

	const links = [
		{ icon: 'ph:house', text: m.homepage(), url: staticPaths.github.url },
		{ icon: 'ph:users', text: m.contributors(), url: staticPaths.contributors.url },
		{ icon: 'ph:paper-plane-tilt', text: m.feedback(), url: staticPaths.feedback.url }
	];

	const openLink = (url: string) => window.open(url, '_blank');
</script>

<Dialog.Root bind:open>
	<Dialog.Content class="sm:max-w-[480px]">
		<Dialog.Header class="gap-4">
			<Dialog.Title class="flex items-center space-x-2">
				<span>{siteConfig.title}</span>
				<Icon icon="ph:git-commit-bold" class="text-muted-foreground" />
				<span class="text-muted-foreground text-sm font-semibold">
					{import.meta.env.PACKAGE_VERSION}
				</span>
			</Dialog.Title>
			<Dialog.Description>{siteConfig.description}</Dialog.Description>
		</Dialog.Header>
		<Dialog.Footer class="flex flex-col-reverse gap-2 sm:flex-row sm:justify-around">
			{#each links as { icon, text, url }}
				<Button variant="link" onclick={() => openLink(url)}>
					<Icon {icon} />
					{text}
				</Button>
			{/each}
		</Dialog.Footer>
	</Dialog.Content>
</Dialog.Root>
