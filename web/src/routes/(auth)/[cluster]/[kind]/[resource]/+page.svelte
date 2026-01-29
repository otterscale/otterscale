<script lang="ts">
	import { page } from '$app/state';
	import { ResourceViewer } from '$lib/components/resource-viewer/index.js';

	const cluster = $derived(page.params.cluster ?? '');
	const resource = $derived(page.params.resource ?? '');
	const group = $derived(page.url.searchParams.get('group') ?? '');
	const version = $derived(page.url.searchParams.get('version') ?? '');
	const kind = $derived(page.url.searchParams.get('kind') ?? '');
	const namespace = $derived(page.url.searchParams.get('namespace') ?? '');
	const name = $derived(page.url.searchParams.get('name') ?? '');
</script>

{#key cluster + namespace + group + version + kind + resource + name}
	<ResourceViewer {cluster} {namespace} {group} {version} {kind} {resource} {name} />
{/key}
