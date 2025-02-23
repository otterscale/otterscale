<script lang="ts">
	import Separator from '$lib/components/ui/separator/separator.svelte';

	import { Avatar, FeatureBar, FeatureSidebar, FunctionBar, SearchBar } from '$lib/components/nav';
</script>

<header class="sticky top-0 z-50 flex h-16 items-center gap-4 border-b bg-header px-4 md:px-6">
	<FeatureBar />
	<FeatureSidebar />
	<div class="flex w-full items-center gap-4 md:ml-auto md:gap-2 lg:gap-4">
		<SearchBar />
		<!-- <Separator class="my-4" /> -->
		<FunctionBar />
		<Avatar />
	</div>
</header>
