<script lang="ts">
	import { page } from '$app/state';
	import ListingCard from '$lib/components/ListingCard.svelte';
	import Pagination from '$lib/components/Pagination.svelte';
	import Seo from '$lib/components/Seo.svelte';

	let { data } = $props();
	const seller = $derived(data.seller);
	const listings = $derived(data.listings);

	const memberSince = $derived(
		new Date(seller.created_at).toLocaleDateString('pt-PT', { month: 'long', year: 'numeric' })
	);

	let copied = $state(false);

	const shopUrl = $derived(`${page.url.origin}/loja/${seller.id}`);

	const whatsappShare = $derived(
		`https://wa.me/?text=${encodeURIComponent(`Vê a loja de ${seller.name} no Despacha Aí: ${shopUrl}`)}`
	);

	async function copyLink() {
		try {
			await navigator.clipboard.writeText(shopUrl);
			copied = true;
			setTimeout(() => (copied = false), 2000);
		} catch {
			// área de transferência indisponível — o URL está na barra de endereço
		}
	}
</script>

<Seo
	title={`Loja de ${seller.name} | Despacha Aí`}
	description={`${listings.total} ${listings.total === 1 ? 'anúncio' : 'anúncios'} de ${seller.name} no Despacha Aí. Compra directamente pelo WhatsApp.`}
	image={seller.avatar_url}
/>

<header
	class="flex flex-col items-start gap-4 rounded-2xl border border-neutral-200 bg-white p-6 sm:flex-row sm:items-center sm:justify-between"
>
	<div class="flex items-center gap-4">
		{#if seller.avatar_url}
			<img src={seller.avatar_url} alt="" class="h-16 w-16 rounded-full object-cover" />
		{:else}
			<span
				class="flex h-16 w-16 items-center justify-center rounded-full bg-brand-100 text-2xl font-bold text-brand-800"
			>
				{seller.name[0]?.toUpperCase() ?? '?'}
			</span>
		{/if}
		<div>
			<h1 class="flex items-center gap-1.5 text-xl font-bold text-neutral-900">
				{seller.name}
				{#if seller.is_verified}
					<svg
						class="h-5 w-5 text-brand-600"
						viewBox="0 0 20 20"
						fill="currentColor"
						aria-label="Vendedor verificado"
					>
						<path
							fill-rule="evenodd"
							d="M16.403 12.652a3 3 0 0 0 0-5.304 3 3 0 0 0-3.75-3.751 3 3 0 0 0-5.305 0 3 3 0 0 0-3.751 3.75 3 3 0 0 0 0 5.305 3 3 0 0 0 3.75 3.751 3 3 0 0 0 5.305 0 3 3 0 0 0 3.751-3.75Zm-2.546-4.46a.75.75 0 0 0-1.214-.883l-3.483 4.79-1.88-1.88a.75.75 0 1 0-1.06 1.061l2.5 2.5a.75.75 0 0 0 1.137-.089l4-5.5Z"
							clip-rule="evenodd"
						/>
					</svg>
				{/if}
			</h1>
			<p class="mt-0.5 text-sm text-neutral-500">
				Membro desde {memberSince} · {listings.total}
				{listings.total === 1 ? 'anúncio' : 'anúncios'}
			</p>
		</div>
	</div>

	<div class="flex gap-2">
		<a
			href={whatsappShare}
			target="_blank"
			rel="noopener noreferrer"
			class="flex h-10 items-center gap-2 rounded-full bg-whatsapp px-4 text-sm font-bold text-white transition hover:brightness-95"
		>
			<svg class="h-4 w-4" viewBox="0 0 24 24" fill="currentColor" aria-hidden="true">
				<path
					d="M17.472 14.382c-.297-.149-1.758-.867-2.03-.967-.273-.099-.471-.148-.67.15-.197.297-.767.966-.94 1.164-.173.199-.347.223-.644.075-.297-.15-1.255-.463-2.39-1.475-.883-.788-1.48-1.761-1.653-2.059-.173-.297-.018-.458.13-.606.134-.133.298-.347.446-.52.149-.174.198-.298.298-.497.099-.198.05-.371-.025-.52-.075-.149-.669-1.612-.916-2.207-.242-.579-.487-.5-.669-.51-.173-.008-.371-.01-.57-.01-.198 0-.52.074-.792.372-.272.297-1.04 1.016-1.04 2.479 0 1.462 1.065 2.875 1.213 3.074.149.198 2.096 3.2 5.077 4.487.709.306 1.262.489 1.694.625.712.227 1.36.195 1.871.118.571-.085 1.758-.719 2.006-1.413.248-.694.248-1.289.173-1.413-.074-.124-.272-.198-.57-.347m-5.421 7.403h-.004a9.87 9.87 0 0 1-5.031-1.378l-.361-.214-3.741.982.998-3.648-.235-.374a9.86 9.86 0 0 1-1.51-5.26c.001-5.45 4.436-9.884 9.888-9.884 2.64 0 5.122 1.03 6.988 2.898a9.825 9.825 0 0 1 2.893 6.994c-.003 5.45-4.437 9.884-9.885 9.884m8.413-18.297A11.815 11.815 0 0 0 12.05 0C5.495 0 .16 5.335.157 11.892c0 2.096.547 4.142 1.588 5.945L.057 24l6.305-1.654a11.882 11.882 0 0 0 5.683 1.448h.005c6.554 0 11.89-5.335 11.893-11.893a11.821 11.821 0 0 0-3.48-8.413Z"
				/>
			</svg>
			Partilhar
		</a>
		<button
			type="button"
			onclick={copyLink}
			class="flex h-10 items-center gap-2 rounded-full border border-neutral-200 px-4 text-sm font-medium text-neutral-600 transition hover:border-brand-300 hover:text-brand-700"
		>
			{copied ? 'Link copiado!' : 'Copiar link'}
		</button>
	</div>
</header>

<section class="mt-6">
	{#if listings.data.length}
		<div class="grid grid-cols-2 gap-4 sm:grid-cols-3 xl:grid-cols-4">
			{#each listings.data as listing (listing.id)}
				<ListingCard {listing} />
			{/each}
		</div>
		<Pagination page={listings.page} totalPages={listings.total_pages} />
	{:else}
		<div class="rounded-2xl border border-dashed border-neutral-300 bg-white p-12 text-center">
			<p class="text-lg font-semibold text-neutral-800">Esta loja ainda não tem anúncios</p>
			<p class="mt-1 text-sm text-neutral-500">Volta mais tarde para veres as novidades.</p>
		</div>
	{/if}
</section>
