<script lang="ts">
	import Icon from '@iconify/svelte';

	import LogoImage from '$lib/assets/logo.png';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import * as Dialog from '$lib/components/ui/dialog';
	import { siteConfig } from '$lib/config/site';
	import { m } from '$lib/paraglide/messages';

	let { open = $bindable(false) }: { open: boolean } = $props();

	const links = [
		{
			icon: 'simple-icons:github',
			url: 'https://github.com/otterscale/otterscale'
		},
		{
			icon: 'ph:paper-plane-tilt-fill',
			url: 'https://github.com/otterscale/otterscale/issues/new/choose'
		}
	];
</script>

<Dialog.Root bind:open>
	<Dialog.Content class="sm:max-w-xl">
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
		<Dialog.Header class="flex-row items-center justify-evenly gap-4">
			<Dialog.Title class="flex flex-col items-center gap-2 font-medium">
				<img src={LogoImage} alt="logo" class="relative bottom-0 -mt-4 size-24" />
				<div class="-mt-5 flex flex-col gap-1">
					<span class=" text-xl">{siteConfig.title}</span>
					<Badge variant="outline">{import.meta.env.PACKAGE_VERSION}</Badge>
				</div>
			</Dialog.Title>
			<Dialog.Description class="flex flex-col justify-between gap-4 py-2">
				<h2 class="text-center text-3xl">
					{m.join_1()}
					<br />
					<span class="text-muted-foreground/80"> {m.join_2()} </span>
				</h2>
				<div class="flex justify-center gap-4 text-primary">
					{#each links as { icon, url }}
						<Button variant="outline" size="icon" href={url} target="_blank">
							<Icon {icon} />
						</Button>
					{/each}
				</div>
			</Dialog.Description>
		</Dialog.Header>
	</Dialog.Content>
</Dialog.Root>
