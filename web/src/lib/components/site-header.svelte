<script lang="ts">
	import type { User } from 'better-auth';
	import Separator from '$lib/components/ui/separator/separator.svelte';
	import {
		Avatar,
		FeatureBar,
		FeatureSidebar,
		FunctionBar,
		SearchBar
	} from '$lib/components/navbar';
	import { features } from './features';

	export let user: User;
</script>

<header class="sticky top-0 z-50 flex h-16 items-center gap-4 border-b bg-sidebar px-4">
	<FeatureBar {features} />
	<FeatureSidebar {features} />
	<div class="flex gap-4 md:ml-auto md:gap-2 lg:gap-4">
		<SearchBar />
		<Separator orientation="vertical" />
		<FunctionBar {user} />
	</div>
	<Avatar {user} />
</header>
