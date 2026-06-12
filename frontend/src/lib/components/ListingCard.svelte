<script lang="ts">
	import type { Listing } from '$lib/types';
	import { formatPrice, locationText, timeAgo } from '$lib/utils';

	let { listing }: { listing: Listing } = $props();

	let favorited = $state(false);

	const isRecent = $derived(Date.now() - new Date(listing.created_at).getTime() < 86_400_000);

	function handleImageError(event: Event) {
		const image = event.currentTarget as HTMLImageElement;
		image.src = '/product-placeholder.svg';
		image.alt = 'Imagem indisponível';
	}

	function toggleFavorite(event: MouseEvent) {
		event.preventDefault();
		event.stopPropagation();
		favorited = !favorited;
	}
</script>

<a
	href={`/anuncio/${listing.id}`}
	class="group flex h-full flex-col overflow-hidden rounded-2xl border border-neutral-200 bg-white transition hover:-translate-y-0.5 hover:border-brand-300 hover:shadow-lg hover:shadow-neutral-900/5"
>
	<div class="relative aspect-4/3 overflow-hidden bg-neutral-100">
		<img
			src={listing.images?.[0]?.image_url ?? '/product-placeholder.svg'}
			alt={listing.title}
			loading="lazy"
			class="h-full w-full object-cover transition duration-300 group-hover:scale-105"
			onerror={handleImageError}
		/>

		<!-- favorito -->
		<button
			type="button"
			onclick={toggleFavorite}
			aria-label={favorited ? 'Remover dos favoritos' : 'Adicionar aos favoritos'}
			aria-pressed={favorited}
			class="absolute top-2.5 right-2.5 flex h-9 w-9 items-center justify-center rounded-full bg-white/90 shadow-sm backdrop-blur transition hover:scale-110 active:scale-95 {favorited
				? 'text-red-500'
				: 'text-neutral-400 hover:text-red-500'}"
		>
			<svg class="h-5 w-5" viewBox="0 0 24 24" fill={favorited ? 'currentColor' : 'none'}>
				<path
					d="M21 8.25c0-2.485-2.099-4.5-4.688-4.5-1.935 0-3.597 1.126-4.312 2.733-.715-1.607-2.377-2.733-4.313-2.733C5.1 3.75 3 5.765 3 8.25c0 7.22 9 12 9 12s9-4.78 9-12Z"
					stroke="currentColor"
					stroke-width="1.8"
					stroke-linecap="round"
					stroke-linejoin="round"
				/>
			</svg>
		</button>

		<div class="absolute top-2.5 left-2.5 flex flex-col items-start gap-1.5">
			{#if listing.is_featured}
				<span
					class="flex items-center gap-1 rounded-full bg-amber-400 px-2.5 py-1 text-xs font-bold text-amber-950 shadow-sm"
				>
					<svg class="h-3 w-3" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
						<path
							fill-rule="evenodd"
							d="M10.868 2.884c-.321-.772-1.415-.772-1.736 0l-1.83 4.401-4.753.381c-.833.067-1.171 1.107-.536 1.651l3.62 3.102-1.106 4.637c-.194.813.691 1.456 1.405 1.02L10 15.591l4.069 2.485c.713.436 1.598-.207 1.404-1.02l-1.106-4.637 3.62-3.102c.635-.544.297-1.584-.536-1.65l-4.752-.382-1.831-4.401Z"
							clip-rule="evenodd"
						/>
					</svg>
					Destaque
				</span>
			{/if}
			{#if listing.condition === 'new'}
				<span class="rounded-full bg-brand-600 px-2.5 py-1 text-xs font-bold text-white shadow-sm"
					>Novo</span
				>
			{/if}
		</div>
		{#if listing.status === 'sold'}
			<div class="absolute inset-0 flex items-center justify-center bg-neutral-900/40">
				<span class="rounded-full bg-neutral-900/90 px-4 py-1.5 text-sm font-bold text-white"
					>Vendido</span
				>
			</div>
		{/if}
	</div>

	<div class="flex flex-1 flex-col gap-1 p-4">
		<p class="text-xl font-extrabold tracking-tight text-brand-700">
			{formatPrice(listing.price)}
		</p>
		<h3 class="line-clamp-2 text-sm leading-snug font-medium text-neutral-800">
			{listing.title}
		</h3>
		<div class="mt-auto flex items-center gap-1.5 pt-2 text-xs">
			<svg
				class="h-3.5 w-3.5 shrink-0 text-neutral-400"
				viewBox="0 0 20 20"
				fill="currentColor"
				aria-hidden="true"
			>
				<path
					fill-rule="evenodd"
					d="m9.69 18.933.003.001C9.89 19.02 10 19 10 19s.11.02.308-.066l.002-.001.006-.003.018-.008a5.741 5.741 0 0 0 .281-.14c.186-.096.446-.24.757-.433.62-.384 1.445-.966 2.274-1.765C15.302 14.988 17 12.493 17 9A7 7 0 1 0 3 9c0 3.492 1.698 5.988 3.355 7.584a13.731 13.731 0 0 0 2.273 1.765 11.842 11.842 0 0 0 .976.544l.062.029.018.008.006.003ZM10 11.25a2.25 2.25 0 1 0 0-4.5 2.25 2.25 0 0 0 0 4.5Z"
					clip-rule="evenodd"
				/>
			</svg>
			<span class="truncate font-medium text-neutral-600"
				>{locationText(listing.province, listing.city)}</span
			>
			<span class="ml-auto flex shrink-0 items-center gap-1 text-neutral-400">
				{#if isRecent && listing.status !== 'sold'}
					<span class="h-1.5 w-1.5 rounded-full bg-brand-500"></span>
				{/if}
				{timeAgo(listing.created_at)}
			</span>
		</div>
	</div>
</a>
