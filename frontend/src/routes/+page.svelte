<script lang="ts">
	import { getListings } from '$lib/remote/listings.remote';
	import { getCategories } from '$lib/remote/categories.remote';
	import ListingCard from '$lib/components/ListingCard.svelte';
	import CategoryIcon from '$lib/components/CategoryIcon.svelte';
	import Seo from '$lib/components/Seo.svelte';
	import { page } from '$app/state';

	const listings = getListings({ limit: 12, sort: 'newest' });
	const featured = getListings({ featured: true, limit: 4 });
	const categories = getCategories();
</script>

<Seo
	title="Despacha Aí — Compra e venda em Benguela e toda Angola"
	description="Classificados grátis de Benguela, Lobito, Catumbela e toda Angola. Compra e vende telemóveis, carros, casas, roupa e muito mais. Publica grátis e negoceia pelo WhatsApp."
	jsonLd={{
		'@context': 'https://schema.org',
		'@type': 'WebSite',
		name: 'Despacha Aí',
		alternateName: 'Despacha Aí Angola',
		url: page.url.origin,
		description: 'Classificados grátis de Benguela e toda Angola',
		potentialAction: {
			'@type': 'SearchAction',
			target: {
				'@type': 'EntryPoint',
				urlTemplate: `${page.url.origin}/search?search={search_term_string}`
			},
			'query-input': 'required name=search_term_string'
		}
	}}
/>

<!-- Hero -->
<section
	class="relative overflow-hidden rounded-3xl bg-linear-to-br from-brand-900 via-brand-700 to-brand-500 text-white"
>
	<!-- decorative shapes -->
	<div
		class="pointer-events-none absolute -top-24 -right-24 h-80 w-80 rounded-full bg-white/10 blur-2xl"
		aria-hidden="true"
	></div>
	<div
		class="pointer-events-none absolute -bottom-32 left-1/3 h-72 w-72 rounded-full bg-brand-300/20 blur-3xl"
		aria-hidden="true"
	></div>
	<svg class="pointer-events-none absolute inset-0 h-full w-full opacity-[0.07]" aria-hidden="true">
		<defs>
			<pattern id="hero-dots" width="28" height="28" patternUnits="userSpaceOnUse">
				<circle cx="2" cy="2" r="2" fill="white" />
			</pattern>
		</defs>
		<rect width="100%" height="100%" fill="url(#hero-dots)" />
	</svg>

	<div class="relative grid items-center gap-8 px-6 py-7 sm:px-10 sm:py-9 lg:grid-cols-[1.2fr_1fr]">
		<div>
			<span
				class="inline-flex items-center gap-2 rounded-full bg-white/15 px-3.5 py-1 text-xs font-semibold backdrop-blur sm:text-sm"
			>
				<span class="relative flex h-2 w-2">
					<span
						class="absolute inline-flex h-2 w-2 animate-ping rounded-full bg-brand-200 opacity-75"
					></span>
					<span class="relative inline-flex h-2 w-2 rounded-full bg-brand-200"></span>
				</span>
				Anuncia de graça — sem taxas, sem comissões
			</span>

			<h1 class="mt-3 text-2xl leading-tight font-extrabold tracking-tight sm:text-3xl">
				Compra e vende em <span class="text-brand-200">Benguela e toda Angola</span>
			</h1>
			<p class="mt-2 max-w-lg text-sm text-brand-100">
				O mercado online de Benguela, Lobito e Catumbela. Publica o teu anúncio em menos de 2
				minutos e negoceia directamente pelo WhatsApp.
			</p>

			<form action="/search" class="mt-4 flex max-w-xl gap-2" role="search">
				<input
					type="search"
					name="search"
					placeholder="O que procuras hoje?"
					class="h-11 flex-1 rounded-full border-transparent bg-white/95 px-5 text-sm text-neutral-900 shadow-lg shadow-brand-950/20 placeholder:text-neutral-400 focus:border-white focus:ring-white"
				/>
				<button
					type="submit"
					class="h-11 shrink-0 rounded-full bg-neutral-900 px-5 text-sm font-semibold text-white shadow-lg shadow-brand-950/20 transition hover:bg-neutral-800"
					>Pesquisar</button
				>
			</form>

			<div class="mt-4 flex flex-wrap items-center gap-3">
				<a
					href="/anunciar"
					class="inline-flex h-10 items-center gap-2 rounded-full bg-white px-5 text-sm font-bold text-brand-800 shadow-lg shadow-brand-950/20 transition hover:bg-brand-50"
				>
					<svg class="h-4 w-4" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
						<path
							d="M10.75 4.75a.75.75 0 0 0-1.5 0v4.5h-4.5a.75.75 0 0 0 0 1.5h4.5v4.5a.75.75 0 0 0 1.5 0v-4.5h4.5a.75.75 0 0 0 0-1.5h-4.5v-4.5Z"
						/>
					</svg>
					Publicar anúncio grátis
				</a>
				<div class="flex items-center gap-4 text-sm text-brand-100">
					<span class="flex items-center gap-1.5">
						<svg class="h-4 w-4 text-brand-200" viewBox="0 0 20 20" fill="currentColor">
							<path
								fill-rule="evenodd"
								d="M16.704 4.153a.75.75 0 0 1 .143 1.052l-8 10.5a.75.75 0 0 1-1.127.075l-4.5-4.5a.75.75 0 0 1 1.06-1.06l3.894 3.893 7.48-9.817a.75.75 0 0 1 1.05-.143Z"
								clip-rule="evenodd"
							/>
						</svg>
						100% grátis
					</span>
					<span class="flex items-center gap-1.5">
						<svg class="h-4 w-4 text-brand-200" viewBox="0 0 20 20" fill="currentColor">
							<path
								fill-rule="evenodd"
								d="M16.704 4.153a.75.75 0 0 1 .143 1.052l-8 10.5a.75.75 0 0 1-1.127.075l-4.5-4.5a.75.75 0 0 1 1.06-1.06l3.894 3.893 7.48-9.817a.75.75 0 0 1 1.05-.143Z"
								clip-rule="evenodd"
							/>
						</svg>
						18 províncias
					</span>
				</div>
			</div>
		</div>

		<!-- hero illustration: phone mockup showing the app -->
		<div class="relative mx-auto hidden w-full max-w-[210px] xl:block" aria-hidden="true">
			<!-- device frame -->
			<div
				class="rotate-3 rounded-[2.2rem] border-[6px] border-neutral-900 bg-neutral-900 shadow-2xl shadow-brand-950/50"
			>
				<div class="relative overflow-hidden rounded-[1.85rem] bg-neutral-50">
					<!-- notch -->
					<div
						class="absolute top-1.5 left-1/2 z-10 h-4 w-16 -translate-x-1/2 rounded-full bg-neutral-900"
					></div>

					<!-- app screen -->
					<div class="px-3 pt-8 pb-3">
						<!-- app header -->
						<div class="flex items-center justify-between">
							<p class="text-[11px] font-extrabold text-brand-700">
								Despacha <span class="text-brand-500">Aí</span>
							</p>
							<span class="h-4 w-4 rounded-full bg-brand-100"></span>
						</div>

						<!-- search pill -->
						<div class="mt-2 flex h-6 items-center gap-1 rounded-full bg-neutral-200/80 px-2">
							<svg class="h-2.5 w-2.5 text-neutral-400" viewBox="0 0 20 20" fill="currentColor">
								<path
									fill-rule="evenodd"
									d="M9 3.5a5.5 5.5 0 1 0 0 11 5.5 5.5 0 0 0 0-11ZM2 9a7 7 0 1 1 12.452 4.391l3.328 3.329a.75.75 0 1 1-1.06 1.06l-3.329-3.328A7 7 0 0 1 2 9Z"
									clip-rule="evenodd"
								/>
							</svg>
							<span class="text-[8px] text-neutral-400">Pesquisar...</span>
						</div>

						<!-- mini listing grid -->
						<div class="mt-2 grid grid-cols-2 gap-1.5">
							<div class="overflow-hidden rounded-lg bg-white shadow-sm">
								<div
									class="flex h-12 items-center justify-center bg-linear-to-br from-sky-200 to-sky-400"
								>
									<svg class="h-7 w-7 text-sky-700" viewBox="0 0 24 24" fill="currentColor">
										<path
											d="M10.5 1.875a1.125 1.125 0 0 1 1.125 1.125v.375h.75V3a1.125 1.125 0 0 1 2.25 0v.375c1.864 0 3.375 1.51 3.375 3.375v12c0 1.864-1.51 3.375-3.375 3.375h-6A3.375 3.375 0 0 1 5.25 18.75v-12c0-1.864 1.51-3.375 3.375-3.375V3c0-.621.504-1.125 1.125-1.125h.75ZM9 6.375A1.875 1.875 0 0 0 7.125 8.25v9A1.875 1.875 0 0 0 9 19.125h6a1.875 1.875 0 0 0 1.875-1.875v-9A1.875 1.875 0 0 0 15 6.375H9Z"
										/>
									</svg>
								</div>
								<div class="p-1">
									<p class="text-[8px] font-extrabold text-brand-700">185.000 Kz</p>
									<p class="text-[7px] text-neutral-500">Smartphone 128GB</p>
								</div>
							</div>
							<div class="overflow-hidden rounded-lg bg-white shadow-sm">
								<div
									class="flex h-12 items-center justify-center bg-linear-to-br from-amber-200 to-amber-400"
								>
									<svg class="h-7 w-7 text-amber-700" viewBox="0 0 24 24" fill="currentColor">
										<path
											d="M3.375 4.5C2.339 4.5 1.5 5.34 1.5 6.375V13.5h12V6.375c0-1.036-.84-1.875-1.875-1.875h-8.25ZM13.5 15h-12v2.625c0 1.035.84 1.875 1.875 1.875h.375a3 3 0 1 1 6 0h3a.75.75 0 0 0 .75-.75V15Z"
										/><path
											d="M8.25 19.5a1.5 1.5 0 1 0-3 0 1.5 1.5 0 0 0 3 0ZM15.75 6.75a.75.75 0 0 0-.75.75v11.25c0 .087.015.17.042.248a3 3 0 0 1 5.958.464.75.75 0 0 0 .75-.75v-4.992a3 3 0 0 0-.879-2.121l-2.56-2.56a3 3 0 0 0-2.121-.879h-.44Z"
										/><path d="M19.5 19.5a1.5 1.5 0 1 0-3 0 1.5 1.5 0 0 0 3 0Z" />
									</svg>
								</div>
								<div class="p-1">
									<p class="text-[8px] font-extrabold text-brand-700">4.500.000 Kz</p>
									<p class="text-[7px] text-neutral-500">Toyota Corolla</p>
								</div>
							</div>
							<div class="overflow-hidden rounded-lg bg-white shadow-sm">
								<div
									class="flex h-12 items-center justify-center bg-linear-to-br from-rose-200 to-rose-400"
								>
									<svg class="h-7 w-7 text-rose-700" viewBox="0 0 24 24" fill="currentColor">
										<path
											fill-rule="evenodd"
											d="M7.5 6a4.5 4.5 0 1 1 9 0 4.5 4.5 0 0 1-9 0ZM3.751 20.105a8.25 8.25 0 0 1 16.498 0 .75.75 0 0 1-.437.695A18.683 18.683 0 0 1 12 22.5c-2.786 0-5.433-.608-7.812-1.7a.75.75 0 0 1-.437-.695Z"
											clip-rule="evenodd"
										/>
									</svg>
								</div>
								<div class="p-1">
									<p class="text-[8px] font-extrabold text-brand-700">12.000 Kz</p>
									<p class="text-[7px] text-neutral-500">Camisa nova</p>
								</div>
							</div>
							<div class="overflow-hidden rounded-lg bg-white shadow-sm">
								<div
									class="flex h-12 items-center justify-center bg-linear-to-br from-violet-200 to-violet-400"
								>
									<svg class="h-7 w-7 text-violet-700" viewBox="0 0 24 24" fill="currentColor">
										<path d="M19.5 6h-15v9h15V6Z" /><path
											fill-rule="evenodd"
											d="M3.375 3C2.339 3 1.5 3.84 1.5 4.875v11.25C1.5 17.16 2.34 18 3.375 18H9.75v1.5H6A.75.75 0 0 0 6 21h12a.75.75 0 0 0 0-1.5h-3.75V18h6.375c1.035 0 1.875-.84 1.875-1.875V4.875C22.5 3.839 21.66 3 20.625 3H3.375Zm0 13.5h17.25a.375.375 0 0 0 .375-.375V4.875a.375.375 0 0 0-.375-.375H3.375A.375.375 0 0 0 3 4.875v11.25c0 .207.168.375.375.375Z"
											clip-rule="evenodd"
										/>
									</svg>
								</div>
								<div class="p-1">
									<p class="text-[8px] font-extrabold text-brand-700">320.000 Kz</p>
									<p class="text-[7px] text-neutral-500">Monitor 27"</p>
								</div>
							</div>
						</div>

						<!-- bottom nav -->
						<div
							class="mt-2 flex items-center justify-around rounded-full bg-white px-3 py-1.5 shadow-sm"
						>
							<span class="h-2 w-2 rounded-full bg-brand-500"></span>
							<span class="h-2 w-2 rounded-full bg-neutral-200"></span>
							<span
								class="flex h-5 w-5 items-center justify-center rounded-full bg-brand-600 text-[10px] font-bold text-white"
								>+</span
							>
							<span class="h-2 w-2 rounded-full bg-neutral-200"></span>
							<span class="h-2 w-2 rounded-full bg-neutral-200"></span>
						</div>
					</div>
				</div>
			</div>

			<!-- floating WhatsApp bubble -->
			<div
				class="absolute -bottom-4 -left-8 flex -rotate-3 items-center gap-2 rounded-xl bg-white px-3 py-2 shadow-xl shadow-brand-950/30"
			>
				<span class="flex h-7 w-7 items-center justify-center rounded-full bg-whatsapp">
					<svg class="h-4 w-4 text-white" viewBox="0 0 24 24" fill="currentColor">
						<path
							d="M17.472 14.382c-.297-.149-1.758-.867-2.03-.967-.273-.099-.471-.148-.67.15-.197.297-.767.966-.94 1.164-.173.199-.347.223-.644.075-.297-.15-1.255-.463-2.39-1.475-.883-.788-1.48-1.761-1.653-2.059-.173-.297-.018-.458.13-.606.134-.133.298-.347.446-.52.149-.174.198-.298.298-.497.099-.198.05-.371-.025-.52-.075-.149-.669-1.612-.916-2.207-.242-.579-.487-.5-.669-.51-.173-.008-.371-.01-.57-.01-.198 0-.52.074-.792.372-.272.297-1.04 1.016-1.04 2.479 0 1.462 1.065 2.875 1.213 3.074.149.198 2.096 3.2 5.077 4.487.709.306 1.262.489 1.694.625.712.227 1.36.195 1.871.118.571-.085 1.758-.719 2.006-1.413.248-.694.248-1.289.173-1.413-.074-.124-.272-.198-.57-.347Z"
						/>
					</svg>
				</span>
				<div>
					<p class="text-xs font-bold text-neutral-900">Ainda está disponível?</p>
					<p class="text-[10px] text-neutral-400">agora mesmo</p>
				</div>
			</div>

			<!-- floating "free" tag -->
			<div
				class="absolute -top-3 -right-5 rotate-6 rounded-full bg-amber-400 px-3 py-1.5 text-xs font-extrabold text-amber-950 shadow-xl shadow-brand-950/30"
			>
				Grátis!
			</div>
		</div>
	</div>
</section>

<!-- Categories -->
<svelte:boundary>
	{@const cats = await categories}
	{#if cats.length}
		<section class="mt-6">
			<div class="flex items-center justify-between">
				<h2 class="text-xl font-bold text-neutral-900">Categorias</h2>
				<a href="/categories" class="text-sm font-semibold text-brand-700 hover:text-brand-800"
					>Ver todas</a
				>
			</div>
			<div
				class="mt-4 grid grid-flow-col grid-rows-1 gap-3 overflow-x-auto pb-2 sm:grid-flow-row sm:grid-cols-4 sm:grid-rows-none lg:grid-cols-8"
			>
				{#each cats as category (category.id)}
					<a
						href={`/search?category_id=${category.id}`}
						class="group flex w-24 shrink-0 flex-col items-center gap-2 rounded-2xl border border-neutral-200 bg-white px-2 py-4 text-center transition hover:-translate-y-0.5 hover:border-brand-300 hover:shadow-md hover:shadow-neutral-900/5 active:scale-95 sm:w-auto"
					>
						<span
							class="flex h-12 w-12 items-center justify-center rounded-xl bg-brand-50 text-brand-600 transition group-hover:bg-brand-100"
						>
							<CategoryIcon slug={category.slug} name={category.name} class="h-6 w-6" />
						</span>
						<span class="line-clamp-2 text-xs leading-tight font-semibold text-neutral-700">
							{category.name}
						</span>
					</a>
				{/each}
			</div>
		</section>
	{/if}

	{#snippet pending()}
		<div class="mt-8 flex gap-2">
			{#each { length: 6 }, i (i)}
				<div class="h-10 w-28 animate-pulse rounded-full bg-neutral-200"></div>
			{/each}
		</div>
	{/snippet}

	{#snippet failed()}{/snippet}
</svelte:boundary>

<!-- Featured listings -->
<svelte:boundary>
	{@const result = await featured}
	{#if result.data.length}
		<section class="mt-8">
			<div class="flex items-center gap-2">
				<svg
					class="h-5 w-5 text-amber-400"
					viewBox="0 0 20 20"
					fill="currentColor"
					aria-hidden="true"
				>
					<path
						fill-rule="evenodd"
						d="M10.868 2.884c-.321-.772-1.415-.772-1.736 0l-1.83 4.401-4.753.381c-.833.067-1.171 1.107-.536 1.651l3.62 3.102-1.106 4.637c-.194.813.691 1.456 1.405 1.02L10 15.591l4.069 2.485c.713.436 1.598-.207 1.404-1.02l-1.106-4.637 3.62-3.102c.635-.544.297-1.584-.536-1.65l-4.752-.382-1.831-4.401Z"
						clip-rule="evenodd"
					/>
				</svg>
				<h2 class="text-xl font-bold text-neutral-900">Em destaque</h2>
			</div>
			<div class="mt-4 grid grid-cols-2 gap-4 sm:grid-cols-3 lg:grid-cols-4">
				{#each result.data as listing (listing.id)}
					<ListingCard {listing} />
				{/each}
			</div>
		</section>
	{/if}

	{#snippet pending()}{/snippet}
	{#snippet failed()}{/snippet}
</svelte:boundary>

<!-- Recent listings -->
<section class="mt-8">
	<div class="flex items-center justify-between">
		<h2 class="text-xl font-bold text-neutral-900">Anúncios recentes</h2>
		<a href="/search" class="text-sm font-semibold text-brand-700 hover:text-brand-800">Ver todos</a
		>
	</div>

	<svelte:boundary>
		{@const result = await listings}
		{#if result.data.length}
			<div class="mt-4 grid grid-cols-2 gap-4 sm:grid-cols-3 lg:grid-cols-4">
				{#each result.data as listing (listing.id)}
					<ListingCard {listing} />
				{/each}
			</div>
		{:else}
			<div
				class="mt-4 rounded-2xl border border-dashed border-neutral-300 bg-white p-10 text-center"
			>
				<p class="font-medium text-neutral-700">Ainda não há anúncios publicados.</p>
				<a
					href="/anunciar"
					class="mt-4 inline-block rounded-full bg-brand-700 px-6 py-2.5 text-sm font-semibold text-white hover:bg-brand-800"
					>Sê o primeiro a anunciar</a
				>
			</div>
		{/if}

		{#snippet pending()}
			<div class="mt-4 grid grid-cols-2 gap-4 sm:grid-cols-3 lg:grid-cols-4">
				{#each { length: 8 }, i (i)}
					<div class="overflow-hidden rounded-2xl border border-neutral-200 bg-white">
						<div class="aspect-4/3 animate-pulse bg-neutral-200"></div>
						<div class="space-y-2 p-4">
							<div class="h-5 w-24 animate-pulse rounded bg-neutral-200"></div>
							<div class="h-4 w-full animate-pulse rounded bg-neutral-100"></div>
						</div>
					</div>
				{/each}
			</div>
		{/snippet}

		{#snippet failed()}
			<div class="mt-4 rounded-2xl border border-red-200 bg-red-50 p-6 text-sm text-red-700">
				Não foi possível carregar os anúncios. Tenta novamente mais tarde.
			</div>
		{/snippet}
	</svelte:boundary>
</section>

<!-- CTA -->
<section
	class="mt-12 flex flex-col items-center gap-4 rounded-3xl border border-brand-200 bg-brand-50 p-8 text-center sm:flex-row sm:justify-between sm:text-left"
>
	<div>
		<h2 class="text-lg font-bold text-brand-900">Tens algo para vender?</h2>
		<p class="mt-1 text-sm text-brand-800/80">
			Publica gratuitamente e chega a milhares de compradores em todo o país.
		</p>
	</div>
	<a
		href="/anunciar"
		class="shrink-0 rounded-full bg-brand-700 px-6 py-3 text-sm font-semibold text-white transition hover:bg-brand-800"
		>Anunciar agora</a
	>
</section>
