<script lang="ts">
	import { page } from '$app/state';
	import { Separator } from '$lib/components/ui/separator';
	import * as Breadcrumb from '$lib/components/ui/breadcrumb';
	import * as Sidebar from '$lib/components/ui/sidebar';
	import { AppSidebar } from '$lib/components/sidebar';
	import { Avatar, FunctionBar, SearchBar } from '$lib/components/navbar';
	import { capitalizeFirstLetter } from 'better-auth';
	import type { PageData } from './$types';
	import type { Snippet } from 'svelte';
	import { i18n } from '$lib/i18n';

	let { data, children }: { data: PageData; children: Snippet<[]> } = $props();

	let breadcrumbs = $derived.by(() => {
		return i18n
			.route(page.url.pathname)
			.split('/')
			.filter((e) => e);
	});

	// TODO: recently visited pages
</script>

<Sidebar.Provider>
	<AppSidebar />
	<Sidebar.Inset>
		<header
			class="flex h-16 shrink-0 items-center gap-4 bg-sidebar transition-[width,height] ease-linear group-has-[[data-collapsible=icon]]/sidebar-wrapper:h-16 md:px-4"
		>
			<div class="flex items-center gap-2">
				<Sidebar.Trigger class="-ml-1" />
				<Separator orientation="vertical" class="mr-2 h-8" />
				<Breadcrumb.Root>
					<Breadcrumb.List>
						{#if breadcrumbs.length > 0}
							<Breadcrumb.Item class="hidden md:block">
								<Breadcrumb.Link href={`/${breadcrumbs[0]}`}
									>{capitalizeFirstLetter(breadcrumbs[0])}</Breadcrumb.Link
								>
							</Breadcrumb.Item>
							{#if breadcrumbs.length == 2}
								<Breadcrumb.Separator class="hidden md:block" />
								<Breadcrumb.Link href={`/${breadcrumbs.join('/')}`}
									>{capitalizeFirstLetter(breadcrumbs[breadcrumbs.length - 1])}</Breadcrumb.Link
								>
							{:else if breadcrumbs.length > 2}
								<Breadcrumb.Item class="hidden md:block">
									<span>...</span>
								</Breadcrumb.Item>
								<Breadcrumb.Separator class="hidden md:block" />
								<Breadcrumb.Item class="hidden md:block">
									<Breadcrumb.Link href={`/${breadcrumbs.slice(0, -1).join('/')}`}
										>{breadcrumbs[breadcrumbs.length - 2]}</Breadcrumb.Link
									>
								</Breadcrumb.Item>
								<Breadcrumb.Separator class="hidden md:block" />
								<Breadcrumb.Item class="hidden md:block">
									<Breadcrumb.Link href={`/${breadcrumbs.join('/')}`}
										>{breadcrumbs[breadcrumbs.length - 1]}</Breadcrumb.Link
									>
								</Breadcrumb.Item>
							{/if}
						{/if}
					</Breadcrumb.List>
				</Breadcrumb.Root>
			</div>
			<div class="flex gap-4 md:ml-auto md:gap-2 lg:gap-4">
				<SearchBar />
				<Separator orientation="vertical" />
				<FunctionBar user={data.user} />
			</div>
			<Avatar user={data.user} />
		</header>
		<div class="relative flex flex-col bg-background">
			<div class="flex-1 p-6">
				{@render children()}
			</div>
		</div>
	</Sidebar.Inset>
</Sidebar.Provider>
