#!/usr/bin/env bash
#
# Backup local rotativo da base de dados + imagens dos anúncios.
#
# Uso manual:   ./scripts/backup.sh
# Via cron:     ver DEPLOY.md (secção Backups)
#
# Restauro:
#   BD:       gunzip -c backups/db-XXXX.sql.gz | docker exec -i despacha_ai_postgres psql -U "$DB_USER" "$DB_NAME"
#   Imagens:  docker exec -i despacha_ai_api sh -c 'tar xzf - -C /app/uploads' < backups/uploads-XXXX.tar.gz

set -euo pipefail

# corre a partir da raiz do projecto, seja qual for o cwd
cd "$(dirname "$0")/.."

# carrega DB_USER / DB_NAME do .env
if [[ ! -f .env ]]; then
	echo "✗ .env não encontrado em $(pwd)" >&2
	exit 1
fi
set -a
# shellcheck disable=SC1091
source .env
set +a

BACKUP_DIR="${BACKUP_DIR:-$(pwd)/backups}"
KEEP_DAYS="${KEEP_DAYS:-7}"
STAMP="$(date +%F-%H%M)"

mkdir -p "$BACKUP_DIR"

echo "→ Backup da base de dados..."
docker exec despacha_ai_postgres pg_dump -U "$DB_USER" "$DB_NAME" \
	| gzip >"$BACKUP_DIR/db-$STAMP.sql.gz"

echo "→ Backup das imagens..."
# o tar corre dentro do container da API (alpine/busybox já traz tar),
# por isso não dependemos do nome prefixado do volume Docker
docker exec despacha_ai_api tar czf - -C /app/uploads . \
	>"$BACKUP_DIR/uploads-$STAMP.tar.gz"

echo "→ A remover backups com mais de $KEEP_DAYS dias..."
find "$BACKUP_DIR" -name 'db-*.sql.gz' -mtime +"$KEEP_DAYS" -delete
find "$BACKUP_DIR" -name 'uploads-*.tar.gz' -mtime +"$KEEP_DAYS" -delete

echo "✓ Backup concluído em $BACKUP_DIR"
ls -lh "$BACKUP_DIR" | tail -n +2
