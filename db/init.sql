-- Create the 'postos' table
CREATE TABLE IF NOT EXISTS postos
(
    id
    SERIAL
    PRIMARY
    KEY,
    razao_social
    TEXT
    NOT
    NULL,
    cnpj
    TEXT
    NOT
    NULL
    UNIQUE,
    distribuidora
    TEXT,
    endereco
    TEXT,
    complemento
    TEXT,
    bairro
    TEXT,
    latitude
    DOUBLE
    PRECISION,
    longitude
    DOUBLE
    PRECISION
);

-- Create the 'amostras' table
CREATE TABLE IF NOT EXISTS amostras
(
    id
    SERIAL
    PRIMARY
    KEY,
    id_numeric
    INT
    UNIQUE
    NOT
    NULL,
    data_coleta
    DATE
    NOT
    NULL,
    grupo_produto
    TEXT
    NOT
    NULL,
    produto
    TEXT
    NOT
    NULL,
    posto_id
    INT
    NOT
    NULL
    REFERENCES
    postos
(
    id
),
    municipio TEXT NOT NULL,
    estado TEXT NOT NULL
    );

-- Create the 'ensaios' table
CREATE TABLE IF NOT EXISTS ensaios
(
    id
    SERIAL
    PRIMARY
    KEY,
    amostra_id
    INT
    NOT
    NULL
    REFERENCES
    amostras
(
    id
),
    nome TEXT NOT NULL,
    resultado TEXT,
    unidade TEXT,
    conforme BOOLEAN
    );
-- Create table to track imported files
CREATE TABLE IF NOT EXISTS imported_files (
                                              id SERIAL PRIMARY KEY,
                                              filename TEXT UNIQUE NOT NULL,
                                              year INT NOT NULL,
                                              month INT NOT NULL,
                                              imported_at TIMESTAMP DEFAULT NOW()
    );


-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_amostras_data_coleta ON amostras (data_coleta);
CREATE INDEX IF NOT EXISTS idx_postos_cnpj ON postos (cnpj);
CREATE INDEX IF NOT EXISTS idx_ensaios_amostra ON ensaios (amostra_id);
CREATE INDEX IF NOT EXISTS idx_imported_files_date ON imported_files (year, month);
