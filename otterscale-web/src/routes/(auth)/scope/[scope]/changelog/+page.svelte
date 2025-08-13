<script lang="ts">
	import { page } from '$app/state';
	import * as Accordion from '$lib/components/ui/accordion';
	import * as Avatar from '$lib/components/ui/avatar';
	import { Badge } from '$lib/components/ui/badge';
	import * as Card from '$lib/components/ui/card';
	import * as HoverCard from '$lib/components/ui/hover-card';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import { formatTimeAgo } from '$lib/formatter';
	import { m } from '$lib/paraglide/messages';
	import { dynamicPaths } from '$lib/path';
	import { breadcrumb } from '$lib/stores';
	import changelogRead from '$lib/stores/changelog';
	import Icon from '@iconify/svelte';
	import type { PageData } from './$types';

	// Set breadcrumb navigation
	breadcrumb.set({ parents: [], current: dynamicPaths.changelog(page.params.scope) });

	interface Props {
		data: PageData;
	}

	let { data }: Props = $props();

	changelogRead.set(true);

	const CHANGE_TYPES = {
		feat: {
			title: m.changelog_feat(),
			icon: 'ph:check-circle',
			colors: {
				bg: 'bg-green-50 dark:bg-green-950/30',
				text: 'text-green-700 dark:text-green-400',
				border: 'border-green-200 dark:border-green-900'
			}
		},
		fix: {
			title: m.changelog_fix(),
			icon: 'ph:warning',
			colors: {
				bg: 'bg-red-50 dark:bg-red-950/30',
				text: 'text-red-700 dark:text-red-400',
				border: 'border-red-200 dark:border-red-900'
			}
		},
		perf: {
			title: m.changelog_perf(),
			icon: 'ph:arrow-circle-up-right',
			colors: {
				bg: 'bg-blue-50 dark:bg-blue-950/30',
				text: 'text-blue-700 dark:text-blue-400',
				border: 'border-blue-200 dark:border-blue-900'
			}
		},
		refactor: {
			title: m.changelog_refactor(),
			icon: 'ph:recycle',
			colors: {
				bg: 'bg-purple-50 dark:bg-purple-950/30',
				text: 'text-purple-700 dark:text-purple-400',
				border: 'border-purple-200 dark:border-purple-900'
			}
		},
		test: {
			title: m.changelog_test(),
			icon: 'ph:test-tube',
			colors: {
				bg: 'bg-yellow-50 dark:bg-yellow-950/30',
				text: 'text-yellow-700 dark:text-yellow-400',
				border: 'border-yellow-200 dark:border-yellow-900'
			}
		},
		style: {
			title: m.changelog_style(),
			icon: 'ph:palette',
			colors: {
				bg: 'bg-pink-50 dark:bg-pink-950/30',
				text: 'text-pink-700 dark:text-pink-400',
				border: 'border-pink-200 dark:border-pink-900'
			}
		},
		docs: {
			title: m.changelog_docs(),
			icon: 'ph:book',
			colors: {
				bg: 'bg-indigo-50 dark:bg-indigo-950/30',
				text: 'text-indigo-700 dark:text-indigo-400',
				border: 'border-indigo-200 dark:border-indigo-900'
			}
		},
		chore: {
			title: m.changelog_chore(),
			icon: 'ph:gear',
			colors: {
				bg: 'bg-gray-50 dark:bg-gray-950/30',
				text: 'text-gray-700 dark:text-gray-400',
				border: 'border-gray-200 dark:border-gray-900'
			}
		}
	} as const;

	function extractPRNumber(url: string): string {
		return `#${url.split('/').pop()}`;
	}

	function hasChanges(release: any, changeKey: keyof typeof CHANGE_TYPES): boolean {
		return release.changes[changeKey]?.length > 0;
	}
</script>

<h2 class="text-center text-3xl font-bold tracking-tight sm:text-4xl">{m.changelog()}</h2>

<p class="text-muted-foreground mt-4 text-center text-lg">
	{m.changelog_description()}
</p>

{#if data.error}
	<span class="mx-auto animate-pulse py-4 text-red-500">
		{data.error}
	</span>
{:else}
	<Accordion.Root
		type="single"
		class="mx-auto w-full py-10 sm:max-w-[80%] md:max-w-[65%]"
		value={data.releases[0].tag_name}
	>
		{#each data.releases as release}
			<Accordion.Item value={release.tag_name}>
				<Accordion.Trigger
					class="hover:bg-accent/50 items-center gap-2 hover:rounded-none hover:no-underline [[data-state=open]]:rounded-b-none [[data-state=open]]:border-b [&>svg:last-child]:hidden"
				>
					<div class="flex w-full flex-col flex-wrap px-4">
						<span
							class="text-foreground flex flex-1 items-center space-x-1 text-lg font-medium text-nowrap"
						>
							<Icon icon="ph:tag" class="size-4.5" />
							<span>{release.name}</span>
						</span>
					</div>

					<Tooltip.Provider>
						<Tooltip.Root>
							<Tooltip.Trigger
								class="text-muted-foreground flex space-x-1 text-xs tracking-tight text-nowrap"
							>
								<Icon icon="ph:git-pull-request" class="text-muted-foreground size-4" />
								<p>{formatTimeAgo(release.created_at)}</p>
							</Tooltip.Trigger>
							<Tooltip.Content>
								{release.created_at.toString()}
							</Tooltip.Content>
						</Tooltip.Root>
					</Tooltip.Provider>

					{#if release.prerelease}
						<Badge href={release.html_url} variant="outline">{m.changelog_prerelease()}</Badge>
					{/if}
					{#if release.latest}
						<Badge href={release.html_url}>{m.changelog_latest()}</Badge>
					{/if}

					<Icon
						icon="ph:caret-down"
						class="text-muted-foreground pointer-events-none mr-4 size-5 shrink-0 translate-y-0.5 transition-transform duration-200"
					/>
				</Accordion.Trigger>

				<Accordion.Content class="flex flex-col space-y-4 p-6 text-balance">
					{#each Object.entries(CHANGE_TYPES) as [key, config]}
						{#if hasChanges(release, key as keyof typeof CHANGE_TYPES)}
							<Card.Root class="gap-0 p-0 {config.colors.border}">
								<Card.Header class="gap-0 rounded-t-xl py-3 {config.colors.bg}">
									<Card.Title
										class="flex items-center space-x-1 font-semibold {config.colors.text}"
									>
										<Icon icon={config.icon} class="size-5" />
										<span>{config.title}</span>
									</Card.Title>
								</Card.Header>

								<Card.Content class="space-y-2 rounded-b-xl py-6">
									<ul class="ml-6 list-disc [&>li]:mt-2">
										{#each release.changes[key as keyof typeof CHANGE_TYPES] as item}
											<li>
												{item.description} by
												<HoverCard.Root>
													<HoverCard.Trigger
														href="https://github.com/{item.author}"
														target="_blank"
														rel="noreferrer noopener"
														class="hover:text-primary/80 rounded-sm font-semibold underline underline-offset-4 focus-visible:outline-2 focus-visible:outline-offset-8 focus-visible:outline-black"
													>
														{item.author}
													</HoverCard.Trigger>
													<HoverCard.Content class="w-70">
														<div class="flex items-center space-x-4">
															<Avatar.Root class="size-10">
																<Avatar.Image src="https://github.com/{item.author}.png" />
																<Avatar.Fallback>{item.author.slice(0, 2)}</Avatar.Fallback>
															</Avatar.Root>
															<div class="space-y-2">
																<span class="text-sm font-semibold">{item.author}</span>
																<span class="text-muted-foreground text-sm">
																	{data.usersMap[item.author].name}
																</span>
																{#if data.usersMap[item.author].company}
																	<div class="flex items-center space-x-1 pt-1">
																		<Icon icon="ph:building-office" class="size-4 opacity-80" />
																		<span class="text-xs">
																			{data.usersMap[item.author].company}
																		</span>
																	</div>
																{/if}
															</div>
														</div>
													</HoverCard.Content>
												</HoverCard.Root>
												in
												<a
													href={item.pull_request}
													target="_blank"
													rel="noreferrer noopener"
													class="text-blue-600 underline underline-offset-4 hover:text-blue-800 focus-visible:outline-2 focus-visible:outline-offset-8 focus-visible:outline-black"
												>
													{extractPRNumber(item.pull_request)}
												</a>
											</li>
										{/each}
									</ul>
								</Card.Content>
							</Card.Root>
						{/if}
					{/each}
				</Accordion.Content>
			</Accordion.Item>
		{/each}
	</Accordion.Root>
{/if}
