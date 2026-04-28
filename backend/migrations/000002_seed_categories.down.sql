DELETE FROM categories
WHERE slug IN (
  'imoveis',
  'casas-a-venda',
  'casas-para-alugar',
  'apartamentos-a-venda',
  'apartamentos-para-alugar',
  'terrenos',

  'viaturas',
  'carros-usados',
  'carros-novos',
  'motas',
  'pecas-e-acessorios',

  'electronicos',
  'telemoveis',
  'computadores',
  'tv-audio',

  'empregos',
  'servicos',

  'casa-e-jardim',
  'moveis',
  'electrodomesticos',

  'moda',
  'roupa',
  'calcado',

  'diversos'
);