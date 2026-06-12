<script lang="ts">
	import { page } from '$app/state';
	import { getListing, updateListing } from '$lib/remote/listings.remote';
	import { PROVINCES } from '$lib/utils';

	const listing = $derived(getListing(page.params.id!));
</script>

<svelte:head>
	<title>Editar anúncio | Despacha Aí</title>
</svelte:head>

<svelte:boundary>
	{@const item = await listing}

	<div class="mx-auto max-w-2xl">
		<a href="/dashboard" class="text-sm font-medium text-brand-700 hover:text-brand-800"
			>← Voltar aos meus anúncios</a
		>
		<h1 class="mt-2 text-2xl font-bold text-neutral-900">Editar anúncio</h1>

		<form
			{...updateListing}
			class="mt-6 space-y-4 rounded-2xl border border-neutral-200 bg-white p-6"
		>
			<input type="hidden" name="id" value={item.id} />

			<label class="block">
				<span class="text-sm font-semibold text-neutral-700">Título *</span>
				<input
					name="title"
					value={item.title}
					class="mt-1 h-12 w-full rounded-xl border-neutral-200 focus:border-brand-500 focus:ring-brand-500"
				/>
				{#each updateListing.fields.title.issues() ?? [] as issue (issue.message)}
					<p class="mt-1 text-sm text-red-600">{issue.message}</p>
				{/each}
			</label>

			<label class="block">
				<span class="text-sm font-semibold text-neutral-700">Descrição *</span>
				<textarea
					name="description"
					rows="5"
					class="mt-1 w-full rounded-xl border-neutral-200 focus:border-brand-500 focus:ring-brand-500"
					>{item.description}</textarea
				>
				{#each updateListing.fields.description.issues() ?? [] as issue (issue.message)}
					<p class="mt-1 text-sm text-red-600">{issue.message}</p>
				{/each}
			</label>

			<div class="grid gap-4 sm:grid-cols-2">
				<label class="block">
					<span class="text-sm font-semibold text-neutral-700">Preço (Kz) *</span>
					<input
						name="price"
						inputmode="decimal"
						value={item.price}
						class="mt-1 h-12 w-full rounded-xl border-neutral-200 focus:border-brand-500 focus:ring-brand-500"
					/>
					{#each updateListing.fields.price.issues() ?? [] as issue (issue.message)}
						<p class="mt-1 text-sm text-red-600">{issue.message}</p>
					{/each}
				</label>

				<label class="block">
					<span class="text-sm font-semibold text-neutral-700">Estado do artigo *</span>
					<select
						name="condition"
						class="mt-1 h-12 w-full rounded-xl border-neutral-200 focus:border-brand-500 focus:ring-brand-500"
					>
						<option value="used" selected={item.condition === 'used'}>Usado</option>
						<option value="new" selected={item.condition === 'new'}>Novo</option>
					</select>
				</label>
			</div>

			<div class="grid gap-4 sm:grid-cols-2">
				<label class="block">
					<span class="text-sm font-semibold text-neutral-700">Província *</span>
					<select
						name="province"
						class="mt-1 h-12 w-full rounded-xl border-neutral-200 focus:border-brand-500 focus:ring-brand-500"
					>
						{#each PROVINCES as province (province)}
							<option value={province} selected={item.province === province}>{province}</option>
						{/each}
					</select>
				</label>

				<label class="block">
					<span class="text-sm font-semibold text-neutral-700">Cidade / Município</span>
					<input
						name="city"
						value={item.city ?? ''}
						class="mt-1 h-12 w-full rounded-xl border-neutral-200 focus:border-brand-500 focus:ring-brand-500"
					/>
				</label>
			</div>

			<label class="block">
				<span class="text-sm font-semibold text-neutral-700">Ponto de referência</span>
				<input
					name="address_reference"
					value={item.address_reference ?? ''}
					class="mt-1 h-12 w-full rounded-xl border-neutral-200 focus:border-brand-500 focus:ring-brand-500"
				/>
			</label>

			<div class="grid gap-4 sm:grid-cols-2">
				<label class="block">
					<span class="text-sm font-semibold text-neutral-700">WhatsApp</span>
					<input
						name="whatsapp_phone"
						value={item.whatsapp_phone ?? ''}
						class="mt-1 h-12 w-full rounded-xl border-neutral-200 focus:border-brand-500 focus:ring-brand-500"
					/>
				</label>
				<label class="block">
					<span class="text-sm font-semibold text-neutral-700">Telefone</span>
					<input
						name="phone"
						value={item.phone ?? ''}
						class="mt-1 h-12 w-full rounded-xl border-neutral-200 focus:border-brand-500 focus:ring-brand-500"
					/>
				</label>
			</div>

			{#each updateListing.fields.allIssues() ?? [] as issue (issue.message)}
				{#if issue.path.length === 0}
					<p class="rounded-xl bg-red-50 p-3 text-sm text-red-700">{issue.message}</p>
				{/if}
			{/each}

			<div class="flex gap-3 pt-2">
				<button
					type="submit"
					disabled={updateListing.pending > 0}
					class="h-12 flex-1 rounded-full bg-brand-700 text-sm font-bold text-white transition hover:bg-brand-800 disabled:opacity-60"
				>
					{updateListing.pending > 0 ? 'A guardar…' : 'Guardar alterações'}
				</button>
				<a
					href="/dashboard"
					class="flex h-12 items-center rounded-full border border-neutral-200 px-6 text-sm font-medium text-neutral-600 hover:border-neutral-300"
					>Cancelar</a
				>
			</div>
		</form>
	</div>

	{#snippet pending()}
		<div class="mx-auto max-w-2xl">
			<div class="h-96 animate-pulse rounded-2xl bg-neutral-200"></div>
		</div>
	{/snippet}
</svelte:boundary>
