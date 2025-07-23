<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import * as Dialog from '$lib/components/ui/dialog';
	import { siteConfig } from '$lib/config/site';
	import { feedbackPath, githubPath, releasePath } from '$lib/path';
	import Icon from '@iconify/svelte';

	let { open = $bindable(false) }: { open: boolean } = $props();

	const links = [
		{ icon: 'ph:git-branch', text: siteConfig.version, url: releasePath },
		{ icon: 'ph:github-logo', text: 'GitHub', url: githubPath },
		{ icon: 'ph:paper-plane-tilt', text: 'Feedback', url: feedbackPath }
	];

	const openLink = (url: string) => window.open(url, '_blank');
</script>

<Dialog.Root bind:open>
	<Dialog.Content class="sm:max-w-[480px]">
		<Dialog.Header>
			<Dialog.Title>{siteConfig.title}</Dialog.Title>
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
