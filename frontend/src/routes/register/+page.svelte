<script lang="ts">
	import { page } from '$app/state';
	import { register } from '$lib/remote/auth.remote';

	const next = $derived(page.url.searchParams.get('next') ?? '/');
</script>

<svelte:head>
	<title>Criar conta | Despacha Aí</title>
</svelte:head>

<div class="mx-auto mt-6 max-w-md">
	<div class="rounded-2xl border border-neutral-200 bg-white p-8">
		<h1 class="text-2xl font-bold text-neutral-900">Criar conta</h1>
		<p class="mt-1 text-sm text-neutral-500">
			Regista-te gratuitamente e começa a vender em minutos.
		</p>

		<form {...register} class="mt-6 space-y-4">
			<input type="hidden" name="next" value={next} />

			<label class="block">
				<span class="text-sm font-semibold text-neutral-700">Nome completo</span>
				<input
					{...register.fields.name.as('text')}
					placeholder="O teu nome"
					autocomplete="name"
					class="mt-1 h-12 w-full rounded-xl border-neutral-200 focus:border-brand-500 focus:ring-brand-500"
				/>
				{#each register.fields.name.issues() ?? [] as issue (issue.message)}
					<p class="mt-1 text-sm text-red-600">{issue.message}</p>
				{/each}
			</label>

			<label class="block">
				<span class="text-sm font-semibold text-neutral-700">Telemóvel</span>
				<input
					{...register.fields.phone.as('tel')}
					placeholder="923 456 789"
					autocomplete="tel"
					class="mt-1 h-12 w-full rounded-xl border-neutral-200 focus:border-brand-500 focus:ring-brand-500"
				/>
				{#each register.fields.phone.issues() ?? [] as issue (issue.message)}
					<p class="mt-1 text-sm text-red-600">{issue.message}</p>
				{/each}
			</label>

			<label class="block">
				<span class="text-sm font-semibold text-neutral-700">
					Email <span class="font-normal text-neutral-400">(opcional)</span>
				</span>
				<input
					{...register.fields.email.as('email')}
					placeholder="email@exemplo.com"
					autocomplete="email"
					class="mt-1 h-12 w-full rounded-xl border-neutral-200 focus:border-brand-500 focus:ring-brand-500"
				/>
				{#each register.fields.email.issues() ?? [] as issue (issue.message)}
					<p class="mt-1 text-sm text-red-600">{issue.message}</p>
				{/each}
			</label>

			<label class="block">
				<span class="text-sm font-semibold text-neutral-700">Palavra-passe</span>
				<input
					{...register.fields.password.as('password')}
					autocomplete="new-password"
					class="mt-1 h-12 w-full rounded-xl border-neutral-200 focus:border-brand-500 focus:ring-brand-500"
				/>
				{#each register.fields.password.issues() ?? [] as issue (issue.message)}
					<p class="mt-1 text-sm text-red-600">{issue.message}</p>
				{/each}
			</label>

			{#each register.fields.allIssues() ?? [] as issue (issue.message)}
				{#if issue.path.length === 0}
					<p class="rounded-xl bg-red-50 p-3 text-sm text-red-700">{issue.message}</p>
				{/if}
			{/each}

			<button
				type="submit"
				disabled={register.pending > 0}
				class="h-12 w-full rounded-full bg-brand-700 text-sm font-bold text-white transition hover:bg-brand-800 disabled:opacity-60"
			>
				{register.pending > 0 ? 'A criar conta…' : 'Criar conta'}
			</button>
		</form>

		<p class="mt-6 text-center text-sm text-neutral-500">
			Já tens conta?
			<a
				href={`/login${next !== '/' ? `?next=${encodeURIComponent(next)}` : ''}`}
				class="font-semibold text-brand-700 hover:text-brand-800">Entrar</a
			>
		</p>
	</div>
</div>
