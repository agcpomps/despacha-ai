<script lang="ts">
	import { page } from '$app/state';
	import { getListings } from '$lib/remote/listings.remote';
	import { getCategories } from '$lib/remote/categories.remote';
	import ListingCard from '$lib/components/ListingCard.svelte';
	import Pagination from '$lib/components/Pagination.svelte';
	import Seo from '$lib/components/Seo.svelte';
	import { PROVINCES } from '$lib/utils';
	import type { ListingFilters, ListingSort } from '$lib/types';

	const categories = getCategories();

	function parseNumber(value: string | null) {
		if (!value) return undefined;
		const parsed = Number(value);
		return Number.isFinite(parsed) && parsed >= 0 ? parsed : undefined;
	}

	const filters: ListingFilters = $derived.by(() => {
		const params = page.url.searchParams;
		const sort = params.get('sort') as ListingSort | null;
		return {
			search: params.get('search') ?? undefined,
			category_id: params.get('category_id') ?? undefined,
			province: params.get('province') ?? undefined,
			min_price: parseNumber(params.get('min_price')),
			max_price: parseNumber(params.get('max_price')),
			sort:
				sort && ['newest', 'oldest', 'price_asc', 'price_desc'].includes(sort) ? sort : undefined,
			page: parseNumber(params.get('page')) || 1,
			limit: 16
		};
	});

	const result = $derived(getListings(filters));
</script>

<Seo
	title={`${filters.search ? `${filters.search} — ` : ''}Anúncios${filters.province ? ` em ${filters.province}` : ' em Benguela e toda Angola'} | Despacha Aí`}
	description={`Encontra ${filters.search || 'artigos novos e usados'}${filters.province ? ` em ${filters.province}` : ' em Benguela, Lobito e toda Angola'}. Classificados grátis com contacto directo pelo WhatsApp.`}
/>

<div class="grid gap-6 lg:grid-cols-[260px_1fr]">
	<!-- Filters sidebar -->
	<aside>
		<form
			method="GET"
			class="space-y-4 rounded-2xl border border-neutral-200 bg-white p-5 lg:sticky lg:top-24"
		>
			<h2 class="text-sm font-bold tracking-wide text-neutral-900 uppercase">Filtros</h2>

			<label class="block">
				<span class="text-xs font-semibold text-neutral-600">Pesquisa</span>
				<input
					type="search"
					name="search"
					value={filters.search ?? ''}
					placeholder="Palavra-chave"
					class="mt-1 h-10 w-full rounded-lg border-neutral-200 text-sm focus:border-brand-500 focus:ring-brand-500"
				/>
			</label>

			<label class="block">
				<span class="text-xs font-semibold text-neutral-600">Categoria</span>
				<select
					name="category_id"
					class="mt-1 h-10 w-full rounded-lg border-neutral-200 text-sm focus:border-brand-500 focus:ring-brand-500"
				>
					<option value="">Todas</option>
					<svelte:boundary>
						{#each await categories as category (category.id)}
							<option value={category.id} selected={filters.category_id === category.id}>
								{category.name}
							</option>
							{#each category.children ?? [] as child (child.id)}
								<option value={child.id} selected={filters.category_id === child.id}>
									&nbsp;&nbsp;{child.name}
								</option>
							{/each}
						{/each}
						{#snippet pending()}{/snippet}
						{#snippet failed()}{/snippet}
					</svelte:boundary>
				</select>
			</label>

			<label class="block">
				<span class="text-xs font-semibold text-neutral-600">Província</span>
				<select
					name="province"
					class="mt-1 h-10 w-full rounded-lg border-neutral-200 text-sm focus:border-brand-500 focus:ring-brand-500"
				>
					<option value="">Todo o país</option>
					{#each PROVINCES as province (province)}
						<option value={province} selected={filters.province === province}>{province}</option>
					{/each}
				</select>
			</label>

			<div class="grid grid-cols-2 gap-2">
				<label class="block">
					<span class="text-xs font-semibold text-neutral-600">Preço mín.</span>
					<input
						type="number"
						name="min_price"
						min="0"
						value={filters.min_price ?? ''}
						placeholder="0"
						class="mt-1 h-10 w-full rounded-lg border-neutral-200 text-sm focus:border-brand-500 focus:ring-brand-500"
					/>
				</label>
				<label class="block">
					<span class="text-xs font-semibold text-neutral-600">Preço máx.</span>
					<input
						type="number"
						name="max_price"
						min="0"
						value={filters.max_price ?? ''}
						placeholder="Kz"
						class="mt-1 h-10 w-full rounded-lg border-neutral-200 text-sm focus:border-brand-500 focus:ring-brand-500"
					/>
				</label>
			</div>

			<label class="block">
				<span class="text-xs font-semibold text-neutral-600">Ordenar por</span>
				<select
					name="sort"
					class="mt-1 h-10 w-full rounded-lg border-neutral-200 text-sm focus:border-brand-500 focus:ring-brand-500"
				>
					<option value="newest" selected={!filters.sort || filters.sort === 'newest'}
						>Mais recentes</option
					>
					<option value="oldest" selected={filters.sort === 'oldest'}>Mais antigos</option>
					<option value="price_asc" selected={filters.sort === 'price_asc'}>Preço: menor</option>
					<option value="price_desc" selected={filters.sort === 'price_desc'}>Preço: maior</option>
				</select>
			</label>

			<div class="flex gap-2 pt-1">
				<button
					type="submit"
					class="h-10 flex-1 rounded-full bg-brand-700 text-sm font-semibold text-white transition hover:bg-brand-800"
					>Aplicar</button
				>
				<a
					href="/search"
					class="flex h-10 items-center rounded-full border border-neutral-200 px-4 text-sm font-medium text-neutral-600 hover:border-neutral-300"
					>Limpar</a
				>
			</div>
		</form>
	</aside>

	<!-- Results -->
	<section>
		<svelte:boundary>
			{@const listings = await result}
			<div class="flex items-baseline justify-between">
				<h1 class="text-xl font-bold text-neutral-900">
					{filters.search ? `Resultados para “${filters.search}”` : 'Explorar anúncios'}
				</h1>
				<p class="text-sm text-neutral-500">{listings.total} resultados</p>
			</div>

			{#if listings.data.length}
				<div class="mt-4 grid grid-cols-2 gap-4 xl:grid-cols-3">
					{#each listings.data as listing (listing.id)}
						<ListingCard {listing} />
					{/each}
				</div>
				<Pagination page={listings.page} totalPages={listings.total_pages} />
			{:else}
				<div
					class="mt-4 rounded-2xl border border-dashed border-neutral-300 bg-white p-12 text-center"
				>
					<p class="text-lg font-semibold text-neutral-800">Nenhum anúncio encontrado</p>
					<p class="mt-1 text-sm text-neutral-500">
						Tenta ajustar os filtros ou pesquisar por outro termo.
					</p>
				</div>
			{/if}

			{#snippet pending()}
				<div class="h-7 w-48 animate-pulse rounded bg-neutral-200"></div>
				<div class="mt-4 grid grid-cols-2 gap-4 xl:grid-cols-3">
					{#each { length: 6 }, i (i)}
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

			{#snippet failed(error)}
				<div class="rounded-2xl border border-red-200 bg-red-50 p-6 text-sm text-red-700">
					{(error as Error)?.message ?? 'Não foi possível carregar os anúncios.'}
				</div>
			{/snippet}
		</svelte:boundary>
	</section>
</div>
