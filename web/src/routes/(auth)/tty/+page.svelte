<script>
	import { page } from '$app/state';
	import { Terminal } from '$lib/components/applications/terminal';
</script>

<Terminal
	scope={page.url.searchParams.get('scope') ?? ''}
	facility={page.url.searchParams.get('facility') ?? ''}
	namespace={page.url.searchParams.get('namespace') ?? ''}
	podName={page.url.searchParams.get('pod') ?? ''}
	containerName={page.url.searchParams.get('container') ?? ''}
	command={(page.url.searchParams.get('command') ?? '/bin/sh').split(',')}
/>
