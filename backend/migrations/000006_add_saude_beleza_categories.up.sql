-- NOTA: INSERTs separados de propósito — um subselect dentro do mesmo
-- INSERT ... VALUES não vê as linhas inseridas por esse próprio statement,
-- pelo que o parent_id resolvia para NULL (bug latente no seed 000002).
INSERT INTO categories (name, slug, parent_id)
VALUES ('Saúde e Beleza', 'saude-e-beleza', NULL)
ON CONFLICT (slug) DO NOTHING;

INSERT INTO categories (name, slug, parent_id)
VALUES ('Suplementos', 'suplementos', (SELECT id FROM categories WHERE slug = 'saude-e-beleza'))
ON CONFLICT (slug) DO NOTHING;

-- Reparação: bases de dados criadas de raiz com o seed 000002 ficaram com a
-- árvore plana (subcategorias sem parent_id), pelo mesmo motivo acima.
-- Repõe o parent_id de todas as subcategorias conhecidas. Idempotente:
-- só toca em linhas cujo parent_id está NULL.
UPDATE categories c
SET parent_id = p.id
FROM categories p
WHERE c.parent_id IS NULL
  AND (p.slug, c.slug) IN (
    ('imoveis', 'casas-a-venda'),
    ('imoveis', 'casas-para-alugar'),
    ('imoveis', 'apartamentos-a-venda'),
    ('imoveis', 'apartamentos-para-alugar'),
    ('imoveis', 'terrenos'),
    ('viaturas', 'carros-usados'),
    ('viaturas', 'carros-novos'),
    ('viaturas', 'motas'),
    ('viaturas', 'pecas-e-acessorios'),
    ('electronicos', 'telemoveis'),
    ('electronicos', 'computadores'),
    ('electronicos', 'tv-audio'),
    ('casa-e-jardim', 'moveis'),
    ('casa-e-jardim', 'electrodomesticos'),
    ('moda', 'roupa'),
    ('moda', 'calcado'),
    ('saude-e-beleza', 'suplementos')
  );
