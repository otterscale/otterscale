<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { setCallback } from '$lib/callback';
	import { i18n } from '$lib/i18n';
	import { Separator } from '$lib/components/ui/separator';
	import * as Breadcrumb from '$lib/components/ui/breadcrumb';
	import * as Sidebar from '$lib/components/ui/sidebar';
	import { AppSidebar } from '$lib/components/sidebar';
	import { Avatar, FunctionBar, SearchBar } from '$lib/components/navbar';
	import { capitalizeFirstLetter } from 'better-auth';
	import { authClient } from '$lib/auth-client';

	let { children } = $props();

	let breadcrumbs = $derived.by(() => {
		return page.url.pathname.split('/').filter((e) => e);
	});

	// TODO: recently visited pages

	$effect(() => {
		const session = authClient.useSession();
		if (!!session) {
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
