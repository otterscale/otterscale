<script lang="ts">
	import * as Card from '$lib/components/ui/card';
	import * as Avatar from '$lib/components/ui/avatar';
	import Icon from '@iconify/svelte';
	import { Button, buttonVariants } from '../ui/button';
	import { Badge } from '../ui/badge';
	import { goto } from '$app/navigation';

	const benefits = ['4 Team member', '4 GB Storage', 'Upto 6 pages'];
</script>

<section class="container grid place-items-center gap-10 py-20 md:py-32 lg:grid-cols-2">
	<div class="space-y-6 text-center lg:text-start">
		<main class="text-5xl font-bold md:text-6xl">
			<h1 class="inline">
				Power Your
				<span
					class="inline bg-gradient-to-r from-[#F596D3] to-[#D247BF] bg-clip-text text-transparent"
				>
					All-in-One
				</span>{' '}
			</h1>
			{' '}
			Data Center with{' '}
			<h2 class="inline">
				<span
					class="inline bg-gradient-to-r from-[#61DAFB] via-[#1fc0f1] to-[#03a3d7] bg-clip-text text-transparent"
				>
					OtterScale
				</span>{' '}
			</h2>
		</main>

		<p class="text-md mx-auto text-muted-foreground md:w-10/12 lg:mx-0">
			Seamless, Scalable, and Future-Ready Solutions
		</p>

		<div class="space-x-2 space-y-4 md:space-x-4 md:space-y-0">
			<Button size="lg" href="/tutorial">
				<Icon icon="ph:cursor-click" />
				<p class="px-2">Get Started</p>
			</Button>
			<Button size="lg" variant="outline" target="_blank" href="https://github.com/openhdc/openhdc">
				<Icon icon="akar-icons:github-fill" />
				<p class="px-2">GitHub</p>
			</Button>
		</div>
	</div>

	<div class="z-10">
		<div class="relative hidden h-[500px] w-[705px] flex-row flex-wrap gap-8 lg:flex">
			<!-- Testimonial -->
			<Card.Root
				class="absolute -top-[15px] w-[340px] px-2 shadow-black/10 drop-shadow-xl hover:scale-105 dark:shadow-white/10"
			>
				<Card.Header class="flex flex-row gap-4">
					<Avatar.Root>
						<Avatar.Image src="https://github.com/iasthc.png" />
						<Avatar.Fallback>PT</Avatar.Fallback>
					</Avatar.Root>
					<Card.Title>
						Paul Tsai
						<Card.Description>@iasthc</Card.Description>
					</Card.Title>
				</Card.Header>
				<Card.Content class="font-medium italic">OtterScale: Integration, But Better!</Card.Content>
			</Card.Root>

			<!-- Team -->
			<Card.Root
				class="absolute right-[20px] top-4 flex w-80 flex-col items-center justify-center shadow-black/10 drop-shadow-xl hover:scale-105 dark:shadow-white/10"
			>
				<Card.Header class="flex items-center justify-center pb-2">
					<img
						src="https://upload.wikimedia.org/wikipedia/commons/6/61/Phison-logo.svg"
						alt="Phison"
					/>
					<Card.Description class="font-normal text-primary">Service Provider</Card.Description>
				</Card.Header>

				<Card.Content class="text-center">
					<p>A World Leader in NAND Controllers & Flash Storage Solutions</p>
				</Card.Content>

				<Card.Footer>
					<Button
						size="icon"
						variant="outline"
						class="border-none shadow-none"
						href="https://www.linkedin.com/company/phison-electronics"
						target="_blank"
					>
						<Icon icon="akar-icons:linkedin-fill" class="h-5 w-5" />
					</Button>
					<Button
						size="icon"
						variant="outline"
						class="border-none shadow-none"
						href="https://www.facebook.com/phisonUS"
						target="_blank"
					>
						<Icon icon="akar-icons:facebook-fill" class="h-5 w-5" />
					</Button>

					<Button
						size="icon"
						variant="outline"
						class="border-none shadow-none"
						href="https://www.youtube.com/c/PhisonElectronicsCorp"
						target="_blank"
					>
						<Icon icon="akar-icons:youtube-fill" class="h-5 w-5" />
					</Button>
				</Card.Footer>
			</Card.Root>

			<!-- Pricing -->
			<Card.Root
				class="absolute left-[50px] top-[145px] w-72 shadow-black/10 drop-shadow-xl hover:scale-105 dark:shadow-white/10"
			>
				<Card.Header>
					<Card.Title class="item-center flex justify-between">
						Free
						<Badge variant="secondary" class="text-sm text-primary">Most popular</Badge>
					</Card.Title>
					<div>
						<span class="text-3xl font-bold">$0</span>
						<span class="text-muted-foreground"> /month</span>
					</div>

					<Card.Description>
						Lorem ipsum dolor sit, amet ipsum consectetur adipisicing elit.
					</Card.Description>
				</Card.Header>

				<Card.Content>
					<Button class="w-full">Start Free Trial</Button>
				</Card.Content>

				<hr class="m-auto mb-4 w-4/5" />

				<Card.Footer class="flex">
					<div class="space-y-4">
						{#each benefits as benefit}
							<span class="flex">
								<Icon icon="lucide:check" class="text-green-500" />
								<h3 class="ml-2">{benefit}</h3>
							</span>
						{/each}
					</div>
				</Card.Footer>
			</Card.Root>

			<!-- Service -->
			<Card.Root
				class="absolute -right-[10px] bottom-[50px] w-[350px] shadow-black/10 drop-shadow-xl hover:scale-105 dark:shadow-white/10"
			>
				<Card.Header class="flex items-start justify-start gap-4 space-y-1 pb-6 md:flex-row">
					<div class="p-1">
						<Icon icon="logos:fabric" class="h-10 w-10" />
					</div>
					<div>
						<Card.Title>Unified Platform Orchestration</Card.Title>
						<Card.Description class="text-md mt-2">
							Integrate and manage with unparalleled ease.
						</Card.Description>
					</div>
				</Card.Header>
			</Card.Root>
		</div>
	</div>

	<div class="shadow">12334534534</div>
</section>

<style>
	/* HeroCards background shadow */
	.shadow {
		position: absolute;
		background: hsl(var(--primary) / 0%);
		border-radius: 48px;
		rotate: 48deg;
		top: 240px;
		z-index: 0;
		width: 480px;
		height: 480px;
		filter: blur(120px);
		will-change: transform;
		animation: shadow-slide infinite 4s linear alternate;
	}

	@keyframes shadow-slide {
		from {
			background: hsl(var(--primary) / 20%);
			right: 460px;
		}

		to {
			background: hsl(var(--primary) / 80%);
			right: 160px;
		}
	}

	@media (max-width: 1024px) {
		.shadow {
			top: 70px;
		}

		@keyframes shadow-slide {
			from {
				background: hsl(var(--primary) / 20%);
				right: 460px;
			}

			to {
				background: hsl(var(--primary) / 50%);
				right: 160px;
			}
		}
	}

	@media (max-width: 768px) {
		.shadow {
			top: 70px;
			width: 100px;
			height: 350px;
			filter: blur(60px);
		}

		@keyframes shadow-slide {
			from {
				background: hsl(var(--primary) / 20%);
				right: 280px;
			}

			to {
				background: hsl(var(--primary) / 30%);
				right: 100px;
			}
		}
	}
</style>
