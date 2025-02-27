<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { SiteSidebar } from '$lib/components';
	import pb, { upsertVisit } from '$lib/pb';
	import { setCallback } from '$lib/callback';
	import { i18n } from '$lib/i18n';
	import { Button, buttonVariants } from '$lib/components/ui/button';
	import { Separator } from '$lib/components/ui/separator';
	import * as Breadcrumb from '$lib/components/ui/breadcrumb';
	import * as Drawer from '$lib/components/ui/drawer';
	import * as Sidebar from '$lib/components/ui/sidebar';
	import { AppSidebar } from '$lib/components/sidebar';
	import { Avatar, FunctionBar, SearchBar } from '$lib/components/navbar';

	let { children } = $props();

	let currentPage = i18n.route(page.url.pathname);

	onMount(async () => {
		await upsertVisit();
	});

	$effect(() => {
		if (currentPage !== i18n.route(page.url.pathname)) {
			(async () => await upsertVisit())();
			currentPage = i18n.route(page.url.pathname);
		}
		if (!pb.authStore.isValid) {
			goto(setCallback(i18n.resolveRoute('/login')));
		}
	});
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
						<Breadcrumb.Item class="hidden md:block">
							<Breadcrumb.Link href="#">Building Your Application</Breadcrumb.Link>
						</Breadcrumb.Item>
						<Breadcrumb.Separator class="hidden md:block" />
						<Breadcrumb.Item>
							<Breadcrumb.Page>Data Fetching</Breadcrumb.Page>
						</Breadcrumb.Item>
					</Breadcrumb.List>
				</Breadcrumb.Root>
			</div>
			<div class="flex gap-4 md:ml-auto md:gap-2 lg:gap-4">
				<SearchBar />
				<Separator orientation="vertical" />
				<FunctionBar />
			</div>
			<Avatar />
		</header>
		<div class="relative flex flex-col bg-background">
			<div class="flex-1">
				{@render children()}
			</div>
		</div>
	</Sidebar.Inset>
</Sidebar.Provider>
