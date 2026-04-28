INSERT INTO categories (name, slug, parent_id)
VALUES
-- Imóveis
('Imóveis', 'imoveis', NULL),
('Casas à venda', 'casas-a-venda', (SELECT id FROM categories WHERE slug = 'imoveis')),
('Casas para alugar', 'casas-para-alugar', (SELECT id FROM categories WHERE slug = 'imoveis')),
('Apartamentos à venda', 'apartamentos-a-venda', (SELECT id FROM categories WHERE slug = 'imoveis')),
('Apartamentos para alugar', 'apartamentos-para-alugar', (SELECT id FROM categories WHERE slug = 'imoveis')),
('Terrenos', 'terrenos', (SELECT id FROM categories WHERE slug = 'imoveis')),

-- Viaturas
('Viaturas', 'viaturas', NULL),
('Carros usados', 'carros-usados', (SELECT id FROM categories WHERE slug = 'viaturas')),
('Carros novos', 'carros-novos', (SELECT id FROM categories WHERE slug = 'viaturas')),
('Motas', 'motas', (SELECT id FROM categories WHERE slug = 'viaturas')),
('Peças e acessórios', 'pecas-e-acessorios', (SELECT id FROM categories WHERE slug = 'viaturas')),

-- Electrónicos
('Electrónicos', 'electronicos', NULL),
('Telemóveis', 'telemoveis', (SELECT id FROM categories WHERE slug = 'electronicos')),
('Computadores', 'computadores', (SELECT id FROM categories WHERE slug = 'electronicos')),
('TV e Áudio', 'tv-audio', (SELECT id FROM categories WHERE slug = 'electronicos')),

-- Empregos e Serviços
('Empregos', 'empregos', NULL),
('Serviços', 'servicos', NULL),

-- Casa
('Casa e Jardim', 'casa-e-jardim', NULL),
('Móveis', 'moveis', (SELECT id FROM categories WHERE slug = 'casa-e-jardim')),
('Electrodomésticos', 'electrodomesticos', (SELECT id FROM categories WHERE slug = 'casa-e-jardim')),

-- Moda
('Moda', 'moda', NULL),
('Roupa', 'roupa', (SELECT id FROM categories WHERE slug = 'moda')),
('Calçado', 'calcado', (SELECT id FROM categories WHERE slug = 'moda')),

-- Diversos
('Diversos', 'diversos', NULL)
ON CONFLICT (slug) DO NOTHING;