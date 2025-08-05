<script lang="ts">
	import Icon from '@iconify/svelte';
	import { page } from '$app/state';
	import { urlIcon, type Path } from '$lib/path';
	import { cephPaths, kubernetesPaths } from '$lib/routes';
	import { currentCeph, currentKubernetes } from '$lib/stores';

	interface Props {
		background: string;
		path: Path;
		description: string;
	}

	let { background, path, description }: Props = $props();

	const disabled = (url: string): boolean => {
		return (
			(!$currentCeph && cephPaths(page.params.scope).some((path) => path.url === url)) ||
			(!$currentKubernetes && kubernetesPaths(page.params.scope).some((path) => path.url === url))
		);
	};
</script>

<a
	href={path.url}
	class="group bg-card flex aspect-square flex-col overflow-hidden rounded-lg border shadow-sm transition-all hover:shadow-md {disabled(
		path.url
	)
		? 'pointer-events-none cursor-not-allowed opacity-50'
		: ''}"
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
</a>
