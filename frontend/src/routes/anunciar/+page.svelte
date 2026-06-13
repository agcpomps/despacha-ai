<script lang="ts">
	import { createListing } from '$lib/remote/listings.remote';
	import { getCategories } from '$lib/remote/categories.remote';
	import { PROVINCES } from '$lib/utils';

	const categories = getCategories();

	const MAX_IMAGES = 8;

	let photos = $state<File[]>([]);
	let previews = $state<string[]>([]);
	let submitInput: HTMLInputElement;

	// reconstrói a FileList do input que é submetido (FileList é só-leitura,
	// por isso usamos um DataTransfer) e refaz as pré-visualizações
	function sync() {
		const dt = new DataTransfer();
		for (const file of photos) dt.items.add(file);
		if (submitInput) submitInput.files = dt.files;

		for (const url of previews) URL.revokeObjectURL(url);
		previews = photos.map((file) => URL.createObjectURL(file));
	}

	// acumula as fotos vindas da câmara ou da galeria (até MAX_IMAGES)
	function addFiles(event: Event) {
		const input = event.currentTarget as HTMLInputElement;
		const incoming = Array.from(input.files ?? []).filter((f) => f.type.startsWith('image/'));
		photos = [...photos, ...incoming].slice(0, MAX_IMAGES);
		input.value = ''; // permite voltar a tirar foto / re-selecionar o mesmo ficheiro
		sync();
	}

	function removePhoto(index: number) {
		photos = photos.filter((_, i) => i !== index);
		sync();
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

			<!-- input realmente submetido; preenchido via DataTransfer em sync() -->
			<input
				bind:this={submitInput}
				type="file"
				name="images"
				accept="image/*"
				multiple
				class="hidden"
				tabindex="-1"
				aria-hidden="true"
			/>

			<div class="mt-3 grid grid-cols-2 gap-3">
				<!-- Tirar foto: capture abre a câmara direto no telemóvel -->
				<label
					class="flex cursor-pointer flex-col items-center justify-center gap-2 rounded-xl border-2 border-dashed border-neutral-300 p-6 text-center transition hover:border-brand-400 hover:bg-brand-50/40 aria-disabled:pointer-events-none aria-disabled:opacity-40"
					aria-disabled={previews.length >= MAX_IMAGES}
				>
					<svg class="h-7 w-7 text-brand-600" viewBox="0 0 24 24" fill="currentColor">
						<path
							fill-rule="evenodd"
							d="M12 16.5a4.5 4.5 0 1 0 0-9 4.5 4.5 0 0 0 0 9Z"
							clip-rule="evenodd"
						/>
						<path
							fill-rule="evenodd"
							d="M9.344 3.071a49.52 49.52 0 0 1 5.312 0c.967.052 1.83.585 2.332 1.39l.821 1.317c.24.383.645.643 1.11.71.386.054.77.113 1.152.177 1.432.239 2.429 1.493 2.429 2.909V18a3 3 0 0 1-3 3h-15a3 3 0 0 1-3-3V9.574c0-1.416.997-2.67 2.429-2.909.382-.064.766-.123 1.151-.178a1.56 1.56 0 0 0 1.11-.71l.822-1.315a2.942 2.942 0 0 1 2.332-1.39ZM6.75 12.75a5.25 5.25 0 1 1 10.5 0 5.25 5.25 0 0 1-10.5 0Zm12-1.5a.75.75 0 1 0 0-1.5.75.75 0 0 0 0 1.5Z"
							clip-rule="evenodd"
						/>
					</svg>
					<span class="text-sm font-medium text-neutral-600">Tirar foto</span>
					<input
						type="file"
						accept="image/*"
						capture="environment"
						class="hidden"
						onchange={addFiles}
					/>
				</label>

				<!-- Galeria: seleção múltipla -->
				<label
					class="flex cursor-pointer flex-col items-center justify-center gap-2 rounded-xl border-2 border-dashed border-neutral-300 p-6 text-center transition hover:border-brand-400 hover:bg-brand-50/40 aria-disabled:pointer-events-none aria-disabled:opacity-40"
					aria-disabled={previews.length >= MAX_IMAGES}
				>
					<svg class="h-7 w-7 text-brand-600" viewBox="0 0 24 24" fill="currentColor">
						<path
							fill-rule="evenodd"
							d="M1.5 6a3 3 0 0 1 3-3h15a3 3 0 0 1 3 3v12a3 3 0 0 1-3 3h-15a3 3 0 0 1-3-3V6Zm3-1.5A1.5 1.5 0 0 0 3 6v12c0 .39.149.745.393 1.011l6.647-6.646a2.25 2.25 0 0 1 3.182 0l1.949 1.948 1.586-1.585a2.25 2.25 0 0 1 3.182 0L21 13.789V6a1.5 1.5 0 0 0-1.5-1.5h-15Zm10.125 4.5a1.875 1.875 0 1 1 3.75 0 1.875 1.875 0 0 1-3.75 0Z"
							clip-rule="evenodd"
						/>
					</svg>
					<span class="text-sm font-medium text-neutral-600">Galeria</span>
					<input type="file" accept="image/*" multiple class="hidden" onchange={addFiles} />
				</label>
			</div>

			{#if previews.length}
				<div class="mt-4 grid grid-cols-4 gap-2">
					{#each previews as url, index (url)}
						<div
							class="group relative aspect-square overflow-hidden rounded-lg border border-neutral-200"
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
							<button
								type="button"
								onclick={() => removePhoto(index)}
								aria-label="Remover foto"
								class="absolute top-1 right-1 flex h-6 w-6 items-center justify-center rounded-full bg-neutral-900/70 text-white transition hover:bg-red-600"
							>
								<svg class="h-3.5 w-3.5" viewBox="0 0 20 20" fill="currentColor">
									<path
										d="M6.28 5.22a.75.75 0 0 0-1.06 1.06L8.94 10l-3.72 3.72a.75.75 0 1 0 1.06 1.06L10 11.06l3.72 3.72a.75.75 0 1 0 1.06-1.06L11.06 10l3.72-3.72a.75.75 0 0 0-1.06-1.06L10 8.94 6.28 5.22Z"
									/>
								</svg>
							</button>
						</div>
					{/each}
				</div>
				<p class="mt-2 text-xs text-neutral-400">{previews.length} de {MAX_IMAGES} fotos</p>
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
