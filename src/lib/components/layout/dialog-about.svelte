<script lang="ts">
	import BugIcon from '@lucide/svelte/icons/bug';
	import GitHubIcon from '@lucide/svelte/icons/github';

	import { version } from '$app/environment';
	import LogoImage from '$lib/assets/logo.svg';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import * as Dialog from '$lib/components/ui/dialog';
	import { m } from '$lib/paraglide/messages';

	let { open = $bindable(false) }: { open: boolean } = $props();
</script>

<Dialog.Root bind:open>
	<Dialog.Content>
		<Dialog.Header>
			<Dialog.Title>{m.about()}</Dialog.Title>
			<div class="flex w-full flex-col items-center justify-center">
				<Badge variant="secondary">{version}</Badge>
				<div class="flex w-full items-center justify-evenly">
					<img src={LogoImage} alt="logo" class="size-24" />
					<div class="flex gap-2">
						<h2 class="text-2xl leading-tight">
							{m.join_1()}
							<br />
							<span class="text-foreground/60"> {m.join_2()} </span>
						</h2>
					</div>
				</div>
				<div class="flex gap-2">
					<Button
						variant="outline"
						size="icon"
						href="https://github.com/otterscale"
						target="_blank"
					>
						<GitHubIcon />
					</Button>
					<Button
						variant="outline"
						size="icon"
						href="https://github.com/otterscale/otterscale/issues/new/choose"
						target="_blank"
					>
						<BugIcon />
					</Button>
				</div>
			</div>
		</Dialog.Header>
	</Dialog.Content>
</Dialog.Root>
