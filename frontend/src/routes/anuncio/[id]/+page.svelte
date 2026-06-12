<script lang="ts">
	import { page } from '$app/state';
	import {
		getListing,
		featureListing,
		unfeatureListing,
		bumpListing
	} from '$lib/remote/listings.remote';
	import { getCurrentUser } from '$lib/remote/auth.remote';
	import Seo from '$lib/components/Seo.svelte';
	import { formatPrice, locationText, timeAgo, CONDITION_LABELS } from '$lib/utils';

	const listing = $derived(getListing(page.params.id!));
	const user = getCurrentUser();

	let activeImage = $state(0);
	let adminBusy = $state(false);
	let adminError = $state<string | null>(null);
	let copied = $state(false);

	async function copyNumber(number: string) {
		try {
			await navigator.clipboard.writeText(number);
			copied = true;
			setTimeout(() => (copied = false), 2000);
		} catch {
			// clipboard indisponível — o número está visível na mesma
		}
	}

	async function runAdminAction(action: () => Promise<unknown>) {
		adminBusy = true;
		adminError = null;
		try {
			await action();
		} catch {
			adminError = 'A acção falhou. Tenta novamente.';
		} finally {
			adminBusy = false;
		}
	}

	function whatsappLink(phone: string, title: string) {
		const digits = phone.replace(/\D/g, '');
		const text = encodeURIComponent(
			`Olá! Vi o seu anúncio "${title}" no Despacha Aí e estou interessado(a).`
		);
		return `https://wa.me/${digits}?text=${text}`;
	}

	function handleImageError(event: Event) {
		const image = event.currentTarget as HTMLImageElement;
		image.src = '/product-placeholder.svg';
	}
</script>

<svelte:boundary>
	{@const item = await listing}

	<Seo
		title={`${item.title} — ${formatPrice(item.price)} em ${locationText(item.province, item.city)} | Despacha Aí`}
		description={`${item.title} por ${formatPrice(item.price)} em ${locationText(item.province, item.city)}, Angola. ${item.description.slice(0, 120)}`}
		type="product"
		image={item.images?.[0]?.image_url}
		jsonLd={{
			'@context': 'https://schema.org',
			'@type': 'Product',
			name: item.title,
			description: item.description.slice(0, 300),
			image: item.images.map((img) => img.image_url),
			offers: {
				'@type': 'Offer',
				price: item.price,
				priceCurrency: 'AOA',
				itemCondition:
					item.condition === 'new'
						? 'https://schema.org/NewCondition'
						: 'https://schema.org/UsedCondition',
				availability:
					item.status === 'active' ? 'https://schema.org/InStock' : 'https://schema.org/SoldOut',
				areaServed: { '@type': 'Place', name: `${item.province}, Angola` }
			}
		}}
	/>

	<nav class="mb-4 text-sm text-neutral-500" aria-label="Navegação">
		<a href="/" class="hover:text-brand-700">Início</a>
		<span class="mx-1">/</span>
		{#if item.category}
			<a href={`/search?category_id=${item.category.id}`} class="hover:text-brand-700"
				>{item.category.name}</a
			>
			<span class="mx-1">/</span>
		{/if}
		<span class="text-neutral-700">{item.title}</span>
	</nav>

	<div class="grid gap-6 lg:grid-cols-[1fr_360px]">
		<!-- Gallery + description -->
		<div>
			<div class="overflow-hidden rounded-2xl border border-neutral-200 bg-white">
				<div class="relative aspect-4/3 bg-neutral-100">
					<img
						src={item.images?.[activeImage]?.image_url ?? '/product-placeholder.svg'}
						alt={item.title}
						class="h-full w-full object-contain"
						onerror={handleImageError}
					/>
					{#if item.status === 'sold'}
						<span
							class="absolute top-4 left-4 rounded-full bg-neutral-900/85 px-3 py-1.5 text-sm font-semibold text-white"
							>Vendido</span
						>
					{/if}
				</div>
				{#if item.images.length > 1}
					<div class="flex gap-2 overflow-x-auto border-t border-neutral-100 p-3">
						{#each item.images as image, index (image.id)}
							<button
								type="button"
								class="h-16 w-16 shrink-0 overflow-hidden rounded-lg border-2 transition {index ===
								activeImage
									? 'border-brand-600'
									: 'border-transparent opacity-70 hover:opacity-100'}"
								onclick={() => (activeImage = index)}
								aria-label={`Imagem ${index + 1}`}
							>
								<img
									src={image.image_url}
									alt=""
									class="h-full w-full object-cover"
									onerror={handleImageError}
								/>
							</button>
						{/each}
					</div>
				{/if}
			</div>

			<section class="mt-6 rounded-2xl border border-neutral-200 bg-white p-6">
				<h2 class="text-base font-bold text-neutral-900">Descrição</h2>
				<p class="mt-3 text-sm leading-relaxed whitespace-pre-line text-neutral-700">
					{item.description}
				</p>
			</section>
		</div>

		<!-- Summary + seller -->
		<aside class="space-y-4">
			<div class="rounded-2xl border border-neutral-200 bg-white p-6">
				<p class="text-3xl font-extrabold text-brand-800">{formatPrice(item.price)}</p>
				<h1 class="mt-2 text-xl leading-snug font-bold text-neutral-900">{item.title}</h1>

				<dl class="mt-4 space-y-2 text-sm">
					<div class="flex justify-between">
						<dt class="text-neutral-500">Estado</dt>
						<dd class="font-medium text-neutral-800">
							{CONDITION_LABELS[item.condition] ?? item.condition}
						</dd>
					</div>
					<div class="flex justify-between">
						<dt class="text-neutral-500">Localização</dt>
						<dd class="font-medium text-neutral-800">
							{locationText(item.province, item.city)}
						</dd>
					</div>
					{#if item.address_reference}
						<div class="flex justify-between gap-4">
							<dt class="shrink-0 text-neutral-500">Referência</dt>
							<dd class="text-right font-medium text-neutral-800">{item.address_reference}</dd>
						</div>
					{/if}
					<div class="flex justify-between">
						<dt class="text-neutral-500">Publicado</dt>
						<dd class="font-medium text-neutral-800">{timeAgo(item.created_at)}</dd>
					</div>
				</dl>
			</div>

			<div class="rounded-2xl border border-neutral-200 bg-white p-6">
				<h2 class="text-sm font-bold tracking-wide text-neutral-500 uppercase">Vendedor</h2>
				<div class="mt-3 flex items-center gap-3">
					{#if item.seller?.avatar_url}
						<img src={item.seller.avatar_url} alt="" class="h-12 w-12 rounded-full object-cover" />
					{:else}
						<span
							class="flex h-12 w-12 items-center justify-center rounded-full bg-brand-100 text-base font-bold text-brand-800"
						>
							{item.seller?.name?.[0]?.toUpperCase() ?? '?'}
						</span>
					{/if}
					<p class="font-semibold text-neutral-900">{item.seller?.name ?? 'Vendedor'}</p>
				</div>

				<div class="mt-4 space-y-2">
					{#if item.whatsapp_phone}
						<a
							href={whatsappLink(item.whatsapp_phone, item.title)}
							target="_blank"
							rel="noopener noreferrer"
							class="flex h-12 w-full items-center justify-center gap-2 rounded-full bg-whatsapp text-sm font-bold text-white transition hover:brightness-95"
						>
							<svg class="h-5 w-5" viewBox="0 0 24 24" fill="currentColor" aria-hidden="true">
								<path
									d="M17.472 14.382c-.297-.149-1.758-.867-2.03-.967-.273-.099-.471-.148-.67.15-.197.297-.767.966-.94 1.164-.173.199-.347.223-.644.075-.297-.15-1.255-.463-2.39-1.475-.883-.788-1.48-1.761-1.653-2.059-.173-.297-.018-.458.13-.606.134-.133.298-.347.446-.52.149-.174.198-.298.298-.497.099-.198.05-.371-.025-.52-.075-.149-.669-1.612-.916-2.207-.242-.579-.487-.5-.669-.51-.173-.008-.371-.01-.57-.01-.198 0-.52.074-.792.372-.272.297-1.04 1.016-1.04 2.479 0 1.462 1.065 2.875 1.213 3.074.149.198 2.096 3.2 5.077 4.487.709.306 1.262.489 1.694.625.712.227 1.36.195 1.871.118.571-.085 1.758-.719 2.006-1.413.248-.694.248-1.289.173-1.413-.074-.124-.272-.198-.57-.347m-5.421 7.403h-.004a9.87 9.87 0 0 1-5.031-1.378l-.361-.214-3.741.982.998-3.648-.235-.374a9.86 9.86 0 0 1-1.51-5.26c.001-5.45 4.436-9.884 9.888-9.884 2.64 0 5.122 1.03 6.988 2.898a9.825 9.825 0 0 1 2.893 6.994c-.003 5.45-4.437 9.884-9.885 9.884m8.413-18.297A11.815 11.815 0 0 0 12.05 0C5.495 0 .16 5.335.157 11.892c0 2.096.547 4.142 1.588 5.945L.057 24l6.305-1.654a11.882 11.882 0 0 0 5.683 1.448h.005c6.554 0 11.89-5.335 11.893-11.893a11.821 11.821 0 0 0-3.48-8.413Z"
								/>
							</svg>
							WhatsApp
						</a>
					{/if}
					{#if item.phone ?? item.whatsapp_phone ?? item.seller?.phone}
						{@const contactPhone = (item.phone ?? item.whatsapp_phone ?? item.seller?.phone)!}
						<a
							href={`tel:${contactPhone}`}
							class="flex h-12 w-full items-center justify-center gap-2 rounded-full border border-brand-700 text-sm font-bold text-brand-800 transition hover:bg-brand-50"
						>
							<svg class="h-4 w-4" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
								<path
									fill-rule="evenodd"
									d="M2 3.5A1.5 1.5 0 0 1 3.5 2h1.148a1.5 1.5 0 0 1 1.465 1.175l.716 3.223a1.5 1.5 0 0 1-1.052 1.767l-.933.267c-.41.117-.643.555-.48.95a11.542 11.542 0 0 0 6.254 6.254c.395.163.833-.07.95-.48l.267-.933a1.5 1.5 0 0 1 1.767-1.052l3.223.716A1.5 1.5 0 0 1 18 15.352V16.5a1.5 1.5 0 0 1-1.5 1.5H15c-1.149 0-2.263-.15-3.326-.43A13.022 13.022 0 0 1 2.43 8.326 13.019 13.019 0 0 1 2 5V3.5Z"
									clip-rule="evenodd"
								/>
							</svg>
							{contactPhone}
						</a>
						<button
							type="button"
							onclick={() => copyNumber(contactPhone)}
							class="flex h-10 w-full items-center justify-center gap-2 rounded-full text-xs font-semibold text-neutral-500 transition hover:text-brand-700"
						>
							{#if copied}
								<svg class="h-4 w-4 text-brand-600" viewBox="0 0 20 20" fill="currentColor">
									<path
										fill-rule="evenodd"
										d="M16.704 4.153a.75.75 0 0 1 .143 1.052l-8 10.5a.75.75 0 0 1-1.127.075l-4.5-4.5a.75.75 0 0 1 1.06-1.06l3.894 3.893 7.48-9.817a.75.75 0 0 1 1.05-.143Z"
										clip-rule="evenodd"
									/>
								</svg>
								Número copiado!
							{:else}
								<svg class="h-4 w-4" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
									<path
										d="M7 3.5A1.5 1.5 0 0 1 8.5 2h3.879a1.5 1.5 0 0 1 1.06.44l3.122 3.12A1.5 1.5 0 0 1 17 6.622V12.5a1.5 1.5 0 0 1-1.5 1.5h-1v-3.379a3 3 0 0 0-.879-2.121L10.5 5.379A3 3 0 0 0 8.379 4.5H7v-1Z"
									/>
									<path
										d="M4.5 6A1.5 1.5 0 0 0 3 7.5v9A1.5 1.5 0 0 0 4.5 18h7a1.5 1.5 0 0 0 1.5-1.5v-5.879a1.5 1.5 0 0 0-.44-1.06L9.44 6.439A1.5 1.5 0 0 0 8.378 6H4.5Z"
									/>
								</svg>
								Copiar número
							{/if}
						</button>
					{/if}
				</div>

				<p class="mt-4 text-xs text-neutral-400">
					{item.views_count}
					{item.views_count === 1 ? 'visualização' : 'visualizações'}
				</p>
			</div>

			<div class="rounded-2xl bg-amber-50 p-4 text-xs leading-relaxed text-amber-800">
				<strong>Dica de segurança:</strong> encontra-te com o vendedor num local público e verifica o
				artigo antes de pagar.
			</div>

			<!-- painel admin: promoção do anúncio -->
			<svelte:boundary>
				{#if (await user)?.role === 'admin'}
					<div class="rounded-2xl border-2 border-dashed border-neutral-300 bg-neutral-100 p-5">
						<h2 class="flex items-center gap-1.5 text-sm font-bold text-neutral-700">
							<svg class="h-4 w-4" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
								<path
									fill-rule="evenodd"
									d="M8.34 1.804A1 1 0 0 1 9.32 1h1.36a1 1 0 0 1 .98.804l.295 1.473c.497.144.971.342 1.416.587l1.25-.834a1 1 0 0 1 1.262.125l.962.962a1 1 0 0 1 .125 1.262l-.834 1.25c.245.445.443.919.587 1.416l1.473.294a1 1 0 0 1 .804.98v1.361a1 1 0 0 1-.804.98l-1.473.295a6.95 6.95 0 0 1-.587 1.416l.834 1.25a1 1 0 0 1-.125 1.262l-.962.962a1 1 0 0 1-1.262.125l-1.25-.834a6.953 6.953 0 0 1-1.416.587l-.294 1.473a1 1 0 0 1-.98.804H9.32a1 1 0 0 1-.98-.804l-.295-1.473a6.957 6.957 0 0 1-1.416-.587l-1.25.834a1 1 0 0 1-1.262-.125l-.962-.962a1 1 0 0 1-.125-1.262l.834-1.25a6.957 6.957 0 0 1-.587-1.416l-1.473-.294A1 1 0 0 1 1 10.68V9.32a1 1 0 0 1 .804-.98l1.473-.295c.144-.497.342-.971.587-1.416l-.834-1.25a1 1 0 0 1 .125-1.262l.962-.962A1 1 0 0 1 5.38 3.03l1.25.834a6.957 6.957 0 0 1 1.416-.587l.294-1.473ZM13 10a3 3 0 1 1-6 0 3 3 0 0 1 6 0Z"
									clip-rule="evenodd"
								/>
							</svg>
							Administração
						</h2>

						{#if item.is_featured}
							<p class="mt-2 text-xs text-neutral-500">
								Em destaque até {item.featured_until
									? new Date(item.featured_until).toLocaleDateString('pt-PT')
									: '—'}
							</p>
						{/if}

						<div class="mt-3 flex flex-wrap gap-2">
							{#if item.is_featured}
								<button
									type="button"
									disabled={adminBusy}
									onclick={() => runAdminAction(() => unfeatureListing(item.id))}
									class="rounded-full border border-neutral-300 bg-white px-3.5 py-2 text-xs font-semibold text-neutral-700 transition hover:border-red-300 hover:text-red-600 disabled:opacity-50"
									>Remover destaque</button
								>
							{:else}
								{#each [7, 15, 30] as days (days)}
									<button
										type="button"
										disabled={adminBusy}
										onclick={() => runAdminAction(() => featureListing({ id: item.id, days }))}
										class="rounded-full bg-amber-400 px-3.5 py-2 text-xs font-bold text-amber-950 transition hover:bg-amber-300 disabled:opacity-50"
										>Destacar {days}d</button
									>
								{/each}
							{/if}
							<button
								type="button"
								disabled={adminBusy}
								onclick={() => runAdminAction(() => bumpListing(item.id))}
								class="rounded-full border border-neutral-300 bg-white px-3.5 py-2 text-xs font-semibold text-neutral-700 transition hover:border-brand-400 disabled:opacity-50"
								>Subir anúncio</button
							>
						</div>

						{#if adminError}
							<p class="mt-2 text-xs text-red-600">{adminError}</p>
						{/if}
					</div>
				{/if}

				{#snippet pending()}{/snippet}
				{#snippet failed()}{/snippet}
			</svelte:boundary>
		</aside>
	</div>

	{#snippet pending()}
		<div class="grid gap-6 lg:grid-cols-[1fr_360px]">
			<div class="aspect-4/3 animate-pulse rounded-2xl bg-neutral-200"></div>
			<div class="space-y-4">
				<div class="h-48 animate-pulse rounded-2xl bg-neutral-200"></div>
				<div class="h-40 animate-pulse rounded-2xl bg-neutral-100"></div>
			</div>
		</div>
	{/snippet}
</svelte:boundary>
