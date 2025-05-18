CREATE EXTENSION IF NOT EXISTS pgcrypto;

DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1 FROM pg_catalog.pg_settings
    WHERE name = 'my.epoch_timestamp'
  ) THEN
    PERFORM set_config('my.epoch_timestamp', '1577836800', false);
  END IF;
END $$;

CREATE OR REPLACE FUNCTION new_id(prefix TEXT)
RETURNS TEXT AS $$
DECLARE
  t INTEGER;
  randb BYTEA;
  buf  BYTEA;
BEGIN
  t := (EXTRACT(EPOCH FROM now())::INT
        - current_setting('my.epoch_timestamp')::INT);
  buf := set_byte(
           set_byte(
             set_byte(
               set_byte(
                 '\000\000\000\000'::bytea,
                 0, (t >> 24) & 255
               ),
               1, (t >> 16) & 255
             ),
             2, (t >> 8)  & 255
           ),
           3, t & 255
         );

  randb := gen_random_bytes(8);
  buf := buf || randb;

  IF prefix = '' THEN
    RETURN encode(buf, 'hex');
  ELSE
    RETURN prefix || '_' || encode(buf, 'hex');
  END IF;
END;
$$ LANGUAGE plpgsql STRICT;

CREATE TABLE IF NOT EXISTS categories (
  id          VARCHAR(255) PRIMARY KEY DEFAULT new_id('categories'),
  name        VARCHAR(255) NOT NULL UNIQUE,
  icon        TEXT,
  created_at  TIMESTAMP   NOT NULL DEFAULT NOW(),
  updated_at  TIMESTAMP,
  deleted_at  TIMESTAMP
);

CREATE TABLE IF NOT EXISTS subcategories (
  id          VARCHAR(255) PRIMARY KEY DEFAULT new_id('subcategories'),
  name        VARCHAR(255) NOT NULL,
  category_id VARCHAR(255) NOT NULL
               REFERENCES categories(id),
  created_at  TIMESTAMP   NOT NULL DEFAULT NOW(),
  updated_at  TIMESTAMP,
  deleted_at  TIMESTAMP,
  UNIQUE (name, category_id)
);

INSERT INTO categories (name, icon) 
VALUES
  ('Serviços Domésticos e Reparos', '🔧'),
  ('Beleza e Estética', '💅'),
  ('Negócios e Consultoria', '💼'),
  ('Arte, Música e Eventos', '🎶'),
  ('Transporte e Veículos', '🚗'),
  ('Educação e Ensino', '📚'),
  ('Saúde e Bem-estar', '💊');

INSERT INTO subcategories (name, category_id)
VALUES
  ('Eletricistas', (SELECT id FROM categories WHERE name = 'Serviços Domésticos e Reparos')),
  ('Encanadores', (SELECT id FROM categories WHERE name = 'Serviços Domésticos e Reparos')),
  ('Pedreiros', (SELECT id FROM categories WHERE name = 'Serviços Domésticos e Reparos')),
  ('Pintores', (SELECT id FROM categories WHERE name = 'Serviços Domésticos e Reparos')),
  ('Jardineiros', (SELECT id FROM categories WHERE name = 'Serviços Domésticos e Reparos')),
  ('Marceneiros', (SELECT id FROM categories WHERE name = 'Serviços Domésticos e Reparos')),
  ('Técnicos de eletrodomésticos', (SELECT id FROM categories WHERE name = 'Serviços Domésticos e Reparos')),
  ('Limpadores de piscina', (SELECT id FROM categories WHERE name = 'Serviços Domésticos e Reparos')),
  ('Diaristas', (SELECT id FROM categories WHERE name = 'Serviços Domésticos e Reparos')),
  ('Babás', (SELECT id FROM categories WHERE name = 'Serviços Domésticos e Reparos')),
  ('Cuidadores de idosos', (SELECT id FROM categories WHERE name = 'Serviços Domésticos e Reparos')),
  ('Seguranças particulares', (SELECT id FROM categories WHERE name = 'Serviços Domésticos e Reparos')),

  ('Cabeleireiros', (SELECT id FROM categories WHERE name = 'Beleza e Estética')),
  ('Manicures', (SELECT id FROM categories WHERE name = 'Beleza e Estética')),
  ('Maquiadores', (SELECT id FROM categories WHERE name = 'Beleza e Estética')),
  ('Esteticistas', (SELECT id FROM categories WHERE name = 'Beleza e Estética')),
  ('Massagistas', (SELECT id FROM categories WHERE name = 'Beleza e Estética')),

  ('Advogados', (SELECT id FROM categories WHERE name = 'Negócios e Consultoria')),
  ('Contadores', (SELECT id FROM categories WHERE name = 'Negócios e Consultoria')),
  ('Consultores', (SELECT id FROM categories WHERE name = 'Negócios e Consultoria')),
  ('Tradutores', (SELECT id FROM categories WHERE name = 'Negócios e Consultoria')),
  ('Revisores de texto', (SELECT id FROM categories WHERE name = 'Negócios e Consultoria')),
  ('Copywriters', (SELECT id FROM categories WHERE name = 'Negócios e Consultoria')),
  ('Social Media', (SELECT id FROM categories WHERE name = 'Negócios e Consultoria')),

  ('Fotógrafos', (SELECT id FROM categories WHERE name = 'Arte, Música e Eventos')),
  ('Músicos', (SELECT id FROM categories WHERE name = 'Arte, Música e Eventos')),
  ('DJs', (SELECT id FROM categories WHERE name = 'Arte, Música e Eventos')),
  ('Cerimonialistas', (SELECT id FROM categories WHERE name = 'Arte, Música e Eventos')),
  ('Decoradores de eventos', (SELECT id FROM categories WHERE name = 'Arte, Música e Eventos')),
  ('Cozinheiros para eventos', (SELECT id FROM categories WHERE name = 'Arte, Música e Eventos')),
  ('Buffets', (SELECT id FROM categories WHERE name = 'Arte, Música e Eventos')),

  ('Motoristas particulares', (SELECT id FROM categories WHERE name = 'Transporte e Veículos')),
  ('Motoboys', (SELECT id FROM categories WHERE name = 'Transporte e Veículos')),
  ('Mecânicos', (SELECT id FROM categories WHERE name = 'Transporte e Veículos')),
  ('Lavadores de carro', (SELECT id FROM categories WHERE name = 'Transporte e Veículos')),
  ('Instrutores de direção', (SELECT id FROM categories WHERE name = 'Transporte e Veículos')),

  ('Professores', (SELECT id FROM categories WHERE name = 'Educação e Ensino')),
  ('Acompanhantes pedagógicos', (SELECT id FROM categories WHERE name = 'Educação e Ensino')),
  ('Monitores de reforço escolar', (SELECT id FROM categories WHERE name = 'Educação e Ensino')),
  ('Intérpretes de Libras', (SELECT id FROM categories WHERE name = 'Educação e Ensino')),

  ('Nutricionistas', (SELECT id FROM categories WHERE name = 'Saúde e Bem-estar')),
  ('Personal Trainers', (SELECT id FROM categories WHERE name = 'Saúde e Bem-estar'));
