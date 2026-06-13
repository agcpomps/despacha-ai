# Deploy na Hetzner

Stack: Caddy (HTTPS automático) → SvelteKit (Node) + API Go → PostgreSQL.
Tudo em Docker Compose num único VPS.

## 1. Criar o servidor

1. [Hetzner Cloud](https://console.hetzner.cloud) → novo projecto → **Add Server**
   - Localização: **Falkenstein** ou **Nuremberg** (melhor latência para África via Europa)
   - Imagem: **Ubuntu 24.04**
   - Tipo: **CX22** (2 vCPU / 4GB — chega de sobra para o MVP)
   - SSH key: adiciona a tua chave pública
2. **Firewall** (na consola Hetzner): permitir apenas TCP 22, 80, 443

## 2. DNS (Cloudflare)

No painel da Cloudflare → o teu domínio → **DNS** → adicionar registos:

```
Tipo  Nome  Conteúdo            Proxy
A     @     <IP-do-servidor>    DNS only (nuvem CINZENTA)
A     www   <IP-do-servidor>    DNS only (nuvem cinzenta, opcional)
```

> ⚠️ **Importante: desliga o proxy (nuvem laranja → cinzenta) pelo menos no
> primeiro arranque.** Com o proxy ligado, o Caddy não consegue emitir o
> certificado Let's Encrypt e o site não sobe.
>
> Depois do site estar no ar com HTTPS, podes (opcional) ligar o proxy
> laranja para ganhares CDN/protecção DDoS da Cloudflare — nesse caso vai a
> **SSL/TLS → Overview** e escolhe o modo **Full (strict)**. Se não quiseres
> complicar, deixa em DNS only: o Caddy já te dá HTTPS na mesma.

## 3. Preparar o servidor (uma vez)

```bash
ssh root@<IP>

# Docker
curl -fsSL https://get.docker.com | sh

# Código
git clone https://github.com/agcpomps/despacha-ai.git /opt/despacha-ai
cd /opt/despacha-ai

# Configuração
cp .env.production.example .env
nano .env        # DOMAIN, passwords, JWT_SECRET (openssl rand -base64 48)
```

## 4. Arrancar

```bash
cd /opt/despacha-ai
docker compose -f docker-compose.prod.yml --env-file .env up -d --build
```

O Caddy obtém o certificado TLS automaticamente no primeiro arranque
(o DNS já tem de estar a apontar para o servidor).

Verificar:

```bash
docker compose -f docker-compose.prod.yml ps
curl https://<dominio>/health
```

## 5. Primeiro admin

```bash
docker exec -it despacha_ai_postgres psql -U despacha -d despacha_ai \
  -c "UPDATE users SET role = 'admin' WHERE phone = '+244XXXXXXXXX';"
```

(Regista-te primeiro no site; depois logout + login para o token apanhar o role.)

## 6. Actualizar (cada deploy)

```bash
cd /opt/despacha-ai
git pull
docker compose -f docker-compose.prod.yml --env-file .env up -d --build
```

## 7. Backups (importante!)

Duas camadas de protecção:

### Camada 1 — Hetzner Backups (recomendado, automático)

Ao criar o servidor, liga a opção **Backups** (~20% do custo do servidor, ~€0.76/mês
no CX22). São snapshots automáticos de todo o disco, guardados pela Hetzner **fora da
máquina** — protegem mesmo que o servidor morra. Zero manutenção.

### Camada 2 — backups locais granulares (script)

O [scripts/backup.sh](scripts/backup.sh) faz dump da BD + tar das imagens, com rotação
(guarda 7 dias). Permite restauro rápido e selectivo (ex: recuperar a BD depois de um
deploy mau). Testa uma vez à mão:

```bash
cd /opt/despacha-ai
./scripts/backup.sh
```

Depois agenda diariamente com o cron:

```bash
crontab -e
```

```cron
# 03h00 todos os dias
0 3 * * * cd /opt/despacha-ai && ./scripts/backup.sh >> /var/log/despacha-backup.log 2>&1
```

**Restauro** (comandos no topo do próprio script):

```bash
# base de dados
gunzip -c backups/db-XXXX.sql.gz | docker exec -i despacha_ai_postgres psql -U "$DB_USER" "$DB_NAME"
# imagens
docker exec -i despacha_ai_api sh -c 'tar xzf - -C /app/uploads' < backups/uploads-XXXX.tar.gz
```

> Para máxima segurança, podes enviar os dumps de `backups/` para um Hetzner
> Storage Box (~€3.20/mês, 1TB) com `rsync` — mas a Camada 1 já cobre o essencial.

## Notas

- A BD não expõe porta nenhuma para fora — só é acessível na rede interna do Docker.
- `PUBLIC_API_BASE_URL` aponta para o domínio público: as remote functions do
  SvelteKit chamam a API através do Caddy, e as URLs das imagens ficam corretas.
- Logs: `docker compose -f docker-compose.prod.yml logs -f api` (ou `frontend`, `caddy`).
