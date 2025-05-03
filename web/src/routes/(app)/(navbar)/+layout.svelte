<script lang="ts">
	import { SiteFooter, SiteHeader } from '$lib/components';
	import type { Snippet } from 'svelte';
	import type { PageData } from '../../$types';

	let { data, children }: { data: PageData; children: Snippet<[]> } = $props();
</script>

<div class="relative flex min-h-screen flex-col bg-background">
	<SiteHeader user={data.user} />
	<div class="flex-1">
		{@render children()}
	</div>
	<SiteFooter />
</div>
