<script lang="ts">
	import { SiteFooter, SiteHeader } from '$lib/components';
	import type { PageData } from './$types';

	let { data }: { data: PageData } = $props();
</script>

<div class="relative flex min-h-screen flex-col bg-background">
	<SiteHeader user={data.user} />
	<div class="flex-1">Home</div>
	<SiteFooter />
</div>
