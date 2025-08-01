<!-- <script>
    import { CarouselApi } from "@/components/ui/carousel";
    let api;

    const tweets = [
        {
            id: "1",
            user: {
                name: "Alex Morgan",
                handle: "@alexmdev",
                avatar:
                    "https://images.unsplash.com/photo-1568602471122-7832951cc4c5?q=80&w=1170&auto=format&fit=crop",
                verified: true,
            },
            content:
                "Just launched my new portfolio website built with Next.js and Tailwind CSS! Check it out and let me know what you think. #webdev #frontend #portfoliowebsite",
            date: "2h ago",
            metrics: {
                replies: 12,
                retweets: 24,
                likes: 89,
            },
            url: "#",
        },
    ];
</script>

<div class="bg-background py-16 md:py-24"></div>
    <div class="container mx-auto px-4 2xl:max-w-[1400px]">
        <div class="mb-10 text-center"></div>
            <div class="mb-3 inline-flex items-center gap-2">
                <Twitter class="h-5 w-5 text-sky-500" />
                <Badge
                    variant="outline"
                    class="border-sky-200 bg-sky-50 px-3 py-1 text-xs font-medium tracking-wide text-sky-600 uppercase dark:border-sky-500/20 dark:bg-sky-500/10 dark:text-sky-400"
                >
                    Twitter
                </Badge>
            </div>
            <h2 class="text-3xl font-bold tracking-tight sm:text-4xl md:text-5xl">
                Latest Tweets
            </h2>
            <p class="text-muted-foreground mx-auto mt-4 max-w-2xl">
                Stay updated with my latest thoughts and interactions in the tech community
            </p>
        </div>

        <div class="relative">
            <Carousel
                opts={{
                    align: "start",
                    loop: true,
                }}
                class="w-full"
                bind:api
            >
                <CarouselContent class="pt-2 pb-6">
                    {#each tweets as tweet (tweet.id)}
                        <CarouselItem class="md:basis-1/2 lg:basis-1/3"></CarouselItem>
                            <Card class="h-full border shadow-sm transition-all hover:shadow-md"></Card>
                                <CardContent class="flex h-full flex-col">
                                    <div class="mb-4 flex items-start justify-between">
                                        <div class="flex items-center gap-3">
                                            <div class="relative h-10 w-10 overflow-hidden rounded-full">
                                                <img
                                                    src={tweet.user.avatar}
                                                    alt={tweet.user.name}
                                                    class="object-cover w-full h-full"
                                                />
                                            </div>
                                            <div>
                                                <div class="flex items-center gap-1"></div>
                                                    <span class="text-sm font-semibold">
                                                        {tweet.user.name}
                                                    </span>
                                                    {#if tweet.user.verified}
                                                        <VerifiedIcon
                                                            class="h-3.5 w-3.5 text-sky-500"
                                                            fill="#0ea5e9"
                                                        />
                                                    {/if}
                                                </div>
                                                <div class="text-muted-foreground text-xs">
                                                    {tweet.user.handle}
                                                </div>
                                            </div>
                                        </div>
                                        <Twitter class="h-5 w-5 text-sky-500" />
                                    </div>

                                    <p class="mb-4 grow text-sm">{tweet.content}</p>

                                    <div class="border-t pt-4">
                                        <div class="text-muted-foreground mb-2 text-xs">
                                            {tweet.date}
                                        </div>
                                        <div class="flex items-center justify-between">
                                            <div class="flex items-center gap-4">
                                                <div class="text-muted-foreground flex items-center gap-1 hover:text-sky-500">
                                                    <MessageCircle class="h-4 w-4" />
                                                    <span class="text-xs">
                                                        {tweet.metrics.replies}
                                                    </span>
                                                </div>
                                                <div class="text-muted-foreground flex items-center gap-1 hover:text-green-500">
                                                    <Repeat2 class="h-4 w-4" />
                                                    <span class="text-xs">
                                                        {tweet.metrics.retweets}
                                                    </span>
                                                </div>
                                                <div class="text-muted-foreground flex items-center gap-1 hover:text-red-500"></div>
                                                    <Heart class="h-4 w-4" />
                                                    <span class="text-xs">
                                                        {tweet.metrics.likes}
                                                    </span>
                                                </div>
                                            </div>
                                            <a
                                                href={tweet.url}
                                                target="_blank"
                                                rel="noopener noreferrer"
                                                class="h-8 w-8 rounded-full hover:bg-sky-50 hover:text-sky-500 dark:hover:bg-sky-950/50"
                                            >
                                                <ExternalLink class="h-4 w-4" />
                                            </a>
                                        </div>
                                    </div>
                                </CardContent>
                            </Card>
                        </CarouselItem>
                    {/each}
                </CarouselContent>
            </Carousel>
        </div>

        <div class="flex flex-col items-center space-y-6">
            <div class="flex items-center justify-center gap-4">
                <button
                    class="btn variant-outline h-9 w-9 rounded-full p-0"
                    on:click={() => api?.scrollPrev()}
                >
                    <ChevronLeft class="h-5 w-5" />
                </button>
                <button
                    class="btn variant-outline h-9 w-9 rounded-full p-0"
                    on:click={() => api?.scrollNext()}
                >
                    <ChevronRight class="h-5 w-5" />
                </button>
            </div>

            <a
                href="#"
                target="_blank"
                rel="noopener noreferrer"
                class="btn variant-outline gap-2 border-sky-200 hover:bg-sky-50 hover:text-sky-600 dark:border-sky-500/20 dark:hover:bg-sky-950/30 dark:hover:text-sky-400"
            >
                <Twitter class="h-4 w-4" />
                <span>Follow me on Twitter</span>
            </a>
        </div>
    </div>
</div> -->
