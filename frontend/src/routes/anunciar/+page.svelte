<script lang="ts">
	import { createListing } from '$lib/remote/listings.remote';
	import { getCategories } from '$lib/remote/categories.remote';
	import { PROVINCES } from '$lib/utils';

	const categories = getCategories();

	let previews = $state<string[]>([]);

	function handleFilesChange(event: Event) {
		const input = event.currentTarget as HTMLInputElement;
		for (const url of previews) URL.revokeObjectURL(url);
		previews = Array.from(input.files ?? [])
			.slice(0, 8)
			.map((file) => URL.createObjectURL(file));
	}
</script>

<svelte:head>
	<title>Anunciar | Despacha Aí</title>
</svelte:head>

<div class="mx-auto max-w-2xl">
	<h1 class="text-2xl font-bold text-neutral-900">Publicar anúncio</h1>
	<p class="mt-1 text-sm text-neutral-500">
		Preenche os dados do artigo. Os campos marcados com * são obrigatórios.
	</p>

	<form {...createListing} enctype="multipart/form-data" class="mt-6 space-y-6">
		<!-- Article info -->
		<section class="space-y-4 rounded-2xl border border-neutral-200 bg-white p-6">
			<h2 class="text-sm font-bold tracking-wide text-neutral-500 uppercase">Artigo</h2>

			<label class="block">
				<span class="text-sm font-semibold text-neutral-700">Título *</span>
				<input
					{...createListing.fields.title.as('text')}
					placeholder="Ex: iPhone 13 Pro 128GB"
					class="mt-1 h-12 w-full rounded-xl border-neutral-200 focus:border-brand-500 focus:ring-brand-500"
				/>
				{#each createListing.fields.title.issues() ?? [] as issue (issue.message)}
					<p class="mt-1 text-sm text-red-600">{issue.message}</p>
				{/each}
			</label>

			<label class="block">
				<span class="text-sm font-semibold text-neutral-700">Descrição *</span>
				<textarea
					name="description"
					rows="5"
					placeholder="Descreve o estado, características e detalhes do artigo"
					class="mt-1 w-full rounded-xl border-neutral-200 focus:border-brand-500 focus:ring-brand-500"
				></textarea>
				{#each createListing.fields.description.issues() ?? [] as issue (issue.message)}
					<p class="mt-1 text-sm text-red-600">{issue.message}</p>
				{/each}
			</label>

			<div class="grid gap-4 sm:grid-cols-2">
				<label class="block">
					<span class="text-sm font-semibold text-neutral-700">Preço (Kz) *</span>
					<input
						{...createListing.fields.price.as('text')}
						inputmode="decimal"
						placeholder="50 000"
						class="mt-1 h-12 w-full rounded-xl border-neutral-200 focus:border-brand-500 focus:ring-brand-500"
					/>
					{#each createListing.fields.price.issues() ?? [] as issue (issue.message)}
						<p class="mt-1 text-sm text-red-600">{issue.message}</p>
					{/each}
				</label>

				<label class="block">
					<span class="text-sm font-semibold text-neutral-700">Categoria</span>
					<select
						name="category_id"
						class="mt-1 h-12 w-full rounded-xl border-neutral-200 focus:border-brand-500 focus:ring-brand-500"
					>
						<option value="">Sem categoria</option>
						<svelte:boundary>
							{#each await categories as category (category.id)}
								<option value={category.id}>{category.name}</option>
								{#each category.children ?? [] as child (child.id)}
									<option value={child.id}>&nbsp;&nbsp;{child.name}</option>
								{/each}
							{/each}
							{#snippet pending()}{/snippet}
							{#snippet failed()}{/snippet}
						</svelte:boundary>
					</select>
				</label>
			</div>

			<fieldset>
				<legend class="text-sm font-semibold text-neutral-700">Estado do artigo *</legend>
				<div class="mt-2 flex gap-3">
					<label
						class="flex flex-1 cursor-pointer items-center justify-center gap-2 rounded-xl border border-neutral-200 px-4 py-3 text-sm font-medium has-checked:border-brand-600 has-checked:bg-brand-50 has-checked:text-brand-800"
					>
						<input type="radio" name="condition" value="used" checked class="text-brand-600" />
						Usado
					</label>
					<label
						class="flex flex-1 cursor-pointer items-center justify-center gap-2 rounded-xl border border-neutral-200 px-4 py-3 text-sm font-medium has-checked:border-brand-600 has-checked:bg-brand-50 has-checked:text-brand-800"
					>
						<input type="radio" name="condition" value="new" class="text-brand-600" />
						Novo
					</label>
				</div>
				{#each createListing.fields.condition.issues() ?? [] as issue (issue.message)}
					<p class="mt-1 text-sm text-red-600">{issue.message}</p>
				{/each}
			</fieldset>
		</section>

		<!-- Images -->
		<section class="rounded-2xl border border-neutral-200 bg-white p-6">
			<h2 class="text-sm font-bold tracking-wide text-neutral-500 uppercase">Fotografias</h2>
			<p class="mt-1 text-xs text-neutral-500">Até 8 imagens. A primeira será a capa do anúncio.</p>

			<label
				class="mt-3 flex cursor-pointer flex-col items-center justify-center gap-2 rounded-xl border-2 border-dashed border-neutral-300 p-8 text-center transition hover:border-brand-400 hover:bg-brand-50/40"
			>
				<svg class="h-8 w-8 text-neutral-400" viewBox="0 0 24 24" fill="currentColor">
					<path
						fill-rule="evenodd"
						d="M1.5 6a3 3 0 0 1 3-3h15a3 3 0 0 1 3 3v12a3 3 0 0 1-3 3h-15a3 3 0 0 1-3-3V6Zm3-1.5A1.5 1.5 0 0 0 3 6v12c0 .39.149.745.393 1.011l6.647-6.646a2.25 2.25 0 0 1 3.182 0l1.949 1.948 1.586-1.585a2.25 2.25 0 0 1 3.182 0L21 13.789V6a1.5 1.5 0 0 0-1.5-1.5h-15Zm10.125 4.5a1.875 1.875 0 1 1 3.75 0 1.875 1.875 0 0 1-3.75 0Z"
						clip-rule="evenodd"
					/>
				</svg>
				<span class="text-sm font-medium text-neutral-600">Clica para escolher imagens</span>
				<input
					type="file"
					name="images"
					accept="image/*"
					multiple
					class="hidden"
					onchange={handleFilesChange}
				/>
			</label>

			{#if previews.length}
				<div class="mt-4 grid grid-cols-4 gap-2">
					{#each previews as url, index (url)}
						<div
							class="relative aspect-square overflow-hidden rounded-lg border border-neutral-200"
						>
							<img
								src={url}
								alt={`Pré-visualização ${index + 1}`}
								class="h-full w-full object-cover"
							/>
							{#if index === 0}
								<span
									class="absolute bottom-1 left-1 rounded bg-brand-700 px-1.5 py-0.5 text-[10px] font-semibold text-white"
									>Capa</span
								>
							{/if}
						</div>
					{/each}
				</div>
			{/if}
			{#each createListing.fields.images.issues() ?? [] as issue (issue.message)}
				<p class="mt-2 text-sm text-red-600">{issue.message}</p>
			{/each}
		</section>

		<!-- Location & contact -->
		<section class="space-y-4 rounded-2xl border border-neutral-200 bg-white p-6">
			<h2 class="text-sm font-bold tracking-wide text-neutral-500 uppercase">
				Localização e contacto
			</h2>

			<div class="grid gap-4 sm:grid-cols-2">
				<label class="block">
					<span class="text-sm font-semibold text-neutral-700">Província *</span>
					<select
						name="province"
						class="mt-1 h-12 w-full rounded-xl border-neutral-200 focus:border-brand-500 focus:ring-brand-500"
					>
						<option value="">Seleccionar</option>
						{#each PROVINCES as province (province)}
							<option value={province}>{province}</option>
						{/each}
					</select>
					{#each createListing.fields.province.issues() ?? [] as issue (issue.message)}
						<p class="mt-1 text-sm text-red-600">{issue.message}</p>
					{/each}
				</label>

				<label class="block">
					<span class="text-sm font-semibold text-neutral-700">Cidade / Município</span>
					<input
						{...createListing.fields.city.as('text')}
						placeholder="Ex: Viana"
						class="mt-1 h-12 w-full rounded-xl border-neutral-200 focus:border-brand-500 focus:ring-brand-500"
					/>
				</label>
			</div>

			<label class="block">
				<span class="text-sm font-semibold text-neutral-700">Ponto de referência</span>
				<input
					{...createListing.fields.address_reference.as('text')}
					placeholder="Ex: Junto ao mercado"
					class="mt-1 h-12 w-full rounded-xl border-neutral-200 focus:border-brand-500 focus:ring-brand-500"
				/>
			</label>

			<div class="grid gap-4 sm:grid-cols-2">
				<label class="block">
					<span class="text-sm font-semibold text-neutral-700">WhatsApp</span>
					<input
						{...createListing.fields.whatsapp_phone.as('tel')}
						placeholder="+244 923 456 789"
						class="mt-1 h-12 w-full rounded-xl border-neutral-200 focus:border-brand-500 focus:ring-brand-500"
					/>
				</label>
				<label class="block">
					<span class="text-sm font-semibold text-neutral-700">Telefone</span>
					<input
						{...createListing.fields.phone.as('tel')}
						placeholder="+244 923 456 789"
						class="mt-1 h-12 w-full rounded-xl border-neutral-200 focus:border-brand-500 focus:ring-brand-500"
					/>
				</label>
			</div>
		</section>

		{#each createListing.fields.allIssues() ?? [] as issue (issue.message)}
			{#if issue.path.length === 0}
				<p class="rounded-xl bg-red-50 p-3 text-sm text-red-700">{issue.message}</p>
			{/if}
		{/each}

		<button
			type="submit"
			disabled={createListing.pending > 0}
			class="h-13 w-full rounded-full bg-brand-700 py-3.5 text-sm font-bold text-white transition hover:bg-brand-800 disabled:opacity-60"
		>
			{createListing.pending > 0 ? 'A publicar…' : 'Publicar anúncio'}
		</button>
	</form>
</div>
