<script lang="ts">
	import Icon from '@iconify/svelte';
	import { urlIcon, type Path } from '$lib/path';

	interface Props {
		background: string;
		path: Path;
		description: string;
	}

	let { background, path, description }: Props = $props();
</script>

<article
	class="group bg-card flex aspect-square flex-col overflow-hidden rounded-lg border shadow-sm transition-all hover:shadow-md"
>
	<header class="relative flex aspect-video">
		<div
			class="absolute inset-0 flex h-full items-center justify-center saturate-60 transition-opacity group-hover:opacity-80 dark:saturate-100 {background}"
		>
			<div
				class="bg-primary/30 mx-auto flex size-18 overflow-hidden rounded-full shadow-md group-hover:scale-105"
			>
				<Icon icon="{urlIcon(path.url)}-fill" class="text-muted m-auto size-8" />
			</div>
		</div>
	</header>

	<div class="flex flex-col items-center justify-center space-y-1 p-4 text-center">
		<h3 class="text-xl font-semibold">
			{path.title}
		</h3>
		<p class="text-muted-foreground text-sm">
			{description}
		</p>
	</div>
</article>
