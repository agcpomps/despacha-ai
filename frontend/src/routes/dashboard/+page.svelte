<script lang="ts">
	import { page } from '$app/state';
	import { env } from '$env/dynamic/public';
	import { getMyListings, setListingStatus, deleteListing } from '$lib/remote/listings.remote';
	import { formatPrice, timeAgo, STATUS_LABELS } from '$lib/utils';
	import ConfirmDialog from '$lib/components/ConfirmDialog.svelte';
	import type { Listing } from '$lib/types';

	const adminWhatsApp = (env.PUBLIC_ADMIN_WHATSAPP ?? '').replace(/\D/g, '');

	function promoteLink(listing: Listing) {
		const text = encodeURIComponent(
			`Olá! Quero destacar o meu anúncio "${listing.title}" no Despacha Aí:\n${page.url.origin}/anuncio/${listing.id}\n\nQuais são os preços?`
		);
		return `https://wa.me/${adminWhatsApp}?text=${text}`;
	}
	const statusFilter = $derived.by(() => {
		const value = page.url.searchParams.get('status');
		return value === 'active' || value === 'sold' || value === 'paused' ? value : undefined;
	});

	const pageNumber = $derived(Number(page.url.searchParams.get('page')) || 1);

	const myListings = $derived(getMyListings({ status: statusFilter, page: pageNumber, limit: 20 }));

	let busyId = $state<string | null>(null);
	let actionError = $state<string | null>(null);

	// confirmação de remoção
	let confirmOpen = $state(false);
	let pendingDeleteId = $state<string | null>(null);

	async function changeStatus(id: string, status: 'active' | 'sold' | 'paused') {
		busyId = id;
		actionError = null;
		try {
			await setListingStatus({ id, status });
			await myListings.refresh();
		} catch {
			actionError = 'Não foi possível actualizar o anúncio.';
		} finally {
			busyId = null;
		}
	}

	function askRemove(id: string) {
		pendingDeleteId = id;
		confirmOpen = true;
	}

	async function confirmRemove() {
		const id = pendingDeleteId;
		pendingDeleteId = null;
		if (!id) return;

		busyId = id;
		actionError = null;
		try {
			await deleteListing(id);
			await myListings.refresh();
		} catch {
			actionError = 'Não foi possível apagar o anúncio.';
		} finally {
			busyId = null;
		}
	}

	const tabs = [
		{ label: 'Todos', value: undefined },
		{ label: 'Activos', value: 'active' },
		{ label: 'Vendidos', value: 'sold' },
		{ label: 'Pausados', value: 'paused' }
	] as const;

	const statusStyles: Record<string, string> = {
		active: 'bg-brand-100 text-brand-800',
		sold: 'bg-neutral-200 text-neutral-700',
		paused: 'bg-amber-100 text-amber-800'
	};
</script>

<svelte:head>
	<title>Os meus anúncios | Despacha Aí</title>
</svelte:head>

<div class="flex flex-wrap items-center justify-between gap-4">
	<div>
		<h1 class="text-2xl font-bold text-neutral-900">Os meus anúncios</h1>
		<p class="mt-1 text-sm text-neutral-500">Gere, edita e acompanha as tuas publicações.</p>
	</div>
	<a
		href="/anunciar"
		class="rounded-full bg-brand-700 px-5 py-2.5 text-sm font-semibold text-white transition hover:bg-brand-800"
		>Novo anúncio</a
	>
</div>

<nav class="mt-6 flex gap-2 overflow-x-auto" aria-label="Filtrar por estado">
	{#each tabs as tab (tab.label)}
		<a
			href={tab.value ? `/dashboard?status=${tab.value}` : '/dashboard'}
			class="shrink-0 rounded-full px-4 py-2 text-sm font-medium transition {statusFilter ===
			tab.value
				? 'bg-brand-700 text-white'
				: 'border border-neutral-200 bg-white text-neutral-600 hover:border-brand-300'}"
		>
			{tab.label}
		</a>
	{/each}
</nav>

{#if actionError}
	<p class="mt-4 rounded-xl bg-red-50 p-3 text-sm text-red-700">{actionError}</p>
{/if}

<svelte:boundary>
	{@const result = await myListings}
	{#if result.data.length}
		<div class="mt-6 space-y-3">
			{#each result.data as listing (listing.id)}
				<article
					class="flex flex-col gap-4 rounded-2xl border border-neutral-200 bg-white p-4 sm:flex-row sm:items-center"
				>
					<a href={`/anuncio/${listing.id}`} class="flex min-w-0 flex-1 items-center gap-4">
						<img
							src={listing.images?.[0]?.image_url ?? '/product-placeholder.svg'}
							alt=""
							class="h-20 w-20 shrink-0 rounded-xl border border-neutral-100 object-cover"
						/>
						<div class="min-w-0">
							<h2 class="truncate font-semibold text-neutral-900">{listing.title}</h2>
							<p class="mt-0.5 text-sm font-bold text-brand-800">{formatPrice(listing.price)}</p>
							<p class="mt-0.5 flex items-center gap-2 text-xs text-neutral-500">
								<span
									class="rounded-full px-2 py-0.5 font-semibold {statusStyles[listing.status] ??
										'bg-neutral-100 text-neutral-600'}"
								>
									{STATUS_LABELS[listing.status] ?? listing.status}
								</span>
								<span>{listing.views_count} visualizações</span>
								<span>· {timeAgo(listing.created_at)}</span>
							</p>
						</div>
					</a>

					<div class="flex shrink-0 flex-wrap items-center gap-2">
						{#if listing.status === 'active' && !listing.is_featured && adminWhatsApp}
							<a
								href={promoteLink(listing)}
								target="_blank"
								rel="noopener noreferrer"
								class="flex items-center gap-1 rounded-full bg-amber-400 px-3.5 py-2 text-xs font-bold text-amber-950 transition hover:bg-amber-300 active:scale-95"
							>
								<svg class="h-3.5 w-3.5" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
									<path
										fill-rule="evenodd"
										d="M10.868 2.884c-.321-.772-1.415-.772-1.736 0l-1.83 4.401-4.753.381c-.833.067-1.171 1.107-.536 1.651l3.62 3.102-1.106 4.637c-.194.813.691 1.456 1.405 1.02L10 15.591l4.069 2.485c.713.436 1.598-.207 1.404-1.02l-1.106-4.637 3.62-3.102c.635-.544.297-1.584-.536-1.65l-4.752-.382-1.831-4.401Z"
										clip-rule="evenodd"
									/>
								</svg>
								Destacar
							</a>
						{/if}
						{#if listing.status === 'active'}
							<button
								type="button"
								onclick={() => changeStatus(listing.id, 'sold')}
								disabled={busyId === listing.id}
								class="rounded-full border border-neutral-200 px-3.5 py-2 text-xs font-semibold text-neutral-700 transition hover:border-brand-400 disabled:opacity-50"
								>Marcar vendido</button
							>
							<button
								type="button"
								onclick={() => changeStatus(listing.id, 'paused')}
								disabled={busyId === listing.id}
								class="rounded-full border border-neutral-200 px-3.5 py-2 text-xs font-semibold text-neutral-700 transition hover:border-amber-400 disabled:opacity-50"
								>Pausar</button
							>
						{:else}
							<button
								type="button"
								onclick={() => changeStatus(listing.id, 'active')}
								disabled={busyId === listing.id}
								class="rounded-full border border-brand-300 px-3.5 py-2 text-xs font-semibold text-brand-800 transition hover:bg-brand-50 disabled:opacity-50"
								>Reactivar</button
							>
						{/if}
						<a
							href={`/dashboard/editar/${listing.id}`}
							class="rounded-full border border-neutral-200 px-3.5 py-2 text-xs font-semibold text-neutral-700 transition hover:border-brand-400"
							>Editar</a
						>
						<button
							type="button"
							onclick={() => askRemove(listing.id)}
							disabled={busyId === listing.id}
							class="rounded-full border border-red-200 px-3.5 py-2 text-xs font-semibold text-red-600 transition hover:bg-red-50 disabled:opacity-50"
							>Apagar</button
						>
					</div>
				</article>
			{/each}
		</div>
	{:else}
		<div class="mt-6 rounded-2xl border border-dashed border-neutral-300 bg-white p-12 text-center">
			<p class="text-lg font-semibold text-neutral-800">
				{statusFilter ? 'Nenhum anúncio neste estado.' : 'Ainda não publicaste nenhum anúncio.'}
			</p>
			<a
				href="/anunciar"
				class="mt-4 inline-block rounded-full bg-brand-700 px-6 py-2.5 text-sm font-semibold text-white hover:bg-brand-800"
				>Publicar o primeiro anúncio</a
			>
		</div>
	{/if}

	{#snippet pending()}
		<div class="mt-6 space-y-3">
			{#each { length: 4 }, i (i)}
				<div class="h-28 animate-pulse rounded-2xl bg-neutral-200"></div>
			{/each}
		</div>
	{/snippet}
</svelte:boundary>

<ConfirmDialog
	bind:open={confirmOpen}
	danger
	title="Apagar anúncio"
	message="Tens a certeza que queres apagar este anúncio? Esta ação não pode ser desfeita."
	confirmLabel="Apagar"
	onconfirm={confirmRemove}
/>
